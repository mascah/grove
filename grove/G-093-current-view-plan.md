---
id: "G-093"
type: plan
title: "Project-wide current view: per-record merge-base projection"
status: current
created: "2026-09-22T21:39:28Z"
updated: "2026-09-22T21:39:56Z"
work: ["G-042"]
---

# G-042 current view plan

**Goal:** the board opens on one project-wide view of work, the same from any
linked checkout. Each card shows a record's current state, derived from Git
ancestry over every local branch tip and checkout. Older copies move into the
card's evidence, and genuine divergence stays visible.

**Spec:** [G-042](G-042-current-view.md). **Base:** `939d090` (main), branch
`worktree-G-042`.

## Owner decisions, 2026-09-22

Asked with examples in the `/grove-work G-042` session:

- **Uncommitted edits count, labelled.** A checkout's uncommitted record is a
  change on top of its HEAD, so it can be the current state, marked
  uncommitted in that checkout. (Selected the recommended option.)
- **Divergence is one card.** The owner selected per-column cards, then added:
  "I think it's important we don't make this too confusing for the user. We
  might just want to show a single card and show the divergence as an
  indicator with more detail on the detail page." Asked where that card sits,
  they chose the earliest lifecycle status among the competing states ("I think
  earliest"), and said the situations were hard to picture. The handoff must
  show concrete examples.
- **No target:** "No target for now." No `target` key and no branch name is
  special. The card detail says where the current state lives. The board
  cannot call anything unintegrated.

## The projection

For one record ID, every valid source gives one **observation**: the record's
exact bytes in that source, or absence. A branch tip observes at its commit.
A checkout whose file matches its HEAD is the same observation as that commit.
One whose file differs (modified, added, deleted) is an **uncommitted**
observation on top of HEAD.

Observation A is **older than** B when their contents differ and the record's
content at the merge base of their commits equals A's: since the two split,
only B's side changed the record. For two observations on one commit, the base
is that commit, so an uncommitted edit is newer than its own HEAD. Then:

| Case | Result |
| --- | --- |
| Same bytes | one state, any number of places |
| Stale branch that never changed the record (merged or not) | older; the other side is current |
| Unmerged branch that changed it, main unchanged since the split | branch's state is current |
| Record only on one branch | current there; absence elsewhere is older (the base lacked it too) |
| Committed deletion on a branch, others unchanged | current state is a deletion: no card, a deleted row in `versions` |
| Both sides changed since the split, different bytes | both current: divergence |
| Revert | a change like any other: judged against the base, never by earlier bytes |
| Record unreadable at the base, or several bases disagreeing | not ordered: both current, with a note naming the pair |
| Older forms a cycle (reverts carried across merges) | current when everything newer, through any chain, is also older: a cycle nothing outside supersedes is current, with a note |
| Invalid or unreadable source | contributes nothing; result incomplete, as today |

**Current** states are the contents held by at least one observation that
nothing is newer than. Timestamps, status order, and branch recency never
decide. Only commits matter, so a tip commit that touches other files changes
nothing. The inputs are the same set of sources from any linked checkout, so
the projection is too.

The record's "supported old/new schema sources and legacy/new IDs" predates
G-052. `schema_version: 3` is now the only schema, so a branch from before the
conversion is an invalid source: the last row of the table, fixture-tested as
one.

Merge bases come from commit objects read through the inspection's existing
`cat-file --batch` process (Git's paint-down-to-common, by committer date).
They are computed only for records whose observations differ, and cached per
commit pair. The base's record is found through the same cached `loadTree`.
There are no new Git processes, per branch or otherwise.

## Surfaces

- `versions.Version` gains `Older` (why it is superseded, naming the newer
  place and the base commit; empty when current). `Group` gains `Notes` for
  pairs that could not be ordered. A branch whose current state is a deletion
  gets a committed row with no record, like a checkout's deleted row.
- `grove versions`: a `CURRENT` column (`yes`/`older`); JSON `current` and
  `older`.
- Board: opens on the current view. `b` offers it first, then each checkout's
  own board as before (explicit source inspection). Cards are placed by their
  current state. For divergence, the card sits in the earliest status among
  the current states and is marked `⑂ N states`. A current state with
  uncommitted edits is marked `uncommitted`. Elsewhere lists work whose
  current state is a deletion. An open card lists current rows first, then
  older ones marked `older`. Its header describes the current state or states
  and where each lives. A row says why it is older. History follows the first
  current version under the header. Selecting a workspace is unchanged: one
  exact row, freshly resolved.

## Steps

1. Projection in `internal/versions/current.go`, with merge-base reading in
   `tree.go`'s objects. Fixture tests for every row of the table, including
   the G-030/G-023 shape (acceptance 4), committed deletion, detached and
   dirty checkouts, and invocation from two checkouts giving equal
   projections.
2. `versions` column and JSON; CLI tests.
3. Board current view; model tests; terminal lifecycle check.
4. Measure load on this repository and a synthetic many-branch repository
   (G-031's shapes) before and after.
5. Docs: README board and `versions`, G-002's note, AGENTS.md's
   checkout-scoped sentence, record model if touched. Then review and hand off.
