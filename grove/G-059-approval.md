---
id: "G-059"
type: term
title: "Approval"
status: settled
created: "2026-09-21T05:01:56Z"
updated: "2026-09-21T14:23:33Z"
relates_to: ["G-057", "G-058", "G-060"]
formerly: "T-006"
---

## Meaning

A person with authority, or someone they delegated to, accepting one specific
[candidate](G-057-candidate.md) as meeting the work's outcome. It is bound to
that candidate: if the candidate changes, the approval does not carry over and
must be reconsidered.

Not a [review](G-058-review.md), which informs it; not passing checks, which
cannot supply judgment; and not [integration](G-060-integration.md), which
follows it. Since G-044, `grove approve` records it: the work record's
`approved` field names the candidate, the verdict is appended to the body,
and `check` rejects an approval of any other commit, so a changed candidate
needs its own. `grove feedback` withdraws it and returns the work to active.

## Relationships

Follows review, precedes integration.
