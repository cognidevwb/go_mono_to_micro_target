// Package acl translates the monolith's on-the-wire shape for Payment into
// payments' own model: the bare status string becomes a typed constant,
// legacy column names are remapped, and invariants the monolith let through
// are rejected at the boundary. Pure functions, no I/O.
package acl

import (
	"errors"

	"github.com/acme/shop/services/payments/internal/payments"
)

// PaymentStatus is payments' own typed status, replacing the monolith's bare string.
type PaymentStatus string

// The statuses the monolith recorded on a payment row.
const (
	PaymentCaptured PaymentStatus = "captured"
	PaymentFailed   PaymentStatus = "failed"
)

var legacyStatus = map[string]PaymentStatus{
	"captured": PaymentCaptured,
	"failed":   PaymentFailed,
}

// LegacyPayment is the monolith's row shape for a payment.
type LegacyPayment struct {
	ID      uint    `json:"id"`
	OrderID uint    `json:"order_id"`
	Amount  float64 `json:"amount"`
	Status  string  `json:"status"`
}

// TranslatePayment maps a legacy payment row into payments' own model,
// rejecting the invariants the monolith let through.
func TranslatePayment(l LegacyPayment) (payments.Payment, error) {
	if l.OrderID == 0 {
		return payments.Payment{}, errors.New("order id is required")
	}
	if l.Amount <= 0 {
		return payments.Payment{}, errors.New("amount must be positive")
	}
	status, ok := legacyStatus[l.Status]
	if !ok {
		return payments.Payment{}, errors.New("unknown payment status: " + l.Status)
	}
	return payments.Payment{
		ID:      l.ID,
		OrderID: l.OrderID,
		Amount:  l.Amount,
		Status:  string(status),
	}, nil
}
