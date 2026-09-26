// Equivalence tests: replay a golden {request, expectedResponse} recorded
// from the monolith's catalog handler against this service via httptest, and
// assert the new response matches the legacy one once volatile fields (the
// server-assigned id) are normalised away. Needs a live Postgres — set
// TEST_DATABASE_URL to run; otherwise skipped, since this module does not pin
// testcontainers-go.
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

	"github.com/acme/shop/services/catalog/internal/catalog"
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
	if err := db.AutoMigrate(&catalog.Category{}, &catalog.Product{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	gin.SetMode(gin.TestMode)
	r := gin.New()
	catalog.RegisterRoutes(r.Group("/api"), catalog.NewService(db))
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

func TestCreateProductEquivalence(t *testing.T) {
	router := newRouter(t)
	g := loadGolden(t, "create_product.json")

	req := httptest.NewRequest(g.Method, g.Path, bytes.NewReader(g.Body))
	req.Header.Set("Content-Type", "application/json")
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

func TestListProductsEquivalence(t *testing.T) {
	router := newRouter(t)

	create := loadGolden(t, "create_product.json")
	req := httptest.NewRequest(create.Method, create.Path, bytes.NewReader(create.Body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(httptest.NewRecorder(), req)

	req = httptest.NewRequest(http.MethodGet, "/api/products", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rec.Code, rec.Body.String())
	}
	var products []catalog.Product
	if err := json.Unmarshal(rec.Body.Bytes(), &products); err != nil {
		t.Fatalf("decode products: %v", err)
	}
	if len(products) == 0 {
		t.Fatal("expected at least one product")
	}
}
