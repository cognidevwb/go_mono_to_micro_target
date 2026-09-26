// Equivalence tests: replay a golden {request, expectedResponse} recorded
// from the monolith's orders handler against this service via httptest, and
// assert the new response matches the legacy one once volatile fields (the
// server-assigned id) are normalised away. Needs a live Postgres — set
// TEST_DATABASE_URL to run; otherwise skipped, since this module does not
// pin testcontainers-go. Peer services are stubbed with httptest servers and
// wired in through their env vars.
package equivalence

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/acme/shop/services/orders/internal/clients/catalog"
	"github.com/acme/shop/services/orders/internal/clients/customers"
	"github.com/acme/shop/services/orders/internal/clients/inventory"
	"github.com/acme/shop/services/orders/internal/clients/payments"
	"github.com/acme/shop/services/orders/internal/orders"
	"github.com/acme/shop/services/orders/internal/platform"
)

type golden struct {
	Method   string          `json:"method"`
	Path     string          `json:"path"`
	Body     json.RawMessage `json:"body"`
	Expected json.RawMessage `json:"expected"`
}

func loadGolden(t *testing.T, name string) golden {
	t.Helper()
	data, err := os.ReadFile("testdata/equivalence/" + name)
	if err != nil {
		t.Fatalf("read golden %s: %v", name, err)
	}
	var g golden
	if err := json.Unmarshal(data, &g); err != nil {
		t.Fatalf("decode golden %s: %v", name, err)
	}
	return g
}

func stubJSON(t *testing.T, status int, body any) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(body)
	}))
	t.Cleanup(srv.Close)
	return srv
}

func newRouter(t *testing.T) http.Handler {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("no container runtime present: set TEST_DATABASE_URL to run against a live Postgres")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	if err := db.AutoMigrate(&orders.Order{}, &orders.OrderLine{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	customersSrv := stubJSON(t, http.StatusOK, map[string]any{"active": true})
	t.Setenv("CUSTOMERS_SERVICE_URL", customersSrv.URL)

	catalogSrv := stubJSON(t, http.StatusOK, []catalog.Product{{ID: 1, SKU: "sku-1", Name: "Widget", Price: 10}})
	t.Setenv("CATALOG_SERVICE_URL", catalogSrv.URL)

	inventorySrv := stubJSON(t, http.StatusOK, map[string]any{"reserved": true})
	t.Setenv("INVENTORY_SERVICE_URL", inventorySrv.URL)

	paymentsSrv := stubJSON(t, http.StatusOK, map[string]any{})
	t.Setenv("PAYMENTS_SERVICE_URL", paymentsSrv.URL)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(platform.AuthMiddleware())
	svc := orders.NewService(db, catalog.NewService(db), customers.NewService(db), inventory.NewService(db),
		payments.NewService(db, payments.NewGateway("")), platform.NewBus())
	orders.RegisterRoutes(r.Group("/api"), svc)
	return r
}

// normalize decodes a response and drops the server-assigned id, which the
// golden fixture cannot pin.
func normalize(t *testing.T, raw json.RawMessage) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	delete(m, "id")
	return m
}

func TestPlaceOrderEquivalence(t *testing.T) {
	router := newRouter(t)
	g := loadGolden(t, "place_order.json")

	req := httptest.NewRequest(g.Method, g.Path, bytes.NewReader(g.Body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "7")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201: %s", rec.Code, rec.Body.String())
	}
	got := normalize(t, rec.Body.Bytes())
	want := normalize(t, g.Expected)
	for k, v := range want {
		if got[k] != v {
			t.Fatalf("field %s = %v, want %v", k, got[k], v)
		}
	}
}

func TestMyOrdersEquivalence(t *testing.T) {
	router := newRouter(t)
	create := loadGolden(t, "place_order.json")

	req := httptest.NewRequest(create.Method, create.Path, bytes.NewReader(create.Body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "7")
	router.ServeHTTP(httptest.NewRecorder(), req)

	req = httptest.NewRequest(http.MethodGet, "/api/orders/mine", nil)
	req.Header.Set("X-User-ID", "7")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rec.Code, rec.Body.String())
	}
	var list []orders.Order
	if err := json.Unmarshal(rec.Body.Bytes(), &list); err != nil {
		t.Fatalf("decode orders: %v", err)
	}
	if len(list) == 0 {
		t.Fatal("expected at least one order")
	}
}
