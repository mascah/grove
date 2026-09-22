---
id: "G-079"
type: work
title: "Update a record by hand without a lookup and commit it in one step"
status: review
created: "2026-09-22T16:01:52Z"
updated: "2026-09-22T16:26:30Z"
kind: tooling
size: small
relates_to: ["G-009", "G-044", "G-062"]
candidate: "9eda6c6"
---

## Outcome

The owner, marking work done or changing a field at a shell, runs one
`update` command and is finished: no `show --json` first, no revision to copy,
and no separate commit when the change is the record alone. Agent sessions
keep the revision guard they have today.

The owner asked for this on 2026-09-22 after two integrations in one week
needed a lookup, a paste and a hand commit
([G-078](G-078-g-039-trial-evidence-for-the-int.md) finding 3, and the G-038
and G-076 done commits `b8fb232` and `e8bcf03`). In the same session they
chose the two designs below: make `--expect` optional, and add an opt-in
`--commit` with a generated message.

## Constraints

Observed at main `ccdc92d`:

- `update` refuses a missing `--expect` as a usage error
  (`internal/cli/cli.go:444`). [G-009](G-009-update-records.md)'s design,
  which the owner approved on 2026-09-19, says "requires exactly one
  `--expect REVISION`", "No force/ignore-revision option", and "Do not change
  configuration, refs, or index". This work revises those three sentences on
  the owner's 2026-09-22 choice; the rest of G-009's contract, including the
  write lock, the pre-publication re-read that refuses an observed change,
  byte preservation and the failure reporting, stays as it is.
- The guard exists for a caller whose read is stale
  ([G-062](G-062-revision.md)). A person at a shell reads and writes within
  seconds, and the write lock already serializes cooperating Grove commands,
  so the shell form `--expect "$(show --json …)"` bypasses the guard while
  keeping its cost. An agent session that read the record through `context`
  minutes earlier has the revision at hand and a real reason to pass it.
- `update` already runs Git for the `done` ancestry check
  (`internal/update/update.go:259` through `repo.Git`), and `new` requires
  Git; a commit adds Git writes to a command that already reads Git.
- Every workflow commit is deliberate today: the work guide commits the
  review handoff "alone" and the done change "there"
  (`docs/work-execution.md`), the shaping guide commits "only when the person
  agrees" (`docs/work-shaping.md`), and the done commits differ in shape:
  `e8bcf03` changed one record by two lines, `ccdc92d` changed two records.
  A commit that swept the working tree would break all three, so the commit
  is scoped to the updated record's file and is never the default.
