# Cycles

## Purpose

A cycle is approved changes the owner means to build together, in a file
`<root>/cycles/<id>-<slug>.yaml`. Its members are the changes its
`changes` lists. A cycle is done when the target holds every member
archived, finished or dropped. A member "needs a person" when the cycle's
`needs_human` names it. A cycle is stale when the folder of a member not
yet archived, or a file one of its tasks names, differs on the target from
the same path at the cycle's `planned` commit, unless every commit on the
target's first-parent history since then that changed it also edited the
cycle file.

## Boundaries

This covers the cycle file, `/grove-plan`, the `unattended` mark, `grove
check`'s reports on cycles, and how `grove changes`, `grove review` and
the project index show them. Building a cycle, with the owner or with
`--unattended`, belongs to `grove/specs/implement/spec.md`, and building
it with no one watching to `grove/specs/runs/spec.md`. What makes a plan,
and the mark, count is the merge into the target (`CONSTITUTION.md` rule
4).

## Requirements

### Requirement: Plan a cycle

When the owner runs `/grove-plan <description>`, the system shall print
instructions for writing a new cycle file holding a `title`, an optional
one-line `goal`, the approved changes the owner agrees to as `changes`, the
time Grove gives the plan as `created`, and the target commit the plan read
as `planned`, editing no change's files.

#### Scenario: Two approved changes

- GIVEN the target approves changes A and B, neither archived
- WHEN the owner runs `/grove-plan 0.4 release` and agrees to A and B
- THEN a new cycle file lists A and B under `changes`, `created` is the
  time Grove gave the plan, `planned` is the target's commit, and A's and
  B's files are unchanged

#### Scenario: A draft change named

- GIVEN change D exists in the working tree but not in the target
- WHEN the owner asks for D in the cycle
- THEN D is not listed, and the agent says D is not approved

### Requirement: Plan one change

When the owner runs `/grove-plan <change-id>` for a change the target
approves, the system shall print instructions for writing a new cycle file
whose only member is that change.

#### Scenario: Send one change to the runner

- GIVEN the target approves change C, in no cycle
- WHEN the owner runs `/grove-plan C`
- THEN a new cycle file lists only C, with the change's title as the
  cycle's title unless the owner gives another

### Requirement: Revise a cycle

When the owner runs `/grove-plan <cycle-id>`, the system shall print
instructions for revising that cycle file as the owner asks and setting
`planned` to the target commit the revision read.

#### Scenario: Re-plan a stale cycle

- GIVEN cycle K is stale because member A was revised on the target
- WHEN the owner runs `/grove-plan K`
- THEN the agent reads A as the target holds it, asks about it as for a new
  plan, and sets `planned` to the target's commit

#### Scenario: Re-plan after a member's file moved

- GIVEN cycle K is stale because a file member B's T-01 names changed on
  the target, and the owner is content to build B on the new text
- WHEN the owner runs `/grove-plan K` and merges the cycle file
- THEN `planned` is the target's commit, K is no longer stale, and B's
  files are unchanged

### Requirement: Ask about work that needs a person

While planning, when a member's proposal, change specs or tasks call for
something an unattended session may not be able to do, the plan
instructions shall have the agent tell the owner what and ask whether the
change needs a person, and record it under `needs_human` with the owner's
reason only when the owner says it does.

#### Scenario: A task needs a human eye

- GIVEN member C has a task whose acceptance is judging how a new screen
  looks
- WHEN the agent asks and the owner says C needs them
- THEN `needs_human` names C with the owner's reason

#### Scenario: The owner says the runner can do it

- GIVEN member B has a task that calls a service needing a login
- WHEN the agent asks and the owner says the runner has that login
- THEN `needs_human` does not name B

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

#### Scenario: The runner needs something set up

- GIVEN member B has a task that calls a service needing a login, and the
  owner said the runner has it
- WHEN the agent finishes the plan
- THEN `needs_human` does not name B, and the verdict names the login the
  runner's checkout must have

### Requirement: Mark a cycle unattended

The system shall treat a cycle as handed to `grove run` only while the
target's cycle file holds `unattended: true`, and the plan instructions
shall never write or remove that field.

#### Scenario: Planned but not marked

- GIVEN the target holds cycle K without `unattended`
- WHEN `grove run` is invoked
- THEN K is not built

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
  `/grove-plan K`, and `grove run` does not build K

### Requirement: Report a malformed cycle

When the owner runs `grove check`, the system shall report a cycle file
with no `title`, no `created`, no `changes` or no `planned`, with an
unknown field, that
lists a change that is no change, open or archived, or that lists a change
another cycle not yet done also lists.

#### Scenario: One change in two cycles

- GIVEN cycles K and L, neither done, both list change A
- WHEN the owner runs `grove check`
- THEN it reports the later of the two, naming A and the other cycle

### Requirement: Show cycles

When the owner runs `grove changes`, the system shall show, with each
approved change not archived that a cycle not yet done lists, that cycle's
ID and whether it is marked `unattended`; and `grove review` shall label a
cycle file as a cycle.

#### Scenario: One marked, one not

- GIVEN change A is in marked cycle K, and change B is in no cycle
- WHEN the owner runs `grove changes`
- THEN A's line names K as unattended, and B's line names no cycle

### Requirement: Count a member changed with its plan

When a commit on the target's first-parent history after a cycle's
`planned` commit changes the folder of a member not yet archived, or a
file one of its tasks names, the system shall treat the cycle as stale, in
`grove check` and in `grove run`, unless that commit also edits the cycle
file.

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

## Known discrepancies

None recorded.

## Open questions

None.
