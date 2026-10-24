# Implementation: Propose asks what matters

<!-- Progress on change 261023-yqfpj, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Propose orders its round, offers to split and writes delegation

Status: in progress

- `internal/instruct/requires/propose.md`: the round is ordered by
  consequence (behavior, scope, compatibility or data loss first; then
  architecture or interfaces costly to reverse; routine choices last, on
  one line of recommended defaults one "yes" accepts), leaving none out;
  while confirming scope, an idea too big for one change, or with parts
  the owner could agree separately, is offered as a split, written on this
  branch as changes from their own `instructions propose`, linked with
  `depends_on`; a `Delegated:` line is written only where the owner agreed
  in the round.
- `internal/instruct/schemas/default/skills/propose.md`: step 2 offers the
  split; step 3 gives an example of each tier of the round, and says
  delegation may be recommended.
- `internal/instruct/schemas/default/documents/tasks/guide.md` and
  `template.md`: the optional `Delegated:` line, right after `Files:`,
  honored with the owner, ignored unattended, not read by Grove; the
  template's example shows one.
- `docs/getting-started.md` step 3, `docs/concepts.md`'s `tasks.md` row and
  `docs/reference.md`'s description of a task.
- `internal/cli/walkthrough_test.go`: `TestPropose` checks the round's
  order, the split offer and the `Delegated:` line in the tasks guide and
  template; `TestMovedFiles` gives T-01 a `Delegated:` line right after its
  `Files:` line, and `grove check` reports no problem on it, nor on
  `tests/test_records.py`, which only that line names, when main adds it
  after base.
