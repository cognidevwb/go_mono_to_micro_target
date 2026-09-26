// Package acl translates the monolith's on-the-wire shapes for Order and
// OrderLine into orders' own model: legacy column names are remapped,
// stringly-typed statuses become typed constants, and invariants the
// monolith let through are rejected at the boundary. Pure functions, no I/O.
package acl

import (
	"errors"

	"github.com/acme/shop/services/orders/internal/orders"
)

// LegacyOrderLine is the monolith's row shape for an order line.
type LegacyOrderLine struct {
	ID        uint    `json:"id"`
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Qty       int     `json:"qty"`
	UnitPrice float64 `json:"unit_price"`
}

// LegacyOrder is the monolith's row shape for an order, using its legacy
// column names and a bare status string.
type LegacyOrder struct {
	ID         uint              `json:"id"`
	CustomerID uint              `json:"customer_id"`
	State      string            `json:"state"`
	Total      float64           `json:"total"`
	Lines      []LegacyOrderLine `json:"lines"`
}

// legacyStatus maps the monolith's bare status strings to orders' typed
// constants.
var legacyStatus = map[string]orders.Status{
	"pending": orders.StatusPending,
	"placed":  orders.StatusPlaced,
	"failed":  orders.StatusFailed,
}

// TranslateOrderLine maps a legacy order line row into orders' own model.
func TranslateOrderLine(l LegacyOrderLine) (orders.OrderLine, error) {
	if l.Qty <= 0 {
		return orders.OrderLine{}, errors.New("quantity must be positive")
	}
	if l.UnitPrice < 0 {
		return orders.OrderLine{}, errors.New("unit price must not be negative")
	}
	return orders.OrderLine{
		ID:        l.ID,
		OrderID:   l.OrderID,
		ProductID: l.ProductID,
		Quantity:  l.Qty,
		UnitPrice: l.UnitPrice,
	}, nil
}

// TranslateOrder maps a legacy order row into orders' own model.
func TranslateOrder(l LegacyOrder) (orders.Order, error) {
	status, ok := legacyStatus[l.State]
	if !ok {
		return orders.Order{}, errors.New("unknown order status: " + l.State)
	}
	lines := make([]orders.OrderLine, 0, len(l.Lines))
	for _, ll := range l.Lines {
		line, err := TranslateOrderLine(ll)
		if err != nil {
			return orders.Order{}, err
		}
		lines = append(lines, line)
	}
	return orders.Order{
		ID:         l.ID,
		CustomerID: l.CustomerID,
		Status:     string(status),
		Total:      l.Total,
		Lines:      lines,
	}, nil
}
