---
id: health-checks
title: Health checks
family: observability
currency: standard
role: invariant
applies_when:
  - always
teaches: /healthz says the process lives; /readyz says it can serve — keep them different or the orchestrator restarts healthy pods.
---

# Health checks

## What it is

Liveness answers only whether the process should be restarted; readiness checks dependencies it cannot serve without (database ping, broker connection).

## Currency (2026)

**STANDARD — invariant.**

## Fit signals (from `.cognidev/feature/context.md`)

- Always.

## What it changes in the generated services

- `/healthz` and `/readyz` on every service and the gateway; probes in k8s/ACA manifests.

## Go 2026 implementation

```go
mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), time.Second); defer cancel()
    if pool.Ping(ctx) != nil { http.Error(w, "db", http.StatusServiceUnavailable); return }
})
```

## Anti-patterns

- Checking downstream services in liveness (cascading restarts).

## Interacts with

- [[chassis]] · [[deployment]]
