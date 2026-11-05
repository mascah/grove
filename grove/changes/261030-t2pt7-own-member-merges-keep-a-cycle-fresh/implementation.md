# Implementation: A cycle's own member merges keep its plan fresh

<!-- Progress on change 261030-t2pt7, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Count a cycle's own archived members as expected

Status: done

- `internal/cli/check.go`: `sincePlan` keeps its signature and now adds
  the archive folders of the changes `cy` lists, read from `cy` as the
  caller passes it, to the folders whose archiving commit is skipped,
  beside those of the cycles it depends on. Its comment says so.
- `internal/cli/walkthrough_test.go`: `TestFollows` gives L's member d a
  task naming `ledger/cli.py`. After c alone lands in a squash editing that
  file, `grove check` is clean, `grove run` selects L and
  `grove instructions implement d` names no moved file; a later commit
  archiving nothing makes L stale (and K) and the implement line names the
  file. With `check.go` as main has it, the test fails at the first of
  these ("check printed ... 261023-p0000 ... changed on main since this
  plan").
- `docs/planning.md` "Re-plan a cycle": a commit archiving a change of
  this cycle does not count, so its changes can reach main one pull
  request at a time.
- `dev/walkthrough.md` step 8 of "Coordinating changes" ("Building a cycle
  unattended") opens with the paragraph decision 6 describes. Its
  `grove check` is clean as written: the squash archives delete, a member
  of the 0.4 cycle, and `check` does not compare a cycle member against
  its base. The step was not run by hand in a harness here.
- `just check`: passed (4030 Go lines, 1692 test lines, 8023 words).

## Scenarios

- cycles "Revised and re-planned in one squash merge": internal/cli/walkthrough_test.go TestReplan ("a squash merge"), passed
- cycles "Revised and re-planned in one merge commit": internal/cli/walkthrough_test.go TestReplan ("a merge commit"), passed
- cycles "Marked after a member changed": internal/cli/walkthrough_test.go TestReplan ("marked after a member changed"), passed
- cycles "Revision and re-plan as two commits on the target": internal/cli/walkthrough_test.go TestReplan ("two commits fast-forwarded"), passed
- cycles "A member lands in its own squash merge": internal/cli/walkthrough_test.go TestFollows ("one of L's changes archived, editing a file its other change names"), passed
- cycles "Then a commit archiving nothing": internal/cli/walkthrough_test.go TestFollows ("a commit archiving nothing after one of L's changes"), passed
- cycles "The cycle it follows lands in a squash merge": internal/cli/walkthrough_test.go TestFollows ("L's last change archived in a squash"), passed
- cycles "Stacked pull requests": internal/cli/walkthrough_test.go TestFollows (L's two changes landed as two commits, then K selected), passed
- cycles "Through another cycle": internal/cli/walkthrough_test.go TestFollows ("M's change archived in a merge commit", K clean), passed
- cycles "Any other commit still counts": internal/cli/walkthrough_test.go TestFollows ("a commit archiving nothing", "an archive of a cycle K does not follow"), passed

## Reconciliation

<!-- One line per change spec, once applied to its maintained spec:

    - record-editing: applied to specs/record-editing/spec.md
-->
