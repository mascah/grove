## T-01 — Unattended edits a file a task needs

Result: the unattended implement instructions, the tasks guide, the docs
and the walkthrough say what the change spec's "Change only the files the
tasks name" says, as the proposal's decisions give it.

Files:
- `internal/instruct/requires/implement-unattended.md`
- `internal/instruct/schemas/default/documents/tasks/guide.md`
- `internal/cli/walkthrough_test.go`
- `docs/coordinating.md`
- `docs/getting-started.md`
- `dev/walkthrough.md`

Interfaces: none.

Constraints: rule 12's instruction-text budget; rule 3 unchanged, so no
`Also changed:` line unattended.

Acceptance:
- `grove instructions implement <change-id> --unattended` no longer says
  a file no task's `Files:` line names stays as it is, nor to ignore
  `Delegated:` lines; it says to edit a file a task needs or makes false,
  name it and why in the task's evidence, and write no `Also changed:`
  line. The wording test checks the new text and the absence of the old.
- The tasks guide no longer says implementing unattended ignores a
  `Delegated:` line.
- `docs/coordinating.md` and `docs/getting-started.md` describe the
  unattended behavior as the change spec gives it, and say `grove review`
  shows such a file as `named nowhere`.
- The walkthrough step "Delegated with me, ignored unattended" is renamed
  and expects `grove run <cycle-id>` to create `tests/test_records.py`,
  with no `Also changed:` line in `implementation.md`.
- `just check` passes.

## T-02 — Trial case (c) unattended expects the edit

Result: trial case (c)'s unattended good run is the one the proposal's
decisions give, in its `good` function and in `dev/trials.md`.

Files:
- `dev/trials/c.sh`
- `dev/trials.md`

Interfaces: none.

Constraints: rule 17; the "Edits outside" count is unchanged.

Acceptance:
- On case (c)'s setup with an unattended end state written by hand,
  `good` passes when `ledger/records.py` is changed, T-01 in
  `implementation.md` is `Status: done` and names `ledger/records.py`, and
  `python3 -m unittest` passes on `unknown-id-unattended`; it prints a line
  naming the check when `ledger/records.py` is unchanged, when T-01 is
  blocked, or when T-01 does not name the file.
- Case (c) in `dev/trials.md` says, under "A good run", that unattended it
  edits `ledger/records.py` and names it in T-01's evidence.
