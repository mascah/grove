# Implementation: Maintained specs list the changes that shaped them

## Tasks

### T-01 — Archive adds a History line

Status: done

`skills/archive.md`, under "Reconcile the specs", says to add
`- [<change-id>](../../changes/archive/<folder>/) <title>` at the end of
`## History` of each maintained spec the change has a change spec for,
adding the section last if it has none, and that `archive` does not
compare History. The spec template ends with `## History`, and the spec
guide says what it holds. No Go code changed.

`internal/cli/walkthrough_test.go` TestArchive checks the archive
instructions give the line's form, that the printed template ends with
`## History`, and fills the new `record-editing` spec with a History line
for the change; `grove check` and `grove archive` pass. TestFold archives
folds whose maintained spec has no History. Removing History from the
template makes TestArchive fail. `just check` passed, 7,814 instruction
words.

### T-02 — Document History

Status: not started

### T-03 — History for Grove's own specs

Status: not started

## Scenarios

## Reconciliation
