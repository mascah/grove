# Implement

## Purpose

`/grove-implement` builds an approved change, or each change of a cycle,
in the owner's session, asking the owner instead of blocking and revising
the change on its branch with their agreement; with `--unattended` it asks
no one. `/grove-implement <change-id>` also builds a draft on the branch
that proposed it, which one merge then approves with its code, and
`/grove-implement <cycle-id>` builds each change of a cycle in turn. Where
the change does not settle how a task is built, it asks the owner instead
of blocking, and with their agreement revises the change on its
implementation branch, so the merge approves the revision with the code.
With `--unattended` it builds a change or a cycle the same way but asks no
one: what it cannot settle is marked blocked, and the change stays as
approved. "The owner" is the person in whose session the skill runs.

## Boundaries

This covers the implement instructions, with the owner and with
`--unattended`, how `grove review` names a draft built on its branch,
`grove.yaml`'s rules for implement, and the removal of the
apply skill that `--unattended` replaced. Starting a session with no one
watching, and what it refuses, belong to `grove run`
(`grove/specs/runs/spec.md`). Revising a change outside its implementation
branch belongs to `/grove-propose`. What makes a revision approved is the
merge into the target (`CONSTITUTION.md` rule 4).

## Requirements

### Requirement: Implement an approved change

When the owner runs `/grove-implement <change-id>`, without
`--unattended`, for a change the target approves and that waits on
nothing the working tree has not archived as finished, the system shall
print instructions for building its tasks in the owner's session, on a
branch other than the target, recording each task's status in
`implementation.md` and archiving the change when every task is done.

#### Scenario: A change with two tasks

- GIVEN the target approves change C with tasks T-01 and T-02, neither
  started
- WHEN the owner runs `/grove-implement C` on the target branch
- THEN the agent creates a branch, builds T-01 and T-02, records both as
  done in `implementation.md`, and archives C on that branch

#### Scenario: A change that waits

- GIVEN the target approves change B with `depends_on: [A]`, and the
  working tree does not hold A archived as finished
- WHEN the owner runs `/grove-implement B`
- THEN nothing is built, and Grove says B waits on A

### Requirement: Ask the owner instead of blocking

When, without `--unattended`, a task is decided by an open question in the
change, or cannot be built as the change says, the implement instructions
shall have the agent ask the owner and wait for the answer, instead of
marking the task blocked.

#### Scenario: An open question decides a task

- GIVEN change C's proposal has an open question that decides how T-02 is
  built
- WHEN the agent reaches T-02
- THEN it asks the owner the question before building T-02

#### Scenario: The owner leaves it open

- GIVEN the agent asked the owner a question deciding T-02
- WHEN the owner says to leave it
- THEN T-02 is marked `Status: blocked` with the question under it, and the
  agent goes on with the other tasks

### Requirement: Revise the change on its branch

When the owner's answer or feedback changes what the change says, the
implement instructions shall have the agent say which of the change's
proposal, change specs and tasks it would edit and how, and, only once the
owner agrees, edit them on the implementation branch, run `grove check`
clean, and then build to the revised change.

#### Scenario: Feedback changes a requirement mid-build

- GIVEN the agent built T-01 of change C, and the owner, trying it, asks for
  different behavior than C's change spec states
- WHEN the owner agrees to the edit the agent describes
- THEN C's change spec on the branch states the new behavior, the code
  follows it, and `grove review` on the branch shows the edited change spec

#### Scenario: The owner does not agree

- GIVEN the agent proposed an edit to C's change spec
- WHEN the owner declines it
- THEN C's files are unchanged and the agent builds what C says

### Requirement: Implement a cycle

When the owner runs `/grove-implement <cycle-id>`, without
`--unattended`, for a cycle the target holds, the system shall print
instructions for building each member not yet archived, in build order,
each as `/grove-implement <change-id>` builds it, in the owner's session.

#### Scenario: A cycle the owner builds

- GIVEN cycle K lists A and B, is not marked `unattended`, and
  `needs_human` names B
- WHEN the owner runs `/grove-implement K`
- THEN A and then B are built in the owner's session, the agent asking the
  owner where `/grove-implement` would

### Requirement: Tell the owner of a moved file

When the owner runs `/grove-implement <change-id>` for a change in no
cycle that is not done, and a task names a file that differs on the target
from the same file at the change's `base`, the system shall say which task
and file, and the implement instructions shall have the agent ask the owner
before building that task, instead of building on the file as it was.

