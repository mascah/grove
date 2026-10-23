# Changes to apply

The maintained spec is [grove/specs/apply/spec.md](../../../../specs/apply/spec.md).

## ADDED Requirements

### Requirement: Verify before finishing a member

When the agent has built every task it can of a member of a cycle, the
apply instructions shall have it, before it finishes that member, run the
project's tests, re-read each done task's evidence against the code as it
now is and set a task whose evidence no longer holds back to `Status: in
progress`, and write a `## Scenarios` section in `implementation.md` with
one line per scenario of the change specs naming the check that shows it
and its result, or `unverified` and why; its report names each unverified
scenario, and it asks no one.

#### Scenario: A scenario no test shows

- GIVEN a member whose change spec has a scenario no test checks
- WHEN the agent finishes the member
- THEN that line says `unverified` and why, the report names it, and the
  agent goes on without asking
