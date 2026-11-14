# Changes to runs

The changes that change 261114-zpgx9 proposes to the requirements of
[runs](../../../../specs/runs/spec.md), as at 999ab56.

## MODIFIED Requirements

### Requirement: Start one attempt

When the operator runs `grove claim` in a clean checkout that is not
paused, and either gives a cycle ID that is a buildable cycle on the
fetched target or gives none and the fetched target holds an eligible
cycle, the system shall, when the checkout has a remote named `origin`,
fetch the target and `grove/runs` from it and fast-forward the local
target branch to the fetched target; record the attempt for the given
cycle, or else for the oldest eligible cycle, the one whose `created` is
earliest, as "Record the attempt" says; create the attempt's branch
`grove/<attempt>` from the fetched target and switch the checkout to it;
print the cycle's ID, alone, on standard output; and start no process but
Git.

#### Scenario: Local target has commits of its own

- GIVEN the checkout's local target branch holds a commit `origin` does
  not
- WHEN the operator runs `grove claim`
- THEN nothing is recorded, no branch is created, and Grove says the local
  target cannot be fast-forwarded

#### Scenario: A marked cycle of two changes, one waiting on the other

- GIVEN `origin`'s target holds cycle K, marked `unattended`, listing
  approved changes A and B, B depending on A, neither archived
- WHEN the operator runs `grove claim`
- THEN standard output is K's ID and nothing else, the checkout is on a
  new branch from the fetched target, `origin`'s `grove/runs` holds a
  record naming K, A then B, the target commit, the host and that branch,
  and no process other than Git was started

#### Scenario: Composed with a harness

- GIVEN eligible cycle K
- WHEN the operator runs
  `K=$(grove claim) && claude -p "/grove-implement $K"`
- THEN the session starts on `/grove-implement <K's ID>` in
  the checkout, on the attempt's branch

#### Scenario: Two marked cycles

- GIVEN `origin`'s target holds eligible cycles K and L planned the same
  day, K's `created` the earlier
- WHEN the operator runs `grove claim`
- THEN the record names K only

#### Scenario: The older marked cycle waits on the newer

- GIVEN `origin`'s target holds cycles K and L, both marked `unattended`,
  K's `created` the earlier, K's `depends_on` naming L, and L not done
- WHEN the operator runs `grove claim`
- THEN the record names L only

#### Scenario: A given cycle not marked

- GIVEN `origin`'s target holds buildable cycle K without `unattended`,
  and eligible cycle L
- WHEN the operator runs `grove claim K`
- THEN the record names K, and standard output is K's ID

#### Scenario: A given cycle while paused

- GIVEN the latest record on `grove/runs` names cycle L, whose members are
  not archived on the target, and K is buildable
- WHEN the operator runs `grove claim K`
- THEN nothing is recorded and Grove says the project is paused; with
  `--retry`, an attempt is recorded for K

#### Scenario: A checkout with no remote

- GIVEN the checkout has no remote named `origin`, and its local target
  branch holds eligible cycle K
- WHEN the operator runs `grove claim`
- THEN nothing is fetched, the checkout is on a new branch from the local
  target, the local branch `grove/runs` holds the record with the local
  target's commit, and standard output is K's ID

#### Scenario: The branch merged locally

- GIVEN the latest record, in a checkout with no remote, names cycle K,
  whose attempt's branch the owner built and merged into the local target,
  archiving K's members
- WHEN the operator runs `grove claim`
- THEN the project is not paused

### Requirement: Nothing to apply

When the fetched target holds no eligible cycle, `grove claim` given no
cycle shall record nothing, create no branch, print nothing on standard
output, say there is nothing to claim, and exit non-zero.

#### Scenario: Every change archived

- GIVEN every change on `origin`'s target is archived
- WHEN the operator runs `grove claim`
- THEN nothing is recorded, no branch is created, and the exit status is
  non-zero

#### Scenario: An approved change in no cycle

- GIVEN `origin`'s target approves change C, which no cycle lists
- WHEN the operator runs `grove claim`
- THEN nothing is recorded and no branch is created

#### Scenario: A planned cycle not marked

- GIVEN `origin`'s target holds cycle K listing approved change C, without
  `unattended`
- WHEN the operator runs `grove claim`
- THEN nothing is recorded

#### Scenario: A marked cycle that went stale

- GIVEN `origin`'s target holds cycle K, marked `unattended`, and member A
  was revised on the target after K's `planned` commit
- WHEN the operator runs `grove claim`
- THEN nothing is recorded

#### Scenario: Composed with a harness, nothing to claim

- GIVEN every change on `origin`'s target is archived
- WHEN the operator runs
  `K=$(grove claim) && claude -p "/grove-implement $K"`
- THEN `claude` is not started, and the shell's status is non-zero
