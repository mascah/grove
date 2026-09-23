---
id: "G-103"
type: plan
title: "G-046 managed runs plan"
status: current
created: "2026-09-23T03:35:15Z"
updated: "2026-09-23T03:35:47Z"
work: ["G-046"]
---

## Design

Prepared 2026-09-22 at `bd2418f` for [G-046](G-046-managed-runs.md)
(record revision `sha256:4dc799d6…`), on G-045's `internal/attempt`
package (plan [G-100](G-100-g-045-durable-attempt-plan.md), decision
[G-101](G-101-attempt-mechanism.md), review
[G-102](G-102-g-045-durable-attempt-review.md) and its notes for G-046) and
G-044's review actions in `internal/tui/review.go`. The plan owns the steps;
G-046 owns outcome, bounds and acceptance. The roadmap's open run-experience
questions ([G-047](G-047-adoption-roadmap-plan.md)) are design tasks decided
here.

### Decisions

- **Separate screens, plus a section.** The attempt is not the record: it
  has its own lifecycle, files and a Stop. Two screens: *Attempts* (every
  attempt of the repository, or one work's, newest first) and *Attempt*
  (one attempt's facts, outcome, final report and recent activity). The
  record detail of work gains one header row naming its latest attempt and
  the keys, so a person returning after closing the TUI finds the attempt
  from the card they already know, and board cards of work with a running
  or orphaned attempt carry a `● running` or `● orphaned` tag (`▶` already
  marks focus).
- **Keys.** `R` on the detail of proposed or active work opens the launch
  prompt; `A` on the board lists every attempt, on a work detail that
  work's; Enter on an attempt row opens it; `x` on an attempt or its row
  asks `y/n` and stops it; `o` on the attempt opens its work's record,
  where the reviews and evidence are. Upper case keeps the launch and list
  away from `a` (abandoned, approve) and `r` (refresh); every launch passes
  through a typed prompt, so no single key starts a process.
- **Launch prompt.** Budget, then permission mode, typed each time with no
  prefill: G-045 made both required because a default would make the
  choice in practice, and a prefilled field is that default. `attempt.
  ValidBudget` checks the budget before the next field; the mode is passed
  through as the CLI does.
- **Where it runs.** The launch runs `attempt.Start` in the board's root.
  When the shown record stands on a branch other than the target, that
  branch is passed, with its live checkout as the worktree when it has one:
  the continuation after feedback runs on the candidate's branch, wherever
  its worktree is (G-044's own lives at `.claude/worktrees/G-044`, not
  Start's default path). Otherwise Start's defaults apply.
- **Changed inputs.** `Request.Expect` carries the record revision the TUI
  showed for the launching checkout; Start refuses when the checkout's
  record no longer hashes to it, naming both, before anything is written.
  The CLI keeps no such flag (its reader is the person at the shell).
- **Honest outcome.** Derived for display, never written: running;
  orphaned; interrupted; and for a finished attempt, in this order, a
  candidate ready (the result's record is in `review` with a candidate),
  waiting on an open question that blocks the work in the current view,
  stopped, failed (no result event, `is_error`, or a nonzero exit, with the
  subtype), or ended without a handoff (the record still proposed or
  active). A clean exit alone never reads as ready.
- **Reads.** A new backend read, separate from the one-at-a-time Git read
  slot so it never cancels an inspection: `attempt.List` of the whole
  repository and, with the attempt screen open, `attempt.Show` of that one
  plus a bounded tail of its events. At most one in flight; collected at
  close with the other reads. It runs once the board is read, on opening
  either screen, after a launch or stop, on `r`, and every 2 s while an
  attempt is running or orphaned; nothing ticks otherwise. The list reads
  the attempts directory under the common directory the board found, so it
  starts no Git process; the open attempt's read starts one, cancellably. When an
  attempt known to be running is read finished, the board is re-read so
  the branch's record shows its new state. `List` stops scanning events
  (`attempts` never printed them); `Show` still does.
- **Activity and the final report.** `attempt.Activity(path, window)` reads
  only the last `window` bytes of `events.jsonl` (MaxLine plus 64 KiB, so
  a final result line always fits), drops the leading partial line, and
  keeps one short line per event: assistant text (first line), tool use
  (name), a tool error, system subtype, the result, plus the result's
  report text. The attempt screen shows facts, outcome, the report
  rendered as the record bodies are (escaped before glamour, filtered
  after), then the last activity lines that fit: the report stays above
  the activity, whatever scrolls. High volume costs one bounded read per
  tick, never the file.
- **Owner reaping** (G-102's note): `Start` waits on the owner in a
  goroutine instead of releasing it, so a long-lived launcher reaps it; a
  CLI launcher exits first and the owner is reparented as before.
- **Stop** runs `attempt.Stop` as an action behind `y/n`; keys wait for it,
  as for G-044's actions. Stop of an orphan can take its grace period.
- **Feedback continuation.** Feedback's result names `R` on the record as
  the next step besides `/grove-work`; the record is then active on its
  branch, and `R` launches there. Evidence stays: the attempt appends
  nothing but what the headless guide writes.

### Placement

- `internal/attempt`: `Request.Expect`, owner reaping, `List` without
  events, `Activity`.
- `internal/tui`: `attempts.go` (screens, reads, prompt, outcome); Backend
  gains `Attempts`, `Launch`, `Stop`; `run.go` wires them; the detail's
  header row; card tags; hints.
- `internal/tui/testdata/terminal.py`: a scenario with a fake provider
  (`GROVE_CLAUDE`) covering launch, observe, quit, reconnect in a new
  session, Stop, and a high-volume provider while keys still move.
- README and `docs/record-model.md` where they describe `grove` and
  attempts.

## Steps

1. Commit this plan; set G-046 active.
2. `internal/attempt` changes with tests (`-short` where a process spawns).
3. TUI screens, reads, launch and stop, outcome, detail row, tags; model
   tests with a fake backend for eligibility, duplicate and changed-input
   refusals, outcomes, feedback continuation, reconnect, the tick's bounds
   and escaping of provider text.
4. Pseudo-terminal scenario with the fake provider, and a manual exercise
   in a disposable clone with a built binary.
5. Docs; verification per AGENTS.md.
6. Independent review of the combined diff; fixes within the cap.
7. Evidence and handoff in G-046; `status=review` with the candidate.
