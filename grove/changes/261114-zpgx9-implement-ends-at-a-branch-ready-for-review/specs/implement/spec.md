# Changes to implement

The changes that change 261114-zpgx9 proposes to the requirements of
[implement](../../../../specs/implement/spec.md), as at 999ab56.

## ADDED Requirements

### Requirement: Record each choice for the reviewer

When a task's acceptance seems to need a person, or the task rests on a
choice the change leaves open, including an open question that decides
it, the implement instructions shall have the agent first run the
project's app or command and look; build its best choice; mark the task
`Status: done` once its evidence holds; write under it a `Check:` line
saying what it chose or could not see, such as a real call, hardware or a
permission prompt, and why; go on without stopping; and list every
`Check:` line first in its report.

#### Scenario: A screen to look at

- GIVEN change C, whose T-01 adds a settings screen and whose T-02 adds a
  command-line flag, and C does not say how the screen shows an error
- WHEN `/grove-implement C` builds it
- THEN the agent runs the app and looks at the screen, T-01 is built with
  its choice and is `Status: done` with a `Check:` line saying what it
  chose and why, T-02 is `Status: done`, C is archived on the branch, and
  the report lists T-01's `Check:` line first

#### Scenario: An open question, low cost

- GIVEN change C's open question asks whether a list is sorted by name or
  by date
- WHEN `/grove-implement C` builds the task it decides
- THEN the agent builds one order, writes a `Check:` line naming it and
  why, asks no one, and the task is `Status: done`

### Requirement: Block only a costly choice

Where a wrong choice in a task would cost data, break an interface other
tasks build on, or be hard to undo, the implement instructions shall have
the agent build nothing of that task, mark it `Status: blocked` with a
`Check:` line holding the question and why, leave unbuilt each task that
builds on it, naming it, build every other task, and finish with the
branch and the question, waiting for no one.

#### Scenario: A format at stake

- GIVEN change C's open question decides how T-02 stores records, so a
  wrong guess would change the format of `records.json`
- WHEN `/grove-implement C` is run
- THEN the agent builds T-01 and T-03, builds nothing of T-02, marks it
  `Status: blocked` with the question on a `Check:` line, does not
  archive C, and finishes without waiting, its report listing T-02's
  question first

#### Scenario: A task that builds on a blocked one

- GIVEN T-01 of change C is `Status: blocked`, and T-03 lays out its page
  inside the screen T-01 adds
- WHEN the agent reaches T-03
- THEN T-03 is not built, and `implementation.md` names T-01 under it

### Requirement: Take a blocked task's answer

When a blocked task is answered, by the owner in the session, in the
words that started the session, or by an `Answer:` line already under the
task, the implement instructions shall have the agent write the owner's
words, quoted, under the task as an `Answer:` line if none is there, and
build on it, revising the change as "Settle what the change does not"
says where the answer changes what the change says.

#### Scenario: An answer in the words that start a session

- GIVEN T-02 of change C is `Status: blocked` with a `Check:` line asking
  which ID format records use
- WHEN a session is started on the branch with the owner's words "use
  six-digit IDs" beside `/grove-implement C`
- THEN T-02 has `Answer: "use six-digit IDs"` under it, is built to it,
  and is `Status: done`

## MODIFIED Requirements

### Requirement: Build a change, a draft or a cycle

When `/grove-implement <id>` is run, the system shall print instructions
for building, on a branch other than the target, the approved change
`<id>`, or each member of cycle `<id>` not yet archived in build order,
each listed with the members of the cycle it depends on and each member's
instructions naming its place in the cycle's plan; or a change only the
working tree holds, as a draft that merging the branch approves with its
code. The instructions shall have the agent ask no one and wait for no
one, and end with the branch ready for review. The agent records each
task's status in `implementation.md`, builds no task already done again
unless the check before finishing finds its evidence no longer holds, and
archives the change when every task is done. The system shall print no
instructions, and say why, for a change that waits on one the working
tree has not archived as finished; and a cycle's instructions shall leave
out, saying why, a member that waits on a change neither being built nor
archived as finished.

