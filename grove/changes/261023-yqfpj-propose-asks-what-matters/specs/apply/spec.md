# Changes to apply

The maintained spec is [grove/specs/apply/spec.md](../../../../specs/apply/spec.md).

## ADDED Requirements

### Requirement: Ignore a task's delegation

When a member's task has a `Delegated:` line, the apply instructions shall
have the agent ignore it: change only the files the tasks' `Files:` lines
name, and leave a path only a `Delegated:` line names as it is, as it
does any file no task names: when the task would have changed it, naming
it, and why, in that task's evidence in `implementation.md`.

#### Scenario: A delegated path under apply

- GIVEN cycle K lists change C, whose T-01 names `ledger/records.py` and
  `tests/test_cli.py` on its `Files:` line and `tests/test_records.py` on
  its `Delegated:` line
- WHEN `/grove-apply K` builds T-01
- THEN `tests/test_records.py` is not created, `find`'s tests go in
  `tests/test_cli.py`, and T-01's evidence names `tests/test_records.py`
  and why it was left
