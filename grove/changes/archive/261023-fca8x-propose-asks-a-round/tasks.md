## T-01 — Propose asks a round

Result: `/grove-propose` behaves as
[specs/propose/spec.md](specs/propose/spec.md) "Ask a round before
writing" says, the default schema's proposal has no assumptions section,
and the documentation specifies it.

Files: `internal/instruct/requires/propose.md`,
`internal/instruct/schemas/default/skills/propose.md`,
`internal/instruct/schemas/default/documents/proposal/guide.md`,
`internal/instruct/schemas/default/documents/proposal/template.md`,
`docs/concepts.md`, `docs/reference.md`, `dev/walkthrough.md`,
`internal/cli/walkthrough_test.go`

Interfaces: none; the round is instruction text, and Grove parses nothing
new.

Constraints: the explore instructions keep "One focused question at a
time"; the round comes after the scope is confirmed, not before.

Acceptance:
- `docs/concepts.md` and `docs/reference.md` describe the round and no
  longer mention assumptions to confirm in review; the bugfix example in
  `docs/reference.md` follows; `dev/walkthrough.md` proposes a change
  through a round (rule 8).
- On the ledger fixture, in Claude Code: `/grove-propose` from an explore
  conversation lists its would-be assumptions with recommended answers,
  waits, and writes a proposal whose decisions quote the answers, with no
  assumptions section.
- `internal/cli/walkthrough_test.go` checks that `instructions propose`
  prints the rule; `just check` passes.
