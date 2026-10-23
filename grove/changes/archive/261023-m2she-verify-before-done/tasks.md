## T-01 — Propose checks the change before presenting it

Result: `/grove-propose` behaves as
[specs/propose/spec.md](specs/propose/spec.md) "Check the change before
presenting it" says.

Files: `internal/instruct/requires/propose.md`,
`internal/instruct/schemas/default/skills/propose.md`,
`docs/getting-started.md`, `docs/reference.md`, `dev/walkthrough.md`,
`internal/cli/walkthrough_test.go`

Interfaces: none; it is instruction text.

Acceptance:
- `docs/getting-started.md`, step 3, and the `instructions propose` row
  of `docs/reference.md` describe the check (rule 8).
- On the ledger fixture, in Claude Code, a change proposed with a
  requirement no task covers is fixed before the agent presents it, and
  the agent says what it changed.
- `walkthrough_test.go` checks that `instructions propose` prints the
  rule.

## T-02 — Implement and apply verify before finishing

Result: implement, apply and archive behave as
[specs/implement/spec.md](specs/implement/spec.md) "Verify before
finishing" and "Archive names unverified scenarios" and
[specs/apply/spec.md](specs/apply/spec.md) "Verify before finishing a
member" say, and the implementation template has the scenario list.

Files: `internal/instruct/requires/implement.md`,
`internal/instruct/requires/apply.md`,
`internal/instruct/requires/archive.md`,
`internal/instruct/schemas/default/skills/implement.md`,
`internal/instruct/schemas/default/skills/apply.md`,
`internal/instruct/schemas/default/documents/implementation/template.md`,
`internal/instruct/schemas/default/documents/implementation/guide.md`,
`docs/concepts.md`, `docs/getting-started.md`, `dev/walkthrough.md`,
`internal/cli/walkthrough_test.go`

Interfaces: `implementation.md`'s `## Scenarios` section, read by people
and agents; Grove parses nothing new.

Constraints: the `Status:` lines stay exactly as the template has them
(rule 3).

Acceptance:
- `docs/concepts.md`'s change table and `docs/getting-started.md`, step
  5, describe the scenario list and the recheck (rule 8).
- On the ledger fixture, in Claude Code, `/grove-implement` writes a
  scenario list naming each scenario's test, and archive names an
  unverified one before folding.
- `walkthrough_test.go` checks that `instructions implement`, `instructions
  apply` and `instructions archive` print their rules.
