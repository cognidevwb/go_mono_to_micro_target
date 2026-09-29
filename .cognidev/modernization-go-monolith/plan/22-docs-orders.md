# Create the docs-orders slice — 1 files (docs)

_edit — change these files, compile, verify_

Create this whole slice as ONE coherent unit: README.md. Read the grounding FIRST — the monolith sources this is ported from + the structural shards — and write every file against the others' real signatures (they are generated together, so make them consistent). Port the real behaviour from the monolith; do not invent APIs. Then make the module type-check.

EVERY change these files need — finish ALL of them, including any the title does not name:
- README.md: Docs: fill orders's <!-- CW-SEAM[kind=business-capability service=orders] --> — 2-3 sentences on the business capability orders owns, grounded in its source (orders-service)

## Restricted access — edit only these (focus set)
- `services/orders/README.md`

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
