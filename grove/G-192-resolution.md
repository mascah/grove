---
id: "G-192"
type: term
title: "Resolution"
status: proposed
created: "2026-09-26T01:39:11Z"
updated: "2026-09-26T01:39:21Z"
relates_to: ["G-056", "G-057", "G-058", "G-059", "G-060", "G-178", "G-182"]
---

## Meaning

The merge of a target commit into the branch of a
[candidate](G-057-candidate.md) that conflicts with the target, with its
conflicts settled, handed off as a new candidate. It is a merge, never a
rebase, so the earlier candidate and everything reviewed at it stay
ancestors of the new one. What is judged is the resolution: the merged
target commit, the files the merge conflicted on and how each was settled.
Taking one side of a file drops the other side's change, and the
resolution names that.

A **resolution attempt** is one [attempt](G-056-attempt.md) whose whole
mandate is that merge, given as feedback that names the target commit and
the files. `grove resolve` starts one by hand
([G-178](G-178-candidate-target-update.md)); a standing owner policy may
start one ([G-182](G-182-standing-policy-delegation.md)). This sense is
unrelated to a question being resolved, which means it was answered.

## Relationships

It produces a candidate, which [review](G-058-review.md) examines and
[approval](G-059-approval.md) accepts as for any other. It changes nothing
about [integration](G-060-integration.md), which merges the new candidate
as usual.
