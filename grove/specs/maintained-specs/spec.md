# Maintained specs

## Purpose

A maintained spec is a folder under `specs/` holding `spec.md`, titled by
its first heading and described by its Purpose's first sentence. This
covers how Grove finds a maintained spec under the document root and what
it reads from it, its ID, title and description, and what `grove check`
reports about a spec folder. Owners writing specs, and agents reading the
project index, rely on it.

## Boundaries

This covers finding a spec and reading its title and description. What a
spec's requirements must say, and how archive folds them, is
[requirements](../requirements/spec.md). A change's `.grove.yaml` is not
covered here; it stays as it is.

## Requirements

### Requirement: A spec is a folder holding spec.md

When Grove reads the documents under the root, the system shall treat
each folder `specs/<id>/` that holds a `spec.md` as maintained spec
`<id>`, whose title is the text of the first line of `spec.md` that
starts with `# `, whether or not the folder holds a `.grove.yaml`.

#### Scenario: A spec with no .grove.yaml

- GIVEN the target holds `grove/specs/record-editing/spec.md`, whose
  first heading is `# Record editing`, and no `.grove.yaml` beside it
- WHEN the owner runs `grove specs`
- THEN it prints `record-editing  Record editing
  grove/specs/record-editing/spec.md`

#### Scenario: A folder with no spec.md

- GIVEN `grove/specs/notes/` holds only `ideas.md`
- WHEN the owner runs `grove specs` and `grove check`
- THEN `grove specs` lists no `notes`, and `grove check` reports nothing
  about it

### Requirement: A spec's description is its Purpose's first sentence

When the system lists a maintained spec in a skill's project index, it
shall give as the spec's description the first sentence of the first
paragraph under its `## Purpose` heading, its lines joined by spaces: the
text up to and including the first full stop that is followed by a space
or ends the paragraph; and with no such paragraph, no description.

#### Scenario: A Purpose of two sentences

- GIVEN `record-editing`'s Purpose reads "Editing a saved record keeps its
  category. A record is one line of the ledger." over two lines
- WHEN the agent runs `grove instructions explore`
- THEN the project index's line for `record-editing` gives its title and
  path, then "— Editing a saved record keeps its category."

#### Scenario: A full stop inside a name

- GIVEN a spec's Purpose starts "A project's `grove.yaml` gives each skill
  its context. It is optional."
- WHEN the agent runs `grove instructions explore`
- THEN the spec's description is "A project's `grove.yaml` gives each
  skill its context."

#### Scenario: A spec with no Purpose

- GIVEN `record-editing`'s `spec.md` has no `## Purpose` heading
- WHEN the agent runs `grove instructions explore` and the owner runs
  `grove check`
- THEN the index lists `record-editing` with its title and path and no
  description, and `grove check` reports nothing about it

### Requirement: Report a spec with no title

If a maintained spec's `spec.md` has no line that starts with `# `, then
`grove check` shall report that file, saying it needs a `# ` heading to
give the spec its title.

#### Scenario: No heading

- GIVEN `grove/specs/record-editing/spec.md` starts with `## Purpose`
  and has no `# ` line
- WHEN the owner runs `grove check`
- THEN it reports `grove/specs/record-editing/spec.md`, saying it has no
  `# ` heading for the spec's title, and exits non-zero

### Requirement: Report a leftover .grove.yaml

When the owner runs `grove check`, the system shall report each
`specs/<id>/.grove.yaml` under the root in one sentence, saying it is no
longer read and to delete it, and shall not read its fields.

#### Scenario: A spec folder from before

- GIVEN `grove/specs/record-editing/` holds `spec.md` and a `.grove.yaml`
  with `title: Record editing` and `owner: me`
- WHEN the owner runs `grove check`
- THEN it reports `grove/specs/record-editing/.grove.yaml` once, saying
  it is no longer read and to delete it, reports no unknown field in it,
  and exits non-zero

#### Scenario: Deleted

- GIVEN the owner deletes that `.grove.yaml`
- WHEN the owner runs `grove check`
- THEN it prints `No problems in grove/.`

### Requirement: A maintained spec has no metadata

When the agent runs `grove instructions archive <change-id>` for a change
with a change spec for a capability, the system shall give the `spec`
template with no `.grove.yaml` to write beside it; and if a project's
change schema gives a kind `metadata` other than
`changes/<change>/.grove.yaml`, then `grove check` shall report it.

#### Scenario: Archiving a new capability

- GIVEN finished change C has a change spec for `record-editing`, which
  has no maintained spec
- WHEN the agent runs `grove instructions archive C`
- THEN the instructions give the template for
  `grove/specs/<capability>/spec.md` and none for
  `grove/specs/<capability>/.grove.yaml`

#### Scenario: A schema with spec metadata

- GIVEN `grove/schemas/bugfix/schema.yaml` lists kind `spec` with
  `metadata: specs/<capability>/.grove.yaml`
- WHEN the owner runs `grove check`
- THEN it reports that schema, saying metadata must be
  `changes/<change>/.grove.yaml`

## Known discrepancies

None recorded.

## Open questions

None.
