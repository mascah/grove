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

## Evidence

Implemented 2026-09-22 on `worktree-G-089` from `main` adce678, the same
day as the incident, at the owner's request, with no plan: the fix is one
function and its callers. Commits: 7924765 (the guard, every caller, the
regression test, the hook, the terminal script, `AGENTS.md`), 263a59a
(status), ca6a599 (review fixes), then the commit carrying this evidence and
[G-090](G-090-review-of-g-089-ambient-git-envi.md), which is the `candidate`.

**Against the acceptance:**

1. `grep -rn 'exec\.Command' --include='*.go'` finds Git spawned only in
   `repo.Command` and the regression test's deliberate unguarded control;
   the other hits run `go`, `python3`, and the test binary.
2. `TestCommandIgnoresAmbientRepository` in `internal/repo/command_test.go`
   points `GIT_DIR` and `GIT_WORK_TREE` at a decoy: `repo.Command` reports
   its own `.git`, the unguarded control reports the decoy's. The reviewer
   showed it fails with the `cmd.Env` line removed.
3. `GIT_DIR=<decoy>/.git GIT_WORK_TREE=<decoy> GOFLAGS=-buildvcs=false go test
   -count=1 -timeout 120s ./...`: all packages ok; the decoy kept one commit,
   one branch, no worktrees, its refs hash unchanged; the worktree's status
   unchanged. Run by this session and again by the reviewer at ca6a599.
4. `lefthook.yml` pre-push runs `env -u` for all seven variables; the
   reviewer ran that line under lefthook 2.1.8 from a linked worktree and
   the child saw none of them while keeping `GIT_COMMITTER_DATE`.
5. `AGENTS.md` states the rule in the development-policy list and names
   this record.

**Decisions:** the seven variables are the ones that name a repository or
its parts; `GIT_COMMITTER_*`, `GIT_AUTHOR_*`, `GIT_TRACE`,
`GIT_CONFIG_GLOBAL`, and `GIT_CONFIG_SYSTEM` still reach the child because
tests and existing evidence rely on them. `GIT_CONFIG_PARAMETERS`,
`GIT_PREFIX`, and `GIT_CEILING_DIRECTORIES` were tested and do not move
`git -C dir`. Pre-commit jobs are not scrubbed, with a comment saying why.
The history test's merge gained the committer identity it lacked, which
Linux CI in G-081 had exposed; that line conflicts textually with G-081's
own fix and is resolved by taking either.

**Verification** at ca6a599 on macOS, Go 1.26.5: `gofmt -l .` empty,
`go vet ./...`, `go run ./cmd/grove check` (`OK: 85 records`), plain
`go test -count=1 -timeout 120s ./...` all ok (`internal/versions` 7.3 s),
plus the decoy run above.

**Review:** G-090, two rounds, no blocking findings; the should-fix items
(hook variable list, `AGENTS.md` scope) and nits are fixed in ca6a599.

## Next

In Review. This branch has no `.github` workflow, so its pull request runs CI
only once [G-081](G-081-github-ci.md) is merged. Integrator, from the main
checkout, after G-081:

1. `git merge worktree-G-089`; the one conflict, if any, is the merge line in
   `internal/versions/history_test.go`: keep the `repo.Command` form with the
   `-c user.name=t -c user.email=t@t` flags, and keep G-081's `MERGE_HEAD`
   assertion beneath it.
2. `go test -count=1 -timeout 120s ./...`, then `git push origin main`; the
   hook is safe again from this merge on.
3. Write the verdict here and `go run ./cmd/grove update G-089 --expect
   REVISION --set status=done`, committed on `main`.
