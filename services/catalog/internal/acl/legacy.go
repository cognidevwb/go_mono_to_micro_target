// Package acl translates the monolith's on-the-wire shapes for Category and
// Product into catalog's own model: legacy column names are remapped and
// invariants the monolith let through are rejected at the boundary. Pure
// functions, no I/O.
package acl

import (
	"errors"

	"github.com/acme/shop/services/catalog/internal/catalog"
)

// LegacyCategory is the monolith's row shape for a category.
type LegacyCategory struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// LegacyProduct is the monolith's row shape for a product, using its legacy
// column names.
type LegacyProduct struct {
	ID         uint    `json:"id"`
	SKUCode    string  `json:"sku_code"`
	Title      string  `json:"title"`
	UnitPrice  float64 `json:"unit_price"`
	CategoryID uint    `json:"category_id"`
}

// TranslateCategory maps a legacy category row into catalog's own model.
func TranslateCategory(l LegacyCategory) (catalog.Category, error) {
	if l.Name == "" {
		return catalog.Category{}, errors.New("category name is required")
	}
	return catalog.Category{ID: l.ID, Name: l.Name}, nil
}

// TranslateProduct maps a legacy product row into catalog's own model,
// rejecting invariants (a blank SKU, a non-positive price) the monolith's
// storage did not enforce.
func TranslateProduct(l LegacyProduct) (catalog.Product, error) {
	if l.SKUCode == "" {
		return catalog.Product{}, errors.New("sku is required")
	}
	if l.UnitPrice <= 0 {
		return catalog.Product{}, errors.New("price must be positive")
	}
	return catalog.Product{
		ID:         l.ID,
		SKU:        l.SKUCode,
		Name:       l.Title,
		Price:      l.UnitPrice,
		CategoryID: l.CategoryID,
	}, nil
}
