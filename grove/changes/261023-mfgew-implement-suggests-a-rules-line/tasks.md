## T-01 — Implement suggests a rules line

Result: `/grove-implement` behaves as
[specs/implement/spec.md](specs/implement/spec.md) "Suggest a rules line"
says.

Files: `internal/instruct/requires/implement-owner.md`, `docs/reference.md`,
`docs/getting-started.md`, `dev/walkthrough.md`,
`internal/cli/walkthrough_test.go`

Interfaces: none; it is instruction text.

Constraints: the instruction text stays within rule 12's word budget as
PR #47 proposes it.

Acceptance:
- `docs/reference.md`, beside `context` and `rules`, and
  `docs/getting-started.md`, step 5, say in plain words that implement may
  suggest a line and writes it only with the owner's agreement (rule 8).
- On the ledger fixture, in Claude Code, after I correct the agent twice,
  in ways that would apply to future work, while it builds a change, it
  quotes a correction and suggests one line, not two; when
  I agree, `grove.yaml` on the branch holds it and `grove check` prints
  "No problems in grove/."; when I decline, `grove.yaml` is unchanged.
- On the ledger fixture, the main line's `/grove-implement`, with no
  correction, suggests no line.
- `walkthrough_test.go` checks that `instructions implement <change-id>`
  prints the rule, and `instructions implement <change-id> --unattended`
  does not.

## T-02 — Unattended names a line, never writes it

Result: `/grove-implement --unattended` behaves as
[specs/implement/spec.md](specs/implement/spec.md) "Name a rules line
unattended" says.

Files: `internal/instruct/requires/implement-unattended.md`,
`docs/getting-started.md`, `dev/walkthrough.md`,
`internal/cli/walkthrough_test.go`

Interfaces: none; it is instruction text.

Acceptance:
- `docs/getting-started.md`, step 5, says `/grove-implement --unattended`
  never edits `grove.yaml` but may name a line in its report (rule 8).
- On the ledger fixture, `claude -p "/grove-implement <cycle-id>
  --unattended"` leaves `grove.yaml` unchanged, and any line it suggests
  is in its last message.
- `walkthrough_test.go` checks that `instructions implement <change-id>
  --unattended` prints the rule, and `instructions implement <change-id>`
  does not.
