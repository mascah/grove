---
id: "G-046"
type: work
title: "Launch and inspect managed attempts from the TUI"
status: active
created: "2026-09-21T00:54:16Z"
updated: "2026-09-23T03:35:48Z"
kind: feature
size: medium
priority: 4
depends_on: ["G-044", "G-045"]
relates_to: ["G-035", "G-043"]
formerly: "W-028"
---

## Outcome

Launch one eligible work item from the board, inspect its durable attempt,
return after closing the TUI, and review the resulting candidate in Grove.

## Scope and bounds

Implement is an explicit assignment action, not a side effect of changing
status. Show a runs list and attempt detail or an integrated detail section,
using G-045's ownership/events rather than making the TUI own the process.
Present progress, waiting/failure, budget and Stop with linked logs as needed.
Retain the final report and reviews when ephemeral activity scrolls away.

## Acceptance

1. Eligible selection creates exactly one bounded assignment; duplicate starts,
   missing readiness and changed inputs produce useful actionable outcomes.
2. Closing and reopening the TUI reconnects to the same attempt. Stop is a
   separate explicit action and preserves partial work/evidence.
3. Running, waiting, failed/interrupted and candidate-ready outcomes are honest;
   a successful process exit alone cannot fabricate Review readiness.
4. Review feedback can launch a new bounded attempt on the correct candidate
   branch, preserve prior evidence and return work to Active.
5. Owner-reviewed connected workflow and terminal lifecycle checks cover launch,
   observe, exit, reconnect, wait, stop, review and continuation. Test high-volume
   output without flooding the main interface or blocking navigation.

## Preparation and next

No implementation plan exists yet. Create a linked plan through the supported
CLI after G-045 and G-044. Decide whether attempt detail is a separate screen using
the demonstrated workflow. Multi-selection, batch scheduling and dependency
graphs remain later candidates with explicit budgets and dependency analysis.
