## T-01 — The implement skill

Result: `/grove-implement <change-id>` behaves as
[specs/implement/spec.md](specs/implement/spec.md) says, and the
documentation specifies it.

Constraints: within the rule 12 budget; messages are one sentence (rule
14). Apply's change-ID and no-ID forms keep working.

Acceptance:
- `grove init` writes the implement skill for each harness under `agents`;
  `rules.skills.implement` is accepted.
- `docs/getting-started.md` builds its change with `/grove-implement`;
  `docs/reference.md` lists the skill and `grove instructions implement
  <change-id>`; `docs/concepts.md` says implement asks the owner and may
  revise the change on its branch, and apply asks no one.
  `dev/walkthrough.md` runs those steps (rule 8).
- On the ledger fixture, in Claude Code: a change with an open question
  deciding a task is implemented; the agent asks; the owner's answer
  changes the change spec, the owner agrees, the agent edits it, `grove
  check` is clean, and `grove review` shows the edit beside the code. The
  same change under `/grove-apply` marks the task blocked.
- `internal/cli/walkthrough_test.go` covers `grove instructions implement`
  refusing a draft and a change that waits; `just check` passes.
