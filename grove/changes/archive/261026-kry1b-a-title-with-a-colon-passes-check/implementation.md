# Implementation: A title with a colon passes check

## Tasks

### T-01 — Quote the title in the change and cycle templates

Status: done

Evidence:
- `internal/instruct/metadata/change.grove.yaml` and
  `internal/instruct/metadata/cycle.yaml` print `title: "<title>"`.
- `docs/planning.md` and `docs/coordinating.md` show the example title
  quoted.
- `internal/cli/walkthrough_test.go` `TestPropose` fills the `.grove.yaml`
  template that `grove instructions propose` prints with
  `Records: delete by ID` and expects `grove check` to report no problem;
  with the templates reverted it fails (`grove check`: status 1, 1 problem
  found). `TestSchemas` expects `title: "<title>"` in the printed change
  template, and `TestCycles` in the printed cycle template.
- `just check` passed: 3902 Go lines, 1618 test lines, 7931 instruction
  words.

## Scenarios

The change has no change spec, so no scenarios (the owner agreed the fix
stays observed behavior, guarded by `TestPropose`).

## Reconciliation

No change spec to apply.
