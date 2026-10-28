# Changes to prototype

The changes that change 261028-rxsfa proposes to the requirements of
[`prototype`](../../../../specs/prototype/spec.md).

## MODIFIED Requirements

### Requirement: Propose uses a prototype's findings

When the conversation holds a prototype's findings, the propose
instructions shall have the agent write them into the proposal so that it
stands without the prototype branch, stating what was tried, what was
observed and the verdict, naming the branch and its commit only as where
they came from, and link no file outside the change and the target.

#### Scenario: Proposing after a prototype

- GIVEN a prototype on `prototype/coaching` at commit `abc1234`, not pushed,
  found that retyping one highlight felt right and three felt tedious
- WHEN the owner proposes from the conversation
- THEN the proposal states what was tried, that finding and the verdict, and
  names `prototype/coaching` at `abc1234` as where it came from

#### Scenario: An experiment not committed

- GIVEN the owner chose not to commit a prototype's experiment
- WHEN the owner proposes from the conversation
- THEN the proposal states the findings and says the experiment was not
  committed
