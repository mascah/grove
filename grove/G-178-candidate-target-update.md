---
id: "G-178"
type: work
title: "Update a conflicting candidate to the moved target through one bounded attempt"
status: done
created: "2026-09-25T21:39:27Z"
updated: "2026-09-26T03:08:42Z"
kind: feature
relates_to: ["G-044", "G-045", "G-046", "G-057", "G-058", "G-059", "G-060", "G-101", "G-134", "G-140", "G-177", "G-179", "G-180", "G-182", "G-192"]
candidate: "4e9b1b26a7b699659ec8d697a67fc8e21121aed3"
approved: "4e9b1b26a7b699659ec8d697a67fc8e21121aed3"
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

## Evidence

Implemented headless on `worktree-G-178` from main `e812672` (G-177
integrated). It started from this record at `sha256:83ec768e…` and plan
[G-191](G-191-plan-for-g-178-candidate-resolution.md), committed at
`38c765c`. Code runs through `7cb7cea`; the candidate adds only this
evidence and the review record.

**The operation.** `grove resolve ID` from any checkout, or the board's `m`,
runs one function for both, `attempt.Resolve` (`internal/attempt/resolve.go`),
so [G-180](G-180-policy-driven-integration.md) calls the same mandate and
refusals.

- It finds the one branch holding the record in review and that branch's
  checkout. It predicts the merge with G-177's `PredictContext` and refuses
  anything but a conflict.
- It records feedback there through `update.Feedback`. The text comes from
  `Mandate`: the target commit in full, the conflicting files, and "merge
  that commit, never a rebase or a later tip; resolve those files; verify;
  hand off the merge; change nothing else; stop at an unsettled choice".
- It then calls `Start` once, on that branch in that checkout, with the
  launch defaults. The group sharing the candidate
  ([G-188](G-188-selected-work-shared-candidate.md)) reopens and runs
  together, the given ID first.
- The assignment stays `/grove-work IDS --interaction headless`. The mandate
  travels in the record, where the work guide reads feedback.
- The name is `resolve`, following G-182's "resolution attempt". The term
  [G-192](G-192-resolution.md) (proposed) defines it apart from a question
  being resolved.

**Acceptance 1.** One action starts exactly one attempt, with the mandate
visible in the record and in the board's line. Refused, with nothing
written:

