# Implementation: Propose asks what matters

<!-- Progress on change 261023-yqfpj, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Propose orders its round, offers to split and writes delegation

Status: done

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
- Trials in Claude Code 2.1.310 on fresh ledger fixtures, headless, with
  this build first on PATH:
  - P1, `/grove-propose Editing an unknown ID prints an error`: the round
    asked the exit status (item 2) before where the check goes (item 4),
    with the test's name and place under routine defaults, but as a
    bulleted block; the round offered no delegation, and when the owner
    asked for it, T-01 got `Delegated: test layout; \`tests/test_records.py\``
    with `Files:` not naming the file. The bugfix schema it chose has no
    Decisions section, so the delegation was not recorded as a decision.
  - The routine-defaults wording became "on one line, not a list" in both
    `requires/propose.md` and the propose guide. P1b, re-run with the
    default schema: the exit status (1) before where the check goes (3),
    the routine defaults on one line (5), which itself recommended
    delegating test layout and private helpers; on "Yes." T-01 got that
    `Delegated:` line, and the proposal's Decisions recorded it quoting the
    owner's "Yes.". `grove check` clean.
  - P2, `/grove-propose Add ledger delete ID, and an error for an unknown
    ID`: while confirming scope it offered two changes (recommending one);
    on "split it" it ran `grove instructions propose` once for each, asked
    one round ordered behavior, interfaces, routine defaults, and wrote both
    on the branch, the delete change with `depends_on` naming the other.
    `grove check` clean.
- `dev/walkthrough.md` records P1 as "The round in order" and P2 as "Split
  a change" under Other paths.

### T-02 — Implement honors delegation and stops on a discovery, and ignores delegation unattended

Status: done

- `internal/instruct/requires/implement-owner.md`: make the choices and
  change the paths a task's `Delegated:` line names without asking; ask,
  as before, about a file neither `Files:` nor `Delegated:` names; and, on
  finding something that makes the agreed approach wrong, even where the
  task could still be built as written, stop that task, build nothing more
  of it, and tell the owner what was found and the revision the change
  needs, naming which documents and how.
- `internal/instruct/requires/implement-unattended.md`: a file no task's
  `Files:` line names stays as it is even if a `Delegated:` line names it;
  `Delegated:` lines are ignored, and a path the task would have changed is
  named, with why, in its evidence. The shared `requires/implement.md` is
  as change 261023-g0a3h left it.
- `docs/getting-started.md` step 5: delegation honored with the owner, the
  stop on a discovery, and `--unattended` ignoring `Delegated:`.
- `internal/cli/walkthrough_test.go`: `TestImplementInstructions` checks
  that `instructions implement` prints the delegation rule and the stop and
  not the rule to ignore `Delegated:`, and that `--unattended` prints that
  rule and neither of the other two.
- Trials: I2, `/grove-implement <change-id>` (no other words) on an
  approved change whose T-01 names `tests/test_records.py` only on its
  `Delegated:` line: it ran the owner's instructions and created
  `tests/test_records.py` for `find`'s tests without asking ("the task
  allowed me to add that file"); 13 tests pass. U1, the same change in a
  cycle of one, `grove run <cycle-id>`: `tests/test_records.py` was not
  created, `find`'s test went in `tests/test_cli.py`, and T-01's evidence
  says "`tests/test_records.py` is named only on the `Delegated:` line,
  which an unattended run ignores". (I1, with the prompt "/grove-implement
  <id>. I'm the owner; build it without stopping for my feedback", passed
  the whole sentence as the ID; the session then chose `--unattended`
  itself and ignored the line, as unattended should.)
- `dev/walkthrough.md` records both as "Delegated with me, ignored
  unattended".
- Recheck before finishing: T-01 and T-02's evidence holds against the code
  as it now is. `just check` passes (3732 Go lines, 1367 test lines);
  `grove check` prints "No problems in grove/.".

## Scenarios

- implement "A delegated path": trial I2, passed (`tests/test_records.py` created without asking); `TestImplementInstructions` checks the rule, passed.
- implement "A path nothing names": unverified, no trial staged a file outside `Files:` and `Delegated:`; the owner's part keeps the rule to ask, checked by `TestImplementInstructions`.
- implement "The approach rests on something false": unverified, no trial staged a false premise; `TestImplementInstructions` checks the stop rule is printed with the owner and not unattended.
- implement "A delegated path, unattended": trial U1, passed (file not created, `find`'s tests in `tests/test_cli.py`, the evidence names the file and why).
- propose "Two parts, one needing the other": trial P2, passed (offered while confirming scope, both written on the branch, the delete change naming the other in `depends_on`).
- propose "The owner keeps it whole": unverified, no trial declined a split; the rule says to write one change otherwise.
- propose "The owner delegates test layout": trial P1, passed for T-01's `Delegated:` line naming test layout and `tests/test_records.py` with `Files:` not naming it; the round in P1 did not itself recommend the delegation, the owner asked for it.
- propose "Nothing delegated": trial P1's first round, which offered none, and `TestPropose` for the rule, passed.
- propose "Check ignores the line": `walkthrough_test.go` `TestMovedFiles`, passed.
- propose "The owner accepts the round": unverified here beyond trials P1 and P2, whose answers were recorded quoting the owner; no assumptions section was written.
- propose "The round is ordered": trial P1, passed for the order (exit status second, before where the check goes fourth); the routine choices came as a bulleted block, which led to the "on one line, not a list" wording, rechecked in trial P1b, passed (one line).
- propose "One "yes" to the routine line": trials P1 and P2, passed (one "yes" accepted the routine defaults).
- propose "The owner declines one": unverified, no trial left a question open; unchanged from change 261023-g0a3h.
- propose "Nothing to assume": unverified, no trial had nothing to ask.

## Reconciliation

- implement: applied to specs/implement/spec.md
- propose: applied to specs/propose/spec.md
