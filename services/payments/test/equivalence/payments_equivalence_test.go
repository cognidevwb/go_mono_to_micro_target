// Package equivalence replays the monolith's payment request shapes against the
// service's own handlers.
package equivalence

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/acme/shop/services/payments/internal/app"
)

func TestChargeRequestValidation(t *testing.T) {
	h, err := app.Routes(&app.Deps{})
	if err != nil {
		t.Fatalf("Routes: %v", err)
	}
	tests := []struct {
		name string
		body string
		want int
	}{
		{"malformed", `{`, http.StatusBadRequest},
		{"missing order", `{"amount": 5}`, http.StatusBadRequest},
		{"negative amount", `{"orderId": 1, "amount": -1}`, http.StatusBadRequest},
		{"valid but no database", `{"orderId": 1, "amount": 5}`, http.StatusServiceUnavailable},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/v1/payments/charge", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Internal-Caller", "orders-service")
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			if rec.Code != tc.want {
				t.Fatalf("status = %d, want %d (%s)", rec.Code, tc.want, rec.Body.String())
			}
		})
	}
}
