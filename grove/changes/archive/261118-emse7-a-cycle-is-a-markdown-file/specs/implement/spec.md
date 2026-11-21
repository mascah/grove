# Changes to implement

The changes that change 261118-emse7 proposes to the requirements of
[implement](../../../../specs/implement/spec.md), as it stands at f7e0140.
This is future work: nothing here describes the target branch until the
change is implemented and the maintained spec is reconciled.

## ADDED Requirements

### Requirement: A cycle's notes are advice

Where a cycle file holds notes below its fields, the instructions
`/grove-implement` prints for that cycle, and for each of its members,
shall print the notes with the cycle file and say that where the notes
and a member's documents disagree, the member's documents hold.

#### Scenario: A note a member contradicts

- GIVEN cycle K lists A and B, K's notes say B's T-01 changes
  `ledger/cli.py`, and B's `tasks.md` names only `ledger/report.py` for
  T-01
- WHEN `/grove-implement K` is run
- THEN K's instructions and B's each print K's file with its notes and say
  a member's documents hold over them, and B's T-01 is built in
  `ledger/report.py`

#### Scenario: A cycle with no notes

- GIVEN cycle K's file holds only its fields
- WHEN `/grove-implement K` is run
- THEN the instructions print K's file and say nothing about notes
