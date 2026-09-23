---
id: "G-116"
type: plan
title: "G-109 Attempts layouts: observed screens, options and plan"
status: current
created: "2026-09-23T19:51:06Z"
updated: "2026-09-23T19:53:31Z"
work: ["G-109"]
---

## Inputs

Prepared on 2026-09-23 for [G-109](G-109-attempts-usability.md) at revision
`sha256:5a72acaa…`, headless, on branch `worktree-G-109` from `main`
`6208e82`. Read in full: G-109, G-046 (its evidence, limits and the owner's
verdict), the G-096 visual-proposal precedent, `internal/tui/attempts.go`,
`internal/attempt/facts.go`, `internal/attempt/activity.go`, the `View`,
`Launch` and `Result` types, the attempt tests, `docs/board.md`'s Attempts
section. No layout is approved; [G-117](G-117-which-attempts-list-and-detail-l.md)
asks the owner to choose, and nothing below is implemented until it is
answered.

## Observed at `6208e82`

Rendered by a throwaway model test (deleted, never committed) that fed this
repository's six real attempts, read with `attempt.ListDir`, `ShowContext`
and `ReadActivity`, into the current screens at 80×24 and 120×36. Home paths
are shortened to `~` here.

The list at 80×24 (the whole attempt set):

```text
Attempts in this repository  (6)
> G-109.20260923T194926Z  running  worktree-G-109
  G-108.20260923T194708Z  running  worktree-G-108
  G-107.20260923T182316Z  candidate ready: G-107 in review on worktree-G-107 wi…
  G-107.20260923T164723Z  candidate ready: G-107 in review on worktree-G-107 wi…
  G-105.20260923T044923Z  candidate ready: G-105 in review on worktree-G-105 wi…
  G-046.20260923T033237Z  candidate ready: G-046 in review on worktree-G-046 wi…
↑↓  Enter show  x stop  o record  Esc back  q quit
```

