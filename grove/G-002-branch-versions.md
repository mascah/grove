---
id: "G-002"
type: question
title: How should the board present differing versions of one record?
status: resolved
blocks: ["G-010", "G-011"]
relates_to: ["G-003", "G-009", "G-035", "G-042"]
created: "2026-09-19T14:08:40Z"
updated: "2026-09-21T01:03:42Z"
formerly: "Q-001"
---

## Current disposition, 2026-09-20

The owner approved [G-035](G-035-interactive-adoption.md), replacing
this question's explicit-versions-first presentation as the future default with
a project-wide current view. Preserve the answer below as the authority behind
the delivered CLI/board, not a constraint against that redesign. Exact source
inspection and fresh workspace targeting remain required. G-042 owns unresolved
projection details; current software still follows the earlier presentation.

## Question

Main and a feature branch can hold different versions of one record. Which
version should the combined view display, and how does a user inspect or open
each version without silently editing the wrong checkout?

## Constraints and evidence

The [restart brief](brief.md#current-view-and-tui)
owns the current branch-context direction. The earlier routing experiment is
preserved in Git at `c9904ea:grove/brief.md`. The answer below records the
original presentation choice. Detailed source and routing contracts belong to
G-010/G-011 and their coordination plan.

This question previously blocked G-010 and G-011. Its resolved status removes
that policy blocker; it does not establish implementation readiness or completed
verification. Claims and run lifetimes need their own design when execution
enters scope.

## Answer

On 2026-09-19 the owner selected: "Group under W-003, show both branch statuses,
and require a version selection to open a workspace."

Group observations by record ID and retain each source explicitly. When main
says `proposed` and a feature checkout says `done`, show both statuses. Do not
invent one authoritative status or infer integration from the feature's state.
Require explicit version selection before opening its workspace. Grouping is a
navigation aid, not proof that independently imported matching IDs share an
origin, and does not authorize merging records automatically.

## Alternatives considered

- Explicit versions (selected): exposes divergence and makes the editing
  destination deliberate. Costs an extra selection when sources differ.
- Main as headline: gives a familiar baseline, but hides branch progress in the
  default summary and assumes a branch should govern all work.
- Latest `updated` as headline: reduces selection, but direct editors need not
  maintain dates and a recent edit does not establish authority or integration.

Reconsider the display policy if real usage makes explicit version selection
too noisy. A default view may summarize identical content while preserving
source selection; timestamps alone must not establish authority.

## Remaining design ownership

[G-010](G-010-record-versions.md) owns the source scope,
committed/live representation, validation, output, and selector contract.
[G-011](G-011-record-workspace.md) owns workspace lookup,
freshness checks, and missing/ambiguous-checkout outcomes. The owner's answer
does not approve every technical proposal in those work records. Those contracts
were subsequently implemented and integrated; the
[coordination plan](G-013-coordination-plan.md) retains their
finalized technical details without changing this settled presentation choice.

## Inspected evidence

At `ee69c42`, `internal/project/project.go` loads live files and validates one
complete source; `internal/project/graph.go` resolves relationships within that
set. Cross-branch inspection needs another source adapter, retaining the same
schema rules. `internal/create/create.go` scans IDs for allocation but is not a
record/version reader and must not define display precedence.

Git documents branch/ref enumeration in
[for-each-ref](https://git-scm.com/docs/git-for-each-ref), committed tree entries
in [ls-tree](https://git-scm.com/docs/git-ls-tree), and worktree branch, HEAD,
path, and detached state in
[worktree porcelain output](https://git-scm.com/docs/git-worktree#_porcelain_format).
These provide discovery inputs, not a transactional selection or a Grove claim.
At that historical inspection the repository had only main and one checkout;
it did not supply new evidence for conflicting branches. The earlier fixture
in the historical brief establishes only basic routing feasibility.

## Next

[G-042](G-042-current-view.md) specifies the new default projection.
G-010/G-011's explicit inspection/routing remain available; G-017/G-030's board
and history are integrated. The earlier reliability repairs also shipped.
Automatic checkout creation and managed agent launching remain future work.
Historical reviews and completed records retain their original acceptance scope.
