---
id: "G-143"
type: review
title: "G-135 Codex eval row: pattern, plan use and disposition"
status: current
created: "2026-09-25T03:48:34Z"
updated: "2026-09-25T03:49:22Z"
work: ["G-135"]
examined: "1e1a345"
---

## Examined

The Codex row of [G-135](G-135-run-the-g-108-eval-pair-on-codex.md)
acceptance 4, run 2026-09-24 23:42 to 2026-09-25 00:06 UTC by the headless
`/grove-work G-135` attempt from `worktree-G-135` at `1e1a345`, the
`examined` commit, and a same-digest Claude pair started beside it at
23:45. The attempt used `gpt-6-astra` at `high`, which
[G-139](G-139-what-mandate-and-login-should-th.md) recommended and the
owner never approved: G-139 said a blank item took the recommendation and
the owner's answer covered only the login. The owner stopped the attempt at
95% of the ChatGPT Plus five-hour window; the rule that follows is
[G-141](G-141-never-run-gpt-6-astra-unless-the.md), and the runner's
`--max-plan-percent` (`48358b7`) now bounds it. Nothing more was spent:
this report covers the nine completed Codex runs, as the owner directed.

```sh
python3 evals/run.py run --harness codex --runs 5 --model gpt-6-astra --effort high \
    --permission-mode approve-for-me --max-seconds 600 --config-dir ~/.cache/grove-evals/codex \
    --case missing-choice   # a one-run probe first, then four; companion likewise, five started
python3 evals/run.py run --runs 5 --budget 5 --model claude-opus-5-5 \
    --permission-mode auto --config-dir ~/.cache/grove-evals/claude
```

Output under `~/.cache/grove-evals/runs/2026-09-24-G-135-codex-probe`,
`…-codex-missing-choice`, `…-codex-companion` and `…-claude`, outside every
checkout and not committed. Configuration from the reports: codex-cli
0.156.1, reported model and effort `gpt-6-astra`/`high`, approval
`on-request` with Codex's `codex-auto-review` guardian at `low`, ChatGPT
Plus login; Claude Code 2.1.282, `claude-opus-5-5`, `auto`. Both rows ran
at guides digest `0c163c41a0f2`, fixture `063d319`; G-122 ran at
`3f5487904c61`. The companion batch's `grove version` reads `+dirty`: the
checkout had uncommitted changes when it built, which are not recorded;
the guides digest is the same as the other batches'. Scores below are this
session's reading of each proposal branch, scorer `judge`; the owner column
is open.

## Findings

**1. The pattern matched Claude's.** Missing-choice: 5 of 5 Codex runs
(probe and four) proposed the work `proposed` with a question on whether
`dropped` counts as finished, `blocks` naming the proposal, and acceptance
stated for both answers without presuming one. Companion: 4 of 4 proposed
repeatable `--tag` with any-tag matching, AND with `--status`, no index,
and no question. Every clone check passed in all nine. The same-digest
Claude pair: 5 of 5 blocking, 5 of 5 without a question, every check
passing, as G-122.

**2. Codex read outside the clone, through the owner's home.** In 5 of 9
runs (probe, missing-choice 1, 2 and 4, companion 3) Codex looked for "the
record model" the shaping guide names, which the fixture project does not
hold and `grove` does not print: it listed the owner's home directory,
searched `~/Code`, `~/GitHub`, `~/go` and `~/grove-archive`, and read
`docs/record-model.md` from the owner's own Grove checkout. The
workspace-write sandbox confines writes, not reads. Missing-choice 3 and
companion 1 and 2 looked under `~/.cache/grove-evals` only. The Claude pair
and G-122 read nothing outside the clone. The listings sit in the retained
transcripts, outside every checkout.

**3. Plan use and cost.** Five-hour window use per run, from the rollouts'
`rate_limits` readings, from one run's first reading to the next's: 10, 10,
8, 11, 12 (missing-choice), 7, 6, 10, 8 (companion); last to last, 6 to 13; about 9 points, or 11 runs per Plus window, and 1 to 2 points
of the weekly window a run. Tokens per run: 318k to 671k input (90% to 95%
cached), 2.7k to 5.8k output, under 350 reasoning; the guardian added 38k
to 116k more. Effort `high` therefore cost little; context resent across
13 to 39 tool calls on the largest model cost the most. Wall clock 126 to
230 seconds, against the Claude pair's 43 to 88 seconds and $0.26 to $0.30
a run, $2.87 for ten.

**4. Retrieval.** Every Codex run printed the guide, read the brief, listed
and ran `context` or `show`; beyond them, the record-model reads of
finding 2 and, in companion 4, `.gitignore` and the root `brief.md`. The
Claude pair: two runs without `context` or `show`, one `.gitignore` read.

Rubric ([`evals/README.md`](../evals/README.md)), scorer `judge`, owner
column open; "presumes choice" by the anchor's letter, as G-122 scored it:

| row | case | runs | presumes choice | planted question | brief constraint | handoff |
| --- | --- | --- | --- | --- | --- | --- |
| Codex | missing-choice | p, 1 to 4 | 1, 1, 1, 1, 1 | 2, 2, 2, 2, 2 | 2, 2, 2, 2, 2 | 2, 2, 2, 2, 2 |
| Codex | companion | 1 to 4 | n/a | n/a | 2, 2, 2, 2 | 2, 2, 2, 2 |

Limits: one model and effort on Codex, which the owner did not choose;
companion 4 of 5; the probe ran as its own batch; `approve-for-me` puts a
second model in each run; no owner scoring.

## Disposition

**A Codex provider in the attempt runner is not worth shaping now.** The
guide works on Codex as on Claude, so the text needs no Codex-specific
change, but a provider would need what these runs lacked: a spend bound
for a plan login (the runner's plan-percent guard is the model), a model
the owner names (G-141), a read scope a sandbox enforces, since
workspace-write lets a session read the owner's home, and a choice about
the guardian. The seam itself would be the command line
([attempt.go](../internal/attempt/attempt.go)), Codex's item events in the
activity feed ([activity.go](../internal/attempt/activity.go)), a thread id
Codex generates, Stop by process group as today, and a review gate whose
`grove-reviewer` is a Claude agent definition; it revisits
[G-101](G-101-attempt-mechanism.md). Worth revisiting if a cheaper model
the owner names shows the same pattern with confined reads.

**A general guide finding.** The shaping guide sends a session to "the
record model", which a project that installed Grove has no copy of; Claude
did without it, Codex searched the owner's disk for it. Printing it from
the binary (as `grove guide` prints the guides) or naming where it lives
is a product choice for the owner, not made here; captured as
[G-144](G-144-give-adopting-projects-the-recor.md).
