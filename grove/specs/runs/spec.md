# Runs

## Purpose

`grove claim` records one attempt at a cycle, the one it is given or else
the oldest marked `unattended`, as a commit on the branch `grove/runs`,
creates the attempt's branch in the checkout and prints the cycle's ID;
whoever called it starts the agent. A run is one such attempt. "The
checkout" is the Git checkout `grove claim` is run in. "The fetched target" is the target branch
as just fetched from `origin`, or, in a checkout with no remote named
`origin`, the local target branch; in both, with the `grove.yaml` it
holds. "A buildable cycle" is a cycle on the fetched target that is not
done, is not stale, has no member `needs_human` names, has no member that
waits on a change neither in the cycle nor archived as finished, and waits
on no cycle, that is, the fetched target holds done every cycle its
`depends_on` names (`grove/specs/cycles/spec.md`); "an eligible cycle" is
a buildable cycle that holds `unattended: true`.

## Boundaries

This covers `grove claim`, its cycle argument and `--retry`: claiming,
refusing, recording and pausing attempts. What the session builds, and how
task progress is recorded, belong to the implement instructions with
`--unattended` (`grove/specs/implement/spec.md`) and `implementation.md`.
Recovering an interrupted attempt, creating checkouts, opening pull
requests and scheduling are outside it (`CONSTITUTION.md` rules 2 and 11);
Grove pushes only `grove/runs`. Starting, bounding, stopping and logging
the agent are the caller's; the agent pushes and opens pull requests only
as the project's `rules` tell it, and launchd or cron calls `grove claim`.

## Requirements

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
  `K=$(grove claim) && claude -p "/grove-implement $K --unattended"`
- THEN the session starts on `/grove-implement <K's ID> --unattended` in
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

### Requirement: Refuse a given cycle that cannot be built

If the operator gives `grove claim` a cycle ID that is not a buildable
cycle on the fetched target, then the system shall record nothing, create
no branch, print nothing on standard output, exit non-zero and say in one
sentence why: there is no such cycle, it is done, a member is under
`needs_human`, it waits on a cycle its `depends_on` names that is not
done, its plan is stale, or a member waits on a change neither in the
cycle nor archived as finished.

#### Scenario: A member needs a person

- GIVEN `origin`'s target holds cycle K whose `needs_human` names member C
- WHEN the operator runs `grove claim K`
- THEN nothing is recorded, no branch is created, and Grove says C needs a
  person

#### Scenario: A stale plan

- GIVEN `origin`'s target holds cycle K, and member A was revised on the
  target after K's `planned` commit
- WHEN the operator runs `grove claim K`
- THEN nothing is recorded, and Grove says K's plan is stale

#### Scenario: A cycle it follows not done

- GIVEN `origin`'s target holds cycle K, whose `depends_on` names cycle
  L, and L's change C is not archived
- WHEN the operator runs `grove claim K`
- THEN nothing is recorded, no branch is created, and Grove says K waits
  on cycle L, which is not done

#### Scenario: No such cycle

- GIVEN `origin`'s target holds no cycle with ID Z
- WHEN the operator runs `grove claim Z`
- THEN nothing is recorded, standard output is empty, the exit status is
  non-zero, and Grove says there is no cycle Z

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
  `K=$(grove claim) && claude -p "/grove-implement $K --unattended"`
- THEN `claude` is not started, and the shell's status is non-zero

### Requirement: Refuse an unclean checkout

If the checkout has uncommitted changes, then `grove claim`, with or
without `--retry`, shall record nothing, create no branch and say the
checkout has uncommitted changes.

#### Scenario: Retrying after an interrupted edit

- GIVEN a session on an attempt's branch was stopped after the agent
  edited a file and before it committed
- WHEN the operator runs `grove claim --retry` without committing,
  stashing or discarding the edit
- THEN nothing is recorded and the edit is untouched

### Requirement: Record the attempt

When `grove claim` selects a cycle, the system shall record the attempt as
one commit on the branch `grove/runs`, with an empty tree and the record
as its message: the attempt's ID, its start time, the selected cycle and
its changes in build order, the target commit, the host's name, the
checkout and the attempt's branch; where the checkout has `origin`, push
it so that the push succeeds only while `origin`'s `grove/runs` is at the
tip just fetched, and otherwise update the local `grove/runs` only while
it is at the tip just read; and write nothing under `.git/grove/`. The
record states that an attempt started, not how any session ended nor
whether any task is done.

#### Scenario: The first attempt of a project

- GIVEN `origin` has no branch `grove/runs`
- WHEN the operator runs `grove claim` with eligible cycle K
- THEN `origin` gains `grove/runs` with one commit that holds no file,
  whose message is K's record, and the checkout has no `.git/grove/`

#### Scenario: A later attempt

- GIVEN `origin/grove/runs` holds one record
- WHEN `grove claim --retry` records an attempt
- THEN `origin/grove/runs` holds two commits, the newer the new record and
  the older as it was

#### Scenario: Read by hand

- GIVEN `grove claim` recorded an attempt for cycle K from host H
- WHEN the operator runs `git log -1 --format=%B origin/grove/runs`
- THEN it prints the record, naming K, H, the attempt's branch and when it
  started

