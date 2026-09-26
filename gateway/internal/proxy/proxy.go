// Package proxy routes path prefixes to upstream services.
package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sort"
	"strings"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// Route sends every request under Prefix to Upstream.
type Route struct {
	Prefix   string
	Service  string
	Upstream string
}

// FromEnv lets <SERVICE>_URL (e.g. ORDERS_SERVICE_URL) override an upstream.
func FromEnv(routes []Route) []Route {
	out := make([]Route, len(routes))
	for i, r := range routes {
		key := strings.ToUpper(strings.ReplaceAll(r.Service, "-", "_")) + "_URL"
		if v := os.Getenv(key); v != "" {
			r.Upstream = v
		}
		out[i] = r
	}
	return out
}

// New builds the routing handler. The longest matching prefix wins; with no
// match the request goes to legacy (the monolith) or, without one, 404s.
func New(routes []Route, legacy string) (http.Handler, error) {
	type target struct {
		prefix string
		proxy  *httputil.ReverseProxy
	}
	build := func(raw string) (*httputil.ReverseProxy, error) {
		u, err := url.Parse(raw)
		if err != nil || u.Host == "" {
			return nil, fmt.Errorf("bad upstream %q", raw)
		}
		p := httputil.NewSingleHostReverseProxy(u)
		p.Transport = otelhttp.NewTransport(http.DefaultTransport)
		return p, nil
	}
	var targets []target
	for _, r := range routes {
		p, err := build(r.Upstream)
		if err != nil {
			return nil, fmt.Errorf("route %s: %w", r.Prefix, err)
		}
		targets = append(targets, target{prefix: r.Prefix, proxy: p})
	}
	sort.Slice(targets, func(i, j int) bool { return len(targets[i].prefix) > len(targets[j].prefix) })
	var fallback http.Handler = http.NotFoundHandler()
	if legacy != "" {
		p, err := build(legacy)
		if err != nil {
			return nil, fmt.Errorf("legacy upstream: %w", err)
		}
		fallback = p
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, t := range targets {
			if r.URL.Path == t.prefix || strings.HasPrefix(r.URL.Path, strings.TrimSuffix(t.prefix, "/")+"/") {
				t.proxy.ServeHTTP(w, r)
				return
			}
		}
		fallback.ServeHTTP(w, r)
	}), nil
}
