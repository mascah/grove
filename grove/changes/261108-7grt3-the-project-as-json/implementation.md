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
