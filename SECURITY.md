# Security

- Identity: jwt-keycloak. The gateway verifies the bearer token (OIDC discovery via `AUTH_ISSUER`); each service re-validates.
- Images: distroless static, non-root (UID 65532), read-only root filesystem in Kubernetes.
- Supply chain: `govulncheck` on every module, Syft SBOM + Grype scan (fail on high) on every image.
- Secrets: `.env.example` documents them; `.env` is never committed (pre-commit hook).
