---
id: "G-119"
type: review
title: "G-108 eval skeleton review"
status: current
created: "2026-09-23T20:11:31Z"
updated: "2026-09-23T20:11:48Z"
work: ["G-108"]
examined: "1f03a12"
---

## Examined

An independent review of the offline half of [G-108](G-108-workflow-evals.md):
the eval runner, fixture, checks and rubric under `evals/`, built on
`worktree-G-108` to plan [G-115](G-115-g-108-eval-skeleton-plan.md). A
separate reviewer agent (Claude Code 2.1.281 subagent, read-only, told not to
edit, commit or run the real `claude`) examined `4b5b345`, then the fixes at
`399cc5f` and `1f03a12`; `examined` is the last. Its evidence came from
running `python3 evals/run.py selftest`, fake runs through `--claude`,
mutations of the checks in throwaway clones, and signals sent to the runner.
No paid run happened: that waits for [G-118](G-118-what-mandate-should-the-g-108-pa.md).

## Findings

Round 1, on `4b5b345`, eight findings:

1. High: Ctrl-C or a crash left the paid `claude` process running.
2. High: the retrieval facts matched words after `grove` anywhere in a
   command, so a `new` title containing "tasks list", or a `grove-evals-`
   path, set flags for subcommands that never ran.
3. Medium: the frontmatter parser crashed the run on `blocks: [G-002]` and
   block-style lists, which `grove check` accepts, losing the run's record.
4. Medium: `report.md` could not tell a harness failure (no login, budget
   stop, timeout) from an agent that proposed nothing.
5. Medium: the fixture brief's "`closed` (the date it left `open`)" hinted
   that finished means not open, weakening the planted choice.
6. Low to medium: the selftest could not see four checks break.
7. Low: versions and the fixture commit were only in the report header; the
   CLI built in a linked worktree stamps main's revision, not HEAD's.
8. Low: records in nested folders missed; the baseline was the clone's
   current `main`; any branch commit satisfied the handoff check; OAuth
   tokens were stripped; an `--out` under a `CLAUDE.md` would leak it.

Round 2, on `399cc5f`: 1 to 8 confirmed fixed, by running them; three new:
SIGHUP still left the session running; a state failure after a paid run
dropped its cost and hid a deleted `main`; `--project=DIR` and
backslash-continued `grove` commands were missed.

Round 3, on `1f03a12`: all three confirmed fixed by running them, each of
the nine checks caught by the selftest when forced to pass, and no remaining
findings. One optional nit: a state failure also drops the retrieval facts,
which depend only on the transcript.

## Disposition

- 1: `try/finally` kills the session's process group on every exit; SIGTERM
  and SIGHUP unwind like Ctrl-C (`1f03a12`).
- 2: only a command whose program is `grove` counts, with its subcommand as
  the next word after `--project`; the known misses (`timeout`, `$(…)`) and
  the heredoc false positive are documented in `evals/README.md`.
- 3: tolerant parsing of both list forms; a runner failure is recorded per
  run with its cost and the next run starts.
- 4: a harness column (exit, timeout, error result, denials) in the report.
- 5: the hint removed; `tasks.py` still closes both statuses through one
  function, code rather than intent, and G-118 asks the owner to approve the
  fixture.
- 6: a `worse` fake mode that pushes, promotes, breaks `check`, names a stale
  commit and leaves two proposal branches.
- 7: versions and commits in every `run.json`; the stamp labelled as able to
  lag, the guide digest as the pin.
- 8: `ls-tree -r`, the pre-run `main`, the tip only, `CLAUDE_CODE_OAUTH_TOKEN`
  passed through; `--out` documented rather than constrained.
- Round-2 findings: fixed in `1f03a12`. The nit is not taken; the transcript
  is retained, so the facts can be recomputed.

This record is evidence, not approval, and says nothing about the paid runs
or acceptance 5.
