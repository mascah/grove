---
id: "W-013"
type: work
title: "Keep a full board load fast with hundreds of branches and worktrees"
status: active
created: "2026-09-20T04:36:57Z"
updated: "2026-09-20T04:57:47Z"
kind: fix
priority: 2
size: small
relates_to: ["W-004", "W-006", "W-008", "W-009", "W-012"]
---

## Outcome

A full load (every local branch tip and every checkout, which the board does
at start, on refresh, and before returning a workspace) stays well under a
second in a large repository. The owner raised this on 2026-09-19: 0.6 s for
this small repository suggested that an enterprise repository with a few
hundred branches would choke. Measured the same day, it did, and the cause
was fixed on branch `worktree-W-009`.

## Evidence, 2026-09-19

Synthetic repositories built with `git fast-import` from this project's own
`grove.yaml` and records plus 20,000 unrelated Markdown files, five linked
worktrees, and one commit per branch. Times are `grove versions` from a built
binary on the owner's Mac, second run; processes counted with `GIT_TRACE`.

| Repository | Before | After |
| --- | --- | --- |
| This repository: 4 branches, 4 checkouts | 0.57 s, 44 processes | 0.24 s, 18 processes |
| 300 branches differing only outside `grove/`, 6 checkouts | 27.7 s, 656 processes | 0.35 s, 24 processes |
| 1,000 branches, each with its own edit to W-001, 6 checkouts | 91.0 s | 0.45 s |

`versions` output before and after is byte-identical on the 300-branch
repository (5,529 lines) and on this one.

What was slow, and what changed:

- Each distinct branch commit ran `git ls-tree -r` over the whole tree below
  the project and then fetched every `.md` blob in it, in two processes. Cost
  grew with branches times repository size. Now one `git cat-file --batch`
  process serves the whole inspection, asked only for object IDs. It reads what
  the project loader reads and nothing else: `grove.yaml`, the entries on the
  path to the record folder, and that folder. Trees and blobs are cached by ID,
  and branches that agree on all of that share one loaded, validated project.
  A branch now costs a few object lookups and no process.
- Each checkout ran eight single-path `rev-parse` processes (twelve with a
  nested project), because a path may contain a newline (W-008). They are now
  asked together, three per checkout (five nested). A combined reply is used
  only when it holds exactly one terminator per option, which is only possible
  when no path holds a newline; otherwise each is asked alone as before.

One behaviour is wider than before: an entry with an unusable name outside the
record folder no longer fails its branch, because it is no longer read.

## Constraints

Reads only. Keep W-006 and W-008 provenance checks, the second worktree
inventory, cancellation of every Git process, and identical `versions` and
`workspace` results. No cache between runs, no index, no daemon.

## Acceptance

1. Identical `versions` output before and after on real and synthetic
   repositories. Met, above.
2. A test shows branches differing only outside the record folder share one
   loaded project through one `cat-file` process with no whole-tree listing,
   and a differing record folder gets its own
   (`TestCommittedReadIsScopedAndShared`). Batched path answers equal single
   answers for paths with newlines, tabs, and a separate Git directory
   (`TestGitPathsMatchSingleAnswers`).
3. Existing cancellation, provenance, symlink, submodule, nested-prefix, and
   unstable-read tests pass unchanged; full and race suites and vet pass.

## Remaining limits

- Checkouts are read one after another at about 35 ms each, so dozens of
  worktrees would again reach a second. Reading them concurrently is the next
  step if that happens; it is marked `ponytail:` in `internal/versions/live.go`.
- A thousand branches produce a thousand rows per record in `versions` output,
  and a fold of a thousand places in the board. Whether a load should read
  every local branch, or only those that touch the record folder or are checked
  out, is an owner decision about meaning, not speed.
- The first run after a build or on a cold disk is about twice the steady time.
  Times above were taken on an otherwise idle machine. With a render using
  every core, the same binary took 0.85 s on the 300-branch repository and
  0.6 s here; process counts do not change with load and are the steadier
  measure.
- Measured on one machine, macOS, SHA-1 repositories. The reader sizes object
  IDs from the commit ID, so SHA-256 repositories should work; not run.

## Next

Owner reviews the numbers. If accepted, mark done with
`go run ./cmd/grove update` when `worktree-W-009` is integrated.
