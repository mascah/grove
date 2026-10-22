# Changes to implement

There is no maintained spec for `implement` at 7ce7b91; 261022-d1y6d, which
this change waits on, adds it. The requirement below is ADDED to it.

## ADDED Requirements

### Requirement: Implement a cycle

When the owner runs `/grove-implement <cycle-id>` for a cycle the target
holds, the system shall print instructions for building each member not yet
archived, in build order, each as `/grove-implement <change-id>` builds it,
in the owner's session.

#### Scenario: A cycle the owner builds

- GIVEN cycle K lists A and B, is not marked `unattended`, and
  `needs_human` names B
- WHEN the owner runs `/grove-implement K`
- THEN A and then B are built in the owner's session, the agent asking the
  owner where `/grove-implement` would

## Observed current behavior

- At 7ce7b91 there is no implement skill and there are no cycles.
