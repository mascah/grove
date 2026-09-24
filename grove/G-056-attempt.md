---
id: "G-056"
type: term
title: "Attempt"
status: settled
created: "2026-09-21T05:01:56Z"
updated: "2026-09-21T14:23:32Z"
relates_to: ["G-054", "G-055", "G-057", "G-058"]
formerly: "T-003"
---

## Meaning

One particular execution of [work](G-054-work.md) by one owner: the inputs it
started from, its checkout, its progress, its result, and what is needed to
recover or resume it. Several attempts can serve one work item, for example
after review feedback.

An attempt ending, even successfully, is not the work being done and does not
by itself put anything into [review](G-058-review.md). What Grove keeps
about an attempt, and where, is documented under
[attempts](../docs/commands.md#attempts), not here.

## Relationships

Includes [preparation](G-055-preparation.md) and implementation. Produces a
[candidate](G-057-candidate.md).
