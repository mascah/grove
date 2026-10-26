# Changes to implement

The maintained spec is [grove/specs/implement/spec.md](../../../../specs/implement/spec.md),
as it stands at ea7fbaf.

## ADDED Requirements

### Requirement: Recommend re-basing after a dependency lands

When `/grove-implement <change-id>`, with the owner, names a task's file
as changed since the change's `base`, and a change in its `depends_on` is
archived on the target, the implement instructions shall have the agent
recommend revising the change with propose, which re-bases it, as the
normal step, beside building on the file as it is now.

#### Scenario: Built after its dependency

- GIVEN change B depends on A, A is archived on the target, and the merge
  that archived A changed `ledger/export.py`, which B's T-01 names
- WHEN the owner runs `/grove-implement B` and the agent reaches T-01
- THEN it tells the owner what changed and recommends revising B with
  `/grove-propose B` to re-base it, or else building on the file as it is
  now, and waits for the answer