#### Scenario: A named file changed since base

- GIVEN the target approves change C, in no cycle, whose T-02 names
  `internal/cli/run.go`, which the target changed after C's `base`
- WHEN the owner runs `/grove-implement C`
- THEN Grove's instructions name T-02 and `internal/cli/run.go` as changed
  since base, and the agent, on reaching T-02, tells the owner what changed
  and asks before building

#### Scenario: A file no task names needs a change

- GIVEN the agent, building T-01, finds `docs/concepts.md`, which no task
  names, states what T-01 makes false
- WHEN it tells the owner and they agree
- THEN T-01's `Files:` line on the branch names `docs/concepts.md`, and the
  agent edits the page; had the owner declined, the page would be left as
  it is

#### Scenario: A member of a cycle

- GIVEN change C is a member of cycle K, and a file C's T-01 names changed
  on the target after K's `planned` commit
- WHEN the owner runs `/grove-implement K`
- THEN C is not reported against its `base`; `grove check` reports K as
  stale, and the owner re-plans K or builds on

### Requirement: Verify before finishing

When the agent has built every task it can of a change, with the owner or
with `--unattended`, the implement instructions shall have it, before it
finishes, run the project's tests, re-read each done task's evidence
against the code as it now is and set a task whose evidence no longer
holds back to `Status: in progress`, and write a `## Scenarios` section in
`implementation.md` with one line per scenario of the change specs naming
the check that shows it and its result, or `unverified` and why; its
finishing report names each unverified scenario, and with `--unattended`
it asks no one.

#### Scenario: Every scenario has a test

- GIVEN a change whose change spec has two scenarios, both shown by tests
  that pass
- WHEN the agent finishes implementing it
- THEN `implementation.md`'s `## Scenarios` names each scenario with its
  test and "passed"

#### Scenario: A done task later broken

- GIVEN T-01 was marked done, and building T-02 changed the code T-01's
  evidence relies on so its test fails
- WHEN the agent rechecks before finishing
- THEN T-01 is `Status: in progress` again until it is rebuilt

#### Scenario: A scenario only a person can judge

- GIVEN a scenario says a list "reads clearly"
- WHEN the agent writes the scenario list
- THEN that line says `unverified` and why, and the finishing report
  names it

#### Scenario: A scenario no test shows, unattended

- GIVEN a member, built with `--unattended`, whose change spec has a
  scenario no test checks
- WHEN the agent finishes the member
- THEN that line says `unverified` and why, the report names it, and the
  agent goes on without asking

### Requirement: Archive names unverified scenarios

When the agent runs the archive instructions for a finished change, they
shall have it tell the owner each scenario `implementation.md` marks
`unverified`, before the fold, without refusing to archive.

#### Scenario: One scenario unverified

- GIVEN a finished change whose scenario list marks one scenario
  `unverified`
- WHEN the agent archives it
- THEN it names that scenario to the owner, then folds and archives

### Requirement: Build unattended

When `/grove-implement <id> --unattended` is run for a change the target
approves or a cycle the target holds, the system shall print instructions
that have the agent build it as `/grove-implement <id>` would but ask no
one: a task an open question decides, or that cannot be built as the
change says, is marked `Status: blocked` with the question or the reason;
in a change in no cycle that is not done, a task naming a file the target
changed since the change's `base` is marked `Status: blocked` with the
file and is not built on it;
a file no task names is left as it is and named in that task's evidence;
the change's files other than `implementation.md` are not edited; and a
cycle is built whether or not it is marked `unattended` and whatever its
`needs_human` says.

#### Scenario: One change, by hand

- GIVEN the target approves change C, in no cycle, with tasks T-01 and
  T-02 that no open question decides
- WHEN `/grove-implement C --unattended` is run on the target branch
- THEN the agent creates a branch, builds T-01 and T-02 without asking
  anyone, records both as done in `implementation.md`, and archives C on
  that branch

#### Scenario: An open question decides a task

- GIVEN cycle K lists change C, whose open question decides T-02
- WHEN `/grove-implement K --unattended` is run and the agent reaches C's
  T-02
- THEN T-02 is marked `Status: blocked` with the question, nobody is
  asked, and C's proposal, change specs and tasks are unchanged

#### Scenario: A named file moved since base

