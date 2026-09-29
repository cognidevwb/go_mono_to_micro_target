// Package store opens payments's OWN database, migrating only the tables this
// service owns.
package store

import (
	"context"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/acme/shop/services/payments/internal/payments"
)

// Open connects to payments's own DSN and migrates the Payment table.
func Open(ctx context.Context, dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	if err := db.WithContext(ctx).AutoMigrate(&payments.Payment{}); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return db, nil
}
