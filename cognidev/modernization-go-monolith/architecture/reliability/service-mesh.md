---
id: service-mesh
title: Service mesh
family: reliability
currency: rising
role: choice
applies_when:
  - services >= 6
teaches: mTLS, retries and traffic policy in the platform — ambient/sidecarless meshes make it cheaper at fleet scale.
---

# Service mesh

## What it is

Istio ambient mode or Cilium service mesh provide mTLS, L7 policy and telemetry without per-pod sidecars.

## Currency (2026)

**RISING (sidecarless).**

## Fit signals (from `.cognidev/feature/context.md`)

- Six or more services; zero-trust or multi-team platform.

## What it changes in the generated services

- Namespace labels / policies; client-side retries reduced to avoid double retries.

## Go 2026 implementation

`kubectl label namespace shop istio.io/dataplane-mode=ambient`

## Anti-patterns

- A mesh for five services on one team.

## Interacts with

- [[service-identity]] · [[resilience]]
