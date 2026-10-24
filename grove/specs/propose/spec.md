# Propose

## Purpose

`/grove-propose` turns a chosen direction into a change the owner can
review: a proposal, change specs and tasks. This covers how it settles
what those documents say, what they must say, and what `grove check`
reports on them. A task's named
files are the repo-relative paths on its `Files:` line.

## Boundaries

This covers the documents `/grove-propose` writes and `grove check`'s
reports on a change. Building a change belongs to
`grove/specs/implement/spec.md`, and a
cycle's staleness to `grove/specs/cycles/spec.md`.

## Requirements

### Requirement: Tasks name files and interfaces

When a change's tasks are written, each task shall name, on a `Files:`
line, the repo-relative paths it creates, modifies or tests, and, on an
`Interfaces:` line, the functions, types or commands it produces that
other tasks or callers use, with no code, commands or steps.

#### Scenario: A task that fixes a function

- GIVEN the owner decided, while proposing, that `find` in
  `ledger/records.py` raises `ValueError` for an unknown ID
- WHEN the agent writes the task
- THEN its `Files:` line names `ledger/records.py` and `tests/test_cli.py`,
  and its `Interfaces:` line says `find` raises `ValueError` with the
  message the requirement states

#### Scenario: A task that updates documentation

- GIVEN the owner wants the README's planned list updated in this change
- WHEN the agent writes the task
- THEN its `Files:` line names `README.md`, and the task has no steps

### Requirement: Report a named file that moved

When the owner runs `grove check`, the system shall report, for each
change in no cycle that is not done, each task whose named file differs on
the target from the same file at the change's `base`, naming the task and
the file and saying to revise the change with propose.

#### Scenario: A file changed since base

- GIVEN change C, in no cycle, proposed at `base` b143991, with T-01 naming
  `internal/cli/run.go`, and the target has since merged a commit touching
  that file
- WHEN the owner runs `grove check`
- THEN it reports C's `tasks.md`, saying T-01's `internal/cli/run.go`
  changed since base b143991 and to revise the change with propose

#### Scenario: A member of a cycle

- GIVEN change C is a member of cycle K, not done, and T-01's named file
  changed on the target after K's `planned` commit
- WHEN the owner runs `grove check`
- THEN C is not reported against its `base`, and K is reported as stale
  (`grove/specs/cycles/spec.md`)

#### Scenario: A file this change's own branch touched

- GIVEN change C is being built on a branch, where T-01 edited its named
  file
- WHEN the owner runs `grove check` on that branch
- THEN C is not reported, since the file is unchanged on the target

### Requirement: Ask a round before writing

When the owner has confirmed the scope of a new change, the propose
instructions shall have the agent list, before writing any document, every
decision it would otherwise assume, numbered, each with a recommended
answer worded so that "yes" accepts it, in this order: those that change
behavior, scope or compatibility or lose data first, then those that fix
architecture or interfaces costly to reverse, and routine choices inside
the agreed approach last, as recommended defaults on one line that a
single "yes" accepts; and wait for the owner. Each answer is recorded as a
decision in the proposal, quoted, and a decision the owner declines to
make is recorded as an open question naming the task it blocks.

#### Scenario: The owner accepts the round

- GIVEN the scope is confirmed and the agent would otherwise assume where
  the check goes and whether the README is updated
- WHEN it lists both with recommended answers and the owner says "yes"
- THEN the proposal records each as a decision quoting the owner's "yes"
  and the recommended answer, and no document lists an assumption

#### Scenario: The round is ordered

- GIVEN the scope of "Editing an unknown ID prints an error" is confirmed,
  and the agent would otherwise assume the exit status, whether `find` in
  `ledger/records.py` raises `ValueError` or `ledger/cli.py` checks the
  ID, and the test's name and place in `tests/test_cli.py`
- WHEN it asks the round
- THEN the exit status comes first, where the check goes second, and the
  test's name and place last, on one line of recommended defaults, and
  none is left out

#### Scenario: One "yes" to the routine line

- GIVEN the round's last line recommends three routine choices
- WHEN the owner says "yes" to the round
- THEN the proposal records each of the three as a decision quoting the
  owner's "yes"

#### Scenario: The owner declines one

- GIVEN the round lists where the check goes
- WHEN the owner says to leave that one open
- THEN the proposal lists it under open questions, naming the task it
  decides, and `/grove-implement --unattended` leaves that task blocked
  while `/grove-implement` asks

#### Scenario: Nothing to assume

- GIVEN the conversation settled every decision the documents need
- WHEN the agent reaches the round
- THEN it says there is nothing to ask and writes the documents

### Requirement: Check the change before presenting it

When the agent has written a change's documents, the propose instructions
shall have it check, before it tells the owner the change is written, that
each requirement has an outcome a user or caller can observe, each ADDED
or MODIFIED requirement is covered by a task's acceptance, the proposal's
scope, the change specs and the tasks do not contradict each other, and
no task does what the scope leaves out; fix what it can, tell the owner
what it changed, and ask the owner about what needs a decision.

