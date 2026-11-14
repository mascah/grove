# Changes to propose

The changes that change 261114-yx2b7 proposes to the requirements of
[propose](../../../../specs/propose/spec.md), as it stands at 999ab56. This
is future work: nothing here describes the target branch until the change
is implemented and the maintained spec is reconciled.

## MODIFIED Requirements

### Requirement: Re-plan the cycle of a revised member

When the owner runs `/grove-propose <change-id>` for a change that a cycle
not yet done lists, the propose instructions shall print that cycle's file
and have the agent, on the same branch as the revision, set the cycle's
`planned` to the target commit the instructions name, change nothing else
in it, and never write or remove `unattended`.

#### Scenario: A blocked member answered

- GIVEN cycle K, marked `unattended`, lists change B, left blocked by an
  open question
- WHEN the owner runs `/grove-propose B` and answers the question
- THEN the agent revises B, sets K's `planned` to the target commit the
  instructions name, asks nothing about whether B needs a person, and
  leaves `unattended: true` as it is

#### Scenario: A change in no cycle

- GIVEN change D, which no cycle not yet done lists
- WHEN the owner runs `/grove-propose D`
- THEN the instructions print no cycle file, and the agent edits no cycle
