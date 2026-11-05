## T-01 — Implement edits a file a task needs and says why

Result: the implement instructions, with the owner and unattended, the
`implementation.md` and tasks guides, `grove review`, the docs and the
walkthrough say what the change spec's "Change only the files the tasks
name" and "Review names what the branch approves and builds" say, as the
proposal's decisions give them.

Files:
- `internal/instruct/requires/implement-owner.md`
- `internal/instruct/requires/implement-unattended.md`
- `internal/instruct/schemas/default/documents/implementation/guide.md`
- `internal/instruct/schemas/default/documents/tasks/guide.md`
- `internal/cli/review.go`
- `internal/cli/walkthrough_test.go`
- `docs/coordinating.md`
- `docs/getting-started.md`
- `docs/reference.md`
- `dev/walkthrough.md`

Interfaces: `grove review` labels a file an `Also changed:` line names
`also changed <change-id>:T-NN: <why>`.

Constraints: rule 12's instruction-text budget; rule 3 as amended on
2026-11-05.

Acceptance:
- `grove instructions implement <change-id>`, with and without
  `--unattended`, says to edit a file a task needs or makes false
  without asking and to write `` Also changed: `<path>` — <why> `` under
  the task. It no longer says to ask the owner first, nor, unattended,
  that such a file stays as it is, to ignore `Delegated:` lines, or to
  write no `Also changed:` line. The wording test checks the new text and
  the absence of the old.
- The `implementation.md` guide says an `Also changed:` line records a
  file the implementing agent changed beyond the task's files, and why;
  the tasks guide no longer says implementing unattended ignores a
  `Delegated:` line.
- `grove review` shows a file an `Also changed:` line names as
  `also changed <change-id>:T-NN: <why>`; a test shows it, and that a
  file no line names is still `named nowhere`.
- `docs/coordinating.md`, `docs/getting-started.md` and
  `docs/reference.md` describe both modes and the label as the change
  spec gives them, with no "once you agree".
- The walkthrough step "Delegated with me, ignored unattended" is renamed
  and expects `grove run <cycle-id>` to create `tests/test_records.py`;
  the walkthrough's steps with the owner expect the agent to edit a file
  a task makes false without asking and write its `Also changed:` line,
  and its review step shows `also changed`.
- `just check` passes.

## T-02 — Trial case (c) expects the edit in both modes

Result: trial case (c)'s good runs, with the owner and unattended, are the
ones the proposal's decisions give, in its `good` function, the comment
above it, its scripted answers, and case (c)'s "A good run" and "Checked
by `good`" paragraphs in `dev/trials.md`, as 261029-wtbc0 left them at
6090cb3.

Files:
- `dev/trials/c.sh`
- `dev/trials.md`

Interfaces: none.

Constraints: rule 17; the "Edits outside Files/Also changed" count is
unchanged.

Acceptance:
- On case (c)'s setup with an end state written by hand on each branch,
  `good` passes when, on both `unknown-id-impl` and
  `unknown-id-unattended`, `ledger/records.py` is changed, T-01 in
  `implementation.md` is `Status: done` with an `Also changed:` line
  naming `ledger/records.py`, and `python3 -m unittest` passes; it prints
  a line naming the check when `ledger/records.py` is unchanged, when
  T-01 is blocked, or when no `Also changed:` line names the file.
- Case (c) no longer has a scripted answer about `ledger/records.py`, and
  `good` no longer checks that the file was unchanged before such an
  answer.
- Case (c) in `dev/trials.md` says, under "A good run" and "Checked by
  `good`", that in both modes the agent edits `ledger/records.py` without
  asking and writes an `Also changed:` line for it.
