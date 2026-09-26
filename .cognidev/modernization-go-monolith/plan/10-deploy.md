# Create the deploy slice — 5 files (deploy)

_edit — change these files, compile, verify_

Create this whole slice as ONE coherent unit: strangler.values.yaml, strangler.values.yaml, strangler.values.yaml, strangler.values.yaml, strangler.values.yaml. The acceptance criteria below are the specification — implement exactly those, nothing more. There is no existing implementation to port from. Use the grounding for the real signatures of the types and helpers you must call, and write every file in the slice against the others (they are generated together, so keep them consistent). Do not invent APIs, and do not survey the wider tree. Then make the module type-check.

EVERY change these files need — finish ALL of them, including any the title does not name:
- strangler.values.yaml: Deployment: set catalog's strangler traffic weight 0..100 + legacyPaths (CW-SEAM), grounded in catalog's coupling — a leaf higher, the hub / saga orchestrator canary low (catalog-service)
- strangler.values.yaml: Deployment: set customers's strangler traffic weight 0..100 + legacyPaths (CW-SEAM), grounded in customers's coupling — a leaf higher, the hub / saga orchestrator canary low (customers-service)
- strangler.values.yaml: Deployment: set inventory's strangler traffic weight 0..100 + legacyPaths (CW-SEAM), grounded in inventory's coupling — a leaf higher, the hub / saga orchestrator canary low (inventory-service)
- strangler.values.yaml: Deployment: set payments's strangler traffic weight 0..100 + legacyPaths (CW-SEAM), grounded in payments's coupling — a leaf higher, the hub / saga orchestrator canary low (payments-service)
- strangler.values.yaml: Deployment: set orders's strangler traffic weight 0..100 + legacyPaths (CW-SEAM), grounded in orders's coupling — a leaf higher, the hub / saga orchestrator canary low (orders-service)

## Restricted access — edit only these (focus set)
- `services/catalog/deploy/strangler.values.yaml`
- `services/customers/deploy/strangler.values.yaml`
- `services/inventory/deploy/strangler.values.yaml`
- `services/payments/deploy/strangler.values.yaml`
- `services/orders/deploy/strangler.values.yaml`

_Expansion: none — do not edit existing files in this step._

## Acceptance criteria (done when)
- All 5 files in the slice are written coherently against each other's real signatures.
- Written to the acceptance criteria above, against the real signatures of the files in the grounding — no invented APIs, reuses the app's own libraries.
- The whole module compiles clean (the verification gate is green) after this task.

## Read first (grounding)
- .cognidev/understand/graph.json
- .cognidev/understand/state.json
- .cognidev/understand/flows.json
- .cognidev/understand/uml.json
- .cognidev/understand/files.json
- services/catalog/internal/catalog/handlers.go
- services/customers/internal/customers/handlers.go
- services/inventory/internal/httpapi/stock_item_handlers.go
- services/payments/internal/httpapi/payment_handlers.go
- services/orders/internal/orders/handlers.go
