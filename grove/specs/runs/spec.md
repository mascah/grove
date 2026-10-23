# Runs

## Purpose

`grove run` builds a cycle the owner marked `unattended` without anyone
typing `/grove-apply`: in a checkout an operator prepared for one project,
it starts one bounded Claude Code session on the apply skill for that
cycle, records the attempt, and then pauses the project until the cycle is
archived. "The checkout" is the Git checkout `grove run` is run in. "The
fetched target" is the target branch as just fetched from `origin`, or,
in a checkout with no remote named `origin`, the local target branch; in
both, with the `grove.yaml` it holds. "An eligible cycle" is a cycle on the fetched
target that holds `unattended: true`, is not done, is not
stale, has no member `needs_human` names, and has no member that waits on
a change neither in the cycle nor archived as finished
(`grove/specs/cycles/spec.md`).

## Boundaries

This covers `grove run`, `--retry` and `--timeout`: starting, refusing,
stopping, recording and pausing attempts. What the session builds, and how
task progress is recorded, belong to the apply instructions and
`implementation.md`. Recovering an interrupted attempt, creating checkouts,
pushing, opening pull requests, scheduling and harnesses other than Claude
Code are outside it (`CONSTITUTION.md` rules 2 and 11); the agent pushes
and opens pull requests only as the project's `rules` tell it, and launchd
or cron calls `grove run`.

## Requirements

### Requirement: Start one attempt

When the operator runs `grove run` in a clean checkout that is not paused,
and the fetched target holds an eligible cycle, the system shall, when the
checkout has a remote named `origin`, fetch the target from it and
fast-forward the local target branch to it; then create the attempt's
branch from the fetched target, record the attempt for the oldest eligible
cycle, the one whose `created` is earliest, launch one Claude Code session
on the apply skill with that cycle's ID, and wait for the session to end.

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

- GIVEN `origin`'s target holds eligible cycles K and L planned the same
  day, K's `created` the earlier
- WHEN the operator runs `grove run`
- THEN the attempt selects K only

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
  target, one session is launched on the apply skill with K's ID, and the
  attempt's record holds the local target's commit

#### Scenario: The branch merged locally

- GIVEN the latest attempt, in a checkout with no remote, built cycle K on
  its branch, and the owner merged that branch into the local target
- WHEN the operator runs `grove run`
- THEN the project is not paused

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

### Requirement: One run per checkout

While a `grove run` is running in the checkout, the system shall make any
other `grove run` in the same checkout exit without launching a session,
saying a run is in progress. A `grove run` killed outright releases the
checkout, though its session may continue; that case is the operator's.

#### Scenario: A timer fires during a run

- GIVEN a session launched by `grove run` is still running
- WHEN `grove run` is invoked again in the same checkout
- THEN it exits without launching a session, and the first run continues

#### Scenario: Two projects

- GIVEN a session running in project P's checkout
- WHEN `grove run` is invoked in project Q's checkout
- THEN Q's run proceeds as if P's were not running

### Requirement: Refuse an unclean checkout

If the checkout has uncommitted changes, then `grove run`, with or without
`--retry`, shall launch no session and say the checkout has uncommitted
changes.

#### Scenario: Retrying after an interrupted edit

- GIVEN an attempt was cancelled after the agent edited a file and before it
  committed
- WHEN the operator runs `grove run --retry` without committing, stashing or
  discarding the edit
- THEN no session is launched and the edit is untouched

### Requirement: Stop a session

When the session runs past the time limit (`--timeout`, 2 hours unless
given) or `grove run` is interrupted or terminated, the system shall stop
every process in the session's process group, record the attempt as timed
out or cancelled, and leave its branch, files and output in place. When
the session exits by itself, the system shall stop any process left in its
group before the run ends. Processes that left the group are not stopped.

#### Scenario: Time limit reached

- GIVEN `grove run --timeout 1m` launched a session that is still running
  after a minute, with a child process that ignores the request to stop
- WHEN the minute passes
- THEN no process in the session's group is left running, the attempt is
  recorded as timed out, and the branch, its uncommitted files and the
  output remain

#### Scenario: Operator cancels

- GIVEN a running `grove run`
- WHEN the operator presses Ctrl-C, or launchd stops the job
- THEN no process in the session's group is left running and the attempt
  is recorded as cancelled

#### Scenario: Session exits leaving a child behind

- GIVEN the session started a child in its process group and then exited
- WHEN `grove run` records the attempt
- THEN the child is no longer running

### Requirement: Show the session as it runs

While the session runs, the system shall write every event the session
emits to the attempt's output file as JSON Lines, and print to standard
output one line for each assistant message, each tool call and each tool
result, as it happens.

#### Scenario: Watching a run

- GIVEN `grove run` launched a session on an eligible cycle
- WHEN the agent runs a tool
- THEN, before the session ends, standard output shows a line naming the
  tool and the first line of its input, and then a line with the first
  line of its result

#### Scenario: A stopped session

- GIVEN a session that made several tool calls was stopped at the time
  limit
- WHEN the operator reads the attempt's output file
- THEN it holds each event up to the stop, one JSON object per line

#### Scenario: Started by a timer

- GIVEN launchd started `grove run` with standard output sent to a file
- WHEN the session runs
- THEN that file receives the same lines a terminal would

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

## Known discrepancies

None recorded.

## Open questions

None.
