# Cycles planned for unattended runs

## Problem

Merging a change approves it, and at 7ce7b91 approval is all `grove run`
looks at: it selects every approved change that is not archived and waits
on nothing (`internal/cli/run.go`, `buildable(ip.Agreed, ip.Agreed)`), and
the session it launches, `/grove-apply` with no ID, selects the same set
again (`internal/instruct/requires/apply-all.md`). The owner runs `grove
run` from launchd on a Mac mini for the keyborg project. In the 2026-10-22
explore session they said: "once I merge this, the launchd item on the mini
is going to automatically try and build this, when really I want to
implement this interactively". Today the only way to keep the mini off a
change is to stop the launchd job.

Nothing checks, before an unattended build, whether a change needs a
person: credentials the runner lacks, manual steps, judging how a screen
looks.

Cycles were built once and parked (D-68, commit `2f9b457`; the build is on
branch `phase-3-cycles`). There a cycle was a planning group: a file with a
title, goal and dates, and a `cycle` field on each member change. It had no
check for human work and no link to `grove run`, which did not exist yet.

## Intended outcome

A cycle is approved changes the owner means to build together, in a file
`grove/cycles/<id>-<slug>.yaml` listing them. `/grove-plan` writes it: from
a description, from one change ID (a cycle of one), or revising a cycle by
its ID. Plan reads each change, asks the owner about any that may need a
person, records the ones the owner says do, and records the target commit
it checked. Merging the cycle file plans the cycle.

Handing a cycle to the runner is a separate, later merge: `unattended:
true` in the cycle file. `grove check` reports a cycle so marked that still
has a change the owner said needs a person, and a cycle whose changes
changed on the target since its plan.

`grove run` builds only cycles: the oldest one marked `unattended` that is
not done, not stale, has no change needing a person and waits on nothing
outside it. It launches `/grove-apply <cycle-id>` and pauses until that
cycle is archived. An approved change in no cycle, or in a cycle not
marked, is never built unattended; the owner builds it with
`/grove-implement`, which now also takes a cycle ID. Apply takes only a
cycle ID.

The rules are in [specs/cycles/spec.md](specs/cycles/spec.md),
[specs/runs/spec.md](specs/runs/spec.md),
[specs/apply/spec.md](specs/apply/spec.md) and
[specs/implement/spec.md](specs/implement/spec.md).

## Scope

### In scope

- The cycle file, its template and how `grove check` reads it.
- The plan skill: `grove instructions plan`, its requirements and the
  default schema's guide, the skill `grove init` writes, and
  `rules.skills.plan`.
- The `unattended` mark, the needs-a-person check and the stale-plan check
  for cycles.
- `grove run` selecting one cycle per attempt, launching apply with its ID,
  and pausing per cycle.
- `apply <cycle-id>` replacing apply's change-ID and no-ID forms;
  `implement <cycle-id>`.
- `grove changes` and `grove review` showing cycles; the project index in
  instructions listing them.
- Documentation: a docs page for planning a cycle and marking it; running,
  coordinating, reference, index, README and the guides that type
  `/grove-apply` updated; the walkthrough.

### Out of scope

- Per-change unattended marks; one change goes to the runner as a cycle of
  one.
- Building more than one cycle per attempt, and a priority field.
- Start and end dates on a cycle.
- A command that sets `unattended`; the owner edits the file and merges it.

## Decisions

- **The runner builds only cycles.** The owner, 2026-10-22, on what the
  runner reads: "I lean towards 1 more than 2 slightly", 1 being cycles
  only, and confirmed the scope stated before this proposal. "Sometimes I
  just want to send an individual change to be implemented by the runner":
  that is a one-change cycle from `/grove-plan <change-id>`. Reason: one
  kind of mark in one kind of file.
- **Plan, then mark separately.** The owner: "2, plan then mark ready
  separately." and "A cycle should only be elligible for execution after
  it's explicitly been marked as ready." Reason: a planned cycle the owner
  never marks is one they build themselves.
- **Plan asks the owner about human work.** The owner: "I think the plan
  flow should ask the user?", and confirmed that only the changes they say
  need a person are noted, and those block the mark.
- **One mark, named `unattended`.** The owner, on splitting "ready" from
  "can run unattended": "1, call it unattended or auto", then "unattended".
- **A stale plan is caught.** The owner chose "1": plan records the target
  commit it checked; a cycle whose changes differ on the target since then
  is skipped by `grove run` and reported by `grove check`.
- **One cycle per attempt, oldest first; the pause is per cycle.** The
  owner chose "1".
- **Apply takes a cycle ID; implement takes a change or cycle ID; no-ID
  apply goes.** The owner chose "1".
- **Membership is listed in the cycle file; the file has a title and an
  optional goal.** The owner's answers, 2026-10-22. Reason: plan never edits
  an approved change's `.grove.yaml`.
- **`phase-3-cycles` is inspiration, not a template.** The owner: "We can
  use the old branch as inspiration we don't have to copy it exactly".
- **The constitution allows it** since #19 (rules 2, 4, 7, 10 and 11).

### Assumptions to confirm in review

- The cycle file's fields are `title`, optional `goal`, `changes` (approved
  change IDs), `planned` (the target commit the plan checked),
  `needs_human` (change ID to the reason the owner gave; optional) and
  `unattended` (optional).
- Plan never writes or removes `unattended`. Re-planning a marked cycle
  keeps the mark; if the owner then says a change needs a person, `grove
  check` reports the mark and plan tells the owner to remove it.
- Plan includes only changes the target approves and that are not
  archived, and tells the owner of a member that waits on a change outside
  the cycle that is not archived as finished.
- "Stale": the folder of a member change not yet archived differs on the
  target from the same folder at `planned`.
- A change may be in only one cycle that is not done; `grove check` reports
  a second. A cycle is done when every member is archived; its file stays.
- Build order inside a cycle is `depends_on`, otherwise oldest first, as
  today.
- Apply builds the cycle it is given whether or not it is marked; `grove
  run` alone enforces the mark.
- When no cycle is eligible, `grove run` says there is nothing to apply, as
  today; `grove check` says why a marked cycle is skipped.
- `grove changes` shows each open change's cycle and whether it is marked
  `unattended`; `grove review` labels a cycle file `cycle`.
- Projects that use `/grove-apply` with a change ID or no ID, or a pipeline
  typing `/grove-apply`, must switch to `/grove-implement` or plan a cycle.
  This is a breaking change.

## Open questions

None must be answered before review. None blocks a task.

## Related documents

- Requirements: [specs/cycles/spec.md](specs/cycles/spec.md),
  [specs/runs/spec.md](specs/runs/spec.md),
  [specs/apply/spec.md](specs/apply/spec.md),
  [specs/implement/spec.md](specs/implement/spec.md)
- Tasks: [tasks.md](tasks.md)
- Waits on: 261022-d1y6d (Implement a change with the owner)
- Constitution amendment: pull request #19
