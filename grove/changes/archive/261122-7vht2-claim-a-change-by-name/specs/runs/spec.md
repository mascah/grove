# Changes to runs

The changes that change 261122-7vht2 proposes to the requirements of
[runs](../../../../specs/runs/spec.md), as it stands at 07d319f. This is
future work: nothing here describes the target branch until the change is
implemented and the maintained spec is reconciled.

The Purpose becomes: `grove claim` records one attempt at a piece of
work, the cycle or change it is given or else the oldest cycle marked
`unattended`, as a commit on the branch `grove/runs`, creates the
attempt's branch in the checkout and prints the piece's ID; whoever
called it starts the agent. A run is one such attempt. "The checkout",
"the fetched target", "a buildable cycle" and "an eligible cycle" keep
their definitions. "A loose change" is an approved change the fetched
target holds, not archived, that no cycle not done lists. "A buildable
change" is a loose change that waits on nothing: the fetched target holds
archived as finished every change its `depends_on` names. "A piece" is a
cycle not done or a loose change, and "a buildable piece" a buildable
cycle or a buildable change. A change is "in flight" while a record on
`grove/runs` names it and the fetched target does not hold it archived,
as finished or dropped; its attempt is the newest record that names it; a
piece is in flight while one of its changes is; and a record is in flight
while it is the attempt of a change in flight. A change's "planned files"
are the paths its tasks' `Files:` lines name and, for each capability it
has a change spec for, that capability's maintained spec; a piece's are
those of its changes not archived. Two pieces "overlap" when a planned
file of one is a planned file of the other or lies under a folder that
is. All of it is read from the fetched target's documents and from
`grove/runs`, never from an attempt's branch.

The Boundaries say this covers `grove claim`, its argument and `--retry`,
and that which piece is claimed next and how many attempts run at once
are the caller's (`CONSTITUTION.md` rule 2).

## ADDED Requirements

### Requirement: Refuse a given piece that cannot be built

If the operator gives `grove claim` an ID that is not a buildable piece
on the fetched target, then the system shall record nothing, create no
branch, print nothing on standard output, exit non-zero and say in one
sentence why: there is no such cycle or change; the cycle is done, waits
on a cycle its `depends_on` names that is not done, has a stale plan, or
has a member that waits on a change neither in the cycle nor archived as
finished; or the change is archived, is listed by a cycle not done, which
it names as the one to claim, or waits on a change not archived as
finished.

#### Scenario: A needs_human left over

- GIVEN `origin`'s target holds cycle K, buildable, whose file still holds
  `needs_human` naming member C
- WHEN the operator runs `grove claim K`
- THEN an attempt is recorded for K, and K's ID is printed

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

#### Scenario: No such cycle or change

- GIVEN `origin`'s target holds no cycle and no change with ID Z
- WHEN the operator runs `grove claim Z`
- THEN nothing is recorded, standard output is empty, the exit status is
  non-zero, and Grove says there is no cycle or change Z

#### Scenario: A change that waits

- GIVEN `origin`'s target approves loose changes A and B, B with
  `depends_on: [A]`, and A is not archived
- WHEN the operator runs `grove claim B`
- THEN nothing is recorded, and Grove says B waits on A, which is not
  archived as finished

#### Scenario: Its dependency finished only on a branch

- GIVEN B depends on A, and A is archived as finished on its attempt's
  branch, whose pull request is open
- WHEN the operator runs `grove claim B`
- THEN nothing is recorded, and Grove says B waits on A

#### Scenario: A member of a cycle

- GIVEN `origin`'s target holds cycle K, not done, listing approved
  change A
- WHEN the operator runs `grove claim A`
- THEN nothing is recorded, and Grove says A is in cycle K and to claim K

#### Scenario: An archived change

- GIVEN `origin`'s target holds change A archived
- WHEN the operator runs `grove claim A`
- THEN nothing is recorded, and Grove says A is archived

### Requirement: Refuse what is in flight or overlaps it

If the operator gives `grove claim`, without `--retry`, a buildable piece
that is in flight or that overlaps a piece in flight, then the system
shall record nothing, create no branch, print nothing on standard output,
exit non-zero and say in one sentence which: that the piece is in flight,
naming the attempt, the host that started it, when, and its branch; or
that it overlaps, naming a path the two share, the change in flight that
plans it, and that change's attempt.

#### Scenario: The same change again

- GIVEN a record names loose change A, whose pull request is open
- WHEN the operator runs `grove claim A`, on the same host or another
- THEN nothing is recorded, and Grove says A is in flight, naming the
  attempt, its host, its start time and its branch

#### Scenario: A file both name

- GIVEN A is in flight, A's T-01 names `ledger/cli.py`, and buildable
  change B's T-02 names `ledger/cli.py`
- WHEN the operator runs `grove claim B`
- THEN nothing is recorded, and Grove says B shares `ledger/cli.py` with
  A, naming A's attempt

