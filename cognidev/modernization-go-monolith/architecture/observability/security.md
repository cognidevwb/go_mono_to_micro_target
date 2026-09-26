---
id: security
title: Security — JWT at the edge and in every service
family: observability
currency: standard
role: invariant
applies_when:
  - always
teaches: Validate the token at the gateway and again in each service; carry identity as a claim, not as an in-process context value.
---

# Security — JWT at the edge and in every service

## What it is

The gateway authenticates; every service re-validates the JWT (issuer, audience, expiry, signature via JWKS) and authorizes. The monolith's auth middleware set `platform.UserIDKey` in `context.Context` — after the cut that value does not cross the network; the user id travels in the token and each service extracts it itself.

## Currency (2026)

**STANDARD — invariant.**

## Fit signals (from `.cognidev/feature/context.md`)

- `ambient-context.json` reads of a request-scoped key (the fixture's `UserIDKey`) in files of several services.

## What it changes in the generated services

- JWT middleware per service (`github.com/golang-jwt/jwt/v5` + `github.com/MicahParks/keyfunc/v3`, or `github.com/coreos/go-oidc/v3`); the ambient key becomes a claim read by that middleware; seam `// CW-SEAM[kind=ambient-context]`.

## Go 2026 implementation

```go
tok, err := jwt.Parse(raw, jwks.Keyfunc, jwt.WithIssuer(iss), jwt.WithAudience(aud), jwt.WithExpirationRequired())
sub, _ := tok.Claims.GetSubject()
ctx = context.WithValue(r.Context(), userIDKey{}, sub)
```

## Anti-patterns

- Trusting a gateway-set `X-User-Id` header without mTLS.
- Sharing the monolith's context key type via `pkg/`.

## Interacts with

- [[api-gateway]] · [[service-identity]] · [[secrets-config]]
