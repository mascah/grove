## T-01 — `grove changes --json` prints the project value

Result: `grove changes --json` prints the object `grove serve` answers
`/api/project` with, every key it prints today kept in shape, so one
export of the project serves scripts and the page alike.

Files:
- `internal/cli/changes.go`
- `internal/cli/serve.go`
- `internal/cli/walkthrough_test.go`

Delegated: test layout and private helpers

Interfaces: the JSON keys `read`, `specs`, `cycles`, `branches`, `runs`,
`history` and `problems`, and `cycle`, `ready` and `files` on each change,
as `serve` names them.

Constraints: the keys and shapes `grove changes --json` prints today do
not change; the Go and test line budgets of rule 12 hold.

Acceptance:
- Every scenario of the changes and serve change specs holds on the
  fixture.
- The guide "Run Grove in CI"'s `jq` over `.changes[]` works unchanged
  on the new output.
- `just check` passes.

## T-02 — Documentation and the walkthrough

Result: "The project as data" in the documentation shows the value's
keys, says the page reads the same object, and the walkthrough runs the
step.

Files:
- `docs/coordinating.md`
- `docs/watching.md`
- `docs/reference.md`
- `dev/walkthrough.md`

Acceptance:
- `docs/coordinating.md` lists `cycles`, `branches`, `runs`, `history`
  and `problems` with one line each and names `grove serve` as reading the
  same value.
- Each documented step is run on the fixture and recorded in
  `dev/walkthrough.md`.
