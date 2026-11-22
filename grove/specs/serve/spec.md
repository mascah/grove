# Serve

## Purpose

`grove serve` shows where a project stands on one local page, read from
the checkout: what `grove claim` would do with each cycle and each change
in no cycle, and who has claimed one, the work in progress on branches and in the working tree, the runs
recorded on the branch `grove/runs`, the
dependency graph, the target's history of the record and the maintained
specs, each document readable beside them. "The target" is the local
target branch `grove.yaml` names. "A blob the target's history never
held" is a file's content that no commit on the target's history held
under the document root, which tells new work from work a squash merge
already landed. "In flight", "a loose change", "a buildable piece" and
"overlap" are as the [runs](../runs/spec.md) spec defines them, read here
from the target the checkout holds and the `grove/runs` it has fetched.

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
at that commit, Markdown rendered to HTML with raw HTML dropped, the
fields a Markdown file starts with between two lines of `---` as text
above the rest, and YAML as text, and refuse any other path or ref.

#### Scenario: A proposal at the target

- GIVEN the target holds change A
- WHEN the page asks for A's `proposal.md` with an empty ref
- THEN the answer is HTML with A's title as its first heading

#### Scenario: A cycle file

- GIVEN cycle K's file holds its fields and, below them, notes with a
  heading and a list
- WHEN the page asks for K's file
- THEN the answer is HTML holding K's fields as text, each on its own
  line as the file has them, and below them the notes rendered, the
  heading as a heading and the list as a list

#### Scenario: A cycle file with no notes

- GIVEN cycle K's file holds only its fields
- WHEN the page asks for K's file
- THEN the answer is HTML holding K's fields as text and nothing below

#### Scenario: A change's metadata

- GIVEN the target holds change A
- WHEN the page asks for A's `.grove.yaml`
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

### Requirement: What is undecided and blocked

When the page asks for `/api/project`, the system shall give each change
`open_questions` as `grove changes --json` does (changes spec, "Show what
is undecided and blocked"), and each change of a branch row the same,
read from the change's main document at that branch's tip or in the
working tree; and the page shall show, in the status column of a cycle's
member row, `open questions` for a change that has them and its tasks
blocked on the target, and, on a branch card, `open questions` beside the
change's task chips.

#### Scenario: A member with questions and a blocked task

- GIVEN cycle K on the target lists change C, whose proposal has an open
  question and whose `implementation.md` on the target gives T-02
  `Status: blocked`
- WHEN the owner opens the page
- THEN C's row under K shows `open questions` and `T-02 blocked`

#### Scenario: A draft with a question

- GIVEN branch `draft-notes` adds change D, which the target does not
  hold, and D's proposal has an open question
- WHEN the page asks for `/api/project`
- THEN the row naming `draft-notes` gives D `open_questions` with the
  question's text, and its card shows `open questions` beside D

#### Scenario: Blocked only on a branch

- GIVEN the target approves B with no `implementation.md`, and branch
  `build-b` holds B's `implementation.md` with T-02 `Status: blocked`
- WHEN the owner opens the page
- THEN `build-b`'s card shows T-02 blocked, and B's row under its cycle
  shows no blocked task

## Known discrepancies

None recorded.

## Open questions

None.

## History

- [261107-9z9tg](../../changes/archive/261107-9z9tg-serve-the-project-as-a-page/) Serve the project's state as a local page
- [261108-jpg8n](../../changes/archive/261108-jpg8n-claim-replaces-run/) grove claim replaces grove run, with attempt records on a branch
- [261108-7grt3](../../changes/archive/261108-7grt3-the-project-as-json/) The project as JSON
- [261108-cns75](../../changes/archive/261108-cns75-show-what-is-undecided-and-blocked/) Show what is undecided and blocked
- [261114-yx2b7](../../changes/archive/261114-yx2b7-remove-needs-human/) Remove needs_human
- [261118-emse7](../../changes/archive/261118-emse7-a-cycle-is-a-markdown-file/) A cycle is a Markdown file with plan's notes
- [261122-7vht2](../../changes/archive/261122-7vht2-claim-a-change-by-name/) Claim a change by name, and more than one piece of work at once
