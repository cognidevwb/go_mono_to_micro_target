// Package acl translates the monolith's catalog shapes into catalog's model.
// The catalog context has no stringly-typed statuses or renamed legacy columns,
// so the only translation is rejecting invalid products at the boundary.
package acl

import "errors"

// ValidatePrice rejects a price the catalog will not store.
func ValidatePrice(price float64) error {
	if price <= 0 {
		return errors.New("price must be positive")
	}
	return nil
}
