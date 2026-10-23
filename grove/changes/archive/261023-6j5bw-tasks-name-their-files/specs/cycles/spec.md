# Changes to cycles

The changes that change 261023-6j5bw proposes to the requirements of
`cycles` (`grove/specs/cycles/spec.md` at b143991). The Purpose's
definition of stale becomes: a cycle is stale when the folder of a member
not yet archived, or a file one of its tasks names, differs on the target
from the same path at the cycle's `planned` commit.

## MODIFIED Requirements

### Requirement: Report a cycle that cannot run unattended

When the owner runs `grove check`, the system shall report a cycle file
that holds `unattended: true` while `needs_human` names a member, and a
cycle file that holds `unattended: true` while the cycle is stale, each
naming what to do.

#### Scenario: Marked with a change that needs a person

- GIVEN cycle K holds `unattended: true` and `needs_human` names C
- WHEN the owner runs `grove check`
- THEN it reports K's `unattended` line, saying C needs a person

#### Scenario: Marked and stale

- GIVEN cycle K holds `unattended: true`, and member A's proposal was
  revised on the target after K's `planned` commit
- WHEN the owner runs `grove check`
- THEN it reports K, saying A changed since the plan and to re-plan K with
  `/grove-plan K`

#### Scenario: Marked and a member's file moved

- GIVEN cycle K holds `unattended: true`, member B's T-01 names
  `ledger/records.py`, and the target merged a commit touching that file
  after K's `planned` commit
- WHEN the owner runs `grove check`
- THEN it reports K, saying B changed since the plan and to re-plan K with
  `/grove-plan K`, and `grove run` does not build K

### Requirement: Revise a cycle

When the owner runs `/grove-plan <cycle-id>`, the system shall print
instructions for revising that cycle file as the owner asks and setting
`planned` to the target commit the revision read.

#### Scenario: Re-plan a stale cycle

- GIVEN cycle K is stale because member A was revised on the target
- WHEN the owner runs `/grove-plan K`
- THEN the agent reads A as the target holds it, asks about it as for a new
  plan, and sets `planned` to the target's commit

#### Scenario: Re-plan after a member's file moved

- GIVEN cycle K is stale because a file member B's T-01 names changed on
  the target, and the owner is content to build B on the new text
- WHEN the owner runs `/grove-plan K` and merges the cycle file
- THEN `planned` is the target's commit, K is no longer stale, and B's
  files are unchanged
