# Changes to runs

The maintained spec is [grove/specs/runs/spec.md](../../../../specs/runs/spec.md).
"Retry" starts from the text change 261023-g0a3h gives it, which this
change depends on; only the scenario's GIVEN differs.

## MODIFIED Requirements

### Requirement: Retry

When the operator runs `grove run --retry` in a clean checkout, the system
shall start a new attempt as an unpaused `grove run` does, whatever the
latest attempt selected.

#### Scenario: A blocked change revised

- GIVEN cycle K's A was merged as finished, B was blocked by an open
  question, and the owner merged, in one squash merge, a revision of B
  answering it that re-planned K
- WHEN the operator runs `grove run --retry`
- THEN a new branch is created from the fetched target and a session is
  launched on the implement skill with K's ID and `--unattended`, which
  builds B
