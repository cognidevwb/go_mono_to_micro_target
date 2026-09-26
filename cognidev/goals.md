# Why you're modernizing (the goals behind the split)

The tool can see your code; it can't see your reasons. Your reasons decide **how much** to
extract. In 2026 the honest default is *modular-monolith-first, extract only against a named
force* — so what you write here directly sets the right-sizing gate. Be blunt about priorities.

## Primary drivers — rank the ones that apply
<!-- The forces that actually justify pulling a service out. Rank or delete.
- Team autonomy — separate teams keep colliding in one deploy pipeline
- Independent scaling — one area's load (search, pricing, media) dwarfs the rest
- Release velocity — one slow/risky area blocks everyone else's releases
- Fault isolation — a failure in X must not take down checkout
- Tech/lifecycle divergence — one area needs a different runtime/cadence/datastore
- Cost / cognitive load — (usually a reason to KEEP things together, say so) -->

## What "done" looks like
<!-- Concrete success criteria. e.g. "the Payments team can deploy without coordinating";
"checkout stays up when reporting is down"; "we can scale search to 10x independently". -->

## Non-goals — what you do NOT want
<!-- e.g. "we are NOT rewriting business logic"; "we are NOT moving off MySQL";
"do not touch the legacy reporting exports". Non-goals prevent scope creep in the plan. -->

## Timeline & appetite
<!-- e.g. "a 2-service pilot this quarter, not a big-bang"; "hard cutover date in Q3";
"low appetite for operational complexity — keep the moving parts few". -->
