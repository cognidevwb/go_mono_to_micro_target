// Package store opens catalog's own *gorm.DB from its own DSN and migrates
// only the tables this service owns — Category and Product. The monolith's
// shared platform.Open and its global AutoMigrate do not come over.
package store

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/acme/shop/services/catalog/internal/catalog"
)

// Open connects catalog's own database and migrates its owned models.
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
