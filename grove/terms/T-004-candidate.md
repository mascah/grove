---
id: "T-004"
type: term
title: "Candidate"
status: proposed
created: "2026-09-21T05:01:56Z"
updated: "2026-09-21T05:02:47Z"
relates_to: ["T-003", "T-005", "T-006", "T-007"]
---

## Meaning

The specific result of an [attempt](T-003-attempt.md) that is offered for
judgment: a Git commit on a work branch, together with the evidence gathered
at that commit. It is exact on purpose, so that what was examined, what was
approved, and what gets integrated can be shown to be the same thing.

A branch name is not a candidate, because it moves. A candidate that changes
is a new candidate.

## Relationships

A [review](T-005-review.md) examines a candidate. [Approval](T-006-approval.md)
is of a candidate. [Integration](T-007-integration.md) puts an approved
candidate into the target.
