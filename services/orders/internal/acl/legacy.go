// Package acl translates the monolith's shapes for Order and OrderLine into
// orders's model. Pure functions, no I/O: legacy column names are mapped and
// invariants are rejected at the boundary.
package acl

import (
	"errors"

	"github.com/acme/shop/services/orders/internal/orders"
)

// LegacyOrderLine is an order line as the monolith's shared database named it.
type LegacyOrderLine struct {
	ID        uint    `json:"id"`
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

// ErrInvalidOrderLine is returned when a legacy row breaks an invariant.
var ErrInvalidOrderLine = errors.New("invalid order line")

// ToOrderLine maps a legacy row to orders's OrderLine. A zero product id, a
// non-positive quantity or a negative price is rejected.
func ToOrderLine(l LegacyOrderLine) (orders.OrderLine, error) {
	if l.ProductID == 0 || l.Quantity <= 0 || l.UnitPrice < 0 {
		return orders.OrderLine{}, ErrInvalidOrderLine
	}
	return orders.OrderLine{ID: l.ID, OrderID: l.OrderID, ProductID: l.ProductID, Quantity: l.Quantity, UnitPrice: l.UnitPrice}, nil
}

// ToLegacy maps an OrderLine back to the monolith's row shape.
func ToLegacy(l orders.OrderLine) LegacyOrderLine {
	return LegacyOrderLine{ID: l.ID, OrderID: l.OrderID, ProductID: l.ProductID, Quantity: l.Quantity, UnitPrice: l.UnitPrice}
}
