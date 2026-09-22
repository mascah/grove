---
id: "G-072"
type: review
title: "G-071 process count review"
status: current
created: "2026-09-22T03:46:56Z"
updated: "2026-09-22T03:53:38Z"
work: ["G-071"]
examined: "418e3d2942ff5741d10dc0d070c62f9cf350bfc2"
---

## Examined

Two rounds by an independent reviewer agent that wrote none of the code and
edited nothing, on `worktree-G-071` from main `5c637b4`, reading
[G-071](G-071-spawn-fewer-git-processes-per-in.md) and the diff.

- Round 1: the uncommitted working tree before `418e3d2`. The reviewer built
  binaries from the base and the tree and ran eleven repository states through
  both, diffing `versions` output.
- Round 2, final: `418e3d2`, this record's `examined`, against the base with
  twenty states and mutation runs of both guards. The commit after it,
  `9a69223`, holds only the final-check test row suggested below.

## Findings

Round 1, consequential: the one-process path never checked that the
registered worktree path is itself a directory. `git worktree list` keeps
reporting a path after it is replaced by a symlink, and Git answers
`rev-parse` through the link exactly as for the checkout, so a symlink at the
registered path, including one leading to a tree outside the repository with
a gitfile pointing at `<common>/worktrees/NAME`, was admitted as a valid live
source, contributed a version with a working selector, and `workspace
--source` resolved a record path through the link. The base refused it with
"worktree path is a symlink or not a directory". The same gap was in the
cheap check before the second inventory's re-entry and, through it, in the
final check of a resolved workspace.

Round 1, minor: the cheap check compares a `.git` directory by shape, not
contents, and does not look for a filesystem boundary appearing between the
root and the project; both need mid-read edits of Git metadata. Accepted as
ceilings if named.

Round 1, process: the "inspect entering a worktree" cancellation scenario
blocks on an argument that the root's first process now also carries, as
the base's second process already did; a fidelity note, not a regression.

Round 2: the fix holds at all three sites. Both symlink variants are refused
with the base's wording, no `live:feature:` selector is offered, and a
selector captured before the swap is refused identically on base and fix.
Twenty states give identical output, wording and order on base and fix. Each
guard is load-bearing: removing the one in `readWorktree` fails
`TestInspectWorktreeReplacedBySymlink`, removing the one in `unchanged` fails
the "checkout root symlinked" case of `TestInspectIdentityChangedDuringRead`.
The final-check path was correct but unpinned; a scratch table row in
`TestResolveFinalCheck` fails with the guard removed. Nothing consequential
remains.

The reviewer's own counts, from a logging `git` wrapper on `PATH` rather
than `GIT_TRACE`: `rev-parse` 2,006 to 642 over the package, all processes
3,811 to 2,291, and 15 to 8 for one `versions` on this repository. The
`rev-parse` row of G-071's table agrees within 0.3%; the all-processes row
differs by about 5% because the two instruments count differently.

## Disposition

- Consequential finding: fixed in `418e3d2` with `isDir` guards in
  `readWorktree` and `unchanged`, pinned by the two tests above.
- Minor finding: named as a ponytail ceiling on `unchanged`; no change.
- Final-check row: added in `9a69223`.
- Cancellation scenario note: left as is.
- Verification the reviewer ran on `418e3d2`: `gofmt -l`, `go vet`,
  `grove check` clean; `go test -count=1 -timeout 180s ./...` green;
  `internal/versions` alone 4.60, 4.73, 4.56 s.
