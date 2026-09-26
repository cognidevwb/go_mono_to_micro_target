// Package acl translates the monolith's on-the-wire shape for Customer into
// customers's own model: legacy column names are remapped and invariants the
// monolith let through are rejected at the boundary. Pure functions, no I/O.
package acl

import (
	"errors"

	"github.com/acme/shop/services/customers/internal/customers"
)

// LegacyCustomer is the monolith's row shape for a customer, using its legacy
// column names.
type LegacyCustomer struct {
	ID        uint   `json:"id"`
	EmailAddr string `json:"email_address"`
	FullName  string `json:"full_name"`
	IsActive  bool   `json:"is_active"`
}

// TranslateCustomer maps a legacy customer row into customers's own model,
// rejecting invariants (a blank email, a blank name) the monolith's storage
// did not enforce.
func TranslateCustomer(l LegacyCustomer) (customers.Customer, error) {
	if l.EmailAddr == "" {
		return customers.Customer{}, errors.New("email is required")
	}
	if l.FullName == "" {
		return customers.Customer{}, errors.New("name is required")
	}
	return customers.Customer{
		ID:     l.ID,
		Email:  l.EmailAddr,
		Name:   l.FullName,
		Active: l.IsActive,
	}, nil
}
