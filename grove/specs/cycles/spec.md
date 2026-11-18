# Cycles

## Purpose

A cycle lists approved changes to build together: `/grove-plan` writes it,
and `grove claim <cycle-id>` claims it for the agent its caller starts, as
does a `grove claim` given no cycle once it holds `unattended: true`. A cycle is
a file `<root>/cycles/<id>-<slug>.md`: its fields, as YAML between two lines
of `---`, and below them the notes its plan wrote. Its members are the
changes its `changes` lists. A cycle is done when the target holds every member
archived, finished or dropped. A cycle "depends on" each cycle its optional
`depends_on` names, and "waits on" each of those the target does not hold
done. A cycle is stale when the folder of a member not yet archived, or a
file one of its tasks names, differs on the target from the same path at
the cycle's `planned` commit, unless every commit on the target's
first-parent history since then that changed it also edited the cycle
file, archived a change the cycle itself lists, or archived a change that
a cycle this one depends on, directly or through other cycles, lists.

## Boundaries

This covers the cycle file, `/grove-plan`, the `unattended` mark, `grove
check`'s reports on cycles, and how `grove changes`, `grove review` and
the project index show them. Building a cycle belongs to
`grove/specs/implement/spec.md`, and building
it with no one watching to `grove/specs/runs/spec.md`. What makes a plan,
and the mark, count is the merge into the target (`CONSTITUTION.md` rule
4).

## Requirements

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

### Requirement: Plan one change

When the owner runs `/grove-plan <change-id>` for a change the target
approves, the system shall print instructions for writing a new cycle file
whose only member is that change.

#### Scenario: Send one change to the runner

- GIVEN the target approves change C, in no cycle
- WHEN the owner runs `/grove-plan C`
- THEN a new cycle file lists only C, with the change's title as the
  cycle's title unless the owner gives another

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

### Requirement: Ask which cycles a new cycle follows

While planning a new cycle, when the project index lists cycles not yet
done, the plan instructions shall have the agent ask the owner whether
the new cycle follows any of them, and write under `depends_on` only the
cycles the owner names.

#### Scenario: The owner names one

- GIVEN the project index lists cycles K and L, neither done
- WHEN the owner runs `/grove-plan 0.5 release` and, asked, says it
  follows L
- THEN the new cycle file holds `depends_on: [L]`

#### Scenario: The owner names none

- GIVEN the project index lists cycle K, not done
- WHEN the agent asks and the owner says the new cycle follows no cycle
- THEN the new cycle file has no `depends_on`

#### Scenario: No other cycle

- GIVEN the project index lists no cycle that is not done
- WHEN the owner runs `/grove-plan 0.5 release`
- THEN the agent asks nothing about other cycles, and the cycle file has
  no `depends_on`

### Requirement: Say whether the cycle can run unattended

When the cycle file is written, the plan instructions shall have the agent
tell the owner that `/grove-implement <cycle-id>` builds the cycle with
them, then whether `grove claim <cycle-id>` would claim it, and that
`unattended: true`, merged on its own, hands it to a `grove claim` given
no cycle, such as a timer's; and, for each obstacle, what removes it.

#### Scenario: Nothing in the way

- GIVEN members A and B, neither of which waits on a change outside the
  cycle
- WHEN the agent finishes the plan
- THEN it names `/grove-implement <cycle-id>` first, says
  `grove claim <cycle-id>` would claim the cycle, and says a `grove claim`
  given no cycle would claim it once `unattended: true` is merged in, and
  how that mark is added

#### Scenario: A member waits on a change outside the cycle

- GIVEN member B depends on change D, neither in the cycle nor archived as
  finished
- WHEN the agent finishes the plan
- THEN it says `grove claim` would not claim the cycle until D is archived
  as finished

#### Scenario: The cycle waits on a cycle

- GIVEN the owner said the new cycle follows cycle L, which is not done
- WHEN the agent finishes the plan
- THEN it says `grove claim` would not claim the cycle until L is done,
  and that the merges archiving L's changes leave its plan as it checked
  it

