# Decomposing this monolith

5 services will be extracted from this repository, owning 7 tables between them, with 2 cross-service calls to resolve and 1 saga.

Read in this order:

1. `what-changes.md` — what moves where, and what it costs.
2. `the-plan.md` — the extraction order, and why that order.
3. `what-you-learn.md` — the architecture this cut teaches, and where to read it.
4. `tasks/` — one document per service, with its files.
5. `your-choices.yaml` — what you were asked, and what was decided for you.
6. `status/` — filled after the run: what each task actually did.

`share/` holds the same documents as standalone pages, for sending to
someone without a checkout.

---

Every number here was measured by an earlier phase of this run. Nothing
in these documents was re-derived by reading the tree a second time.
