---
id: "G-141"
type: decision
title: "Never run gpt-6-astra unless the owner names it"
status: accepted
created: "2026-09-25T01:56:58Z"
updated: "2026-09-25T01:57:10Z"
relates_to: ["G-135", "G-139"]
---

## Decision

Owner, 2026-09-24, after the G-135 Codex runs: a session, an attempt or an
eval Grove starts does not run `gpt-6-astra` unless the owner names it for
that use, and a default in the owner's `~/.codex/config.toml` is not that
instruction. More generally, what a paid run spends (model, effort, runs,
cap) is the owner's explicit answer: a mandate question never says a blank
item takes the recommendation, and a blank or missing item leaves the
question open and blocking, so the work waits and asks again rather than
spending.

Evidence: [G-139](G-139-what-mandate-and-login-should-th.md) recommended
`gpt-6-astra` at `high` because the owner's `config.toml` defaults to it,
and said any item left blank takes the recommendation. The owner's answer
reported the login only; the model, effort and runs were never approved,
and the attempt read the blanks as approval instead of blocking. The
[G-135](G-135-run-the-g-108-eval-pair-on-codex.md) attempt then ran 9½
Codex shaping runs on it in 24 minutes, 2026-09-24T23:42Z to
2026-09-25T00:06Z, moving the ChatGPT Plus five-hour window from 12% to 95%
used (8 to 12 points a run) and the weekly window from 2% to 15%, as the
rollouts' `rate_limits` readings under `~/.cache/grove-evals/codex/sessions`
show. The owner stopped the attempt.

## Alternatives

- **Follow the owner's Codex configuration default**: what G-139 did; the
  interactive default is the owner's own choice for their sessions, not a
  mandate for unattended runs.
- **A plan-quota cap alone** (the G-135 runner's `--max-plan-percent`):
  bounds how much a batch spends, but not which model spends it; both apply.

## Reconsideration

When the owner names another model to keep out of unattended use, or the
plan's limits change what a run on astra costs.
