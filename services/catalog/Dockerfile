# catalog-service — built from the workspace root so the shared kernel (pkg/)
# is in the build context:  docker build -f services/catalog/Dockerfile .
ARG GO_IMAGE=golang:1.26-bookworm
FROM ${GO_IMAGE} AS build
WORKDIR /src
COPY pkg/ pkg/
COPY services/catalog/ services/catalog/
WORKDIR /src/services/catalog
ENV GOWORK=off CGO_ENABLED=0
RUN go mod download
RUN go build -trimpath -ldflags="-s -w" -o /out/catalog ./cmd/catalog

# distroless static, non-root (UID 65532): no shell, no package manager.
FROM gcr.io/distroless/static-debian12:nonroot
USER nonroot:nonroot
COPY --from=build /out/catalog /app
EXPOSE 8080
ENTRYPOINT ["/app"]
