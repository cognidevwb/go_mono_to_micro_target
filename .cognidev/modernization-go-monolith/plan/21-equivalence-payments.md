# Create the payments slice — 1 files (equivalence)

_edit — change these files, compile, verify_

Create this whole slice as ONE coherent unit: payments_equivalence_test.go. Read the grounding FIRST — the monolith sources this is ported from + the structural shards — and write every file against the others' real signatures (they are generated together, so make them consistent). Port the real behaviour from the monolith; do not invent APIs. Then make the module type-check.

EVERY change these files need — finish ALL of them, including any the title does not name:
- payments_equivalence_test.go: Equivalence payments: for each endpoint derive a {request, seed, expectedResponse} golden from the monolith's handler into testdata/equivalence/, replay it against the service with httptest, normalise volatile fields, assert new == legacy (payments-service)

## Restricted access — edit only these (focus set)
- `services/payments/test/equivalence/payments_equivalence_test.go`

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
- internal/payments/gateway.go
- internal/payments/model.go
- internal/payments/service.go
- services/payments/internal/httpapi/payment_handlers.go
