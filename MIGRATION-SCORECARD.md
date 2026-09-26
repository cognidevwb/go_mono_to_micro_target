# Decomposition scorecard

Each stage of the cut graded on three axes. `—` means the axis was not
measured, and it never reads as a pass: a stage nothing looked at is
reported that way rather than counted as clean.

## Roll-up

| KPI | Result | Target | |
|---|---|---|---|
| Build | green | green | PASS |
| Services in go.work | 5 of 5 | all | PASS |
| Services independent of the monolith | yes | all | PASS |
| Seams still to fill | 0 | 0 | PASS |
| Tests preserved | not measured | no loss | — |
| Gates not evaluated | 2 of 25 | 0 | WATCH |
| Categories failing | 0 of 12 | 0 | PASS |

- **Build** — verify.json: build_ok
- **Services in go.work** — all 5 planned service(s) are modules in go.work
- **Services independent of the monolith** — none of 5 service module(s) depend on anything outside the generated tree
- **Seams still to fill** — no marker left across 184 generated file(s)
- **Tests preserved** — the before-state or the verify step recorded no test count
- **Gates not evaluated** — gates.json: a gate that could not read its input is counted here, never as a pass
- **Categories failing** — 0 category/categories had nothing to measure

## By stage

| Stage | Completeness | Correctness | Restraint |
|---|---|---|---|
| Service boundaries | PASS | PASS | WATCH |
| Service independence | PASS | PASS | — |
| Code movement | PASS | PASS | — |
| Cross-service calls | PASS | PASS | — |
| Transactions | PASS | PASS | — |
| Data ownership | PASS | PASS | — |
| Seams finished | PASS | PASS | — |
| Platform wiring | PASS | PASS | — |
| Runtime readiness | PASS | PASS | PASS |
| Events and contracts | PASS | PASS | — |
| Tests | PASS | PASS | — |
| Dependencies and diagnostics | — | — | PASS |

### Service boundaries

- **Completeness** (PASS) — all 5 planned service(s) are modules in go.work
- **Correctness** (PASS) — all 5 service name(s) describe a domain
- **Restraint** (WATCH) — 8 of 627 changed file(s) were not proposed by the plan (619 were): .dockerignore, .editorconfig, .githooks/pre-commit, contracts/events/envelope.schema.json, contracts/events/order.placed.schema.json

A boundary go.work does not build is a directory. A boundary named after a source folder is the monolith's tree with a prefix.

### Service independence

- **Completeness** (PASS) — none of 5 service module(s) depend on anything outside the generated tree
- **Correctness** (PASS) — none of 116 service file(s) or 5 go.mod(s) reach into another service
- **Restraint** (—) — independence is a property of the result, not of the diff

The decisive row. A service whose go.mod still requires the monolith cannot be built, versioned or deployed without it.

### Code movement

- **Completeness** (PASS) — 46 file(s) written across 5 service(s); 0 construct(s) were deferred rather than guessed at
- **Correctness** (PASS) — 97% of the moved lines were placed verbatim or generated; 12 seam point(s) were left for the develop loop
- **Restraint** (—) — the port writes only into the generated tree

### Cross-service calls

- **Completeness** (PASS) — all 2 cross-context call edge(s) resolve to a client that performs an HTTP call
- **Correctness** (PASS) — all 4 generated client(s) are constructed at a composition root
- **Restraint** (—) — a client is generated, not edited

A call that crossed a package boundary in the monolith is an HTTP call now, or it is nothing. A client nothing constructs is a nil pointer at the first request.

### Transactions

- **Completeness** (PASS) — all 1 declared saga(s) have an orchestrator that references every participant
- **Correctness** (PASS) — all 2 remote saga step(s) have a compensation naming their service
- **Restraint** (—) — a saga is written, not edited

One db.Transaction across four packages is not one transaction across four services. A saga that cannot undo a step commits halfway.

### Data ownership

- **Completeness** (PASS) — each of 5 data-owning service(s) opens its own database; 9 table(s), one owner each
- **Correctness** (PASS) — 2 piece(s) of carried state, each in its one listed owner
- **Restraint** (—) — ownership is a property of the result, not of the diff

Tables and package-level state both have exactly one owner after the cut, or they diverge silently.

### Seams finished

- **Completeness** (PASS) — no marker left across 184 generated file(s)
- **Correctness** (PASS) — all 2 cross-context call edge(s) resolve to a client that performs an HTTP call
- **Restraint** (—) — a seam is filled in place

Completeness is whether the markers are gone. Correctness is whether the code that replaced them does the thing — deleting a marker is the cheaper way to make the first one pass.

### Platform wiring

- **Completeness** (PASS) — all 5 service(s) have a route and an upstream
- **Correctness** (PASS) — all 8 generated module file(s) target go 1.26
- **Restraint** (—) — the platform modules are generated

### Runtime readiness

- **Completeness** (PASS) — all 5 service(s) expose liveness and readiness
- **Correctness** (PASS) — all 5 service image(s) are distroless and non-root
- **Restraint** (PASS) — all 5 service(s) build from their own module with no vet finding

Restraint here is the build itself: a service that only compiles inside go.work, beside its siblings, was never cut.

### Events and contracts

- **Completeness** (PASS) — all 1 publishing service(s) have a transactional outbox
- **Correctness** (PASS) — every provider contract still has the shape its consumers were built against
- **Restraint** (—) — a contract is recorded, not edited

An in-process bus never lost a message between commit and publish; a broker does, unless an outbox sits in between.

### Tests

- **Completeness** (PASS) — all 5 service(s) on disk have a test function
- **Correctness** (PASS) — the tests of all 5 service(s) ran and passed (59 test(s))
- **Restraint** (—) — tests are added, not rewritten, by this playbook

### Dependencies and diagnostics

- **Completeness** (—) — not evaluated — the advisory scan did not complete: the advisory database was not queried — this run is offline by default; enable the registry lookup to scan
- **Correctness** (—) — not evaluated — the advisory scan did not complete: the advisory database was not queried — this run is offline by default; enable the registry lookup to scan
- **Restraint** (PASS) — no wholesale suppression across 184 generated file(s)

A //nolint added to make the port pass is a restraint failure, not a correctness one: it changed a setting the migration never required.

## Restraint — what changed outside the new services

Diffed against `source`: 627 file(s) changed, +79642 −28 lines. 619 were named by the plan or are the run's own output; 8 were not.

Changed without being proposed:

- `.dockerignore`
- `.editorconfig`
- `.githooks/pre-commit`
- `contracts/events/envelope.schema.json`
- `contracts/events/order.placed.schema.json`
- `contracts/openapi/catalog-service.json`
- `contracts/openapi/customers-service.json`
- `contracts/openapi/orders-service.json`

