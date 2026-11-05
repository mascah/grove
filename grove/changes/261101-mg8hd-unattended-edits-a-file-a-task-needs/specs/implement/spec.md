# Changes to implement

The changes that change 261101-mg8hd proposes to the requirements of
[implement](../../../../specs/implement/spec.md), as it stands at 6090cb3.

## MODIFIED Requirements

### Requirement: Change only the files the tasks name

The implement instructions shall have the agent change the files the
tasks' `Files:` lines name, make the choices and change the paths a
task's `Delegated:` line names, and, for any other file a task needs
changed to be built as the change says or makes false, edit it without
asking anyone and without marking the task blocked, and write
`` Also changed: `<path>` — <why> `` under that task in
`implementation.md`, with the owner and with `--unattended` alike; and
edit none of the change's documents other than `implementation.md`, nor
`grove.yaml`.

#### Scenario: A delegated path

- GIVEN T-01's `Files:` line names `tests/test_cli.py` and its
  `Delegated:` line names `tests/test_records.py`
- WHEN the agent puts `find`'s tests in a new `tests/test_records.py`
- THEN it creates the file without asking the owner, with the owner or
  unattended

#### Scenario: A file nothing names

- GIVEN the agent, building T-01 with the owner, finds `README.md`, which
  no line of T-01 names, states what T-01 makes false
- WHEN it builds T-01
- THEN it edits `README.md` without asking and writes
  `` Also changed: `README.md` — T-01 makes its usage example false ``
  under T-01

#### Scenario: Unattended

- GIVEN the same T-01, built by `/grove-implement K --unattended`
- WHEN the agent builds it
- THEN it creates `tests/test_records.py`, edits `README.md` and writes
  its `Also changed:` line with why, and `grove review` shows `README.md`
  as also changed by T-01 with that reason

#### Scenario: Unattended, a file the task needs

- GIVEN T-01 names only `ledger/cli.py` and `tests/test_cli.py`, and its
  `Interfaces:` says `find` in `ledger/records.py` raises `ValueError` for
  an unknown ID
- WHEN `/grove-implement K --unattended` builds it
- THEN `ledger/records.py` is changed, T-01 is `Status: done` with an
  `Also changed:` line naming `ledger/records.py` and why, and T-01 is not
  blocked

### Requirement: Review names what the branch approves and builds

When the current branch holds a change the target does not hold, `grove
review` shall say whether merging the branch approves it as future work
or, once it is archived as finished on the branch, with its code; and when
the branch holds a change open, or archived where the target does not
hold it archived, `grove review` shall list each file outside the
document root that the branch changes against the target, as planned,
naming the task whose `Files:` line names it or a folder above it, as
also changed, naming the task and the reason its `Also changed:` line in
`implementation.md` gives, or as named nowhere.

#### Scenario: A draft, before and after building

- GIVEN branch `rename` holds change D, which the target does not hold
- WHEN the owner runs `grove review` on `rename` with D open, and again
  once D is archived as finished there
- THEN, with D open, it says merging `rename` approves change D as future
  work and does not implement it; with D archived, its first line is
  "Merging rename into main approves change D with its code."

#### Scenario: Planned, also changed and named nowhere

- GIVEN branch `c-impl` holds approved change C, open or archived, whose
  T-01 names `ledger/records.py` and `tests/fixtures/` on its `Files:`
  line, whose `implementation.md` has `` Also changed: `README.md` — T-01
  makes its usage example false `` under T-01, and the branch changes
  `ledger/records.py`, `tests/fixtures/records.json`, `README.md` and
  `Cargo.lock`
- WHEN the owner runs `grove review` on `c-impl`
- THEN it names no change approved by the merge, and lists
  `ledger/records.py` and `tests/fixtures/records.json` as planned by
  C:T-01, `README.md` as also changed by C:T-01 with that reason, and
  `Cargo.lock` as named nowhere
