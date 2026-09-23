---
id: "G-109"
type: work
title: "Make Attempts easy to scan and act on"
status: active
created: "2026-09-23T16:05:07Z"
updated: "2026-09-23T20:10:09Z"
relates_to: ["G-045", "G-046", "G-107", "G-110"]
---

## Outcome

Let a person returning to Grove quickly understand which work an attempt
concerns, what happened or is happening, whether attention is needed, and what
they can do next, while retaining access to execution evidence.

Owner intent, shaping conversation 2026-09-23: the Attempts screens contain too
much information in an insufficiently usable presentation. Improve them as
part of preparing a small external preview. The layout below is proposed.

## Constraints

Observed at main `f27444e`: [attempts.go](../internal/tui/attempts.go) renders
list rows from attempt ID, outcome, optional final cost and branch. Detail
renders the outcome, every line of [CLI diagnostic facts](../internal/attempt/facts.go),
then the final report and recent activity. This puts session IDs, commands,
paths and event statistics ahead of the report.
[G-046](G-046-managed-runs.md) owns the implemented launch/reconnect/stop
workflow and its evidence. It records terminal tests with a fake provider;
its owner verdict accepted the views while leaving real TUI launch testing
for fresh work. Do not turn those observations into a claim of a real trial.

Proposed hierarchy: work title and honest outcome; required attention and
available actions; final result or recent progress; then expandable diagnostic
details. Compare concrete list and detail layouts using representative real
or sanitized attempts and the current keyboard workflow before selecting a
design. Preserve stable attempt identity even when human-readable titles lead.

Reuse the current process owner, bounded event reads, outcome derivation,
explicit launch/stop prompts, exact source targeting and stale-input checks.
Process exit, committed candidate readiness and human acceptance remain
different facts. Keep provider text escaping and the rendering filters intact.
Retain raw evidence access; disclose missing, stale or truncated evidence.
No generated activity summaries requiring an additional model call.

Scope covers Attempts list/detail and the directly connected work-detail
navigation/help. Small UI improvements in that flow belong here when tied to
a concrete task; unrelated board polish needs its own example and proposal.
No new execution provider, scheduling, automatic retry, bulk action, permission
default or lifecycle change. A broader report-retention repair, if preparation
finds one necessary, must be surfaced as a scope choice rather than hidden
inside presentation work.

## Acceptance

1. The owner can identify the work, state, attention needed and available next
   action in list/detail without scrolling through diagnostic metadata first.
   Demonstrate running, waiting on a question, failed, stopped/interrupted,
   orphaned, candidate-ready and ended-without-handoff examples.
2. A final report or useful current activity is easy to reach. Full provenance,
   budget/permission facts and raw evidence remain discoverable; absent or
   stale data is labelled rather than filled with inferred success.
3. Navigation remains usable with long titles, many attempts, long reports,
   high-volume output, empty/error states and narrow terminals. Record the
   tested terminal sizes and keyboard paths; inspect actual rendering.
4. Connected-workflow and terminal-lifecycle checks cover launch, observe,
   leave/reconnect, stop, inspect result, and return to the correct work and
   candidate. Preserve duplicate-start refusal, bounded reads and escaping.
5. The owner judges the revised screens in a terminal against the tasks above.
   Automated checks are separate evidence. Any real-provider trial requires
   an explicitly bounded mandate, and fake-provider tests are labelled.

## Evidence

Headless attempt of 2026-09-23 (`G-109.20260923T200725Z`), branch
`worktree-G-109` from `main` at `1a56af3`. It started from this record at
`sha256:254e7f7f…`, plan [G-116](G-116-g-109-attempts-layouts-observed.md)
at `sha256:97cac3d8…` and [G-117](G-117-which-attempts-list-and-detail-l.md)
as the owner answered it in this checkout. The owner's answer is committed
unchanged in `559f898`. The plan's section "Revised after G-117's answer"
(`4f721e1`) records the design that was built, with one choice the owner
left open: the left column holds State, Next and the final report.

