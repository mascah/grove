<!-- Each task is a section headed "## T-NN — title". IDs are never
renumbered or reused. -->

## T-01 — The project value starts a bounded number of git processes

Result: `grove changes --json`, `grove serve`'s `/api/project` and `grove
check` start no more `git` processes for an added change or run record. An
added unmerged branch that holds no Grove document the target lacks adds
none either. The counts are those the proposal's Expected section allows.

Files:
- `internal/cli/walkthrough_test.go`
- `internal/cli/serve.go`
- `internal/cli/check.go`
- `internal/cli/claim.go`
- `internal/gitx/gitx.go`

Interfaces: `internal/gitx` runs Git with given standard input, for the
batched `diff-tree --stdin` and `cat-file --batch` calls.

Constraints:
- Every command prints what it printed at the change's base, except the
  rename case the proposal's Expected section names.
- Git runs only in `internal/gitx`, by argv.
- `just check` stays within the Go and test line budgets.

Acceptance:
- A walkthrough test, written first and failing at the change's base, puts
  a `git` shim on `PATH` that logs each call and runs the real Git. It
  runs `grove changes --json` once on a fixture project, then again after
  adding approved changes whose tasks name files and that modify a
  maintained requirement, run records, and an unmerged branch with no
  Grove document. Both runs start the same number of `git` processes.
- The existing walkthrough tests pass unchanged, except where they cover
  a rename inside a folder a task names.
- On a change whose base is a short commit ID, `grove check` still reports
  a named file changed since base.
- `grove check` in a checkout whose target branch does not exist, with no
  cycle and no task naming a file, succeeds as it does at the change's
  base.
- `just check` passes.
