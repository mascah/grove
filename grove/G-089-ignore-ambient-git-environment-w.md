---
id: "G-089"
type: work
title: "Ignore ambient Git environment when Grove and its tests run Git"
status: active
created: "2026-09-22T17:11:12Z"
updated: "2026-09-22T17:19:18Z"
size: small
---

## Outcome

Grove and its test suite act only on the repository they are pointed at,
whatever Git environment they inherit. A Git hook, an editor, or a shell that
exports `GIT_DIR`, `GIT_WORK_TREE`, or the other variables that name a
repository can no longer redirect a Git process Grove starts, and the tests
prove it. The owner asked for this on 2026-09-22 after the incident below and
called it a significant risk to the code.

## Constraints

**Observed on 2026-09-22:**

- A `git push` from the linked worktree `.claude/worktrees/G-081` ran
  lefthook's pre-push hook, `go test ./...`. Git exports
  `GIT_DIR=<repo>/.git/worktrees/G-081` to the hook (shown with a scratch
  `core.hooksPath` hook on `push --dry-run`), and every test fixture's
  `git init`, `commit`, `checkout`, and `worktree add` in a temp directory
  inherited it. The fixtures therefore acted on the real repository: the
  branch tip was replaced by a fixture commit and pushed to origin as PR #2's
  head, HEAD moved to a fixture branch, branches `bad-yaml`, `code`,
  `feature`, and `records` appeared, two temp worktrees were registered, and
  `core.bare = true` landed in the shared config, which made the main checkout
  refuse work-tree commands. The commits were intact and the owner repaired
  it by hand ([G-081](G-081-github-ci.md) has the commands).
- Product code spawned Git in three places, all `git -C dir …`, which
  `GIT_DIR` overrides just the same: `repo.GitContext`, the `cat-file --batch`
  reader in `internal/versions`, and `git grep` in `internal/create`. Nine
  test files spawned Git directly, and `internal/tui/testdata/terminal.py`
  spawns it from Python.

**Selected:** one guard at the root, `repo.Command`, that every Git process
goes through, rather than a per-call or hook-only fix; the hook keeps a
second, independent scrub. Variables dropped: `GIT_DIR`, `GIT_WORK_TREE`,
`GIT_INDEX_FILE`, `GIT_COMMON_DIR`, `GIT_OBJECT_DIRECTORY`,
`GIT_ALTERNATE_OBJECT_DIRECTORIES`, `GIT_NAMESPACE`. Setting `GIT_DIR` to
point Grove at a repository other than its project path was never supported.

## Acceptance

1. No `exec.Command("git", …)` or `exec.CommandContext(…, "git", …)` remains
   outside `repo.Command` and its regression test; `grep` shows it.
2. A regression test in `internal/repo` sets `GIT_DIR` and `GIT_WORK_TREE`
   to a decoy repository and shows `repo.Command` acting on its own
   directory while an unguarded child lands in the decoy.
3. The whole suite, run with `GIT_DIR` and `GIT_WORK_TREE` exported to a
   decoy repository, passes and leaves the decoy's refs, branches, and
   worktrees unchanged.
4. lefthook's pre-push hook runs the suite with those variables unset.
5. `AGENTS.md` states the rule and names this record.

## Next