Commits: `0315df9` (implementation), `611f7ec` and `feb6b42` (review
fixes). What changed:

- **The list** is G-117's option A. It has three groups, `Needs you`,
  `Running` and `Settled`, newest first within each, with G-117 part 2's
  rule. Each row gives the work ID, the title (dropped below 60 columns), a
  short state (cut last) and a time (`14m so far`, `2h ago`). Every settled
  row says why it is settled, such as `done: candidate 1614e89` or
  `candidate 71a650e, superseded`. The derivation is `standingOf` in
  `internal/tui/attempts.go`, on top of the unchanged `outcome` rules.
- **One attempt**, following the owner's direction:
  - The header has the work ID and title, a coloured state badge (`◆` needs
    you, `●` running, `✓` done, `·` otherwise) and the time.
  - A run panel follows, with coloured values: attempt, model and provider
    version, budget with the cost spent, permission mode, branch and base,
    start and end.
  - Metric chips: turns, tokens, context against the window with a gauge,
    subagents, compactions, tools and errors.
  - Then `State` and `Next`, and the details (`attempt.Facts`, unchanged)
    behind `d`.
  - From 100 columns, the final report sits on the left and the activity
    timeline on the right. Narrower, they stack.
  - Timeline rows show the local `HH:MM:SS` from the event's own
    timestamp, or blanks when it has none. Consecutive repeats collapse to
    one row with a count, such as `· system: thinking_tokens (×14)`, and the
    count is never cut.
- **Reader** (`internal/attempt/activity.go`): `Activity.Entries` carry
  kind, time, text and count. `Activity.Metrics` is provider-neutral: the
  TUI never reads Claude's events, and another harness's reader fills the
  same fields. It uses the same bounded window, 1 MiB + 64 KiB.
  - Tokens are the result's `modelUsage` totals. While the run is going,
    the input is summed and the output is `–`.
  - Turns sum the result events' counts, `≈` until a result ends the run.
  - Subagents are distinct subagent task IDs.
  - A count from a cut window shows `≥`, and an unknown value `–`.
  - A failed read keeps the last good read under the failure line.
- **Work detail row**: `Attempts N · latest: SHORT STATE, TIME · A lists
  them · R launches one`. The attempt screen gets a `d` hint, and
  `docs/board.md`'s Attempts section and key table are rewritten.

Against the acceptance:

1. Every named state is derived, and `TestAttemptStandings` checks group,
   short state, State and Next for each. The states are running, question,
   failed, interrupted, orphaned, candidate awaiting judgment, ended
   without a handoff and stopped, plus done now, superseded, a later
   attempt, feedback given, a candidate moved on by hand and a title that
   was not read. Real renders from this repository's seven attempts, 80×24,
   list:

   ```text
   Attempts in this repository (7) · 1 need you · 1 running · 5 settled
    Needs you
   > G-108  Establish behavioral evalu…  answer question G-118              20m ago
    Running
     G-109  Make Attempts easy to scan…  running                         25m so far
    Settled
     G-109  Make Attempts easy to scan…  ended, no handoff, superseded      38m ago
     G-107  Reconcile current document…  done: candidate 1614e89             1h ago
     G-107  Reconcile current document…  candidate 71a650e, superseded       3h ago
   ```

   One attempt at 80×24: the report starts on the first page, where it used
   to start on page two after 18 lines of provenance.

   ```text
   G-107  Reconcile current documentation and give each fact one owner
   ✓ done: candidate 1614e89  1h ago

     Attempt  G-107.20260923T182316Z
     Model    claude-opus-5-5 · 2.1.280 (Claude Code)
     Budget   $6.03 of $50 · permission mode auto
     Branch   worktree-G-107 from 4159e79ef76d (reused)
     Started  12:23:16 Wed 23 Sep · ended 12:36:35 after 13m
     Turns ≥64   Tokens 10.6M in · 93k out   Context 165k of 1M ▰▰▱▱▱▱▱▱▱▱ 17%
     Subagents ≥3   Compactions ≥0   Tools ≥90   Errors ≥2

   State  Candidate ready: G-107 in review on worktree-G-107 with candidate
          1614e89. Candidate 1614e89 was integrated: G-107 is done.
   Details: bounds, provenance, events and raw files · d shows them

   Final report
    G-107 is back in **review** with a new candidate, 1614e89, …
   ```

   The real attempts cover running, question, done, superseded and ended
   without a handoff. Every other state is from fixtures only.
