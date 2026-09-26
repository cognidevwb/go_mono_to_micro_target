# Create the orders slice — 1 files (equivalence)

_edit — change these files, compile, verify_

Create this whole slice as ONE coherent unit: orders_equivalence_test.go. Read the grounding FIRST — the monolith sources this is ported from + the structural shards — and write every file against the others' real signatures (they are generated together, so make them consistent). Port the real behaviour from the monolith; do not invent APIs. Then make the module type-check.

EVERY change these files need — finish ALL of them, including any the title does not name:
- orders_equivalence_test.go: Equivalence orders: for each endpoint derive a {request, seed, expectedResponse} golden from the monolith's handler into testdata/equivalence/, replay it against the service with httptest, normalise volatile fields, assert new == legacy (orders-service)

## Restricted access — edit only these (focus set)
- `services/orders/test/equivalence/orders_equivalence_test.go`

_Expansion: none — do not edit existing files in this step._

## Acceptance criteria (done when)
- All 1 files in the slice are written coherently against each other's real signatures.
- Ported from the monolith sources in the grounding — no invented APIs, reuses the app's ORM/auth libraries.
- The whole module compiles clean (the verification gate is green) after this task.

## Read first (grounding)
- .cognidev/understand/graph.json
- .cognidev/understand/state.json
- .cognidev/understand/flows.json
- .cognidev/understand/uml.json
- .cognidev/understand/files.json
- internal/orders/handlers.go
- internal/orders/model.go
- internal/orders/service.go
- services/orders/internal/orders/handlers.go
