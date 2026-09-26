package httpx

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReadinessReportsAFailedDependency(t *testing.T) {
	h := NewHealth()
	h.AddReadiness("db", func(context.Context) error { return errors.New("down") })
	mux := http.NewServeMux()
	h.Mount(mux)
	for path, want := range map[string]int{"/healthz": 200, "/readyz": 503} {
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		if rec.Code != want {
			t.Fatalf("%s = %d, want %d", path, rec.Code, want)
		}
	}
}
