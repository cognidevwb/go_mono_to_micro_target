---
id: contract-testing
title: Contract testing
family: communication
currency: standard
role: default
applies_when:
  - services >= 2
teaches: Each consumer pins what it reads from a provider; the provider’s CI fails when it would break a consumer — independent deploys stay safe.
---

# Contract testing

## What it is

Consumer-driven contracts: each consumer records the requests it makes and the fields it reads; the provider verifies it still satisfies them before it ships.

## Currency (2026)

**STANDARD.** Pact (pact-go v2) and schema-snapshot approaches; `contracts-go` generates the lock/drift check this playbook uses.

## Fit signals (from `.cognidev/feature/context.md`)

- ≥2 services with client edges between them.

## What it changes in the generated services

- `contracts.lock.json` snapshot of each provider's DTOs; `contracts-drift.json` after develop.
- Optional pact-go tests per client in `internal/clients/*_contract_test.go`.

## Go 2026 implementation

```go
// pact-go v2 consumer test (shape)
mock.AddInteraction().Given("product 42 exists").
    UponReceiving("a price request").WithRequest("GET", "/v1/products/42/price").
    WillRespondWith(200, func(b *consumer.V2ResponseBuilder) { b.JSONBody(matchers.Like(map[string]any{"amount": 9.5})) })
```

## Anti-patterns

- End-to-end environments as the only safety net.
- Contracts generated from the provider (they test nothing a consumer needs).

## Interacts with

- [[anti-corruption-layer]] · [[testing]] · [[expand-contract]]
