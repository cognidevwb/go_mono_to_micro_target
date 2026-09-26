---
id: secrets-config
title: Secrets and configuration
family: observability
currency: standard
role: invariant
applies_when:
  - always
teaches: Configuration from the environment, secrets from the platform’s store via workload identity — never in git, never in the image.
---

# Secrets and configuration

## What it is

Twelve-factor config: every setting an env var parsed into one struct at startup; secrets injected by the platform (Kubernetes Secrets/External Secrets, Azure Key Vault references in Container Apps via managed identity).

## Currency (2026)

**STANDARD — invariant.**

## Fit signals (from `.cognidev/feature/context.md`)

- Always; the monolith's connection strings and payment-provider keys (`externals.json` outbound).

## What it changes in the generated services

- `.env.example` listing every variable; Helm values/ACA secrets referencing Key Vault.

## Go 2026 implementation

```go
type Config struct{ DatabaseURL string `env:"DATABASE_URL,required"`; NATSURL string `env:"NATS_URL" envDefault:"nats://nats:4222"` }
// github.com/caarlos0/env/v11: cfg, err := env.ParseAs[Config]()
```

## Anti-patterns

- Secrets in compose files committed to git.
- Config read lazily from anywhere in the code.

## Interacts with

- [[security]] · [[chassis]]
