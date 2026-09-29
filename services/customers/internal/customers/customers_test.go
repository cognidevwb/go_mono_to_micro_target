package customers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRegisterCustomerValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cases := []struct {
		name string
		body string
	}{
		{"missing fields", `{}`},
		{"bad email", `{"email":"nope","name":"Ann"}`},
		{"missing name", `{"email":"a@b.co"}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := gin.New()
			RegisterRoutes(r.Group("/api"), NewService(nil))
			req := httptest.NewRequest(http.MethodPost, "/api/customers", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
			}
		})
	}
}

func TestGetCustomerBadID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	RegisterRoutes(r.Group("/api"), NewService(nil))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/customers/abc", nil))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}
