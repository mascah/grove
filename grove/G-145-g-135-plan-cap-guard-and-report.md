---
id: "G-145"
type: review
title: "G-135 plan-cap guard and report: independent review"
status: current
created: "2026-09-25T04:01:37Z"
updated: "2026-09-25T04:01:49Z"
work: ["G-135"]
examined: "4513c69"
---

## Examined

Three rounds by fresh `grove-reviewer` agents, 2026-09-25, of the changes
made after the owner stopped the G-135 attempt, range `1e1a345..4513c69`
on `worktree-G-135`: the plan-cap guard in `evals/run.py` and its README
section, [G-141](G-141-never-run-gpt-6-astra-unless-the.md),
[G-142](G-142-keep-a-blank-mandate-answer-from.md), the report
[G-143](G-143-g-135-codex-eval-row-pattern-pla.md) and
[G-144](G-144-give-adopting-projects-the-recor.md). Each ran `python3
evals/run.py selftest` and `grove check` and read the retained runs under
`~/.cache/grove-evals` without writing; none ran a paid command. Round 1
examined `48358b7`, round 2 `1e1a345..643d29d`, round 3 `4513c69`, the
`examined` commit.

## Findings

Round 1 (`48358b7`), no blocking defect: the first run was gated on a
possibly stale reading with no per-run allowance; a stale prior reading
could be counted as the first run's use; a reading without `resets_at`
read as a reset window; use across a reset could go negative; the stop
message was inexact; README and G-141 said 8 to 12 points where the
rollouts show 6 to 13. Fixed in `f37be4e`: a 13-point floor per run,
same-window subtraction clamped at 0, a missing `resets_at` read as
current, the README saying the cap is not a ceiling. The guardian rollout
picked by mtime was checked safe on all nine pairs.

Round 2 (`1e1a345..643d29d`), no blocking defect: a new home with no
reading skipped the floor on its first run; `resets_at` jitters by a second
within a window, so exact equality misjudged it; the selftest covered only
the happy path and the stop; G-143 understated the reads outside the clone
and named one fixture commit for all batches; G-144 said the binary embeds
only the guides. Fixed in `4513c69`, with selftest checks for the new-home
floor, a reset window and a missing `resets_at`. G-143's numbers, the
pattern per case and the Claude pair were verified against the evidence.

Round 3 (`4513c69`): every earlier finding resolved; nothing blocks.

## Disposition

Open, low, left at the review cap: no selftest for the stop after a run
that leaves no reading or for the clamp at 0; a run whose readings lack
`resets_at` is measured from its own first reading (conservative for the
gate, the floor bounds it, no retained reading lacks it); G-143 calls all
of `~/.cache/grove-evals` the eval's output where one search was of the
Codex home; one README line is longer than its neighbours.
