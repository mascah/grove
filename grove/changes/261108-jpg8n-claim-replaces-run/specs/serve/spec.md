# Changes to serve

The changes that change 261108-jpg8n proposes to the requirements of
[serve](../../../../specs/serve/spec.md). This is future work.

The Purpose's "what `grove run` would do with each cycle" becomes "what
`grove claim` would do with each cycle, and who has claimed one", and
"the runs recorded" reads the branch `grove/runs`. The Boundaries name
the runs spec for the verdict as before.

## MODIFIED Requirements

### Requirement: Readiness of each cycle

When the page asks for `/api/project`, the system shall give every cycle
on the target with its members in build order, the members archived, and
the verdict `grove claim` would give it: done when every member is
archived; otherwise the member that needs a person, as `needs_human`
names it; otherwise the cycle it waits on that is not done; otherwise the
member that waits on a change neither in the cycle nor archived as
finished; otherwise the members that changed on the target since its
plan, making it stale; otherwise, when the latest record on `grove/runs`
names it, that it is claimed, with the host, the start time and the
attempt's branch; otherwise that `grove claim` would claim it, saying
when it is not marked `unattended` and when the project is paused after
another cycle's attempt.

#### Scenario: A marked cycle that is fresh

- GIVEN the target holds cycle K, marked `unattended`, listing approved
  changes A and B, B with `depends_on: [A]`, neither archived, nothing of
  theirs changed since the plan, and no record on `grove/runs` names a
  cycle not done
- WHEN the page asks for `/api/project`
- THEN K's members are A then B and its verdict is that `grove claim`
  would claim it

#### Scenario: A member needs a person

- GIVEN cycle L's `needs_human` names member C
- WHEN the page asks for `/api/project`
- THEN L's verdict names C as needing a person

#### Scenario: A member changed since the plan

- GIVEN cycle M lists member D, and a commit on the target after M's
  `planned` revised D's proposal without editing M's file
- WHEN the page asks for `/api/project`
- THEN M's verdict says it is stale and names D

#### Scenario: A claimed cycle

- GIVEN the latest record on `grove/runs` names cycle K, from host H, and
  K's members are not archived on the target
- WHEN the page asks for `/api/project`
- THEN K's verdict says it is claimed, naming H, when, and the attempt's
  branch

### Requirement: Recorded runs

When the page asks for `/api/project`, the system shall list each attempt
record on `origin/grove/runs`, or on the local `grove/runs` when the
checkout has no `origin`, newest first, each with its attempt, cycle,
start time, host and branch and, where the checkout holds the attempt's
branch locally or on `origin`, that branch's latest commit time.

#### Scenario: An attempt from another host

- GIVEN `origin/grove/runs`, fetched into the checkout, holds a record from
  host H, and `origin` holds the attempt's branch with a commit after the
  start
- WHEN the page asks for `/api/project`
- THEN the runs hold that record, with its cycle, branch, start time and
  H, and the branch's latest commit time

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
and its members' status and where each is being built; in progress
listing the branches with work, each change's task progress and new
files, and the recorded attempts with their hosts; the graph drawing each
cycle as a band with arrows for `depends_on` between changes and between
cycles; history drawing a timeline of the marks and listing the commits;
and the pane keeping the loaded document and both columns' scroll
positions across a refresh.

#### Scenario: Readiness at a glance

- GIVEN the fixture with cycle K fresh and marked, and cycle L with a
  member that needs a person
- WHEN the owner opens the URL
- THEN the readiness tab shows K as a cycle `grove claim` would claim and
  L with the member that needs a person named, and the tab's count is 2

#### Scenario: A refresh while reading

- GIVEN the owner has opened a change's proposal in the pane and scrolled
  it
- WHEN 10 seconds pass
- THEN the pane still shows the proposal at the same place
