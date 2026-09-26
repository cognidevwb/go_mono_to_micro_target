# Architecture preferences (keyed to the pattern catalog)

The binding pattern choices are made in the questionnaire and routed **deterministically** —
this file is where you add the *why* and the *per-service nuance* the questionnaire can't hold.
The tool's full pattern library, with fit signals and trade-offs, is in the playbook's
`architecture/` catalog (`architecture/00-catalog.md`). Reference it by name here.

## Right-sizing ambition
<!-- modular (keep one deployable) | selective (extract a few, keep a core) | full (split all).
Selective is the default. e.g. "selective — extract Billing + Search only." -->

## Read/write model — per service (CQRS?)
<!-- Default is plain CRUD per service. Name only the services where a read model actually hurts.
e.g. "CQRS on Catalog (read-heavy, search shapes fight the write model); everything else CRUD."
Do NOT event-source unless you truly need an audit/temporal log (Billing ledger, maybe). -->

## Coordination — sagas
<!-- orchestration (default; a coordinator owns the flow) | choreography (peers react to events).
e.g. "the CreateOrder flow is an orchestrated saga owned by Ordering." -->

## Edge — gateway vs BFF
<!-- Single Go reverse-proxy gateway (default) | a BFF per client (web / mobile / partner).
e.g. "we have a web SPA + a mobile app with different needs → a web BFF and a mobile BFF." -->

## Platform toggles
<!-- Deployment target (Kubernetes / Azure Container Apps / docker-compose), service mesh
(usually NO for a small estate), Dapr (only if you want broker portability). State your reality. -->

## Patterns to AVOID
<!-- Anything you've been burned by or won't operate. e.g. "no service mesh — no platform team";
"no event sourcing"; "no Kafka, NATS only". The tool will steer around these. -->
