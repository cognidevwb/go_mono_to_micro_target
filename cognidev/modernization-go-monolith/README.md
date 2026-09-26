# modernization-go-monolith — Decompose a Go monolith into microservices

The Go twin of `modernization-dotnet-monolith` (v1.1.0), phase for phase. It reads a
gin / echo / chi + gorm / sqlx / pgx monolith through the IDE's Go structural shards,
maps it deterministically onto bounded-context services, lets a model refine the cut
one checked decision at a time, shows a briefing and a task table for approval, then
scaffolds a Go 1.26 `go.work` platform, MOVES the monolith's packages into it verbatim,
has the develop loop fill only the marked seams under a `go build` + `go vet` gate, and
closes with build/test proof, measured gates, a scorecard and a computed sign-off.

Version 1.0.0 · intent `modernize` · archetype `transformer` · transform_shape `replace`.

## Phases (35)

| # | id | tool | notes |
|---|----|------|-------|
| 1 | playbook-baseline | playbook-integrity record | |
| 2 | seed-context | seed-context | `onboarding/` → `cognidev/` |
| 3 | capture | questionnaire | `questionnaire/` 00–07 |
| 4 | validate-requirement | validate-requirement | |
| 5 | knowledge | index-go-knowledge --emit | `knowledge/*.md` → knowledge.json |
| 6 | understand | structural (`go-structural`) | shards under `.cognidev/understand/` |
| 7 | setup-environment | go-setup-env | **gate** |
| 8 | feasibility | go-feasibility (G3) | never fails the run |
| 9 | cicd-facts | scan-cicd-go (G3) | cicd-facts.json |
| 10 | map-services | map-services-go | service-map / plan / context / decomposition-facts |
| 11 | plan-decomposition | decompose-go (tool-loop, optional) | prompts in `context/decompose/` |
| 12 | reconcile-map | map-services-go --reconcile | |
| 13 | expand-plan | expand-plan-go | plan.json + decomposition.json |
| 14 | patterns | index-go-patterns | `architecture/` catalog → patterns.json |
| 15 | briefing | brief-go-monolith --tables-only | **review gate** |
| 16 | plan-tasks | plan-tasks | **review gate**, editable + selectable (pilot-then-continue) |
| 17 | publish-briefing | brief-go-monolith | cognidev/docs/ + DECOMPOSITION-BRIEFING.html |
| 18 | baseline | go-verify --baseline | go-baseline.json |
| 19 | scaffold-services | scaffold-services-go (G2) | |
| 20 | verify-build-inheritance | verify-cicd-go (G3) | |
| 21 | plan-apply | feature-plan | |
| 22 | port-monolith | port-monolith-go (G2) | port-report.json |
| 23 | contracts-manifest | contracts-go gen (G2) | contracts.lock.json |
| 24 | develop | execute-plan (tool-loop) | `CW_TYPECHECK_CMD` = go-typecheck |
| 25 | contracts-check | contracts-go check (G2) | contracts-drift.json |
| 26 | record-status | brief-go-monolith --status | MIGRATION-LOG.md |
| 27 | verify | go-verify | verify.json |
| 28 | understand-delivered | structural | |
| 29 | cve | go-cve | cve.json + cve-findings.json |
| 30 | gates | go-gates (G3) | gates.json; executes `gates/validate.yaml`; repair hand-back to develop |
| 31 | scorecard | go-scorecard (G3) | scorecard.json + MIGRATION-SCORECARD.md |
| 32 | closing-review | closing-review | NEXT-STEPS.md |
| 33 | playbook-integrity | playbook-integrity verify | |
| 34 | signoff | signoff-report | `signoff.yaml` → VALIDATION-REPORT.html |
| 35 | finalize | cutover-go (G2) | executes `git_strategy` |

## What changed from the .NET playbook, and why

- **`scan-runtimes` (scan-dotnet-runtimes) — DROPPED.** It exists because a classic
  net48 monolith and its net10 services need two different SDKs, so one machine can build
  the after-state and not the before-state. Go has no such split: the Go 1 compatibility
  promise means one current toolchain builds both sides, `go-setup-env` already resolves the
  highest `go`/`toolchain` directive across every go.mod (and installs it user-locally when
  missing), and `GOTOOLCHAIN` fetches a newer toolchain on demand. The inventory the scan
  would write is exactly what `.cognidev/env.json` records. The sign-off checks
  `runtimes-scanned` and `no-unplaceable-frameworks` are dropped with it; `runtime-target-chosen`
  reads `env.json`'s `go_version`.
- **Aspire AppHost / ServiceDefaults → `go.work` + `pkg/`.** Go has no distributed-app model
  to generate; the workspace file plays the "every service builds together" role and `pkg/`
  holds the contract DTOs (kept minimal — no behaviour).
- **YARP → the `gateway` answer** (generated Go reverse proxy, Envoy or Traefik).
- **Wolverine/MassTransit → the `messaging_lib` answer** (Watermill or the native NATS /
  Kafka client) plus a generated SQL outbox.
