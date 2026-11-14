# Changes to cycles

The changes that change 261114-yx2b7 proposes to the requirements of
[cycles](../../../../specs/cycles/spec.md), as it stands at 999ab56. This
is future work: nothing here describes the target branch until the change
is implemented and the maintained spec is reconciled.

The Purpose loses "asking which need a person," and the sentence defining
"needs a person".

## MODIFIED Requirements

### Requirement: Plan a cycle

When the owner runs `/grove-plan <description>`, the system shall print
instructions for writing a new cycle file holding a `title` the agent
writes from its members without asking for one, an optional one-line
`goal`, the approved changes the owner agrees to as `changes`, the time
Grove gives the plan as `created`, and the target commit the plan read as
`planned`, editing no change's files except, for each member `grove
check` reports stale on a requirement or for a named file changed since
its base, telling the owner that this edits approved work and bringing
that member up to date on the same branch as the propose instructions for
it say, before writing the cycle file.

#### Scenario: Two approved changes

- GIVEN the target approves changes A and B, neither archived, and
  `grove check` reports neither
- WHEN the owner runs `/grove-plan 0.4 release` and agrees to A and B
- THEN a new cycle file lists A and B under `changes`, `created` is the
  time Grove gave the plan, `planned` is the target's commit, and A's and
  B's files are unchanged

#### Scenario: The agent titles the cycle

- GIVEN the target approves changes A and B
- WHEN the owner runs `/grove-plan 0.4 release` and agrees to A and B
- THEN the agent writes the cycle's `title` itself and asks the owner
  nothing about it

#### Scenario: A draft change named

- GIVEN change D exists in the working tree but not in the target
- WHEN the owner asks for D in the cycle
- THEN D is not listed, and the agent says D is not approved

#### Scenario: A member stale on its base

- GIVEN the target approves change B, whose T-01 names `ledger/report.py`,
  and the merge that archived B's dependency A changed `ledger/report.py`
  after B's `base`
- WHEN the owner runs `/grove-plan 0.4 release` and agrees to B
- THEN the agent says bringing B up to date edits approved work, brings
  T-01 up to date with `ledger/report.py` as the target has it, sets B's
  `base` to the target's commit and writes the cycle file, all on the same
  branch, without stopping for another session, and `grove check` is clean

### Requirement: Say whether the cycle can run unattended

When the cycle file is written, the plan instructions shall have the agent
tell the owner that `/grove-implement <cycle-id>` builds the cycle with
them, then whether `grove claim <cycle-id>` would claim it, and that
`unattended: true`, merged on its own, hands it to a `grove claim` given
no cycle, such as a timer's; and, for each obstacle, what removes it.

#### Scenario: Nothing in the way

- GIVEN members A and B, neither of which waits on a change outside the
  cycle
- WHEN the agent finishes the plan
- THEN it names `/grove-implement <cycle-id>` first, says
  `grove claim <cycle-id>` would claim the cycle, and says a `grove claim`
  given no cycle would claim it once `unattended: true` is merged in, and
  how that mark is added

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

### Requirement: Report a cycle that cannot run unattended

When the owner runs `grove check`, the system shall report a cycle file
that holds `unattended: true` while the cycle is stale, naming what to do.

#### Scenario: Marked and stale

- GIVEN cycle K holds `unattended: true`, and member A's proposal was
  revised on the target after K's `planned` commit
- WHEN the owner runs `grove check`
- THEN it reports K, saying A changed since the plan and to re-plan K with
  `/grove-plan K`

#### Scenario: Marked and a member's file moved

- GIVEN cycle K holds `unattended: true`, member B's T-01 names
  `ledger/records.py`, and the target merged a commit touching that file
  after K's `planned` commit
- WHEN the owner runs `grove check`
- THEN it reports K, saying B changed since the plan and to re-plan K with
  `/grove-plan K`, and `grove claim` does not claim K

### Requirement: Report a malformed cycle

When the owner runs `grove check`, the system shall report a cycle file
with no `title`, no `created`, no `changes` or no `planned`, with an
unknown field, that lists a change that is no change, open or archived, or
that lists a change another cycle not yet done also lists; and, for a
cycle not yet done, a `depends_on` that names no cycle, or a cycle that
depends on it in turn, directly or through other cycles.

#### Scenario: One change in two cycles

- GIVEN cycles K and L, neither done, both list change A
- WHEN the owner runs `grove check`
- THEN it reports the later of the two, naming A and the other cycle

#### Scenario: A needs_human left over

- GIVEN cycle K's file still holds `needs_human` naming member C
- WHEN the owner runs `grove check`
- THEN it reports K's `needs_human` line as an unknown field, listing the
  fields a cycle file takes

#### Scenario: A dependency on no cycle

- GIVEN cycle K, not done, whose `depends_on` names 261099-zzzzz, and no
  cycle file has that ID
- WHEN the owner runs `grove check`
- THEN it reports K's `depends_on` line, saying 261099-zzzzz is no cycle

#### Scenario: Cycles that wait on each other

- GIVEN cycle K depends on cycle L, and L depends on K, neither done
- WHEN the owner runs `grove check`
- THEN it reports each one's `depends_on` line, saying the other waits on
  it in turn, so neither can run

## REMOVED Requirements

### Requirement: Ask about work that needs a person

Removed with `needs_human`: whether a change needs a person is settled in
review of the pull request implement ends at, not by a question and a gate.
