---
id: "G-071"
type: work
title: "Spawn fewer Git processes per inspection"
status: proposed
created: "2026-09-22T03:19:16Z"
updated: "2026-09-22T03:19:42Z"
kind: refactor
size: medium
relates_to: ["G-030", "G-031", "G-052"]
---

## Outcome

`versions`, `workspace`, `context` and the board answer with fewer Git
processes per inspection, so an inspection of a repository with a few
worktrees is cheaper and `internal/versions` runs inside the owner's
five-second package budget. Behaviour observable through the CLI, including
every diagnostic the G-030 and G-031 tests pin, stays the same. This is the
owner's intent from the 2026-09-21 test-suite triage: the suite must stay
near ten seconds as features are added, and the remaining cost is the
product's, not the tests'.

## Scope and constraints

Observed on main `fb2b398`, 2026-09-21, with `GIT_TRACE` over
`go test -short -count=1 ./internal/versions`: 3,480 Git processes per run,
of which 2,001 are `rev-parse --path-format=absolute`, 338 `worktree list`,
186 `cat-file`, 150 `for-each-ref`, and about 760 fixture commands. With
150 inspections or resolutions that is about 13 `rev-parse` each for two
worktrees. `internal/versions/live.go` documents the per-checkout cost as a
deliberate ceiling: entering (`identity`), the prefix and record folder
(`owns` through `locate`), and the second inventory's re-entry, one after
another; `inspect` in `versions.go` also asks the root for `--git-dir` after
`LocateContext` already ran one `rev-parse`. Each Git process costs about
9 ms CPU under load and this machine sustains roughly 600 per second, so
package wall time is process count; `-parallel` and a shared fixture
template were measured and change well under a second.

In: cut the processes one inspection or resolution spawns. Proposed design,
labelled proposed: answer `--git-dir`, `--show-prefix` and
`--git-common-dir` for a checkout and its project directory from one
`rev-parse` where Git allows, reuse the root's identity from `LocateContext`
instead of a second call, and keep the second inventory's re-entry only where
it detects a real change between reads. Keep the `owns` guarantee: a symlink
or foreign repository at the project location must still be refused with the
same diagnostic. Update the ponytail notes in `live.go` and `workspace.go`
that state the old costs.

Out: reducing the fixture processes in tests, caching across invocations,
changing what `versions` reports or how sources are ordered, and any
concurrency across checkouts (the note in `live.go` reserves that for dozens
of worktrees; it does not lower CPU).

## Acceptance

1. `GIT_TRACE` over the same package run shows the `rev-parse` count at half
   or less of 2,001, with the before and after counts and the per-inspection
   process list recorded as evidence.
2. `go test -count=1 ./internal/versions` runs under five seconds on the
   owner's machine, measured three times, and the whole suite stays under the
   evidence run's timeout with no test or fixture changed for speed.
3. Every existing test in `internal/versions`, `internal/cli` and
   `internal/tui` passes unchanged; the symlink, foreign-repository,
   moved-worktree and change-between-reads cases keep their diagnostics.
4. The cost notes in `live.go` and `workspace.go` state the new
   per-inspection process count.

## Next

Assign. Needs no plan: the change is inside `internal/versions` and
`internal/repo`, and the evidence is the traced process count.
