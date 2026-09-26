# 3. Typed, resilient HTTP clients

Status: accepted

A cross-context call goes through a generated client that mirrors the provider's Go API, over `pkg/httpx.NewClient` (bounded timeout, trace propagation). The consumer owns its copy of every contract type.
