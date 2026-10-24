# Revising a change re-plans its cycle

## Problem

Getting a paused `grove run` past a blocked change takes two merges. At
57af2e3, `docs/running.md` ("The pause") says to answer the question with
`/grove-propose <change-id>`, merge the revision, re-plan the cycle with
`/grove-plan <cycle-id>` and merge that, "since the revision made its plan
stale", then run `grove run --retry`. The plan's finish says the same for a
member the owner said needs a person (`internal/instruct/requires/plan.md`).

The second merge exists because of how the two steps are built.
`instruct.Revise` prints only the change being revised
(`internal/instruct/instruct.go`), so the session that revises a member
leaves its cycle alone. And `changedSincePlan` (`internal/cli/check.go`),
which `grove check` and `grove run` both use, diffs each open member's
folder and named files between the cycle's `planned` commit and the
target, so the revision itself makes the cycle stale, even when the same
merge set `planned` again.

The owner wants cycles kept, as the unit of delivery and of permission to
build unattended, and this loop shortened by one merge.

## Intended outcome

`/grove-propose <change-id>`, for a member of a cycle that is not done,
also prints the cycle's file. On the same branch, the agent asks the owner
whether the change needs a person, as `/grove-plan` does, and sets the
cycle's `planned` to the target commit the instructions name. It never
writes or removes `unattended`.

The stale check no longer counts a member's change that arrived on the
target in the same commit, on its first-parent history, as an edit to the
cycle file. So a squash merge or a merge commit of the revision and the
re-plan leaves the cycle as its plan checked it, and `grove run --retry`
builds it. A member changed by any other commit still makes the cycle
stale.

Past a pause, the steps become: revise the change with `/grove-propose
<change-id>`, merge, `grove run --retry`. The rules are in
[specs/propose/spec.md](specs/propose/spec.md),
[specs/cycles/spec.md](specs/cycles/spec.md) and
[specs/runs/spec.md](specs/runs/spec.md).

## Scope

### In scope

- The propose instructions for revising a member of a cycle not done:
  the cycle's file, the needs-a-person question for that change, and
  `planned`.
- The plan instructions' finish, which names the shorter way past a
  member that needs a person.
- The stale check shared by `grove check` and `grove run`.
- `docs/running.md`, `docs/planning.md`, `docs/coordinating.md`,
  `docs/reference.md` and `dev/walkthrough.md`.

### Out of scope

- Dropping cycles or the stale check: a cycle stays the unit of delivery
  and of permission to build unattended (`CONSTITUTION.md` rules 4 and
  10).
- `/grove-plan <cycle-id>`: it stays, for a re-plan with no revision, such
  as after a member's named file moved.
- The other members of the cycle: propose asks about the revised change
  only.
- Grove committing or merging anything; the owner merges as before.
- A revision and a re-plan that reach the target as separate commits, by
  rebase or fast-forward: the cycle stays stale, as decision 2 says.
- Revising a change in no cycle: unchanged.
- How `grove run` launches its session: change 261023-g0a3h, which this
  change depends on, has it type `/grove-implement <cycle-id>
  --unattended`. Runs' "Retry" is the one requirement both modify; this
  change starts from that change's text.

## Decisions

The owner agreed to the direction on 2026-10-23 ("yes to all") and to
every answer below in one round ("yes").

- **`/grove-propose <change-id>`, for a change that is a member of a cycle
  that is not done, also prints that cycle's file, and the agent re-plans
  it on the same branch:** it asks the owner the needs-a-person question
  for that change, as `/grove-plan` does, and sets `planned` to the target
  commit the instructions name. It never writes or removes `unattended`.
  Reason: the revision and the re-plan then go in one branch and one
  merge, and the owner still answers the question the plan exists to ask,
  so `grove run` still builds only what a plan checked (rule 10). The mark
  stays the owner's, in a merge of its own (rule 4).
- **The stale-plan check ignores a change to a member, its folder or a
  file its tasks name, when that change arrived in the same commit on the
  target's first-parent history as an edit to the cycle file;** squash
  merges and merge commits both count. A member changed by any other
  commit still makes the cycle stale, such as a member changed on main and
  the cycle marked `unattended` in a later merge. Reason: a commit that
  edits the cycle file with its member is the plan checking that member;
  reading first parents treats a merge commit as one commit, as a squash
  is. It is about 20 lines of Go in `internal/cli/check.go`, which
  `internal/cli/run.go` already shares.
- **The documented steps past a pause drop the separate re-plan merge,**
  in `docs/running.md` ("The pause"), `docs/planning.md` ("Re-plan a
  cycle") and `docs/coordinating.md` where it names re-planning. Reason:
  the documentation is the specification (rule 8).

## Open questions

None must be answered before review. None blocks a task.

## Related documents

- Requirements: [specs/propose/spec.md](specs/propose/spec.md),
  [specs/cycles/spec.md](specs/cycles/spec.md),
  [specs/runs/spec.md](specs/runs/spec.md)
- Tasks: [tasks.md](tasks.md)
