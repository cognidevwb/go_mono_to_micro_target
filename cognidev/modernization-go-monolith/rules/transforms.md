# Permitted drift — Go monolith → microservices

Read by the sign-off stage (docs/73 §4 L2). Each `transform` block names one difference
this decomposition *intends* between a corresponded source⇄target pair. A difference
matching a rule becomes `PORTED_EQUIVALENT` citing the rule id; anything unmatched
stays `SUSPECT`. This is a Class B migration (docs/73 §2): same language, correspondence
is a *partition* — many monolith files map onto one service, via
`service-map.json` (file → service) and `port-report.json` (`ported_from` → `path`).

## What this migration legitimately changes

`port-monolith-go` is a mover, not a rewriter: each owned file is copied into
`services/<ctx>/internal/<pkg>/` **verbatim**, package clause kept. Only four things are
rewritten mechanically, and they are the whole source of intended drift:

1. **Import paths** — a package that moved with the file gets its new module path. No
   fingerprint token is derived from an import path, so this is invisible.
2. **The shared `*gorm.DB`** (`platform.Open`) becomes the service's own handle on its
   own database; `AutoMigrate`/`Register` lists only the service's owned models. A unit of
   work that used to span several contexts therefore commits fewer tables.
3. **Every foreign context's service type** (`*catalog.Service` held in a struct field)
   becomes the generated client interface in `internal/clients`, and the call site keeps
   its method name and gains a leading `ctx context.Context` argument plus an error
   return where the in-process call had none.
4. **Seams are marked** — `CW-SEAM[kind=cross-call]` on each rewritten call,
   `CW-SEAM[kind=saga]` above each cross-context `db.Transaction`, `CW-SEAM[kind=event]`
   on each in-process publish/subscribe that now crosses services — and develop fills
   them.

Everything else is generated platform: go.work, `cmd/<ctx>/main.go`, the gateway, the
outbox relay. Those have no source ancestor and belong in the `UNCLAIMED` pile for L4
triage, not here.

## What it must conserve

- **`decisions`, `arithmetic`, `states`** — ported byte-for-byte. Any change on a ported
  pair is a defect. The one exception is additive request validation below.
- **Routes.** A ported handler keeps its method + path registration; the `/<service>/`
  prefix exists only in the gateway's route table. A changed `ep:`/`sig:` is a defect.
- **The compensation leg.** The saga rule permits the consistency boundary to narrow to
  one service's own tables. It does not permit compensations to be absent — no rule here
  lets an `e:compensate:` token be missing.
- **Provider read-model shape.** `contracts-go` locks each provider's exported fields and
  json tags into `contracts.lock.json`; a changed DTO shape is drift that tool reports.
- **A `seam(...)` "not implemented yet" error (or `ErrSeamNotImplemented`) is never permitted** in a signed-off tree.
- **Reads dropped by the database split** get no rule — they must reappear as a client
  call plus a consumer-owned `internal/contracts/<provider>` DTO, and are left `SUSPECT` otherwise.

---

## Cross-context calls

```transform
id: in-process-call-became-a-client-call-with-context
title: An in-process method call on another context's service became a typed-client call taking a context
field: calls_out
op: rename
from: *
to: *
note: port-monolith-go retargets the receiver from *<ctx>.Service to the generated clients.<Ctx>Client interface and keeps the method name, adding ctx as the first argument. The callee leaf name is unchanged, so this rule only matches when the same method name appears on both sides of the pair; a call that was deleted matches nothing and blocks.
```

```transform
id: remote-call-added-transport-error-handling
title: A remote call added transport error handling the in-process call could not have had
field: errors
op: allow-added
match: e:if:err*
note: An in-process call to a method with no error return has no failure mode; an HTTP call does. Checking the returned error at the call site is part of the intended design. Additive only — every error path the monolith had must still correspond.
```

## The saga

```transform
id: cross-context-transaction-became-a-saga
title: One cross-context db.Transaction became a single-service commit inside a saga
field: writes
op: rename
from: cb:transaction(*|*
to: cb:transaction(*
note: map-services-go detects a db.Transaction whose effective write set spans contexts (dataflow.json consistency boundaries) and port-monolith-go marks it CW-SEAM[kind=saga]. Database-per-service means the orchestrator can only commit its own tables. The `from:` glob requires a multi-table boundary, so a single-table commit is never touched. Licenses the narrowing ONLY.
```

```transform
id: saga-added-compensation-paths
title: The saga added compensation paths for the writes it can no longer roll back
field: errors
op: allow-added
match: e:compensate*
note: With the transaction split, a failed participant is undone by an explicit compensating command, hardest-to-reverse step last. Additive-only: nothing here permits a compensation to be missing.
```

```transform
id: saga-added-an-outbox-write
title: The saga added a transactional outbox write alongside its own commit
field: writes
op: allow-added
match: w:*outbox*
note: The saga persists its state and its outgoing messages in one local transaction through the service's outbox table (generated by scaffold-services-go). Scoped to outbox-named write targets so it cannot absorb a business table appearing out of nowhere.
```

```transform
id: saga-added-event-publication
title: The saga publishes messages where the monolith called the next step directly
field: calls_out
op: allow-added
match: publish*
note: Where the monolith completed the whole unit of work inline, the orchestrator now publishes through the broker (NATS JetStream by default). Added external effect only; it licenses no removal.
```

