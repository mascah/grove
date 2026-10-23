# Changes to cycles

The maintained spec is [grove/specs/cycles/spec.md](../../../../specs/cycles/spec.md)
(at ab885af). This change depends on change 261023-bgtcb, and starts from
the text it gives.

The Purpose gains: a cycle "depends on" each cycle its optional
`depends_on` names, and "waits on" each of those the target does not hold
done. Its definition of stale, as change 261023-bgtcb leaves it, becomes:
a cycle is stale when the folder of a member not yet archived, or a file
one of its tasks names, differs on the target from the same path at the
cycle's `planned` commit, unless every commit on the target's first-parent
history since then that changed it also edited the cycle file, or archived
a change that a cycle this one depends on, directly or through other
cycles, lists.

"Count a member changed with its plan" is ADDED by change 261023-bgtcb,
so the maintained spec holds it before this change is archived. This
change MODIFIES it rather than adding a second requirement: two
requirements each saying when a member's change makes a cycle stale would
contradict each other, since each would count the commit the other
excuses. "Say whether the cycle can run unattended" starts from change
261023-bgtcb's text too; only its new scenario differs.

## ADDED Requirements

### Requirement: A cycle can depend on cycles

Where a cycle file's `depends_on` names cycles, the system shall have
`grove run` build none of that cycle, given or marked, while the target
does not hold one of them done, and have `/grove-implement <cycle-id>`
build it as it would with no `depends_on`.

#### Scenario: The cycle it follows is not done

- GIVEN cycle K's `depends_on` names cycle L, and L's change C is not
  archived on the target
- WHEN the operator runs `grove run K`
- THEN no session is launched, and Grove says K waits on cycle L, which
  is not done

#### Scenario: The cycle it follows is done

- GIVEN cycle K's `depends_on` names cycle L, every member of L is
  archived on the target, and nothing else keeps K from being built
- WHEN the operator runs `grove run K`
- THEN a session is launched on K

#### Scenario: Built with the owner

- GIVEN cycle K's `depends_on` names cycle L, which is not done
- WHEN the owner runs `/grove-implement K`
- THEN the instructions list K's members in build order, as they would
  with no `depends_on`

### Requirement: Ask which cycles a new cycle follows

While planning a new cycle, when the project index lists cycles not yet
done, the plan instructions shall have the agent ask the owner whether
the new cycle follows any of them, and write under `depends_on` only the
cycles the owner names.

#### Scenario: The owner names one

- GIVEN the project index lists cycles K and L, neither done
- WHEN the owner runs `/grove-plan 0.5 release` and, asked, says it
  follows L
- THEN the new cycle file holds `depends_on: [L]`

#### Scenario: The owner names none

- GIVEN the project index lists cycle K, not done
- WHEN the agent asks and the owner says the new cycle follows no cycle
- THEN the new cycle file has no `depends_on`

#### Scenario: No other cycle

- GIVEN the project index lists no cycle that is not done
- WHEN the owner runs `/grove-plan 0.5 release`
- THEN the agent asks nothing about other cycles, and the cycle file has
  no `depends_on`

## MODIFIED Requirements

### Requirement: Count a member changed with its plan

When a commit on the target's first-parent history after a cycle's
`planned` commit changes the folder of a member not yet archived, or a
file one of its tasks names, the system shall treat the cycle as stale, in
`grove check` and in `grove run`, unless that commit also edits the cycle
file, or archives a change that a cycle this one depends on, directly or
through other cycles, lists.

#### Scenario: Revised and re-planned in one squash merge

- GIVEN cycle K, marked `unattended`, lists B, and on a branch B's
  proposal is revised and K's `planned` set to the target's commit
- WHEN the branch is squash-merged into the target and the owner runs
  `grove check`
- THEN it reports nothing for K, and `grove run` builds K

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
- WHEN the operator runs `grove run K`
- THEN no session is launched, and Grove says K's plan is stale

#### Scenario: The cycle it follows lands in a squash merge

- GIVEN cycle K, marked `unattended`, depends on cycle L, K's member B
  has a task naming `docs/reference.md`, and a squash merge into the
  target archives L's only change and edits `docs/reference.md`
- WHEN the owner runs `grove check`
- THEN it reports nothing for K, and `grove run` builds K

#### Scenario: Stacked pull requests

- GIVEN cycle K depends on cycle L, whose changes C and D were built as
  stacked pull requests, and the target took them as two commits on its
  first-parent history, each archiving one of them and editing a file K's
  member B names
- WHEN the operator runs `grove run K`
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

### Requirement: Say whether the cycle can run unattended

When the cycle file is written, the plan instructions shall have the agent
tell the owner that `/grove-implement <cycle-id>` builds the cycle with
them, then whether `grove run <cycle-id>` would build it, and that
`unattended: true`, merged on its own, hands it to a `grove run` given no
cycle, such as a timer's; and, for each obstacle, what removes it.

#### Scenario: Nothing in the way

- GIVEN members A and B, neither of which the owner said needs a person,
  and neither waits on a change outside the cycle
- WHEN the agent finishes the plan
- THEN it names `/grove-implement <cycle-id>` first, says
  `grove run <cycle-id>` would build the cycle, and says a `grove run`
  given no cycle would build it once `unattended: true` is merged in, and
  how that mark is added

#### Scenario: A member needs a person

- GIVEN the owner said member C needs a person, for judging how a new
  screen looks
- WHEN the agent finishes the plan
- THEN it names `/grove-implement <cycle-id>`, says `grove run` would
  refuse the cycle, given or marked, and that revising C with
  `/grove-propose C` so that its acceptance can be checked without a
  person, which re-plans the cycle on the same branch, and merging that
  removes it

#### Scenario: A member waits on a change outside the cycle

- GIVEN member B depends on change D, neither in the cycle nor archived as
  finished
- WHEN the agent finishes the plan
- THEN it says `grove run` would not build the cycle until D is archived as
  finished

#### Scenario: The cycle waits on a cycle

- GIVEN the owner said the new cycle follows cycle L, which is not done
- WHEN the agent finishes the plan
- THEN it says `grove run` would not build the cycle until L is done, and
  that the merges archiving L's changes leave its plan as it checked it

#### Scenario: The runner needs something set up

- GIVEN member B has a task that calls a service needing a login, and the
  owner said the runner has it
- WHEN the agent finishes the plan
- THEN `needs_human` does not name B, and the verdict names the login the
  runner's checkout must have

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

### Requirement: Show cycles

When the owner runs `grove changes`, the system shall show, with each
approved change not archived that a cycle not yet done lists, that cycle's
ID, whether it is marked `unattended`, and each cycle it waits on; and
`grove review` shall label a cycle file as a cycle.

#### Scenario: One marked, one not

- GIVEN change A is in marked cycle K, and change B is in no cycle
- WHEN the owner runs `grove changes`
- THEN A's line names K as unattended, and B's line names no cycle

#### Scenario: A cycle that waits on a cycle

- GIVEN change A is in cycle K, whose `depends_on` names cycle L, which
  is not done
- WHEN the owner runs `grove changes`
- THEN A's line names K and says it waits on cycle L
