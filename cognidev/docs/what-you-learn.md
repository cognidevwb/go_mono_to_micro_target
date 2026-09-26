# What you learn

This decomposition applies 40 of the 59 patterns in the catalog. Each row says
what the pattern teaches, why it applies HERE, and the document to read.

### Communication

| Pattern | What it teaches | Why it applies here | Read |
| --- | --- | --- | --- |
| **Gateway aggregation** (default) | When one screen needs three services, compose once at the edge with bounded fan-out — not three round-trips from the client. | 2 cross-context calls | `architecture/communication/aggregator.md` |
| **API gateway** (invariant) | One edge for every client — routing, auth, rate limits and the strangler route table live there, not in each service. | applies to every decomposition | `architecture/communication/api-gateway.md` |
| **Contract testing** (default) | Each consumer pins what it reads from a provider; the provider’s CI fails when it would break a consumer — independent deploys stay safe. | 5 services | `architecture/communication/contract-testing.md` |
| **Event-driven messaging** (default) | The in-process bus becomes a broker: every topic a context publishes becomes a durable subject, published through the outbox. | you chose event_backbone = nats, which is not none | `architecture/communication/event-driven.md` |
| **Saga — orchestration** (default) | One coordinator owns the multi-service write: each step a local transaction, each failure a compensating step, the riskiest step last. | you chose coordination = orchestration | `architecture/communication/saga-orchestration.md` |
| **Service discovery** (default) | Let the platform resolve names — compose service names, Kubernetes Services, Container Apps app names — and pass base URLs by environment. | applies to every decomposition | `architecture/communication/service-discovery.md` |
| **Synchronous transport — REST at the edge, gRPC inside** (invariant) | Every cross-service read is a typed net/http client with a context deadline; gRPC is ported where the monolith already speaks it, not imposed. | applies to every decomposition | `architecture/communication/sync-transport.md` |

### Data

| Pattern | What it teaches | Why it applies here | Read |
| --- | --- | --- | --- |
| **API composition** (default) | Answer a cross-service query by calling the owners and joining in memory — try this before building a read model. | 2 cross-context calls | `architecture/data/api-composition.md` |
| **Connection-pool budget** (invariant) | Pods × pool size must fit under the database’s max_connections — set it, don’t inherit a default. | you chose orchestrator = azure-containerapps, which is not compose | `architecture/data/connection-pool-budget.md` |
| **Domain events and aggregates** (invariant) | An aggregate is the consistency boundary; what it records as having happened is the event other services react to. | applies to every decomposition | `architecture/data/domain-event-aggregate.md` |
| **Expand–contract schema change** (invariant) | Two versions of the code always share the schema during a rollout — add, migrate, then remove, never rename in place. | applies to every decomposition | `architecture/data/expand-contract.md` |
| **Idempotency key** (invariant) | The same guarantee at the HTTP edge — a caller-supplied key makes a retried write a no-op. | 5 services | `architecture/data/idempotency-key.md` |
| **Idempotent consumer** (invariant) | At-least-once is reality — record each processed message id in the same transaction as its effect. | you chose event_backbone = nats, which is not none | `architecture/data/idempotent-consumer.md` |
| **Optimistic concurrency** (invariant) | A version column and a conditional UPDATE — the split widened every write window, and last-writer-wins now loses data. | applies to every decomposition | `architecture/data/optimistic-concurrency.md` |
| **Outbox relay topology** (invariant) | Decide who drains the outbox — one relay per service, leader-elected or SKIP LOCKED — and make sure every replica is allowed to run it. | you chose outbox = sql-outbox, which is not none | `architecture/data/outbox-relay-topology.md` |
| **Saga isolation countermeasures** (invariant) | A saga has no isolation — other requests see its intermediate states, so design semantic locks and reversible steps. | 1 saga candidate | `architecture/data/saga-isolation.md` |
| **Transactional outbox** (invariant) | Write the row and the message in one local transaction, relay after commit — the only honest way out of the dual write. | you chose outbox = sql-outbox, which is not none | `architecture/data/transactional-outbox.md` |

### Decomposition

