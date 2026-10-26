## T-01 — Revising with propose re-bases a change

Result: `grove check`, propose and the owner's implement instructions
behave as [specs/propose/spec.md](specs/propose/spec.md) "Say that
revising re-bases" and "Re-base a change whose named file moved", and
[specs/implement/spec.md](specs/implement/spec.md) "Recommend re-basing
after a dependency lands", say.

Files: `internal/cli/check.go`, `internal/instruct/requires/propose.md`,
`internal/instruct/requires/implement-owner.md`, `docs/coordinating.md`,
`dev/walkthrough.md`, `internal/cli/walkthrough_test.go`

Interfaces: `grove check`'s moved-file message, as the requirement
quotes it.

Acceptance:
- On the fixture, `grove check` prints the new message for a file changed
  since base, and still exits non-zero.
- With the change B and A of "A dependency landed", `/grove-propose B`
  sets `base` to the target's commit and `grove check` is then clean.
- `/grove-implement B` with the owner recommends revising with propose at
  T-01; `--unattended` still marks T-01 blocked.
- `docs/coordinating.md` and the walkthrough quote the new message and
  say to re-base with propose once a dependency lands (rule 8).
- `walkthrough_test.go`'s expected message is the new one.
- `just check` passes, the instruction text included.
