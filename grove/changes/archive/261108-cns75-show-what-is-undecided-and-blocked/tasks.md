## T-01 — Mark open questions and blocked tasks

Result: `grove changes`, `grove changes --json` and the `grove serve`
page show which changes have open questions and which tasks are blocked,
as the changes and serve change specs say.

Files:
- `internal/doc/doc.go`
- `internal/cli/changes.go`
- `internal/cli/serve.go`
- `internal/cli/page.html`
- `internal/cli/walkthrough_test.go`

Delegated: test layout and private helpers; a small direct test of the
"None" rule in `internal/doc/`

Interfaces: the JSON key `open_questions` on each change and on each
change of a branch row.

Constraints: no key `grove changes --json` or `/api/project` prints
today changes shape; the Go, test and page budgets of rule 12 hold.

Acceptance:
- Every scenario of the changes and serve change specs holds on the
  fixture, the page's markers checked by hand in a browser as the serve
  spec's other page scenarios are.
- The guide "Run Grove in CI"'s `jq` over `.changes[]` works unchanged.
- `just check` passes.

## T-02 — Documentation and the walkthrough

Result: the documentation shows the markers and the field, and the
walkthrough runs those steps on the fixture.

Files:
- `docs/coordinating.md`
- `docs/watching.md`
- `docs/reference.md`
- `dev/walkthrough.md`

Acceptance:
- `docs/coordinating.md` shows a `grove changes` line with `open
  questions` and a blocked task, says the line reports what the target
  holds, and names `open_questions` under "The project as data".
- `docs/watching.md` says where the page shows the markers.
- `docs/reference.md`'s `grove changes` row names the markers.
- Each documented step is run on the fixture and recorded in
  `dev/walkthrough.md`.
