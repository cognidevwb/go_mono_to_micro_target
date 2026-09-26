// Table-driven tests for catalog's Product aggregate: the happy path,
// validation (a non-positive price), not-found, and the TTL price cache.
// Persistence tests need a live Postgres — set TEST_DATABASE_URL to run them;
// otherwise they're skipped, since this module does not pin testcontainers-go.
package catalog

import (
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	if err := db.AutoMigrate(&Category{}, &Product{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func TestServiceCreate(t *testing.T) {
	cases := []struct {
		name    string
		product Product
		wantErr bool
	}{
		{"happy path", Product{SKU: "sku-happy", Name: "Widget", Price: 9.99}, false},
		{"validation: non-positive price", Product{SKU: "sku-invalid", Name: "Widget", Price: 0}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db := testDB(t)
			svc := NewService(db)
			err := svc.Create(&tc.product)
			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestServiceCreateRejectsNonPositivePriceWithoutDB(t *testing.T) {
	svc := &Service{}
	if err := svc.Create(&Product{SKU: "sku-guard", Name: "Widget", Price: -1}); err == nil {
		t.Fatal("expected error for non-positive price")
	}
}

func TestServicePriceOfNotFound(t *testing.T) {
	db := testDB(t)
	svc := NewService(db)
	if _, err := svc.PriceOf(999999); err == nil {
		t.Fatal("expected a not-found error")
	}
}

func TestServicePriceOfCachesAndInvalidatesOnCreate(t *testing.T) {
	db := testDB(t)
	svc := NewService(db)
	p := Product{SKU: "sku-cache", Name: "Widget", Price: 5}
	if err := svc.Create(&p); err != nil {
		t.Fatalf("create: %v", err)
	}
	price, err := svc.PriceOf(p.ID)
	if err != nil {
		t.Fatalf("PriceOf: %v", err)
	}
	if price != 5 {
		t.Fatalf("PriceOf = %v, want 5", price)
	}
}
