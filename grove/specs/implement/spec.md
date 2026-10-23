# Implement

## Purpose

`/grove-implement <change-id>` builds an approved change in the owner's own
session, and `/grove-implement <cycle-id>` each change of a cycle in turn. Where the change does not settle how a task is built, it asks the
owner instead of blocking, and with their agreement revises the change on
its implementation branch, so the merge approves the revision with the
code. "The owner" is the person in whose session the skill runs.

## Boundaries

This covers the implement instructions and, by contrast, the rule that
apply asks no one. Building a cycle without asking anyone, and building
without anyone watching, belong to apply (`grove/specs/apply/spec.md`)
and to `grove run`
(`grove/specs/runs/spec.md`). Revising a change outside its implementation
branch belongs to `/grove-propose`. What makes a revision approved is the
merge into the target (`CONSTITUTION.md` rule 4).

## Requirements

### Requirement: Implement an approved change

When the owner runs `/grove-implement <change-id>` for a change the target
approves and that waits on nothing the working tree has not archived as
finished, the system shall print instructions for building its tasks in
the owner's session, on a branch other than the target, recording each
task's status in `implementation.md` and archiving the change when every
task is done.

#### Scenario: A change with two tasks

- GIVEN the target approves change C with tasks T-01 and T-02, neither
  started
- WHEN the owner runs `/grove-implement C` on the target branch
- THEN the agent creates a branch, builds T-01 and T-02, records both as
  done in `implementation.md`, and archives C on that branch

#### Scenario: A change that is only a draft

- GIVEN change D exists in the working tree but not in the target
- WHEN the owner runs `/grove-implement D`
- THEN nothing is built, and Grove says the target holds no approved change
  D

#### Scenario: A change that waits

- GIVEN the target approves change B with `depends_on: [A]`, and the
  working tree does not hold A archived as finished
- WHEN the owner runs `/grove-implement B`
- THEN nothing is built, and Grove says B waits on A

### Requirement: Ask the owner instead of blocking

When a task is decided by an open question in the change, or cannot be
built as the change says, the implement instructions shall have the agent
ask the owner and wait for the answer, instead of marking the task
blocked.

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

### Requirement: Apply asks no one

The apply instructions shall keep their rule: a task an open question
decides, or that cannot be built as the change says, is marked
`Status: blocked` without asking anyone, and the change's files other than
`implementation.md` are not edited.

#### Scenario: The same open question under apply

- GIVEN cycle K lists change C, whose open question decides T-02
- WHEN K is built with `/grove-apply K` and the agent reaches C's T-02
- THEN T-02 is marked `Status: blocked` with the question, nobody is asked,
  and C's proposal, change specs and tasks are unchanged

### Requirement: Implement a cycle

When the owner runs `/grove-implement <cycle-id>` for a cycle the target
holds, the system shall print instructions for building each member not yet
archived, in build order, each as `/grove-implement <change-id>` builds it,
in the owner's session.

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

When the agent has built every task it can of a change with the owner,
the implement instructions shall have it, before it finishes, run the
project's tests, re-read each done task's evidence against the code as it
now is and set a task whose evidence no longer holds back to `Status: in
progress`, and write a `## Scenarios` section in `implementation.md` with
one line per scenario of the change specs naming the check that shows it
and its result, or `unverified` and why; its finishing report names each
unverified scenario.

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

### Requirement: Archive names unverified scenarios

When the agent runs the archive instructions for a finished change, they
shall have it tell the owner each scenario `implementation.md` marks
`unverified`, before the fold, without refusing to archive.

#### Scenario: One scenario unverified

- GIVEN a finished change whose scenario list marks one scenario
  `unverified`
- WHEN the agent archives it
- THEN it names that scenario to the owner, then folds and archives

## Known discrepancies

None recorded.
