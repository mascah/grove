# Implement reconciles a requirement rewritten on the target

## Problem

Two approved changes can each rewrite the same requirement. Once the
first is archived, the second still holds the text it was written
against, and its fold puts that older text over the newer agreed one.

- `grove check` reports the second change: "requirement "Keep a category"
  was rewritten since base 4f1c2a9; revise the change with propose"
  (`grove/specs/propose/spec.md` at 07d319f, "Report a stale
  requirement").
- The implement instructions name a task's file that changed on the
  target since the change's `base`, and say what to do about it
  (`grove/specs/implement/spec.md` at 07d319f, "A task's file changed on
  the target"). They name no requirement that was rewritten, and nothing
  has the agent run `grove check` before it builds
  (`internal/instruct/requires/implement.md` at 07d319f).
- The archive instructions run `grove check` only after the fold
  (`internal/instruct/schemas/default/skills/archive.md` at 07d319f). By
  then the maintained spec holds the change's own text, which the check
  allows, so it reports nothing (`internal/cli/check.go` at 07d319f, the
  comment on `stale`).
- `grove archive` accepts the fold, because the change names that
  requirement (`CONSTITUTION.md` rule 3 at 07d319f).

So a change built after another rewrote its requirement drops the other's
text from the maintained spec with no conflict in Git, no failing test
and no report. An example: a requirement says editing a record keeps its
category. Change A rewrites it to keep the category and the tags. Change
B, written before A landed, rewrites it to keep the category unless the
user clears it. Built after A, B's fold removes the tags, which nobody
decided.

Today the owner meets this when planning a cycle, where plan brings a
stale member up to date (`grove/specs/cycles/spec.md` at 07d319f, "Plan a
cycle"). Change 261122-7vht2 lets a change be claimed without a plan, and
then nothing stands between the stale change and its fold.

## Intended outcome

When a change modifies or removes a requirement that was rewritten or
removed on the target since the change's `base`, the instructions
`/grove-implement` prints name that requirement, as they name a changed
file today. Before building any task, the agent revises the change on its
branch so the requirement holds both what the maintained spec now says
and what the change alters, and records what it reconciled for the
reviewer. The merge of the branch approves the revision with the code, as
`CONSTITUTION.md` rule 4 already says.

## Scope

### In scope

- The implement instructions naming each such requirement, for an
  approved change, a draft and a member of a cycle, with or without
  `--runner`.
- What the instructions have the agent do about it.
- A new requirement in the implement spec.
- `docs/coordinating.md`, `dev/walkthrough.md` and the walkthrough test.

### Out of scope

- `grove claim`: it refuses nothing for staleness against `base`. Whether
  a change in no cycle can be claimed at all is change 261122-7vht2.
- `grove check`'s report and the propose instructions for a stale
  requirement: unchanged. The agent revises by following
  `grove instructions propose <change-id>`, which already says how.
- `grove archive`: it compares the fold as it does now.
- A cycle's stale plan, which still stops `grove claim`
  (`grove/specs/cycles/spec.md`, "Count a member changed with its plan").
- Any amendment to `CONSTITUTION.md`: rule 4 already has the merge
  approve a change revised on its branch by implement.

## Decisions

The owner settled these in an explore on 2026-11-21 and answered the
round before writing with "yes" to each recommended answer.

1. The second change's text is reconciled during the build, not before
   it. Asked whether `grove claim` should refuse a change whose
   requirement was rewritten until a revision is merged, or take it and
   have the builder reconcile it on the branch, the owner said: "during
   the build". Reason: a refusal would need a merge for each such pair,
   which is scheduling by merge again; the revision is still approved by
   a merge, the one that brings the code.
2. This is its own change, and 261122-7vht2 names it in `depends_on`.
   Reason: it stands alone, since a change built with the owner has the
   same gap today, and it is what makes claiming a stale change safe.
3. The implement instructions name each rewritten or removed requirement
   themselves, from the check `grove check` already runs; the agent is
   not told to run `grove check` to find out. Reason: it does not depend
   on the agent remembering a command, and it reads like the line that
   names a changed file.
4. Routine, accepted together: private helpers, test names and layout,
   and the exact wording of the line the instructions print are the
   implementer's to choose.

The owner's reasoning drew on how the projects that influenced this
handle work that went stale: beads' `bd ready` lists work with no open
blocker and never asks whether it is still current, and Gas Town's
Refinery rebases each finished branch on the target and runs the tests
before merging (https://github.com/gastownhall/beads,
https://github.com/gastownhall/gastown). Both check after the work, not
before it. A rewritten requirement is the one case that kind of check
cannot see, since Git and the tests both pass.

## Open questions

None.

## Related documents

- Requirements: [implement](specs/implement/spec.md)
- Tasks: [tasks.md](tasks.md)
