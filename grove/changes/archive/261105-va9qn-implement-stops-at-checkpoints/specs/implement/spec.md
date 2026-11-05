# Changes to implement

The changes that change 261105-va9qn proposes to the requirements of
[implement](../../../../specs/implement/spec.md), as it stands at 8a8d064.

## ADDED Requirements

### Requirement: Stop at a checkpoint

When a task's acceptance can be shown by no check the agent can run, such
as how something looks or a manual step, or the task rests on a choice
the change does not settle, including an open question that decides it,
the implement instructions shall have the agent treat the task as a
checkpoint: build its best guess and mark the task `Status: to check`,
or, where a wrong guess would be costly, such as to data, to an interface
other tasks build on, or to something hard to undo, build nothing of it
and mark it `Status: blocked`; write under it a `Check:` line saying what
to try or decide and why; leave each task that builds on it unbuilt,
naming the checkpoint; build every other task; and then, with the owner,
stop once and put every checkpoint to them, and, with `--unattended`,
finish; in both modes the report lists the checkpoints first.

#### Scenario: A screen to look at, unattended

- GIVEN cycle K lists change C, whose T-01 adds a settings screen and
  whose T-02 adds a command-line flag, and C does not say how the
  screen shows an error
- WHEN `/grove-implement K --unattended` builds C
- THEN T-01 is built with the agent's choice and is `Status: to check`
  with a `Check:` line saying what to try and why, T-02 is
  `Status: done`, C is not archived, and the report lists T-01's
  checkpoint first

#### Scenario: An open question, with the owner

- GIVEN change C's open question decides how T-02 stores records, so a
  wrong guess would change the format of `records.json`
- WHEN the owner runs `/grove-implement C`
- THEN the agent builds T-01 and T-03 without stopping, builds nothing
  of T-02 and marks it `Status: blocked` with a `Check:` line, then stops
  once and asks the owner about T-02

#### Scenario: No checkpoint

- GIVEN every task of change C has acceptance the project's tests show
- WHEN the owner runs `/grove-implement C`
- THEN the agent builds every task without stopping between them, and
  stops only to report and archive

#### Scenario: A task that builds on a checkpoint

- GIVEN T-01 of change C is `Status: to check`, and T-03 lays out its
  page inside the screen T-01 adds
- WHEN the agent reaches T-03
- THEN T-03 is not built, and `implementation.md` names T-01's checkpoint
  under it

### Requirement: Take a checkpoint's answer

When a checkpoint is answered, by the owner in the session, with
`--unattended` in the words that started the session, or by an `Answer:`
line already under the task, the implement instructions shall have the
agent write the owner's words, quoted, under the task as an `Answer:`
line if none is there, and build on it, marking a task whose answer
confirms its guess `Status: done` once its evidence holds; except that
where the answer changes what the change says, with the owner the agent
revises the change as "Settle what the change does not" says, and with
`--unattended` it leaves the task `Status: blocked`, naming the revision
the change needs.

#### Scenario: Confirmed in a short session

- GIVEN T-01 of change C is `Status: to check` on branch `c-impl`, and
  every other task is done
- WHEN the owner runs `/grove-implement C` on `c-impl`, tries the screen
  and says "the inline error is fine"
- THEN the agent writes `Answer: "the inline error is fine"` under T-01,
  marks it `Status: done`, and archives C

#### Scenario: An answer in the words that start an unattended session

- GIVEN T-02 of change C is `Status: blocked` with a `Check:` line asking
  which ID format records use
- WHEN an unattended session is started on the branch with the owner's
  words "use six-digit IDs" beside `/grove-implement C --unattended`
- THEN T-02 has `Answer: "use six-digit IDs"` under it, is built to it,
  and is `Status: done`

#### Scenario: An answer that changes the change, unattended

- GIVEN T-01 of change C is `Status: to check`, and the owner's answer
  asks for an undo button C's change spec does not state
- WHEN an unattended session takes the answer
- THEN T-01 is `Status: blocked`, naming the revision C needs, and C's
  proposal, change specs and tasks are unchanged

## MODIFIED Requirements

### Requirement: Settle what the change does not

When a task cannot be built as the change says, or shows the change's
approach is wrong, the implement instructions shall have the agent, with
the owner, stop that task at once and ask, naming which of the proposal,
change specs and tasks it would edit and how, and edit them on the
branch only once the owner agrees, by following
`grove instructions propose <change-id>` and then running `grove check`,
marking the task `Status: blocked` with the question if the owner leaves
it open; and, with `--unattended`, mark the task `Status: blocked` with
the reason, ask no one, and edit none of the change's files but
`implementation.md`. With the owner, they shall also have it set a done
task whose result an agreed edit changes back to in progress until it is
rebuilt, and stop when `grove check` reports a requirement of the change
as rewritten or removed since its base, telling the owner the change
must be revised with propose and merged again first.

#### Scenario: Feedback changes the change

- GIVEN the agent built T-01 of change C as a checkpoint and told the
  owner what to try, and the owner, trying it, asks for behavior C's
  change spec does not state
- WHEN the owner agrees to the edit the agent describes
- THEN the agent edits C's change spec on the branch by following
  `grove instructions propose C`, T-01 is in progress until rebuilt, and
  the agent builds to the edit; had the owner declined, C would be
  unchanged and built as it says

#### Scenario: The approach rests on something false

- GIVEN C's T-01 says `find` raises `ValueError` for an unknown ID, and
  the agent finds `report` relies on `find` returning nothing
- WHEN it makes the finding while building T-01
- THEN it stops T-01 at once and tells the owner what it found and which
  of C's documents it would edit and how

#### Scenario: Unattended

- GIVEN cycle K lists change C, and building T-01 the agent finds
  `report` relies on `find` returning nothing, where T-01 says `find`
  raises `ValueError`
- WHEN `/grove-implement K --unattended` reaches T-01
- THEN T-01 is `Status: blocked` with the reason, nobody is asked, and
  C's proposal, change specs and tasks are unchanged
