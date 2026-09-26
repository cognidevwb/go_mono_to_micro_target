---
id: testing
title: Integration testing with real dependencies
family: observability
currency: standard
role: invariant
applies_when:
  - always
teaches: Test against real Postgres and NATS in containers, not mocks — the split moved the risk into the edges mocks hide.
---

# Integration testing with real dependencies

## What it is

Each service's tests run its handlers against real dependencies started by testcontainers-go; clients are tested against `httptest.Server`; the before/after test counts are compared with the monolith baseline.

## Currency (2026)

**STANDARD — invariant.** testcontainers-go modules for postgres and nats.

## Fit signals (from `.cognidev/feature/context.md`)

- Always; `run_tests` answer decides whether they execute in verify.

## What it changes in the generated services

- `internal/<pkg>/<pkg>_test.go` per service; `-race` in CI.

## Go 2026 implementation

```go
pg, err := postgres.Run(ctx, "postgres:17-alpine", postgres.WithDatabase("orders"), postgres.BasicWaitStrategies())
t.Cleanup(func() { _ = testcontainers.TerminateContainer(pg) })
```

## Anti-patterns

- sqlite as a stand-in for Postgres.
- Deleting the monolith's tests because they no longer compile.

## Interacts with

- [[contract-testing]] · [[chaos-engineering]]
