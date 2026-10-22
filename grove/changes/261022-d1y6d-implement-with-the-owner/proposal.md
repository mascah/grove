# Implement a change with the owner

## Problem

Grove has one skill for building a change, apply, and it is written for a
session nobody watches. At 7ce7b91, `internal/instruct/requires/apply.md`
tells the agent not to settle an open question "even if you could ask",
and not to edit the change's files: a task that cannot be built as the
change says is marked `Status: blocked` and waits for a revision merged
through `/grove-propose`. That is right for `grove run`, which launches the
same skill unattended (`grove/specs/runs/spec.md`).

It is wrong for the owner building in their own session. In the 2026-10-22
explore session the owner said of a keyborg change: "I want to implement
this interactively since I will want to interact with it mid build and give
immediate feedback." With apply as it is, each piece of feedback that
changes behavior costs a blocked task, a proposal revision, a merge and
another apply.

## Intended outcome

`/grove-implement <change-id>` builds an approved change in the owner's
session. Where apply would block, implement asks the owner. When the
owner's answer or feedback changes what the change says, the agent, with
the owner's agreement, revises the change's proposal, change specs and
tasks on the implementation branch and then builds to the revision. Merging
that branch approves the revision along with the code (amended rule 4), and
`grove review` shows it. Implement records progress in `implementation.md`
and archives the change when every task is done, as apply does.

Apply stays the skill that asks no one. The rules are in
[specs/implement/spec.md](specs/implement/spec.md).

## Scope

### In scope

- The implement skill: `grove instructions implement <change-id>`, its
  requirements and the default schema's guide for it, the skill `grove
  init` writes for each harness, and `rules.skills.implement`.
- Documentation: getting started builds with `/grove-implement`; the
  reference lists the skill and the command; concepts says how implement
  and apply differ.
- Apply's wording, where it needs to say it is the skill that asks no one.

### Out of scope

- `implement <cycle-id>`, and apply taking a cycle ID: change 261022-1grb0,
  which waits on this one.
- Removing or narrowing apply's current forms (a change ID, or no ID for
  every approved change). They stay until 261022-1grb0, so `grove run` keeps
  working.
- Revising a change outside its implementation branch; that stays
  `/grove-propose`.

## Decisions

- **Implement and apply are separate skills.** The owner, 2026-10-22: "we
  should maybe restructure things and have a separate `implement` and
  `apply` skill". Reason: an interactive session and an unattended one need
  opposite rules about asking.
- **Implement may revise the change on its branch, with the owner's
  agreement.** The owner chose, 2026-10-22, "1" from: change the spec on the
  same branch; clarify only and revise through propose; build freely and
  fix the spec at archive. Reason: immediate feedback, while the merge
  stays the approval. The constitution allows it since #19 (rule 4).
- **Apply asks no one.** The owner chose, 2026-10-22, that `apply` is the
  runner's skill and `implement` is for everything the owner builds.

### Assumptions to confirm in review

- The skill is named `implement` (`/grove-implement`), as the owner wrote
  it.
- Implement refuses what apply refuses today: a change the target does not
  approve, and one that waits on a change the working tree has not archived
  as finished.
- Implement edits only the change's `proposal.md`, change specs and
  `tasks.md`, beside `implementation.md`; never `.grove.yaml`, which
  `depends_on` and `base` live in.
- After editing the change, implement runs `grove check` and fixes what it
  reports in the change before building.
- When the owner declines to answer, or says to leave it, the task is
  marked `Status: blocked` with the question, as apply does.
- A task renumbered by a revision keeps the rules `tasks.md` already has:
  IDs are never renumbered or reused.

## Open questions

None must be answered before review. None blocks a task.

## Related documents

- Requirements: [specs/implement/spec.md](specs/implement/spec.md)
- Tasks: [tasks.md](tasks.md)
- Constitution amendment: pull request #19