| Pattern | What it teaches | Why it applies here | Read |
| --- | --- | --- | --- |
| **Anti-corruption layer** (invariant) | Translate at every boundary — the consumer maps the provider’s contract into its own types, so one service’s model never leaks into another. | 2 cross-context calls | `architecture/decomposition/anti-corruption-layer.md` |
| **Branch by abstraction** (choice) | Put an interface in front of the in-process call, swap the implementation for a client behind it — the in-code half of the strangler. | you chose migration_approach = strangler | `architecture/decomposition/branch-by-abstraction.md` |
| **Database per service** (invariant) | Each service owns its tables; nobody else reads them — the one invariant that makes independent deployment real. | applies to every decomposition | `architecture/decomposition/database-per-service.md` |
| **Decompose by business capability** (default) | A service is a thing the business does and owns data for — catalog, orders, payments — never a technical layer. | you chose boundary_strategy = domain | `architecture/decomposition/decompose-by-business-capability.md` |
| **Right-sizing** (invariant) | Extract a service only against a named force — team, scale, isolation or lifecycle; everything else stays in a bounded core. | applies to every decomposition | `architecture/decomposition/right-sizing.md` |
| **Strangler fig** (invariant) | Put a gateway in front of the monolith and move one route group at a time; the monolith shrinks while production keeps working. | applies to every decomposition | `architecture/decomposition/strangler-fig.md` |

### Method

| Pattern | What it teaches | Why it applies here | Read |
| --- | --- | --- | --- |
| **How the cut is computed** (reference) | Boundaries come from data ownership and co-writes measured in the shards; the model only refines names and edges it can be checked against. | applies to every decomposition | `architecture/_stack/decomposition-method.md` |
| **Go microservice conventions** (reference) | Every service looks the same when you open it — cmd/, internal/, migrations/, one go.mod — so a reader learns the fleet once. | applies to every decomposition | `architecture/_stack/microservices-conventions.md` |
| **The 2026 Go target architecture** (reference) | One pinned stack, stated once — Go 1.26, a module per service in one go.work, slog + OpenTelemetry, NATS JetStream, distroless nonroot images. | applies to every decomposition | `architecture/_stack/target-architecture-2026.md` |

### Observability

| Pattern | What it teaches | Why it applies here | Read |
| --- | --- | --- | --- |
| **Microservice chassis** (invariant) | One small shared startup path — config, slog, OTel, health, graceful shutdown — so every main.go is the same ten lines. | applies to every decomposition | `architecture/observability/chassis.md` |
| **Health checks** (invariant) | /healthz says the process lives; /readyz says it can serve — keep them different or the orchestrator restarts healthy pods. | applies to every decomposition | `architecture/observability/health-checks.md` |
| **Observability — slog + OpenTelemetry** (invariant) | Traces, metrics and logs over one OTLP pipe; the trace id in every log line; the backend is a swap. | applies to every decomposition | `architecture/observability/observability.md` |
| **Secrets and configuration** (invariant) | Configuration from the environment, secrets from the platform’s store via workload identity — never in git, never in the image. | applies to every decomposition | `architecture/observability/secrets-config.md` |
| **Security — JWT at the edge and in every service** (invariant) | Validate the token at the gateway and again in each service; carry identity as a claim, not as an in-process context value. | applies to every decomposition | `architecture/observability/security.md` |
| **Integration testing with real dependencies** (invariant) | Test against real Postgres and NATS in containers, not mocks — the split moved the risk into the edges mocks hide. | applies to every decomposition | `architecture/observability/testing.md` |

### Reliability

| Pattern | What it teaches | Why it applies here | Read |
| --- | --- | --- | --- |
| **Autoscaling — HPA and KEDA** (choice) | Scale HTTP services on CPU/RPS and consumers on queue depth — CPU is the wrong signal for a JetStream consumer. | you chose orchestrator = azure-containerapps, which is not compose | `architecture/reliability/autoscaling.md` |
| **Chaos engineering** (choice) | Prove the timeouts, breakers and compensations work by breaking things on purpose, in a test environment first. | 5 services | `architecture/reliability/chaos-engineering.md` |
| **Deployment — compose, Kubernetes, Azure Container Apps** (default) | One static binary per service in a distroless nonroot image, deployed the same way everywhere — compose locally, Kubernetes/Helm or Container Apps in production. | applies to every decomposition | `architecture/reliability/deployment.md` |
| **Progressive delivery** (choice) | Move traffic to an extracted service gradually — weighted routes and flags, with a fast way back. | you chose migration_approach = strangler | `architecture/reliability/progressive-delivery.md` |
| **Rate limiting** (default) | Limit at the edge per client, and inside services on expensive endpoints — shed load before the database does. | applies to every decomposition | `architecture/reliability/rate-limiting.md` |
| **Resilience — timeouts, retries, breakers** (invariant) | Every hop has a deadline, bounded retries with jitter for idempotent calls, and a breaker per upstream. | applies to every decomposition | `architecture/reliability/resilience.md` |

