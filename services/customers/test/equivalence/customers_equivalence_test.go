package equivalence

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/acme/shop/services/customers/internal/customers"
)

// TestRegisterCustomerValidation replays the monolith's invalid-request golden:
// a body missing required fields answers 400 before the service is reached.
func TestRegisterCustomerValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	customers.RegisterRoutes(r.Group("/api"), customers.NewService(nil))

	req := httptest.NewRequest(http.MethodPost, "/api/customers", strings.NewReader(`{"name":""}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}
