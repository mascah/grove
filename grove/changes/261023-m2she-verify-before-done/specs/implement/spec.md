# Changes to implement

The maintained spec is [grove/specs/implement/spec.md](../../../../specs/implement/spec.md).

## ADDED Requirements

### Requirement: Verify before finishing

When the agent has built every task it can of a change with the owner,
the implement instructions shall have it, before it finishes, run the
project's tests, re-read each done task's evidence against the code as it
now is and set a task whose evidence no longer holds back to `Status: in
progress`, and write a `## Scenarios` section in `implementation.md` with
one line per scenario of the change specs naming the check that shows it
and its result, or `unverified` and why; its finishing report names each
unverified scenario.

#### Scenario: Every scenario has a test

- GIVEN a change whose change spec has two scenarios, both shown by tests
  that pass
- WHEN the agent finishes implementing it
- THEN `implementation.md`'s `## Scenarios` names each scenario with its
  test and "passed"

#### Scenario: A done task later broken

- GIVEN T-01 was marked done, and building T-02 changed the code T-01's
  evidence relies on so its test fails
- WHEN the agent rechecks before finishing
- THEN T-01 is `Status: in progress` again until it is rebuilt

#### Scenario: A scenario only a person can judge

- GIVEN a scenario says a list "reads clearly"
- WHEN the agent writes the scenario list
- THEN that line says `unverified` and why, and the finishing report
  names it

### Requirement: Archive names unverified scenarios

When the agent runs the archive instructions for a finished change, they
shall have it tell the owner each scenario `implementation.md` marks
`unverified`, before the fold, without refusing to archive.

#### Scenario: One scenario unverified

- GIVEN a finished change whose scenario list marks one scenario
  `unverified`
- WHEN the agent archives it
- THEN it names that scenario to the owner, then folds and archives
