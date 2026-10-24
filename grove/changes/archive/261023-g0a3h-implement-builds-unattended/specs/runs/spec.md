# Changes to runs

The maintained spec is [grove/specs/runs/spec.md](../../../../specs/runs/spec.md).

## MODIFIED Requirements

### Requirement: Start one attempt

When the operator runs `grove run` in a clean checkout that is not paused,
and either gives a cycle ID that is a buildable cycle on the fetched target
or gives none and the fetched target holds an eligible cycle, the system
shall, when the checkout has a remote named `origin`, fetch the target from
it and fast-forward the local target branch to it; then create the
attempt's branch from the fetched target, record the attempt for the given
cycle, or else for the oldest eligible cycle, the one whose `created` is
earliest, launch one Claude Code session on the implement skill with that
cycle's ID and `--unattended`, typed `/grove-implement <cycle-id>
--unattended` where the project lists Claude Code as skills and
`/grove:implement <cycle-id> --unattended` where it lists it as commands,
and wait for the session to end.

#### Scenario: Local target has commits of its own

- GIVEN the checkout's local target branch holds a commit `origin` does not
- WHEN the operator runs `grove run`
- THEN no branch is created and no session is launched, and Grove says the
  local target cannot be fast-forwarded

#### Scenario: A marked cycle of two changes, one waiting on the other

- GIVEN `origin`'s target holds cycle K, marked `unattended`, listing
  approved changes A and B, B depending on A, neither archived
- WHEN the operator runs `grove run`
- THEN the checkout is on a new branch from the fetched target, one Claude
  Code session is launched on `/grove-implement K --unattended`, and the
  attempt's record names K, A then B, and the target commit

#### Scenario: Claude Code listed as commands

- GIVEN the fetched target's `grove.yaml` has `agents` with `claude:
  commands`, and eligible cycle K
- WHEN the operator runs `grove run`
- THEN one session is launched on `/grove:implement K --unattended`

#### Scenario: Two marked cycles

- GIVEN `origin`'s target holds eligible cycles K and L planned the same
  day, K's `created` the earlier
- WHEN the operator runs `grove run`
- THEN the attempt selects K only

#### Scenario: A given cycle not marked

- GIVEN `origin`'s target holds buildable cycle K without `unattended`,
  and eligible cycle L
- WHEN the operator runs `grove run K`
- THEN one session is launched on the implement skill with K's ID and
  `--unattended`, and the attempt's record names K

#### Scenario: A given cycle while paused

- GIVEN the latest attempt selected cycle L, whose members are not
  archived on the target, and K is buildable
- WHEN the operator runs `grove run K`
- THEN no session is launched and Grove says the project is paused; with
  `--retry`, an attempt starts and selects K

#### Scenario: A project that does not list Claude Code

- GIVEN the fetched target's `grove.yaml` has `agents` without `claude`
- WHEN the operator runs `grove run`
- THEN no session is launched, and Grove says `run` needs Claude Code
  listed under `agents`

#### Scenario: A checkout with no remote

- GIVEN the checkout has no remote named `origin`, and its local target
  branch holds eligible cycle K
- WHEN the operator runs `grove run`
- THEN nothing is fetched, the checkout is on a new branch from the local
  target, one session is launched on the implement skill with K's ID and
  `--unattended`, and the attempt's record holds the local target's commit

#### Scenario: The branch merged locally

- GIVEN the latest attempt, in a checkout with no remote, built cycle K on
  its branch, and the owner merged that branch into the local target
- WHEN the operator runs `grove run`
- THEN the project is not paused

### Requirement: Retry

When the operator runs `grove run --retry` in a clean checkout, the system
shall start a new attempt as an unpaused `grove run` does, whatever the
latest attempt selected.

#### Scenario: A blocked change revised

- GIVEN cycle K's A was merged as finished, B was blocked by an open
  question, the owner merged a revision of B answering it, and re-planned K
- WHEN the operator runs `grove run --retry`
- THEN a new branch is created from the fetched target and a session is
  launched on the implement skill with K's ID and `--unattended`, which
  builds B
