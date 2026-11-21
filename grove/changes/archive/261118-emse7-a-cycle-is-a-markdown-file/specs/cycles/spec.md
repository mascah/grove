# Changes to cycles

The changes that change 261118-emse7 proposes to the requirements of
[cycles](../../../../specs/cycles/spec.md), as it stands at f7e0140. This
is future work: nothing here describes the target branch until the change
is implemented and the maintained spec is reconciled.

In the Purpose, "A cycle is a file `<root>/cycles/<id>-<slug>.yaml`, and
its members are the changes its `changes` lists" becomes "A cycle is a
file `<root>/cycles/<id>-<slug>.md`: its fields, as YAML between two lines
of `---`, and below them the notes its plan wrote. Its members are the
changes its `changes` lists."

## ADDED Requirements

### Requirement: Report a cycle file left as YAML

When the owner runs `grove check` and `<root>/cycles/` holds a file whose
name ends `.yaml`, the system shall report that file, saying a cycle file
is Markdown now and to rename it to `.md` with its fields between two
lines of `---`, and exit non-zero, and shall read no cycle from it.

#### Scenario: An open cycle left as YAML

- GIVEN `grove/cycles/261020-m4k8t-0-4-release.yaml` lists approved change
  A and holds `unattended: true`
- WHEN the owner runs `grove check`, `grove changes` and `grove claim`
- THEN `grove check` reports the file, saying a cycle file is Markdown now
  and to rename it to `261020-m4k8t-0-4-release.md` and put its fields
  between two lines of `---`, and exits non-zero; `grove changes` shows A
  with no "in cycle"; and `grove claim` claims no cycle

#### Scenario: A done cycle left as YAML

- GIVEN `grove/cycles/261019-k3m9q-0-3-release.yaml` lists only changes
  the target holds archived
- WHEN the owner runs `grove check`
- THEN it reports that file the same way

#### Scenario: Converted

- GIVEN that open cycle's file renamed to
  `grove/cycles/261020-m4k8t-0-4-release.md`, a line of `---` added above
  its fields and another below, and nothing else changed
- WHEN the commit holding it is on the target and the owner runs
  `grove check` and `grove claim`
- THEN `grove check` reports nothing for the cycle, `planned` is what it
  was, and `grove claim` claims it

## MODIFIED Requirements

### Requirement: Plan a cycle

When the owner runs `/grove-plan <description>`, the system shall print
instructions for writing a new cycle file holding, as its fields, a
`title` the agent writes from its members without asking for one, an
optional one-line `goal`, the approved changes the owner agrees to as
`changes`, the time Grove gives the plan as `created`, and the target
commit the plan read as `planned`; and, below its fields and only where
there is something to say, notes on what the plan found on the target at
`planned`, what whoever builds the cycle should know above any one
change, and why these changes are planned together in this order and what
was left out, none of them changing a member's requirements or tasks;
editing no change's files except, for each member `grove check` reports
stale on a requirement or for a named file changed since its base,
telling the owner that this edits approved work and bringing that member
up to date on the same branch as the propose instructions for it say,
before writing the cycle file.

#### Scenario: Two approved changes

- GIVEN the target approves changes A and B, neither archived, and
  `grove check` reports neither
- WHEN the owner runs `/grove-plan 0.4 release` and agrees to A and B
- THEN a new cycle file, `grove/cycles/<cycle-id>-<slug>.md`, lists A and
  B under `changes` between two lines of `---`, `created` is the time
  Grove gave the plan, `planned` is the target's commit, and A's and B's
  files are unchanged

#### Scenario: The agent titles the cycle

- GIVEN the target approves changes A and B
- WHEN the owner runs `/grove-plan 0.4 release` and agrees to A and B
- THEN the agent writes the cycle's `title` itself and asks the owner
  nothing about it

#### Scenario: A draft change named

- GIVEN change D exists in the working tree but not in the target
- WHEN the owner asks for D in the cycle
- THEN D is not listed, and the agent says D is not approved

#### Scenario: A member stale on its base

- GIVEN the target approves change B, whose T-01 names `ledger/report.py`,
  and the merge that archived B's dependency A changed `ledger/report.py`
  after B's `base`
- WHEN the owner runs `/grove-plan 0.4 release` and agrees to B
- THEN the agent says bringing B up to date edits approved work, brings
  T-01 up to date with `ledger/report.py` as the target has it, sets B's
  `base` to the target's commit and writes the cycle file, all on the same
  branch, without stopping for another session, and `grove check` is clean

#### Scenario: What the plan found is noted

