package httpx

import (
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// NewClient returns an HTTP client for service-to-service calls: a bounded
// timeout and trace propagation on every request.
func NewClient() *http.Client {
	return &http.Client{
		Timeout:   5 * time.Second,
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
}
