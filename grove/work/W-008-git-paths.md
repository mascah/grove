---
id: "W-008"
type: work
title: "Preserve Git paths through discovery and coordination"
status: done
created: "2026-09-19T20:13:57Z"
updated: "2026-09-19T22:07:50Z"
kind: fix
priority: 2
size: medium
relates_to: ["W-002", "W-003", "W-004", "W-005"]
---

## Outcome

Git paths containing control characters retain their exact identity through
inspection, workspace resolution, and mutation coordination. All mutation locks
and counters stay under the actual Git common directory. Repairs existing
path/side-effect contracts (implemented; see Evidence). [Review R6](../../docs/reviews/2026-09-19-integrated-cli.md)
reproduced incorrect provenance and a stray write lock for a main checkout whose
name contains a newline. [Implementation plan](../../docs/plans/W-008-git-paths.md).

## Constraints

Keep selector grammar, allocation semantics, Git-required mutation boundary,
read-only command behavior, and human escaping/JSON path round trips. No shell
eval, path restrictions introduced to avoid parsing, automatic cleanup of old
stray directories, new dependencies, or changes to independent-clone scope.

## Design

Run each path-producing `rev-parse` query independently. Remove exactly its one
output newline, preserving embedded/trailing path whitespace; never split a
path-bearing result on newline or use TrimSpace. Centralize discovery used by
versions and mutations. Parse worktree inventory from `--porcelain -z` in a
shared repo helper and use it in versions and the allocator, preserving raw
paths and fields. Keep absent live record roots normal and inaccessible scans
explicit errors; retain allocator's accepted floor/recovery policy and source
scope. This work does not replace its ID prefilter with the full record parser.

## Acceptance

1. Main checkout and separate common/Git directory paths containing newline,
   tab, spaces, and trailing whitespace round-trip exactly. Nested project
   prefixes and linked worktrees do too. Exact repository/prefix/locator values
   and successful live/committed resolution are asserted, not only exit codes.
2. New/update use only the actual `<common>/grove` locks/counters. Snapshot the
   enclosing temporary parent to catch the demonstrated stray `.../new/grove`
   directory; inspection creates no coordination state. Two linked processes
   still share the same lock and ID sequence with unusual paths.
3. A highest existing live ID in an oddly named worktree contributes to allocation
   even with no persistent counter. Ref scans, unrelated prefixes, missing roots,
   and inaccessible worktrees keep their documented behavior. No path is parsed
   as shell text or silently dropped due to Git's display quoting.
4. Existing W-002/W-003 and W-004/W-005 tests plus full/race suites, vet,
   formatting, Grove `check`, and independent path/side-effect review pass.

## Dependencies and handoff

Shared repo helpers and all affected commands are integrated. No unfinished
product prerequisite. Work serially after W-006/W-007 to avoid shared helper
conflicts; this is implementation scheduling, not a semantic dependency.
Use one Fable agent in an isolated worktree; preserve the retained old worktree.
Do not clean up hypothetical stray state automatically.

## Evidence

Implemented 2026-09-19 on branch `worktree-W-006-W-008` (base `2d6de36`), not
yet integrated into main. Code: `2e18866` (`repo.GitPath`, `Locate`,
`repo.Worktrees`, versions), `9d11bdb` (allocator), `40e882f` and its correction (review fixes).
The [plan](../../docs/plans/W-008-git-paths.md#implementation-notes-2026-09-19)
records the bounded adjustments. Selector grammar, allocation and counter
protocol, the Git-required mutation boundary, and path escaping are unchanged.

Reproduced before repair on a pristine export of `2d6de36` holding only the
new tests: for a main checkout named `new\nline` the repository was reported
as `.../new` with prefix `line/.git\nsub\nproject/` and both live sources
invalid; `lock`, `next-ids`, and `write.lock` were created under the stray
sibling `.../new/grove`; allocation returned 2 instead of 91 with a live
W-090 in an oddly named checkout.

Acceptance: (1) `TestLocateExactPaths` asserts exact common directory and
prefix for a newline-named main checkout, a nested `sub\nproject/` prefix, a
linked checkout with a tab and trailing space, and a separate Git directory
named `git\tdir\nmid `; `TestOddGitPathsRoundTrip` asserts repository, prefix,
linked Git directory and locator, and resolves both live and both committed
selectors to exact checkout, project, and record paths, identically from the
linked checkout; `TestNewlineCheckoutCLI` proves raw paths in JSON and
escaped paths on stderr. (2) `TestCommonDirCreatesStateOnlyUnderTheRealCommonDirectory`
and `TestCoordinationStateStaysUnderTheCommonDirectory` list the fixtures'
enclosing directory: updates and creations from the main and a linked
checkout leave exactly `<common>/grove/{lock,next-ids,write.lock}`, share one
ID sequence, and reads create nothing. (3) `TestAllocateScansOddlyNamedWorktrees`
returns 91 then 92 at the root and at a nested prefix, ignores an unrelated
prefix, and treats a checkout without the record folder as normal;
`TestAllocateConcurrentAcrossWorktrees` now allocates concurrently from a
path Git quotes for display; `TestParseWorktrees` covers bare, detached,
locked, and prunable entries with exact fields and order. The inaccessible
scan refusal is unchanged. (4) Existing W-002 to W-005 tests pass unchanged.

Independent review (separate reviewer agent, commits `2e18866` and `9d11bdb`,
own export): no correctness finding. It confirmed on a base export that the
enclosing-directory assertion fails with the stray paths, classified every
Git call in `cmd/` and `internal/` and found no remaining newline-unsafe
parsing of path-bearing output (`git grep -h` prints no paths), and probed
record folders containing glob characters, spaces, and newlines without
lowering the ID floor. `40e882f` added concurrent allocation coverage under odd names. It also
removed a common-directory query the review called redundant; the combined
review disproved that, and the query is restored (see the combined evidence).

Limits: Git cannot reopen a separate Git directory whose name ends in a
newline, so that form cannot reach Grove; a record folder written as pathspec
magic (`:(x)dir`) makes allocation refuse loudly; the allocator still tells an
absent record root from an unreadable one by error text; old stray
directories from the defect are not cleaned up; inspection costs more Git
processes than before (10 to 39 in the measured case, 55 with a nested
prefix). Suite results are in
the [combined repair evidence](../../docs/reviews/2026-09-19-repairs-W-006-W-008.md).

## Next

Integrate branch `worktree-W-006-W-008` into main as a separate, explicit
step; this record being done asserts completion on that branch only. If a
checkout with a newline in its path was ever used before this repair, look
for and remove a stray sibling `grove/` directory by hand.
