## T-01 — Implement and claim name a missing base or planned commit

Result: `grove instructions implement <change-id>`, attended or
`--unattended`, `grove instructions implement <cycle-id>` and
`grove claim <cycle-id>` refuse with the line `grove check` prints for the
missing `base` or `planned` commit, and print no instructions, when that
commit is not in the repository; a `grove claim` given no cycle still
skips such a cycle silently and records nothing.

Files:
- `internal/cli/instructions.go`
- `internal/cli/claim.go`
- `internal/cli/walkthrough_test.go`
- `dev/walkthrough.md`
- `docs/guides/run-grove-in-ci.md`

Delegated: test layout and private helpers; `internal/cli/check.go`,
should the `missing` sentence need a path prefix

Interfaces: the error text of `grove instructions implement` and
`grove claim <cycle-id>` for a missing commit is `check`'s line for that
document, prefixed `grove: `.

Constraints: one sentence per message (rule 14); handling the error Git
returns, no guard beyond it (rule 13); the Go and test line budgets of
rule 12 hold.

Acceptance:
- In a `git clone --depth 1` of the fixture holding a change in no cycle,
  a change in a cycle not done, and that cycle, each command above exits
  non-zero and prints only `grove: ` and the line `grove check` prints
  for that document, in the shallow-clone wording; `grove claim
  <cycle-id>` records nothing and creates no branch.
- In a full clone with a `base` or `planned` set to a commit that does
  not exist, the same, in the "is not a commit in this repository"
  wording.
- `grove claim` given no cycle prints no Git error for such a cycle and
  records nothing.
- A test in `internal/cli/walkthrough_test.go` fails before the fix and
  passes with it.
- `dev/walkthrough.md`'s shallow-clone step shows the implement and claim
  lines run on the fixture, and `docs/guides/run-grove-in-ci.md` "What to
  expect" says the implement instructions fail as `check` does.
- `just check` passes.
