# Changes to implement

The changes that change 261122-z1kh0 proposes to the requirements of
[implement](../../../../specs/implement/spec.md), as it stands at 15816d7.
This is future work: nothing here describes the target branch until the
change is implemented and the maintained spec is reconciled.

## ADDED Requirements

### Requirement: Record work found outside the change

When the agent, building a task, sees work that the change does not ask
for and that no task needs, such as a bug, a false document, a missing
test, or code that differs from a maintained spec, the implement
instructions shall have the agent change nothing for it, write under that
task in `implementation.md`
`Found: <what is wrong or missing, where, and how it was seen>`, and list
each `Found:` line in its report after every `Check:` line, each blocked
task's question and each unverified scenario, and before what was done;
and the system shall read nothing from a `Found:` line.

#### Scenario: A bug next to the work

- GIVEN change C, whose T-01 adds `ledger delete ID` and names only
  `ledger/cli.py` and `tests/test_cli.py`, and `ledger edit` with an
  unknown ID prints a `StopIteration` traceback, which C does not mention
- WHEN `/grove-implement C` builds T-01 and the agent sees that traceback
- THEN `edit`'s code is unchanged, T-01 is `Status: done` with a `Found:`
  line saying `ledger edit` prints a traceback for an unknown ID and how
  the agent saw it, and the report lists that line after T-01's `Check:`
  lines and before what was done

#### Scenario: A file the task makes false

- GIVEN T-01 changes what `ledger delete` prints, and the README's usage
  example would then be false
- WHEN `/grove-implement C` builds T-01
- THEN the agent edits `README.md` and writes an `Also changed:` line for
  it under T-01, and writes no `Found:` line for it

#### Scenario: Nothing found

- GIVEN the agent sees no work outside change C
- WHEN `/grove-implement C` builds it
- THEN `implementation.md` has no `Found:` line, and the report names
  none

#### Scenario: Grove does not read it

- GIVEN change C, archived on the target with a `Found:` line under T-01
- WHEN the owner runs `grove changes --json` and `grove review`
- THEN their output is what it would be without that line
