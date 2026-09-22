---
id: "G-071"
type: work
title: "Spawn fewer Git processes per inspection"
status: done
created: "2026-09-22T03:19:16Z"
updated: "2026-09-22T03:59:15Z"
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

## Evidence

Implemented on `worktree-G-071` from main `5c637b4`, code in `418e3d2`.
Design as proposed, with one addition found while measuring: skipping the
second inventory's re-entry was needed to reach half, since one process
per checkout at each read still left about 1,200 `rev-parse`.

- The root's `--git-dir` rides on the one `rev-parse` that locates the
  repository (`repo.IdentifyContext`).
- `readWorktree` enters and owns a checkout with one `rev-parse` in the
  record folder: Git's prefix there, with a matching common directory and
  a main- or linked-shaped Git directory, proves every directory above it.
  Any mismatch, prunable entry, symlink, missing or invalid project takes
  the old step-by-step path, so diagnostics keep their wording and order.
- `Source.unchanged` skips the re-entry after the second inventory when
  the `.git` entry at the root, its `commondir`, the absence of `.git` or
  `HEAD` in the directories down to the project, the walk to the project
  and the `grove.yaml` bytes all read as before; otherwise the old
  re-entry runs and reports as before.

Acceptance 1, `GIT_TRACE` over `go test -short -count=1 ./internal/versions`
on this machine, 2026-09-21:

| | before (`fb2b398`) | after (`418e3d2`) |
| --- | --- | --- |
| `rev-parse --path-format=absolute …` | 2,001 | 640 |
| `worktree list` | 338 | 338 |
| `cat-file` | 186 | 186 |
| `for-each-ref` | 150 | 150 |
| all Git processes | 3,607 | 2,209 |

One `versions G-071` on this repository (three worktrees, project at the
root): 15 processes to 8. After: one `rev-parse` for the root, one per
checkout, two `worktree list`, one `for-each-ref`, one `cat-file`.

Acceptance 2: `go test -count=1 ./internal/versions` three times: 4.63 s,
4.55 s, 4.59 s (8.1 s before). `go test -count=1 -timeout 120s ./...`
passed in one run; no test or fixture was changed for speed.

Acceptance 3: every test in `internal/versions`, `internal/cli` and
`internal/tui` passes; the symlink, foreign-repository, moved-worktree and
change-between-reads cases keep their diagnostics. One test was changed in
scenario, not for speed: `TestCancellationKillsGitAndWritesNothing`'s
"re-entering a worktree after the second inventory" blocked a Git process
that an unchanged checkout no longer starts, so it now plants a `HEAD`
file in the project directory between the reads. New tests:
`TestInspectIdentityChangedDuringRead` (a repository begun in the project
directory, the checkout root symlinked, the checkout replaced by a foreign
one, and `commondir` repointed while the branch and commit stay the same,
each between the reads) and `TestInspectWorktreeReplacedBySymlink`.

Acceptance 4: the notes on `identity` in `live.go` and `Resolve` in
`workspace.go` state the new counts.

Independent review: [G-072](G-072-g-071-process-count-review.md). Round 1
found that the one-process path had dropped the check that the registered
path is itself a directory, so a symlink at it was admitted; fixed in
`418e3d2` with the two symlink tests above. Round 2 on `418e3d2` found
nothing consequential; `9a69223` adds the same case to the final-check
table. Process counts here are from `GIT_TRACE`; the reviewer's wrapper
counts all processes about 5% higher and agrees on `rev-parse`.

## Next

Done. The owner accepted the acceptance-3 reading on 2026-09-21 (the
cancellation scenario had to change once an unchanged checkout is not
re-entered) and asked for the merge; `worktree-G-071` was fast-forwarded
onto main.