#### Scenario: A requirement no task covers

- GIVEN the agent wrote a change whose change spec ADDs "Delete a record"
  and whose only task's acceptance covers editing
- WHEN it checks the change before presenting it
- THEN it adds the missing acceptance or task, and tells the owner it did

#### Scenario: Scope and a task disagree

- GIVEN the proposal puts the README out of scope and a task names
  `README.md`
- WHEN the agent checks the change
- THEN it asks the owner which holds, and changes nothing until they say

### Requirement: Offer one merge or two

When the agent finishes writing a change the target does not hold, the
propose instructions shall have it offer the owner both next steps:
merge the proposal now to approve it, or run `/grove-implement
<change-id>` on this branch so that one merge approves the change with
its code; and shall have it choose neither for the owner.

#### Scenario: A new change on its branch

- GIVEN the agent wrote change D on branch `rename`, and `grove check`
  reports no problem
- WHEN it finishes
- THEN it names D's files, task IDs and open questions, and offers to
  merge `rename` to approve D, or to build D with `/grove-implement D` on
  `rename` so one merge approves it with its code

#### Scenario: Revising an approved change

- GIVEN the target holds change C, and the owner revises it with
  `/grove-propose C`
- WHEN the agent finishes
- THEN it does not offer to implement C on this branch, since
  `/grove-implement C` builds C as the target holds it

### Requirement: Offer to split a change

When the agent confirms the scope of a new change whose idea is too big
for one change or has parts the owner could agree separately, the propose
instructions shall have it say so and ask; if the owner agrees, it writes
several changes on the same branch, each from its own
`grove instructions propose`, linked with `depends_on`.

#### Scenario: Two parts, one needing the other

- GIVEN the owner asks to propose `ledger delete ID` and an error instead
  of a traceback for an unknown ID, and `delete` finds its record the way
  `edit` does
- WHEN the agent confirms the scope
- THEN it says these are two changes the owner could agree separately,
  and, once the owner agrees, writes the unknown-ID change and the delete
  change on this branch, each from its own `grove instructions propose`,
  the delete change naming the other in `depends_on`

#### Scenario: The owner keeps it whole

- GIVEN the agent offered to split the same idea in two
- WHEN the owner says to keep it as one change
- THEN the agent writes one change

### Requirement: Delegate choices in a task

Where the owner agrees in the round that the implementer may make some
choices without asking, or change some paths beyond a task's `Files:`,
the propose instructions shall have the agent write them on an optional
`Delegated:` line of that task, and write no `Delegated:` line the owner
did not agree to; Grove shall read nothing on that line.

#### Scenario: The owner delegates test layout

- GIVEN the round recommends that the implementer may lay out tests and
  name private helpers as it likes, and may add `tests/test_records.py`,
  and the owner says "yes"
- WHEN the agent writes T-01
- THEN T-01 has a `Delegated:` line naming those choices and
  `tests/test_records.py`, and its `Files:` line does not name
  `tests/test_records.py`

#### Scenario: Nothing delegated

- GIVEN the round asked nothing about delegation, or the owner declined it
- WHEN the agent writes the tasks
- THEN no task has a `Delegated:` line

#### Scenario: Check ignores the line

- GIVEN a change whose T-01 has a `Delegated:` line right after its
  `Files:` line, and the target changed `tests/test_records.py` after the
  change's `base`
- WHEN the owner runs `grove check`
- THEN it reports nothing about T-01: its `Files:` line reads as before,
  and the delegated path is not watched

### Requirement: Re-plan the cycle of a revised member

When the owner runs `/grove-propose <change-id>` for a change that a cycle
not yet done lists, the propose instructions shall print that cycle's file
and have the agent, on the same branch as the revision, ask the owner
whether the change needs a person as the plan instructions do, set the
cycle's `planned` to the target commit the instructions name, and never
write or remove `unattended`.

#### Scenario: A blocked member answered

- GIVEN cycle K, marked `unattended`, lists change B, left blocked by an
  open question
- WHEN the owner runs `/grove-propose B` and answers the question
- THEN the agent revises B, asks whether B needs a person, sets K's
  `planned` to the target commit the instructions name, and leaves
  `unattended: true` as it is

#### Scenario: A member that no longer needs a person

- GIVEN cycle K's `needs_human` names C, for judging how a new screen
  looks
- WHEN the owner runs `/grove-propose C`, has its acceptance checked by a
  test, and says C no longer needs a person
- THEN the agent takes C out of `needs_human` and sets K's `planned`, and
  once the branch is squash-merged `grove run K` builds K

#### Scenario: A change in no cycle

- GIVEN change D, which no cycle not yet done lists
- WHEN the owner runs `/grove-propose D`
- THEN the instructions print no cycle file, and the agent edits no cycle

## Known discrepancies

None recorded.

## Open questions

None.
