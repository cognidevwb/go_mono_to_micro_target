// contracts · Contract tests orders->customers: for each testdata/contracts/customers/*.json golden response assert it decodes into orders's view of customers-service's Customer with no unknown fields (json.Decoder.DisallowUnknownFields) and satisfies the value invariants; record one fixture (orders-service)

package clients

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	customersclient "github.com/acme/shop/services/orders/internal/clients/customers"
)

// customerFixture mirrors customers-service's Customer as orders reads it off
// the wire: exported fields and json tags copied from the provider's model.
type customerFixture struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"createdAt"`
}

func TestCustomersContract_CustomerDecodesWithInvariants(t *testing.T) {
	raw, err := os.ReadFile("testdata/contracts/customers/customer.json")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	var c customerFixture
	if err := dec.Decode(&c); err != nil {
		t.Fatalf("decode fixture: %v", err)
	}

	if c.ID == 0 {
		t.Error("expected non-zero customer id")
	}
	if c.Email == "" {
		t.Error("expected non-empty email")
	}
	if c.Name == "" {
		t.Error("expected non-empty name")
	}
	if c.CreatedAt.IsZero() {
		t.Error("expected non-zero createdAt")
	}
}

func TestCustomersContract_IsActiveReplaysGoldenResponse(t *testing.T) {
	raw, err := os.ReadFile("testdata/contracts/customers/customer.json")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	var fixture customerFixture
	if err := json.Unmarshal(raw, &fixture); err != nil {
		t.Fatalf("unmarshal fixture: %v", err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/customers/5" || r.Method != http.MethodGet {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(fixture)
	}))
	defer srv.Close()

	t.Setenv("CUSTOMERS_SERVICE_URL", srv.URL)
	svc := customersclient.NewService(nil)

	active, err := svc.IsActive(fixture.ID)
	if err != nil {
		t.Fatalf("IsActive: %v", err)
	}
	if active != fixture.Active {
		t.Errorf("expected active %v, got %v", fixture.Active, active)
	}
}
