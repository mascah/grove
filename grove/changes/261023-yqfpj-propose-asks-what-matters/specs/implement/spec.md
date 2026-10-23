# Changes to implement

The maintained spec is [grove/specs/implement/spec.md](../../../../specs/implement/spec.md).

## ADDED Requirements

### Requirement: Honor a task's delegation

When a task the owner is building has a `Delegated:` line, the implement
instructions shall have the agent make the choices it names and change the
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
