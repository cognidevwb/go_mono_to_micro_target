---
id: api-gateway
title: API gateway
family: communication
currency: standard
role: invariant
applies_when:
  - always
teaches: One edge for every client — routing, auth, rate limits and the strangler route table live there, not in each service.
---

# API gateway

## What it is

A single entry point in front of the fleet: routes each path to its service (and, during migration, everything else to the monolith), validates the JWT, applies rate limits and CORS, and stamps a request id / trace context.

## Currency (2026)

**STANDARD — invariant.** Every reference architecture; in Go estates the edge is a small `net/http/httputil.ReverseProxy` service, or Envoy / Traefik when the platform team runs one.

## Fit signals (from `.cognidev/feature/context.md`)

- More than one service; external clients calling routes that now belong to different services (`http.json`).

## What it changes in the generated services

- `gateway/` module: route table generated from `plan.json` (one prefix per service + monolith fallback), JWT middleware, `/healthz`.
- Or Envoy/Traefik config when `gateway` = envoy/traefik.

## Go 2026 implementation

```go
proxy := &httputil.ReverseProxy{Rewrite: func(r *httputil.ProxyRequest) {
    r.SetURL(target(r.In.URL.Path)); r.SetXForwarded()
}}
mux.Handle("/", otelhttp.NewHandler(auth(rateLimit(proxy)), "gateway"))
```
`Rewrite` (Go 1.20+) replaces `Director` and strips hop-by-hop headers correctly.

## Anti-patterns

- Business logic or aggregation of many services in the gateway (that is [[aggregator]] or a BFF, deliberately).
- Services exposed directly to the internet alongside the gateway.

## Interacts with

- [[strangler-fig]] · [[security]] · [[rate-limiting]] · [[backend-for-frontend]]
