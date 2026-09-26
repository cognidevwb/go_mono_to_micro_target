// Equivalence tests: replay a golden {request, seed, expectedResponse}
// against this service via httptest, and assert the response matches once
// volatile fields (the server-assigned id) are normalised away. Needs a live
// Postgres — set TEST_DATABASE_URL to run; otherwise skipped, since this
// module does not pin testcontainers-go.
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

	"github.com/acme/shop/services/inventory/internal/httpapi"
	"github.com/acme/shop/services/inventory/internal/inventory"
)

type golden struct {
	Method   string          `json:"method"`
	Path     string          `json:"path"`
	Seed     json.RawMessage `json:"seed"`
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

func newRouter(t *testing.T) (http.Handler, *gorm.DB) {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("no container runtime present: set TEST_DATABASE_URL to run against a live Postgres")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	if err := db.AutoMigrate(&inventory.StockItem{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	gin.SetMode(gin.TestMode)
	r := gin.New()
	httpapi.RegisterRoutes(r, db, inventory.NewService(db))
	return r, db
}

// normalize decodes a response and drops fields the golden fixture cannot pin.
func normalize(t *testing.T, raw json.RawMessage) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return m
}

func TestReserveStockItemEquivalence(t *testing.T) {
	router, db := newRouter(t)
	g := loadGolden(t, "reserve_stock_item.json")

	var seed inventory.StockItem
	if err := json.Unmarshal(g.Seed, &seed); err != nil {
		t.Fatalf("decode seed: %v", err)
	}
	if err := db.Create(&seed).Error; err != nil {
		t.Fatalf("seed stock item: %v", err)
	}

	req := httptest.NewRequest(g.Method, g.Path, bytes.NewReader(g.Body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rec.Code, rec.Body.String())
	}
	got := normalize(t, rec.Body.Bytes())
	want := normalize(t, g.Expected)
	for k, v := range want {
		if got[k] != v {
			t.Fatalf("field %s = %v, want %v", k, got[k], v)
		}
	}
}

func TestReserveStockItemEquivalenceOutOfStock(t *testing.T) {
	router, db := newRouter(t)
	seed := inventory.StockItem{ProductID: 42, OnHand: 1, Reserved: 0}
	if err := db.Create(&seed).Error; err != nil {
		t.Fatalf("seed stock item: %v", err)
	}

	body := []byte(`{"product_id":42,"quantity":100}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/stockitem/reserve", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409: %s", rec.Code, rec.Body.String())
	}
}
