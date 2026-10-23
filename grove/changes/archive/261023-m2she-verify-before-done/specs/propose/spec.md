# Changes to propose

The maintained spec is [grove/specs/propose/spec.md](../../../../specs/propose/spec.md).

## ADDED Requirements

### Requirement: Check the change before presenting it

When the agent has written a change's documents, the propose instructions
shall have it check, before it tells the owner the change is written, that
each requirement has an outcome a user or caller can observe, each ADDED
or MODIFIED requirement is covered by a task's acceptance, the proposal's
scope, the change specs and the tasks do not contradict each other, and
no task does what the scope leaves out; fix what it can, tell the owner
what it changed, and ask the owner about what needs a decision.

#### Scenario: A requirement no task covers

- GIVEN the agent wrote a change whose change spec ADDs "Delete a record"
  and whose only task's acceptance covers editing
- WHEN it checks the change before presenting it
- THEN it adds the missing acceptance or task, and tells the owner it did

#### Scenario: Scope and a task disagree

- GIVEN the proposal puts the README out of scope and a task names
  `README.md`
- WHEN the agent checks the change
- THEN it asks the owner which holds, and changes nothing until they say
