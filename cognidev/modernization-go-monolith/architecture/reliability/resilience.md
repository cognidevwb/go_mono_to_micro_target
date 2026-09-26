---
id: resilience
title: Resilience — timeouts, retries, breakers
family: reliability
currency: standard
role: invariant
applies_when:
  - always
teaches: Every hop has a deadline, bounded retries with jitter for idempotent calls, and a breaker per upstream.
---

# Resilience — timeouts, retries, breakers

## What it is

The method calls that became network calls fail in new ways. Each typed client applies a timeout from the caller's context, retries only idempotent requests with exponential backoff + jitter, and opens a circuit breaker per upstream.

## Currency (2026)

**STANDARD — invariant.** `github.com/sony/gobreaker/v2`, `github.com/cenkalti/backoff/v5` or `hashicorp/go-retryablehttp`.

## Fit signals (from `.cognidev/feature/context.md`)

- Every client edge; outbound calls to external providers (the payments gateway).

## What it changes in the generated services

- Wrapped `http.RoundTripper` per client: otel → breaker → retry → transport.

## Go 2026 implementation

```go
cb := gobreaker.NewCircuitBreaker[*http.Response](gobreaker.Settings{Name: "catalog", Timeout: 30 * time.Second})
resp, err := cb.Execute(func() (*http.Response, error) { return hc.Do(req) })
```

## Anti-patterns

- Retrying POSTs without an idempotency key.
- Retries at every layer (gateway + client + mesh) multiplying load.

## Interacts with

- [[sync-transport]] · [[idempotency-key]] · [[rate-limiting]]
