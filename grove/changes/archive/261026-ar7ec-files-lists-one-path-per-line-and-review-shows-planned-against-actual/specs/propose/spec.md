# Changes to propose

The changes that change 261026-ar7ec proposes to the requirements of
[propose](../../../../specs/propose/spec.md), as they stand at ea7fbaf.

## ADDED Requirements

### Requirement: Report a malformed Files line

If a task's `Files:` line holds anything but repo-relative paths in
backticks, one per bullet or separated by commas, then `grove check` shall
report it at its line in `tasks.md`, naming the task and both forms and
saying that a folder ends in `/`, and Grove shall read no file from that
line.

#### Scenario: A path outside the repository

- GIVEN T-01's `Files:` line names `` `../ledger/cli.py` ``
- WHEN the owner runs `grove check`
- THEN it reports the line, saying T-01's `Files:` line should list
  repo-relative paths in backticks, one per bullet or separated by commas,
  and that a folder ends in `/`, and exits non-zero

#### Scenario: A folder

- GIVEN T-01's `Files:` line names `` `tests/fixtures/` ``
- WHEN the owner runs `grove check`
- THEN it reports nothing for the line

## MODIFIED Requirements

### Requirement: Tasks name files and interfaces

When a change's tasks are written, each task shall name, on a `Files:`
line, the repo-relative paths it creates, modifies or tests, one
backticked path per bullet or separated by commas, a path ending in `/`
naming a folder and every file under it, and, on an `Interfaces:` line,
the functions, types or commands it produces that other tasks or callers
use, with no code, commands or steps.

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

#### Scenario: One path per bullet

- GIVEN a task touching `ledger/records.py`, `tests/test_cli.py` and the
  fixture folder `tests/fixtures/`
- WHEN the agent writes the task
- THEN `Files:` is followed by three bullets, each one backticked path,
  the last `tests/fixtures/`, and `grove check` reports nothing for it

#### Scenario: A comma line still reads

- GIVEN a change written before this one, whose T-01 has
  `` Files: `ledger/records.py`, `tests/test_cli.py` ``
- WHEN the owner runs `grove check`
- THEN it reports nothing for the line, and Grove reads both files as
  T-01's

### Requirement: Report a named file that moved

When the owner runs `grove check`, the system shall report, for each
change in no cycle that is not done, each task whose named file, or a file
under a folder it names, differs on the target from the same file at the
change's `base`, naming the task and the file and saying to revise the
change with propose.

#### Scenario: A file changed since base

- GIVEN change C, in no cycle, proposed at `base` b143991, with T-01 naming
  `internal/cli/run.go`, and the target has since merged a commit touching
  that file
- WHEN the owner runs `grove check`
- THEN it reports C's `tasks.md`, saying T-01's `internal/cli/run.go`
  changed since base b143991 and to revise the change with propose

#### Scenario: A file under a named folder

- GIVEN change C, in no cycle, with T-01 naming `tests/fixtures/`, and the
  target has since merged a commit touching `tests/fixtures/records.json`
- WHEN the owner runs `grove check`
- THEN it reports C's `tasks.md`, saying T-01's
  `tests/fixtures/records.json` changed since base and to revise the
  change with propose

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
