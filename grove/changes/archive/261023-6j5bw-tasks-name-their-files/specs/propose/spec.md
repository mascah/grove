# Changes to propose

The changes that change 261023-6j5bw proposes to the requirements of
`propose`: the documents `/grove-propose` writes and what `grove check`
reports on them. There is no maintained spec for `propose` at b143991.

## ADDED Requirements

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

## Observed current behavior

- At b143991, `tasks.md` has no `Files:` line; the tasks template
  (`internal/instruct/schemas/default/documents/tasks/template.md`) shows
  `Result:`, `Constraints:` and `Acceptance:`, and the propose guide says
  "no file-by-file plan".
- `internal/doc/doc.go` reads a task's ID and title from its heading
  (`taskHeading`) and nothing else from `tasks.md`.
- `internal/cli/check.go`, `stale`, compares a change's MODIFIED and
  REMOVED requirements with the maintained spec at `base`;
  `changedSincePlan` diffs each open member's folder between `planned` and
  the target.
