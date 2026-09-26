# gateway — built from the workspace root:  docker build -f gateway/Dockerfile .
ARG GO_IMAGE=golang:1.26-bookworm
FROM ${GO_IMAGE} AS build
WORKDIR /src
COPY pkg/ pkg/
COPY gateway/ gateway/
WORKDIR /src/gateway
ENV GOWORK=off CGO_ENABLED=0
RUN go mod download
RUN go build -trimpath -ldflags="-s -w" -o /out/gateway ./cmd/gateway

# distroless static, non-root (UID 65532): no shell, no package manager.
FROM gcr.io/distroless/static-debian12:nonroot
USER nonroot:nonroot
COPY --from=build /out/gateway /app
EXPOSE 8080
ENTRYPOINT ["/app"]
