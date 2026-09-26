# Architecture catalog — Go monolith → microservices (2026)

This folder is the playbook's **architectural knowledge base**: the
microservices patterns actually practised in 2026, each as one Markdown card the
decomposition points to as data. `index-go-patterns` reads every card's front
matter, evaluates its `applies_when` conditions against the MEASURED map
(`.cognidev/feature/plan.json`) and the captured answers
(`.cognidev/intent.yml`), and writes `.cognidev/feature/patterns.json` — which
patterns this cut applies, the sentence (with the number) that selected each,
and why every other card was left out.

## The 2026 stance — right-size first, then decompose

"Microservices by default" is a named anti-pattern. A context is extracted
against a named force — team autonomy, independent scaling, fault isolation, a
divergent lifecycle — and everything else stays in a bounded core. Go makes the
modular step cheap: packages, `internal/` and the import graph already enforce
boundaries, and `go.work` lets a package graduate into its own module without
rewriting call sites.

## Card grammar (what the selector parses)

```yaml
---
id: transactional-outbox
title: Transactional outbox
family: data                       # method · structure · decomposition · communication · data · observability · reliability
currency: standard                 # standard · rising · declining · niche · reference
role: invariant                    # invariant · default · choice · reference
applies_when:                      # ANDed; `[]` = never selected, then why_not is required
  - answer:outbox != none
  - answer:event_backbone != none
teaches: The one line a developer takes away.
---
```

Conditions: `always` · `answer:<key> = <value>` / `!= <value>` (an unanswered
key satisfies nothing) · `<fact> <op> <n>` over the measured facts `services`,
`sagas`, `calls`, `entities`, `reads`, `writes`, `read_heavy_services`. An
unknown key or an unparseable condition fails the phase — a typo must not
silently un-select a pattern.

Answer keys the cards read: `boundary_strategy` (domain · subdomain ·
technical), `migration_approach` (strangler · parallel-run · big-bang),
`service_internal_structure` (flat-package · hexagonal), `read_write_model`
(crud · cqrs), `event_backbone` (nats · kafka · none), `coordination`
(orchestration · choreography), `outbox` (… · none), `orchestrator`
(kubernetes · azure-containerapps · compose), `load_profile` (internal · steady
· bursty · high-volume).

Every card body has the same headings: **What it is · Currency (2026) · Fit
signals · What it changes in the generated services · Go 2026 implementation ·
Anti-patterns · Interacts with**.

## Gate 0 — right-sizing

| Option | When | What it changes |
| --- | --- | --- |
| [Modular monolith](decomposition/modular-monolith-first.md) | ≤2 data-owning contexts, no named force | keep one binary; depguard-enforced package boundaries |
| [Selective extraction](decomposition/right-sizing.md) **(default)** | a few contexts carry a force | strangler extraction of the justified services (`pilot_service`) |
| Full decomposition | multi-team, scaling, isolation everywhere | every data-owning context becomes a module in `go.work` |

## Gate 1 — internal service structure

| Option | When | What it changes |
| --- | --- | --- |
| [Flat package](structure/flat-package.md) **(default)** | right-sized service; monolith already package-per-context | ported package kept verbatim under `internal/<pkg>`; depguard rules |
| [Hexagonal](structure/hexagonal.md) | large, long-lived, infra-heavy service | `internal/domain` (stdlib only) + `app` + `adapters/*` |
| [Styles not applied](structure/styles-not-applied.md) | — | the omissions, with reasons |

## The families

- **Method** — [target architecture](_stack/target-architecture-2026.md) ·
  [conventions](_stack/microservices-conventions.md) ·
  [how the cut is computed](_stack/decomposition-method.md)
- **Decomposition** — business capability · subdomain · database per service ·
  strangler fig · branch by abstraction · anti-corruption layer · modular
  monolith first · parallel run · right-sizing
- **Communication** — API gateway · BFF (reference) · aggregation · sync
  transport · event-driven · saga orchestration · saga choreography · sync
  compensated chain · service discovery · contract testing · Dapr (reference)
- **Data** — transactional outbox · relay topology · idempotent consumer ·
  idempotency key · saga isolation · optimistic concurrency · expand–contract ·
  domain events · API composition · CQRS · materialized view · event sourcing ·
  CDC · connection-pool budget · read replica · read-through cache · TCC
  (reference) · durable execution (reference)
- **Observability & security** — chassis · slog + OpenTelemetry · health
  checks · JWT security · secrets & config · service identity · testing
- **Reliability** — resilience · rate limiting · deployment · autoscaling ·
  progressive delivery · chaos engineering · service mesh · sidecar (reference)

## How Go realizes microservices in 2026 (canonical library per concern)

Toolchain → **Go 1.26**, `go.work` multi-module · HTTP → **net/http** (+ the
monolith's gin/echo, or **chi v5**) · data → **pgx v5 + sqlc** or **gorm** ·
migrations → **golang-migrate** / **goose** · events → **nats.go jetstream**
(default) or **franz-go** (Kafka) · outbox → an owned table + relay (or
**Watermill SQL**) · logging → **log/slog** · telemetry → **OpenTelemetry Go**
(OTLP) · resilience → **sony/gobreaker v2** + **cenkalti/backoff v5** ·
auth → **golang-jwt/jwt v5** / **coreos/go-oidc v3** · config →
**caarlos0/env v11** · tests → **testcontainers-go** + `httptest` · CI →
**golangci-lint**, **govulncheck**, **Syft** + **Grype** (never
aquasecurity/trivy-action) · images → **distroless static nonroot** · deploy →
compose · Kubernetes + Helm · Azure Container Apps.

## Invariants no pattern may break

1. **Database per service.** No shared schema, no cross-service JOIN, no second
   `*gorm.DB` on another service's tables.
2. **Cross-service writes are sagas over an outbox**, never a transaction handle
   passed across a boundary; consumers are idempotent.
3. **No service imports another service's module.** Only `pkg/contracts`.
4. **Every hop has a deadline** — `context.Context` first, no `http.DefaultClient`.
5. **No package-level mutable state** carried into a replicated service.
6. **Right-size before you split.**
