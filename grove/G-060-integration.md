---
id: "G-060"
type: term
title: "Integration"
status: settled
created: "2026-09-21T05:01:57Z"
updated: "2026-09-21T14:23:34Z"
relates_to: ["G-054", "G-057", "G-059"]
formerly: "T-007"
---

## Meaning

An approved [candidate](G-057-candidate.md) becoming part of the project's
configured target, such as a merge into the main branch, shown by Git
ancestry rather than by a status.

Since G-038, done means exactly this for work with a `candidate`: `update`
writes done only where that commit is already an ancestor of HEAD. A done
record without a candidate is older and claims its outcome only in its own
branch; it may never have been merged, and it keeps that meaning.

Since G-044, `grove integrate` performs it from the target's checkout: a
plain merge of the one branch holding an approved candidate, aborted on
conflict, then done written and committed there, with the branch and its
worktree removed only on request and only where Git agrees. Since G-177,
a conflict is predicted with `git merge-tree` in objects only and refused
before the merge starts, naming the files; the abort remains for what the
prediction cannot see.

## Relationships

Follows [approval](G-059-approval.md). Completes implementation
[work](G-054-work.md) under the target lifecycle.
