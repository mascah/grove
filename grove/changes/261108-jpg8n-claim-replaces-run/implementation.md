# Implementation: grove claim replaces grove run, with attempt records on a branch

<!-- Progress on change 261108-jpg8n, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — `grove claim` selects, records and prepares

Status: done

`internal/cli/run.go` became `internal/cli/claim.go`. `claim` fetches the
target and `grove/runs` from `origin`, using a pattern refspec so that a
project with no record yet fetches without error. It refuses an unclean
checkout, a local target that cannot be fast-forwarded and, without
`--retry`, a latest record whose cycle has a member not archived, naming
that attempt, its host, start time and branch. It selects the given or
oldest eligible cycle as `run` did. It writes the record as a commit with
the empty tree (`hash-object -t tree -w /dev/null`, then `commit-tree`) and
pushes it with `--force-with-lease=refs/heads/grove/runs:<tip just
fetched>`, or, with no `origin`, runs `update-ref` against the tip just
read. When the update is refused and `grove/runs` has moved, it decides
again as a fresh claim; otherwise it returns Git's error. Then it creates
`grove/<attempt>` from the fetched target, fast-forwards the local target
and prints the cycle ID alone. With nothing to claim it prints nothing on
stdout and exits 1. The lock, `--timeout`, the Claude Code launcher, the
relay, the process-group stop and `.git/grove/` are gone, and so is the
requirement that `agents` list Claude Code. `cli.go` has
`claim [--retry] [<cycle-id>]`, and `run` falls through to the usage with
status 2. `check.go`'s message for a marked stale cycle says `grove claim`
skips it. `internal/gitx/gitx.go` needed no change.

`TestRun` became `TestClaim` in `internal/cli/walkthrough_test.go`, with a
bare `origin` and two clones, h1 and h2. For two claims at once, h2's
`pre-push` hook runs h1's claim (the test binary run as grove through
`TestMain`) while h2's push is under way. h1's record wins, h2's lease is
refused, and h2 reads again, finds the project paused and creates no
branch. The test also covers:

- `grove run` returning the usage's status 2;
- nothing to claim: exit 1, empty stdout;
- `agents` listing only Codex;
- the local-target and unclean-checkout refusals, the latter with
  `--retry`, the edit kept;
- a push refused by `origin`'s `pre-receive` hook: nothing recorded, no
  remote-tracking ref, no branch;
- the record's fields and empty tree, and no `.git/grove/`;
- the pause on the same host, and its end after one member archived as
  finished and one as dropped;
- a given cycle refused (no such cycle, a member needs a person, paused)
  and claimed with `--retry`, keeping the earlier record and its branch;
- no `origin`: the record on the local `grove/runs`, and the pause ended by
  merging the attempt's branch.

TestFollows, TestReplan and TestMovedFiles name `grove claim`.

`just check` passes: 4303 Go lines, 1941 test lines, 408 page lines, 8614
instruction words. On a fresh fixture, `bin/grove run` exits 2 with the
usage, and `bin/grove claim` prints nothing on stdout, says "nothing to
claim on main at 5334772" and exits 1.

Also changed: `internal/cli/serve.go` — it read `.git/grove/`, which is gone; it reads the local `grove/runs` with `claim`'s `records` until T-02 has it prefer `origin/grove/runs`.
Also changed: `internal/doc/doc.go` — check's `needs_human` message and a comment said `grove run` skips the cycle, which is false once `run` is no command.

### T-02 — `serve` reads the records and names `grove claim`

Status: done

`internal/cli/serve.go` reads the records on `origin/grove/runs`, or on the
local `grove/runs` when there is no `origin/grove/runs`. Each run on the
page is the record plus `moved`: the latest commit time of its branch,
local or on `origin`. A cycle's verdict, after the existing refusals from
`runnable`, is one of these:

- "claimed by attempt X, which H started at T on branch B", when the latest
  record names it and a member is still open;
- otherwise `grove claim would claim it`, with the cycle's ID when it is
  not marked and "; it is not marked unattended";
- while the project is paused after another cycle's attempt, `--retry` in
  the command and "; the project is paused after attempt X of cycle Y",
  with `runs` false.

The pause test is now one helper, `pausing` in `claim.go`, used by both
`claim` and `serve`. `serve`'s branches no longer list `grove/runs`, local
or on `origin`. Without that exclusion, a `grove/runs` branch broke
`/api/project`, because the record branch shares no history with the
target: `git diff: ... no merge base`. `internal/cli/page.html` heads
readiness "What grove claim would do". Its runs table, "Attempts recorded on
grove/runs", shows attempt, cycle, started, host, branch and last commit.

`TestServe` checks the following:

- K's verdict is "grove claim would claim it" with no runs.
- After a record from host `h2.example` is pushed to `origin/grove/runs`,
  with its branch pushed to `origin` and no local attempt, K is claimed,
  naming the host, time and branch. The run carries its host, cycle,
  branch, start and the branch's commit time.
- With a local copy of `grove/runs` (not listed as a branch) and a second
  record naming L, K says `grove claim --retry would claim it; the project
  is paused after attempt 261106-r0000 of cycle L`, and the runs are newest
  first.

`just check` passes: 4338 Go lines, 1977 test lines, 408 page lines, 8614
instruction words. How the runs table looks in a browser was not looked
at. Its cells are the ones the old table had, apart from the host and last
commit columns.

Also changed: `internal/cli/claim.go` — the pause test moved into `pausing`, so `claim` and `serve` share it.

### T-03 — Instruction text names `grove claim`

Status: not started

### T-04 — Documentation and the walkthrough

Status: not started

## Scenarios

<!-- Before finishing, one line per scenario of the change specs: the check
that shows it and its result, or `unverified` and why:

    - record-editing "Update the note on a legacy record": tests/test_cli.py test_edit_legacy, passed
-->

## Reconciliation

<!-- One line per change spec, once applied to its maintained spec:

    - record-editing: applied to specs/record-editing/spec.md
-->