#### Scenario: Nothing in the working tree

- WHEN `grove claim` records an attempt
- THEN `git status` shows nothing of it, and the attempt's branch holds
  only the fetched target's commits

### Requirement: Pause after an attempt

While the latest record on the fetched `grove/runs`, or on the local
`grove/runs` in a checkout with no `origin`, names a cycle a member of
which is not archived, as finished or dropped, on the fetched target,
`grove claim` without `--retry` shall record nothing, create no branch,
print nothing on standard output, exit non-zero and say the project is
paused, naming the attempt, the host that started it, when, and its
branch.

#### Scenario: Pull request not merged yet

- GIVEN the latest record names cycle K, built on its branch, and the pull
  request is open
- WHEN `grove claim` is invoked, on the same host or another
- THEN nothing is recorded, and Grove names the attempt, its host, its
  start time and its branch

#### Scenario: Another host

- GIVEN host H1's `grove claim` recorded cycle K, and a fresh clone on
  host H2
- WHEN the operator runs `grove claim` on H2
- THEN nothing is recorded, and Grove says the project is paused after
  H1's attempt

#### Scenario: Another cycle marked meanwhile

- GIVEN the latest record names cycle K, none of whose members is archived
  on the target, and cycle L was marked `unattended` and merged
- WHEN `grove claim` is invoked
- THEN nothing is recorded

#### Scenario: Cycle merged

- GIVEN the pull request archiving all of cycle K's members was merged and
  cycle L is eligible
- WHEN `grove claim` is invoked
- THEN a new attempt is recorded, for L

#### Scenario: Blocked change dropped

- GIVEN the latest record names cycle K with A and B, A is archived as
  finished on the target and B is archived as dropped
- WHEN `grove claim` is invoked
- THEN the project is not paused

### Requirement: Retry

When the operator runs `grove claim --retry` in a clean checkout, the
system shall record a new attempt as an unpaused `grove claim` does,
whatever the latest record names, and leave earlier records and their
branches as they are.

#### Scenario: A blocked change revised

- GIVEN cycle K's A was merged as finished, B was blocked by an open
  question, and the owner merged, in one squash merge, a revision of B
  answering it that re-planned K
- WHEN the operator runs `grove claim --retry`
- THEN a new branch is created from the fetched target, a new record
  naming K is pushed, and standard output is K's ID

#### Scenario: The earlier attempt kept

- GIVEN the latest record names K from host H1, whose session died
- WHEN the operator runs `grove claim --retry` on host H2
- THEN `grove/runs` holds H1's record and the new one, and H1's branch is
  untouched

### Requirement: One claim wins

When two `grove claim` commands record an attempt for one project at once,
the system shall keep the record whose ref update `origin`, or the local
repository where there is no `origin`, accepts first, and have the other,
whose update is refused because `grove/runs` moved, read `grove/runs`
again and decide again as a fresh `grove claim` would; a push refused for
any other reason shall be reported in one sentence quoting Git, and record
nothing.

#### Scenario: Two hosts at once

- GIVEN checkouts of one project on hosts H1 and H2, the project not
  paused, and eligible cycle K
- WHEN both run `grove claim` within the same second
- THEN `origin`'s `grove/runs` holds one new record, for K; one host
  prints K's ID and is on its attempt's branch; and the other prints
  nothing on standard output, says the project is paused after that
  attempt, naming the first host, and creates no branch

#### Scenario: A push refused for permissions

- GIVEN the checkout's account may not push to `origin`
- WHEN the operator runs `grove claim` with eligible cycle K
- THEN Grove says the push was refused, quoting Git's reason, records
  nothing on `origin`, leaves the local `grove/runs` as `origin` has it,
  and creates no branch

## Known discrepancies

None recorded.

## Open questions

None.

## History

- [261021-2aj5w](../../changes/archive/261021-2aj5w-run-one-bounded-apply-session/) Run one bounded apply session
- [261022-1grb0](../../changes/archive/261022-1grb0-cycles-for-unattended-runs/) Cycles planned for unattended runs
- [261023-rcxv2](../../changes/archive/261023-rcxv2-plan-verdict-no-remote-streamed-output/) Plan's verdict, runs without a remote, streamed output
- [261023-3f9h6](../../changes/archive/261023-3f9h6-run-builds-a-given-cycle/) grove run builds a cycle it is given
- [261023-yz5bk](../../changes/archive/261023-yz5bk-styled-run-output/) Styled grove run output
- [261023-g0a3h](../../changes/archive/261023-g0a3h-implement-builds-unattended/) Implement builds unattended, and apply goes
- [261023-bgtcb](../../changes/archive/261023-bgtcb-revising-a-change-re-plans-its-cycle/) Revising a change re-plans its cycle
- [261023-2ej3w](../../changes/archive/261023-2ej3w-cycles-depend-on-cycles/) A cycle can depend on cycles
- [261108-jpg8n](../../changes/archive/261108-jpg8n-claim-replaces-run/) grove claim replaces grove run, with attempt records on a branch
