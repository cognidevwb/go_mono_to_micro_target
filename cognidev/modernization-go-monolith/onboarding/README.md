# Your CogniDev brief — tell the tool what the code can't

This `cognidev/` folder is **your** channel to the tool. CogniDev already reads your
source: it parses the Go packages, the gorm/sqlx models, the routes, the call graph — so it knows
*what the code is*. What it cannot know is *why* — your goals, your team shape, what a
word means in your business, the constraints you live under, and where the bodies are
buried. **That is what this folder is for.**

Everything here is **optional**. An empty file is ignored — the tool falls back to what
it derives from the code. But **whatever you do write here is treated as authoritative**:
it is injected into the planning and generation steps as high-priority human context,
and it can override the tool's heuristic guesses (within the safety invariants). Filling
even one of these files materially changes the plan. That is the trade: low effort in,
high leverage out.

## What's here

| File | What it captures | Where the tool uses it |
| --- | --- | --- |
| `CLAUDE.md` | Coding standards + house rules the generator must obey | **Develop** — every line it writes into a seam |
| `goals.md` | Why you're modernizing; priorities; non-goals | **Right-sizing gate** + planning |
| `constraints.md` | Hard limits — tech you must keep, compliance, timeline, team size | Planning + scaffold defaults |
| `domain.md` | Your ubiquitous language — what terms and entities *mean* | **Boundary naming + merging** in planning |
| `decomposition.md` | Team ownership + desired boundaries + the force behind each split | **Right-sizing gate + the service map** |
| `architecture.md` | Pattern preferences, keyed to the `architecture/` catalog | Pattern selection + planning |
| `known-issues.md` | Pain points, god-classes, flaky seams, risky areas | **Strangler order** + risk flags |
| `decisions/` | Append-only ADRs — decisions you want held as binding | Planning treats accepted ADRs as constraints |

## How to fill it

- **Prose is fine.** Write like you'd brief a new senior engineer. Bullet points, half
  sentences, "we think X but aren't sure" — all useful.
- **Be specific.** "Billing must be its own service so the Payments team can own it" beats
  "we want good architecture."
- **Delete the prompts you don't use.** Each file has guiding sub-headings; keep the ones
  you fill, drop the rest. What's left is what the tool reads.
- **Commit it.** This is your domain knowledge — version it with your repo. It survives
  re-runs and gets better each time you touch it.

You never have to open this folder. But an hour spent here is the single highest-leverage
thing you can do to make the decomposition match *your* system instead of a generic one.
