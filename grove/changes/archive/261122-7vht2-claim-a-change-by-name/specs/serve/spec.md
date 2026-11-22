# Changes to serve

The changes that change 261122-7vht2 proposes to the requirements of
[serve](../../../../specs/serve/spec.md), as it stands at 07d319f. This
is future work: nothing here describes the target branch until the change
is implemented and the maintained spec is reconciled.

The Purpose's "what `grove claim` would do with each cycle, and who has
claimed one" becomes "what `grove claim` would do with each cycle and
each change in no cycle, and who has claimed one". "In flight", "a loose
change", "a buildable piece" and "overlap" are as the
[runs](../../../../specs/runs/spec.md) spec defines them, read here from
the target the checkout holds and the `grove/runs` it has fetched.

## ADDED Requirements

### Requirement: Readiness of each change in no cycle

When the page asks for `/api/project`, the system shall give every loose
change the verdict `grove claim` given its ID would give it: the change
it waits on that is not archived as finished; otherwise, when it is in
flight, that it is, with the attempt, the host, the start time and the
attempt's branch; otherwise the path it shares with a change in flight,
naming that change and its attempt; otherwise that `grove claim <its ID>`
would claim it; and `claimable`, true only for the last.

#### Scenario: Nothing in the way

- GIVEN the target approves C, in no cycle, waiting on nothing, and no
  change is in flight
- WHEN the page asks for `/api/project`
- THEN C's verdict is that `grove claim <C's ID>` would claim it, and
  its `claimable` is true

#### Scenario: A change that waits

- GIVEN loose change B depends on A, which is not archived
- WHEN the page asks for `/api/project`
- THEN B's verdict names A as what it waits on, and its `claimable` is
  false

#### Scenario: A change in flight

- GIVEN a record on `grove/runs` names loose change C, from host H, and
  the target does not hold C archived
- WHEN the page asks for `/api/project`
- THEN C's verdict says it is in flight, naming H, when, and the
  attempt's branch, and its `claimable` is false

#### Scenario: A change that overlaps one in flight

- GIVEN loose change A is in flight, and loose change B, which waits on
  nothing, names `ledger/cli.py`, as A does
- WHEN the page asks for `/api/project`
- THEN B's verdict says it shares `ledger/cli.py` with A, naming A's
  attempt, and its `claimable` is false

#### Scenario: A member of a cycle

- GIVEN cycle K, not done, lists approved change A
- WHEN the page asks for `/api/project`
- THEN A has no verdict and no `claimable` of its own

## MODIFIED Requirements

### Requirement: The project value

When the page asks for `/api/project`, the system shall read the target
branch's documents and answer with the one JSON object `grove changes
--json` prints (changes spec, "Print the project as data"), holding the
root, the target, its full commit, the time of the read, the maintained
specs, every change as `grove changes --json` gives it, each with the
cycle not done that lists it, whether it is approved and waits on
nothing and, for a loose change, what `grove claim` given its ID would
do, every cycle, the branches with work, the recorded attempts, the
history and the problems `grove check` would report on the target's
documents, as the requirements below say.

#### Scenario: An approved change in no cycle

- GIVEN the target approves change C, listed by no cycle, with
  `depends_on: [A]` and A archived as finished
- WHEN the page asks for `/api/project`
- THEN C appears with no cycle, as ready, and with a verdict

#### Scenario: The same value as the command

- GIVEN a checkout of a Grove project with `grove serve` running
- WHEN the page asks for `/api/project` and the owner runs
  `grove changes --json` in the checkout
- THEN the two objects are equal but for the time of the read

### Requirement: Readiness of each cycle

When the page asks for `/api/project`, the system shall give every cycle
on the target with its members in build order, the members archived, and
the verdict `grove claim` given its ID would give it: done when every
member is archived; otherwise the cycle it waits on that is not done;
otherwise the member that waits on a change neither in the cycle nor
archived as finished; otherwise the members that changed on the target
since its plan, making it stale; otherwise, when it is in flight, that it
is, with the attempt, the host, the start time and the attempt's branch;
otherwise the path it shares with a change in flight, naming that change
and its attempt; otherwise that `grove claim <cycle-id>` would claim it,
saying when it is not marked `unattended` and when a `grove claim` given
no piece is paused after another piece's attempt; with `claimable`, true
only for the last, and `runs`, true only when it is claimable and no
change is in flight.

#### Scenario: A marked cycle that is fresh

- GIVEN the target holds cycle K, marked `unattended`, listing approved
  changes A and B, B with `depends_on: [A]`, neither archived, nothing of
  theirs changed since the plan, and no change is in flight
- WHEN the page asks for `/api/project`
- THEN K's members are A then B, its verdict is that `grove claim` would
  claim it, and its `claimable` and `runs` are true

#### Scenario: A member waits outside the cycle

- GIVEN cycle L's member C depends on change Z, neither in L nor archived
  as finished
- WHEN the page asks for `/api/project`
- THEN L's verdict names C as waiting on Z

#### Scenario: A member changed since the plan

- GIVEN cycle M lists member D, and a commit on the target after M's
  `planned` revised D's proposal without editing M's file
- WHEN the page asks for `/api/project`
- THEN M's verdict says it is stale and names D

#### Scenario: A cycle in flight

- GIVEN a record on `grove/runs` names cycle K, from host H, and K's
  members are not archived on the target
- WHEN the page asks for `/api/project`
- THEN K's verdict says it is in flight, naming H, when, and the
  attempt's branch, and its `claimable` and `runs` are false

#### Scenario: Another piece in flight that shares nothing

- GIVEN loose change A is in flight, and cycle K, marked `unattended`
  and buildable, overlaps A nowhere
- WHEN the page asks for `/api/project`
- THEN K's verdict is that `grove claim <K's ID>` would claim it, saying
  a `grove claim` given no piece is paused after A's attempt; its
  `claimable` is true and its `runs` is false

