# Changes to cycles

The changes that change 261118-5mqc1 proposes to the requirements of
[cycles](../../../../specs/cycles/spec.md), as it stands at dc8d335. This
is future work: nothing here describes the target branch until the change
is implemented and the maintained spec is reconciled.

## MODIFIED Requirements

### Requirement: Say whether the cycle can run unattended

When the cycle file is written, the plan instructions shall have the agent
tell the owner that `/grove-implement <cycle-id>` builds the cycle with
them, then whether `grove claim <cycle-id>` would claim it, then, where
the cycle file holds `unattended: true`, that merging it hands the cycle
to a `grove claim` given no cycle, such as a timer's, and otherwise that
`unattended: true`, merged on its own, does; and, for each obstacle, what
removes it.

#### Scenario: Nothing in the way

- GIVEN members A and B, neither of which waits on a change outside the
  cycle, and the owner did not ask for an unattended cycle
- WHEN the agent finishes the plan
- THEN it names `/grove-implement <cycle-id>` first, says
  `grove claim <cycle-id>` would claim the cycle, and says a `grove claim`
  given no cycle would claim it once `unattended: true` is merged in, and
  how that mark is added

#### Scenario: Marked as the owner asked

- GIVEN the owner asked for an unattended cycle, and the agent wrote
  `unattended: true` in the cycle file
- WHEN the agent finishes the plan
- THEN it names `/grove-implement <cycle-id>` first, says
  `grove claim <cycle-id>` would claim the cycle, and says merging the
  cycle file hands it to a `grove claim` given no cycle

#### Scenario: A member waits on a change outside the cycle

- GIVEN member B depends on change D, neither in the cycle nor archived as
  finished
- WHEN the agent finishes the plan
- THEN it says `grove claim` would not claim the cycle until D is archived
  as finished

#### Scenario: The cycle waits on a cycle

- GIVEN the owner said the new cycle follows cycle L, which is not done
- WHEN the agent finishes the plan
- THEN it says `grove claim` would not claim the cycle until L is done,
  and that the merges archiving L's changes leave its plan as it checked
  it

#### Scenario: The runner needs something set up

- GIVEN member B has a task that calls a service needing a login, and the
  owner said the runner has it
- WHEN the agent finishes the plan
- THEN the verdict names the login the runner's checkout must have

### Requirement: Mark a cycle unattended

The system shall treat a cycle as handed to `grove claim` only while the
target's cycle file holds `unattended: true`, and the plan instructions
shall have the agent write or remove that field only where the owner asks
for it, in the description or the conversation, never offering it, and
keep it as it is otherwise.

#### Scenario: Planned but not marked

- GIVEN the target holds cycle K without `unattended`
- WHEN `grove claim` is invoked
- THEN K is not claimed

#### Scenario: The owner asks for an unattended cycle

- GIVEN the target approves changes A and B, neither in a cycle
- WHEN the owner runs `/grove-plan an unattended cycle for A and B`,
  agrees to A and B, and merges the cycle file
- THEN the cycle file holds `unattended: true`, and a `grove claim` given
  no cycle claims it with no other merge

#### Scenario: The owner asks nothing about it

- GIVEN the target approves changes A and B, neither in a cycle
- WHEN the owner runs `/grove-plan 0.4 release` and agrees to A and B
- THEN the cycle file has no `unattended`, and the agent does not ask
  whether to add it

#### Scenario: An obstacle does not stop the mark

- GIVEN the owner asks for an unattended cycle, and member B depends on
  change D, neither in the cycle nor archived as finished
- WHEN the agent writes the cycle file
- THEN the file holds `unattended: true`, and the agent says `grove claim`
  would not claim the cycle until D is archived as finished

#### Scenario: Re-plan keeps the mark

- GIVEN the target holds cycle K with `unattended: true`, stale because
  member A was revised
- WHEN the owner runs `/grove-plan K` and says nothing about the mark
- THEN the agent sets `planned` and K still holds `unattended: true`

#### Scenario: The owner asks to remove the mark

- GIVEN the target holds cycle K with `unattended: true`
- WHEN the owner runs `/grove-plan K` and asks that K no longer run
  unattended
- THEN K's file holds no `unattended`
