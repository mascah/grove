---
id: "W-027"
type: work
title: "Run one bounded implementation independently of the viewing terminal"
status: proposed
created: "2026-09-21T00:54:16Z"
updated: "2026-09-21T01:03:41Z"
kind: feature
size: large
priority: 4
depends_on: ["W-020", "W-022"]
relates_to: ["D-004", "W-026", "W-028"]
---

## Outcome

One explicitly assigned implementation runs independently of the viewing
terminal, leaves a durable result, and can be inspected, reconnected or stopped
without accidentally starting it again.

## Scope and bounds

Reuse the interactive assignment/preparation/review contract. Begin with one
work item and one attempt, then bounded independent review where the workflow
requires it. No schedules, dream loop, automatic scope expansion, arbitrary
multi-item fan-out or automatic merge. Keep ordinary CLI access daemon-free.

Define attempt identity, source/input revisions, working directory, process
owner, budgets, permission profile, events/logs, checkpoint, result, wait,
cancellation and interruption/recovery. Persist raw output separately from
durable review evidence. Record actual model/role configuration. Bound subagents
and retries explicitly. Process exit and streamed claims are not acceptance.

## Acceptance

1. Compare a Grove-owned `claude -p` process with available native Claude
   background sessions; select the smallest mechanism that meets this contract.
   Evaluate tmux only as an option. Pin capabilities to the installed version.
2. Demonstrate competing-start refusal, TUI/terminal exit survival, reconnect
   without replayed work, Stop after reconnect and explicit owner-loss handling.
3. Bound and handle partial/unknown/oversized events, provider errors, budget
   exhaustion, stale inputs and missing results. Preserve partial code/evidence.
4. A missing human decision persists a question and wait; an unchanged wait
   cannot trigger an uncontrolled retry. A result identifies exact revisions.
5. Exercise fake-process lifecycle failures first and one explicitly bounded
   real-provider trial. Distinguish observed recovery from untested shutdown
   behavior; machine-reboot recovery is not a selected requirement.

## Preparation and next

No implementation plan exists yet. Read the runner research linked by the
roadmap at this boundary, inspect the current provider CLI, and create a linked
plan through the supported Grove CLI. Technical prerequisites are W-020/W-022;
preferred investment is after the useful review experience in W-026. Worktree
creation must have explicit destination, provenance and failure handling.
