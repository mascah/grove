---
id: "G-083"
type: review
title: "G-079 review of optional --expect and --commit"
status: current
created: "2026-09-22T16:18:22Z"
updated: "2026-09-22T16:18:30Z"
work: ["G-079"]
examined: "5e4b3a1"
---

## Examined

Independent review by a reviewer agent that did not edit the interfaces under
review, on 2026-09-22, of `git diff 17667d6 5e4b3a1` on `worktree-G-079`
against [G-079](G-079-update-a-record-by-hand-without.md)'s Design and
Acceptance: `internal/update/update.go`, `internal/cli/cli.go`, their tests,
`README.md`, `docs/record-model.md`, `docs/work-execution.md` and
`docs/work-shaping.md`. The reviewer ran `go test -short -count=1
./internal/update ./internal/cli` (both ok), `gofmt -l .`, `go vet ./...`,
`go run ./cmd/grove check` (OK: 80 records), and built the branch's CLI to
drive six scratch repositories: an untracked record, an unborn HEAD, a project
root in a subdirectory of the repository, a merge in progress, a tree with a
staged modification plus a further unstaged edit and a staged deletion, and a
non-Git project. A second pass examined the fix commit `34ca80f`.

## Findings

No blocking findings. The isolation property held under every tree state the
reviewer built: `git show --name-only HEAD` listed the record alone, `git
status --porcelain` was byte-identical before and after, a staged version of
another file stayed the staged version, and a staged deletion did not ride
along. The omitted `--expect` runs every other check inside the same write
lock, and a stale one is refused with the unchanged message.

1. Nit. `--commit` outside a Git work tree exits 1, not the Design's usage
   error 2, because `update` has always refused to run outside Git with exit 1
   before any write, and `--commit` is checked after that. Accepted: exit 2
   there would be inconsistent with plain `update`; the Design sentence is
   read as "refused", which it is.
2. Should-fix. A commit that fails after publication (a merge in progress, a
   refusing hook) leaves the record staged by the preceding `git add`, and
   "nothing was committed" read as if the index were untouched. Fixed in
   `34ca80f`: the message says the file is staged but nothing was committed,
   the README says the same, and no reset is attempted, since `git add` may
   have replaced a partially staged version of that file.
3. Nit. The pathspec was not literal, so a hand-made directory name holding
   `[`, `*` or `?` was reasoned to become a glob. On re-review the reviewer
   tried `grove/weird[1]/` and `grove/x[!x]y/` against the old binary and
   both committed the correct file: the finding did not reproduce. The
   `:(literal)` prefix from `34ca80f` is kept at zero cost so nobody has to
   re-derive why the bare form is safe.
4. Nit. The write lock is held across the commit, so a slow hook blocks other
   Grove writes in every linked worktree. Accepted for a single user; the
   upgrade is to unlock before the commit.
5. Nit. A newline in a `--set` value made a multi-line commit subject. Fixed in
   `34ca80f`: `\r` and `\n` in values become spaces in the message.
6. Nit, for the owner. The result carries `commit` only when `--commit` was
   given, `null` on a no-op, a SHA otherwise; the Design's `{id, path,
   revision, changed, commit}` reads as always present. Deliberate, documented
   in the record model and asserted both ways in the CLI test; not changed.

A pre-existing exposure, not introduced here: the fixtures set identity and
`core.hooksPath` locally but not `GIT_CONFIG_GLOBAL`, so a developer with a
global hooks path or commit template would see these commits fail, as with
every other Git test in the package.

## Disposition

Findings 2, 3 and 5 fixed in `34ca80f`; 1, 4 and 6 accepted with the reasons
above. Re-review of `34ca80f`: clear. The reviewer reran the untracked-record,
subdirectory-root, staged-plus-unstaged, staged-deletion, no-op and
merge-in-progress cases against the fixed binary with the same results, saw
the staged message be true under a real merge conflict, and a newline value
produce a one-line subject while the record keeps the newline. One remaining
nit: the README sentence about a refused `--commit` said the file "is
staged", which is false on the near-unreachable branch where `git add` itself
fails; reworded in the evidence commit to say the message tells which.
