# Changes to implement

The changes that change 261101-mg8hd proposes to the requirements of
[implement](../../../../specs/implement/spec.md), as it stands at 1a24167.

## MODIFIED Requirements

### Requirement: Change only the files the tasks name

The implement instructions shall have the agent change only the files the
tasks' `Files:` lines name, make the choices and change the paths a task's
`Delegated:` line names without asking, and, for any other file a task
needs changed to be built as the change says or makes false, with the
owner, tell the owner and, only with their agreement, edit it and write
`` Also changed: `<path>` — <reason> `` under that task in
`implementation.md`; with `--unattended`, edit it without marking the task
blocked, name it, and why, in that task's evidence, and write no
`Also changed:` line; in neither mode the change's documents other than
`implementation.md`, nor `grove.yaml`.

#### Scenario: A delegated path

- GIVEN T-01's `Files:` line names `tests/test_cli.py` and its
  `Delegated:` line names `tests/test_records.py`
- WHEN the agent puts `find`'s tests in a new `tests/test_records.py`
- THEN it creates the file without asking the owner

#### Scenario: A file nothing names

- GIVEN the agent, building T-01, finds `README.md`, which no line of T-01
  names, states what T-01 makes false
- WHEN it tells the owner and they agree
- THEN it edits `README.md` and writes `` Also changed: `README.md` ``
  under T-01 with the reason; had they declined, nothing is written

#### Scenario: Unattended

- GIVEN the same T-01, built by `/grove-implement K --unattended`
- WHEN the agent builds it
- THEN it creates `tests/test_records.py` and edits `README.md`, T-01's
  evidence names `README.md` and why, there is no `Also changed:` line,
  and `grove review` marks `README.md` `named nowhere`

#### Scenario: Unattended, a file the task needs

- GIVEN T-01 names only `ledger/cli.py` and `tests/test_cli.py`, and its
  `Interfaces:` says `find` in `ledger/records.py` raises `ValueError` for
  an unknown ID
- WHEN `/grove-implement K --unattended` builds it
- THEN `ledger/records.py` is changed, T-01 is `Status: done` with
  evidence naming `ledger/records.py` and why, and T-01 is not blocked
