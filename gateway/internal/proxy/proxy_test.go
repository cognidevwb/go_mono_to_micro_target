package proxy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func upstream(name string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, name)
	}))
}

func TestEveryRouteReachesItsService(t *testing.T) {
	legacy := upstream("legacy")
	defer legacy.Close()
	var routes []Route
	for _, r := range Routes() {
		up := upstream(r.Service)
		defer up.Close()
		r.Upstream = up.URL
		routes = append(routes, r)
	}
	h, err := New(routes, legacy.URL)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range routes {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, r.Prefix, nil))
		if got := rec.Body.String(); got != r.Service {
			t.Fatalf("%s reached %q, want %q", r.Prefix, got, r.Service)
		}
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/not-yet-extracted", nil))
	if rec.Body.String() != "legacy" {
		t.Fatalf("unmatched path reached %q, want legacy", rec.Body.String())
	}
}

func TestNoOutsideCallerReachesAServiceClaimingToBeOne(t *testing.T) {
	var seen http.Header
	up := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) { seen = r.Header.Clone() }))
	defer up.Close()
	h, err := New([]Route{{Prefix: "/api/x", Service: "x-service", Upstream: up.URL}}, "")
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodGet, "/api/x", nil)
	req.Header.Set("X-Internal-Caller", "orders-service")
	req.Header.Set("X-Internal-Token", "guess")
	h.ServeHTTP(httptest.NewRecorder(), req)
	if seen.Get("X-Internal-Caller") != "" || seen.Get("X-Internal-Token") != "" {
		t.Fatalf("the gateway forwarded the internal-caller headers: %v", seen)
	}
}
