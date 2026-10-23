## T-01 — Archive adds a History line

Result: the maintained spec template ends with `## History`, and the
archive instructions add a line to it, as
[specs/requirements/spec.md](specs/requirements/spec.md) "Archive adds a
History line" and "Archive does not compare History" say.

Files: `internal/instruct/schemas/default/documents/spec/template.md`,
`internal/instruct/schemas/default/documents/spec/guide.md`,
`internal/instruct/schemas/default/skills/archive.md`,
`internal/cli/walkthrough_test.go`

Interfaces: none; it is instruction text.

Constraints: `grove archive` and `grove check` gain no code for History.

Acceptance:
- `grove instructions archive` for a finished change prints the History
  line's form, says it goes at the end of History, under "Reconcile the
  specs", the finishing work a dropped change skips, and prints a
  maintained spec template that ends with `## History`.
- In `walkthrough_test.go`, on the fixture: the archive instructions
  name History; the new `record-editing` spec, filled from the printed
  template with a History line for the change, passes `grove check` and
  `grove archive`; a fold with no History line archives too.

## T-02 — Document History

Result: the documentation and the walkthrough say a maintained spec's
History lists the changes that shaped it.

Files: `docs/concepts.md`, `docs/reference.md`, `dev/walkthrough.md`

Interfaces: none.

Acceptance:
- `docs/concepts.md`, "Why archive folds a change into the specs", says
  each maintained spec's History links the archived changes that shaped
  it, where their reasons are.
- `docs/reference.md`'s `grove instructions archive` row says the fold
  adds the change's History line.
- `dev/walkthrough.md` step 5 has the session add a History line for
  261015-k3m9q to the new `record-editing` spec, and step 6 shows it in
  the merged spec; "A bug" has `/grove-archive` add the bugfix change's
  line below it.

## T-03 — History for Grove's own specs

Result: each maintained spec in `grove/specs/` ends with a History
listing the changes in `grove/changes/archive/` that have a change spec
for it, as "Archive adds a History line" gives the line.

Files: `grove/specs/apply/spec.md`, `grove/specs/config/spec.md`,
`grove/specs/cycles/spec.md`, `grove/specs/implement/spec.md`,
`grove/specs/propose/spec.md`, `grove/specs/prototype/spec.md`,
`grove/specs/runs/spec.md`, `grove/specs/requirements/spec.md`

Interfaces: none.

Constraints: no requirement or other section of these specs changes.

Acceptance:
- Each spec's History names exactly the archived changes with a
  `specs/<capability>/` folder for it, in the order they were merged into
  main, by ID within one merge. At 0c6f419 that is: `apply` 261022-1grb0;
  `config` 261023-cnnhm; `cycles` 261022-1grb0, 261023-rcxv2,
  261023-6j5bw, 261023-3f9h6; `implement` 261022-d1y6d, 261022-1grb0,
  261023-4044p, 261023-6j5bw; `propose` 261023-6j5bw, 261023-fca8x;
  `prototype` 261022-ka8bx; `runs` 261021-2aj5w, 261022-1grb0,
  261023-rcxv2, 261023-3f9h6, 261023-yz5bk. Changes archived since then
  are listed too, such as 261023-bfwp5 in the `requirements` spec its
  archive creates.
- Every link opens an archived change's folder.
- `grove check` prints "No problems in grove/.".