#### Scenario: An approved change

- GIVEN the target approves change C with tasks T-01 and T-02
- WHEN `/grove-implement C` is run on the target branch
- THEN the agent creates a branch, builds T-01 and T-02, records both as
  done in `implementation.md`, and archives C on that branch

#### Scenario: A draft

- GIVEN change D was proposed on branch `rename` and committed, and the
  target does not hold D
- WHEN `/grove-implement D` is run on `rename`
- THEN the instructions say D is a draft that merging `rename` approves
  with its code, and the agent builds and archives D on `rename`

#### Scenario: A change that waits

- GIVEN change B, approved or a draft, with `depends_on: [A]`, and the
  working tree does not hold A archived as finished
- WHEN `/grove-implement B` is run
- THEN nothing is built, and Grove says B waits on A

#### Scenario: A cycle

- GIVEN cycle K lists A and B, B depending on A, and K is not marked
  `unattended`
- WHEN `/grove-implement K` is run
- THEN K's instructions list B `after A`, and A and then B are built on
  one branch, each by `grove instructions implement <change-id>`, B's
  instructions naming it the second of two in K's plan; only `grove claim`
  reads `unattended`

#### Scenario: A member left out

- GIVEN cycle K lists A and C, and C depends on Z, which K does not list
  and the target has not archived as finished
- WHEN `/grove-implement K` is run
- THEN the instructions leave C out, saying it waits on Z, and A is built

#### Scenario: The flag is gone

- GIVEN approved change C
- WHEN `grove instructions implement C --unattended` is run
- THEN Grove prints its usage and exits 2

### Requirement: Build a cycle's changes in their own contexts

Where the agent's harness lets it start its own subagents, the cycle's
implement instructions shall have the session building a cycle build each
member in a fresh context of its own by
`grove instructions implement <change-id>`, told to neither archive nor
commit; start at once every member whose dependencies are archived and
whose tasks' `Files:` name no file a member being built names; itself
archive and commit each returned member whose tasks are all done, one at
a time; where two members built together changed the same file, run the
project's tests once both returned and, if they fail, rebuild the
later-started member on its own before archiving it; and list every
member's `Check:` lines and blocked tasks first in its report; and, where
the harness has none, the session shall build each member itself, one
after another.

#### Scenario: Two members that share nothing

- GIVEN cycle K lists A, B and C, C depending on A and B, and no file in
  A's tasks' `Files:` is in B's
- WHEN `/grove-implement K` is run in a harness with subagents
- THEN A and B are built at the same time, each in its own subagent by
  `grove instructions implement <change-id>`; the session archives and
  commits each as it returns, one at a time; and C is started once both
  are archived

#### Scenario: Members that share a file

- GIVEN cycle K lists A and B, which depend on nothing, and both name
  `dev/walkthrough.md` in their tasks' `Files:`
- WHEN `/grove-implement K` is run in a harness with subagents
- THEN B's subagent is started only once A is archived

#### Scenario: A choice in a member

- GIVEN A and B of cycle K are built at the same time, and A's subagent
  returns T-01 as `Status: done` with a `Check:` line
- WHEN both have returned
- THEN the session archives both, asks no one, and its report lists A's
  T-01 `Check:` line first

#### Scenario: Two members changed the same file

- GIVEN A and B of cycle K are built at the same time, A's subagent
  changes `internal/cli/check.go` with an `Also changed:` line, B's tasks
  name that file, and the project's tests fail once both have returned
- WHEN the session checks them before archiving
- THEN B is built again on its own, after A is archived, and is archived
  only once its evidence holds

#### Scenario: A harness with no subagents

- GIVEN cycle K lists A and B, which share nothing
- WHEN `/grove-implement K` is run in a harness that cannot start
  subagents
- THEN the session builds A and then B itself, each by
  `grove instructions implement <change-id>`

### Requirement: Name the change's documents

