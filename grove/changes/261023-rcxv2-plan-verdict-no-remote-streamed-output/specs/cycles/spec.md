# Changes to cycles

The changes that change 261023-rcxv2 proposes to the requirements of
`cycles` (`grove/specs/cycles/spec.md` at 5383ad8).

## ADDED Requirements

### Requirement: Say whether the cycle can run unattended

When the cycle file is written, the plan instructions shall have the agent
tell the owner whether `grove run` would build the cycle once it is marked
`unattended` and, for each obstacle, what removes it.

#### Scenario: Nothing in the way

- GIVEN members A and B, neither of which the owner said needs a person,
  and neither waits on a change outside the cycle
- WHEN the agent finishes the plan
- THEN it says `grove run` would build the cycle once `unattended: true`
  is merged in, and how that mark is added

#### Scenario: A member needs a person

- GIVEN the owner said member C needs a person, for judging how a new
  screen looks
- WHEN the agent finishes the plan
- THEN it says `grove run` would skip the cycle if marked, and that
  revising C with `/grove-propose C` so that its acceptance can be checked
  without a person, merging it, and re-planning the cycle removes that

#### Scenario: A member waits on a change outside the cycle

- GIVEN member B depends on change D, neither in the cycle nor archived as
  finished
- WHEN the agent finishes the plan
- THEN it says `grove run` would not build the cycle until D is archived as
  finished

#### Scenario: The runner needs something set up

- GIVEN member B has a task that calls a service needing a login, and the
  owner said the runner has it
- WHEN the agent finishes the plan
- THEN `needs_human` does not name B, and the verdict names the login the
  runner's checkout must have
