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

Status: done

`docs/concepts.md`, "Why archive folds a change into the specs", says each
maintained spec's History links the archived changes that shaped it, where
the reasons are. `docs/reference.md`'s `grove instructions archive` row
says the fold adds the change's History line. `dev/walkthrough.md` step 5
has the session end the new `record-editing` spec with 261015-k3m9q's
History line, step 6 shows it in the merged spec, and "A bug" has
`/grove-archive` add the bugfix change's line below it. Read against the
acceptance; documentation only, no test.

### T-03 — History for Grove's own specs

Status: done

The owner agreed to fix T-03's `Files:` line: `grove/specs/apply/spec.md`
no longer exists, and `changes`, `init`, `maintained-specs` and `schemas`
were missing. Each of the 11 specs in `grove/specs/` now ends with a
History listing the archived changes with a `specs/<capability>/` folder
for it, ordered by the first-parent commit on main that archived each, by
ID within one merge; 261023-psd1b, archived on this branch, is last in
`propose`. Each list begins with the acceptance's list for 0c6f419 (less
`apply`, which has no spec; 261022-1grb0, 261023-g0a3h and 261023-m2she
had `apply` change specs). `git diff` adds lines only; every link names a
folder in `grove/changes/archive/`; `grove check`: "No problems in
grove/."

## Scenarios

## Reconciliation
