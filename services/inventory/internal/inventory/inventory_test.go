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
		t.Skip("TEST_DATABASE_URL not set; no Postgres available")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := db.AutoMigrate(&StockItem{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func TestReserve(t *testing.T) {
	db := testDB(t)
	const pid = 900001
	db.Where("product_id = ?", pid).Delete(&StockItem{})
	if err := db.Create(&StockItem{ProductID: pid, OnHand: 5}).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}
	t.Cleanup(func() { db.Where("product_id = ?", pid).Delete(&StockItem{}) })

	cases := []struct {
		name    string
		product uint
		qty     int
		want    bool
		wantErr bool
	}{
		{"within stock", pid, 3, true, false},
		{"out of stock", pid, 3, false, false},
		{"unknown product", 900002, 1, false, true},
	}
	svc := NewService(db)
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := svc.Reserve(db, tc.product, tc.qty)
			if (err != nil) != tc.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tc.wantErr)
			}
			if got != tc.want {
				t.Fatalf("reserved = %v, want %v", got, tc.want)
			}
		})
	}
}
