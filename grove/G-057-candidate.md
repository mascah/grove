---
id: "G-057"
type: term
title: "Candidate"
status: settled
created: "2026-09-21T05:01:56Z"
updated: "2026-09-21T14:23:32Z"
relates_to: ["G-056", "G-058", "G-059", "G-060"]
formerly: "T-004"
---

## Meaning

The specific result of an [attempt](G-056-attempt.md) that is offered for
judgment: a Git commit on a work branch, together with the evidence gathered
at that commit. It is exact on purpose, so that what was examined, what was
approved, and what gets integrated can be shown to be the same thing.

A branch name is not a candidate, because it moves. A candidate that changes
is a new candidate.

## Relationships

A [review](G-058-review.md) examines a candidate. [Approval](G-059-approval.md)
is of a candidate. [Integration](G-060-integration.md) puts an approved
candidate into the target.
