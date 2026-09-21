---
id: "G-007"
type: work
title: Create records with shared sequential ID allocation
status: done
kind: feature
priority: 2
size: medium
members: []
depends_on: []
relates_to: ["G-003", "G-004", "G-006", "G-002"]
created: "2026-09-19T15:25:51Z"
updated: "2026-09-19T15:52:00Z"
formerly: "W-002"
---

## Outcome

Run `grove new <type> <title>` in a Git checkout to allocate the next
`W-`, `Q-`, or `D-` number, write a valid record file without overwriting
anything, and print its path. Two Grove commands in different linked worktrees
of one repository never receive the same ID. Grove's own records are the first
data; this work should be the first record created by the command's successor
test, not by hand.

## Why now

G-003 is done and every record so far was numbered by hand. The next planned
work, combined views and execution, depends on records that agents can create
without coordinating numbers in conversation. Priority 2 follows the brief's
sequence; medium reflects locking, durable state, recovery, and file creation
with tests, without any update or move behavior.

## Constraints

- Meet the [allocation requirements](../docs/record-model.md#identity-and-dates):
  exclusive lock across worktrees and types, durable reservation before release,
  no ID without persistence, non-overwriting creation, gaps accepted.
- Coordination state lives in a Grove-owned directory under
  `git rev-parse --path-format=absolute --git-common-dir`. It is never a
  tracked record and the read-only commands never create it.
- Require Git for creation; plain directories still read. Do not add automatic
  allocation outside Git or a `--id` bypass in this work.
- Never silently restart at 1. Recovery must consider committed records on every
  local ref and live records in every worktree, and must state that it cannot
  recover reservations for records that were never written or were deleted.
- Generate filenames and slugs per the [folders and files](../docs/record-model.md#folders-and-files)
  contract. Write `created` and `updated` with one UTC timestamp.
- Keep G-003's reader unchanged in behavior; reuse it to validate the project
  after creation. Keep the installed sibling CLI untouched.
- Defer status/field updates, renames, moves, deletes, imports, clone merging,
  Windows locking, and cross-branch views.

## Acceptance

- `new work "Title"` in this repository creates `grove/work/W-NNN-<slug>.md`
  with the next number, matching frontmatter, and a body skeleton; `check`
  passes afterward and `list` shows the record.
- `new` from a linked worktree and from main share one counter; a fixture with
  concurrent invocations across two worktrees yields distinct IDs.
- With no counter state, the first allocation initializes from the highest ID
  found across local refs and worktrees, reports that on stderr, and never
  returns a number at or below one already in use.
- A counter file lower than an observed ID is corrected upward and reported.
- Lock or persistence failure creates no file and exits 1 with a diagnostic.
- An existing target filename is never overwritten; `--slug` selects the name.
- A process killed while holding the lock does not block the next command.
- Tests run in temporary Git repositories; nothing writes to this checkout's
  records or Git directory during tests.

## Design

[G-006](G-006-allocator-mechanism.md) proposes the mechanism:
`flock` on `<common-dir>/grove/lock`, a plain-text `<common-dir>/grove/next-ids`
file replaced atomically, and a scan of local refs and worktrees as a floor on
every allocation. Read it before planning. Command shape:

```text
grove [--project DIR] new work|question|decision TITLE [--slug SLUG]
```

Allocation steps: lock; read counter; scan floor; next = max(counter, floor+1);
write counter via temp file, fsync, rename; unlock; then create the file with
`O_CREATE|O_EXCL`; then reload the project and fail if validation fails without
deleting the file. Git commands used: `rev-parse --git-common-dir`,
`worktree list --porcelain`, `for-each-ref`, `grep`.

## Evidence

Closed 2026-09-19 on branch `worktree-W-002` at `8d359e3`. The
[implementation plan](G-008-create-plan.md) records the tests,
race/vet/format results, the real creation of G-009 through the command, and
the independent review. Every acceptance line above has a fixture test; the
review's blocking finding, an unreadable file in another worktree silently
lowering the floor, was fixed and covered before closure. Structural
validation of the created record is proven; usefulness of the body skeleton
is the owner's judgment.

## Next

Integrate this branch into main, then shape
[G-009](G-009-update-records.md) for status and field updates.