#### Scenario: A folder one names

- GIVEN A is in flight, A's T-01 names `ledger/`, and B's T-01 names
  `ledger/report.py`
- WHEN the operator runs `grove claim B`
- THEN nothing is recorded, and Grove names `ledger/report.py`, or
  `ledger/`, as shared with A

#### Scenario: One capability, different requirements

- GIVEN A is in flight with a change spec for `record-editing` that
  MODIFIES "Keep a category", B has a change spec for `record-editing`
  that ADDS "Keep the tags", and their tasks name no file in common
- WHEN the operator runs `grove claim B`
- THEN nothing is recorded, and Grove says B shares the maintained spec
  of `record-editing` with A

#### Scenario: A cycle that overlaps a change in flight

- GIVEN loose change A is in flight naming `ledger/cli.py`, and buildable
  cycle K's member C names `ledger/cli.py`
- WHEN the operator runs `grove claim K`
- THEN nothing is recorded, and Grove says K shares `ledger/cli.py` with
  A

#### Scenario: A cycle with a member in flight

- GIVEN a record names cycle K with A and B, A is archived as finished on
  the target and B is not archived
- WHEN the operator runs `grove claim K`
- THEN nothing is recorded, and Grove says K is in flight, naming that
  attempt

#### Scenario: What an attempt changed on its branch

- GIVEN A is in flight, the target's A names only `ledger/cli.py`, the
  attempt revised A on its branch to name `ledger/report.py` too, and B's
  tasks name only `ledger/report.py`
- WHEN the operator runs `grove claim B`
- THEN an attempt is recorded for B

#### Scenario: The attempt it would overlap is merged

- GIVEN a record names A, the target holds A archived as finished, and
  B's tasks name a file A's named
- WHEN the operator runs `grove claim B`
- THEN an attempt is recorded for B

## MODIFIED Requirements

### Requirement: Start one attempt

When the operator runs `grove claim` in a clean checkout, and either
gives the ID of a buildable piece on the fetched target that is not in
flight and overlaps no piece in flight, or gives none while no change is
in flight and the fetched target holds an eligible cycle, the system
shall, when the checkout has a remote named `origin`, fetch the target
and `grove/runs` from it and fast-forward the local target branch to the
fetched target; record the attempt for the given piece, or else for the
oldest eligible cycle, the one whose `created` is earliest, as "Record
the attempt" says; create the attempt's branch `grove/<attempt>` from the
fetched target and switch the checkout to it; print the piece's ID,
alone, on standard output; and start no process but Git.

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

#### Scenario: A change in no cycle

- GIVEN `origin`'s target approves change C, which no cycle lists and
  which waits on nothing, and no change is in flight
- WHEN the operator runs `grove claim C`
- THEN standard output is C's ID and nothing else, the checkout is on a
  new branch from the fetched target, and `origin`'s `grove/runs` holds a
  record naming C alone and no cycle

#### Scenario: A change composed with a harness

- GIVEN buildable change C
- WHEN the operator runs
  `C=$(grove claim <C's ID>) && claude -p "/grove-implement $C"`
- THEN the session starts on `/grove-implement <C's ID>` in the checkout,
  on the attempt's branch

#### Scenario: A given cycle while another is in flight

- GIVEN a record names cycle L, whose members are not archived on the
  target, K is buildable, and K and L overlap nowhere
- WHEN the operator runs `grove claim K`
- THEN an attempt is recorded for K, and `grove/runs` holds two records
  in flight

#### Scenario: A second change that shares nothing

- GIVEN loose change A is in flight, and buildable change B's planned
  files are none of A's
- WHEN the operator runs `grove claim B` in another checkout
- THEN an attempt is recorded for B, and A's record and branch are as
  they were

#### Scenario: A change stale against its base

- GIVEN buildable change C, whose T-01 names `ledger/cli.py`, which the
  target changed after C's `base`, so that `grove check` reports C
- WHEN the operator runs `grove claim C`
- THEN an attempt is recorded for C

#### Scenario: A change with an open question

- GIVEN buildable change C, whose proposal's "Open questions" holds a
  question
- WHEN the operator runs `grove claim C`
- THEN an attempt is recorded for C

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

### Requirement: Record the attempt

When `grove claim` selects a piece, the system shall record the attempt
as one commit on the branch `grove/runs`, with an empty tree and the
record as its message: the attempt's ID, its start time, the selected
cycle and its changes in build order or, for a loose change, that change
alone and no cycle, the target commit, the host's name, the checkout and
the attempt's branch; where the checkout has `origin`, push it so that
the push succeeds only while `origin`'s `grove/runs` is at the tip just
fetched, and otherwise update the local `grove/runs` only while it is at
the tip just read; and write nothing under `.git/grove/`. The record
states that an attempt started, not how any session ended nor whether any
task is done.

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

