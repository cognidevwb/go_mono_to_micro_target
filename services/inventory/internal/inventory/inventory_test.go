// Table-driven tests for inventory's StockItem aggregate: the happy path,
// not-found, and the concurrency / out-of-stock guard on Reserve, plus
// Restock. Persistence tests need a live Postgres — set TEST_DATABASE_URL to
// run them; otherwise they're skipped, since this module does not pin
// testcontainers-go.
package inventory

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
	if err := db.AutoMigrate(&StockItem{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func TestServiceReserve(t *testing.T) {
	db := testDB(t)
	svc := NewService(db)
	item := StockItem{ProductID: 1, OnHand: 10, Reserved: 0}
	if err := db.Create(&item).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}

	cases := []struct {
		name   string
		qty    int
		wantOK bool
	}{
		{"happy path", 5, true},
		{"out of stock guard", 100, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var ok bool
			err := db.Transaction(func(tx *gorm.DB) error {
				var e error
				ok, e = svc.Reserve(tx, item.ProductID, tc.qty)
				return e
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if ok != tc.wantOK {
				t.Fatalf("Reserve ok = %v, want %v", ok, tc.wantOK)
			}
		})
	}
}

func TestServiceReserveNotFound(t *testing.T) {
	db := testDB(t)
	svc := NewService(db)
	err := db.Transaction(func(tx *gorm.DB) error {
		_, e := svc.Reserve(tx, 999999, 1)
		return e
	})
	if err == nil {
		t.Fatal("expected a not-found error")
	}
}

func TestServiceRestockTopsUpBelowReserved(t *testing.T) {
	db := testDB(t)
	svc := NewService(db)
	item := StockItem{ProductID: 2, OnHand: 1, Reserved: 5}
	if err := db.Create(&item).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}
	if err := svc.Restock(); err != nil {
		t.Fatalf("Restock: %v", err)
	}
	var got StockItem
	if err := db.First(&got, item.ID).Error; err != nil {
		t.Fatalf("reload: %v", err)
	}
	if got.OnHand != got.Reserved+10 {
		t.Fatalf("OnHand = %d, want %d", got.OnHand, got.Reserved+10)
	}
}
