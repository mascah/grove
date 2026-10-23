# A guide to writing specs of existing code

## Problem

Grove starts with no specs and gives no way to catch up. At 0c6f419,
`docs/getting-started.md` says "A new project has no specs, and that is
fine", and a maintained spec comes only from archiving a change that
alters behavior. An owner who wants the code they already have written
down as agreed behavior has no recipe for it.

The pieces exist. The change-spec guide says where there is no
maintained spec, put the proposed rules under ADDED and what the code
does under "Observed current behavior"
(`internal/instruct/schemas/default/documents/change-spec/guide.md`); the
spec guide says to keep the rule and record where the code differs under
"Known discrepancies"
(`internal/instruct/schemas/default/documents/spec/guide.md`). But the
tasks guide says nothing about a change with no work to do
(`internal/instruct/schemas/default/documents/tasks/guide.md`), so an
agent proposing such a change invents tasks. `grove check` and `grove
archive` already accept a change with no tasks (see "Observed current
behavior" in [specs/propose/spec.md](specs/propose/spec.md)).

In the 2026-10-23 review against OpenSpec, Spec Kit and AI-DLC, AI-DLC
has a reverse-engineering step for existing code; OpenSpec (issue #724)
and Spec Kit leave it out.

## Intended outcome

A guide, `docs/guides/spec-existing-code.md`, shows how to record
behavior the code already has: propose a change with the rules the owner
agrees under ADDED, what the code does today under "Observed current
behavior", and each difference the owner does not accept kept for the
maintained spec's "Known discrepancies"; give it no tasks; and archive
it on the same branch, so one merge records it. The guide is tried on
the ledger fixture before it is published. The tasks guide says such a
change has no tasks. The rule is in
[specs/propose/spec.md](specs/propose/spec.md).

## Scope

### In scope

- `docs/guides/spec-existing-code.md`, trialed on the ledger fixture,
  with the trial's result in the implementing pull request.
- One sentence in the tasks guide.
- Links from `docs/index.md`'s guides list, `docs/getting-started.md`
  and `mkdocs.yml`.

### Out of scope

- A `grove bootstrap` command or any new command, flag or field: parked
  (decision 2).
- Changes to `grove check` or `grove archive`: both already accept a
  change with no tasks.
- Archiving a draft on its branch with `/grove-implement`: change
  261023-6yafg, which this change waits on.
- A step in `dev/walkthrough.md`: guides are trialed recipes, not
  specification (rule 8).

## Decisions

The owner agreed to the direction on 2026-10-23 ("yes to all") and to
every answer below in one round ("yes to C & D").

- **A guide, `docs/guides/spec-existing-code.md`: propose a change that
  only records behavior the code already has, with the agreed rules
  under ADDED, what the code does today under "Observed current
  behavior", and differences the owner does not accept under the
  maintained spec's "Known discrepancies", with no tasks; it archives on
  the same branch, through 261023-6yafg's draft path or
  `/grove-archive`, so one merge records it. It is trialed on the ledger
  fixture before it is published, and the trial is a task with its
  result in the implementing pull request.** Reason: AI-DLC reverse
  engineers existing code; OpenSpec (#724) and Spec Kit leave it out;
  and guides are trialed recipes (rule 8).
- **The tasks guide says a change that only records behavior the code
  already has has no tasks; `grove bootstrap` stays parked, with no new
  command.** Reason: `grove check` and `grove archive` already accept
  zero tasks, so one sentence of instruction text is enough.

## Open questions

None must be answered before review. None blocks a task.

Once 261023-6yafg is built, `grove review` will say a merge approves
such a change "with its code", though it changes no code. Whether that
wording should differ belongs to 261023-6yafg and can wait.

## Related documents

- Requirements: [specs/propose/spec.md](specs/propose/spec.md)
- Tasks: [tasks.md](tasks.md)
