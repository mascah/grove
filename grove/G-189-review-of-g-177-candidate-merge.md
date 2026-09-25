---
id: "G-189"
type: review
title: "Review of G-177 candidate: merge prediction"
status: current
created: "2026-09-25T23:24:09Z"
updated: "2026-09-25T23:24:25Z"
work: ["G-177"]
examined: "829f08c63aaa8506dae3c838b5484b1abe2f108d"
---

## Examined

The G-177 implementation on `worktree-G-177`, base `370df4a` (plan
[G-183](G-183-plan-for-g-177-merge-prediction.md) and status active, on
main `6b14141`), by three fresh `grove-reviewer` agents, one per round:
`0fbf4c3`, then fixes to `502b27f`, then to `829f08c`. The record's
acceptance and constraints, the plan and `CLAUDE.md` were the criteria.
The reviewers ran `go vet`, `gofmt -l`, `grove check`, `go test -short` per
package, and the built binary in disposable repositories, including a PATH
shim that makes `git merge-tree --write-tree` fail as on Git before 2.38.

## Findings

Round 1 (`0fbf4c3`):

1. Medium. On Git before 2.38, a failed prediction refused every
   `integrate` and failed the Review block's whole changes read.
2. Low-medium. A conflict Git names no file for (a split directory rename)
   printed an empty file list everywhere.
3. Low. No test reached `integrate`'s `git merge` refusal path any more.
4. Low. The board and `deps` predict the candidate; `integrate` predicts
   the branch tip it merges.
5. Low, interpretive. A prediction is `rev-parse`, `merge-base` and
   `merge-tree` (plus `commit-tree` per clean step of an order), not one
   process per candidate as acceptance 4 words it.
6. Low, knowledge. The settled term G-060's mechanics paragraph still said
   integrate aborts on conflict.

Round 2 (`502b27f`): all six verified resolved or dispositioned; new:

7. Low-medium. The refusal's `grove feedback` command was `%q`-quoted,
   which the CLI's escaping turned into unpasteable `\"`.
8. Low. `record-model.md` and `work-execution.md` overstated the refusal
   for the unpredicted fallback; `board.md` lacked the unpredicted text.
9. Cosmetic. "cannot" versus "could not be predicted".

Round 3 (`829f08c`): 7 to 9 verified resolved, including pasting the
printed feedback command into a disposable repository (exit 0); nothing
consequential. Acceptance 1 to 5 met.

## Disposition

1. Fixed in `502b27f`: `integrate` falls through to `git merge`;
   `Changes.Unpredicted` keeps the files and the block says it could not be
   predicted (`TestReviewMergeUnpredicted`). No automated shim test for
   integrate's fallthrough (a PATH shim is not parallel-safe); round 2
   verified it by running.
2. Fixed in `502b27f`: `Merge.Where` says "where Git names no file", with a
   real split-rename test.
3. Fixed in `502b27f`: subtest "what prediction cannot see".
4. Kept, a stated limit: `integrate` already refuses commits after the
   candidate that change anything but the record, so the two can differ
   only when a record-only commit after the candidate conflicts on the
   target.
5. Left to the owner: reads stay bounded, per candidate, on demand, never
   during the board load.
6. Fixed in `502b27f`: G-060 names the pre-merge refusal.
7. to 9. Fixed in `829f08c`.
