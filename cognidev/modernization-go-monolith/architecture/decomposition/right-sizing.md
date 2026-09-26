---
id: right-sizing
title: Right-sizing
family: decomposition
currency: rising
role: invariant
applies_when:
  - always
teaches: Extract a service only against a named force — team, scale, isolation or lifecycle; everything else stays in a bounded core.
---

# Right-sizing

## What it is

The Gate-0 judgement. A context becomes a service because something specific demands it — a team that ships on its own cadence, a hot path that scales alone, a failure domain to isolate, a divergent lifecycle. Without one, the context stays in a bounded modular core. Also: a self-contained service answers its requests without a synchronous chain of other services.

## Currency (2026)

**RISING.** Microservice-envy is a named anti-pattern in the 2025–2026 radar; cost and cognitive load drive consolidation.

## Fit signals (from `.cognidev/feature/context.md`)

- `driver` answer; coupling weights between candidates; sync call chains in hot routes.

## What it changes in the generated services

- The briefing reports per service the force that justifies it and the hops it adds.
- `pilot_service` lets a team carve one service, read it, then continue.

## Go 2026 implementation

No code — the measurement is the cut: services, calls and sagas in `plan.json`, and the per-endpoint reach in `decomposition-facts.json`.

## Anti-patterns

- Splitting every package because it is a package.
- A service that needs three synchronous calls to answer its main route.

## Interacts with

- [[modular-monolith-first]] · [[decompose-by-business-capability]] · [[api-composition]]
