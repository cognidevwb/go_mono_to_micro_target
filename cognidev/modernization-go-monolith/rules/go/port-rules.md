# Go port rules — moving a monolith package into a service module

These bind every task the develop loop runs and every repair the gates hand back.
They exist because the cheapest correct decomposition is a MOVE, and every rule
below is a way a well-meaning edit turns a move into a rewrite.

## 1. Move, don't rewrite
- `port-monolith-go` copies each owned file verbatim into
  `services/<ctx>/internal/<pkg>/<file>.go`. The package clause (`package orders`)
  is KEPT. Do not rename the package, split the file, reorder declarations, or
  "tidy" logic that is not at a seam.
- A ported file carries a `// Code ported by cognidev from <path>. DO NOT EDIT
  outside CW-SEAM markers.` header. Keep it until the last seam in the file is
  filled; then leave it — it tells a reviewer the body is the monolith's.
- Business logic (conditions, arithmetic, status transitions) is conserved
  byte-for-byte. A change there is a defect, not drift (see `transforms.md`).

## 2. Imports
- Imports of packages that moved WITH the file are rewritten mechanically to
  `<base_module>/services/<ctx>/internal/<pkg>`. Never hand-edit these.
- An import of ANOTHER context's package (`internal/catalog` inside orders) is
  replaced by the generated client `internal/clients/<provider>/client.go` + the
  consumer's OWN copy of the provider's DTOs in
  `internal/contracts/<provider>/contracts.go`. A service module must never import another
  service module — Go's `internal/` rule makes that a compile error anyway
  (`use of internal package … not allowed`), and adding a `replace` directive to
  get around it is forbidden.
- `pkg/` is a minimal shared kernel of infrastructure (httpx, otelx, events
  envelope, outbox, saga runner). It carries NO service's DTOs — each consumer
  owns its copy, and `contracts-go` reports a copy that drifts from its provider.
  No database types, no gorm models.

## 3. The CW-SEAM grammar
`// CW-SEAM[kind=<kind> key=value …]` on its own line, immediately above the code it
governs. Kinds:
- `cross-call target=<ctx> member=<Method>` — a typed client method whose body is
  returns the `seam(...)` "CW-SEAM: … is not implemented yet" error. Implement it as an HTTP call to
  `<ctx>-service` (base URL from config, `http.Client` with a timeout, retry with
  backoff on 5xx/transport errors only, OpenTelemetry propagation via otelhttp),
  decode into the `internal/contracts/<provider>` DTO, return typed errors (404 → ErrNotFound).
- `saga name=<Name> participants=<a,b,…>` — the `db.Transaction(func(tx *gorm.DB)
  error {...})` that used to span contexts. Only this service's own tables stay in
  the local transaction; each participant step becomes a command through the
  outbox; compensations run in reverse; the hardest-to-reverse step (charging
  money) runs LAST.
- `event topic=<topic> from=<ctx> to=<ctx>` — an in-process `Bus.Publish` /
  `Subscribe` that now crosses services. Publish by writing to the outbox in the
  same transaction; subscribe with an idempotent consumer keyed by message id.
- `shared-state name=<var>` — a package-level map/slice (guarded by a sync.Mutex)
  that two services would each hold a copy of. Replace with a read-through cache
  owned by one service, or an explicit lookup; never leave two drifting copies.
- `ambient key=<Key>` — a `ctx.Value(<Key>)` read whose value was set by middleware
  in another context. It now arrives as a verified JWT claim or a propagated header.
- `job name=<Func>` — a ticker goroutine; it runs in exactly ONE service (the owner
  of the tables it writes) and must stop on ctx.Done().
- `contract-impact provider=<p> item=<Type>` — see `develop-feature.md`.
Delete each marker as you implement it. A remaining marker, or any
`ErrSeamNotImplemented`, re-opens the task.

## 4. Dependencies
Only the pinned 2026 set may appear in a service go.mod: the monolith's own modules
(same versions), `github.com/go-chi/chi/v5 v5.3.0` / `github.com/gin-gonic/gin
v1.12.0` / `github.com/labstack/echo/v4 v4.13.3` (the chosen router),
`gorm.io/gorm` (kept) or `github.com/jackc/pgx/v5 v5.8.0`, `github.com/nats-io/nats.go`
or `github.com/twmb/franz-go` (the chosen broker), `github.com/ThreeDotsLabs/watermill`
(when chosen), `go.opentelemetry.io/otel v1.44.0` + contrib `v0.62.0`,
`github.com/stretchr/testify v1.10.0`, `github.com/testcontainers/testcontainers-go
v0.39.0`. The develop loop never runs `go get` for anything else and never edits
go.mod/go.work — those are scaffold-owned; a missing module is a scaffold defect to
report, not a dependency to add.

## 5. Guardrails — do not carry the monolith's bugs across the seam
- Stock reservation is ONE guarded statement (`UPDATE stock_items SET reserved =
  reserved + ? WHERE product_id = ? AND on_hand - reserved >= ?` and check
  RowsAffected) — never read-check-then-write.
- Status fields become typed string constants (`type PaymentStatus string`), not
  bare strings.
- Handlers return RFC 9457 problem JSON; errors never leak stack traces.
- Every goroutine started in a service has a stop path (ctx) and is joined on
  shutdown.
