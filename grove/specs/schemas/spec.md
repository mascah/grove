# Schemas

## Purpose

Grove has the change schemas `default` and `bugfix` in every project, a
project replaces any of their files, and explore and propose suggest
`bugfix` for wrong behavior. A change schema says which documents a change
holds and how each skill goes for it. `default` is the software change and
`bugfix` is for behavior that is wrong; a project replaces any of their
files, and adds schemas of its own, in `<root>/schemas/`. Owners choosing
what kind of change to propose, and agents following a schema's guides,
rely on this.

## Boundaries

This covers which schemas a project has, where a schema's files come from,
what `bugfix` proposes, and explore and propose suggesting it. What Grove
requires of every skill in every schema, and the document conventions
every schema keeps, are fixed (`CONSTITUTION.md` rules 3 and 10); building
and archiving a change belong to `grove/specs/implement/spec.md`.

## Requirements

### Requirement: Grove has a built-in bugfix schema

The system shall have the change schemas `default` and `bugfix` in every
project, listed in that order before the project's own, whether or not the
project has a `grove/schemas/bugfix/` folder.

#### Scenario: A project with no schemas of its own

- GIVEN the ledger fixture after `grove init`, with no `grove/schemas/`
- WHEN the agent runs `grove instructions propose "Editing a missing record crashes"`
- THEN the instructions list, under "Change schemas", `default` with its
  description and then `bugfix` with "A fix for behavior that is wrong,
  whether or not a maintained spec describes it."

#### Scenario: Proposing with bugfix

- GIVEN the same fixture
- WHEN the agent runs `grove instructions propose --schema bugfix
  "Editing a missing record crashes"`
- THEN the instructions follow change schema `bugfix`, and the
  `.grove.yaml` template says `schema: bugfix`

#### Scenario: A change with schema bugfix

- GIVEN an open change whose `.grove.yaml` says `schema: bugfix`, in a
  project with no `grove/schemas/`
- WHEN the owner runs `grove check`
- THEN it reports nothing about the change's schema

#### Scenario: A schema that does not exist

- GIVEN a project with its own schema `short`
- WHEN the agent runs `grove instructions propose --schema bugfx`
- THEN it prints `grove: schema "bugfx" is not a change schema; schemas
  are default, bugfix and short` and exits non-zero

### Requirement: Bugfix proposes a bug

When the agent runs `grove instructions propose --schema bugfix`, the
system shall print bugfix's propose guide and its proposal template and
guide: a proposal with the sections Observed, Expected, Reproduction, The
requirement it breaks and Open questions; a first task that adds a test
that fails before the fix; and a change spec only when the owner agrees
behavior that no maintained requirement gives, none when the fix restores
a maintained requirement.

#### Scenario: The bugfix proposal template

- GIVEN the ledger fixture with no `grove/schemas/`
- WHEN the agent runs `grove instructions propose --schema bugfix "x"`
- THEN the proposal template has the headings `## Observed`,
  `## Expected`, `## Reproduction`, `## The requirement it breaks` and
  `## Open questions`, and the change-spec and tasks templates are
  `default`'s

#### Scenario: A fix that restores a requirement

- GIVEN the same instructions
- WHEN the agent reads the propose guide
- THEN it says to cite the maintained requirement the bug breaks and write
  no change spec, and that the first task adds a test that fails before
  the fix

#### Scenario: Behavior no requirement gives

- GIVEN the same instructions
- WHEN the agent reads the propose guide
- THEN it says to ask the owner, in the round before writing, whether the
  corrected behavior should become agreed, and to write a change spec that
  ADDs it only if the owner says yes

### Requirement: Where a bugfix file comes from

When Grove reads a file of the `bugfix` schema, the system shall take it
from the project's `grove/schemas/bugfix/`, else from Grove's built-in
`bugfix`, else from the project's `grove/schemas/default/`, else from
Grove's `default`; and it shall take each field of `bugfix`'s
`schema.yaml` from the first of those places that gives it.

#### Scenario: The project's own default file

- GIVEN `grove/schemas/default/documents/tasks/template.md` with a
  `Rollback:` line, and no `grove/schemas/bugfix/`
- WHEN the agent runs `grove instructions propose --schema bugfix "x"`
- THEN the tasks template is the project's, with its `Rollback:` line

#### Scenario: The project's own bugfix file

- GIVEN `grove/schemas/bugfix/documents/proposal/template.md` holding
  `# <title>` and `## Seen`
- WHEN the agent runs `grove instructions propose --schema bugfix "x"`
- THEN the proposal template is the project's, and the propose guide is
  Grove's built-in `bugfix`'s

#### Scenario: A project that copied the example

- GIVEN `grove/schemas/bugfix/` written from `docs/reference.md`'s example
  at 57af2e3
- WHEN the owner runs `grove check`, and the agent runs `grove
  instructions propose --schema bugfix "x"`
- THEN `grove check` reports nothing about the schema, and the
  instructions use the project's files

#### Scenario: A mistake in the project's bugfix

- GIVEN `grove/schemas/bugfix/schema.yaml` holding `instrution: x` on its
  second line, and a file at `grove/schemas/bugfix/documents/propsal/template.md`
- WHEN the owner runs `grove check`
- THEN it reports `grove/schemas/bugfix/schema.yaml:2: unknown field
  "instrution"; schema.yaml takes description, documents and skills` and
  `grove/schemas/bugfix/documents/propsal/template.md: schema bugfix does
  not use this file`, and exits non-zero

### Requirement: Explore and propose suggest bugfix

When the agent runs `grove instructions explore`, or `grove instructions
propose` for a new change with a schema other than `bugfix`, the system
shall print instructions to suggest `bugfix` when the owner describes
behavior that is wrong, and, once the owner confirms, to run `grove
instructions propose --schema bugfix` and follow what it prints.

#### Scenario: Proposing a bug fix

- GIVEN the ledger fixture with no `grove/schemas/`
- WHEN the agent runs `grove instructions propose "Editing a missing
  record crashes"`
- THEN the instructions say to suggest `bugfix` for behavior that is
  wrong and to switch to it only once the owner confirms

#### Scenario: Exploring a bug

- GIVEN the same fixture
- WHEN the agent runs `grove instructions explore`
- THEN the instructions list `default` and `bugfix` under "Change
  schemas" and say to suggest `bugfix` for behavior that is wrong, adding
  `--schema bugfix` to `instructions propose` once the owner confirms

## Known discrepancies

None recorded.

## Open questions

None.