- GIVEN the target approves change C, in no cycle, whose T-02 names
  `internal/cli/run.go`, which the target changed after C's `base`
- WHEN `/grove-implement C --unattended` is run
- THEN Grove's instructions name T-02 and `internal/cli/run.go` as changed
  since base, T-01 is built, and T-02 is marked `Status: blocked` naming
  that file, unbuilt, with nobody asked

#### Scenario: Two members, one waiting on the other

- GIVEN cycle K lists A and B, B depending on A, neither archived
- WHEN `/grove-implement K --unattended` is run
- THEN A is built and archived, then B, on one branch, each by
  `grove instructions implement <change-id> --unattended`

#### Scenario: A cycle that needs a person, by hand

- GIVEN cycle K lists A and B, is not marked `unattended`, and
  `needs_human` names B
- WHEN `/grove-implement K --unattended` is run by hand
- THEN A and B are built asking no one; only `grove run` refuses K

#### Scenario: A change that waits

- GIVEN the target approves change B with `depends_on: [A]`, and the
  working tree does not hold A archived as finished
- WHEN `/grove-implement B --unattended` is run
- THEN nothing is built, and Grove says B waits on A

### Requirement: A member gets its cycle's plan

When `/grove-implement <change-id>` is run, with or without
`--unattended`, for a member of a cycle the target holds that is not
done, the system shall print that change's instructions with the cycle's
plan, naming the member's place in build order.

#### Scenario: The second member, unattended

- GIVEN cycle K lists A and B, B depending on A, and A is archived as
  finished on this branch
- WHEN `grove instructions implement B --unattended` is run
- THEN B's instructions are printed with K's plan, naming B as the second
  of two

#### Scenario: A member built with the owner

- GIVEN cycle K lists A and B, neither archived
- WHEN the owner runs `/grove-implement A`
- THEN A's instructions are printed with K's plan, naming A as the first
  of two

### Requirement: Implement's rules in both modes

Where `grove.yaml` holds lines under `rules.skills.implement`, the system
shall print them with the implement instructions, for a change and for a
cycle, with and without `--unattended`.

#### Scenario: A rule for both modes

- GIVEN `rules.skills.implement` holds "Commit after each task, naming its
  task ID."
- WHEN `grove instructions implement C` runs, and then
  `grove instructions implement C --unattended`
- THEN both print that line with the implement instructions

### Requirement: Report rules for apply

When the owner runs `grove check` and `grove.yaml` holds a
`rules.skills.apply` key, the system shall report that line, saying its
lines go under `rules.skills.implement`.

#### Scenario: A project that kept its apply rules

- GIVEN `grove.yaml` holds two lines under `rules.skills.apply`
- WHEN the owner runs `grove check`
- THEN it reports the `apply` line, pointing to `implement`

### Requirement: No apply skill

The system shall write no apply skill or command for any harness, `grove
init` shall remove the apply skill and command files an earlier Grove
wrote, and `grove instructions apply` shall print the usage and exit with
status 2.

#### Scenario: Init after upgrading

- GIVEN a project whose `.claude/skills/grove-apply/SKILL.md` an earlier
  Grove wrote
- WHEN the owner runs `grove init`
- THEN that folder is removed, and init names `grove-apply` among what it
  removed

#### Scenario: The old command

- GIVEN cycle K on the target
- WHEN `grove instructions apply K` is run
- THEN no instructions are printed, and Grove prints the usage

### Requirement: Implement a draft with the owner

When the owner runs `/grove-implement <change-id>`, without
`--unattended`, for a change this working tree holds and the target does
not, and that waits on nothing the working tree has not archived as
finished, the system shall print instructions for building it as the
working tree holds it, in the owner's session, that say the change is a
draft, that merging the branch into the target approves it with its
code, and that its proposal, change specs and tasks are edited only with
the owner's agreement, as when revising a change on its branch.

#### Scenario: Drafted and built on one branch

- GIVEN the owner proposed change D on branch `rename` and committed it,
  and the target does not hold D
- WHEN the owner runs `/grove-implement D` on `rename`
- THEN Grove's instructions say D is a draft that merging `rename`
  approves with its code, and the agent builds D's tasks, records them as
  done in `implementation.md`, and archives D on `rename`

#### Scenario: Unattended refuses a draft

- GIVEN change D exists in the working tree but not in the target
- WHEN `/grove-implement D --unattended` is run
- THEN nothing is built, and Grove says the target holds no approved
  change D

