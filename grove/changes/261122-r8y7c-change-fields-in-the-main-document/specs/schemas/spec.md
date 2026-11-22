# Changes to schemas

The changes that change 261122-r8y7c proposes to the requirements of
[schemas](../../../../specs/schemas/spec.md), as it stands at 15816d7.
This is future work: nothing here describes the target branch until the
change is implemented and the maintained spec is reconciled.

## MODIFIED Requirements

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
- THEN the instructions follow change schema `bugfix`, and the proposal
  template's fields say `schema: bugfix`

#### Scenario: A change with schema bugfix

- GIVEN an open change whose proposal's fields say `schema: bugfix`, in a
  project with no `grove/schemas/`
- WHEN the owner runs `grove check`
- THEN it reports nothing about the change's schema

#### Scenario: A schema that does not exist

- GIVEN a project with its own schema `short`
- WHEN the agent runs `grove instructions propose --schema bugfx`
- THEN it prints `grove: schema "bugfx" is not a change schema; schemas
  are default, bugfix and short` and exits non-zero
