# ADR-0000: (template — copy to 0001-your-decision.md)

Architecture Decision Records capture a choice you want the tool to treat as **binding**.
Copy this file to `0001-<slug>.md`, fill it in, set the status to `Accepted`, and the planning
step will honor it as a constraint (not a suggestion). Keep them short; append, never edit
history.

- **Status:** Proposed | Accepted | Superseded by ADR-XXXX
- **Date:** YYYY-MM-DD

## Context
<!-- The forces at play — the situation and constraints that make this decision necessary.
e.g. "The Payments team must ship independently; Billing currently shares the monolith's DB." -->

## Decision
<!-- The choice, stated plainly. e.g. "Extract Billing into its own service with its own
Postgres schema; Ordering talks to it via an async CreateInvoice command, never a sync call." -->

## Consequences
<!-- What follows — good and bad. e.g. "Ordering no longer sees Invoice rows directly; a
read model is needed for the order-history screen; the CreateOrder flow becomes a saga." -->

## Alternatives considered
<!-- What you rejected and why. e.g. "Shared DB with a Billing schema — rejected, keeps the
deploy coupling we're trying to remove." -->
