# Propose

## Purpose

`/grove-propose` turns a chosen direction into a change the owner can
review: a proposal, change specs and tasks. This covers how it settles
what those documents say, what they must say, and what `grove check`
reports on them. A task's named
files are the repo-relative paths on its `Files:` line.

## Boundaries

This covers the documents `/grove-propose` writes and `grove check`'s
reports on a change. Building a change belongs to
`grove/specs/implement/spec.md` and `grove/specs/apply/spec.md`, and a
cycle's staleness to `grove/specs/cycles/spec.md`.

## Requirements

### Requirement: Tasks name files and interfaces

When a change's tasks are written, each task shall name, on a `Files:`
line, the repo-relative paths it creates, modifies or tests, and, on an
`Interfaces:` line, the functions, types or commands it produces that
other tasks or callers use, with no code, commands or steps.

#### Scenario: A task that fixes a function

- GIVEN the owner decided, while proposing, that `find` in
  `ledger/records.py` raises `ValueError` for an unknown ID
- WHEN the agent writes the task
- THEN its `Files:` line names `ledger/records.py` and `tests/test_cli.py`,
  and its `Interfaces:` line says `find` raises `ValueError` with the
  message the requirement states

#### Scenario: A task that updates documentation

- GIVEN the owner wants the README's planned list updated in this change
- WHEN the agent writes the task
- THEN its `Files:` line names `README.md`, and the task has no steps

### Requirement: Report a named file that moved

When the owner runs `grove check`, the system shall report, for each
change in no cycle that is not done, each task whose named file differs on
the target from the same file at the change's `base`, naming the task and
the file and saying to revise the change with propose.

#### Scenario: A file changed since base

- GIVEN change C, in no cycle, proposed at `base` b143991, with T-01 naming
  `internal/cli/run.go`, and the target has since merged a commit touching
  that file
- WHEN the owner runs `grove check`
- THEN it reports C's `tasks.md`, saying T-01's `internal/cli/run.go`
  changed since base b143991 and to revise the change with propose

#### Scenario: A member of a cycle

- GIVEN change C is a member of cycle K, not done, and T-01's named file
  changed on the target after K's `planned` commit
- WHEN the owner runs `grove check`
- THEN C is not reported against its `base`, and K is reported as stale
  (`grove/specs/cycles/spec.md`)

#### Scenario: A file this change's own branch touched

- GIVEN change C is being built on a branch, where T-01 edited its named
  file
- WHEN the owner runs `grove check` on that branch
- THEN C is not reported, since the file is unchanged on the target

### Requirement: Ask a round before writing

When the owner has confirmed the scope of a new change, the propose
instructions shall have the agent list, before writing any document, every
decision it would otherwise assume, numbered, each with a recommended
answer worded so that "yes" accepts it, and wait for the owner; each
answer is recorded as a decision in the proposal, quoted, and a decision
the owner declines to make is recorded as an open question naming the task
it blocks.

#### Scenario: The owner accepts the round

- GIVEN the scope is confirmed and the agent would otherwise assume where
  the check goes and whether the README is updated
- WHEN it lists both with recommended answers and the owner says "yes"
- THEN the proposal records each as a decision quoting the owner's "yes"
  and the recommended answer, and no document lists an assumption

#### Scenario: The owner declines one

- GIVEN the round lists where the check goes
- WHEN the owner says to leave that one open
- THEN the proposal lists it under open questions, naming the task it
  decides, and `/grove-apply` leaves that task blocked while
  `/grove-implement` asks

#### Scenario: Nothing to assume

- GIVEN the conversation settled every decision the documents need
- WHEN the agent reaches the round
- THEN it says there is nothing to ask and writes the documents

## Known discrepancies

None recorded.

## Open questions

None.
