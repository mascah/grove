# Changes to cycles

The changes that change 261105-gmpa4 proposes to the requirements of
[cycles](../../../../specs/cycles/spec.md), as it stands at 6090cb3.

## MODIFIED Requirements

### Requirement: Plan a cycle

When the owner runs `/grove-plan <description>`, the system shall print
instructions for writing a new cycle file holding a `title`, an optional
one-line `goal`, the approved changes the owner agrees to as `changes`, the
time Grove gives the plan as `created`, and the target commit the plan read
as `planned`, editing no change's files except, for each member `grove
check` reports stale on a requirement or for a named file changed since
its base, telling the owner that this edits approved work and bringing
that member up to date on the same branch as the propose instructions for
it say, before writing the cycle file.

#### Scenario: Two approved changes

- GIVEN the target approves changes A and B, neither archived, and
  `grove check` reports neither
- WHEN the owner runs `/grove-plan 0.4 release` and agrees to A and B
- THEN a new cycle file lists A and B under `changes`, `created` is the
  time Grove gave the plan, `planned` is the target's commit, and A's and
  B's files are unchanged

#### Scenario: A draft change named

- GIVEN change D exists in the working tree but not in the target
- WHEN the owner asks for D in the cycle
- THEN D is not listed, and the agent says D is not approved

#### Scenario: A member stale on its base

- GIVEN the target approves change B, whose T-01 names `ledger/report.py`,
  and the merge that archived B's dependency A changed `ledger/report.py`
  after B's `base`
- WHEN the owner runs `/grove-plan 0.4 release` and agrees to B
- THEN the agent says bringing B up to date edits approved work, brings
  T-01 up to date with `ledger/report.py` as the target has it, sets B's
  `base` to the target's commit and writes the cycle file, all on the same
  branch, without stopping for another session, and `grove check` is clean
