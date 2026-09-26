---
id: context-deadline-exceeded-on-the-typed-client
runtime: go
since: 0
category: runtime-failure
tags: [http-client, timeouts, service-discovery, resilience, typed-client]
symptoms:
  - "context deadline exceeded"
  - "Client.Timeout exceeded while awaiting headers"
  - "no such host"
  - "connection refused"
remedy: { kind: environment }
reference: "net/http Client.Timeout and context cancellation (pkg.go.dev/net/http#Client)"
---

# The method call that became a network call has no deadline budget

## The signature

```
Get "http://catalog-service:8080/v1/products/42/price": context deadline exceeded (Client.Timeout exceeded while awaiting headers)
```

or, before the service is even reachable:

```
Get "http://catalog-service:8080/v1/products/42/price": dial tcp: lookup catalog-service on 127.0.0.11:53: no such host
```

## What it means

`s.catalog.PriceOf(id)` took nanoseconds and could not fail. As a typed client
it is a network hop with a timeout, a DNS lookup and a peer that can be down.
Inside a loop over order lines (PlaceOrder calls `PriceOf` once per item) N
sequential hops multiply the latency, and a single slow peer exhausts the
caller's whole request budget.

`no such host` is service discovery: the client uses a name the environment
does not define (compose service name vs Kubernetes Service name vs Container
Apps app name).

## What to do

- Every client takes the caller's `context.Context` and sets a per-call timeout
  below the inbound request's; never `http.DefaultClient`.
- Retries only on idempotent reads, with backoff and a small cap; a breaker
  (`sony/gobreaker`) per upstream.
- Batch the per-line call (`GET /v1/prices?ids=…`) where the loop exists.
- Base URLs come from env (`CATALOG_URL`), named identically across compose, k8s
  and ACA.
