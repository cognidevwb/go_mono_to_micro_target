# Update the modify slice — 1 files (services)

_edit — change these files, compile, verify_

Update this whole slice as ONE coherent unit: service.go. The acceptance criteria below are the specification — implement exactly those, nothing more. There is no existing implementation to port from. Use the grounding for the real signatures of the types and helpers you must call, and write every file in the slice against the others (they are generated together, so keep them consistent). Do not invent APIs, and do not survey the wider tree. Then make the module type-check.

EVERY change these files need — finish ALL of them, including any the title does not name:
- service.go: replace the cross-service unit of work in PlaceOrder with the CreateOrder saga (outbox + compensations) — writes in inventory, orders, payments

## Restricted access — edit only these (focus set)
- `internal/orders/service.go`

_Expansion: none — do not edit existing files in this step._

## Acceptance criteria (done when)
- All 1 files in the slice are written coherently against each other's real signatures.
- Written to the acceptance criteria above, against the real signatures of the files in the grounding — no invented APIs, reuses the app's own libraries.
- The whole module compiles clean (the verification gate is green) after this task.

## Read first (grounding)
- .cognidev/understand/graph.json
- .cognidev/understand/state.json
- .cognidev/understand/flows.json
- .cognidev/understand/uml.json
- .cognidev/understand/files.json
