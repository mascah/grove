# Implementation: The project index names requirements two changes both change

<!-- Progress on change 261118-40s1a, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — The index names a shared requirement, and plan tells the owner

Status: done

`doc.Share` and `doc.Change.Shares` (`json:"-"`, so `grove changes --json`
is unchanged); `doc.Load` fills it in `Collection.shares` from the ADDED,
MODIFIED and REMOVED blocks of each open change's capability files, by
capability and requirement name, leaving out a pair where one depends on
the other through `depends_on`, directly or through other changes, open or
archived. The index prints one line per share under the change. The plan
requirements gain one bullet (56 words; `just check` prints 8199
instruction words, 8143 when 261118-emse7 was archived).

Checks: `internal/doc/doc_test.go` `TestShares` (a MODIFIED/REMOVED pair,
a pair ordered through a third change, the same name in two capabilities),
passed. `internal/cli/walkthrough_test.go` `TestCoordinating`: with
recategorize and show-saved approved, `instructions plan "0.5 release"`
and `instructions propose` name each under the other; a cycle holding
both checks clean and `grove claim` claims it; with show-saved's
`depends_on: [recategorize]` committed, or recategorize archived, no line;
passed. By hand with `bin/grove` on a fresh fixture with two changes
modifying "Keep a category": the index shows the line under each, `grove
check` and `grove changes` print what they did before.

Check: the index line reads `also changed by <change-id>: <capability> "<requirement>"`, the tasks' suggested words, since the proposal's open question on the wording had no answer; the requirement name is printed with Go's `%q`, as `grove check` prints requirement names.
Check: 261118-emse7, built and archived on this branch (b26f250), changed `internal/doc/doc.go`, `internal/instruct/instruct.go`, `internal/instruct/requires/plan.md` and `internal/cli/walkthrough_test.go` since this change's base f7e0140; this builds on them as they are now. The new plan bullet sits beside emse7's bullet on notes; it does not ask plan to write the pair in the cycle's notes, which the change does not say.
Check: in the walkthrough test every change `propose` writes holds an ADDED "Edit older records" block, so delete and budgets now show an `also changed by` line for it too; that is the behavior the requirement asks for, and the test asserts only the "Keep a category" lines.

## Scenarios

<!-- Before finishing, one line per scenario of the change specs: the check
that shows it and its result, or `unverified` and why:

    - record-editing "Update the note on a legacy record": tests/test_cli.py test_edit_legacy, passed
-->

## Reconciliation

<!-- One line per change spec, once applied to its maintained spec:

    - record-editing: applied to specs/record-editing/spec.md
-->
