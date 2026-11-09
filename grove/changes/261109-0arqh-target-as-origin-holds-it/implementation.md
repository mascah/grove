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

Status: done

- `docs/concepts.md`: "Approval is a merge" says Grove reads main at the
  newest commit it shares with origin's as last fetched, fetching nothing;
  that a local merge is approved once pushed, and until then main is read
  behind the local branch; that a checkout with no local `main` reads
  origin's and one with no `origin` the local branch. The glossary's
  Target says how it is read.
- `docs/coordinating.md`: `base` is the commit main and origin's both
  held when the change was proposed; revising sets `base` to that commit.
- `docs/planning.md`: `planned` is that commit, in the template note and
  both re-plan paragraphs.
- `docs/reference.md`: the `target` row says how the target is read, and
  that with an `origin` a merge approves once pushed.
- `docs/guides/run-grove-in-ci.md`: a pull-request job's checkout, with
  `origin/main` and no local `main`, needs no extra step.
- `dev/walkthrough.md`: "The target as origin holds it" under Other paths.
  Run 2026-11-09 on a fresh `just fixture` with this build, the proposals
  written by hand from `grove instructions propose`: with an unpushed
  commit on `main`, propose printed `base:` as `origin/main`'s commit
  (63bd59e, main at 3d19717); after the squash merge and a pull, the
  second clone's `grove check` said no problems and `grove changes` listed
  budgets; a delete change merged there and not pushed was not listed,
  and was once pushed; a checkout with `origin/main` fetched, no local
  `main` and a detached HEAD gave `grove check` no problems and the same
  `grove changes` as the second clone.
- `just check`: passed; 4406 Go, 2100 test, 410 page lines, 8613
  instruction words.

## Scenarios

- config "Proposed from a local main that is ahead": `internal/cli/walkthrough_test.go` `TestOriginTarget` (propose's `base` is origin's commit; the second clone's `check` clean after the squash merge and pull), passed; walkthrough "The target as origin holds it" on the fixture, held.
- config "A checkout with no local target branch": `internal/cli/walkthrough_test.go` `TestOriginTarget` (`check` clean, `changes` as the second clone prints it, with no local `main`), passed; walkthrough step, held.
- config "A local merge not yet pushed": `internal/cli/walkthrough_test.go` `TestOriginTarget` (the change unlisted until `main` is pushed), passed; walkthrough step, held.
- config "No remote": the other tests in `internal/cli/walkthrough_test.go` that run with no `origin`, such as `TestPropose`, `TestCoordinating` and `TestMovedFiles`, passed unchanged.

## Reconciliation

<!-- One line per change spec, once applied to its maintained spec:

    - record-editing: applied to specs/record-editing/spec.md
-->
