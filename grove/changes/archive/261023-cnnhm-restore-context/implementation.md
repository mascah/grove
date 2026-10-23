# Implementation: Restore context in grove.yaml

<!-- Progress on change 261023-cnnhm, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Restore context

Status: done

- Code: the code part of e3826d0, reverse-applied cleanly, brings back
  `project.Context`, `instruct.Project.Context`, the "Project context"
  section before "What Grove requires of <skill>" in
  `Instructions.Markdown`, `context` in `configFields`, and the `init`
  template's `context: ""` with its comment. `check`'s message keeps
  `and(configFields)`.
- `internal/cli/walkthrough_test.go`: `TestInit` sets
  `context: Run python3 -m unittest before saying a task is done.` and sees
  the section, then "## What Grove requires of propose", in the propose
  instructions. `TestCheck` reports no unknown field for `context`, which
  the `init` template writes, and lists it in the unknown-field message.
- On the ledger fixture: with that `context` in `grove.yaml`,
  `grove check` is clean, and `grove instructions apply <cycle> <change>`
  prints "## Project context" with the line before
  "## What Grove requires of apply".
- Docs: `docs/reference.md` (example, field table with the decided
  wording, and the note on adding to Grove's instructions),
  `docs/getting-started.md` §1 step, `docs/concepts.md` glossary, and
  `dev/walkthrough.md` setup step, check field list and assumptions.
- `just check`: passes; 3532 Go lines, 1172 test lines.

## Reconciliation

- config: applied to specs/config/spec.md (new)