- `--expect` is written into `README.md`, `docs/record-model.md` ("First
  commands"), `docs/work-execution.md`, `docs/work-shaping.md` and the term
  G-062, and the G-009 record; the record and the term describe history and
  are not rewritten.

Design, proposed:

- `--expect` optional. Omitted, the update applies to the content the file
  holds under the write lock, with every other check unchanged; a stale
  `--expect` that is given is refused exactly as now. The guides keep the
  flag for agent sessions in their `update` invocations, with one sentence
  saying why: their read may be old.
- `--commit`, boolean, opt-in. After a changed publication it runs
  `git add -- PATH` and `git commit -m MESSAGE -- PATH` for the record's
  file alone, leaving other staged or unstaged files as they are, and prints
  the commit in the result (`{id, path, revision, changed, commit}`). The
  message is generated from the request, `docs(G-076): set status=done`,
  with several fields joined by spaces and `--unset size` as `unset size`;
  a reasoned message is still written with Git by hand. A no-op update
  commits nothing and says so. Uncommitted edits already in that file, such
  as a hand-written verdict, are part of the commit, which is the point.
- Refusals: `--commit` without `update` or outside a Git work tree is a
  usage error; a commit that fails after the file was replaced (identity
  unset, a hook, a rebase in progress) reports the update as applied and
  not committed, with the revision, in the same way G-009 reports a failure
  after the rename. `--commit` never runs `git commit -a`, never touches
  other paths, and never pushes.
- Out of scope: a `--message` flag, `--commit` on `new`, any approval or
  merge action ([G-044](G-044-review-integration.md) owns those), and a
  `--no-commit` default flip.

## Acceptance

1. From a clean `main`, `update G-NNN --set status=done --commit` with no
   `--expect` on a record whose candidate is an ancestor of HEAD writes the
   status, commits that one file, and prints the commit; `git show --stat`
   of it lists only that record.
2. `update … --expect REVISION` with a stale revision is still refused with
   exit 1 and no write, and the omitted form applies to a record that a
   direct editor changed after the last commit.
3. `--commit` with another file modified in the working tree commits only
   the record; with a staged change to another file, that change stays
   staged and uncommitted. A no-op update with `--commit` makes no commit.
4. A commit failure after publication reports the applied revision and that
   nothing was committed, exit 1, and the working file holds the update.
5. The guides and the README/record model say `--expect` is optional for a
   person and kept for agent sessions, and every `update` example in
   `docs/work-execution.md` and `docs/work-shaping.md` still runs as written.
6. The owner marks one real record done with the single command from a
   checkout of `main` and judges that the lookup and the hand commit are
   gone.

## Evidence

Implemented on `worktree-G-079` (`.claude/worktrees/G-079`), base `main` at
`17667d6`, from this record at revision `ada396e6…` with no plan, as Next
said. Commits: `4ef3129` active, `713de02` the code and tests, `5e4b3a1` the
four documents, `34ca80f` the review fixes, then this evidence. The
candidate is the commit this record's frontmatter names.

Behavior against the acceptance, at `34ca80f`:

1. `update G-NNN --set status=done --set candidate=C --commit` with no
   `--expect` writes the status, commits that one file, and prints
   `{…, "commit": SHA}`. Fixture: `TestUpdateOptionalExpectAndCommit`
   (`internal/update`) and `TestUpdateCommitResultAndFailure`
   (`internal/cli`). Real use: in a disposable clone of this branch at
   `5e4b3a1`, `update G-079 --set status=done --set candidate=5e4b3a1
   --commit` made `e5e8c87` "docs(G-079): set status=done candidate=5e4b3a1"
   whose `git show --stat` lists this record alone.
2. A stale `--expect` is refused, exit 1, no write, with the unchanged
   message (fixture and the same clone). The omitted form applied over a
   body edited after the last commit (fixture).
3. With another file modified and a third staged, the commit holds the
   record alone and `git status --porcelain` is unchanged for the others
   (fixture; the reviewer also tried a staged-plus-unstaged edit and a
   staged deletion). A no-op with `--commit` makes no commit and prints
   `"commit": null` (fixture and clone).
4. A refusing `pre-commit` hook after publication returns the `*Failure`
   with the applied revision, exit 1 through the CLI as
   `… the file is staged but nothing was committed (the update was applied
   to PATH; revision …)`, and the file holds the update (fixtures).
5. README, `docs/record-model.md`, `docs/work-execution.md` and
   `docs/work-shaping.md` say `--expect` is optional for a person and kept
   by a session because its read may be old; every `update` example in the
   two guides is unchanged and still runs (the reviewer checked each and
   every link). G-009 and G-062 are history and were not rewritten.
6. Open: the owner's own use from a checkout of `main`, at integration.

Decisions taken: the message is `docs(ID): set a=b c=d unset e`, the
request's fields in order, with `\r` and `\n` in values replaced by
spaces; the JSON key is `commit`, present only when `--commit` was given,
`null` when nothing changed. `--commit` outside Git is refused by `update`'s
existing Git requirement, exit 1, not a usage error (review finding 1).
`git add` precedes `git commit -- PATH` so a record `new` created and never
committed can be committed too; the pathspec is `:(literal)`. The commit runs
under the write lock, so a slow hook delays other Grove writes (finding 4).

Verification at `34ca80f`: `go test -short -count=1 ./internal/update
./internal/cli` ok; `go vet ./...`, `gofmt -l .` clean; `go run ./cmd/grove
check` OK (80 records); `go test -count=1 -timeout 120s ./...` all ok.
`versions` (8.0 s) exceeds the five-second budget under the parallel suite
as before this change; `update` (3.5 s) and `cli` (4.6 s) stay under.

Review: [G-083](G-083-g-079-review.md), independent, examined `5e4b3a1` and
the fix `34ca80f`: no blocking findings; 2, 3, 5 fixed, 1, 4, 6 accepted.

## Next

Candidate awaits the owner's judgment. Runnable demo from any checkout of
the branch, in a disposable clone so the ID counter and the branch stay
clean:

```sh
git clone -q -b worktree-G-079 /Users/mascah/GitHub/mascah/grove /tmp/g079 && cd /tmp/g079
go run ./cmd/grove update G-079 --set status=done --set candidate=$(git rev-parse --short HEAD~1) --commit
git show --stat HEAD
```

To integrate (acceptance 6 is this very use): in the `main` checkout,

```sh
git merge --ff-only worktree-G-079
# edit grove/G-079-update-a-record-by-hand-without.md: quote the verdict in this Next
go run ./cmd/grove update G-079 --set status=done --commit
```

The candidate in the frontmatter is an ancestor of `main` after the
fast-forward, so `done` needs no `--set candidate`; a squash or rebase that
landed another commit names it with `--set candidate=COMMIT` in the same
call.
