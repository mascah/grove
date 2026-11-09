# Changes to config

The changes that change 261109-0arqh proposes to the requirements of
[config](../../../../specs/config/spec.md). This is future work.

The Purpose gains a sentence: "The `target` in `grove.yaml` is read as
this checkout and origin both hold it." The Boundaries' "The rest of
`grove.yaml` has no maintained spec yet" becomes "`target` is covered
below; the rest of `grove.yaml` has no maintained spec yet".

## ADDED Requirements

### Requirement: Read the target as origin holds it

When a command reads the branch `grove.yaml` names as `target`, the
system shall read it at the newest commit that both the local branch and
the checkout's `refs/remotes/origin/<target>` hold, or at whichever of
the two exists when only one does, fetching nothing.

#### Scenario: Proposed from a local main that is ahead

- GIVEN two clones of one origin, and in the first the local `main` holds
  a commit that origin's `main` does not
- WHEN its owner proposes a change on a branch there and the branch is
  squash-merged into origin's `main`
- THEN the change's `base` is the commit origin's `main` held, and in the
  second clone, after a pull, `grove check` reports nothing for the change

#### Scenario: A checkout with no local target branch

- GIVEN a checkout whose `origin/main` was fetched and which has no local
  `main`, as a pull-request job's checkout is
- WHEN `grove check` or `grove changes` runs
- THEN it reads the target at `origin/main` and reports as a clone with a
  local `main` at that commit would

#### Scenario: A local merge not yet pushed

- GIVEN the local `main` holds a change's folder that `origin/main` does
  not
- WHEN `grove changes` runs
- THEN the change is not listed, and is listed once `main` is pushed

#### Scenario: No remote

- GIVEN a repository with no `origin`
- WHEN any command reads the target
- THEN it reads the local branch, as before this change

## Observed current behavior

- `ResolveBranch` in `internal/gitx/gitx.go` at ccaee4d resolves
  `refs/heads/<name>` only; `agreed` in `internal/cli/project.go`,
  `movedSinceBase` in `internal/cli/check.go`, `internal/cli/serve.go`
  and `runAttempt` in `internal/cli/run.go` call it.
- On the ledger fixture on 2026-11-08, not committed: a `base` recorded
  from a local `main` one unpushed commit ahead, then squash-merged, was
  reported by `grove check` as not a commit in a second clone of the same
  origin; a checkout with `origin/main` and no local `main` failed
  `grove check`, `grove changes` and `grove instructions propose` with
  `grove: git rev-parse: fatal: Needed a single revision: exit status 128`.
