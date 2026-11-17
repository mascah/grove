# Changes to implement

The changes that change 261117-w31az proposes to the requirements of
[implement](../../../../specs/implement/spec.md). This is future work:
nothing here describes the target branch until the change is implemented
and the maintained spec is reconciled.

## ADDED Requirements

### Requirement: Delivery rules and a caller that delivers

Where `grove.yaml` holds lines under `rules.delivery.<skill>`, the system
shall print them in that skill's instructions, after the skill's rules,
under a heading of their own. When `grove instructions implement` is given
`--runner` before a change's or a cycle's ID, it shall leave out the lines
under `rules.delivery.implement` and print, in Grove's own words, that the
session that started this build delivers it, and that this session
pushes nothing, opens no pull request, creates no branch or worktree, and
leaves its work on the current branch in this working tree. It shall
still print the lines under `rules.skills.implement`, and `grove check`
shall accept `rules.delivery` with the same skill keys as `rules.skills`.

#### Scenario: A session the owner starts

- GIVEN `rules.delivery.implement` holds "Push the branch and open a pull
  request with gh pr create."
- WHEN `grove instructions implement C` runs
- THEN the instructions print that line under the delivery heading, after
  the project's rules for the skill

#### Scenario: A caller that delivers

- GIVEN the same `grove.yaml`, whose `rules.skills.implement` also holds
  "Commit after each task, naming its task ID."
- WHEN `grove instructions implement --runner C` runs
- THEN the instructions print the commit line, do not print the push
  line, and print Grove's paragraph saying this session delivers nothing

#### Scenario: A cycle under a runner

- GIVEN `rules.delivery.implement` holds a push line, and cycle K lists A
  and B
- WHEN `grove instructions implement --runner K` runs
- THEN the instructions do not print the push line, print Grove's
  paragraph, and name `instructions implement --runner <change-id>` for
  A and B

#### Scenario: Delivery lines left under the skill's rules

- GIVEN `rules.skills.implement` holds "Push the branch and open a pull
  request." and `rules.delivery` is absent
- WHEN `grove instructions implement --runner C` runs
- THEN the instructions print that line, as before this change, beside
  Grove's paragraph

#### Scenario: Checking the key

- GIVEN `grove.yaml` holds `rules.delivery.implement` and
  `rules.delivery.apply`
- WHEN the owner runs `grove check`
- THEN it reports `rules.delivery.apply` as an unknown key and nothing for
  `rules.delivery.implement`

#### Scenario: The flag on a skill without it

- GIVEN any project
- WHEN `grove instructions plan --runner` runs
- THEN Grove prints its usage and exits 2

## MODIFIED Requirements

### Requirement: Build a cycle's changes in their own contexts

Where the agent's harness lets it start its own subagents, the cycle's
implement instructions shall have the session building a cycle build each
member in a fresh context of its own by
`grove instructions implement --runner <change-id>`, told to neither
archive nor commit; start at once every member whose dependencies are
archived and whose tasks' `Files:` name no file a member being built
names; itself archive and commit each returned member whose tasks are all
done, one at a time; where two members built together changed the same
file, run the project's tests once both returned and, if they fail,
rebuild the later-started member on its own before archiving it; and list
every member's `Check:` lines and blocked tasks first in its report; and,
where the harness has none, the session shall build each member itself,
one after another.

#### Scenario: Two members that share nothing

- GIVEN cycle K lists A, B and C, C depending on A and B, and no file in
  A's tasks' `Files:` is in B's
- WHEN `/grove-implement K` is run in a harness with subagents
- THEN A and B are built at the same time, each in its own subagent by
  `grove instructions implement --runner <change-id>`; the session
  archives and commits each as it returns, one at a time; and C is
  started once both are archived

#### Scenario: Members that share a file

- GIVEN cycle K lists A and B, which depend on nothing, and both name
  `dev/walkthrough.md` in their tasks' `Files:`
- WHEN `/grove-implement K` is run in a harness with subagents
- THEN B's subagent is started only once A is archived

#### Scenario: A choice in a member

- GIVEN A and B of cycle K are built at the same time, and A's subagent
  returns T-01 as `Status: done` with a `Check:` line
- WHEN both have returned
- THEN the session archives both, asks no one, and its report lists A's
  T-01 `Check:` line first

#### Scenario: Two members changed the same file

- GIVEN A and B of cycle K are built at the same time, A's subagent
  changes `internal/cli/check.go` with an `Also changed:` line, B's tasks
  name that file, and the project's tests fail once both have returned
- WHEN the session checks them before archiving
- THEN B is built again on its own, after A is archived, and is archived
  only once its evidence holds

#### Scenario: A harness with no subagents

- GIVEN cycle K lists A and B, which share nothing
- WHEN `/grove-implement K` is run in a harness that cannot start
  subagents
- THEN the session builds A and then B itself, each by
  `grove instructions implement --runner <change-id>`, and delivers the
  cycle as its own instructions say
