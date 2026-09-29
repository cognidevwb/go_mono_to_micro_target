// Package acl translates the monolith's shapes for StockItem into inventory's
// model. Pure functions, no I/O: legacy column names are mapped and invariants
// are rejected at the boundary.
package acl

import (
	"errors"

	"github.com/acme/shop/services/inventory/internal/inventory"
)

// LegacyStockItem is a stock row as the monolith's shared database named it.
type LegacyStockItem struct {
	ID        uint `json:"id"`
	ProductID uint `json:"product_id"`
	OnHand    int  `json:"on_hand"`
	Reserved  int  `json:"reserved"`
}

// ErrInvalidStockItem is returned when a legacy row breaks an invariant.
var ErrInvalidStockItem = errors.New("invalid stock item")

// ToStockItem maps a legacy row to inventory's StockItem. A zero product id or
// a negative count is rejected. on_hand below reserved is accepted: the
// monolith's Restock job exists to repair exactly that state.
func ToStockItem(l LegacyStockItem) (inventory.StockItem, error) {
	if l.ProductID == 0 || l.OnHand < 0 || l.Reserved < 0 {
		return inventory.StockItem{}, ErrInvalidStockItem
	}
	return inventory.StockItem{ID: l.ID, ProductID: l.ProductID, OnHand: l.OnHand, Reserved: l.Reserved}, nil
}

// ToLegacy maps a StockItem back to the monolith's row shape.
func ToLegacy(s inventory.StockItem) LegacyStockItem {
	return LegacyStockItem{ID: s.ID, ProductID: s.ProductID, OnHand: s.OnHand, Reserved: s.Reserved}
}