#### Scenario: Another piece in flight that shares a file

- GIVEN loose change A is in flight naming `ledger/cli.py`, and buildable
  cycle K's member C names `ledger/cli.py`
- WHEN the page asks for `/api/project`
- THEN K's verdict says it shares `ledger/cli.py` with A, naming A's
  attempt, and its `claimable` and `runs` are false

### Requirement: Recorded runs

When the page asks for `/api/project`, the system shall list each attempt
record on `origin/grove/runs`, or on the local `grove/runs` when the
checkout has no `origin`, newest first, each with its attempt, its cycle
where it has one, its changes, start time, host and branch, `in_flight`,
whether the record is in flight, and, where the checkout holds the
attempt's branch locally or on `origin`, that branch's latest commit
time.

#### Scenario: An attempt from another host

- GIVEN `origin/grove/runs`, fetched into the checkout, holds a record from
  host H, and `origin` holds the attempt's branch with a commit after the
  start
- WHEN the page asks for `/api/project`
- THEN the runs hold that record, with its cycle, branch, start time and
  H, and the branch's latest commit time

#### Scenario: Two attempts, one merged

- GIVEN records name loose changes A and then B, the target holds A
  archived as finished, and B is not archived
- WHEN the page asks for `/api/project`
- THEN the runs hold B's record with `in_flight` true and no `cycle`, and
  A's with `in_flight` false

#### Scenario: An attempt retried

- GIVEN two records name loose change A, which is not archived
- WHEN the page asks for `/api/project`
- THEN the newer record's `in_flight` is true and the older's is false

#### Scenario: No record anywhere

- GIVEN neither `origin` nor the checkout has a branch `grove/runs`
- WHEN the page asks for `/api/project`
- THEN the runs are empty

### Requirement: The page

When the owner opens the URL, the system shall show one page with tabs
for readiness, in progress, graph, history and specs, carrying counts
where they help, a pane that shows a document when a change, cycle, spec
or file is clicked, and a refresh of the value every 10 seconds that a
checkbox stops; readiness listing each cycle not done with its verdict
and its members' status and where each is being built, and each loose
change with its verdict; in progress listing the branches with work,
each change's task progress and new files, and the recorded attempts
with their hosts, marking those in flight; the graph drawing each cycle
as a band with arrows for `depends_on` between changes and between
cycles; history drawing a timeline of the marks and listing the commits;
and the pane keeping the loaded document and both columns' scroll
positions across a refresh.

#### Scenario: Readiness at a glance

- GIVEN the fixture with cycle K fresh and marked, and cycle L with a
  member that waits on a change outside it
- WHEN the owner opens the URL
- THEN the readiness tab shows K as a cycle `grove claim` would claim and
  L with the waiting member named, and the tab's count is 2

#### Scenario: A change in no cycle

- GIVEN that fixture, and approved change C in no cycle, waiting on
  nothing
- WHEN the owner opens the URL
- THEN the readiness tab also shows C as a change `grove claim <C's ID>`
  would claim, and the tab's count is 3

#### Scenario: Attempts in flight

- GIVEN records for loose changes A, archived as finished on the target,
  and B, not archived
- WHEN the owner opens the in progress tab
- THEN B's attempt is marked in flight and A's is not

#### Scenario: A refresh while reading

- GIVEN the owner has opened a change's proposal in the pane and scrolled
  it
- WHEN 10 seconds pass
- THEN the pane still shows the proposal at the same place
