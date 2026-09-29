// Package acl translates the monolith's shapes for Payment into payments's own:
// stringly-typed statuses become typed constants and invariants are rejected at
// the boundary. Pure functions, no I/O.
package acl

import (
	"errors"
	"fmt"
	"math"
)

// PaymentStatus is the typed form of the monolith's bare Payment.Status string.
type PaymentStatus string

const (
	PaymentCaptured PaymentStatus = "captured"
	PaymentVoided   PaymentStatus = "voided"
)

// ParseStatus maps a legacy status string to its typed constant.
func ParseStatus(s string) (PaymentStatus, error) {
	switch PaymentStatus(s) {
	case PaymentCaptured, PaymentVoided:
		return PaymentStatus(s), nil
	}
	return "", fmt.Errorf("unknown payment status %q", s)
}

// ChargeRequest is the body orders sends to POST /v1/payments/charge.
type ChargeRequest struct {
	OrderID uint    `json:"orderId"`
	Amount  float64 `json:"amount"`
}

// Validate rejects a charge the provider must never see.
func (r ChargeRequest) Validate() error {
	if r.OrderID == 0 {
		return errors.New("orderId is required")
	}
	if math.IsNaN(r.Amount) || math.IsInf(r.Amount, 0) || r.Amount < 0 {
		return errors.New("amount must be a non-negative number")
	}
	return nil
}
