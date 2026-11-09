# Implementation: The project as JSON

<!-- Progress on change 261108-7grt3, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — `grove changes --json` prints the project value

Status: done

- `internal/cli/changes.go`: with `--json`, `listChanges` builds the
  server's `value` (the one `/api/project` answers) and prints it indented.
  Every key printed before keeps its shape: `pageChange` embeds
  `doc.Change` and only adds `cycle`, `ready` and `files`.
- `internal/cli/serve.go`: `value` starts `changes` as `[]`, so a project
  with no changes still prints `"changes": []` as before, not `null`
  (found on a fresh fixture: the CI guide's `jq '.changes[]'` failed on
  `null`). The comment on `value` names both readers.
- `internal/cli/walkthrough_test.go`: `TestServe`, after the attempts are
  recorded on origin, compares `grove changes --json` with `/api/project`
  as JSON: equal once `read` is removed, and the working tree's row
  `date`, which is also the time of the read (`time.Now()` in
  `branches`). The existing `TestCoordinating` check of `changes --json`
  (paths, states, `waits_on`, task status) passes unchanged on the new
  output.
- The guide "Run Grove in CI"'s `jq -r '.changes[] | [.id, .state, .path,
  .title] | @tsv'` run on a fresh fixture (no changes: prints nothing,
  exit 0) and on this repository (prints its changes).
- `just check`: passed; budgets 4335 Go, 1991 test, 408 page, 8613
  instruction words.

### T-02 — Documentation and the walkthrough

Status: done

- `docs/coordinating.md`, "The project as data": says the command prints
  the project as one object and that `grove serve`'s page reads the same
  one; the example gains `read`, `specs`, each change's `cycle`, `ready`
  and `files`, and the five keys after `changes`; one line each for
  `read`/`specs`, `cycle`/`ready`/`files`, `cycles`, `branches`, `runs`,
  `history` and `problems`, and that an empty list is `null` but for
  `changes`, which is `[]`.
- `docs/watching.md`: the intro names `grove changes --json` as printing
  the object the page reads. `docs/reference.md`: the `grove changes
  --json` row says it prints the project value.
- `dev/walkthrough.md`: coordinating step 6 names `ready` and the new
  keys; watching step 3 compares `grove changes --json` with
  `/api/project` and names K's verdict, B's `cycle`, the attempt in
  `runs` and `build-b` in `branches`.
- Run on a fresh fixture (`just fixture`, `grove init`, a bare `origin`):
  changes A and B (`depends_on: [A]`) in cycle K marked `unattended`,
  pushed; `grove claim` in a second clone; `git fetch`; `build-b` with
  B's T-01 `Status: done`. `grove changes --json` and `/api/project` from
  `grove serve --addr 127.0.0.1:7391`, compared with `jq -S 'del(.read)'`
  and `diff`: equal. K's verdict "claimed by attempt …", B's `cycle` K,
  `runs` the attempt with its host and branch, `branches` `build-b` with
  T-01 done. The guide "Run Grove in CI"'s `jq` printed both changes.
  Then A archived with T-01 done and C archived with `Status: dropped`:
  A `archived` under `changes/archive/` with task `done`, C `dropped`,
  B with no `waits_on`; before that, A had no `description` key and its
  task no `status`.
- `just check`: passed; budgets unchanged (4335 Go, 1991 test, 408
  page, 8613 instruction words).

## Scenarios

- changes "One approved, one archived": `internal/cli/walkthrough_test.go`
  TestCoordinating (`grove changes --json` block: states, `waits_on`
  absent, task `done`), passed; the archived `path` under
  `changes/archive/` shown on the fixture (T-02), passed.
- changes "A dropped change" (Print the project as data): shown on the
  fixture (T-02), C's `state` `dropped`, passed; no test asserts it.
- changes "No implementation yet": shown on the fixture (T-02), no task
  `status` and no `description` key, passed; no test asserts it.
- changes "Cycles and attempts": `internal/cli/walkthrough_test.go`
  TestServe (cycle K's members and claimed verdict, `runs` with host
  h2.example, A's `cycle` K; `grove changes --json` equal to the page),
  passed; B's `cycle` K shown on the fixture (T-02), passed.
- changes "A branch with work": `internal/cli/walkthrough_test.go`
  TestServe (branches, `build-b` with T-01 done; `grove changes --json`
  equal to the page), passed.
- changes "The same value as the page": `internal/cli/walkthrough_test.go`
  TestServe (`grove changes --json` against `/api/project`), passed. The
  objects are equal once `read` and the working tree's row `date` are
  left out: that `date` is also the time of the read (`time.Now()`), so
  the two can differ by a second there.
- changes "A task not done": `internal/cli/walkthrough_test.go`
  TestArchive, passed.
- changes "Every task done": `internal/cli/walkthrough_test.go`
  TestArchive, passed.
- changes "A dropped change" (Archive only a finished or dropped change):
  `internal/cli/walkthrough_test.go` TestDropChange, passed.
- serve "An approved change in no cycle": `internal/cli/walkthrough_test.go`
  TestServe (E with no cycle and ready), passed.
- serve "The same value as the command": `internal/cli/walkthrough_test.go`
  TestServe, passed, as for "The same value as the page".

## Reconciliation

- changes: applied to specs/changes/spec.md
- serve: applied to specs/serve/spec.md
