# Implementation: Reading the project value starts a bounded number of git processes

<!-- Progress on change 261122-g4243, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — The project value starts a bounded number of git processes

Status: done

What changed:

- `internal/gitx/gitx.go`: `Input` runs Git with given standard input;
  `Objects` reads names in one `git cat-file --batch`, giving each object's
  ID and contents, or neither for a name that is no object.
- `internal/cli/check.go`: `atBases` reads every change's
  `<base>^{commit}` and the maintained specs it touches at base in one
  `cat-file --batch`; `rewritten` reads from that instead of a whole tree
  per change. `movedSinceBase` diffs the target against each distinct full
  base in one `git diff-tree --stdin --always -r --name-only --no-renames`,
  through `diffTree`, which matches each block of output to its request by
  the header `--always` prints. `moved` takes `--no-renames`. `diagnose`
  takes the target's commit and documents; `stalePlans` and
  `movedSinceBase` read the target only when they lack it. `changedSincePlan`
  remembers `sincePlan`'s paths for each cycle while one project value is
  built.
- `internal/cli/serve.go`: unmerged branches are diffed in one
  `diff-tree --stdin --raw` (tip, then target, for each branch not scanned
  before); run records' branch dates come from one `for-each-ref`; the value
  passes its target to `diagnose`.

Checks run:

- `TestGitProcesses` (new, `internal/cli/walkthrough_test.go`): with a git
  shim on `PATH`, `grove changes --json` and `grove check` start the same
  number of processes with one change, record and plain branch as with
  three of each, and the second value holds the third change's moved-file
  and rewritten-requirement problems. At 522f060 it failed: 19 against 29
  processes for `changes --json`, 6 against 12 for `check`. Now it passes.
- `TestMovedFiles` gained three cases: `grove check` with no target branch,
  no cycle and no task naming a file says "No problems"; a base written as
  a 7-character ID still reports a moved file; a renamed file whose old
  path T-01 names and new path T-02 names gives the problem
  "ledger/report.py and ledger/reports.py changed on main since base …",
  and `grove instructions implement` lists "T-01's `ledger/report.py`,
  T-02's `ledger/reports.py`". At 522f060 the rename case failed (only the
  new path); the other two passed there and pass now.
- The proposal's reproduction, on a clone of this repository with its
  remote removed and 5 extra branches: at 522f060, 16 processes, then 46
  after moving the 6 newest archived changes back; with this change, 13,
  then 16. The 3 extra are one cycle's plan check (`rev-parse`, `log`,
  `cat-file -e`): reopening the changes made one cycle not done, which the
  proposal's Expected section allows.
- `just check`: passes; 4,839 Go lines, 2,714 test lines.

Check: the open question on where the batching lives. I added two exported
`gitx` functions rather than the one recommended: `Input`, which runs Git
with standard input, and `Objects`, the `cat-file --batch` reader, which
`TreeFS` now shares so that reply format is parsed in one place. Parsing
`diff-tree` output stays next to its callers, in `diffTree` in
`internal/cli/check.go`.
Check: `changedSincePlan` remembers its answers in an unexported `plans`
field on `project`, which only `server.value` sets, so `grove check`,
`grove claim` and `grove instructions` compute as before. It remembers only
answers that succeeded: a `planned` that is no commit costs one failed
`rev-parse` each time it is asked, at most 3 per cycle.
Check: the project value now passes the documents it read with the
target's schemas to the cycle check, where before it re-read the target
with the working tree's schemas, as the proposal's decisions say.
Also changed: `internal/cli/instructions.go` — `rewritten` now takes what
`atBases` read, so `grove instructions implement` reads its one change's
base through it.
Also changed: `internal/cli/project.go` — the `plans` field that lets
`changedSincePlan` run once per cycle in one project value.
Also changed: `internal/gitx/tree.go` — `TreeFS` reads its blobs through
`Objects` rather than its own `cat-file --batch` parser.

## Scenarios

- None: the change has no change specs. The owner decided the process
  bound stays observed behavior, guarded by `TestGitProcesses`, which
  passes.

## Reconciliation

- None: the change has no change specs.
