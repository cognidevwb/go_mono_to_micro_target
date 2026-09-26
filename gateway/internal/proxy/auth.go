package proxy

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
)

// Authenticate returns middleware that verifies the caller's bearer token
// against issuer (OIDC discovery). With no issuer configured it passes every
// request through — local development only; every deployed environment sets
// AUTH_ISSUER, and each service re-validates the token itself.
func Authenticate(ctx context.Context, issuer, audience string) (func(http.Handler) http.Handler, error) {
	if issuer == "" {
		return func(h http.Handler) http.Handler { return h }, nil
	}
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, fmt.Errorf("oidc discovery %s: %w", issuer, err)
	}
	verifier := provider.Verifier(&oidc.Config{ClientID: audience, SkipClientIDCheck: audience == ""})
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
			if !ok {
				http.Error(w, "missing bearer token", http.StatusUnauthorized)
				return
			}
			if _, err := verifier.Verify(r.Context(), raw); err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}, nil
}
