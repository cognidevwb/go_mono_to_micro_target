---
id: observability
title: Observability — slog + OpenTelemetry
family: observability
currency: standard
role: invariant
applies_when:
  - always
teaches: Traces, metrics and logs over one OTLP pipe; the trace id in every log line; the backend is a swap.
---

# Observability — slog + OpenTelemetry

## What it is

Every service emits OTLP traces and metrics via the OpenTelemetry Go SDK, instruments HTTP (`otelhttp`), database (`otelpgx` / gorm `tracing` plugin) and NATS calls, and logs with `log/slog` JSON including `trace_id`. A collector routes to Grafana LGTM, Azure Monitor or another OTLP backend.

## Currency (2026)

**STANDARD — invariant.** OpenTelemetry Go traces/metrics are stable; logs bridge via `otelslog`.

## Fit signals (from `.cognidev/feature/context.md`)

- Always.

## What it changes in the generated services

- `internal/telemetry` (or chassis) setup; `OTEL_EXPORTER_OTLP_ENDPOINT`, `OTEL_SERVICE_NAME` env; an otel-collector in compose/Helm.

## Go 2026 implementation

```go
exp, _ := otlptracegrpc.New(ctx)
tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exp), sdktrace.WithResource(res))
otel.SetTracerProvider(tp); otel.SetTextMapPropagator(propagation.TraceContext{})
```

## Anti-patterns

- Logging with `log.Printf` (the fixture's jobs do) — no level, no trace id.
- Propagating trace context by hand.

## Interacts with

- [[chassis]] · [[health-checks]]
