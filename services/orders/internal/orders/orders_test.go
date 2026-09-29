package orders

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPlaceOrderRejectsInvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cases := []struct {
		name string
		body string
	}{
		{"empty object", `{}`},
		{"no items", `{"customerId":1,"items":[]}`},
		{"zero quantity", `{"customerId":1,"items":[{"productId":1,"quantity":0}]}`},
		{"not json", `nope`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := gin.New()
			RegisterRoutes(r.Group("/api"), &Service{})
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/orders", bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}
