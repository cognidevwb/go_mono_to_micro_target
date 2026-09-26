// Table-driven tests for orders' Order aggregate: the happy path, the
// inactive-customer guard, and the out-of-stock compensation path.
// Persistence tests need a live Postgres — set TEST_DATABASE_URL to run
// them; otherwise they're skipped, since this module does not pin
// testcontainers-go. Peer services (catalog, customers, inventory, payments)
// are stubbed with httptest servers and wired in through their env vars.
package orders

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/acme/shop/services/orders/internal/clients/catalog"
	"github.com/acme/shop/services/orders/internal/clients/customers"
	"github.com/acme/shop/services/orders/internal/clients/inventory"
	"github.com/acme/shop/services/orders/internal/clients/payments"
	"github.com/acme/shop/services/orders/internal/platform"
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
	if err := db.AutoMigrate(&Order{}, &OrderLine{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func stubJSON(t *testing.T, status int, body any) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(body)
	}))
	t.Cleanup(srv.Close)
	return srv
}

func newTestService(t *testing.T, db *gorm.DB, customerActive bool, price float64, reserve, charge bool) *Service {
	t.Helper()

	customersSrv := stubJSON(t, http.StatusOK, map[string]any{"active": customerActive})
	t.Setenv("CUSTOMERS_SERVICE_URL", customersSrv.URL)

	catalogSrv := stubJSON(t, http.StatusOK, []catalog.Product{{ID: 1, SKU: "sku-1", Name: "Widget", Price: price}})
	t.Setenv("CATALOG_SERVICE_URL", catalogSrv.URL)

	invStatus := http.StatusOK
	if !reserve {
		invStatus = http.StatusConflict
	}
	inventorySrv := stubJSON(t, invStatus, map[string]any{"reserved": reserve})
	t.Setenv("INVENTORY_SERVICE_URL", inventorySrv.URL)

	payStatus := http.StatusOK
	if !charge {
		payStatus = http.StatusInternalServerError
	}
	paymentsSrv := stubJSON(t, payStatus, map[string]any{})
	t.Setenv("PAYMENTS_SERVICE_URL", paymentsSrv.URL)

	return NewService(db, catalog.NewService(db), customers.NewService(db), inventory.NewService(db),
		payments.NewService(db, payments.NewGateway("")), platform.NewBus())
}

func TestServicePlaceOrderHappyPath(t *testing.T) {
	db := testDB(t)
	svc := newTestService(t, db, true, 10, true, true)

	order, err := svc.PlaceOrder(PlaceOrderRequest{CustomerID: 1, Items: []LineItem{{ProductID: 1, Quantity: 2}}})
	if err != nil {
		t.Fatalf("PlaceOrder: %v", err)
	}
	if order.Status != string(StatusPlaced) {
		t.Fatalf("status = %s, want %s", order.Status, StatusPlaced)
	}
	if order.Total != 20 {
		t.Fatalf("total = %v, want 20", order.Total)
	}
}

func TestServicePlaceOrderRejectsInactiveCustomer(t *testing.T) {
	db := testDB(t)
	svc := newTestService(t, db, false, 10, true, true)

	if _, err := svc.PlaceOrder(PlaceOrderRequest{CustomerID: 1, Items: []LineItem{{ProductID: 1, Quantity: 1}}}); err == nil {
		t.Fatal("expected error for an inactive customer")
	}
}

func TestServicePlaceOrderCompensatesOnOutOfStock(t *testing.T) {
	db := testDB(t)
	svc := newTestService(t, db, true, 10, false, true)

	if _, err := svc.PlaceOrder(PlaceOrderRequest{CustomerID: 1, Items: []LineItem{{ProductID: 1, Quantity: 1}}}); err == nil {
		t.Fatal("expected an out-of-stock error")
	}
	var stored Order
	if err := db.Order("id DESC").First(&stored).Error; err != nil {
		t.Fatalf("load order: %v", err)
	}
	if stored.Status != string(StatusFailed) {
		t.Fatalf("status = %s, want %s", stored.Status, StatusFailed)
	}
}

func TestServiceForCustomerListsOwnOrders(t *testing.T) {
	db := testDB(t)
	svc := newTestService(t, db, true, 5, true, true)

	if _, err := svc.PlaceOrder(PlaceOrderRequest{CustomerID: 42, Items: []LineItem{{ProductID: 1, Quantity: 1}}}); err != nil {
		t.Fatalf("PlaceOrder: %v", err)
	}
	list, err := svc.ForCustomer(42)
	if err != nil {
		t.Fatalf("ForCustomer: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("orders = %d, want 1", len(list))
	}
}
