---
id: "Q-001"
type: question
title: How should the board present differing versions of one record?
status: open
blocks: ["W-004", "W-005"]
relates_to: ["W-001", "W-003"]
created: "2026-09-19T14:08:40Z"
updated: "2026-09-19T17:50:16Z"
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

This blocks the proposed version inspection and workspace-location contracts
in W-004 and W-005. It does not block W-003 or existing single-checkout
inspection. Claims and run lifetimes need their own design when execution
enters scope.

## Proposed answer for review

The owner selected cross-branch coordination as the next experience on
2026-09-19. The following policy remains proposed; selection of the experience
does not itself approve these rules.

Group observations by record ID and retain each source explicitly. When main
says `proposed` and a feature checkout says `done`, show both statuses. Do not
invent one authoritative status or infer integration from the feature's state.
Identical bytes may be summarized together, but every source remains selectable.
Grouping is a navigation aid, not proof that independently imported matching
IDs share an origin. Do not merge records automatically.

Three alternatives:

- Explicit versions (recommended): exposes divergence and makes the editing
  destination deliberate. Costs an extra selection when sources differ.
- Main as headline: gives a familiar baseline, but hides branch progress in the
  default summary and assumes a branch should govern all work.
- Latest `updated` as headline: reduces selection, but direct editors need not
  maintain dates and a recent edit does not establish authority or integration.

Start with local branch tips and registered live worktrees in the same Git
repository. Branch snapshots and live files are separate observations. Read each
source's own configuration at the selected project's repository-relative path.
Do not combine relationships across sources or let a missing dependency on one
branch resolve from another. Remote-tracking refs, tags, historical timelines,
cross-clone reconciliation, and ownership claims are outside this first view.

Keep committed observations visible when live records differ or disappear;
report the live difference or absence. Invalid or inaccessible sources remain
visible as diagnostics, with the combined result marked incomplete. Never
silently substitute committed bytes for invalid live data. Existing local
inspection commands keep their current strict validation behavior.

Opening an existing workspace requires explicit source identity and freshness
checks. If the selected committed record differs from the live file, show that
difference and require selection of the live observation before returning it
as the editing context. A missing checkout is a distinct outcome; it does not
authorize switching another directory or creating one implicitly.

Reconsider the display policy if real usage makes explicit version selection
too noisy. A default view may summarize identical content while preserving
source selection; timestamps alone must not establish authority.

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
The current repository has only main and one checkout; this inspection does
not supply new evidence for conflicting branches. The earlier fixture in the
brief establishes only basic routing feasibility.

## Next

Review the explicit-version policy with the owner, then reconcile this answer
with the W-004/W-005 specifications and retain this question as resolved.
Implementation fixtures must exercise conflicting edits, dirty and missing live
records, branch changes after selection, and a branch without a checkout.
