# Changes to propose

The changes that change 261023-fca8x proposes to the requirements of
`propose`: the documents `/grove-propose` writes and how it settles what
they say. There is no maintained spec for `propose` at b143991.

## ADDED Requirements

### Requirement: Ask a round before writing

When the owner has confirmed the scope of a new change, the propose
instructions shall have the agent list, before writing any document, every
decision it would otherwise assume, numbered, each with a recommended
answer worded so that "yes" accepts it, and wait for the owner; each
answer is recorded as a decision in the proposal, quoted, and a decision
the owner declines to make is recorded as an open question naming the task
it blocks.

#### Scenario: The owner accepts the round

- GIVEN the scope is confirmed and the agent would otherwise assume where
  the check goes and whether the README is updated
- WHEN it lists both with recommended answers and the owner says "yes"
- THEN the proposal records each as a decision quoting the owner's "yes"
  and the recommended answer, and no document lists an assumption

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

## Observed current behavior

- At b143991, `internal/instruct/requires/propose.md` has "A question the
  owner defers is not answered by guess. It goes into the change's open
  questions"; the proposal guide and template direct inference to "a list
  headed 'Assumptions to confirm in review'"; the propose guide's step 2
  asks the owner to confirm the scope and asks "If an unanswered question
  would change the requirements or how a task is built".
- Change 261023-rcxv2's proposal, at b143991, carries seven assumptions
  under that heading.
