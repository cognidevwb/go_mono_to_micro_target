// contracts · Contract tests orders->catalog: for each testdata/contracts/catalog/*.json golden response assert it decodes into internal/contracts/catalog with no unknown fields (json.Decoder.DisallowUnknownFields) and satisfies the value invariants; record one fixture (orders-service)

package clients

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	catalogclient "github.com/acme/shop/services/orders/internal/clients/catalog"
	contract "github.com/acme/shop/services/orders/internal/contracts/catalog"
)

func TestCatalogContract_ProductDecodesWithInvariants(t *testing.T) {
	raw, err := os.ReadFile("testdata/contracts/catalog/product.json")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	var p contract.Product
	if err := dec.Decode(&p); err != nil {
		t.Fatalf("decode fixture: %v", err)
	}

	if p.ID == 0 {
		t.Error("expected non-zero product id")
	}
	if p.SKU == "" {
		t.Error("expected non-empty sku")
	}
	if p.Name == "" {
		t.Error("expected non-empty name")
	}
	if p.Price <= 0 {
		t.Errorf("expected positive price, got %v", p.Price)
	}
}

func TestCatalogContract_PriceOfReplaysGoldenResponse(t *testing.T) {
	raw, err := os.ReadFile("testdata/contracts/catalog/product.json")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	var fixture contract.Product
	if err := json.Unmarshal(raw, &fixture); err != nil {
		t.Fatalf("unmarshal fixture: %v", err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/products" || r.Method != http.MethodGet {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]contract.Product{fixture})
	}))
	defer srv.Close()

	t.Setenv("CATALOG_SERVICE_URL", srv.URL)
	svc := catalogclient.NewService(nil)

	price, err := svc.PriceOf(fixture.ID)
	if err != nil {
		t.Fatalf("PriceOf: %v", err)
	}
	if price != fixture.Price {
		t.Errorf("expected price %v, got %v", fixture.Price, price)
	}
}
