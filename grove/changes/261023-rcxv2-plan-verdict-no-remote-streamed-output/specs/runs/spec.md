# Changes to runs

The changes that change 261023-rcxv2 proposes to the requirements of
`runs` (`grove/specs/runs/spec.md` at 5383ad8). The Purpose's definition
of "the fetched target" becomes: the target branch as just fetched from
`origin`, or, in a checkout with no remote named `origin`, the local
target branch; in both, with the `grove.yaml` it holds. The other
requirements, which read "the fetched target", keep their words.

## MODIFIED Requirements

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

## ADDED Requirements

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