#### Scenario: Building shows the draft is wrong

- GIVEN the agent, building draft D's T-01, finds that D's change spec
  states behavior the owner, trying it, does not want
- WHEN the agent says which of D's files it would edit and how
- THEN it edits them only once the owner agrees, runs `grove check`
  clean, and builds to the edited change; had the owner declined, D's
  files would be unchanged

#### Scenario: A large draft

- GIVEN draft D has eight tasks across four capabilities
- WHEN the owner runs `/grove-implement D`
- THEN Grove prints the instructions for D as for a draft of one task,
  with nothing about its size

#### Scenario: A draft that waits

- GIVEN draft D has `depends_on: [A]`, and the working tree does not
  hold A archived as finished
- WHEN the owner runs `/grove-implement D`
- THEN nothing is built, and Grove says D waits on A

### Requirement: Review names a change approved with its code

When the current branch holds a change archived as finished that the
target holds neither open nor archived, `grove review` shall say that
merging the branch into the target approves that change with its code,
instead of as future work.

#### Scenario: Drafted, built and archived on one branch

- GIVEN branch `rename` holds change D archived as finished, and the
  target does not hold D
- WHEN the owner runs `grove review` on `rename`
- THEN its first line is "Merging rename into main approves change D
  with its code."

#### Scenario: A draft not yet built

- GIVEN branch `rename` holds change D open, and the target does not hold
  D
- WHEN the owner runs `grove review`
- THEN it says merging `rename` approves change D as future work and does
  not implement it

#### Scenario: An approved change archived on its branch

- GIVEN the target holds change C open, and branch `c-impl` holds C
  archived as finished
- WHEN the owner runs `grove review` on `c-impl`
- THEN it names no change approved by the merge, and lists C's archived
  folder and the maintained specs it changes

### Requirement: Honor a task's delegation

When a task the owner is building with `/grove-implement` has a
`Delegated:` line, the implement instructions shall have the agent make the choices it names and change the
paths it names without asking the owner, and ask, as for any other, about
a choice or a file neither the task's `Files:` line nor its `Delegated:`
line names.

#### Scenario: A delegated path

- GIVEN T-01 names `ledger/records.py` and `tests/test_cli.py` on its
  `Files:` line, and its `Delegated:` line names test layout and
  `tests/test_records.py`
- WHEN the agent, building T-01, puts `find`'s own tests in a new
  `tests/test_records.py`
- THEN it creates the file without asking the owner

#### Scenario: A path nothing names

- GIVEN the same T-01, and the agent finds `README.md` states what T-01
  makes false
- WHEN it would edit `README.md`
- THEN it tells the owner and asks first, as for any file no task names

### Requirement: Stop on a discovery that breaks the approach

If, while building a task, the agent finds something that makes the
change's agreed approach wrong, even where the task could still be built
as written, then the implement instructions shall have it stop that task,
build nothing more of it, and tell the owner what it found and the
revision the change needs, naming which of its proposal, change specs and
tasks it would edit and how.

#### Scenario: The approach rests on something false

- GIVEN change C's T-01 says `find` in `ledger/records.py` raises
  `ValueError` for an unknown ID, and the agent finds `report` relies on
  `find` returning nothing for a missing ID
- WHEN it makes the finding while building T-01
- THEN it stops T-01, tells the owner what it found and which of C's
  documents it would edit and how, and builds nothing more of T-01 until
  the owner answers

### Requirement: Ignore a task's delegation unattended

When a task built with `/grove-implement <id> --unattended` has a
`Delegated:` line, the implement instructions shall have the agent ignore
it: change only the files the tasks' `Files:` lines name, and leave a path
only a `Delegated:` line names as it is, as it does any file no task
names: when the task would have changed it, naming it, and why, in that
task's evidence in `implementation.md`.

#### Scenario: A delegated path, unattended

- GIVEN cycle K lists change C, whose T-01 names `ledger/records.py` and
  `tests/test_cli.py` on its `Files:` line and `tests/test_records.py` on
  its `Delegated:` line
- WHEN `/grove-implement K --unattended` builds T-01
- THEN `tests/test_records.py` is not created, `find`'s tests go in
  `tests/test_cli.py`, and T-01's evidence names `tests/test_records.py`
  and why it was left

## Known discrepancies

None recorded.
