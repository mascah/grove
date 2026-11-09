# Serve

## Purpose

`grove serve` shows where a project stands on one local page, read from
the checkout: what `grove claim` would do with each cycle, and who has claimed
one, the work in progress on branches and in the working tree, the runs
recorded on the branch `grove/runs`, the
dependency graph, the target's history of the record and the maintained
specs, each document readable beside them. "The target" is the local
target branch `grove.yaml` names. "A blob the target's history never
held" is a file's content that no commit on the target's history held
under the document root, which tells new work from work a squash merge
already landed.

## Boundaries

This covers `grove serve`, its `--addr`, the page it serves and the value
and documents the page reads. Readiness repeats the verdict
`grove claim` gives (`grove/specs/runs/spec.md`) and the problems
`grove check` reports; the value is the one `grove changes --json` prints
(`grove/specs/changes/spec.md`). The
server listens only on loopback, holds nothing and writes nothing
(`CONSTITUTION.md` rule 7); fetching from `origin`, editing a document,
anything hosted and any authentication are outside it.

## Requirements

### Requirement: Start a local server

When the owner runs `grove serve` in a checkout of a Grove project, the
system shall listen on `127.0.0.1:7070`, or on the address `--addr`
gives, print the page's URL, serve the page at `/` and its data under
`/api/`, and run until it is interrupted, fetching from no remote and
writing no file.

#### Scenario: The default address

- GIVEN a checkout of a Grove project
- WHEN the owner runs `grove serve`
- THEN Grove prints `http://127.0.0.1:7070/`, and a request for `/`
  returns the page

#### Scenario: An address given

- WHEN the owner runs `grove serve --addr 127.0.0.1:8080`
- THEN Grove prints `http://127.0.0.1:8080/` and listens there

#### Scenario: The address is taken

- GIVEN another process listens on `127.0.0.1:7070`
- WHEN the owner runs `grove serve`
- THEN Grove says the address is in use and exits non-zero

### Requirement: The project value

When the page asks for `/api/project`, the system shall read the target
branch's documents and answer with the one JSON object `grove changes
--json` prints (changes spec, "Print the project as data"), holding the
root, the target, its full commit, the time of the read, the maintained
specs, every change as `grove changes --json` gives it, each with the
cycle not done that lists it and whether it is approved and waits on
nothing, every cycle, the branches with work, the recorded attempts, the
history and the problems `grove check` would report on the target's
documents, as the requirements below say.

#### Scenario: An approved change in no cycle

- GIVEN the target approves change C, listed by no cycle, with
  `depends_on: [A]` and A archived as finished
- WHEN the page asks for `/api/project`
- THEN C appears with no cycle and as ready

#### Scenario: The same value as the command

- GIVEN a checkout of a Grove project with `grove serve` running
- WHEN the page asks for `/api/project` and the owner runs
  `grove changes --json` in the checkout
- THEN the two objects are equal but for the time of the read

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

### Requirement: Work in progress by branch

When the page asks for `/api/project`, the system shall list, for every
local branch and every `origin` branch the target does not hold merged,
and for the working tree, each change whose files under the root hold a
blob the target's history never held, as a draft when the target does not
hold the change and as in progress when the target approves it, with each
task's status from that branch's `implementation.md` and the files that
are new there, leaving out a branch whose such changes are all archived on
the target, giving a local branch and its origin twin at the same tip as
one row that says origin holds it, and sorting the rows newest first.

#### Scenario: A branch squash-merged and overtaken

- GIVEN branch `apply/A` built change A and edited a maintained spec, its
  commits were squash-merged into the target, A was archived there, and
  the target then edited the same spec again
- WHEN the page asks for `/api/project`
- THEN no row names `apply/A`

#### Scenario: A build in progress

- GIVEN branch `build-B` from the target holds B's `implementation.md`
  with T-01 `Status: done` and T-02 `Status: blocked`, and the target
  approves B
- WHEN the page asks for `/api/project`
- THEN a row names `build-B` with B in progress, T-01 done, T-02 blocked,
  and `implementation.md` as new there

#### Scenario: A draft

- GIVEN branch `draft-notes` adds change D, which the target does not
  hold
- WHEN the page asks for `/api/project`
- THEN a row names `draft-notes` with D as a draft and its documents as
  new there

#### Scenario: The working tree

- GIVEN the working tree holds an uncommitted edit to C's
  `implementation.md`, and the target approves C
- WHEN the page asks for `/api/project`
- THEN a row for the working tree names C in progress

#### Scenario: A teammate's branch

- GIVEN `origin/teammate` holds B's `implementation.md` with T-01
  `Status: to check`, and no local branch has that tip
- WHEN the page asks for `/api/project`
- THEN one row names `origin/teammate`, says only origin holds it, and
  gives T-01 as `to check`

#### Scenario: A pushed branch

- GIVEN local branch `build-B` and `origin/build-B` share a tip
- WHEN the page asks for `/api/project`
- THEN one row names `build-B` and says origin holds it

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

### Requirement: History of the record

When the page asks for `/api/project`, the system shall list, newest
first, each commit on the target's first-parent history that approved a
change (added its `.grove.yaml` under `changes/`), revised one (changed
another of its files there), archived one (added its `.grove.yaml` under
`changes/archive/`), planned a cycle (added its file), re-planned one
(changed it) or edited a maintained spec, giving the commit, its date, its
subject and each such mark with the change, cycle or spec it names.

#### Scenario: A change's life

- GIVEN the target's history holds a commit that added change A, a later
  one that changed A's `proposal.md`, a later one that added cycle K
  listing A, and a later one that moved A to `changes/archive/`
- WHEN the page asks for `/api/project`
- THEN those commits appear, newest first, marked archived A, planned K,
  revised A and approved A

#### Scenario: A commit outside the record

- GIVEN a commit on the target that changed only `ledger/records.py`
- WHEN the page asks for `/api/project`
- THEN it does not appear in the history

### Requirement: A document rendered

When the page asks for `/api/doc` with a path under the root ending
`.md` or `.yaml` and a ref that is empty, `worktree` or a full commit, the
system shall answer with that file at the target, in the working tree or
at that commit, Markdown rendered to HTML with raw HTML dropped and YAML
as text, and refuse any other path or ref.

#### Scenario: A proposal at the target

- GIVEN the target holds change A
- WHEN the page asks for A's `proposal.md` with an empty ref
- THEN the answer is HTML with A's title as its first heading

#### Scenario: A cycle file

- WHEN the page asks for cycle K's file
- THEN the answer is the file's text

#### Scenario: A document on a branch

- GIVEN branch `build-B` holds B's `implementation.md` as the target does
  not
- WHEN the page asks for it with the branch tip's commit as the ref
- THEN the answer is that branch's version

#### Scenario: Raw HTML in a document

- GIVEN a proposal that contains a `<script>` element
- WHEN the page asks for it
- THEN the answer holds no `<script>` element

#### Scenario: A path outside the root

- WHEN the page asks for `../grove.yaml`
- THEN the answer is a refusal, and no file is read

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

## Known discrepancies

None recorded.

## Open questions

None.

## History

- [261107-9z9tg](../../changes/archive/261107-9z9tg-serve-the-project-as-a-page/) Serve the project's state as a local page
- [261108-jpg8n](../../changes/archive/261108-jpg8n-claim-replaces-run/) grove claim replaces grove run, with attempt records on a branch
- [261108-7grt3](../../changes/archive/261108-7grt3-the-project-as-json/) The project as JSON
