---
id: chaos-engineering
title: Chaos engineering
family: reliability
currency: standard
role: choice
applies_when:
  - services >= 4
teaches: Prove the timeouts, breakers and compensations work by breaking things on purpose, in a test environment first.
---

# Chaos engineering

## What it is

Inject latency and failures (kill a participant mid-saga, delay the broker) and verify the system degrades as designed.

## Currency (2026)

**STANDARD.** Chaos Mesh / Litmus in Kubernetes; toxiproxy in tests.

## Fit signals (from `.cognidev/feature/context.md`)

- Four or more services with sagas and client edges.

## What it changes in the generated services

- toxiproxy-backed integration tests for each client; a saga test killing the orchestrator between steps.

## Go 2026 implementation

`github.com/Shopify/toxiproxy/v2/client`: `proxy.AddToxic("latency", "latency", "downstream", 1, toxiproxy.Attributes{"latency": 2000})`

## Anti-patterns

- Chaos in production before it passes in test.

## Interacts with

- [[resilience]] · [[saga-orchestration]]
