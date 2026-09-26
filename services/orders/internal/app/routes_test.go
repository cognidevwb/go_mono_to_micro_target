package app

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoutesBuildWithoutLiveDependencies(t *testing.T) {
	h, err := Routes(&Deps{})
	if err != nil {
		t.Fatalf("Routes: %v", err)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/__no_such_route__", nil))
	// 404 before the port; after it, the monolith's own middleware may answer
	// first (401). Either way an unknown route is a client error, never served.
	if rec.Code < 400 || rec.Code >= 500 {
		t.Fatalf("unknown route = %d, want a 4xx", rec.Code)
	}
}
