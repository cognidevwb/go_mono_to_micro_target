---
id: sync-transport
title: Synchronous transport — REST at the edge, gRPC inside
family: communication
currency: standard
role: invariant
applies_when:
  - always
teaches: Every cross-service read is a typed net/http client with a context deadline; gRPC is ported where the monolith already speaks it, not imposed.
---

# Synchronous transport — REST at the edge, gRPC inside

## What it is

How services call each other synchronously. The generated default is JSON over HTTP through a typed client per upstream context (`internal/clients`), each call bounded by the caller's context. gRPC (google.golang.org/grpc or Connect `connectrpc.com/connect`) is a valid internal transport and is carried across when the monolith already has `.proto` services.

## Currency (2026)

**STANDARD** (REST); gRPC/Connect **RISING** internally.

## Fit signals (from `.cognidev/feature/context.md`)

- Every edge in `plan.json.calls`.
- Signatures that are not portable (a `*gorm.DB`, a channel, a `func` argument) — `decomposition-facts.json` marks those edges `portable: false`.

## What it changes in the generated services

- One `<to>_client.go` per edge with only the methods the consumer calls; base URL from env.
- Non-portable edges get `// CW-SEAM[kind=signature-change]` — the call cannot move as written.

## Go 2026 implementation

```go
type CatalogClient struct{ base string; hc *http.Client }
func NewCatalogClient(base string) *CatalogClient {
    return &CatalogClient{base: base, hc: &http.Client{Timeout: 3 * time.Second, Transport: otelhttp.NewTransport(http.DefaultTransport)}}
}
```

## Anti-patterns

- `http.DefaultClient` (no timeout).
- Passing a transaction handle or channel across the boundary.
- A synchronous chain three services deep on a hot route.

## Interacts with

- [[resilience]] · [[anti-corruption-layer]] · [[contract-testing]] · [[service-discovery]]
