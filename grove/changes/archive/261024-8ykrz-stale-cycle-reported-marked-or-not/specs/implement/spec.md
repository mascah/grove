# Changes to implement

The maintained spec is [grove/specs/implement/spec.md](../../../../specs/implement/spec.md),
at 2d07433. Which commits count since a cycle's plan is
"Count a member changed with its plan" in
[grove/specs/cycles/spec.md](../../../../specs/cycles/spec.md).

## ADDED Requirements

### Requirement: Tell the owner of a file moved since the plan

When `grove instructions implement <change-id>` is run, with the owner or
with `--unattended`, for a member of a cycle that the target holds and
that is not done, and a task names a file that a commit counted by "Count
a member changed with its plan" changed after the cycle's `planned`
commit, the system shall name that task and file as changed since the
cycle's plan.

#### Scenario: With the owner

- GIVEN change C is a member of cycle K, not marked `unattended` and not
  done, C's T-01 names `ledger/records.py`, and the target merged a commit
  touching that file after K's `planned` commit
- WHEN the owner runs `/grove-implement C`
- THEN Grove's instructions name T-01 and `ledger/records.py` as changed
  on the target since K's plan, and the agent, on reaching T-01, tells the
  owner what changed and asks before building it

#### Scenario: Unattended

- GIVEN the same C and K
- WHEN `/grove-implement K --unattended` is run and reaches C
- THEN Grove's instructions for C name T-01 and `ledger/records.py` as
  changed since K's plan, and T-01 is marked `Status: blocked` naming that
  file, unbuilt, with nobody asked

#### Scenario: The cycle it depends on lands

- GIVEN cycle K depends on cycle L, K's member C has T-01 naming
  `ledger/records.py`, and the merge into the target that archives L's
  change also edits that file
- WHEN `grove instructions implement C` is run
- THEN it names no file as changed since K's plan

#### Scenario: Re-planned

- GIVEN C's T-01 names `ledger/records.py`, the target changed that file
  after K's `planned` commit, and K was then re-planned, setting `planned`
  to the target's commit
- WHEN `grove instructions implement C` is run
- THEN it names no file as changed since K's plan
