---
id: "W-008"
type: work
title: "Preserve Git paths through discovery and coordination"
status: proposed
created: "2026-09-19T20:13:57Z"
updated: "2026-09-19T20:21:17Z"
kind: fix
priority: 2
size: medium
relates_to: ["W-002", "W-003", "W-004", "W-005"]
---

## Outcome

Git paths containing control characters retain their exact identity through
inspection, workspace resolution, and mutation coordination. All mutation locks
and counters stay under the actual Git common directory. Proposed repair of
existing path/side-effect contracts. [Review R6](../../docs/reviews/2026-09-19-integrated-cli.md)
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

## Next

Fable: reproduce R6 in a disposable repository, including the outside-repository
lock path, then implement the linked plan and prove the surrounding parent is
unchanged except for expected fixture outputs. Return focused commits and
review evidence; integration is separate.
