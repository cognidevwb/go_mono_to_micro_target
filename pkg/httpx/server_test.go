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

func TestAServiceCallPassesTheUserCheckAndAnOutsideCallDoesNot(t *testing.T) {
	deny := func(http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusUnauthorized) })
	}
	h := UnlessInternal(deny)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	srv := httptest.NewServer(h)
	defer srv.Close()

	t.Setenv("INTERNAL_TOKEN", "s3cret")
	resp, err := NewClient().Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("service call = %d, want 200", resp.StatusCode)
	}
	req, _ := http.NewRequest(http.MethodGet, srv.URL, nil)
	req.Header.Set(CallerHeader, "orders-service")
	req.Header.Set(TokenHeader, "guess")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("a wrong token = %d, want 401", resp.StatusCode)
	}
}
