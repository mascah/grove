## T-01 — Read a spec from its spec.md

Result: Grove finds, titles, describes and checks maintained specs as
[specs/maintained-specs/spec.md](specs/maintained-specs/spec.md) says,
and no schema or guide gives a maintained spec a `.grove.yaml`.

Files: `internal/doc/doc.go`, `internal/doc/meta.go`,
`internal/instruct/instruct.go`, `internal/instruct/metadata/spec.grove.yaml`,
`internal/instruct/schemas/default/schema.yaml`,
`internal/instruct/schemas/default/documents/spec/guide.md`,
`internal/instruct/schemas/default/documents/spec/template.md`,
`internal/instruct/schemas/default/skills/archive.md`,
`internal/cli/walkthrough_test.go`

Interfaces: `doc.Load` finds specs by `specs/*/spec.md` and fills a
spec's `Meta.Title` from its first `# ` heading and `Meta.Description`
from its Purpose, so `grove specs`, the reader and the project index read
them as today; a spec folder's `.grove.yaml` is a `Diagnostic`.
`internal/instruct/metadata/spec.grove.yaml` is deleted, and the `spec`
kind has no `metadata`.

Acceptance:
- In `walkthrough_test.go`, on the fixture: the archived change's new
  spec is written from the `spec.md` template alone, `grove check` is
  clean, and `grove specs` and the explore index list it with the title
  from its heading and the first sentence of its Purpose, a full stop
  inside a name such as `grove.yaml` not ending it; with its Purpose
  deleted, the index gives no description and `grove check` stays clean.
- In `walkthrough_test.go`: a folder under `grove/specs/` with no
  `spec.md` is neither listed nor reported; a spec folder holding a `.grove.yaml` with an
  unknown field is reported once as no longer read, and no unknown field
  is reported; a `spec.md` with no `# ` line is reported.
- `grove instructions archive` gives no `.grove.yaml` template for a
  spec, and the spec guide says the first heading is the title and the
  Purpose's first sentence the description.
- A project schema whose `spec` kind keeps `metadata:
  specs/<capability>/.grove.yaml` is reported by `grove check`.

## T-02 — No design.md

Result: the propose instructions put every decision in the proposal, as
[specs/propose/spec.md](specs/propose/spec.md) "Decisions go in the
proposal" says.

Files: `internal/instruct/schemas/default/skills/propose.md`,
`internal/instruct/schemas/default/documents/proposal/template.md`

Interfaces: none.

Acceptance:
- `grove instructions propose "x"` on the ledger fixture contains no
  `design.md`.
- The proposal template's Decisions comment says a decision that needs
  room is written out there.

## T-03 — Grove's own specs move over

Result: Grove's seven maintained specs have no `.grove.yaml`, and each
Purpose's first sentence reads as its one-line description.

Files: `grove/specs/apply/.grove.yaml`, `grove/specs/config/.grove.yaml`,
`grove/specs/cycles/.grove.yaml`, `grove/specs/implement/.grove.yaml`,
`grove/specs/propose/.grove.yaml`, `grove/specs/prototype/.grove.yaml`,
`grove/specs/runs/.grove.yaml`, `grove/specs/config/spec.md`,
`grove/specs/runs/spec.md`

Interfaces: none.

Constraints: only the Purpose sections change; no requirement, scenario
or other section of a maintained spec is edited.

Acceptance:
- Every `.grove.yaml` under `grove/specs/` is deleted, and `grove check`
  on Grove's own repository prints `No problems in grove/.`
- In `grove instructions explore` on Grove's repository, each spec's
  index line gives a description that says what its `.grove.yaml`
  description said at 0c6f419 in substance: `config`'s opens with what
  `context` does, not "A project's `grove.yaml` configures Grove for
  it.", and `runs`' is one sentence short enough to read as a line.

## T-04 — Document specs without .grove.yaml

Result: the documentation says a maintained spec is a folder holding
`spec.md`, titled by its heading and described by its Purpose, and the
walkthrough writes no spec `.grove.yaml`.

Files: `docs/reference.md`, `docs/concepts.md`, `dev/walkthrough.md`

Interfaces: none.

Acceptance:
- `docs/reference.md` "Files" shows no `.grove.yaml` under
  `specs/<capability>/`, says where a spec's title and description come
  from and that `grove check` reports a leftover one; "Change schemas"
  says only a change's main document has `metadata` (rule 8).
- `docs/concepts.md`, where it says what a maintained spec is, says its
  title is its heading and its Purpose's first sentence describes it.
- `dev/walkthrough.md` step 5 writes the new spec with no `.grove.yaml`
  beside it, and a step shows `grove check` reporting a leftover one
  until it is deleted.