#### Scenario: The runner needs something set up

- GIVEN member B has a task that calls a service needing a login, and the
  owner said the runner has it
- WHEN the agent finishes the plan
- THEN the verdict names the login the runner's checkout must have

### Requirement: Mark a cycle unattended

The system shall treat a cycle as handed to `grove claim` only while the
target's cycle file holds `unattended: true`, and the plan instructions
shall never write or remove that field.

#### Scenario: Planned but not marked

- GIVEN the target holds cycle K without `unattended`
- WHEN `grove claim` is invoked
- THEN K is not claimed

### Requirement: A cycle can depend on cycles

Where a cycle file's `depends_on` names cycles, the system shall have
`grove claim` claim none of that cycle, given or marked, while the target
does not hold one of them done, and have `/grove-implement <cycle-id>`
build it as it would with no `depends_on`.

#### Scenario: The cycle it follows is not done

- GIVEN cycle K's `depends_on` names cycle L, and L's change C is not
  archived on the target
- WHEN the operator runs `grove claim K`
- THEN nothing is recorded, no branch is created, and Grove says K waits
  on cycle L, which is not done

#### Scenario: The cycle it follows is done

- GIVEN cycle K's `depends_on` names cycle L, every member of L is
  archived on the target, and nothing else keeps K from being claimed
- WHEN the operator runs `grove claim K`
- THEN an attempt is recorded for K, and K's ID is printed

#### Scenario: Built with the owner

- GIVEN cycle K's `depends_on` names cycle L, which is not done
- WHEN the owner runs `/grove-implement K`
- THEN the instructions list K's members in build order, as they would
  with no `depends_on`

### Requirement: Report a cycle that cannot run unattended

When the owner runs `grove check`, the system shall report a cycle file
that holds `unattended: true` while the cycle is stale, naming what to do.

#### Scenario: Marked and stale

- GIVEN cycle K holds `unattended: true`, and member A's proposal was
  revised on the target after K's `planned` commit
- WHEN the owner runs `grove check`
- THEN it reports K, saying A changed since the plan and to re-plan K with
  `/grove-plan K`

#### Scenario: Marked and a member's file moved

- GIVEN cycle K holds `unattended: true`, member B's T-01 names
  `ledger/records.py`, and the target merged a commit touching that file
  after K's `planned` commit
- WHEN the owner runs `grove check`
- THEN it reports K, saying B changed since the plan and to re-plan K with
  `/grove-plan K`, and `grove claim` does not claim K

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

### Requirement: Show cycles

When the owner runs `grove changes`, the system shall show, with each
approved change not archived that a cycle not yet done lists, that cycle's
ID, whether it is marked `unattended`, and each cycle it waits on; the
instructions' project index shall name, under each cycle not yet done,
each cycle its `depends_on` names; and `grove review` shall label a cycle
file as a cycle.

#### Scenario: One marked, one not

- GIVEN change A is in marked cycle K, and change B is in no cycle
- WHEN the owner runs `grove changes`
- THEN A's line names K as unattended, and B's line names no cycle

#### Scenario: A cycle that waits on a cycle

- GIVEN change A is in cycle K, whose `depends_on` names cycle L, which
  is not done
- WHEN the owner runs `grove changes`
- THEN A's line names K and says it waits on cycle L

#### Scenario: The index names the cycle followed

- GIVEN cycle K, not done, whose `depends_on` names cycle L
- WHEN the owner runs `grove instructions plan "0.5 release"`
- THEN the project index lists, under K, `follows L`

### Requirement: Name a requirement two approved changes both change

When Grove prints a skill's instructions and two approved changes, neither
archived, each hold in their change specs an ADDED, MODIFIED or REMOVED
block for the same requirement of the same capability, and neither
depends on the other, directly or through other changes, the project
index shall name, under each of the two, the other change, the capability
and the requirement.

#### Scenario: Two changes modify one requirement

- GIVEN the target approves changes A and B, each with a change spec for
  `record-editing` that holds a MODIFIED block for "Keep a category", and
  neither names the other in `depends_on`