### Structure

| Pattern | What it teaches | Why it applies here | Read |
| --- | --- | --- | --- |
| **Flat package per context** (default) | Keep the ported package as the unit — one package per context inside the service, wiring in cmd/, adapters beside it; Go rewards fewer layers. | you chose service_internal_structure = flat-package | `architecture/structure/flat-package.md` |
| **Architecture styles this playbook does not apply** (reference) | Omissions are decisions — plugin architectures, space-based grids, serverless-per-function and a shared-library "platform" are not targets here. | applies to every decomposition | `architecture/structure/styles-not-applied.md` |

## Considered and not applied

An architecture review that lists only what you did is half a review.
These were evaluated against this repository and left out; the reason
is the decision.

| Pattern | Why not | Read |
| --- | --- | --- |
| Backend for frontend | A single Go gateway is what this run scaffolds; per-client backends are a later refinement you add by hand. | `architecture/communication/backend-for-frontend.md` |
| Dapr | This run talks to NATS/Kafka and HTTP directly from Go; Dapr is an additional sidecar runtime nothing here installs. | `architecture/communication/dapr.md` |
| Saga — choreography | needs coordination = choreography — you chose orchestration | `architecture/communication/saga-choreography.md` |
| Synchronous compensated call chain | needs event_backbone = none — you chose nats | `architecture/communication/sync-compensated-chain.md` |
| Change data capture | needs migration_approach = parallel-run — you chose strangler | `architecture/data/change-data-capture.md` |
| CQRS | needs read_write_model = cqrs — you chose crud | `architecture/data/cqrs.md` |
| Durable execution | This platform generates a saga whose state lives in the orchestrating service’s own database — no extra runtime to operate. Durable execution (Temporal Go SDK) is the upgrade path when flows outgrow it. | `architecture/data/durable-execution.md` |
| Event sourcing | needs read_write_model = event-sourced — you chose crud | `architecture/data/event-sourcing.md` |
| Materialized view / read model | needs read_write_model = cqrs — you chose crud | `architecture/data/materialized-view.md` |
| Read replica | needs at least 1 read_heavy_services — this run has 0 | `architecture/data/read-replica.md` |
| Read-through cache | needs at least 1 read_heavy_services — this run has 0 | `architecture/data/read-through-cache.md` |
| Try–Confirm–Cancel | TCC requires every participant to expose Try/Confirm/Cancel; this platform generates compensating sagas because a strangler extraction has participants (and a legacy monolith) that cannot be reshaped that way. Adopt it per participant where isolation is worth it. | `architecture/data/try-confirm-cancel.md` |
| Decompose by subdomain (DDD) | needs boundary_strategy = subdomain — you chose domain | `architecture/decomposition/decompose-by-subdomain.md` |
| Modular monolith first | needs at most 2 services — this run has 5 | `architecture/decomposition/modular-monolith-first.md` |
| Parallel run | needs migration_approach = parallel-run — you chose strangler | `architecture/decomposition/parallel-run.md` |
| Service-to-service identity | needs at least 6 services — this run has 5 | `architecture/observability/service-identity.md` |
| Service mesh | needs at least 6 services — this run has 5 | `architecture/reliability/service-mesh.md` |
| Sidecar / ambassador | Sidecarless (ambient / eBPF) meshes are displacing it, and nothing in this run scaffolds a sidecar. | `architecture/reliability/sidecar.md` |
| Hexagonal (ports and adapters) | needs service_internal_structure = hexagonal — you chose flat-package | `architecture/structure/hexagonal.md` |

