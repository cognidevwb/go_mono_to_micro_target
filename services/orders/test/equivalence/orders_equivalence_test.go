package equivalence

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/acme/shop/services/orders/internal/orders"
)

// The monolith answers a malformed POST /api/orders with 400 before any
// collaborator is touched; the service must answer the same.
func TestPlaceOrderValidationMatchesMonolith(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	orders.RegisterRoutes(r.Group("/api"), &orders.Service{})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/orders", bytes.NewBufferString(`{"customerId":1,"items":[]}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
