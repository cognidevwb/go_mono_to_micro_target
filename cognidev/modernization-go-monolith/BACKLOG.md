# Backlog — Go monolith → microservices

Tracked, deliberately-unbuilt work. Everything here is documented as a reference pattern
but NOT generated, and **not offered in any questionnaire**, so no team can pick it and hit
a gap. An option is added to the questionnaire only once its contract is emitted by the
deterministic tools and proven by a live e2e.

## 1. Framework switch during the port (`web_framework` ≠ the monolith's)
**Status:** offered, but a gin monolith ported onto chi turns every handler into a develop
task (gin.Context → http.Request/ResponseWriter). **To build:** a mechanical handler adapter
in `port-monolith-go` for the gin↔echo↔chi signatures, so a switch stays a move.

## 2. gorm → pgx + sqlc as a mechanical port
**Status:** `data_access: pgx-sqlc` marks every ported gorm call as a seam. **To build:** emit
`query.sql` + `sqlc.yaml` from the gorm call sites the shards already record
(entity-access.json) so sqlc generates the replacement; keep the model call only for the
writes the saga owns.

## 3. Backend-for-Frontend edge axis
Documented in `architecture/communication/backend-for-frontend.md`; one generated gateway is
always scaffolded. Needs a `clients` question and a gateway module per client.

## 4. Google Cloud target (`cloud: gcp`)
Cloud Run + Artifact Registry + Secret Manager + Pub/Sub (Watermill has a Pub/Sub adapter).
Mirror the shape the AWS deploy has.

## 5. gRPC / Connect as the fleet-wide sync transport
Carried over only where the monolith already has a `.proto`. A fleet-wide switch needs
buf generation in the scaffold and HTTP/2 routes in the gateway.

## 6. Honour a purely synchronous split
`event_backbone: none` still scaffolds the saga when map-services-go finds a cross-context
`db.Transaction`. To offer it truthfully, expand-plan + scaffold must suppress saga/outbox
files when there is no cross-service write.

## 7. `scan-runtimes` analogue
Dropped (see README). Revisit only if a monolith turns up that needs a toolchain Go's
forward compatibility cannot provide (a cgo dependency pinned to an old C toolchain).

## Verification debt
- Live e2e on `go-monolith-shop` and `golang-gin-realworld-example-app` with a model
  configured; then a second variation (Kafka + hexagonal + repo-per-service).
