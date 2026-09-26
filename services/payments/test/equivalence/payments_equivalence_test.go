// Equivalence tests: replay a golden {request, expectedResponse} recorded
// from the monolith's payments handler against this service via httptest,
// and assert the new response matches the legacy one once volatile fields
// (the server-assigned id) are normalised away. Needs a live Postgres — set
// TEST_DATABASE_URL to run; otherwise skipped, since this module does not
// pin testcontainers-go. The provider is stubbed with an httptest server.
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

	"github.com/acme/shop/services/payments/internal/httpapi"
	"github.com/acme/shop/services/payments/internal/payments"
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

func stubProvider(t *testing.T, status int) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
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
	if err := db.AutoMigrate(&payments.Payment{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	provider := stubProvider(t, http.StatusOK)
	svc := payments.NewService(db, payments.NewGateway(provider.URL))

	gin.SetMode(gin.TestMode)
	r := gin.New()
	httpapi.RegisterRoutes(r, db, svc)
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
	delete(m, "ID")
	return m
}

func TestChargeEquivalence(t *testing.T) {
	router := newRouter(t)
	g := loadGolden(t, "charge.json")

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