When `grove instructions implement <change-id>` or `grove instructions
archive <change-id>` is run, the system shall name each document of the
change by its path, with the instruction to read every one first, print
none of their text and no template but `implementation.md`'s, and, for an
approved change, name each document whose working-tree copy differs from
the target's with the `git show` command that prints the approved text;
the agent builds the branch's revision, which the merge approves.

#### Scenario: A large change

- GIVEN approved change C whose documents hold 45 KB of text
- WHEN `grove instructions implement C` is run
- THEN the output names each of C's documents by path, holds none of
  their text, and prints the implementation template and no other

#### Scenario: A document that differs

- GIVEN approved change C, whose `tasks.md` an earlier session revised on
  this branch
- WHEN `grove instructions implement C` is run
- THEN it names `tasks.md` as differing, with `git show <commit>:<path>`,
  and the agent builds the revision

### Requirement: Settle what the change does not

When a task cannot be built as the change says, shows the change's
approach is wrong, or the owner's words or `grove check` call for a
revision, the implement instructions shall have the agent revise the
change on the branch, asking no one, by following
`grove instructions propose <change-id>`, which re-bases it, and then
running `grove check`; set a done task whose result the revision changes
back to in progress until it is rebuilt; write a `Check:` line under each
task the revision touches; and build to the revision, the merge of the
branch approving it with the code.

#### Scenario: The approach rests on something false

- GIVEN C's T-01 says `find` raises `ValueError` for an unknown ID, and
  the agent finds `report` relies on `find` returning nothing
- WHEN it makes the finding while building T-01
- THEN it revises C's tasks on the branch by following
  `grove instructions propose C`, writes a `Check:` line under T-01 saying
  what it found and how it revised C, builds T-01 to the revision, and
  asks no one

#### Scenario: Feedback changes the change

- GIVEN the owner, in the session, asks for behavior C's change spec does
  not state
- WHEN the agent takes it
- THEN it edits C's change spec on the branch by following
  `grove instructions propose C`, sets a done task the edit changes to
  in progress until rebuilt, and builds to the edit

### Requirement: Change only the files the tasks name

The implement instructions shall have the agent change the files the
tasks' `Files:` lines name, make the choices and change the paths a
task's `Delegated:` line names, and, for any other file a task needs
changed to be built as the change says or makes false, edit it without
asking anyone and without marking the task blocked, and write
`` Also changed: `<path>` — <why> `` under that task in
`implementation.md`; edit the change's other documents only by revising
it as "Settle what the change does not" says; and never edit
`grove.yaml` for it.

#### Scenario: A delegated path

- GIVEN T-01's `Files:` line names `tests/test_cli.py` and its
  `Delegated:` line names `tests/test_records.py`
- WHEN the agent puts `find`'s tests in a new `tests/test_records.py`
- THEN it creates the file without asking anyone

#### Scenario: A file nothing names

- GIVEN the agent, building T-01, finds `README.md`, which no line of
  T-01 names, states what T-01 makes false
- WHEN it builds T-01
- THEN it edits `README.md` without asking and writes
  `` Also changed: `README.md` — T-01 makes its usage example false ``
  under T-01, and `grove review` shows `README.md` as also changed by
  T-01 with that reason

#### Scenario: A file the task needs

- GIVEN T-01 names only `ledger/cli.py` and `tests/test_cli.py`, and its
  `Interfaces:` says `find` in `ledger/records.py` raises `ValueError` for
  an unknown ID
- WHEN `/grove-implement C` builds it
- THEN `ledger/records.py` is changed, T-01 is `Status: done` with an
  `Also changed:` line naming `ledger/records.py` and why, and T-01 is not
  blocked

### Requirement: A task's file changed on the target

When a task names a file, or a folder holding one, that the target changed
since the change's `base`, or, for a member of a cycle not done, since
the cycle's `planned` commit as the cycles spec's "Count a member changed
with its plan" counts it, the system shall name that task and file in the
implement instructions, and they shall have the agent read the file's
diff from that commit, build on the file as it is now and never as it
was, write under the task a `Check:` line naming the change on the target
and what it meant for the task, and, where the diff changes what the
change says, revise the change as "Settle what the change does not" says,
asking no one.

