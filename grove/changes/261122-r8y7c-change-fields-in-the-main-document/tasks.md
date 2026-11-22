<!-- Each task is a section headed "## T-NN — title". IDs are never
renumbered or reused. -->

## T-01 — Grove reads and writes a change's fields at the top of its main document

Result: Grove finds a change by its main document, reads its fields from
the lines between that document's two `---` lines, and reports a leftover
`.grove.yaml` and a folder with two main documents. Propose's template
puts the fields there. The reader and `grove serve` show them as text.
This is what the changes spec's ADDED requirements and the modified
requirements of maintained-specs, propose, schemas and serve say.

Files:
- `internal/doc/doc.go`
- `internal/doc/meta.go`
- `internal/doc/doc_test.go`
- `internal/instruct/instruct.go`
- `internal/instruct/metadata/`
- `internal/instruct/schemas/default/schema.yaml`
- `internal/instruct/schemas/default/documents/proposal/guide.md`
- `internal/instruct/schemas/bugfix/documents/proposal/guide.md`
- `internal/instruct/requires/propose.md`
- `internal/instruct/requires/implement.md`
- `internal/cli/check.go`
- `internal/cli/instructions.go`
- `internal/cli/specs.go`
- `internal/cli/serve.go`
- `internal/cli/walkthrough_test.go`

Delegated: test layout, private helpers, the name of the fields template
under `internal/instruct/metadata/`, and the exact wording of each
one-sentence message and of the instruction text.

Interfaces: in `internal/doc`, a change's main document, its `Path`, is
the file its fields are read from, and a problem in a field is reported
at that file and line. In `internal/instruct`, a kind's `metadata` takes
`frontmatter`, and the main kind's printed template starts with Grove's
fields between two `---` lines.

Constraints:
- The fields, their names, their meanings and what is required stay as
  they are. `Front` in `internal/doc/meta.go` splits the file, and
  `parse` reads the fields.
- Grove reads no change from a `.grove.yaml`, and has no fallback for one.
- `/api/project` and `grove changes --json` gain no key. `path` names the
  main document, and `files` no longer list `.grove.yaml`.
- The history of the record still marks approvals, revisions and archives
  in commits made before the fields moved.
- The report of a leftover `.grove.yaml` in a spec folder stays as it is.
- `just check` stays within the Go and word budgets.

Acceptance:
- `grove instructions propose "<request>"` gives the proposal's template
  starting with `---`, Grove's fields with `created` and `base` filled,
  and a line `---`, followed by the schema's template. It gives no
  `.grove.yaml` file, and the same holds with `--schema bugfix`.
- A change whose `proposal.md` starts with its fields is listed by
  `grove changes` and read by `grove check`, `grove claim`, `grove review`
  and `grove serve` as before. A problem in a field is reported at its
  line in `proposal.md`.
- `grove check` reports a change folder's leftover `.grove.yaml`, and a
  folder with two Markdown files starting with `---`, each in one
  sentence, and exits non-zero. It reports a schema whose `metadata` is
  not `frontmatter` or belongs to a kind outside `changes/<change>/`.
- `grove check` still reports a proposal whose first heading is not its
  `title` field.
- The implement and propose instructions say never to edit the fields
  above the second `---` by hand, where they said `.grove.yaml`, and the
  error for changing a change's schema names those fields.
- In the reader and on `grove serve`'s page, a proposal shows its fields
  as text above its rendered body.
- `/api/project`'s history marks the scenarios of "History of the record",
  the one from before the fields moved included.
- Tests show each scenario of the changes spec's ADDED requirements, and
  `just check` passes.

## T-02 — Convert Grove's own changes

Result: every change under `grove/changes/`, open and archived, this one
included, keeps its fields at the top of its `proposal.md` and has no
`.grove.yaml`.

Files:
- `grove/changes/`

Constraints:
- Only the fields move. Each field keeps its value and order, and each
  document's body is unchanged.
- Converted in one commit, apart from T-01's code, so the history marks
  it once.

Acceptance:
- `git ls-files 'grove/changes/**/.grove.yaml'` prints nothing.
- `grove check` built from T-01 reports no problems in `grove/`.
- `grove changes --json` built from T-01 gives every change the same `id`,
  `title`, `created`, `description`, `schema`, `depends_on`, `base` and
  `state` as `grove changes --json` at this change's `base` gives it,
  apart from changes added since.
- `git diff --stat` for the conversion commit touches only `proposal.md`
  and `.grove.yaml` files.

## T-03 — Docs, the walkthrough, the trials and the tracker guide

Result: no page, step or script mentions a change's `.grove.yaml`, except
to say how to convert one. The docs show a proposal's fields at its top
and give the prompt that converts a project.

Files:
- `docs/concepts.md`
- `docs/coordinating.md`
- `docs/reference.md`
- `docs/getting-started.md`
- `docs/watching.md`
- `docs/guides/run-grove-in-ci.md`
- `docs/guides/push-to-a-tracker.md`
- `dev/walkthrough.md`
- `dev/trials/`

Delegated: the exact wording, and where in each page the text goes.

Acceptance:
- `docs/concepts.md`'s table of a change's files has no `.grove.yaml`
  row, and its `proposal.md` row names the fields at its top.
- `docs/coordinating.md`'s `depends_on` example shows the fields at the
  top of `proposal.md`.
- `docs/reference.md` says a schema's main kind has
  `metadata: frontmatter`. It gives the prompt that converts a project's
  changes, open and archived: move each `.grove.yaml`'s lines to the top
  of the main document between two `---` lines, delete the file, and run
  `grove check`. It also says `grove check` reports a leftover one.
- `docs/guides/push-to-a-tracker.md` makes the issue body the proposal
  without its fields.
- `dev/walkthrough.md`'s steps that read or write a change's `.grove.yaml`
  read or write its proposal's fields instead. One step converts a change
  by hand after `grove check` reports its leftover `.grove.yaml`.
- `dev/trials/*.sh` set up changes with their fields in `proposal.md`, and
  `just trials` sets up each case without error.
- `grep -rn '\.grove\.yaml' docs dev README.md` finds only the conversion
  prompt and the report of a leftover file.
