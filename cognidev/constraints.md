# Hard constraints (the walls the plan must stay inside)

Non-negotiables the tool must respect. A constraint here **overrides** a pattern default —
e.g. "must stay on MySQL" turns off the Postgres default; "PCI in scope" pushes an
audit trail into the design. State the limit and, if useful, why.

## Technology you must keep
<!-- e.g. "MySQL (the DBA team only runs MySQL)"; "on-prem only, no public cloud"; "the legacy reporting module stays in the monolith — leave it"; "must run on RHEL 8 hosts". -->

## Compliance & data
<!-- e.g. "PCI-DSS on the payments path"; "GDPR — customer PII stays in-region";
"HIPAA audit trail required" (implies immutable audit / event log); "no data leaves the VPC". -->

## Organizational limits
<!-- e.g. "team of 6 — keep operational surface small"; "no platform team, so no service
mesh"; "one shared CI system"; "we cannot run Kafka, NATS only". -->

## Do not touch
<!-- Systems/areas that are off-limits for this effort. e.g. "the nightly batch job";
"the third-party ERP integration"; "the mobile API contract is frozen until v3". -->
