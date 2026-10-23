# Changes to propose

The maintained spec is [grove/specs/propose/spec.md](../../../../specs/propose/spec.md).

## ADDED Requirements

### Requirement: A change that records existing behavior has no tasks

When a change only records behavior the code already has, the propose
instructions shall have the agent write it with no tasks, and `grove
check` and `grove archive` shall accept it with none.

#### Scenario: Recording what edit does

- GIVEN the owner asks for a change that records what `ledger edit` does
  today, and agrees each rule, with no code to change
- WHEN the agent writes the change
- THEN its change spec holds the rules under ADDED and what the code does
  under "Observed current behavior", its `tasks.md` holds no `T-NN`
  task, and `grove check` reports no problem in it

#### Scenario: Archived with no tasks

- GIVEN change R has no tasks, and its change spec is folded into
  `grove/specs/record-editing/spec.md`
- WHEN the owner runs `grove archive R`
- THEN Grove moves R to `changes/archive/`, since no task waits on
  `implementation.md`

#### Scenario: The owner wants the code changed too

- GIVEN, while recording what `ledger edit` does, the owner decides that
  `ledger edit 99 --note x` should print an error instead of a traceback
- WHEN the agent writes the change
- THEN fixing it is a task with its `Files:` line, or the owner leaves
  the traceback under "Known discrepancies" for a later change; the
  change is not written as one that only records behavior

## Observed current behavior

At 0c6f419:

- `grove archive` refuses a change only for a task in `tasks.md` that is
  not `Status: done` and loops over none when there are none
  (`internal/cli/instructions.go`, `archive`). Tried on the ledger
  fixture: a change with no `tasks.md`, and one whose `tasks.md` holds
  only a comment, each archived.
- `grove check` reads tasks only from `T-NN` headings and reports nothing
  for a change with none (`internal/doc/doc.go`, `tasks`). Tried on the
  ledger fixture with both changes above: "No problems in grove/."
- The tasks guide says nothing about a change with no tasks
  (`internal/instruct/schemas/default/documents/tasks/guide.md`).
