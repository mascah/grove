# Changes to propose

The maintained spec is [grove/specs/propose/spec.md](../../../../specs/propose/spec.md),
as it stands at ea7fbaf.

## ADDED Requirements

### Requirement: Say that revising re-bases

When `grove check` reports a task's named file as changed on the target
since the change's `base`, the system shall say, in the same sentence,
that revising the change with propose also re-bases it after a dependency
lands.

#### Scenario: The message

- GIVEN change C, in no cycle, proposed at `base` b143991, whose T-01
  names `ledger/cli.py`, which the target has since changed
- WHEN the owner runs `grove check`
- THEN it reports "T-01's ledger/cli.py changed on main since base
  b143991; revise the change with propose, which also re-bases it after a
  dependency lands"

### Requirement: Re-base a change whose named file moved

When the owner revises, with `/grove-propose <change-id>`, a change that
`grove check` reports for a named file changed since its `base`, the
propose instructions shall have the agent bring the change up to date with
the file as the target holds it and set `base` to the target's commit.

#### Scenario: A dependency landed

- GIVEN change B depends on change A, and the merge that archived A
  created `ledger/export.py`, which B's T-01 names
- WHEN the owner runs `/grove-propose B`
- THEN B's `base` is the target's commit, B's tasks hold for
  `ledger/export.py` as the target has it, and `grove check` reports
  nothing for B
