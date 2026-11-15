# Changes to schemas

The changes that change 261115-z5kj4 proposes to the requirements of
[`schemas`](../../../../specs/schemas/spec.md). This is future work:
nothing here describes the target branch until the change is implemented
and the maintained spec is reconciled.

## MODIFIED Requirements

### Requirement: Explore and propose suggest bugfix

When the agent runs `grove instructions explore`, or `grove instructions
propose` for a new change with a schema other than `bugfix`, the system
shall print instructions to suggest `bugfix` when the owner describes
behavior that is wrong, and, once the owner confirms, to run `grove
instructions propose --schema bugfix` and follow what it prints; with
`--unattended`, the instructions shall instead have the agent keep the
schema they follow and name `bugfix` as an open question.

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

#### Scenario: A bug reported with no one in the session

- GIVEN the same fixture
- WHEN the agent runs `grove instructions propose --unattended "Editing a
  missing record crashes"` and writes the change
- THEN the change's schema is `default`, and an open question says
  `bugfix` fits behavior that is wrong
