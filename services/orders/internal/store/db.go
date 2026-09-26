// Package store opens orders' own *gorm.DB from its own DSN and migrates
// only the tables this service owns — Order and OrderLine. The monolith's
// shared platform.Open and its global AutoMigrate do not come over.
package store

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/acme/shop/services/orders/internal/orders"
)

// Open connects orders' own database and migrates its owned models.
func Open(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&orders.Order{}, &orders.OrderLine{}); err != nil {
		return nil, err
	}
	return db, nil
}
