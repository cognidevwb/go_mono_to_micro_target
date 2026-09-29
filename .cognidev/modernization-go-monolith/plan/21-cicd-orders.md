# Create the orders slice — 1 files (cicd)

_edit — change these files, compile, verify_

Create this whole slice as ONE coherent unit: deploy.targets.yml. The acceptance criteria below are the specification — implement exactly those, nothing more. There is no existing implementation to port from. Use the grounding for the real signatures of the types and helpers you must call, and write every file in the slice against the others (they are generated together, so keep them consistent). Do not invent APIs, and do not survey the wider tree. Then make the module type-check.

EVERY change these files need — finish ALL of them, including any the title does not name:
- deploy.targets.yml: CI/CD: fill orders's deploy environments + post-deploy smoke-check path (# TODO(cognidev)) targeting Azure Container Apps (deploy/azure — azd), images in ACR; build (go build, go vet, go test -race), the distroless image, the Syft SBOM and the Grype scan are already wired (orders-service)

## Restricted access — edit only these (focus set)
- `services/orders/deploy/deploy.targets.yml`

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