- **`nuget-cve` → `go-cve`** (govulncheck when installed, OSV otherwise).
- **Questionnaire:** `dotnet_version/service_framework/aot` → `go_version/web_framework`
  (keep the monolith's router by default — a port is a move); `orm` → `data_access`
  (keep gorm by default); `package_source` → `goproxy`; `aspire_deploy` dropped (no
  Aspire); added `gateway` and `data_ownership` (database- vs schema-per-service vs a
  transitional shared DB); `base_namespace` → `base_module`; `service_internal_structure`
  offers the ported flat package or hexagonal. `pilot_service` and the table's
  `selectable` tick work exactly as in .NET (pilot-then-continue).

## Target layout (what expand-plan-go plans; scaffold-services-go / port-monolith-go write)

```
go.work                                        go 1.26; ./gateway ./pkg ./services/<ctx> …
services/<ctx>/go.mod                          module <base_module>/services/<ctx>
services/<ctx>/cmd/<ctx>/main.go               entrypoint (scaffold)
services/<ctx>/internal/<pkg>/*.go             PORTED monolith package, verbatim, package name kept (cmd/ is not ported)
services/<ctx>/internal/<pkg>/<pkg>_test.go    tests
services/<ctx>/internal/store/db.go            one *gorm.DB / pgxpool per service (+ query.sql for pgx-sqlc)
services/<ctx>/internal/httpapi/<agg>_handlers.go
services/<ctx>/internal/clients/<to>_client.go planned client (replaced by the port's per-package client below) + _contract_test.go
services/<ctx>/internal/events/<ev>.go         publisher; <ev>_consumer.go on the subscribing side
services/<ctx>/internal/outbox/outbox.go       transactional outbox (event_backbone != none)
services/<ctx>/internal/saga/<saga>.go         orchestrator (+ _test.go); <saga>_coordinator.go with no broker
services/<ctx>/internal/jobs/scheduler.go      inherited background jobs (e.g. the restock ticker)
services/<ctx>/internal/acl/legacy.go          anti-corruption layer to the shrinking monolith
services/<ctx>/migrations/0001_init.sql        owned tables only
services/<ctx>/test/equivalence/<ctx>_equivalence_test.go
services/<ctx>/deploy/{deploy.targets.yml,strangler.values.yaml}
services/<ctx>/internal/clients/<provider>/client.go   typed client (one per provider package called)
services/<ctx>/internal/contracts/<provider>/      the consumer's OWN copy of the provider's DTOs
pkg/                                           minimal shared kernel: httpx, otelx, events, outbox, saga (no DTOs)
gateway/                                       reverse proxy module
deploy/                                        compose / k8s / helm / Azure Container Apps
```
CQRS adds `internal/app/<agg>_{commands,queries}.go` + `internal/readmodel/<agg>_view.go`;
hexagonal moves httpapi/clients/store/outbox/readmodel/events/acl under `internal/adapters/…`
and adds `internal/ports/ports.go`; choreography gives each participant
`internal/events/<saga>_reaction.go`.
`product_dirs: [services, gateway, pkg, deploy]`. `base_module` = the answer, else the
monolith's go.mod `module` line, else `example.com/platform`.

## Ownership

- G1: this directory's manifest, questionnaire, rules, context, onboarding, learn, config,
  signoff, knowledge, architecture; tools `map-services-go`, `decompose-go`, `expand-plan-go`,
  `brief-go-monolith`, `index-go-knowledge`, `index-go-patterns`.
- G2: `scaffold-services-go`, `port-monolith-go`, `contracts-go`, `cutover-go`, `scaffolds/**`.
- G3: `go-feasibility`, `go-gates`, `go-scorecard`, `scan-cicd-go`, `verify-cicd-go`, `gates/**`
  (the gates phase reads `gates/validate.yaml`, the same file name as .NET). The sign-off keys
  on these `go-gates` scalars: `cross_context_calls_converted`, `clients_wired_in_main`,
  `sagas_at_commit_points`, `no_residual_seam_markers`, `no_service_references_the_monolith`,
  `no_cross_service_module_import`, `service_names_are_domains_not_folders`,
  `every_service_owns_its_database`, `every_planned_service_is_in_the_workspace`,
  `go_version_is_the_chosen_one`, `gateway_routes_cover_every_service`,
  `no_monolith_import_left_in_services`, `no_wholesale_warning_suppression`, `tests_came_over`,
  `sagas_have_compensations`, `services_build_and_vet`, `services_tests_pass`, `health_endpoints_present`,
  `containers_distroless_nonroot`, `shared_state_has_an_owner`, `events_use_outbox`, `contracts_no_drift`,
  `declarative_readiness_gates` (22 gates), `known_vulnerable_dependencies`. A skipped gate publishes no scalar → NOT VERIFIED.

## Routing

`playbook-map.yaml` `modernize:` — go, golang, gin, echo, chi, fiber, gorilla →
this playbook. `modernization-map.yaml` — card `go-monolith` (go.mod + a web framework +
a persistence signal; SI-first `language: [go]`, `min_entities: 2`).
