---
id: "G-188"
type: decision
title: "Review an explicitly selected set of work together, on one shared candidate"
status: accepted
created: "2026-09-25T23:22:15Z"
updated: "2026-09-25T23:22:27Z"
relates_to: ["G-162", "G-163", "G-185", "G-057", "G-059", "G-060", "G-101"]
---

## Decision

On 2026-09-25 the owner answered
[G-163](G-163-selected-work-review-boundary.md) with its option 1: when
the owner explicitly selects several work items, they are implemented in
order on one branch, and a later item may build on an earlier item's
unmerged changes. Human review and integration stop the selection once,
at its end, not between items. The owner approved the consequences in plan
[G-185](G-185-g-162-selected-work-plan.md) by launching its continuation:

- One Grove-owned attempt runs the whole selection sequentially under one
  aggregate budget; Grove never retries it.
- The complete members enter review together with one shared
  [candidate](G-057-candidate.md): the records on a branch whose
  `candidate` is the same commit are one group, and nothing else defines
  it.
- [Approval](G-059-approval.md) stays per member, so each keeps its own
  acceptance and verdict. [Integration](G-060-integration.md) is per group:
  merging the commit merges all of it, so it waits until every member is
  approved.
- Feedback on any member reopens the group, dropping every member's
  approval, since the next candidate replaces the shared one.
- A started member that is incomplete holds the whole branch out of
  review; members that never started keep their wait and do not.

A dependency edge still means only that the prerequisite is needed: it
does not itself authorize shared execution, which the owner's explicit
selection does. The meanings of work, attempt, candidate, approval and
integration do not change; a group of one behaves exactly as before.

## Alternatives

- **Review and merge between dependent items** (G-163 option 2): smaller
  reviews, but an unattended chain stops at its first boundary.
- **A boundary chosen per assignment** (option 3): flexible, but doubles
  the execution, recovery and verification contracts before either is
  proven.

## Reconsideration

Reconsider when combined reviews prove too large to judge, when owners
need intermediate judgment inside a selection, or when parallel
implementation is selected.
