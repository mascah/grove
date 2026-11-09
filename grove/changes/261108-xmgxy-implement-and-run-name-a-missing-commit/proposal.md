# Implement and run name a missing base or planned commit

<!-- A bug fix describes future work. It becomes approved future work when
it is merged into the target branch, which says nothing about
implementation. -->

## Observed

Found on 2026-11-08 while exploring
[issue #140](https://github.com/mascah/grove/issues/140), which reports
`grove check` after a history rewrite, and reproduced on the ledger
fixture with Grove built at 8174594. The files cited are unchanged at
ccaee4d.

1. **`grove check` names the commit in one sentence.** In a
   `git clone --depth 1` of a fixture holding change A in cycle K, change
   B in no cycle, and K, each `base` and `planned` older than the clone's
   one commit, `grove check` printed, and exited 1:

       grove/changes/<A>/.grove.yaml: base: 6f7ec537d975c7a57c507a546bd849616f3211ba is not in this shallow clone; fetch its full history, such as with git fetch --unshallow
       grove/changes/<B>/.grove.yaml: base: 6f7ec537d975c7a57c507a546bd849616f3211ba is not in this shallow clone; fetch its full history, such as with git fetch --unshallow
       grove/cycles/<K>.yaml: planned: 8f52bece87a4a93aed166a1330d2fc586072fd89 is not in this shallow clone; fetch its full history, such as with git fetch --unshallow
       grove: 3 problems found

   The sentence is `missing` in `internal/cli/check.go` at ccaee4d, used
   by `stale` and `stalePlans`; in a full clone whose commit is gone it
   reads "is not a commit in this repository". The walkthrough's step 6,
   "A shallow clone", and the CI guide's "What to expect" show it.

2. **`grove instructions implement <change-id>` prints Git's error
   instead.** In the same clone, each exiting 1:

       $ grove instructions implement <B>
       grove: git diff: fatal: bad object 6f7ec537d975c7a57c507a546bd849616f3211ba: exit status 128
       $ grove instructions implement <A>
       grove: git rev-parse: fatal: Needed a single revision: exit status 128
       $ grove instructions implement <A> --unattended
       grove: git rev-parse: fatal: Needed a single revision: exit status 128

   `instructionsChange` in `internal/cli/instructions.go` at ccaee4d runs
   `moved`, a `git diff` from the change's `base`, for a change in no
   cycle, or `sincePlan`, a `git rev-parse` of the cycle's `planned`, for
   a member of a cycle not done, and returns their error as the command's.
   A full clone made by importing the tree as a fresh repository, the
   issue's own reproduction, printed the same two lines.

3. **`grove instructions implement <K>` prints the cycle's
   instructions**, 129 lines, exit 0, although each member's own
   instructions, which those have each member's session run, fail as
   above: `instructionsCycle`, same file, reads no `planned`.

4. **`grove run <cycle-id>`** returns the same error: `runnable` in
   `internal/cli/run.go` at ccaee4d returns `changedSincePlan`'s error as
   its own. `eligible`, for a `run` given no cycle, skips such a cycle,
   with the comment "A planned that is no commit is check's to report".
   Not run on the fixture, since `run` starts a session.

## Expected

The owner, 2026-11-08: "Im fine with fixing these independently lets just
confirm first", then "yes to the bugfix", accepting the round recorded
under Decisions.

- `grove instructions implement <change-id>`, attended or `--unattended`,
  `grove instructions implement <cycle-id>` and `grove run <cycle-id>`
  refuse with the line `grove check` prints for the missing commit,
  prefixed `grove: ` as every error is, exit non-zero and print no
  instructions: for a change in no cycle, its `base`; for a cycle, or a
  member of a cycle not done, the cycle's `planned`; in the shallow-clone
  wording in a shallow clone, otherwise "is not a commit in this
  repository".
- `grove run` given no cycle keeps skipping such a cycle silently.

## Reproduction

From a clean checkout, with `grove` built from it:

1. `just fixture /tmp/f`; in `/tmp/f`, `grove init` and commit.
2. Write two approved changes, A and B, each with `base:` main's commit
   and a task whose `Files:` names `ledger/cli.py`, and commit. Write
   cycle K listing A with `planned:` main's commit, and commit.
3. `git clone --depth 1 file:///tmp/f /tmp/s`. In `/tmp/s`, run
   `grove check`, `grove instructions implement <B>`,
   `grove instructions implement <A>`,
   `grove instructions implement <A> --unattended` and
   `grove instructions implement <K>`.
4. For the other wording, in `/tmp/f` set A's `base:` to `0000000` and
   run the same commands.

## The requirement it breaks

None at ccaee4d: no maintained spec says what implement or run do when a
`base` or `planned` commit is missing. `check`'s sentence is observed
behavior, kept so by
[261105-3qdz2](../archive/261105-3qdz2-misnamed-cycle-and-shallow-clone-reported/proposal.md)
(decision 1). The owner kept this observed too (decision 3), so there is
no change spec.

## Decisions

The owner accepted every item below with "yes to the bugfix", 2026-11-08.

1. **Scope:** `grove instructions implement`, for a change and for a
   cycle, attended and `--unattended`, and `grove run <cycle-id>`.
   Approved change
   [261108-jpg8n](../261108-jpg8n-claim-replaces-run/proposal.md) also
   names `internal/cli/run.go`; whichever of the two lands second
   re-bases with propose.
2. **Behavior:** each refuses with `check`'s exact line and exits
   non-zero, printing no instructions. The cycle form refuses when the
   cycle's `planned` is missing. A `run` given no cycle keeps skipping
   the cycle silently.
3. **Observed, not agreed:** no change spec, as 261105-3qdz2 decided for
   `check`'s message. The CI guide's "What to expect" says the implement
   instructions fail the same way; walkthrough step 6 gains the implement
   and run lines.
4. **No `depends_on`.**
5. **Routine defaults:** one task starting with a test in
   `TestShallowClone`; its files as `tasks.md` names them; test layout,
   private helpers and `internal/cli/check.go` delegated, the last should
   the `missing` sentence need a path prefix.

## Open questions

None. A related observation is outside this change: a checkout with
`origin/main` fetched but no local `main`, as a pull-request job's
checkout is, fails every command that reads the target with
`grove: git rev-parse: fatal: Needed a single revision: exit status 128`.

## Related documents

- Tasks: [tasks.md](tasks.md)
- Requirements: none; the owner kept this observed behavior (decision 3).
