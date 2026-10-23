# Specs without .grove.yaml, and no design.md

## Problem

A maintained spec's `.grove.yaml` says nothing its `spec.md` does not.
At 0c6f419 it takes a `title`, which `grove check` requires to repeat
the spec's first heading, and an optional `description`
(`internal/doc/doc.go`, `Load`; `internal/doc/meta.go`, `heading`;
`internal/instruct/metadata/spec.grove.yaml`). Grove finds a spec by that
file (`specs/*/.grove.yaml`), so archive has the agent write it beside
every new spec (`internal/instruct/schemas/default/skills/archive.md`),
and a title edited in one place and not the other is a `check` problem.
The description is a second summary beside the spec's `## Purpose`, and
the two drift: Grove's own seven specs each keep one
(`grove/specs/*/.grove.yaml`).

The propose guide and the proposal template offer a `design.md` "only
when a decision needs more room than the proposal's decisions section"
(`internal/instruct/schemas/default/skills/propose.md`,
`internal/instruct/schemas/default/documents/proposal/template.md`). No
schema lists it as a kind, so nothing writes it from a template, no
instructions include it, and no change in `grove/changes/` has one.

## Intended outcome

A maintained spec is a folder under `specs/` holding `spec.md`. Its title
is its first `# ` heading, and its description, in the instructions'
project index, is the first sentence of its `## Purpose`. Archive writes
no `.grove.yaml` for a new spec, and `grove check` reports one left in a
spec folder as no longer read. Grove's own seven specs lose theirs, each
Purpose opening with a sentence that reads as its one-line description.
Propose and the proposal template no longer mention `design.md`: a
decision that needs room is written out in the proposal. Changes keep
their `.grove.yaml`. The rules are in
[specs/maintained-specs/spec.md](specs/maintained-specs/spec.md) and
[specs/propose/spec.md](specs/propose/spec.md).

## Scope

### In scope

- Finding a maintained spec by its `spec.md`, its title from its heading
  and its description from its Purpose (`internal/doc`).
- The leftover `.grove.yaml` report and the no-heading report in `grove
  check`.
- Dropping the spec's metadata template and the `spec` kind's `metadata`
  from the default schema, so a schema gives metadata only to a change's
  main document.
- The spec guide and template, the archive guide, the propose guide and
  the proposal template.
- Grove's own specs in `grove/specs/`.
- `docs/reference.md`, `docs/concepts.md`, `dev/walkthrough.md` and
  `internal/cli/walkthrough_test.go`.

### Out of scope

- A change's `.grove.yaml`, which holds `created`, `base` and
  `depends_on`; it stays as it is.
- `design.md` as a document kind; it is dropped, not added.
- Showing a description in `grove specs` or the reader, which show none
  at 0c6f419 (see Open questions).
- Any requirement of Grove's own specs: only their Purpose sections are
  reworded.

## Decisions

The owner agreed to the direction on 2026-10-23 ("yes to all") and to
every answer below in one round ("yes to C & D").

- **A maintained spec is a folder under `specs/` holding `spec.md`. Its
  title is the spec's first heading; its description is the first
  sentence of its `## Purpose` section.** Reason: the spec's `.grove.yaml`
  holds only a title that must repeat the heading, plus an optional
  description.
- **`grove check` reports a spec folder's leftover `.grove.yaml` as no
  longer read, in one sentence saying to delete it.** It is a new `check`
  problem on documents that were clean and removes a `.grove.yaml`, so
  the change is breaking (`CONTRIBUTING.md`, "Breaking changes"). Grove's
  own seven specs in `grove/specs/` move over in the same change: their
  `.grove.yaml` is deleted, and each Purpose's first sentence is made to
  read well as a one-line description, compared with today's.
- **`design.md` leaves the propose guide and the proposal template;
  decisions that need room go in the proposal. It does not become a
  document kind.**
- **Changes keep their `.grove.yaml`.** Reason: it holds `created`,
  `base` and `depends_on`, which no Markdown file carries.

## Open questions

- Can wait: decision 1 speaks of the description "in lists (`grove
  specs`, the reader)", but at 0c6f419 neither shows any description,
  a spec's or a change's (`internal/cli/specs.go`, `internal/tui/tui.go`);
  only the instructions' project index does. This change keeps the
  description where Grove shows it today, the index. If the owner wants
  `grove specs` and the reader to show it too, that is a new task here or
  a later change. It blocks no task.

## Related documents

- Requirements: [specs/maintained-specs/spec.md](specs/maintained-specs/spec.md),
  [specs/propose/spec.md](specs/propose/spec.md)
- Tasks: [tasks.md](tasks.md)
