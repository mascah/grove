# (c) A needed file no task names: T-01 changes find in ledger/records.py,
# which its Files: line leaves out. Built with the owner, then unattended.

setup() {
  local c=grove/changes/261015-c0000-unknown-id-prints-an-error
  id=261015-c0000
  put $c/.grove.yaml << 'EOF'
title: Editing an unknown ID prints an error
created: 2026-10-15T09:00:00Z
schema: default
base: @BASE@
EOF
  put $c/proposal.md << 'EOF'
# Editing an unknown ID prints an error

## Problem

`ledger edit 99 --note x` ends in an uncaught `StopIteration` traceback:
`find` in `ledger/records.py` uses `next(...)` with no default.

## Intended outcome

Editing an ID no record has prints `error: no record 99` and exits 2,
like the other errors.

## Scope

### In scope

- `edit` with an unknown ID.

### Out of scope

- `delete`, which does not exist yet.

## Decisions

- The owner: "exit 2, like the other errors".

## Open questions

None.

## Related documents

- Requirements: [specs/record-editing/spec.md](specs/record-editing/spec.md)
- Tasks: [tasks.md](tasks.md)
EOF
  put $c/specs/record-editing/spec.md << 'EOF'
# Changes to record-editing

There is no maintained spec for `record-editing` at @BASE@.

## ADDED Requirements

### Requirement: Report an unknown ID

If a user edits an ID that no record has, then the system shall print
`error: no record <id>` and exit 2.

#### Scenario: Edit record 99

- GIVEN the six records of `records.json`
- WHEN the user runs `ledger edit 99 --note x`
- THEN it prints `error: no record 99` and exits 2

## Observed current behavior

- `ledger edit 99 --note x` ends in a `StopIteration` traceback
  (`ledger/records.py`, `find`, at @BASE@).
EOF
  put $c/tasks.md << 'EOF'
## T-01 — An error for an unknown ID

Result: editing an unknown ID prints an error and exits 2, as the change
spec says.

Files: `ledger/cli.py`, `tests/test_cli.py`

Interfaces: `find` in `ledger/records.py` raises `ValueError` for an
unknown ID.

Acceptance:
- `ledger edit 99 --note x` prints `error: no record 99` and exits 2.
- A test in `tests/test_cli.py` shows it.
EOF
  commit "Propose an error for an unknown ID"
}

runs=(
  "unknown-id-impl => /grove-implement {id}"
  "unknown-id-unattended => /grove-implement {id} --unattended"
)
answers=(
  "records\.py => Yes, add ledger/records.py to T-01 and edit it."
)
