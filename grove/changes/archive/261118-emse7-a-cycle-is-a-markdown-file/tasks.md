<!-- Each task is a section headed "## T-NN — title". IDs are never
renumbered or reused. -->

## T-01 — Grove reads, writes and shows a cycle as Markdown

Result: a cycle is read from `<root>/cycles/<id>-<slug>.md`, its fields
between two lines of `---` and its notes below; plan writes it so, with
notes; the implement instructions say a member's documents hold over
them; `grove serve` shows both parts; and `grove check` reports a `.yaml`
left under `cycles/`.

Files:
- `internal/doc/doc.go`
- `internal/doc/meta.go`
- `internal/doc/doc_test.go`
- `internal/cli/serve.go`
- `internal/instruct/instruct.go`
- `internal/instruct/metadata/cycle.yaml`
- `internal/instruct/metadata/cycle.md`
- `internal/instruct/requires/plan.md`
- `internal/instruct/schemas/default/skills/plan.md`
- `internal/cli/walkthrough_test.go`

Delegated: test layout and private helpers.

Interfaces: `doc.Cycle.Notes` is the text below a cycle file's fields,
trimmed, empty when there is none. The cycle template is
`internal/instruct/metadata/cycle.md`, and `cycle.yaml` is removed.
`grove instructions plan` names `<root>/cycles/<id>-<slug>.md` as the file
to write.

Constraints:
- The seven fields, their checks and their messages stay as they are; the
  parser that reads them is the one that reads a `.grove.yaml`, given the
  text between the two lines of `---`.
- No cycle is read from a `.yaml` file, and nothing converts one
  (`CONSTITUTION.md` rule 13).
- The two new messages are, after the file's path: "a cycle file is
  Markdown now; rename it to `<name>.md` and put its fields between two
  lines of ---" and "a cycle file starts with its fields between two
  lines of ---" (`CONSTITUTION.md` rule 14).
- The sentence saying a member's documents hold over the notes is printed
  where the instructions print the cycle file, so the cycle's session and
  each member's read it from one place in `internal/instruct/instruct.go`.
- `/api/project` gains no key, and `internal/cli/page.html` does not
  change.
- `just check` prints at most 8,147 instruction words.

Acceptance:
- `grove check` reports nothing for a `.md` cycle file holding the seven
  fields between two lines of `---`, with or without notes below, and
  reports a problem in a field by its line in the file (cycles, "Report a
  malformed cycle").
- `grove check` reports a `.md` cycle file that does not start with its
  fields, and each `.yaml` file under `cycles/`, done or not, with the
  messages above, and exits non-zero; `grove changes` and `grove claim`
  show and claim nothing for a `.yaml` cycle (cycles, "Report a cycle file
  left as YAML").
- A `.md` file under `cycles/` whose name starts with no cycle ID is
  reported as a `.yaml` one was (cycles, "Report a cycle file named
  without an ID").
- A cycle renamed to `.md` with the two lines added, and nothing else
  changed, keeps its `planned`, is not stale, and is claimed as before
  (cycles, "Report a cycle file left as YAML", "Converted").
- `grove instructions plan "0.4 release"` names a `.md` path and prints a
  template with the fields between two lines of `---`, and its text says
  what the notes hold, that there may be none, and that they change no
  member; `grove instructions plan <cycle-id>` says to bring the notes up
  to date with the new `planned` (cycles, "Plan a cycle", "Revise a
  cycle").
- `grove instructions implement <cycle-id>` and `grove instructions
  implement --runner <member-id>` print the cycle file with its notes and
  the sentence that a member's documents hold over them, and print no
  such sentence for a cycle file with no notes (implement, "A cycle's
  notes are advice").
- `/api/doc` for a cycle file answers HTML with the fields as text and the
  notes rendered, and still answers a change's `.grove.yaml` as text
  (serve, "A document rendered").
- A direct test in `internal/doc` covers the split: fields and notes,
  fields only, no opening `---`, no closing `---`.
- `just check` passes.

## T-02 — Docs, walkthrough and Grove's own cycle files

Result: the documentation and the walkthrough describe the cycle file as
Markdown with notes and say what a project holding `.yaml` cycle files
must do, and Grove's own cycles are converted.

Files:
- `docs/planning.md`
- `docs/reference.md`
- `docs/watching.md`
- `dev/walkthrough.md`
- `grove/cycles/`

Constraints: Grove's own cycle files keep their fields as they are and
gain no notes; the archived changes keep their text.

Acceptance:
- `docs/planning.md` shows a cycle file as plan writes it, fields between
  two lines of `---` and notes below, says what the notes hold and that
  only plan writes them, and that re-planning brings them up to date.
- `docs/reference.md` gives the layout `cycles/<id>-<slug>.md`, the
  fields, the notes below them, and one paragraph for a project that still
  has `.yaml` cycle files: what `grove check` reports and how to convert
  one.
- `docs/watching.md` says the pane shows a cycle's fields as text and its
  notes rendered.
- `dev/walkthrough.md` has a step for each of: plan writing a `.md` cycle
  with notes; a `.yaml` cycle file left over; a cycle file that does not
  start with its fields; and the implement instructions for a cycle whose
  file holds notes. Its other steps name `.md` cycle files.
- No file under `docs/`, and no step of `dev/walkthrough.md`, names a
  cycle file ending `.yaml`, except where `docs/reference.md` and the
  walkthrough describe one left over.
- Every file in `grove/cycles/` ends `.md`, and `grove check` on this
  repository is clean.
