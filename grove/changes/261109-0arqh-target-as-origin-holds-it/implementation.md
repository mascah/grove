# Implementation: The target is what this checkout and origin both hold

<!-- Progress on change 261109-0arqh, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — The target's commit is the merge base with origin's

Status: done

- `gitx.ResolveBranch` returns the merge base of `refs/heads/<name>` and
  `refs/remotes/origin/<name>`, or the one that exists, and Git's error when
  neither does; it fetches nothing. Every reader of the target already calls
  it: `agreed` (project index, `{{commit}}` for `base` and `planned`,
  implement, archive), `movedSinceBase`, `serve`.
- `gitx.LocalBranch` is the read of the local tip (delegated name);
  `claim`'s fast-forward refusal and its `update-ref` of the local target
  use it, so they still compare the local tip with the fetched origin.
- `internal/cli/walkthrough_test.go` `TestOriginTarget`: a bare origin, two
  clones and a checkout with `origin/main` fetched and no local `main`.
  The first clone's `main` holds an unpushed commit; `grove instructions
  propose` there prints `base:` as origin's commit; after a squash merge
  onto origin's `main` and a pull, the second clone's `grove check` is
  clean; a change merged on the second clone's `main` is not listed by
  `grove changes` until pushed; the checkout with no local `main` runs
  `grove check` clean and `grove changes` as the second clone does.
  Against the old `ResolveBranch` the test fails on all three.
- `TestClaim`'s "the local main has commits origin's does not" refusal
  still passes.
- `just check`: passed; 4406 Go, 2100 test, 410 page lines, 8613
  instruction words.

Also changed: `internal/cli/claim.go` — `run.go`, which the task names, became `claim.go` in 261108-jpg8n; its fast-forward check reads the local tip through `gitx.LocalBranch`.
Also changed: `internal/cli/review.go` — review's diff span was `main...HEAD`, reading the local branch directly; it now uses the target's commit, as decision 1 says the rule governs review.

Archive's "where this branch left the target" (`merge-base HEAD refs/heads/<target>` in `archive.go`) is unchanged, as decision 8 says.

### T-02 — Documentation and the walkthrough

Status: not started

## Scenarios

<!-- Before finishing, one line per scenario of the change specs: the check
that shows it and its result, or `unverified` and why:

    - record-editing "Update the note on a legacy record": tests/test_cli.py test_edit_legacy, passed
-->

## Reconciliation

<!-- One line per change spec, once applied to its maintained spec:

    - record-editing: applied to specs/record-editing/spec.md
-->
