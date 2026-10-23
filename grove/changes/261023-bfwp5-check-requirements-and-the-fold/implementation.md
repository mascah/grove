# Implementation: Check requirements and the archive fold

<!-- Progress on change 261023-bfwp5, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Check each requirement's shape

Status: done

- `internal/doc/requirements.go`: `Requirement.Problems` reports a missing
  "shall" (a whole word, any case) or `#### Scenario:`, nothing for REMOVED;
  `ChangeHeadings` gives the lines of ADDED, MODIFIED and REMOVED headings
  outside fenced code, sharing the fence test with `Requirements`.
- `internal/cli/check.go`: `shapes` reports both for each maintained spec,
  and the requirement shape for each open change's change specs; archived
  changes are not read.
- `internal/doc/requirements_test.go`: `TestShape` covers a fenced heading,
  "SHALL", a REMOVED requirement and both problems.
- `internal/cli/walkthrough_test.go` (named by T-02): `TestCheck` on the
  fixture reports a maintained requirement with no scenario at line 5, one
  with no "shall" at line 9 and a pasted `## ADDED Requirements` at line 3,
  and not a REMOVED requirement of an open change nor an archived change's
  requirement with no scenario. The fixture requirements the other tests
  write now say "shall" and have a scenario, so their `check` stays clean.
- `grove check` on Grove's own repository: "No problems in grove/."
- `go test ./...` passes.

### T-02 — Archive refuses a fold that does not match

Status: not started

### T-03 — Document the checks

Status: not started

## Reconciliation

<!-- One line per change spec, once applied to its maintained spec:

    - record-editing: applied to specs/record-editing/spec.md
-->
