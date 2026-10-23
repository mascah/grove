# Changes to propose

The maintained spec is [grove/specs/propose/spec.md](../../../../specs/propose/spec.md).

## ADDED Requirements

### Requirement: Re-plan the cycle of a revised member

When the owner runs `/grove-propose <change-id>` for a change that a cycle
not yet done lists, the propose instructions shall print that cycle's file
and have the agent, on the same branch as the revision, ask the owner
whether the change needs a person as the plan instructions do, set the
cycle's `planned` to the target commit the instructions name, and never
write or remove `unattended`.

#### Scenario: A blocked member answered

- GIVEN cycle K, marked `unattended`, lists change B, left blocked by an
  open question
- WHEN the owner runs `/grove-propose B` and answers the question
- THEN the agent revises B, asks whether B needs a person, sets K's
  `planned` to the target commit the instructions name, and leaves
  `unattended: true` as it is

#### Scenario: A member that no longer needs a person

- GIVEN cycle K's `needs_human` names C, for judging how a new screen
  looks
- WHEN the owner runs `/grove-propose C`, has its acceptance checked by a
  test, and says C no longer needs a person
- THEN the agent takes C out of `needs_human` and sets K's `planned`, and
  once the branch is squash-merged `grove run K` builds K

#### Scenario: A change in no cycle

- GIVEN change D, which no cycle not yet done lists
- WHEN the owner runs `/grove-propose D`
- THEN the instructions print no cycle file, and the agent edits no cycle
