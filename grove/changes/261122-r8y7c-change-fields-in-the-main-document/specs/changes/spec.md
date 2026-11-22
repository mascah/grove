# Changes to changes

The changes that change 261122-r8y7c proposes to the requirements of
[changes](../../../../specs/changes/spec.md), as it stands at 15816d7.
This is future work: nothing here describes the target branch until the
change is implemented and the maintained spec is reconciled.

## ADDED Requirements

### Requirement: A change's fields lead its main document

When Grove reads the documents under the root, the system shall treat
each folder `changes/<id>-<slug>/` or `changes/archive/<id>-<slug>/` that
holds exactly one Markdown file directly in it whose first line is `---`
as change `<id>`, whose main document is that file, and shall read the
change's fields, `title`, `created`, `description`, `schema`,
`depends_on` and `base`, from the YAML between that line and the next
line that is `---`, reporting a problem in them at its line in that file.

#### Scenario: A proposal with its fields

- GIVEN `grove/changes/261019-r2w7k-add-ledger-delete/proposal.md` starts
  with a line `---`, then `title: "Add ledger delete ID"`, `created`,
  `schema: default` and `base` lines, then a line `---` and
  `# Add ledger delete ID`, and the folder holds no `.grove.yaml`
- WHEN the target holds it and the owner runs `grove changes`
- THEN it prints `261019-r2w7k  Add ledger delete ID  approved`

#### Scenario: A main document of another name

- GIVEN a project schema `rfc` whose main kind is
  `changes/<change>/rfc.md` with `metadata: frontmatter`, and a change
  whose `rfc.md` starts with its fields
- WHEN the owner runs `grove changes --json`
- THEN that change's `path` is its `rfc.md`, and its `files` hold no
  `.grove.yaml`

#### Scenario: A field with a problem

- GIVEN change B's `proposal.md` has `depends_on: 260101-zzzzz`, not a
  list, on its fourth line
- WHEN the owner runs `grove check`
- THEN it reports B's `proposal.md` at line 4, saying `depends_on`
  expects a list of change IDs, and exits non-zero

### Requirement: Report a leftover change .grove.yaml

When the owner runs `grove check` and a file `changes/<folder>/.grove.yaml`
or `changes/archive/<folder>/.grove.yaml` exists under the root, the
system shall report that file in one sentence saying a change's fields
now lead its main document, read no change from it, and exit non-zero.

#### Scenario: A project not converted

- GIVEN `grove/changes/261019-r2w7k-add-ledger-delete/` holds
  `.grove.yaml` and a `proposal.md` that does not start with `---`
- WHEN the owner runs `grove check`
- THEN it reports that `.grove.yaml` once, `grove changes` lists no
  261019-r2w7k, and `grove check` exits non-zero

#### Scenario: Converted

- GIVEN the owner moves that file's fields to the top of `proposal.md`
  between two `---` lines and deletes `.grove.yaml`
- WHEN the owner runs `grove check`
- THEN it reports nothing about the change's fields, and `grove changes`
  lists 261019-r2w7k

### Requirement: Report a change with two main documents

When the owner runs `grove check` and a change folder holds more than one
Markdown file directly in it whose first line is `---`, the system shall
report the folder in one sentence naming those files, read no change from
it, and exit non-zero.

#### Scenario: Fields in two files

- GIVEN change C's folder holds `proposal.md` and `notes.md`, each
  starting with a line `---`
- WHEN the owner runs `grove check`
- THEN it reports C's folder, naming both files, and exits non-zero

#### Scenario: A task list with no fields

- GIVEN change C's `proposal.md` starts with its fields, and its
  `tasks.md` starts with an HTML comment
- WHEN the owner runs `grove check`
- THEN it reports nothing about C's main document