- WHEN the owner runs `grove instructions plan "0.5 release"`
- THEN the project index names, under A, change B, `record-editing` and
  "Keep a category", and under B, change A with the same

#### Scenario: One removes what the other modifies

- GIVEN A's change spec for `record-editing` holds a REMOVED block for
  "Keep a category" and B's a MODIFIED block for it
- WHEN the owner runs `grove instructions propose`
- THEN the project index names each under the other

#### Scenario: One depends on the other

- GIVEN A and B both modify "Keep a category", and B holds
  `depends_on: [A]`
- WHEN the owner runs `grove instructions plan "0.5 release"`
- THEN the project index names neither under the other

#### Scenario: Through another change

- GIVEN A and C both modify "Keep a category", C depends on B, and B
  depends on A
- WHEN the owner runs `grove instructions plan "0.5 release"`
- THEN the project index names neither under the other

#### Scenario: The same name in two capabilities

- GIVEN A modifies "Keep a category" of `record-editing`, and B adds a
  requirement "Keep a category" to `reports`
- WHEN the owner runs `grove instructions plan "0.5 release"`
- THEN the project index names neither under the other

#### Scenario: One of the two is archived

- GIVEN A and B both modified "Keep a category", and the target holds A
  archived
- WHEN the owner runs `grove instructions plan "0.5 release"`
- THEN the project index names no change under B

### Requirement: Plan tells the owner of members that share a requirement

While agreeing a cycle's members with the owner, when the project index
names a requirement two of them both change, the plan instructions shall
have the agent tell the owner which two changes and which requirement,
and that the one built second is stale on it once the first is archived,
ask whether both stay in this cycle or one goes to a later cycle, and
list under `changes` only what the owner then agrees to.

#### Scenario: The owner moves one to a later cycle

- GIVEN the project index names "Keep a category" of `record-editing`
  under A and B
- WHEN the owner runs `/grove-plan 0.5 release`, the agent proposes A and
  B, and the owner, told of the pair, keeps A and leaves B for later
- THEN the new cycle file lists A and not B

#### Scenario: The owner keeps both

- GIVEN that index
- WHEN the owner, told of the pair, keeps both
- THEN the new cycle file lists A and B, `grove check` reports nothing for
  the cycle, and `grove claim <cycle-id>` claims it when nothing else
  stops it

#### Scenario: No pair among the members

- GIVEN the project index names the pair A and B, and the agent proposes
  A and C for the cycle
- WHEN the owner agrees to A and C
- THEN the agent says nothing about requirements two members share

### Requirement: Count a member changed with its plan

When a commit on the target's first-parent history after a cycle's
`planned` commit changes the folder of a member not yet archived, or a
file one of its tasks names, the system shall treat the cycle as stale, in
`grove check` and in `grove claim`, unless that commit also edits the
cycle file, archives a change the cycle itself lists, or archives a change
that a cycle this one depends on, directly or through other cycles, lists.

#### Scenario: Revised and re-planned in one squash merge

- GIVEN cycle K, marked `unattended`, lists B, and on a branch B's
  proposal is revised and K's `planned` set to the target's commit
- WHEN the branch is squash-merged into the target and the owner runs
  `grove check`
- THEN it reports nothing for K, and `grove claim` claims K

#### Scenario: Revised and re-planned in one merge commit

- GIVEN a branch with one commit revising member B and another setting
  cycle K's `planned`
- WHEN the branch is merged into the target with a merge commit
- THEN K is not stale

#### Scenario: Marked after a member changed

- GIVEN member B's proposal was revised in one merge into the target, and
  cycle K was marked `unattended` in a later merge
- WHEN the owner runs `grove check`
- THEN it reports K, saying B changed since the plan and to re-plan K with
  `/grove-plan K`

#### Scenario: Revision and re-plan as two commits on the target

- GIVEN the revision of member B and the edit to cycle K's file reached
  the target as two commits, by rebase or fast-forward
