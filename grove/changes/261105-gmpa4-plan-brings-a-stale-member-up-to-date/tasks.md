## T-01 — Plan brings a stale member up to date

Result: the plan instructions, the planning docs and the walkthrough say
what "Plan a cycle" in [specs/cycles/spec.md](specs/cycles/spec.md) says,
as the proposal's decisions give it.

Files:
- `internal/instruct/requires/plan.md`
- `internal/cli/walkthrough_test.go`
- `docs/planning.md`
- `dev/walkthrough.md`

Delegated: test layout and private helper names

Interfaces: none.

Constraints: rule 12's instruction-text budget; the plan instructions
point to `grove instructions propose <change-id>` for how to bring a
member up to date rather than restating it.

Acceptance:
- `grove instructions plan` no longer says the cycle file is the only
  thing written; it says that for a new cycle, before writing it, the
  agent runs `grove check` and, for each member it reports stale on a
  requirement or for a named file changed since its base, tells the owner
  this edits approved work and brings it up to date on this branch as
  `grove instructions propose <change-id>` says, then goes on with the
  plan. The wording test checks the new text and the absence of the old.
- `docs/planning.md` "Plan a cycle" says so, and that one merge then
  brings the revision and the plan.
- The walkthrough's "Plan a cycle" step has a member `grove check`
  reports changed since its base; the session brings it up to date and
  sets its `base`, `git status` shows that member's files beside the cycle
  file, and `grove check` is clean.
- `just check` passes.
