// Package store opens inventory's own *gorm.DB from its own DSN and
// migrates only the table this service owns — StockItem. The monolith's
// shared platform.Open and its global AutoMigrate do not come over.
package store

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/acme/shop/services/inventory/internal/inventory"
)

// Open connects inventory's own database and migrates its owned model.
func Open(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&inventory.StockItem{}); err != nil {
		return nil, err
	}
	return db, nil
}
