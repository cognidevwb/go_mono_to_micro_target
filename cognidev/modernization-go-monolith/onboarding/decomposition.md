# Decomposition intent (who owns what, and the force behind each split)

This is the highest-signal file for the plan. The tool computes candidate boundaries from
coupling; you supply the two things it can't compute: **team ownership** and **the named
force** that justifies pulling a service out. Every extraction should trace to a force here —
no force means the tool keeps it in the modular core (that's the mature 2026 default).

## Teams & ownership
<!-- Who owns (or would own) which area. Team ownership is the #1 legitimate reason to split.
e.g.
- Payments team → Billing, Payments
- Catalog team → Product, Search, Pricing
- (no dedicated team) → everything else stays in the core host -->

## Extract these — with the force
<!-- Services you want pulled out, each with WHY. e.g.
- Billing — force: team-autonomy (Payments team owns its release cadence)
- Search — force: independent-scaling (indexing load dwarfs the rest)
An entry with no force will be questioned by the plan. -->

## Keep these together (do NOT split)
<!-- Things that look separable but shouldn't be, and why. e.g. "Inventory + Ordering stay
together — they share the reservation transaction; splitting them creates a saga we don't
want yet." -->

## Extraction order preference
<!-- Strangler sequencing. e.g. "extract Notifications first (low-risk, leaf); leave Ordering
(the hub) for last." Leave blank to let the tool sequence by coupling. -->
