---
id: idempotency-key
title: Idempotency key
family: data
currency: standard
role: invariant
applies_when:
  - services >= 2
teaches: The same guarantee at the HTTP edge — a caller-supplied key makes a retried write a no-op.
---

# Idempotency key

## What it is

Write endpoints accept an `Idempotency-Key` header; the first response is stored against it and replayed for a retry with the same key. Retries in typed clients are then safe for writes.

## Currency (2026)

**STANDARD — invariant** (IETF draft `Idempotency-Key` header; Stripe practice).

## Fit signals (from `.cognidev/feature/context.md`)

- Any write route called by another service or retried by the gateway.

## What it changes in the generated services

- Middleware + `idempotency_keys` table per service that exposes writes.

## Go 2026 implementation

```go
func Idempotent(store Store) func(http.Handler) http.Handler { /* key → cached status+body; 409 on in-flight duplicate */ }
```

## Anti-patterns

- Keys generated server-side (the client cannot retry with them).
- Replaying responses across different request bodies.

## Interacts with

- [[sync-compensated-chain]] · [[resilience]]
