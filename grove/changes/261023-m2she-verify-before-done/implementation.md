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

Status: not started

## Reconciliation

<!-- One line per change spec, once applied to its maintained spec:

    - record-editing: applied to specs/record-editing/spec.md
-->
