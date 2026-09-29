// Package store opens orders's own database. Only Order and OrderLine are
// registered; the monolith's shared platform.Open and global AutoMigrate do not
// come over.
package store

import (
	"context"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/acme/shop/services/orders/internal/orders"
)

// Open connects to orders's own DSN and migrates the tables it owns.
func Open(ctx context.Context, dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	if err := db.WithContext(ctx).AutoMigrate(&orders.Order{}, &orders.OrderLine{}); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return db, nil
}
