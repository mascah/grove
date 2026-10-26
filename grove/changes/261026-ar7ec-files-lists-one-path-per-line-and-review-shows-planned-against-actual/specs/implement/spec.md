# Changes to implement

The maintained spec is [grove/specs/implement/spec.md](../../../../specs/implement/spec.md),
as it stands at ea7fbaf.

## ADDED Requirements

### Requirement: Review shows planned against actual

When the current branch holds a change open, or holds it archived where
the target does not hold it archived, `grove review` shall list each file
outside the document root that the branch changes against the target,
each as planned, naming the change's task whose `Files:` line names it or
a folder above it, as agreed, naming the task and the reason its
`Also changed:` line in `implementation.md` gives, or as named nowhere.

#### Scenario: Planned, agreed and named nowhere

- GIVEN branch `c-impl` holds change C, whose T-01 names
  `ledger/records.py` and `tests/fixtures/` on its `Files:` line, whose
  `implementation.md` has `` Also changed: `README.md` — T-01 makes its
  usage example false `` under T-01, and the branch changes
  `ledger/records.py`, `tests/fixtures/records.json`, `README.md` and
  `Cargo.lock`
- WHEN the owner runs `grove review` on `c-impl`
- THEN it lists `ledger/records.py` and `tests/fixtures/records.json` as
  planned by C:T-01, `README.md` as agreed by C:T-01 with that reason, and
  `Cargo.lock` as named nowhere, and exits 0

#### Scenario: After archiving

- GIVEN the same branch after `grove archive C`, the target holding C open
- WHEN the owner runs `grove review`
- THEN it lists the same files the same way

#### Scenario: No change on the branch

- GIVEN a branch that holds no change open, and none archived that the
  target does not hold archived
- WHEN the owner runs `grove review`
- THEN it lists no files against a plan

## MODIFIED Requirements

### Requirement: Revise the change on its branch

When the owner's answer or feedback changes what the change says, the
implement instructions shall have the agent say which of the change's
proposal, change specs and tasks it would edit and how, and, only once the
owner agrees, edit them on the implementation branch, leaving every
`Files:` line as it is, revising a task's files only through
`grove instructions propose <change-id>`, run `grove check` clean, and
then build to the revised change.

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

#### Scenario: A revision that needs a new task

- GIVEN the owner agrees that change C needs a new task touching
  `ledger/export.py`
- WHEN the agent revises C
- THEN it runs `grove instructions propose C` and follows it to add the
  task with its `Files:` line, instead of writing the line itself

### Requirement: Tell the owner of a moved file

When the owner runs `/grove-implement <change-id>` for a change in no
cycle that is not done, and a task names a file, or a folder holding a
file, that differs on the target from the same file at the change's
`base`, the system shall say which task and file, and the implement
instructions shall have the agent ask the owner before building that
task, instead of building on the file as it was.

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
- THEN `implementation.md` has `` Also changed: `docs/concepts.md` `` under
  T-01 with the reason, `tasks.md` is unchanged, and the agent edits the
  page; had the owner declined, the page would be left as it is and no
  line written

#### Scenario: A member of a cycle

- GIVEN change C is a member of cycle K, and a file C's T-01 names changed
  on the target after K's `planned` commit
- WHEN the owner runs `/grove-implement K`
- THEN C is not reported against its `base`; `grove check` reports K as
  stale, and the owner re-plans K or builds on

### Requirement: Ignore a task's delegation unattended

When a task built with `/grove-implement <id> --unattended` has a
`Delegated:` line, the implement instructions shall have the agent ignore
it: change only the files the tasks' `Files:` lines name, write no
`Also changed:` line, and leave a path only a `Delegated:` line names as it
is, as it does any file no task names: when the task would have changed
it, naming it, and why, in that task's evidence in `implementation.md`.

#### Scenario: A delegated path, unattended

- GIVEN cycle K lists change C, whose T-01 names `ledger/records.py` and
  `tests/test_cli.py` on its `Files:` line and `tests/test_records.py` on
  its `Delegated:` line
- WHEN `/grove-implement K --unattended` builds T-01
- THEN `tests/test_records.py` is not created, `find`'s tests go in
  `tests/test_cli.py`, T-01's evidence names `tests/test_records.py`
  and why it was left, and `implementation.md` has no `Also changed:` line
