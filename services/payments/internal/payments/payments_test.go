// Table-driven tests for payments' Payment aggregate: the happy charge path,
// the provider-declined guard, and acl's validation at the boundary.
// Persistence tests need a live Postgres — set TEST_DATABASE_URL to run
// them; otherwise they're skipped, since this module does not pin
// testcontainers-go. The provider is stubbed with an httptest server.
package payments_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/acme/shop/services/payments/internal/acl"
	. "github.com/acme/shop/services/payments/internal/payments"
)

func testDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("no container runtime present: set TEST_DATABASE_URL to run against a live Postgres")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	if err := db.AutoMigrate(&Payment{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func stubProvider(t *testing.T, status int) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestServiceChargeHappyPath(t *testing.T) {
	db := testDB(t)
	provider := stubProvider(t, http.StatusOK)
	svc := NewService(db, NewGateway(provider.URL))

	if err := svc.Charge(db, 1, 42.5); err != nil {
		t.Fatalf("Charge: %v", err)
	}

	var stored Payment
	if err := db.Where("order_id = ?", 1).First(&stored).Error; err != nil {
		t.Fatalf("load payment: %v", err)
	}
	if stored.Status != "captured" {
		t.Fatalf("status = %s, want captured", stored.Status)
	}
	if stored.Amount != 42.5 {
		t.Fatalf("amount = %v, want 42.5", stored.Amount)
	}
}

func TestServiceChargeRejectsProviderDecline(t *testing.T) {
	db := testDB(t)
	provider := stubProvider(t, http.StatusInternalServerError)
	svc := NewService(db, NewGateway(provider.URL))

	if err := svc.Charge(db, 2, 10); err == nil {
		t.Fatal("expected the provider decline to surface as an error")
	}
	var count int64
	db.Model(&Payment{}).Where("order_id = ?", 2).Count(&count)
	if count != 0 {
		t.Fatalf("payments for order 2 = %d, want 0 — no row on decline", count)
	}
}

func TestTranslatePaymentValidation(t *testing.T) {
	cases := []struct {
		name    string
		legacy  acl.LegacyPayment
		wantErr bool
	}{
		{"valid captured", acl.LegacyPayment{ID: 1, OrderID: 1, Amount: 10, Status: "captured"}, false},
		{"zero order id", acl.LegacyPayment{ID: 2, OrderID: 0, Amount: 10, Status: "captured"}, true},
		{"non-positive amount", acl.LegacyPayment{ID: 3, OrderID: 1, Amount: 0, Status: "captured"}, true},
		{"unknown status", acl.LegacyPayment{ID: 4, OrderID: 1, Amount: 10, Status: "pending"}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := acl.TranslatePayment(tc.legacy)
			if (err != nil) != tc.wantErr {
				t.Fatalf("TranslatePayment(%+v) error = %v, wantErr %v", tc.legacy, err, tc.wantErr)
			}
		})
	}
}
