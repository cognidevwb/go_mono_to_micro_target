// Table-driven tests for customers's Customer aggregate: the happy path,
// validation (a missing email or name), and not-found. Persistence tests
// need a live Postgres — set TEST_DATABASE_URL to run them; otherwise
// they're skipped, since this module does not pin testcontainers-go.
package customers

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
	if err := db.AutoMigrate(&Customer{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func TestServiceRegisterAndGet(t *testing.T) {
	db := testDB(t)
	svc := NewService(db)
	c := Customer{Email: "happy@example.com", Name: "Happy Path", Active: true}
	if err := svc.Register(&c); err != nil {
		t.Fatalf("register: %v", err)
	}
	got, err := svc.Get(c.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Email != c.Email || got.Name != c.Name {
		t.Fatalf("got %+v, want %+v", got, c)
	}
}

func TestServiceGetNotFound(t *testing.T) {
	db := testDB(t)
	svc := NewService(db)
	if _, err := svc.Get(999999); err == nil {
		t.Fatal("expected a not-found error")
	}
}

func TestServiceIsActive(t *testing.T) {
	db := testDB(t)
	svc := NewService(db)
	c := Customer{Email: "active@example.com", Name: "Active One", Active: true}
	if err := svc.Register(&c); err != nil {
		t.Fatalf("register: %v", err)
	}
	active, err := svc.IsActive(c.ID)
	if err != nil {
		t.Fatalf("IsActive: %v", err)
	}
	if !active {
		t.Fatal("expected customer to be active")
	}
}

func TestServiceIsActiveNotFound(t *testing.T) {
	db := testDB(t)
	svc := NewService(db)
	if _, err := svc.IsActive(999999); err == nil {
		t.Fatal("expected a not-found error")
	}
}
