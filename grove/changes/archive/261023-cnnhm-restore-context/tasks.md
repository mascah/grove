## T-01 — Restore context

Result: `context` behaves as [specs/config/spec.md](specs/config/spec.md)
says, and the documentation specifies it.

Files: `internal/cli/project.go`, `internal/cli/check.go`, `internal/cli/init.go`, `internal/instruct/instruct.go`, `internal/cli/walkthrough_test.go`, `docs/reference.md`, `docs/getting-started.md`, `docs/concepts.md`, `dev/walkthrough.md`

Interfaces: `instruct.Project.Context`, `grove.yaml`'s `context`.

Constraints: undo e3826d0 where it removed `context`, keeping its `and`
field list in `check`'s message.

Acceptance:
- `docs/reference.md`'s field table and example have `context` with the
  decided wording; `docs/getting-started.md` has a step writing it;
  `docs/concepts.md`'s glossary has it; `dev/walkthrough.md`'s setup step
  writes the decided line and its check step lists `context` (rule 8).
- `internal/cli/walkthrough_test.go` sets a context in `TestInit` and sees
  the "Project context" section in an instructions output, and `TestCheck`
  sees no unknown field for it.
- `just check` passes.