- WHEN the operator runs `grove claim K`
- THEN nothing is recorded, and Grove says K's plan is stale

#### Scenario: A member lands in its own squash merge

- GIVEN cycle K, marked `unattended`, lists A and B, each with a task
  naming `ledger/cli.py`, and a squash merge into the target archives A,
  edits `ledger/cli.py` and leaves K's file as it is
- WHEN the owner runs `grove check`
- THEN it reports nothing for K, `grove claim` claims K, and
  `grove instructions implement B` names no file changed since K's plan

#### Scenario: Then a commit archiving nothing

- GIVEN that squash merge, and a later commit on the target that archives
  no change and edits `ledger/cli.py`
- WHEN the owner runs `grove check`
- THEN it reports K, saying B changed since the plan and to re-plan K with
  `/grove-plan K`

#### Scenario: The cycle it follows lands in a squash merge

- GIVEN cycle K, marked `unattended`, depends on cycle L, K's member B
  has a task naming `docs/reference.md`, and a squash merge into the
  target archives L's only change and edits `docs/reference.md`
- WHEN the owner runs `grove check`
- THEN it reports nothing for K, and `grove claim` claims K

#### Scenario: Stacked pull requests

- GIVEN cycle K depends on cycle L, whose changes C and D were built as
  stacked pull requests, and the target took them as two commits on its
  first-parent history, each archiving one of them and editing a file K's
  member B names
- WHEN the operator runs `grove claim K`
- THEN K's plan is not stale

#### Scenario: Through another cycle

- GIVEN cycle K depends on cycle L, L depends on cycle M, and a merge
  commit into the target archives M's change E and edits a file K's
  member B names
- WHEN the owner runs `grove check`
- THEN K is not stale

#### Scenario: Any other commit still counts

- GIVEN cycle K, marked `unattended`, depends on cycle L, and a commit on
  the target archives no change of L, nor of any cycle L depends on, and
  edits a file K's member B names
- WHEN the owner runs `grove check`
- THEN it reports K, saying B changed since the plan and to re-plan K with
  `/grove-plan K`

## Known discrepancies

None recorded.

## Open questions

None.

## History

- [261022-1grb0](../../changes/archive/261022-1grb0-cycles-for-unattended-runs/) Cycles planned for unattended runs
- [261023-rcxv2](../../changes/archive/261023-rcxv2-plan-verdict-no-remote-streamed-output/) Plan's verdict, runs without a remote, streamed output
- [261023-6j5bw](../../changes/archive/261023-6j5bw-tasks-name-their-files/) Tasks name their files
- [261023-3f9h6](../../changes/archive/261023-3f9h6-run-builds-a-given-cycle/) grove run builds a cycle it is given
- [261023-bgtcb](../../changes/archive/261023-bgtcb-revising-a-change-re-plans-its-cycle/) Revising a change re-plans its cycle
- [261023-2ej3w](../../changes/archive/261023-2ej3w-cycles-depend-on-cycles/) A cycle can depend on cycles
- [261030-t2pt7](../../changes/archive/261030-t2pt7-own-member-merges-keep-a-cycle-fresh/) A cycle's own member merges keep its plan fresh
- [261105-gmpa4](../../changes/archive/261105-gmpa4-plan-brings-a-stale-member-up-to-date/) Plan brings a stale member up to date
- [261105-3qdz2](../../changes/archive/261105-3qdz2-misnamed-cycle-and-shallow-clone-reported/) A misnamed cycle and a shallow clone are reported plainly
- [261108-jpg8n](../../changes/archive/261108-jpg8n-claim-replaces-run/) grove claim replaces grove run, with attempt records on a branch
- [261114-yx2b7](../../changes/archive/261114-yx2b7-remove-needs-human/) Remove needs_human
- [261118-emse7](../../changes/archive/261118-emse7-a-cycle-is-a-markdown-file/) A cycle is a Markdown file with plan's notes
- [261118-40s1a](../../changes/archive/261118-40s1a-index-names-requirement-overlaps/) The project index names requirements two changes both change
