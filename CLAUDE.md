# Working in this workspace

- One Go module per service under `services/`; the shared kernel is `pkg/` (transport only — no domain types).
- Build every module: `for m in $(go list -m -f '{{.Dir}}'); do (cd $m && go build ./... && go vet ./...); done`.
- `// CW-SEAM[kind=...]` marks work the decomposition left: implement it, then delete the marker.
