# Implementation: Cite what the project's readers can open

<!-- Progress on change 261028-rxsfa, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Every skill is told to cite what readers can open

Status: done

- `internal/instruct/instruct.go`: the `meanings` constant gains, after the
  paragraph on drafts, the sentence agreed in the proposal: "Write Grove
  documents for whoever reads this project: cite only what they can open,
  and quote or summarise what only this machine holds, such as an unpushed
  branch, a file outside the checkout or a private project."
- `bin/grove instructions explore`, `bin/grove instructions propose` and
  `bin/grove instructions implement 261028-rxsfa` each print it at line 17,
  before `## What Grove requires of <skill>` at line 24.
- `just check` passed: 4024 Go lines, 1681 test lines, 7997 instruction words.

## Scenarios

<!-- Before finishing, one line per scenario of the change specs: the check
that shows it and its result, or `unverified` and why:

    - record-editing "Update the note on a legacy record": tests/test_cli.py test_edit_legacy, passed
-->

## Reconciliation

<!-- One line per change spec, once applied to its maintained spec:

    - record-editing: applied to specs/record-editing/spec.md
-->
