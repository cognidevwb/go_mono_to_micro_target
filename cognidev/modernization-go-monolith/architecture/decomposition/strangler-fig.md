---
id: strangler-fig
title: Strangler fig
family: decomposition
currency: standard
role: invariant
applies_when:
  - always
teaches: Put a gateway in front of the monolith and move one route group at a time; the monolith shrinks while production keeps working.
---

# Strangler fig

## What it is

Incremental migration: a facade (the gateway) routes each request either to the monolith or to the extracted service. Contexts move one at a time, leaf services first, the orchestrating hub last; the monolith keeps serving everything not yet moved and is deleted only when it serves nothing.

## Currency (2026)

**STANDARD — invariant.** Fowler's pattern, Azure Architecture Center; the default for every non-trivial migration.

## Fit signals (from `.cognidev/feature/context.md`)

- The strangler order in `plan.json.order` (fewest outbound calls first).
- Route groups that map one-to-one onto a context (`http.json`).

## What it changes in the generated services

- `gateway/` routes `/v1/products/*` → catalog-service etc.; unmatched → the monolith.
- Extraction order = `plan.json.order`; the pilot answer (`pilot_service`) carves a subset first.

## Go 2026 implementation

```go
// gateway/main.go
routes := map[string]*url.URL{"/v1/products": catalogURL, "/v1/orders": ordersURL}
fallback := httputil.NewSingleHostReverseProxy(monolithURL)
```
Or Envoy / Traefik with the same route table.

## Anti-patterns

- Big-bang switch of all routes at once.
- Extracting the hub (orders) first — every call it makes becomes a hop before its callees exist.
- Leaving the monolith writing to tables a service now owns.

## Interacts with

- [[api-gateway]] · [[branch-by-abstraction]] · [[anti-corruption-layer]] · [[progressive-delivery]]