- no target;
- no conflict (the fact's own text is given);
- not in review, or in review on several branches;
- no checkout of the branch;
- a running or orphaned attempt of any member;
- no budget or mode;
- every member waiting once active;
- a missing or incompatible skill or reviewer;
- a missing or failing provider;
- from the board, a candidate or target commit that changed since the fact
  was shown.

`TestResolveRefusals` and `TestResolveAGroup` check that nothing was
written and no attempt started. A second `resolve` finds the work active
and is refused. A launch that still fails after the feedback says the
feedback stands and gives the `grove run … --branch … --worktree …` that
launches it. `integrate`'s conflict refusal now names `grove resolve ID`.

**Acceptance 2.** The work guide's step 5 gains "A target that moved":

- merge the named commit, never a rebase, so the earlier candidate and every
  review's `examined` stay ancestors;
- resolve only those files;
- verify, and hand off with Evidence naming each conflict and how it was
  settled;
- a choice the record does not settle is a missing human decision;
- scope the review with `git show --remerge-diff MERGE`.

`TestResolveCleanly`: the new candidate's ancestors include the previous
candidate and the merged target commit. `TestResolveNeedsAChoice`: a clean
provider exit with nothing committed leaves the work active, with its
candidate.

**Acceptance 3.** `versions.Changes.Resolution` is read on demand with the
other changes, never in the board load:

- the latest first-parent merge whose second parent the target holds;
- the target commit it merged;
- the candidate the record named before it, read at the merge's first
  parent;
- each file that merging the parents again conflicts on, with the side it
  kept when it kept one.

The Review block prints this as a row, for example "Resolution: merge M of
main at T into candidate P, whose reviews stay comparable · resolved: a.go,
b.go (took main's side, dropping the branch's change)", and the Changes list
marks those files. `TestResolution` covers a combined resolution, a side
taken, a file Git merged by itself (not listed), a project under a prefix,
a rename/rename, and a merge of another branch or of unrelated history
(none).

**Acceptance 4.** Fake-provider tests in `resolve_test.go`:

- `TestResolveCleanly`: a clean resolution;
- `TestResolveNeedsAChoice`: a resolution that needs a choice;
- `TestResolveWhileTheTargetMoves`: main moves during the attempt; the
  attempt merges the commit its feedback names, and the prediction then
  names the moved main;
- `TestResolveStopped`: Stop during the attempt leaves the work active with
  its candidate.

An end-to-end `grove resolve` with the built binary and a fake provider, on
a disposable repository, printed the feedback, the reused worktree and the
attempt, and a second `resolve` was refused. No real-provider trial was run;
the record bounds it separately.

**Acceptance 5.** These now describe the contract:

- `docs/work-execution.md` (the procedure, and `resolve` among the judging
  dispositions);
- `docs/commands.md` ("Resolving a conflict");
- `docs/board.md` (`m`, the resolution row, the key table);
- `docs/record-model.md` (the lifecycle);
- `grove --help` and `README.md`.

**A fix to G-177 found on the way.** `git merge-tree --name-only` names files
relative to the current directory, so for a project under a prefix G-177's
conflicts were not from the repository's top, as `Merge` documents.
`resolveCommits` now reads the prefix in its one `rev-parse`, and `predict`
joins it, for every caller.

**Verification at `7cb7cea`:**

- `gofmt -l .`: empty.
- `go vet ./...`: clean.
- `go run ./cmd/grove check`: OK, 186 records.
- `go test -count=1 -timeout 120s ./...`: all ok.
- `python3 internal/tui/testdata/terminal.py BINARY`: all 12 ok.

Package times: the machine was loaded (load average 6 to 15) during the
final runs. Under the same load, `internal/versions -short` took 13.9s on
main and 15.1s here. When the machine was idle earlier, it took 4.3s against
4.4s and `internal/attempt -short` about 2.3s. No new test builds, sleeps or
waits on a shim without a `-short` skip.

**Review:** [G-193](G-193-review-of-g-178-candidate-resolution.md), three
rounds by fresh `grove-reviewer` agents. Two consequential findings, the
side-taking gap and the prefix paths, were fixed and re-reviewed. Round 3
found nothing consequential.

**Limits for the owner:**

- A user's `diff.relative=true` leaves resolved files unlabelled under a
  prefix, never mislabelled. The fix is `--no-relative`, which the older
  `numstat` and `Others` diffs lack too.
- No suite test covers a branch whose history began apart and then merged
  the target; a reviewer verified that case.
- A prediction failure on the merge's parents, such as Git before 2.38,
  drops the resolution row without a reason.
- Only the latest merge is read, so a later merge of another branch hides
  the resolution.
- `Resolve` takes no policy attribution, which G-182 requires of G-180's
  automatic acts; G-180 adds it.
- `grove resolve` on the CLI checks the fact it computes itself. Only the
  board passes a shown fact.

## Next

In review. The owner judges the candidate on `worktree-G-178`:

```sh
grove approve G-178 "VERDICT"   # in /Users/mascah/GitHub/mascah/grove/.claude/worktrees/worktree-G-178
grove integrate G-178           # then in the main checkout
```

Or give feedback with `grove feedback G-178 "TEXT"` in the worktree. To
demo: on a candidate in review that conflicts with main, open its card on
the board, press `m`, then Enter; or run `grove resolve ID`. G-180 builds
on `attempt.Resolve`.

Verdict on candidate 4e9b1b2, 2026-09-26: approved
