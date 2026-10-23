# Changes to propose

The maintained spec is [grove/specs/propose/spec.md](../../../../specs/propose/spec.md).

## ADDED Requirements

### Requirement: Offer to split a change

When the agent confirms the scope of a new change whose idea is too big
for one change or has parts the owner could agree separately, the propose
instructions shall have it say so and ask; if the owner agrees, it writes
several changes on the same branch, each from its own
`grove instructions propose`, linked with `depends_on`.

#### Scenario: Two parts, one needing the other

- GIVEN the owner asks to propose `ledger delete ID` and an error instead
  of a traceback for an unknown ID, and `delete` finds its record the way
  `edit` does
- WHEN the agent confirms the scope
- THEN it says these are two changes the owner could agree separately,
  and, once the owner agrees, writes the unknown-ID change and the delete
  change on this branch, each from its own `grove instructions propose`,
  the delete change naming the other in `depends_on`

#### Scenario: The owner keeps it whole

- GIVEN the agent offered to split the same idea in two
- WHEN the owner says to keep it as one change
- THEN the agent writes one change

### Requirement: Delegate choices in a task

Where the owner agrees in the round that the implementer may make some
choices without asking, or change some paths beyond a task's `Files:`,
the propose instructions shall have the agent write them on an optional
`Delegated:` line of that task, and write no `Delegated:` line the owner
did not agree to; Grove shall read nothing on that line.

#### Scenario: The owner delegates test layout

- GIVEN the round recommends that the implementer may lay out tests and
  name private helpers as it likes, and may add `tests/test_records.py`,
  and the owner says "yes"
- WHEN the agent writes T-01
- THEN T-01 has a `Delegated:` line naming those choices and
  `tests/test_records.py`, and its `Files:` line does not name
  `tests/test_records.py`

#### Scenario: Nothing delegated

- GIVEN the round asked nothing about delegation, or the owner declined it
- WHEN the agent writes the tasks
- THEN no task has a `Delegated:` line

#### Scenario: Check ignores the line

- GIVEN a change whose T-01 has a `Delegated:` line right after its
  `Files:` line, and the target changed `tests/test_records.py` after the
  change's `base`
- WHEN the owner runs `grove check`
- THEN it reports nothing about T-01: its `Files:` line reads as before,
  and the delegated path is not watched

## MODIFIED Requirements

### Requirement: Ask a round before writing

When the owner has confirmed the scope of a new change, the propose
instructions shall have the agent list, before writing any document, every
decision it would otherwise assume, numbered, each with a recommended
answer worded so that "yes" accepts it, in this order: those that change
behavior, scope or compatibility or lose data first, then those that fix
architecture or interfaces costly to reverse, and routine choices inside
the agreed approach last, as recommended defaults on one line that a
single "yes" accepts; and wait for the owner. Each answer is recorded as a
decision in the proposal, quoted, and a decision the owner declines to
make is recorded as an open question naming the task it blocks.

#### Scenario: The owner accepts the round

- GIVEN the scope is confirmed and the agent would otherwise assume where
  the check goes and whether the README is updated
- WHEN it lists both with recommended answers and the owner says "yes"
- THEN the proposal records each as a decision quoting the owner's "yes"
  and the recommended answer, and no document lists an assumption

#### Scenario: The round is ordered

- GIVEN the scope of "Editing an unknown ID prints an error" is confirmed,
  and the agent would otherwise assume the exit status, whether `find` in
  `ledger/records.py` raises `ValueError` or `ledger/cli.py` checks the
  ID, and the test's name and place in `tests/test_cli.py`
- WHEN it asks the round
- THEN the exit status comes first, where the check goes second, and the
  test's name and place last, on one line of recommended defaults, and
  none is left out

#### Scenario: One "yes" to the routine line

- GIVEN the round's last line recommends three routine choices
- WHEN the owner says "yes" to the round
- THEN the proposal records each of the three as a decision quoting the
  owner's "yes"

#### Scenario: The owner declines one

- GIVEN the round lists where the check goes
- WHEN the owner says to leave that one open
- THEN the proposal lists it under open questions, naming the task it
  decides, and `/grove-apply` leaves that task blocked while
  `/grove-implement` asks

#### Scenario: Nothing to assume

- GIVEN the conversation settled every decision the documents need
- WHEN the agent reaches the round
- THEN it says there is nothing to ask and writes the documents
