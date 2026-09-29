// Package store opens catalog's own database and migrates only its own tables.
package store

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/acme/shop/services/catalog/internal/catalog"
)

// Open connects catalog's database from its own DSN and migrates Category and Product.
func Open(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&catalog.Category{}, &catalog.Product{}); err != nil {
		return nil, err
	}
	return db, nil
}
