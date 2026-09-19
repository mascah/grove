---
id: "Q-001"
type: question
title: How should the board present differing versions of one record?
status: resolved
blocks: ["W-004", "W-005"]
relates_to: ["W-001", "W-003"]
created: "2026-09-19T14:08:40Z"
updated: "2026-09-19T17:54:10Z"
---

## Question

Main and a feature branch can hold different versions of one record. Which
version should the combined view display, and how does a user inspect or open
each version without silently editing the wrong checkout?

## Constraints and evidence

The [restart brief](../../docs/restart-brief.md#cross-branch-view-and-branch-context-editing-2026-09-18)
owns the accepted branch-context direction and the earlier routing experiment.
The answer below settles presentation and explicit version selection. Detailed
source and routing contracts belong to W-004/W-005 and their coordination plan.

This question previously blocked W-004 and W-005. Its resolved status removes
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

[W-004](../work/W-004-record-versions.md) owns the source scope,
committed/live representation, validation, output, and selector contract.
[W-005](../work/W-005-record-workspace.md) owns workspace lookup,
freshness checks, and missing/ambiguous-checkout outcomes. The owner's answer
does not approve every technical proposal in those work records. Those contracts
were subsequently implemented and integrated; the
[coordination plan](../../docs/plans/W-004-W-005-coordination.md) retains their
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
in the brief establishes only basic routing feasibility.

## Next

Retain this accepted answer for all future workspace presentation. W-004/W-005
are integrated at `5041ae1`; their CLI exposes explicit versions and resolves
existing workspaces. The [integrated review](../../docs/reviews/2026-09-19-integrated-cli.md)
records proposed reliability repairs before interactive actions. An interactive
view, automatic checkout creation, and agent launching remain unimplemented.
