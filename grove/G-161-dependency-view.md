---
id: "G-161"
type: work
title: "See work dependencies and preview a selected assignment"
status: proposed
created: "2026-09-25T20:35:15Z"
updated: "2026-09-25T20:38:13Z"
kind: feature
relates_to: ["G-035", "G-042", "G-043", "G-023", "G-047", "G-054", "G-060", "G-162"]
---

## Outcome

After several shaping sessions, the owner can return to Grove, see the
prerequisite structure across unfinished work, and choose what to do next or
assign together without asking an agent to reconstruct the intended order.

Owner intent, conversation 2026-09-25: uncertainty about dependencies among
roughly ten proposed items is preventing adoption in other hobby projects.
The owner suggested a visual DAG and recalled selecting several items for the
predecessor's work agent, while wanting more human control over that selection.
The owner then requested shaping this exploration. The design below remains
proposed; no layout or execution policy was selected.

## Constraints

### Observed evidence

At main `3f2b923`, the [record model](../docs/record-model.md) already defines
`depends_on`, `members`, question `blocks`, and priority separately.
[Validation](../internal/project/graph.go) rejects dependency and membership
cycles independently. [Context](../internal/handoff/context.go) orders selected
work by transitive prerequisites, including paths through unselected work,
without adding prerequisites to the assignment. The [detail
sidebar](../internal/tui/detail.go) exposes direct relationships in both
directions, but neither the board nor the CLI provides a dependency overview
or an interactive selection preview.

The main checkout contains 50 work records and 24 dependency edges; its largest
connected component contains 13 items. Only four records there are unfinished,
so this is layout evidence, not a representative ten-item proposed backlog.
`versions --json` also shows active G-154 on its work branch while main still
holds its proposed copy: the current view matters here. The [adoption
roadmap](G-047-adoption-roadmap-plan.md) names dependency visualization as a
later candidate; no existing proposal owns it.

### Proposed design and scope

- Add a focused dependency view to the board. Start with unfinished work;
  make a selected item's upstream prerequisites and downstream dependants
  easy to follow, with completed prerequisite history collapsed and
  expandable. Keep dependencies outside the visible filter discoverable and
  counted rather than making a truncated graph appear complete.
- Explore compact terminal lanes or layers with a readable list fallback.
  Arrows run from prerequisite to dependent work, with a visible legend.
  The owner judges concrete layouts before one is settled. A browser view
  or export remains an alternative if actual terminal readability warrants
  it, not a required second interface.
- Let the person select explicit work IDs and preview their dependency order,
  prerequisites outside the selection, and blocking questions. Expose the
  same interpretation through noninteractive inspection; choose the command
  and output contract during preparation. Browsing or selecting starts no
  attempt, adds no prerequisite to the assignment, and changes no status.
- Explain recorded facts relative to a named checkout/base and, when
  configured, the integration target: awaiting implementation, awaiting
  review, candidate present in the base, or candidate integrated. Status
  alone does not establish delivery. Old done records without candidates
  and unavailable evidence retain their uncertainty. The preview does not
  invent a requirement to merge between selected items; that execution
  choice belongs to [G-163](G-163-selected-work-review-boundary.md).
- Keep membership, preferred sequence, and actual prerequisites distinct.
  Equal graph depth means no declared ordering between those items, not
  permission or evidence for parallel writes. An absent dependency field
  means no declared prerequisites, not proven implementation readiness.
- Extend the [shaping workflow](../docs/work-shaping.md) at the relevant
  handoff: capture real prerequisite edges and explain their reasons in
  their owning work records; preserve preferred order as such. Do not
  manufacture an edge to encode priority or duplicate relationships in a
  separately editable graph file. Reconcile the owning command/board docs
  when behavior ships.

Reuse [G-042's current view](G-042-current-view.md), including visible
divergence, and the existing checkout selector. Conflicting branch versions
must not be silently combined into a supposedly authoritative DAG or order;
expose the ambiguity and bind a preview to a coherent source before acting.
Preserve exact source targeting, freshness, text escaping, batched Git reads,
and on-demand history. This work neither launches batches nor edits graph
edges by gesture. New merge-gate fields and an editable scheduling system
are outside this initial proposal.

[G-162](G-162-bounded-work-selection.md) owns managed execution of several
selected items. Graph-first is the assistant's suggested investment order,
not a technical prerequisite between the two proposals.

## Acceptance

1. From a representative ten-item backlog, the owner can identify a chain,
   a shared prerequisite, a convergence point, unrelated work, and the work
   a selected item unlocks without opening every record or asking an agent.
2. A selection preview preserves the explicit selection, orders its members
   correctly through transitive dependencies, and identifies outside
   prerequisites and relevant open questions without silently adding work.
   It explains facts and uncertainty rather than claiming authorization,
   safe parallelism, or a guaranteed completion time.
3. Review candidates not integrated into the intended base, historical done
   records without candidates, abandoned prerequisites, hidden dependencies,
   uncommitted records, divergent edges, and incomplete source reads each
   produce an honest, actionable explanation. A changed source invalidates
   any stale selection handoff.
4. Membership and priority cannot alter dependency order. A same-level pair
   is not labelled safe to run concurrently; a blocking question remains
   distinct from a work-to-work edge.
5. The owner judges 80- and 120-column layouts using the real 13-item
   component and a synthetic unfinished backlog with long titles, branches,
   convergence, and disconnected items. A larger case demonstrates bounded
   navigation rather than shrinking the whole project into illegibility.
6. Connected terminal checks cover board to graph, focus and selection,
   preview, record detail, return, refresh, resize, and exit. CLI and TUI
   agree on the same source's relationships and ordering. Loading retains
   the repository's Git-read and responsiveness constraints.
7. A fresh shaping-session handoff leaves real dependencies and their reasons
   discoverable in records, with preferred sequence identified separately;
   the resulting view is useful before multi-item execution exists.

## Next

The owner can assign this independently through `$grove-work G-161` after the
proposal is committed. Preparation should produce concrete terminal layouts
and a small shared relationship/preview design, with owner judgment of
readability. G-163 does not block this read-only outcome. No implementation
has been assigned by this shaping session.
