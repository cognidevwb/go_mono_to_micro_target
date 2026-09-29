package equivalence

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/acme/shop/services/inventory/internal/httpapi"
)

// TestStockItemBadIDIsProblemJSON replays an invalid-request golden: a
// non-numeric product id answers 400 problem+json before any database access.
func TestStockItemBadIDIsProblemJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	httpapi.RegisterRoutes(r, nil, nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/v1/stockitem/abc", nil))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/problem+json" {
		t.Fatalf("content-type = %q, want application/problem+json", ct)
	}
}
