# Changes to implement

The changes that change 261105-pr77s proposes to the requirements of
[implement](../../../../specs/implement/spec.md), as at bc9c29b.

## ADDED Requirements

### Requirement: Build a cycle's changes in their own contexts

Where the agent's harness lets it start its own subagents, the cycle's
implement instructions shall have the session building a cycle, with the
owner or with `--unattended`, build each member in a fresh context of its
own by `grove instructions implement <change-id> --unattended`, which
neither archives nor commits; start at once every member whose
dependencies are archived and whose tasks' `Files:` name no file a member
being built names; itself archive and commit each returned member whose
tasks are all done, one at a time; where two members built together
changed the same file, run the project's tests once both returned and,
if they fail, rebuild the later-started member on its own before
archiving it; and, with the owner, put the checkpoints each returned
member holds to the owner at once and rerun that member with their words,
or, with `--unattended`, list every checkpoint first in its report; and,
where the harness has none, the session shall build each member itself,
one after another, in its own mode.

#### Scenario: Two members that share nothing

- GIVEN cycle K lists A, B and C, C depending on A and B, and no file in
  A's tasks' `Files:` is in B's
- WHEN `/grove-implement K` is run in a harness with subagents
- THEN A and B are built at the same time, each in its own subagent by
  `grove instructions implement <change-id> --unattended`; the session
  archives and commits each as it returns, one at a time; and C is
  started once both are archived

#### Scenario: Members that share a file

- GIVEN cycle K lists A and B, which depend on nothing, and both name
  `dev/walkthrough.md` in their tasks' `Files:`
- WHEN `/grove-implement K` is run in a harness with subagents
- THEN B's subagent is started only once A is archived

#### Scenario: A checkpoint, with the owner

- GIVEN A and B of cycle K are built at the same time, and A's subagent
  returns T-01 as `Status: to check` with a `Check:` line
- WHEN both have returned
- THEN the session archives B, puts A's T-01 to the owner, and reruns
  A's subagent with the owner's words, which it writes under T-01 as
  `Answer:`

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
  `grove instructions implement <change-id>` in the session's own mode

## MODIFIED Requirements

### Requirement: Build a change, a draft or a cycle

When `/grove-implement <id>` is run, the system shall print instructions
for building, on a branch other than the target, the approved change
`<id>`, or each member of cycle `<id>` not yet archived in build order,
each listed with the members of the cycle it depends on and each member's
instructions naming its place in the cycle's plan; or, with
the owner, a change only the working tree holds, as a draft that merging
the branch approves with its code. The agent records each task's status
in `implementation.md`, builds no task already done again unless the
check before finishing finds its evidence no longer holds, and archives
the change when every task is done. The system shall print no
instructions, and say why, for a change that waits on one the working
tree has not archived as finished, or for a draft given `--unattended`;
and a cycle's instructions shall leave out, saying why, a member that
waits on a change neither being built nor archived as finished.

#### Scenario: An approved change

- GIVEN the target approves change C with tasks T-01 and T-02
- WHEN `/grove-implement C` is run on the target branch, with or without
  `--unattended`
- THEN the agent creates a branch, builds T-01 and T-02, records both as
  done in `implementation.md`, and archives C on that branch

#### Scenario: A draft

- GIVEN the owner proposed change D on branch `rename` and committed it,
  and the target does not hold D
- WHEN the owner runs `/grove-implement D` on `rename`
- THEN the instructions say D is a draft that merging `rename` approves
  with its code, and the agent builds and archives D on `rename`; with
  `--unattended`, Grove says the target holds no approved change D

#### Scenario: A change that waits

- GIVEN change B, approved or a draft, with `depends_on: [A]`, and the
  working tree does not hold A archived as finished
- WHEN `/grove-implement B` is run, with or without `--unattended`
- THEN nothing is built, and Grove says B waits on A

#### Scenario: A cycle

- GIVEN cycle K lists A and B, B depending on A, K not marked
  `unattended` and its `needs_human` naming B
- WHEN `/grove-implement K` is run, with or without `--unattended`
- THEN K's instructions list B `after A`, and A and then B are built on
  one branch, each by `grove instructions implement <change-id>`, B's
  instructions naming it the second of two in K's plan; only `grove run`
  reads `unattended` and `needs_human`

#### Scenario: A member left out

- GIVEN cycle K lists A and C, and C depends on Z, which K does not list
  and the target has not archived as finished
- WHEN `/grove-implement K` is run
- THEN the instructions leave C out, saying it waits on Z, and A is built
