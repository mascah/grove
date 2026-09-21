---
id: "G-020"
type: plan
title: "G-016 Git path identity implementation plan"
status: current
formerly: "docs/plans/W-008-git-paths.md"
work: ["G-016"]
updated: "2026-09-21T21:11:16Z"
---

# G-016 Git path identity implementation plan

> For Fable: implement one task at a time in an isolated worktree using the
> repository instructions and an inline execution workflow (for Codex,
> `superpowers:executing-plans`). Do not dispatch or merge automatically.

Review base: `9b7f730`; include the review/shaping documents in the execution
base. Read the linked work record first; it owns outcome, constraints, and
acceptance. This is a proposed technical repair plan, not implementation evidence.
Go 1.26+, current standard library and existing YAML dependency only.

**Goal:** preserve exact path identity and keep all coordination writes under
the actual common Git directory, including paths containing newlines.

**Architecture:** centralize single-path Git queries and NUL-delimited worktree
inventory in `internal/repo`. Keep `repo.Locate(root)`'s signature and all public
CLI/selector formats. Versions and creation consume the shared inventory;
source validation remains with their existing owners.

**Spec:** [G-016](G-016-git-paths.md).

## Global constraints and review focus

No display-path parsing, TrimSpace on paths, shell eval, schema change, migration,
or automatic stray-directory cleanup. Preserve recovery/ID-floor scope and
reader side-effect guarantees. Review: newline in the main/common directory,
trailing whitespace in a separate Git directory, newline in a project prefix,
untracked high IDs in a quoted worktree path, and read/write state containment.

## Task 1: Read each Git path without ambiguous framing

Files: `internal/repo/repo.go`, new `internal/repo/path_test.go`,
`internal/versions/versions.go`, `versions_test.go`, `workspace_test.go`,
`internal/update/update_test.go`.
Consumes `repo.Git`; keep `Locate` returning `(common, prefix string, err error)`.
Introduce `repo.GitPath(dir, option string) (string, error)` for one rev-parse
path query, reused for `--git-dir`, `--git-common-dir`, and `--show-prefix`.

- [ ] Build fixtures under a temporary parent with main checkout `new\nline`,
  linked checkout `feature\tline`, nested prefix `sub\nproject`, and a separate
  Git directory ending in a space/newline. Use `git init --separate-git-dir`
  where appropriate; generate all arguments as strings, never shell text.
- [ ] Assert exact common/prefix/Git directory values and resolve both live and
  committed selectors. On the base, `new\nline` becomes repository `.../new`
  and prefix `line/.git`, so a live source is incorrectly invalid.
- [ ] Implement single-path framing:

```go
out, err := Git(dir, "rev-parse", "--path-format=absolute", option)
if err != nil { return "", err }
if !strings.HasSuffix(out, "\n") {
    return "", fmt.Errorf("git rev-parse %s: missing output terminator", option)
}
return strings.TrimSuffix(out, "\n"), nil
```

Call separately for every path; retain an empty prefix for a root project.
Replace versions' combined Git-directory query with these exact path values.
Validate path discovery before `CommonDir` creates coordination directories.
- [ ] Assert `update` creates only the correct `.git/grove/write.lock` under
  `new\nline`; the sibling `new/grove` must not exist. Snapshot the fixture's
  enclosing parent, not just repository files. Repeat from a linked checkout
  and prove both use the same common lock. Run repo/update/versions targeted
  tests; commit `fix(repo): preserve exact paths returned by Git`.

## Task 2: Share NUL-delimited worktree inventory with the allocator

Files: new `internal/repo/worktrees.go` and `worktrees_test.go`,
`internal/versions/versions.go`, `internal/create/create.go`, `create_test.go`.
Produce `repo.Worktrees(root string) ([]repo.Worktree, error)` with raw
`Path`, `Head`, `Branch`, `Prunable` strings and `Bare` bool. Move the existing
versions NUL parser without dropping administrative fields; preserve consumers'
error/absence policies. No path quoting/unquoting heuristics.

- [ ] Test parser input with records separated by double NUL and fields by NUL,
  paths containing newline/tab/backslash/quotes/spaces, and bare, detached,
  locked, prunable entries. Confirm fields are exact and ordering preserved.
- [ ] In an allocation fixture, commit W-001, create an oddly named linked
  checkout, and place a normal generated-format W-090 record there untracked
  with no counter file. `Allocate` from main must return 91; repeat with a
  nested project prefix and source configuration matching the accepted allocator
  scope. Also verify absent record roots versus inaccessible scan errors.
  Use fixture-only record bytes, not hand-numbered project work records.
- [ ] Switch versions and allocation to `repo.Worktrees`; remove the allocator's
  display-oriented `--porcelain` split. Keep the committed-ref ID prefilter and
  counter protocol unchanged. Reuse G-014's new administrative validity checks
  if that work is in the execution base.
- [ ] Run concurrent allocation and lock tests from oddly named worktrees, the
  entire create/versions suites, and side-effect snapshots. Commit
  `fix(create): scan worktree paths from NUL-delimited inventory`.

## Verification and return

- [ ] Run targeted package/CLI regressions with `-count=1`, then
  `go test -count=1 ./...`, `go test -race -count=1 ./...`, `go vet ./...`,
  `gofmt -l cmd internal`, and `go run ./cmd/grove check`.
- [ ] Obtain an independent review of the final diff, with the original
  reproducer run against both base and candidate. Preserve concrete failing
  and passing outputs; test names or a done status alone are not evidence.
- [ ] Reconcile the work record, README/model if behavior wording changes,
  and the restart brief's concrete next action. Change fields with
  `go run ./cmd/grove update`; keep evidence in the record/plan bodies.
- [ ] Return focused Conventional Commits, tested candidate revision, review
  disposition, and remaining limits. Do not merge or remove old worktrees.

## Implementation notes, 2026-09-19

Implemented on branch `worktree-W-006-W-008`; the
[work record](G-016-git-paths.md) owns the evidence.
Bounded adjustments to this plan, and why:

- The separate Git directory fixture is named `git\tdir\nmid ` (embedded
  newline, tab, trailing space). Git itself cannot reopen a separate Git
  directory whose name ends in a newline: `git init --separate-git-dir
  $'trail\n'` succeeds, but the `.git` file it writes loses that terminator
  and every later command fails with `not a git repository: .../trail`. Grove
  cannot be given such a repository, so there is nothing to round-trip.
- G-014's `identity` helper had introduced a third newline-splitting
  `rev-parse`; it now uses `repo.GitPath` per path. Comparing the Git directory
  is enough for ownership of a location, because a Git directory belongs to
  one worktree of one repository, so the common directory is asked for only
  when entering a worktree.
- `repo.Worktree` carries `Path`, `Head`, `Branch`, `Prunable`, and `Bare`.
  `locked` is parsed past and not retained: nothing consumes it.
- The allocator keeps its policy of scanning every inventory entry; a bare
  entry, a prunable entry, or a checkout without the record folder is a
  missing root and normal, while an unreadable record stays an error.