#### Scenario: A loose change's record

- GIVEN `grove claim C` recorded an attempt for loose change C
- WHEN the operator runs `git log -1 --format=%B origin/grove/runs`
- THEN the record lists C as its only change and holds no `cycle`

#### Scenario: Nothing in the working tree

- WHEN `grove claim` records an attempt
- THEN `git status` shows nothing of it, and the attempt's branch holds
  only the fetched target's commits

### Requirement: Pause after an attempt

While any change is in flight, on the fetched `grove/runs`, or on the
local `grove/runs` in a checkout with no `origin`, `grove claim` given no
piece and without `--retry` shall record nothing, create no branch, print
nothing on standard output, exit non-zero and say the project is paused,
naming the newest record in flight: the attempt, the host that started
it, when, and its branch.

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

#### Scenario: An older attempt still open

- GIVEN records name cycle K and then loose change C, the target holds C
  archived as finished, K's members are not archived, and cycle L is
  eligible
- WHEN `grove claim` is invoked
- THEN nothing is recorded, and Grove says the project is paused, naming
  K's attempt

#### Scenario: A change claimed by name

- GIVEN the only record names loose change C, not archived on the target,
  and cycle L is eligible
- WHEN `grove claim` is invoked
- THEN nothing is recorded, and Grove says the project is paused after
  C's attempt

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
system shall record a new attempt as `grove claim` does when no change is
in flight, whatever is in flight and whatever overlaps it, refusing still
a given piece that cannot be built, and leave earlier records and their
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

#### Scenario: The same change again

- GIVEN a record names loose change A from host H1, whose session died
- WHEN the operator runs `grove claim A --retry`
- THEN a new record names A, it is A's attempt, H1's record is no longer
  in flight, and H1's branch is untouched

#### Scenario: Despite an overlap

- GIVEN A is in flight naming `ledger/cli.py`, and buildable change B
  names `ledger/cli.py`
- WHEN the operator runs `grove claim B --retry`
- THEN an attempt is recorded for B, and A is still in flight

#### Scenario: A change that waits

- GIVEN loose change B depends on A, which is not archived
- WHEN the operator runs `grove claim B --retry`
- THEN nothing is recorded, and Grove says B waits on A

### Requirement: One claim wins

When two `grove claim` commands record an attempt for one project at once,
the system shall keep the record whose ref update `origin`, or the local
repository where there is no `origin`, accepts first, and have the other,
whose update is refused because `grove/runs` moved, read `grove/runs`
again and decide again as a fresh `grove claim` would; a push refused for
any other reason shall be reported in one sentence quoting Git, and record
nothing.

#### Scenario: Two hosts at once

- GIVEN checkouts of one project on hosts H1 and H2, no change in flight,
  and eligible cycle K
- WHEN both run `grove claim` within the same second
- THEN `origin`'s `grove/runs` holds one new record, for K; one host
  prints K's ID and is on its attempt's branch; and the other prints
  nothing on standard output, says the project is paused after that
  attempt, naming the first host, and creates no branch

#### Scenario: Two changes that share nothing

- GIVEN buildable changes A and B that do not overlap, and no change in
  flight
- WHEN H1 runs `grove claim A` and H2 runs `grove claim B` within the
  same second
- THEN `origin`'s `grove/runs` holds two new records, one for A and one
  for B, each host prints its change's ID, and each is on its attempt's
  branch

#### Scenario: The same change from two hosts

- GIVEN buildable change A, not in flight
- WHEN H1 and H2 both run `grove claim A` within the same second
- THEN `origin`'s `grove/runs` holds one new record; one host prints A's
  ID; and the other prints nothing on standard output, says A is in
  flight, naming the first host, and creates no branch

#### Scenario: Two changes that overlap

- GIVEN buildable changes A and B, both naming `ledger/cli.py`, and no
  change in flight
- WHEN H1 runs `grove claim A` and H2 runs `grove claim B` within the
  same second
- THEN `origin`'s `grove/runs` holds one new record, and the host whose
  update was refused says its change shares `ledger/cli.py` with the
  other and creates no branch

#### Scenario: A push refused for permissions

- GIVEN the checkout's account may not push to `origin`
- WHEN the operator runs `grove claim` with eligible cycle K
- THEN Grove says the push was refused, quoting Git's reason, records
  nothing on `origin`, leaves the local `grove/runs` as `origin` has it,
  and creates no branch

## REMOVED Requirements

### Requirement: Refuse a given cycle that cannot be built

"Refuse a given piece that cannot be built" replaces it, keeping its
scenarios and adding those for a change given by its ID.

## Observed current behavior

At 07d319f, `grove claim <id>` looks the ID up among cycles only and says
the target "holds no cycle" for a change's ID, and the pause reads the
changes of the latest record alone (`internal/cli/claim.go`, `claim` and
`pausing`).
