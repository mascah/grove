# Changes to runs

Changes to the maintained spec [`runs`](../../../../specs/runs/spec.md) as
it reads at 7ce7b91. "An eligible cycle" replaces "eligible changes": a
cycle on the fetched target that holds `unattended: true`, is not done, is
not stale, has no member `needs_human` names, and has no member that waits
on a change neither in the cycle nor archived as finished (see
[cycles](../cycles/spec.md)). The maintained spec's Purpose should say so
when this change is archived.

## MODIFIED Requirements

### Requirement: Start one attempt

When the operator runs `grove run` in a clean checkout that is not paused,
and the fetched target holds an eligible cycle, the system shall
fast-forward the local target branch to the fetched target, create the
attempt's branch from it, record the attempt for the oldest eligible
cycle, launch one Claude Code session on the apply skill with that cycle's
ID, and wait for the session to end.

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
  Code session is launched on the apply skill with K's ID, and the
  attempt's record names K, A then B, and the target commit

#### Scenario: Two marked cycles

- GIVEN `origin`'s target holds eligible cycles K and L, K the older
- WHEN the operator runs `grove run`
- THEN the attempt selects K only

#### Scenario: A project that does not list Claude Code

- GIVEN the fetched target's `grove.yaml` has `agents` without `claude`
- WHEN the operator runs `grove run`
- THEN no session is launched, and Grove says `run` needs Claude Code
  listed under `agents`

### Requirement: Nothing to apply

When the fetched target holds no eligible cycle, `grove run` shall launch
no session, create no branch and say there is nothing to apply.

#### Scenario: Every change archived

- GIVEN every change on `origin`'s target is archived
- WHEN the operator runs `grove run`
- THEN no session is launched and no branch is created

#### Scenario: An approved change in no cycle

- GIVEN `origin`'s target approves change C, which no cycle lists
- WHEN the operator runs `grove run`
- THEN no session is launched and no branch is created

#### Scenario: A planned cycle not marked

- GIVEN `origin`'s target holds cycle K listing approved change C, without
  `unattended`
- WHEN the operator runs `grove run`
- THEN no session is launched

#### Scenario: A marked cycle that went stale

- GIVEN `origin`'s target holds cycle K, marked `unattended`, and member A
  was revised on the target after K's `planned` commit
- WHEN the operator runs `grove run`
- THEN no session is launched

### Requirement: Record the attempt

When an attempt ends, the system shall record its end time and how the
session ended (its exit status, timed out, cancelled, or failed to start)
beside its start time, selected cycle and its changes, target commit,
checkout and branch, under the checkout's `.git/grove/`, and print one
sentence naming the attempt, how it ended and its output file. The record
states how the process ended, not whether any task is done.

#### Scenario: A session exits normally with a task blocked

- GIVEN the agent finished A, marked a task of B `Status: blocked` and
  exited with status 0
- WHEN `grove run` returns
- THEN the record says the session exited 0, `implementation.md` alone says
  which tasks are done or blocked, and nothing under `.git/grove/` appears
  in `git status`

#### Scenario: Claude Code cannot be started

- GIVEN `claude` is not on `PATH`
- WHEN the operator runs `grove run` with an eligible cycle
- THEN the attempt is recorded as failed to start, its branch remains, and
  the next `grove run` without `--retry` says the project is paused

### Requirement: Pause after an attempt

While a member of the cycle the latest attempt selected is not archived, as
finished or dropped, on the fetched target, `grove run` without `--retry`
shall launch no session and say the project is paused, naming the attempt
and its output file.

#### Scenario: Pull request not merged yet

- GIVEN the latest attempt built cycle K's A and B on its branch and the
  pull request is open
- WHEN `grove run` is invoked
- THEN no session is launched

#### Scenario: Another cycle marked meanwhile

- GIVEN the latest attempt selected cycle K, none of whose members is
  archived on the target, and cycle L was marked `unattended` and merged
- WHEN `grove run` is invoked
- THEN no session is launched

#### Scenario: Cycle merged

- GIVEN the pull request archiving all of cycle K's members was merged and
  cycle L is eligible
- WHEN `grove run` is invoked
- THEN a new attempt starts and selects L

#### Scenario: Blocked change dropped

- GIVEN the latest attempt selected cycle K with A and B, A is archived as
  finished on the target and B is archived as dropped
- WHEN `grove run` is invoked
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
  launched on the apply skill with K's ID, which builds B
