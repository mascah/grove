## T-01 — Files takes bullets and folders

Result: `grove check` and the implement instructions read `Files:` as
[specs/propose/spec.md](specs/propose/spec.md) "Tasks name files and
interfaces", "Report a malformed Files line" and "Report a named file that
moved", and [specs/implement/spec.md](specs/implement/spec.md) "Tell the
owner of a moved file", say; a folder entry counts in a cycle's staleness.

Files: `internal/doc/doc.go`, `internal/cli/check.go`,
`internal/instruct/schemas/default/documents/tasks/guide.md`,
`internal/instruct/schemas/default/documents/tasks/template.md`,
`docs/reference.md`, `dev/walkthrough.md`,
`internal/cli/walkthrough_test.go`

Interfaces: `doc.Task.Files` holds the paths of either form, a folder
keeping its trailing `/`; a named file matches a changed path when they
are equal or the named file is a folder above it, in `check`'s moved-file
report, implement's changed-on-target line and `changedSincePlan`.

Acceptance:
- On the fixture, a task with `Files:` followed by bullets naming
  `ledger/records.py` and `tests/fixtures/` passes `grove check`; the same
  task with a comma line still passes.
- `` `../ledger/cli.py` `` on a `Files:` line is reported with the
  message that names both forms and says a folder ends in `/`.
- A commit on the target touching a file under a named folder is reported
  by `grove check`, named in `grove instructions implement`, and makes a
  cycle listing the change stale.
- The tasks guide and template teach one path per bullet and folders; the
  docs and the walkthrough say the same (rule 8).
- `walkthrough_test.go` covers both forms, a folder, and the message.
- `just check` passes, the instruction text included.

## T-02 — Implement records agreed files on Also changed lines

Result: the implement instructions behave as
[specs/implement/spec.md](specs/implement/spec.md) "Revise the change on
its branch", "Tell the owner of a moved file" and "Ignore a task's
delegation unattended" say.

Files: `internal/instruct/requires/implement-owner.md`,
`internal/instruct/requires/implement-unattended.md`,
`internal/instruct/schemas/default/documents/implementation/guide.md`,
`docs/coordinating.md`, `docs/getting-started.md`, `dev/walkthrough.md`,
`internal/cli/walkthrough_test.go`

Interfaces: an `Also changed:` line under a task's heading in
`implementation.md`, `` Also changed: `path` — reason ``, one per file.

Acceptance:
- On the fixture, with the owner, an agreed file beyond `Files:` and
  `Delegated:` ends as an `Also changed:` line with the reason, and
  `tasks.md` is unchanged on the branch.
- A revision adding a task goes through `grove instructions propose
  <change-id>`.
- Unattended, a file no `Files:` line names is named in the task's
  evidence and no `Also changed:` line is written.
- The docs and the walkthrough say the same (rule 8).
- `walkthrough_test.go` checks that the owner's implement instructions
  name `Also changed:` and that the unattended ones forbid it.
- `just check` passes, the instruction text included.

## T-03 — Review shows planned against actual

Result: `grove review` behaves as
[specs/implement/spec.md](specs/implement/spec.md) "Review shows planned
against actual" says.

Files: `internal/cli/review.go`, `internal/doc/doc.go`,
`docs/reference.md`, `dev/walkthrough.md`,
`internal/cli/walkthrough_test.go`

Interfaces: `doc.Task` gains the paths and reasons of its `Also changed:`
lines, not printed in the project data; `grove review` prints, before the
diff, one line per file the branch changes outside the document root:
`planned <change-id>:T-NN`, `agreed <change-id>:T-NN: <reason>` or
`named nowhere`, then the path.

Constraints: rule 3 must name `Also changed:` among what Grove reads
before this task is built.

Acceptance:
- On the fixture, the branch of "Planned, agreed and named nowhere" shows
  the four files as that scenario says, before and after archiving, and
  `grove review` exits 0.
- A branch holding no change prints no such list.
- `docs/reference.md` and the walkthrough show the list (rule 8).
- `walkthrough_test.go` covers the three marks and a folder entry.
- `just check` passes.

## T-04 — Trials read bullets and Also changed

Result: the trials' "Edits outside Files/Delegated" count reads the
bullet form and does not count a path an `Also changed:` line names.

Files: `dev/trials/run.sh`, `dev/trials.md`

Acceptance:
- A trial run whose `tasks.md` lists files as bullets, and whose
  `implementation.md` names a file on an `Also changed:` line, counts
  neither as outside.
- `dev/trials.md` describes the count so.
