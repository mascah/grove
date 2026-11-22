# Changes to propose

<!-- The changes that change 261122-m78w6 proposes to the requirements of
`propose`. This is future work: nothing here describes the target branch
until the change is implemented and the maintained spec is reconciled. -->

Maintained spec: [propose](../../../../specs/propose/spec.md).

## MODIFIED Requirements

### Requirement: Ask a round before writing

When the owner asks for a new change, the propose instructions shall have
the agent send, before writing any document, one message that states the
scope, offers a split and asks what must land first where either applies,
and asks at most five of the decisions it would otherwise assume, the most
consequential, numbered, each with a recommended answer and its reason
worded so that "yes" accepts it, those that change behavior, scope or
compatibility or lose data first, then those that fix architecture or
interfaces costly to reverse; and wait for the owner. Each answer is
recorded as a decision in the proposal, quoted; a decision the owner
declines to make is recorded as an open question naming the task it
blocks; and every other decision the agent would otherwise assume,
routine choices among them, is recorded as an open question with its
recommended answer, its reason and the task it decides.

#### Scenario: The owner accepts the round

- GIVEN the agent would otherwise assume whether an unknown ID exits 1 or
  2, and where the check goes
- WHEN it asks both with recommended answers and the owner says "yes"
- THEN the proposal records each as a decision quoting the owner's "yes"
  and the recommended answer, and no document lists an assumption

#### Scenario: The round is ordered, and routine choices are not asked

- GIVEN the owner asks for "Editing an unknown ID prints an error", and
  the agent would otherwise assume the exit status, whether `find` in
  `ledger/records.py` raises `ValueError` or `ledger/cli.py` checks the
  ID, and the test's name and place in `tests/test_cli.py`
- WHEN it sends the round
- THEN the message states the scope, asks the exit status first and where
  the check goes second, and does not ask the test's name and place; the
  proposal lists that as an open question with its recommended answer,
  naming T-01

#### Scenario: More than five decisions

- GIVEN the agent would otherwise assume seven decisions that change
  behavior or are costly to reverse
- WHEN it sends the round
- THEN it asks the five most consequential, and the proposal lists the
  other two as open questions, each with its recommended answer, its
  reason and the task it decides

#### Scenario: The owner declines one

- GIVEN the round asks where the check goes
- WHEN the owner says to leave that one open
- THEN the proposal lists it under open questions, naming the task it
  decides, and `/grove-implement` builds its best answer to it with a
  `Check:` line, or blocks the task where a wrong answer would be costly

#### Scenario: Nothing to assume

- GIVEN the conversation settled every decision the documents need
- WHEN the agent reaches the round
- THEN its message states the scope and says there is nothing to ask, and
  it writes the documents once the owner confirms the scope
