---
id: rate-limiting
title: Rate limiting
family: reliability
currency: standard
role: default
applies_when:
  - always
teaches: Limit at the edge per client, and inside services on expensive endpoints — shed load before the database does.
---

# Rate limiting

## What it is

Token-bucket limits at the gateway per client/key, plus per-endpoint limits on expensive routes.

## Currency (2026)

**STANDARD.** `golang.org/x/time/rate`.

## Fit signals (from `.cognidev/feature/context.md`)

- Always at the edge.

## What it changes in the generated services

- Gateway middleware; 429 with `Retry-After`.

## Go 2026 implementation

```go
lim := rate.NewLimiter(rate.Limit(50), 100)
if !lim.Allow() { w.Header().Set("Retry-After", "1"); http.Error(w, "slow down", http.StatusTooManyRequests); return }
```

## Anti-patterns

- Per-replica limits that multiply with the replica count.

## Interacts with

- [[api-gateway]] · [[resilience]]
