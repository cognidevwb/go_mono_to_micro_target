---
id: read-through-cache
title: Read-through cache
family: data
currency: standard
role: choice
applies_when:
  - read_heavy_services >= 1
  - answer:load_profile != internal
teaches: Remove the read instead of moving it — and replace any package-level cache map, which is per-replica after the split.
---

# Read-through cache

## What it is

A cache in front of hot reads, with single-flight to stop stampedes. The monolith's package-level `priceCache` map guarded by a `sync.RWMutex` becomes per-replica and inconsistent after scaling out; a shared cache (Redis/Valkey) or explicit TTL + invalidation event replaces it.

## Currency (2026)

**STANDARD.**

## Fit signals (from `.cognidev/feature/context.md`)

- Read-heavy services; `datastructures.json` shared_mutable package vars; cache keys touched from two services.

## What it changes in the generated services

- `golang.org/x/sync/singleflight` + TTL cache, or Valkey via `github.com/redis/go-redis/v9`; invalidation on the owner's change event.

## Go 2026 implementation

```go
v, err, _ := s.sf.Do(key, func() (any, error) { return s.load(ctx, id) })
```

## Anti-patterns

- A package-level map as a cache in a replicated service.
- Caching another service's data with no invalidation contract.

## Interacts with

- [[materialized-view]] · [[read-replica]]
