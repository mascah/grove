# Changes to implement

The changes that change 261023-6j5bw proposes to the requirements of
`implement` (`grove/specs/implement/spec.md` at b143991).

## ADDED Requirements

### Requirement: Tell the owner of a moved file

When the owner runs `/grove-implement <change-id>` for a change in no
cycle that is not done, and a task names a file that differs on the target
from the same file at the change's `base`, the system shall say which task
and file, and the implement instructions shall have the agent ask the owner
before building that task, instead of building on the file as it was.

#### Scenario: A named file changed since base

- GIVEN the target approves change C, in no cycle, whose T-02 names
  `internal/cli/run.go`, which the target changed after C's `base`
- WHEN the owner runs `/grove-implement C`
- THEN Grove's instructions name T-02 and `internal/cli/run.go` as changed
  since base, and the agent, on reaching T-02, tells the owner what changed
  and asks before building

#### Scenario: A file no task names needs a change

- GIVEN the agent, building T-01, finds `docs/concepts.md`, which no task
  names, states what T-01 makes false
- WHEN it tells the owner and they agree
- THEN T-01's `Files:` line on the branch names `docs/concepts.md`, and the
  agent edits the page; had the owner declined, the page would be left as
  it is

#### Scenario: A member of a cycle

- GIVEN change C is a member of cycle K, and a file C's T-01 names changed
  on the target after K's `planned` commit
- WHEN the owner runs `/grove-implement K`
- THEN C is not reported against its `base`; `grove check` reports K as
  stale, and the owner re-plans K or builds on
