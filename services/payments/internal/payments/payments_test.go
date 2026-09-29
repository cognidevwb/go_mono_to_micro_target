package payments

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGatewayCharge(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		wantErr bool
	}{
		{"captured", http.StatusOK, false},
		{"declined", http.StatusPaymentRequired, true},
		{"provider down", http.StatusInternalServerError, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var got map[string]any
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/v1/charges" {
					t.Errorf("path = %s", r.URL.Path)
				}
				_ = json.NewDecoder(r.Body).Decode(&got)
				w.WriteHeader(tc.status)
			}))
			defer srv.Close()
			err := NewGateway(srv.URL).Charge(7, 12.5)
			if (err != nil) != tc.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tc.wantErr)
			}
			if got["amount"] != 12.5 {
				t.Fatalf("amount = %v", got["amount"])
			}
		})
	}
}

// A declined charge must not touch the database (tx is nil here).
func TestServiceChargeDeclinedWritesNothing(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusPaymentRequired)
	}))
	defer srv.Close()
	svc := NewService(nil, NewGateway(srv.URL))
	if err := svc.Charge(nil, 1, 10); err == nil {
		t.Fatal("expected provider decline error")
	}
}
