---
id: "G-109"
type: work
title: "Make Attempts easy to scan and act on"
status: proposed
created: "2026-09-23T16:05:07Z"
updated: "2026-09-23T16:08:00Z"
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

## Next

Waiting on question [G-117](G-117-which-attempts-list-and-detail-l.md): the
owner chooses a layout. Nothing is implemented.

Checkpoint, 2026-09-23, headless attempt `G-109.20260923T194926Z`:
selected G-109 alone, on branch `worktree-G-109` from `main` `6208e82`,
starting from this record at `sha256:5a72acaa…`. Preparation is done:
plan [G-116](G-116-g-109-attempts-layouts-observed.md) records the current
screens, rendered from this repository's six real attempts by a throwaway
model test that was deleted and never committed. It lists six findings,
draws options A and B at 80×24 with every state G-109's acceptance names
(sanitized `X-` examples where no real attempt exists), and gives the
implementation steps. It also records that no report-retention repair is
needed. G-117 asks for the layout, the rule for what needs the owner, and
whether provider `system:` notices are counted rather than listed; it
recommends A with both rules. No commands are still running. Once G-117 is
resolved, rerun `/grove-work G-109`: the record goes active and G-116's
Steps apply to the chosen option.
