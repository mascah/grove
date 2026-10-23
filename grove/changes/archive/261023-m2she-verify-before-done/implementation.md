# Implementation: Verify before done

<!-- Progress on change 261023-m2she, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Propose checks the change before presenting it

Status: done

- `internal/instruct/requires/propose.md`: before telling the owner the
  change is written, check observable outcomes, that a task's acceptance
  covers each ADDED or MODIFIED requirement, that scope, specs and tasks do
  not contradict, and that no task is out of scope; fix what it can, say
  what it changed, ask about what needs a decision.
- `internal/instruct/schemas/default/skills/propose.md`: "5. Check the
  change", how to read the documents against each other.
- `docs/getting-started.md` step 3 and the `instructions propose` row of
  `docs/reference.md` describe the check; `dev/walkthrough.md` main line
  step 3 and a new "Propose checks the change" under Other paths.
- `internal/cli/walkthrough_test.go`: `TestPropose` checks that
  `instructions propose` prints the rule and the guide's step.
- Trials in Claude Code 2.1.309 on fresh ledger fixtures, with this build's
  `grove` first on PATH, run headless with the owner's answers given up front:
  - Revising a hand-written draft whose "Print the deleted record" no
    acceptance covered, told to change only T-01's title: the session ran
    the check and named the uncovered requirement and the fix it would make,
    but did not make it, since the owner had asked for only the title
    (the revise rule, "change what the owner asks for and leave the rest").
  - A new change from "Add ledger delete ID, which removes the record and
    prints the deleted record": the session ran the check, found T-02's
    acceptance contradicting the scope, fixed it, and said what it changed
    under "Fix made while checking".
  - Not shown: a requirement no task covers, fixed before presenting, as
    the acceptance asks. The new change's draft had no such gap, and in the
    revision the revise rule won.
- `just check` passes.

### T-02 — Implement and apply verify before finishing

Status: done

- `internal/instruct/requires/implement.md` and `requires/apply.md`: before
  finishing, run the project's tests, re-read each done task's evidence
  against the code as it now is, set one that no longer holds back to
  `Status: in progress`, and write `## Scenarios`; the finishing report names
  each unverified scenario; apply asks no one. A task marked done is not
  rebuilt unless that check finds its evidence no longer holds.
- `internal/instruct/requires/archive.md`: for a finished change, name each
  `unverified` scenario to the owner first; it does not stop the archive.
- `internal/instruct/schemas/default/skills/implement.md` and `skills/apply.md`:
  what makes a good check for a scenario, and when it is `unverified`. The
  guides' "Run the project's tests before you finish" moved into the
  requires text, so it is said once.
- `internal/instruct/schemas/default/documents/implementation/template.md`
  and `guide.md`: a `## Scenarios` section with an example line; the
  `Status:` lines are unchanged.
- `docs/concepts.md`'s change table and `docs/getting-started.md` step 5
  describe the scenario list and the recheck; `dev/walkthrough.md` main line
  step 5.
- `internal/cli/walkthrough_test.go`: `TestApplyInstructions` checks the rule
  in `instructions apply` (with "Ask no one." and the template's
  `## Scenarios`) and `instructions implement`; `TestArchive` checks it in
  `instructions archive`.
- Trial in Claude Code 2.1.309 on a fresh ledger fixture, headless, with an
  approved change "Unknown ID error" with two scenarios, one a test can show
  and one, "The message reads clearly", only a person can judge, and the
  owner's "build without waiting, then archive" given up front:
  `/grove-implement` built T-01 with `test_edit_unknown_id` and wrote
  `## Scenarios` with `tests/test_cli.py test_edit_unknown_id, passed` and
  `unverified, only a person can judge how the message reads`. Its report
  named the unverified scenario, and the archive created
  `specs/record-editing/spec.md` with the requirement copied word for word
  and archived the change. From a headless transcript that only has the
  final report, it cannot be shown that the scenario was named before the
  fold rather than after.
- Recheck before finishing: T-01's evidence still holds (its rule is in
  `requires/propose.md` and `TestPropose` passes). `just check` passes;
  `grove check` is clean.

## Scenarios

- propose "A requirement no task covers": trial, unverified. The check ran
  and named the gap, but in a revision limited to one edit it was not fixed.
- propose "Scope and a task disagree": unverified, no trial staged it. In
  the new-change trial a contradiction between scope and acceptance was
  fixed without asking, since the round's answers already settled it.
- implement "Every scenario has a test": trial, passed (one testable
  scenario named with its test and "passed").
- implement "A done task later broken": unverified, no trial staged it.
- implement "A scenario only a person can judge": trial, passed.
- implement "One scenario unverified" (archive): trial, passed in part. It
  was named in the report, but whether that came before the fold is not shown.
- apply "A scenario no test shows": unverified, no apply trial was run.
  The rule's text is checked by `TestApplyInstructions`.

## Reconciliation

<!-- One line per change spec, once applied to its maintained spec:

    - record-editing: applied to specs/record-editing/spec.md
-->

- apply: applied to specs/apply/spec.md
- implement: applied to specs/implement/spec.md
- propose: applied to specs/propose/spec.md
