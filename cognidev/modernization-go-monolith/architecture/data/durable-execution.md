---
id: durable-execution
title: Durable execution
family: data
currency: rising
role: reference
applies_when: []
why_not: This platform generates a saga whose state lives in the orchestrating service’s own database — no extra runtime to operate. Durable execution (Temporal Go SDK) is the upgrade path when flows outgrow it.
teaches: When sagas outgrow a table and a loop, a workflow engine persists every step for you — Temporal has first-class Go.
---

# Durable execution

## What it is

Workflow engines (Temporal, Restate, DBOS) persist each step and replay deterministically after a crash; the saga becomes ordinary Go code.

## Currency (2026)

**RISING.** Temporal's Go SDK is its reference SDK.

## Fit signals (from `.cognidev/feature/context.md`)

- Many long-running, multi-step flows with timers.

## What it changes in the generated services

Not generated — an extra cluster to operate.

## Go 2026 implementation

`workflow.ExecuteActivity(ctx, ReserveStock, lines).Get(ctx, &res)`

## Anti-patterns

- Adopting it for one three-step saga.

## Interacts with

- [[saga-orchestration]]
