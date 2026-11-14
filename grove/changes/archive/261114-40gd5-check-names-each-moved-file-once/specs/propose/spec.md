# Changes to propose

<!-- The changes that change 261114-40gd5 proposes to the requirements of
`propose` (`grove/specs/propose/spec.md` at 999ab56). -->

## MODIFIED Requirements

### Requirement: Report a named file that moved

When the owner runs `grove check`, the system shall report, once for each
change in no cycle that is not done, each file a task names, or a file
under a folder a task names, that differs on the target from the same file
at the change's `base`, naming each such file once and saying to revise
the change with propose.

#### Scenario: A file changed since base

- GIVEN change C, in no cycle, proposed at `base` b143991, with T-01 naming
  `internal/cli/run.go`, and the target has since merged a commit touching
  that file
- WHEN the owner runs `grove check`
- THEN it reports C's `tasks.md` once, saying `internal/cli/run.go`
  changed since base b143991 and to revise the change with propose

#### Scenario: A file under a named folder

- GIVEN change C, in no cycle, with T-01 naming `tests/fixtures/`, and the
  target has since merged a commit touching `tests/fixtures/records.json`
- WHEN the owner runs `grove check`
- THEN it reports C's `tasks.md`, saying `tests/fixtures/records.json`
  changed since base and to revise the change with propose

#### Scenario: Several files, several tasks

- GIVEN change C, in no cycle, with T-01 naming `ledger/cli.py` and
  `tests/test_cli.py` and T-02 naming `ledger/cli.py`, and the target has
  since changed both files
- WHEN the owner runs `grove check`
- THEN it reports one problem, on C's `tasks.md`, naming
  `ledger/cli.py` and `tests/test_cli.py` once each

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

### Requirement: Say that revising re-bases

When `grove check` reports a change's named files as changed on the target
since the change's `base`, the system shall say, in the same sentence,
that revising the change with propose also re-bases it after a dependency
lands.

#### Scenario: The message

- GIVEN change C, in no cycle, proposed at `base` b143991, whose T-01
  names `ledger/cli.py`, which the target has since changed
- WHEN the owner runs `grove check`
- THEN it reports "ledger/cli.py changed on main since base b143991;
  revise the change with propose, which also re-bases it after a
  dependency lands"
