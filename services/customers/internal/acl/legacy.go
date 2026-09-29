// Package acl translates the monolith's customer shapes into customers's model.
// The customers context has no stringly-typed statuses or renamed legacy columns,
// so the only translation is rejecting invalid customers at the boundary.
package acl

import (
	"errors"
	"net/mail"
	"strings"
)

// ValidateCustomer rejects a customer the service will not store.
func ValidateCustomer(email, name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("name is required")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("email is invalid")
	}
	return nil
}
