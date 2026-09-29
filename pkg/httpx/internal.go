package httpx

import (
	"crypto/subtle"
	"net/http"
	"os"
)

// The two headers a service-to-service call carries. The monolith made these
// calls in process, where no user check applied; the services keep the
// monolith's user check on their routes, so a call from another service of
// this platform has to say so. The gateway strips both headers from every
// request it forwards, so no outside caller can claim to be a service.
const (
	CallerHeader = "X-Internal-Caller"
	TokenHeader  = "X-Internal-Token"
)

// internalToken is the shared secret services present to each other
// (INTERNAL_TOKEN). Unset, the caller header alone is trusted — the services
// are reachable only through the gateway or the cluster network.
func internalToken() string { return os.Getenv("INTERNAL_TOKEN") }

// IsInternal reports whether r is a call from another service of this
// platform.
func IsInternal(r *http.Request) bool {
	if r.Header.Get(CallerHeader) == "" {
		return false
	}
	want := internalToken()
	if want == "" {
		return true
	}
	return subtle.ConstantTimeCompare([]byte(r.Header.Get(TokenHeader)), []byte(want)) == 1
}

// UnlessInternal applies mw (a user check ported from the monolith) to every
// request except a service-to-service call, which goes straight to next.
func UnlessInternal(mw func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		checked := mw(next)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if IsInternal(r) {
				next.ServeHTTP(w, r)
				return
			}
			checked.ServeHTTP(w, r)
		})
	}
}

// internalCaller marks every request a service sends to another service.
type internalCaller struct{ next http.RoundTripper }

func (t internalCaller) RoundTrip(r *http.Request) (*http.Response, error) {
	r = r.Clone(r.Context())
	name := os.Getenv("OTEL_SERVICE_NAME")
	if name == "" {
		name = "service"
	}
	r.Header.Set(CallerHeader, name)
	if tok := internalToken(); tok != "" {
		r.Header.Set(TokenHeader, tok)
	}
	return t.next.RoundTrip(r)
}
