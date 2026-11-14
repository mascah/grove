# Changes to propose

The changes that change 261114-zpgx9 proposes to the requirements of
[propose](../../../../specs/propose/spec.md), as at 999ab56.

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
  decides, and `/grove-implement` builds its best answer to it with a
  `Check:` line, or blocks the task where a wrong answer would be costly

#### Scenario: Nothing to assume

- GIVEN the conversation settled every decision the documents need
- WHEN the agent reaches the round
- THEN it says there is nothing to ask and writes the documents
