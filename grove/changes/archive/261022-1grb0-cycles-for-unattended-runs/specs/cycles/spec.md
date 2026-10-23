# Changes to cycles

There is no maintained spec for `cycles` at 7ce7b91. Every requirement
below is ADDED. A cycle file is `<root>/cycles/<id>-<slug>.yaml`. Its
members are the changes its `changes` lists. A cycle is done when the
target holds every member archived, finished or dropped. A member "needs a
person" when the cycle's `needs_human` names it. A cycle is stale when the
folder of a member not yet archived differs on the target from the same
folder at the cycle's `planned` commit.

## ADDED Requirements

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

## Observed current behavior

- At 7ce7b91 there are no cycles. Rule 11 listed them as out of scope
  until #19. The parked build on branch `phase-3-cycles` (inspiration only)
  had `grove/cycles/<id>-<slug>.yaml` with `title`, `goal`, `start` and
  `end`, a `cycle` field in each member's `.grove.yaml`, `/grove-plan`, and
  `grove changes` grouping changes by cycle.
- At 7ce7b91 `grove check` reports a change's spec that is stale against
  its `base` (`docs/coordinating.md`, "When the plan is stale").
