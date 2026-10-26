## T-01 — Quote the title in the change and cycle templates

Result: a title with a colon, written into the printed `.grove.yaml` or
cycle template, passes `grove check`.

Files:
`internal/instruct/metadata/change.grove.yaml`,
`internal/instruct/metadata/cycle.yaml`,
`internal/cli/walkthrough_test.go`, `docs/planning.md`,
`docs/coordinating.md`

Interfaces: both templates print `title: "<title>"`.

Acceptance:
- A test fills the template that `grove instructions propose` prints with
  a title holding `: `, and `grove check` reports no problem; it fails
  before the fix with "invalid YAML".
- The test's assertion on the printed template expects `title: "<title>"`.
- The cycle template that `grove instructions plan` prints shows
  `title: "<title>"`.
- The example `.grove.yaml` and cycle file in `docs/planning.md` and
  `docs/coordinating.md` show the title quoted.
- `just check` passes, the instruction text included.
