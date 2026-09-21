---
id: "T-003"
type: term
title: "Attempt"
status: proposed
created: "2026-09-21T05:01:56Z"
updated: "2026-09-21T05:02:47Z"
relates_to: ["T-001", "T-002", "T-004", "T-005"]
---

## Meaning

One particular execution of [work](T-001-work.md) by one owner: the inputs it
started from, its checkout, its progress, its result, and what is needed to
recover or resume it. Several attempts can serve one work item, for example
after review feedback.

An attempt ending, even successfully, is not the work being done and does not
by itself put anything into [review](T-005-review.md). Grove has no attempt
record yet; today an attempt is visible as a branch, a worktree, and the
checkpoint in the work's Next.

## Relationships

Includes [preparation](T-002-preparation.md) and implementation. Produces a
[candidate](T-004-candidate.md).
