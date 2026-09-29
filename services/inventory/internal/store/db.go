// Package store opens inventory's own database. Only StockItem is registered;
// the monolith's shared platform.Open and global AutoMigrate do not come over.
package store

import (
	"context"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/acme/shop/services/inventory/internal/inventory"
)

// Open connects to inventory's own DSN and migrates the tables it owns.
func Open(ctx context.Context, dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	if err := db.WithContext(ctx).AutoMigrate(&inventory.StockItem{}); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return db, nil
}
