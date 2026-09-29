package httpx

import (
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// NewClient returns an HTTP client for service-to-service calls: a bounded
// timeout, trace propagation, and the internal-caller headers the receiving
// service's user check lets through (see internal.go).
func NewClient() *http.Client {
	return &http.Client{
		Timeout:   5 * time.Second,
		Transport: internalCaller{otelhttp.NewTransport(http.DefaultTransport)},
	}
}
