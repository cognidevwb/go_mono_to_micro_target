// Package acl translates the monolith's on-the-wire shapes for StockItem
// into inventory's own model: legacy column names are remapped, the
// stringly-typed status the monolith let callers pass is replaced by a typed
// constant derived from the item itself, and invariants the monolith's
// storage never enforced are rejected at the boundary. Pure functions, no I/O.
package acl

import (
	"errors"

	"github.com/acme/shop/services/inventory/internal/inventory"
)

// StockStatus is inventory's typed replacement for the monolith's ad-hoc
// status strings.
type StockStatus string

const (
	StatusInStock    StockStatus = "in_stock"
	StatusOutOfStock StockStatus = "out_of_stock"
)

// LegacyStockItem is the monolith's row shape for a stock item, using its
// legacy column names.
type LegacyStockItem struct {
	ProductID   uint `json:"product_id"`
	QtyOnHand   int  `json:"qty_on_hand"`
	QtyReserved int  `json:"qty_reserved"`
}

// StockItemView is the shape returned to callers: inventory's model plus the
// derived, typed status.
type StockItemView struct {
	ProductID uint        `json:"product_id"`
	OnHand    int         `json:"on_hand"`
	Reserved  int         `json:"reserved"`
	Available int         `json:"available"`
	Status    StockStatus `json:"status"`
}

// TranslateStockItem maps a legacy stock item row into inventory's own
// model, rejecting invariants (a missing product, negative quantities, more
// reserved than on hand) the monolith's storage let through.
func TranslateStockItem(l LegacyStockItem) (inventory.StockItem, error) {
	if l.ProductID == 0 {
		return inventory.StockItem{}, errors.New("product_id is required")
	}
	if l.QtyOnHand < 0 {
		return inventory.StockItem{}, errors.New("qty_on_hand must not be negative")
	}
	if l.QtyReserved < 0 {
		return inventory.StockItem{}, errors.New("qty_reserved must not be negative")
	}
	if l.QtyReserved > l.QtyOnHand {
		return inventory.StockItem{}, errors.New("qty_reserved cannot exceed qty_on_hand")
	}
	return inventory.StockItem{
		ProductID: l.ProductID,
		OnHand:    l.QtyOnHand,
		Reserved:  l.QtyReserved,
	}, nil
}

// ViewOf derives the caller-facing view of item, including its typed status.
func ViewOf(item inventory.StockItem) StockItemView {
	available := item.OnHand - item.Reserved
	status := StatusInStock
	if available <= 0 {
		status = StatusOutOfStock
	}
	return StockItemView{
		ProductID: item.ProductID,
		OnHand:    item.OnHand,
		Reserved:  item.Reserved,
		Available: available,
		Status:    status,
	}
}