- GIVEN that plan, in which B was brought up to date
- WHEN the agent writes the cycle file
- THEN the file's notes, below its fields, say B was brought up to date
  with `ledger/report.py` as the target has it at `planned`, and what in
  T-01 moved

#### Scenario: A change left out is noted

- GIVEN the target approves changes A, B and C, and the owner leaves C out
  because it waits on a login the runner does not have
- WHEN the agent writes the cycle file
- THEN its notes say C was left out and why

#### Scenario: Nothing to note

- GIVEN the target approves changes A and B, `grove check` reports
  neither, and the plan found nothing to say beyond its fields
- WHEN the agent writes the cycle file
- THEN the file holds its fields between two lines of `---` and no notes

### Requirement: Revise a cycle

When the owner runs `/grove-plan <cycle-id>`, the system shall print
instructions for revising that cycle file as the owner asks, setting
`planned` to the target commit the revision read, and bringing its notes
up to date with that commit.

#### Scenario: Re-plan a stale cycle

- GIVEN cycle K is stale because member A was revised on the target
- WHEN the owner runs `/grove-plan K`
- THEN the agent reads A as the target holds it, asks about it as for a new
  plan, and sets `planned` to the target's commit

#### Scenario: Re-plan after a member's file moved

- GIVEN cycle K is stale because a file member B's T-01 names changed on
  the target, and the owner is content to build B on the new text
- WHEN the owner runs `/grove-plan K` and merges the cycle file
- THEN `planned` is the target's commit, K is no longer stale, and B's
  files are unchanged

#### Scenario: A note the target made false

- GIVEN cycle K's notes say every file B's T-01 names is as it was at B's
  base, and K is stale because the target changed one of them
- WHEN the owner runs `/grove-plan K`
- THEN the notes say what holds at the new `planned` commit, and no longer
  what the target made false

### Requirement: Report a malformed cycle

When the owner runs `grove check`, the system shall report a cycle file
that does not start with its fields between two lines of `---`, with no
`title`, no `created`, no `changes` or no `planned`, with an unknown
field, that lists a change that is no change, open or archived, or that
lists a change another cycle not yet done also lists; and, for a cycle
not yet done, a `depends_on` that names no cycle, or a cycle that depends
on it in turn, directly or through other cycles; each problem in a field
by that field's line in the file.

#### Scenario: One change in two cycles

- GIVEN cycles K and L, neither done, both list change A
- WHEN the owner runs `grove check`
- THEN it reports the later of the two, naming A and the other cycle

#### Scenario: A needs_human left over

- GIVEN cycle K's file still holds `needs_human` naming member C
- WHEN the owner runs `grove check`
- THEN it reports K's `needs_human` line as an unknown field, listing the
  fields a cycle file takes

#### Scenario: A dependency on no cycle

- GIVEN cycle K, not done, whose `depends_on` names 261099-zzzzz, and no
  cycle file has that ID
- WHEN the owner runs `grove check`
- THEN it reports K's `depends_on` line, saying 261099-zzzzz is no cycle

#### Scenario: Cycles that wait on each other

- GIVEN cycle K depends on cycle L, and L depends on K, neither done
- WHEN the owner runs `grove check`
- THEN it reports each one's `depends_on` line, saying the other waits on
  it in turn, so neither can run

#### Scenario: No fields at the top

- GIVEN `grove/cycles/261020-m4k8t-0-4-release.md` starts with a heading
  and holds no line of `---`
- WHEN the owner runs `grove check`
- THEN it reports the file, saying a cycle file starts with its fields
  between two lines of `---`, and exits non-zero

#### Scenario: The line is the file's

- GIVEN cycle K's file holds `---` on its first line and `needs_human` on
  its sixth
- WHEN the owner runs `grove check`
- THEN the report names line 6

#### Scenario: Notes are not checked

- GIVEN cycle K's file holds valid fields and, below them, notes that
  contain a line `needs_human: [C]`
- WHEN the owner runs `grove check`
- THEN it reports nothing for K

### Requirement: Report a cycle file named without an ID

When the owner runs `grove check` and a file under `<root>/cycles/` has a
name that does not start with a cycle ID, the system shall report the
file, saying its name must start with a cycle ID, and exit non-zero, and
`grove changes` shall name no cycle for the changes that file lists.

#### Scenario: A letter IDs leave out

- GIVEN `grove/cycles/261105-proto-unknown-id.md` lists approved change
  C, and `o` is a letter cycle IDs leave out
- WHEN the owner runs `grove check` and `grove changes`
- THEN `grove check` reports the file, saying the file name must start
  with a cycle ID, and exits non-zero, and `grove changes` shows C with
  no "in cycle"