## In-process events

```transform
id: in-process-bus-became-a-broker-topic
title: A platform.Bus publish became an outbox-backed broker publish
field: calls_out
op: rename
from: publish*
to: publish*
note: The monolith's in-process Bus.Publish(topic, payload) becomes a publish of the same topic name through the outbox. Topic names are conserved (order.placed stays order.placed); the rule matches only a publish on both sides.
```

## Cross-context reads

```transform
id: foreign-model-read-became-a-contract-dto
title: A foreign context's model became a caller-side contract DTO
field: io
op: allow-added
match: dto:*
note: The consumer now receives its own internal/contracts/<provider> structs — exported fields and json tags copied from the provider's model, never the gorm model itself. Additive only; a diverging copy is reported by contracts-go against contracts.lock.json.
```

## Request surface

```transform
id: request-validation-became-declarative
title: An inline request guard was also expressed as declarative binding validation
field: decisions
op: allow-added
match: d:validate*
note: Request structs may gain `binding:`/`validate:` tags (go-playground/validator) that restate at the boundary what the ported handler already checks. Additive only; a dropped inline check still shows as a decisions shortfall.
```

---

## The composition root

The monolith's `main` opens the database, builds the router, subscribes the in-process
handlers and starts its background jobs. After the carve each service's composition root
(`internal/app/wire.go`, `internal/app/app.go`, `cmd/<ctx>/main.go`) does its own share,
and three of those responsibilities change shape by design:

- an in-process `bus.Subscribe(topic, svc.Handler)` becomes a durable broker consumer,
  `events.Subscribe<Topic>(ctx, bus, seen, svc)`, generated in `internal/events/` of the
  service that owns the handler; the consumer calls the same handler method;
- an in-process ticker (`StartXxxJob(ctx, svc)`) becomes a leased scheduled job,
  `jobs.Run(ctx, db, interval, svc.Work)`, so one replica runs it — the same work function;
- `gin.Default()` becomes `gin.New()` plus `gin.Recovery()`; request logging moves to the
  OpenTelemetry handler in `cmd/<ctx>/main.go`.

```transform
id: in-process-subscription-became-a-broker-consumer
title: A bus.Subscribe in main became the owning service's durable consumer
field: calls_out
op: rename
from: subscribe
to: subscribe*
note: events.Subscribe<Topic> subscribes on the broker and calls the same handler method; the topic and the handler are read inside the generated consumer.
```

```transform
id: subscription-arguments-moved-into-the-consumer
title: The topic constant and handler method value are read by the generated consumer
field: reads
op: allow-removed
match: r:*.topic*
```

```transform
id: subscription-handler-moved-into-the-consumer
title: The handler method value is referenced by the generated consumer
field: reads
op: allow-removed
match: r:service.on*
```

```transform
id: ticker-job-became-a-leased-scheduled-job
title: A StartXxxJob ticker became jobs.Run over the same work function
field: calls_out
op: allow-removed
match: start*job
note: The owning service's wire.go schedules the same work method with jobs.Run under a database lease.
```

```transform
id: default-router-became-new-plus-recovery
title: gin.Default became gin.New with Recovery; access logs moved to OpenTelemetry
field: calls_out
op: allow-removed
match: default
```

## Shared state

`develop` fills each `CW-SEAM[kind=shared-state]` by the rule in
`context/develop-feature.md`: a package-level cache stays only in the owning service and,
because each replica now holds its own copy, gets a TTL and is invalidated on write. The
cached values and where they come from do not change; how long a replica may serve one does.
Added time and map calls (`time.Now`, `Before`, `delete`) implement that expiry.

```transform
id: shared-cache-gained-a-ttl
title: A process-local cache became TTL-bound (its entry type carries an expiry)
field: states
op: allow-added
match: st:*cachettl*
```

```transform
id: shared-cache-entry-type
title: The cache's value type became an entry with an expiry
field: io
op: allow-added
match: dto*:*cacheentry*
```

```transform
id: shared-cache-entry-reads
title: Reads of the cache entry's value and expiry
field: reads
op: allow-added
match: r:*cacheentry.*
```

```transform
id: shared-cache-miss-db-read
title: The cache-miss path still reads the owning service's own database
field: reads
op: allow-added
match: r:db
```

```transform
id: shared-cache-expiry-guard
title: The cache hit check also tests the entry's expiry
field: decisions
op: rename
from: d:guard()*
to: d:guard(*cacheentry.expiresat*
```

```transform
id: shared-cache-expiry-clock
title: The cache reads the clock to stamp and test an entry's expiry
field: calls_out
op: allow-added
match: now
note: time.Now() stamps an entry when it is stored and is compared when it is read. Added calls only; the monolith's lookups and the owning database read stay as they were.
```

```transform
id: shared-cache-expiry-compare
title: The cache hit compares now with the entry's expiry
field: calls_out
op: allow-added
match: before
```

```transform
id: shared-cache-expiry-stamp
title: The stored entry's expiry is now plus the TTL
field: calls_out
op: allow-added
match: add
note: time.Now().Add(ttl). A method named Add that is not the clock would also match; the rule only permits calls that were ADDED, so nothing the monolith did can disappear under it.
```
