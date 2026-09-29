# Create the orders-catalog slice — 1 files (contracts)

_edit — change these files, compile, verify_

Create this whole slice as ONE coherent unit: catalog_contract_test.go. The acceptance criteria below are the specification — implement exactly those, nothing more. There is no existing implementation to port from. Use the grounding for the real signatures of the types and helpers you must call, and write every file in the slice against the others (they are generated together, so keep them consistent). Do not invent APIs, and do not survey the wider tree. Then make the module type-check.

EVERY change these files need — finish ALL of them, including any the title does not name:
- catalog_contract_test.go: Contract tests orders->catalog: for each testdata/contracts/catalog/*.json golden response assert it decodes into internal/contracts/catalog with no unknown fields (json.Decoder.DisallowUnknownFields) and satisfies the value invariants; record one fixture (orders-service)

## Restricted access — edit only these (focus set)
- `services/orders/internal/clients/catalog_contract_test.go`

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
- services/orders/internal/orders/handlers.go
