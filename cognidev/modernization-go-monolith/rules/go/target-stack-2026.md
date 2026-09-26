# Target stack — Go microservices, 2026 (pinned)

Every generated service is the same shape. The develop loop writes code against
exactly this; nothing here is a suggestion.

| Concern | Pinned choice | Notes |
| --- | --- | --- |
| Toolchain | Go 1.26 (`go 1.26` in every go.mod and go.work) | `CGO_ENABLED=0` unless the service declares cgo |
| Workspace | `go.work` at the root: `./gateway`, `./pkg`, `./services/<ctx>` … | one module per service; no `replace` between services |
| Router | the monolith's (gin v1.12.0 / echo v4.13.3 / chi v5.3.0) or the `web_framework` answer | handlers keep their signatures |
| Data | gorm kept (one `*gorm.DB` per service, own DSN) or pgx v5.8.0 + sqlc | migrations under `services/<ctx>/migrations/` (goose SQL) |
| Messaging | NATS JetStream (nats.go) default, Kafka (franz-go) opt; Watermill when chosen | transactional outbox table per service, relay goroutine |
| Sync calls | generated `internal/clients/<provider>/client.go` + consumer-owned `internal/contracts/<provider>/`, `net/http` + otelhttp, timeout + backoff | base URL from env `<CTX>_SERVICE_URL` |
| Config | env → struct (caarlos0/env v11.3.1) | secrets only from env / Key Vault refs |
| Logs / traces / metrics | log/slog JSON; OpenTelemetry v1.44.0 OTLP/HTTP | `OTEL_EXPORTER_OTLP_ENDPOINT` |
| Health | `/healthz` (liveness), `/readyz` (DB + broker ping) | |
| Gateway | generated Go reverse proxy (httputil) or Envoy / Traefik | per-route strangler weight |
| Images | multi-stage `golang:1.26` → `gcr.io/distroless/static-debian12:nonroot` | non-root UID 65532 |
| Deploy | Azure Container Apps (default), Kubernetes/Helm, docker compose | GOMEMLIMIT = 90% of the memory limit |
| CI | GitHub Actions / Azure DevOps / GitLab: vet, golangci-lint, test, govulncheck, Syft SBOM, Grype | never `aquasecurity/trivy-action` |
| Tests | testify v1.10.0; testcontainers-go v0.39.0 for the DB | table-driven |

## Rules the shape implies
- A service owns its tables. It never opens another service's DSN, never joins
  across ownership, and never imports another service module.
- Cross-context reads are a client call or a locally-held read model fed by events.
  Cross-context writes are a saga over the outbox.
- `context.Context` flows from the handler into every client, store and publish
  call; request-scoped values (user id, tenant) come from verified token claims.
- Graceful shutdown: `signal.NotifyContext` → stop accepting → drain HTTP →
  stop consumers → flush the outbox relay → close the DB.