#### Scenario: Changed since base

- GIVEN approved change C, in no cycle, whose T-02 names
  `internal/cli/run.go`, which the target changed after C's `base`
- WHEN `/grove-implement C` is run
- THEN the instructions name T-02 and `internal/cli/run.go` as changed
  since base, and the agent builds T-02 on the file as it is now without
  asking, with a `Check:` line under T-02 naming the change

#### Scenario: Its dependency landed

- GIVEN change B depends on A, and the merge that archived A on the
  target changed `ledger/export.py` so that B's T-01, as written, no
  longer fits it
- WHEN `/grove-implement B` reaches T-01
- THEN the agent revises B on the branch by following
  `grove instructions propose B`, which re-bases it, and builds T-01 to
  the revision without asking

#### Scenario: A member of a cycle

- GIVEN change C is a member of cycle K, and the target changed a file
  C's T-01 names after K's `planned` commit
- WHEN `/grove-implement C` is run
- THEN the instructions name T-01 and the file as changed since K's plan,
  not since C's `base`; once K is re-planned, or when the change came by
  the merge that archived a cycle K depends on, they name nothing

### Requirement: Verify before finishing

When the agent has built every task it can of a change, the implement
instructions shall have it run the project's tests, re-read each done
task's evidence against the code as it now is and set a task whose
evidence no longer holds back to `Status: in progress`, and write a
`## Scenarios` section in `implementation.md` with one line per scenario
of the change specs naming the check that shows it and its result, or
`unverified` and why; its report, and the archive instructions before the
fold, shall name each unverified scenario, refusing nothing.

#### Scenario: Every scenario has a test

- GIVEN a change whose two scenarios are shown by passing tests
- WHEN the agent finishes it
- THEN `## Scenarios` names each with its test and "passed"

#### Scenario: A done task later broken

- GIVEN T-01 was done, and building T-02 broke the test T-01's evidence
  relies on
- WHEN the agent rechecks before finishing
- THEN T-01 is `Status: in progress` again until it is rebuilt

#### Scenario: A scenario only a person can judge

- GIVEN a scenario says a list "reads clearly"
- WHEN the agent finishes, and later archives, the change
- THEN that line says `unverified` and why, the report names it, and
  archive names it to the owner before folding

### Requirement: Rules for implement

Where `grove.yaml` holds lines under `rules.skills.implement`, the system
shall print them with the implement instructions; when the agent finishes
a change in which the owner corrected it, or it found something, in a way
that would apply to future work, the instructions shall have its report
name at most one line for `context` or a key under `rules`, with what
prompted it, and have it add the line on the branch only if the owner
agrees in the session, naming it in no task's `Files:` line.

#### Scenario: A rule

- GIVEN `rules.skills.implement` holds "Commit after each task, naming its
  task ID."
- WHEN `grove instructions implement C` is run
- THEN the instructions print that line

#### Scenario: The owner agrees

- GIVEN the owner said "Name each test after the requirement it checks"
  while C was built
- WHEN the agent finishes C and the owner agrees to the line it names
- THEN `grove.yaml` on the branch holds that one line, and no task's
  `Files:` line names `grove.yaml`; had the owner declined, or not been
  in the session, `grove.yaml` is unchanged

#### Scenario: No one in the session

- GIVEN A's tests run only with a command `context` does not name
- WHEN `/grove-implement K` builds A in a session no one watches
- THEN its report names a `context` line giving that command, and
  `grove.yaml` is unchanged

## REMOVED Requirements

### Requirement: Stop at a checkpoint

Implement no longer stops: "Record each choice for the reviewer" and
"Block only a costly choice" replace it.

### Requirement: Take a checkpoint's answer

With no stop and no `Status: to check`, only a blocked task waits on an
answer: "Take a blocked task's answer" replaces it.
