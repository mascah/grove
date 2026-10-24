# Implementation: Specs without .grove.yaml, and no design.md

<!-- Progress on change 261023-d6r79, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Read a spec from its spec.md

Status: done

- `internal/doc/doc.go`: `Load` finds specs by `specs/*/spec.md`, takes
  the title and description from it (`specMeta` in `meta.go`), reports a
  spec with no `# ` line, and puts each `specs/*/.grove.yaml` in the new
  `Collection.Diagnostics` as no longer read, without parsing it.
- `internal/cli/check.go` prints `Collection.Diagnostics`, so a
  `.grove.yaml` in a folder with no `spec.md` is reported too. Added to
  T-01's `Files:` with the owner's agreement.
- `internal/instruct`: `metadata/spec.grove.yaml` deleted; the default
  schema's `spec` kind has no `metadata`; `problems` reports any metadata
  other than `changes/<change>/.grove.yaml`. The spec guide says the first
  heading is the title and Purpose's first sentence the description, the
  template's Purpose comment says the same, and the archive guide writes
  the new spec from its template alone.
- `internal/cli/walkthrough_test.go`: `TestCheck` (leftover `.grove.yaml`
  reported once with no unknown field, also in a folder with no `spec.md`;
  `notes/ideas.md` neither listed nor reported; `untitled/spec.md` with no
  `# ` line reported), `TestArchive` (only the `spec.md` template; title
  and Purpose sentence in `grove specs` and the explore index, the full
  stop in `records.json` not ending it; Purpose deleted, no description
  and `grove check` clean), `TestSchemas` (a `spec` kind with spec
  metadata reported). `just check`: passed.

### T-02 — No design.md

Status: done

- The propose guide's "Write the change" says to write every decision in
  the proposal, however much room it needs; the proposal template's
  Decisions comment says a decision that needs room is written out there.
- On a fresh ledger fixture with this build, `grove instructions propose
  "x"` and `grove instructions propose --schema bugfix "x"` each contain
  no `design.md` (`grep -c` printed 0), and `grep -rn design.md internal`
  finds nothing.

### T-03 — Grove's own specs move over

Status: not started

### T-04 — Document specs without .grove.yaml

Status: not started

## Scenarios

## Reconciliation
