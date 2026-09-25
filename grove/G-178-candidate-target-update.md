---
id: "G-178"
type: work
title: "Update a conflicting candidate to the moved target through one bounded attempt"
status: proposed
created: "2026-09-25T21:39:27Z"
updated: "2026-09-25T21:41:22Z"
kind: feature
relates_to: ["G-044", "G-045", "G-046", "G-057", "G-058", "G-059", "G-060", "G-101", "G-134", "G-140", "G-177", "G-179", "G-180"]
---

## Outcome

When a candidate in review conflicts with the target, the owner can, with
one action, have a bounded attempt merge the target into the candidate's
branch, resolve the conflicts, rerun verification and hand off a new
candidate, and then judge only the resolution, instead of writing feedback,
relaunching, and re-reviewing the whole change.

Owner intent, conversation 2026-09-25: the feedback and launch loop for a
merge conflict is wasted time; the owner wants an agent to attempt the
resolution and complete, or abort when it finds a blocker, and leans toward
automating it. This record is the human-triggered form;
[G-180](G-180-policy-driven-integration.md) owns triggering it by policy,
which [G-179](G-179-standing-policy-question.md) must answer first.

## Constraints

### Observed evidence

At main `47852e3` the loop exists in pieces. [`feedback`](../internal/update/review.go)
returns the record to active, keeps `candidate` and unsets `approved`; `R`
on the board and `grove run` relaunch on the candidate's branch
([G-046](G-046-managed-runs.md)); the attempt reads the record's Next and
the feedback paragraph as its checkpoint; and `run` refuses a record still
in review ("judge that candidate (approve, feedback) before another
attempt", `internal/attempt/attempt.go`).

The [work guide](../docs/work-execution.md)'s judging section says nothing
about a target that moved since the branch was cut; its only nearby text
treats a squash or rebase as a manual merge followed by `update … --set
status=done --set candidate=COMMIT`. Approval is of one commit
([G-057](G-057-candidate.md), [G-059](G-059-approval.md)): a commit after
the candidate that changes any file but the record makes the tip a new
candidate, and `approve` and `integrate` refuse it (verified on a
disposable repository). A merge of the target into the branch is such a
commit, so the approval is withdrawn by design and the resolution needs its
own judgment. The Review block lists the candidate's changed files against
the target with per-file diffs ([G-044](G-044-review-integration.md)) and
has no notion of a previous candidate or of which content is a resolution.

### Proposed design and scope

- One operation, as a command and a board key, on a record in review whose
  candidate conflicts with the target (from
  [G-177](G-177-merge-prediction.md)'s fact or an `integrate` refusal). It
  records the feedback itself, with generated text naming the target commit
  and the conflicting files, and starts one attempt on the candidate's
  branch with the launch defaults ([G-140](G-140-default-an-attempt-s-budget-mode.md)) and
  the narrow mandate: merge that target commit, resolve those files, rerun
  the repository's verification, hand off a new candidate, change nothing
  else. Choose the name in preparation, avoiding `reconcile` (the
  predecessor's name for another operation) and `refresh` (the board's
  reload).
- Work guide: a procedure for a resumed attempt whose target moved. Merge
  the target rather than rebase, so the earlier candidate and every
  review's `examined` stay ancestors of the new tip; resolve; verify;
  commit; hand off with the previous candidate, the merged target commit
  and the resolved files named in Evidence. When a resolution needs a choice
  the record does not settle, stop with a checkpoint that names it: a
  headless attempt never enters review by itself.
- Judging the resolution: the handoff and the Review block name the previous
  candidate and the target commit merged and mark which changed files carry
  a resolution rather than the target's or the candidate's own content, so
  the owner and the independent reviewer judge the resolution. Approval
  stays the owner's act; re-review is scoped, not skipped.
- Bounds: one attempt per operation, no loop. If the target moves again,
  G-177 shows it and the owner decides. Duplicate-start protection, Stop and
  owner loss as [G-045](G-045-durable-attempt.md).
- Out of scope: automatic triggering, delegated approval and automatic
  integration (G-180); any change to the meaning of candidate, approval or
  integration ([G-060](G-060-integration.md)).

## Acceptance

1. From a review record whose candidate conflicts, one action starts exactly
   one attempt on the candidate's branch with the generated feedback and a
   visible narrow mandate. A record without a conflict, without a target,
   with an attempt running, or whose candidate changed since the fact was
   computed is refused with the reason.
2. The attempt ends in review with a new candidate whose ancestors include
   the previous candidate and the merged target commit, verification rerun
   and recorded, and Evidence naming the resolved files; or stays active
   with a checkpoint naming the choice it could not make. A provider exit
   alone marks nothing.
3. The owner can judge the resolution without re-reading the whole change:
   the handoff and the Review block name the previous candidate, the target
   commit and the files whose content is a resolution, and earlier reviews
   remain comparable to the earlier candidate.
4. Exercised with a fake provider on a clean resolution, a resolution that
   needs a choice, a target that moves during the attempt, and Stop during
   the attempt; a real-provider trial on a disposable project is separately
   bounded.
5. The work guide, the command documentation and the board documentation
   own the changed contract.

## Next

Assignable after G-177 or alongside it: `integrate`'s refusal already
supplies the conflicting files, so the operation works before prediction
exists. If the automatic form is what the owner wants, answer G-179 first
so this record's mandate and G-180's policy are shaped together. No
implementation has been assigned by this shaping session.
