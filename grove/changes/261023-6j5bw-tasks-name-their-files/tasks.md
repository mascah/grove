## T-01 — Tasks name their files

Result: a task names its files and interfaces as
[specs/propose/spec.md](specs/propose/spec.md) "Tasks name files and
interfaces" says; propose writes them, implement and apply touch only the
files the tasks name, and the documentation specifies it.

Files: `internal/instruct/schemas/default/documents/tasks/guide.md`,
`internal/instruct/schemas/default/documents/tasks/template.md`,
`internal/instruct/schemas/default/skills/propose.md`,
`internal/instruct/schemas/default/skills/implement.md`,
`internal/instruct/schemas/default/skills/apply.md`,
`internal/instruct/requires/implement.md`,
`internal/instruct/requires/apply.md`, `docs/concepts.md`,
`docs/reference.md`, `docs/coordinating.md`, `dev/walkthrough.md`,
`internal/cli/walkthrough_test.go`

Interfaces: the `Files:` line, after `Result:`, listing repo-relative
paths in backticks, comma-separated; the `Interfaces:` line, free text.

Constraints: the propose guide's "no file-by-file plan" line goes, and
nothing asks for code, commands or steps; the bugfix example guide in
`docs/reference.md` follows.

Acceptance:
- `docs/concepts.md`'s table says a task names its files and interfaces,
  and `docs/reference.md`'s guide texts and bugfix example show the two
  lines; `dev/walkthrough.md` proposes a change whose task names a test
  file and the file under fix (rule 8).
- On the ledger fixture, in Claude Code: `/grove-propose` writes a task
  whose `Files:` line names the files decided in the conversation, and
  `/grove-implement` on it edits no other file without asking.
- `internal/cli/walkthrough_test.go` checks that `instructions propose`,
  `instructions implement` and `instructions apply` print the rule;
  `just check` passes.

## T-02 — Check reports a named file that moved

Result: `grove check`, `grove run` and `grove instructions implement`
behave as [specs/propose/spec.md](specs/propose/spec.md) "Report a named
file that moved", [specs/cycles/spec.md](specs/cycles/spec.md) and
[specs/implement/spec.md](specs/implement/spec.md) say, and the
documentation specifies it.

Files: `internal/doc/doc.go`, `internal/cli/check.go`,
`internal/cli/instructions.go`, `internal/cli/run.go`,
`internal/instruct/requires/implement.md`, `docs/coordinating.md`,
`docs/planning.md`, `docs/reference.md`, `dev/walkthrough.md`,
`internal/cli/walkthrough_test.go`

Interfaces: `doc.Task` gains `Files []string`, read from the task's
`Files:` line; `changedSincePlan` also diffs each open member's named
files; `check`'s report on a change names the task and file and says to
revise with propose, one sentence (rule 14); the cycle report keeps its
words.

Constraints: within the rule 12 budget; one Git call per change, with
paths; no new gate beyond the stale cycle (rule 10 as amended by #28).

Acceptance:
- `docs/coordinating.md`'s stale section shows the report for a moved
  file and says a member's moved file makes its cycle stale instead;
  `docs/planning.md` says a cycle is stale when a member's named file
  moved; `dev/walkthrough.md` runs a change whose named file the target
  changes, for a change alone and for a member of a marked cycle.
- `internal/cli/walkthrough_test.go` covers: a change in no cycle with a
  named file changed since `base` reported, the same change unreported
  when the file is unchanged or only the branch changed it, a member not
  reported against `base` while its cycle is reported stale, `grove run`
  skipping that cycle, `instructions implement` naming the task and file,
  and a `Files:` line Grove cannot read reported by `check`; `just check`
  passes.
