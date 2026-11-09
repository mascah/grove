# Changes to cycles

The changes that change 261108-jpg8n proposes to the requirements of
[cycles](../../../../specs/cycles/spec.md): `grove claim` takes the place
of `grove run` wherever a requirement names it. This is future work.

The Purpose's sentence "and `grove run <cycle-id>` builds it, as does a
`grove run` given no cycle once it holds `unattended: true`" becomes "and
`grove claim <cycle-id>` claims it for the agent its caller starts, as does
a `grove claim` given no cycle once it holds `unattended: true`".

## MODIFIED Requirements

### Requirement: Say whether the cycle can run unattended

When the cycle file is written, the plan instructions shall have the agent
tell the owner that `/grove-implement <cycle-id>` builds the cycle with
them, then whether `grove claim <cycle-id>` would claim it, and that
`unattended: true`, merged on its own, hands it to a `grove claim` given
no cycle, such as a timer's; and, for each obstacle, what removes it.

#### Scenario: Nothing in the way

- GIVEN members A and B, neither of which the owner said needs a person,
  and neither waits on a change outside the cycle
- WHEN the agent finishes the plan
- THEN it names `/grove-implement <cycle-id>` first, says
  `grove claim <cycle-id>` would claim the cycle, and says a `grove claim`
  given no cycle would claim it once `unattended: true` is merged in, and
  how that mark is added

#### Scenario: A member needs a person

- GIVEN the owner said member C needs a person, for judging how a new
  screen looks
- WHEN the agent finishes the plan
- THEN it names `/grove-implement <cycle-id>`, says `grove claim` would
  refuse the cycle, given or marked, and that revising C with
  `/grove-propose C` so that its acceptance can be checked without a
  person, which re-plans the cycle on the same branch, and merging that
  removes it

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
- THEN `needs_human` does not name B, and the verdict names the login the
  runner's checkout must have

### Requirement: Mark a cycle unattended

The system shall treat a cycle as handed to `grove claim` only while the
target's cycle file holds `unattended: true`, and the plan instructions
shall never write or remove that field.

#### Scenario: Planned but not marked

- GIVEN the target holds cycle K without `unattended`
- WHEN `grove claim` is invoked
- THEN K is not claimed

### Requirement: A cycle can depend on cycles

Where a cycle file's `depends_on` names cycles, the system shall have
`grove claim` claim none of that cycle, given or marked, while the target
does not hold one of them done, and have `/grove-implement <cycle-id>`
build it as it would with no `depends_on`.

#### Scenario: The cycle it follows is not done

- GIVEN cycle K's `depends_on` names cycle L, and L's change C is not
  archived on the target
- WHEN the operator runs `grove claim K`
- THEN nothing is recorded, no branch is created, and Grove says K waits
  on cycle L, which is not done

#### Scenario: The cycle it follows is done

- GIVEN cycle K's `depends_on` names cycle L, every member of L is
  archived on the target, and nothing else keeps K from being claimed
- WHEN the operator runs `grove claim K`
- THEN an attempt is recorded for K, and K's ID is printed

#### Scenario: Built with the owner

- GIVEN cycle K's `depends_on` names cycle L, which is not done
- WHEN the owner runs `/grove-implement K`
- THEN the instructions list K's members in build order, as they would
  with no `depends_on`

### Requirement: Report a cycle that cannot run unattended

When the owner runs `grove check`, the system shall report a cycle file
that holds `unattended: true` while `needs_human` names a member, and a
cycle file that holds `unattended: true` while the cycle is stale, each
naming what to do.

#### Scenario: Marked with a change that needs a person

- GIVEN cycle K holds `unattended: true` and `needs_human` names C
- WHEN the owner runs `grove check`
- THEN it reports K's `unattended` line, saying C needs a person

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

### Requirement: Count a member changed with its plan

When a commit on the target's first-parent history after a cycle's
`planned` commit changes the folder of a member not yet archived, or a
file one of its tasks names, the system shall treat the cycle as stale, in
`grove check` and in `grove claim`, unless that commit also edits the
cycle file, archives a change the cycle itself lists, or archives a change
that a cycle this one depends on, directly or through other cycles, lists.

#### Scenario: Revised and re-planned in one squash merge

- GIVEN cycle K, marked `unattended`, lists B, and on a branch B's
  proposal is revised and K's `planned` set to the target's commit
- WHEN the branch is squash-merged into the target and the owner runs
  `grove check`
- THEN it reports nothing for K, and `grove claim` claims K

#### Scenario: Revised and re-planned in one merge commit

- GIVEN a branch with one commit revising member B and another setting
  cycle K's `planned`
- WHEN the branch is merged into the target with a merge commit
- THEN K is not stale

#### Scenario: Marked after a member changed

- GIVEN member B's proposal was revised in one merge into the target, and
  cycle K was marked `unattended` in a later merge
- WHEN the owner runs `grove check`
- THEN it reports K, saying B changed since the plan and to re-plan K with
  `/grove-plan K`

#### Scenario: Revision and re-plan as two commits on the target

- GIVEN the revision of member B and the edit to cycle K's file reached
  the target as two commits, by rebase or fast-forward
- WHEN the operator runs `grove claim K`
- THEN nothing is recorded, and Grove says K's plan is stale

#### Scenario: A member lands in its own squash merge

- GIVEN cycle K, marked `unattended`, lists A and B, each with a task
  naming `ledger/cli.py`, and a squash merge into the target archives A,
  edits `ledger/cli.py` and leaves K's file as it is
- WHEN the owner runs `grove check`
- THEN it reports nothing for K, `grove claim` claims K, and
  `grove instructions implement B` names no file changed since K's plan

#### Scenario: Then a commit archiving nothing

- GIVEN that squash merge, and a later commit on the target that archives
  no change and edits `ledger/cli.py`
- WHEN the owner runs `grove check`
- THEN it reports K, saying B changed since the plan and to re-plan K with
  `/grove-plan K`

#### Scenario: The cycle it follows lands in a squash merge

- GIVEN cycle K, marked `unattended`, depends on cycle L, K's member B
  has a task naming `docs/reference.md`, and a squash merge into the
  target archives L's only change and edits `docs/reference.md`
- WHEN the owner runs `grove check`
- THEN it reports nothing for K, and `grove claim` claims K

#### Scenario: Stacked pull requests

- GIVEN cycle K depends on cycle L, whose changes C and D were built as
  stacked pull requests, and the target took them as two commits on its
  first-parent history, each archiving one of them and editing a file K's
  member B names
- WHEN the operator runs `grove claim K`
- THEN K's plan is not stale

#### Scenario: Through another cycle

- GIVEN cycle K depends on cycle L, L depends on cycle M, and a merge
  commit into the target archives M's change E and edits a file K's
  member B names
- WHEN the owner runs `grove check`
- THEN K is not stale

#### Scenario: Any other commit still counts

- GIVEN cycle K, marked `unattended`, depends on cycle L, and a commit on
  the target archives no change of L, nor of any cycle L depends on, and
  edits a file K's member B names
- WHEN the owner runs `grove check`
- THEN it reports K, saying B changed since the plan and to re-plan K with
  `/grove-plan K`
