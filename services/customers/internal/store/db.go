// Package store opens customers's own database and migrates only its own tables.
package store

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/acme/shop/services/customers/internal/customers"
)

// Open connects customers's database from its own DSN and migrates Customer.
func Open(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&customers.Customer{}); err != nil {
		return nil, err
	}
	return db, nil
}