The top of one finished attempt at 80×24 (G-107's second):

```text
Attempt G-107.20260923T182316Z of G-107
Outcome: candidate ready: G-107 in review on worktree-G-107 with candidate
1614e89
Attempt: G-107.20260923T182316Z (finished)
Work: G-107 at grove/G-107-current-documentation.md, record
sha256:c411c9a83fb1ded4ab77c90338c70f75db52b056098e1ce5a8b5d868f7fe70bc
Worktree: ~/grove/.claude/worktrees/worktree-G-107 on
worktree-G-107 from 4159e79ef76d (reused)
Started: 2026-09-23T18:23:16Z by grove v0.0.0-20260923043516-29f0ceaaa9b3
29f0ceaaa9b382a120af2aa61bcd9c533e8290e9 with claude (2.1.280 (Claude Code)),
owner pid 25773
Bounds: budget 50 USD, permission mode auto, prompts none; one process, no
retries; subagents share the budget
Session: 95b5aaee-8a73-4e51-bde5-5effcb3d05a0
Command: ~/.local/bin/claude -p /grove-work G-107 --interaction
headless --output-format stream-json --verbose --session-id 95b5aaee-8a73-4e51-
bde5-5effcb3d05a0 --max-budget-usd 50 --permission-mode auto --permission-
prompts none
Events: 593 lines, 1502057 bytes (system 270, assistant 205, user 106,
rate_limit_event 8, result 4); unknown 0, malformed 0, oversized 0, partial
↑↓ scroll  x stop  o record  Esc back  q quit
```

A running attempt (G-108's), second page at 80×24:

```text
Activity, newest first (151)
  Now the fixture files.
  system: thinking_tokens
  system: thinking_tokens
  tool: Bash go run ./cmd/grove update G-108 --expect sha256:ff01c9f092aa3f43cf…
  tool: Bash python3 - <<'EOF'
  tool: Bash cat grove/G-115-g-108-eval-skeleton-plan.md; go run ./cmd/grove up…
  Plan record G-115 created; writing its design now.
  tool: Bash go run ./cmd/grove new plan "G-108 eval skeleton plan" 2>&1;
  system: thinking_tokens
  system: thinking_tokens
  … (the next two pages are only system: thinking_tokens)
```

Findings, each tied to an acceptance item of G-109:

1. **No work title anywhere** (acceptance 1). Rows lead with a timestamped
   attempt ID; the work's title is only one `o` away.
2. **Stale attention** (1). Four rows say `candidate ready`, yet G-107,
   G-105 and G-046 are all `done` on `main` now, and G-107's first attempt's
   candidate `71a650e` was superseded by the second attempt's. The outcome is
   what the owner process found at exit, correctly, but nothing says whether
   it still needs the owner.
3. **Diagnostics before the result** (1, 2). A finished attempt's report
   starts on the second page at 80×24, after 18 wrapped lines of provenance;
   a running attempt's activity starts at the bottom of the first page.
4. **Activity noise** (2, 3). Claude Code 2.1.281 emits a `system:
   thinking_tokens` event per thinking step; a little later G-108's log held
   181 system events other than `init` among 237, so useful activity
   scrolls away.
5. **Truncated outcome** (3). At 80 columns the outcome cuts off before the
   candidate, and the cost and branch fall off the row; at 120 columns the
   branch is still cut.
6. **No elapsed time** (1). A running row gives no sign of how long it has
   run; a finished one says neither when nor how long.

Observed, outside this work: the board's ASCII Markdown style renders
`**New `docs/board.md`**` as `**New **docs/board.md` (the strong markers
misplace around a code span, backticks lost). It affects record bodies as
much as reports, so it is board polish for its own proposal, not G-109.
Report retention needs no repair: the report is the result event's text,
the last event of a run: all four finished logs here end with it,
including three over 1 MiB (G-046's 2.6 MB, G-107's 1.6 and 1.5 MB), so it
sits inside the 1 MiB window; G-046's recorded limit stands.

## Representative states

Real attempts cover running (G-108, G-109) and candidate ready whose work
is now done (G-107, G-105, G-046), including a superseded candidate
(G-107's first). The other states come from `TestAttemptOutcomes`' cases,
drawn with sanitized IDs `X-1`… and invented titles, so no mockup claims a
real run that did not happen: waiting on a question, failed (budget
exhausted; killed, no result event), stopped, orphaned, interrupted, ended
without a handoff, and candidate ready still awaiting judgment.

## Option A: attention first, details folded (recommended)

List, grouped, newest first within a group. `Needs you` holds only the
latest attempt of work that is still proposed, active or review, and only
when its state asks the owner for something; `Running` holds running
attempts; `Settled` holds the rest, each saying why nothing is needed.
80×24:

```text
Attempts in this repository (11) · 5 need you · 2 running · 4 settled
 Needs you
> X-1    Tidy the export command        judge candidate c0ffee1       2h ago
  X-2    Choose a colour scheme         answer question X-9           3h ago
  X-3    Speed up the loader            failed: budget exhausted      5h ago
  X-4    Rename the settings file       orphaned: x stops it         20m ago
  X-5    Split the parser               ended, no handoff             1d ago
 Running
  G-109  Make Attempts easy to scan a…  running                      14m so far
  G-108  Establish behavioral evaluat…  running                      16m so far
 Settled
  G-107  Reconcile current documentat…  done: candidate 1614e89        1h ago
  G-107  Reconcile current documentat…  candidate 71a650e, superseded  3h ago
  G-105  Restore green CI: Linux buil…  done: candidate e30f90c       15h ago
  X-6    Retry flaky probe              stopped by x                  2d ago
↑↓ move  Enter show  o work  x stop  r refresh  Esc back  q quit
```

The title column shrinks first; below 60 columns it goes, leaving ID, state
and time, and the state is cut last. Enter still opens the attempt itself,
whose full ID leads its screen, so identity stays exact.

One attempt: who, state, what to do, then the result or the latest
activity; the facts the CLI's `attempt` prints are one key away, unchanged
and complete. Candidate awaiting judgment, 80×24:

```text
X-1  Tidy the export command
Attempt X-1.20260923T160000Z · finished 2h ago after 13m · $6.03 of $50
State  Candidate c0ffee1 is in review on worktree-X-1; you have not judged it.
Next   o opens X-1: a approves, f gives feedback

Final report
 X-1 is in review with candidate c0ffee1 on branch worktree-X-1. …
 …
Details: bounds, provenance, events and raw files · d shows them
↑↓ scroll  d details  o work  Esc back  q quit
```

Running, 80×24, with system notices counted instead of listed (the counts
are G-108's log when this plan was written):

```text
G-108  Establish behavioral evaluations for Grove context and workflows
Attempt G-108.20260923T194708Z · running for 16m · budget $50 · auto
State  Running; nothing needs you until it ends. x stops it.

Latest activity, newest first (237 events; 181 system notices not listed)
  Now the fixture files.
  tool: Bash go run ./cmd/grove update G-108 --expect sha256:ff01c9f092aa3f43…
  tool: Bash python3 - <<'EOF'
  tool: Bash cat grove/G-115-g-108-eval-skeleton-plan.md; go run ./cmd/grove …
  Plan record G-115 created; writing its design now.
  tool: Bash go run ./cmd/grove new plan "G-108 eval skeleton plan" 2>&1;
  …
Details: bounds, provenance, events and raw files · d shows them
↑↓ scroll  d details  x stop  o work  Esc back  q quit
```

`d` expands the details in place, where the line stands, as today's facts
verbatim, the raw `events.jsonl` and `stderr` paths included. The other
states' State and Next lines:

| State | State line | Next line |
| --- | --- | --- |
| waiting on a question | Waiting on open question X-9 (its title). | o opens X-2, whose detail lists the question |
| failed | Failed: error_max_budget_usd, exit 1, after $2.00 of $2. Uncommitted changes remain on worktree-X-3. | d shows the details and raw log; R on X-3 launches again |
| failed, no result | Failed: no result event (killed). | as above |
| orphaned | Its owner is gone and the provider still runs. | x stops it |
| interrupted | The owner and the provider are gone without a result. | d shows the raw log; R on the work launches again |
| ended without a handoff | Ended cleanly, but X-5 is active on worktree-X-5 with no candidate. | o opens X-5; the report says why |
| stopped | Stopped by x (exit 130); partial work kept on worktree-X-6. | none |
| done now | Candidate 1614e89 was integrated: G-107 is done. | none |
| superseded | A later attempt of G-107 followed this one. | Enter on the later attempt |
| uncommitted handoff | Its record says review with candidate c0ffee1, uncommitted, so no candidate is ready. | as ended without a handoff |

Work's detail row becomes `Attempts 2 · latest: judge candidate c0ffee1,
2h ago · A lists them · R launches one` (the same short state as the list).

## Option B: chronological, details below

The smaller change. The list keeps one newest-first order and the same row
(work ID, title, short state, time), marking rows that need the owner with
`!` rather than grouping them:

```text
Attempts in this repository (11) · 5 need you (!) · 2 running
  G-109  Make Attempts easy to scan a…  running                      14m so far
  G-108  Establish behavioral evaluat…  running                      16m so far
> X-4  ! Rename the settings file       orphaned: x stops it         20m ago
  G-107  Reconcile current documentat…  done: candidate 1614e89        1h ago
  X-1  ! Tidy the export command        judge candidate c0ffee1       2h ago
  …
```

One attempt has the same header, State, Next and report or activity as A,
then the facts always shown below them, with no `d`. A long report pushes
the raw paths down, so reaching them means `End` or paging.

## Decisions for the owner

G-117 asks them together: A or B; whether `Needs you` uses the rule above
(latest attempt of unfinished work only, so an old failure after a later
success, or any attempt of done work, never asks for attention); and
whether counting `system:` notices other than `init` instead of listing
them is acceptable (the raw log keeps them all).

## Steps once answered

For the chosen option; B drops steps 2's grouping and 3's toggle.

1. **State.** In `internal/tui/attempts.go`, derive a structured state
   (group, short label, State sentence, Next line) from the facts
   `outcomeOf` reads today plus the work's current record in the board's
   read (status, candidate) and the work's later attempts. `outcomeOf`'s
   derivations and their honesty rules stay; `TestAttemptOutcomes` gains
   each field per case and the done-now, superseded and later-attempt
   cases.
2. **List.** Title from the board's records (escaped like any record text,
   `(title unread)` when the work is on no readable branch), grouped rows,
   time column (`so far` from `Launch.Started` while running; `ago` from
   `Result.Finished`), widths at 40/60/80/160.
3. **One attempt.** Header, State, Next, report or activity, then `d`
   toggling `attempt.Facts` unchanged. Hidden system notices: `note` in
   `internal/attempt/activity.go` counts non-`init` system events instead
   of listing them, and the count is shown.
4. **Work detail row and help.** The row above, the hint bars, and
   `docs/board.md`'s Attempts section and key table (`d`).
5. **Checks.** Model tests for every state's rendering, escaping of titles
   and provider text, widths and scrolling with long titles, 200 attempts,
   long reports and 20,000 events; update `attempt_lifecycle` in
   `internal/tui/testdata/terminal.py` for the new text (fake provider,
   labelled so); the repository's verification; one independent review on
   the final revision; actual renders at 40×10, 80×24 and 160×48 recorded
   in G-109's evidence. The owner's judgment in a terminal (acceptance 5) is
   theirs, and no real-provider trial is run without a bounded mandate.
