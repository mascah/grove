---
id: "Q-001"
type: question
title: How should the board present differing versions of one record?
status: open
blocks: []
relates_to: ["W-001"]
created: "2026-09-19T14:08:40Z"
updated: "2026-09-19T14:32:32Z"
---

## Question

Main and a feature branch can hold different versions of one record. Which
version should the combined view display, and how does a user inspect or open
each version without silently editing the wrong checkout?

## Constraints and evidence

The [restart brief](../../docs/restart-brief.md#cross-branch-view-and-branch-context-editing-2026-09-18)
owns the accepted branch-context direction and the earlier routing experiment.
Resolve version selection, visible source labels, live versus committed data,
and routing to the selected checkout before implementing a combined board.

This does not block inspecting one checkout, so `blocks` is empty. No combined
board work record exists yet. Claims and run lifetimes need their own design
when execution enters scope.

## Next

After local inspection works, exercise a record with conflicting edits in main
and a linked worktree. Include a worktree whose branch changes after selection,
a branch without a checkout, and a dirty record. Compare explicit version
selection with automatic precedence and record the chosen behavior and evidence.
Retain this record with its answer when resolved.