2. The report or the latest activity is on the first page, and the full
   facts, with the raw `events.jsonl` and `stderr` paths, are behind `d`.
   `TestAttemptScreenHonesty` covers the absent and stale labels: model from
   the exit scan or "not in the part of the log read", "spend unknown: no
   result event", `≈`/`≥` turns, and a failed read. The reviewer checked
   turns and subagents against Claude's own totals on every real log.
3. Rendering: `TestAttemptListFits` covers 200 attempts with long wide-glyph
   titles and mixed groups at 40×10, 59×24, 60×24, 80×24 and 160×48.
   `TestAttemptActivityIsBounded` covers 200 wide-glyph entries, a long
   report and metrics at widths 40, 80, 99, 100 and 160. Both assert every
   row is exactly w cells, with scrolling and Esc. The real attempts were
   also rendered at 40×10, 80×24 and 160×48 by a throwaway probe test
   (deleted, never committed) and inspected. The 160×48 screens put the
   report and timeline side by side. Keyboard paths tested: `A`, ↑/↓,
   Enter, `d`, PgDn, `x` with `y`/`n`, `o`, Esc, `R`.
4. The `attempt_lifecycle` scenario in `internal/tui/testdata/terminal.py`
   uses a fake provider, labelled so, and its new text is updated. It
   covers launch, observe 20,000 events, quit and reconnect, stop, the
   question wait and launch refusal, the continuation to a candidate, and a
   return to the right work. Duplicate-start refusal, bounded reads and
   escaping are unchanged. `TestAttemptScreensReconnectAndStop` asserts no
   raw escape sequence reaches the terminal.
5. This is for the owner, in a terminal. No real-provider trial was run.

Verification at `feb6b42`: `go vet ./...`, `gofmt -l .`, `go run
./cmd/grove check` (OK: 112 records), `go test -count=1 -timeout 120s ./...`
and `python3 internal/tui/testdata/terminal.py` (all 10 scenarios), all
passing. At `0315df9` one full run failed `TestOwnerLost` (`Events.Init`
nil). That is a timing race in an existing test of the fake provider; this
change does not touch it, and the package passed on rerun. Under `-short`
in a parallel suite, `internal/versions` (8.3 s) and `internal/cli` (6.1 s)
run over the 5 s limit. This change touches neither.

Review: [G-120](G-120-g-109-attempts-redesign-review.md), by an independent
reviewer subagent in two rounds. The eight round-1 findings (4 should-fix,
4 nits) were fixed in `611f7ec`, and round 2 found none. Its two optional
nits were fixed in `feb6b42`, which was not reviewed again.

Limits:

- The left column's content, the colours and the glyphs are this attempt's
  choice within the owner's direction, for the owner to judge.
- Compactions are counted from Claude's `compact_boundary` notice, which
  appears in none of the real logs.
- Turns and subagents follow Claude's stream only. Another harness needs
  its own reader for `Metrics`.
- Deep in a long list, the group and page headers scroll off, as the page
  header did before.

## Next

In review: candidate on branch `worktree-G-109` (the `candidate` field),
base `main` at `1a56af3`. The owner judges the screens in a terminal
(acceptance 5), for example with `go run ./cmd/grove` in this worktree, then
`A`, Enter, `d`. They judge especially the attempt screen's left column,
colours and glyphs, which G-117 left open. Then, in this worktree:

```sh
go run ./cmd/grove approve G-109 "VERDICT"   # or: go run ./cmd/grove feedback G-109 "TEXT"
```

and in the `main` checkout:

```sh
go run ./cmd/grove integrate G-109
```
