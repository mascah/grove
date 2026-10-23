## T-01 — Embed the bugfix schema

Result: every project has `bugfix` as
[specs/schemas/spec.md](specs/schemas/spec.md) "Grove has a built-in
bugfix schema", "Bugfix proposes a bug" and "Where a bugfix file comes
from" say, its files taken from `docs/reference.md`'s "Example: a bug
schema" at 57af2e3 as Decisions trims it.

Files: `internal/instruct/schemas/bugfix/schema.yaml`,
`internal/instruct/schemas/bugfix/skills/propose.md`,
`internal/instruct/schemas/bugfix/documents/proposal/template.md`,
`internal/instruct/schemas/bugfix/documents/proposal/guide.md`,
`internal/instruct/instruct.go`, `internal/cli/walkthrough_test.go`

Interfaces: `Project.SchemaNames` returns `default`, `bugfix`, then the
project's other schemas; a schema's files and `schema.yaml` are read from
the four places in the order the requirement gives, for `Project.Check`,
`Layout`, `Kinds` and every skill's instructions alike.

Constraints: the built-in `bugfix` names only what differs from `default`,
so it keeps `default`'s kinds and each skill's `writes`; it passes the
checks `grove check` makes of any schema.

Acceptance:
- In `walkthrough_test.go`'s `TestSchemas`, on the fixture with no
  `grove/schemas/`: `grove instructions propose` lists `default` then
  `bugfix`; `--schema bugfix` prints the built-in propose guide, which
  says the first task adds a test that fails before the fix and that a fix
  restoring a maintained requirement writes no change spec, and the
  proposal template with Observed, Expected, Reproduction, The requirement
  it breaks and Open questions; `grove check` is clean with a change of
  `schema: bugfix`.
- In the same test: a project's `grove/schemas/default/documents/tasks/template.md`
  reaches `bugfix`; a project's `grove/schemas/bugfix/documents/proposal/template.md`
  replaces the built-in one while the built-in propose guide stays; the
  existing checks of a mistake in `grove/schemas/bugfix/` and of an unknown
  schema keep their messages.
- On the ledger fixture, `grove/schemas/bugfix/` copied from the example in
  `docs/reference.md` at 57af2e3 leaves `grove check` clean, and `grove
  instructions propose --schema bugfix "x"` uses its files.
- `grove check` on Grove's own repository stays clean.

## T-02 — Explore and propose suggest bugfix

Result: the explore and propose instructions say what
[specs/schemas/spec.md](specs/schemas/spec.md) "Explore and propose
suggest bugfix" says.

Files: `internal/instruct/requires/explore.md`,
`internal/instruct/requires/propose.md`, `internal/cli/walkthrough_test.go`

Interfaces: none.

Acceptance:
- On the ledger fixture, `grove instructions explore` and `grove
  instructions propose "Editing a missing record crashes"` list both
  schemas and say to suggest `bugfix` for behavior that is wrong and to
  use it once the owner confirms; `walkthrough_test.go` checks both.
- In Claude Code on a fresh fixture, `/grove-propose Editing a missing
  record crashes` suggests `bugfix` and, once confirmed, writes a change
  with `schema: bugfix`.

## T-03 — Document the built-in

Result: the documentation describes `bugfix` as Grove's, as
Decisions says for `docs/reference.md`, and the walkthrough runs it
without writing the schema.

Files: `docs/reference.md`, `docs/concepts.md`, `dev/walkthrough.md`,
`CLAUDE.md`

Interfaces: none.

Acceptance:
- `docs/reference.md`, "Change schemas", names both built-in schemas and
  the four places a `bugfix` file comes from; "Example: a bug schema"
  describes the built-in `bugfix` and how to replace one of its files in
  `grove/schemas/bugfix/`; "A schema of your own" keeps a short example
  (rule 8).
- `docs/concepts.md`, "Change schemas" and the glossary, says Grove has
  `default` and `bugfix`.
- `dev/walkthrough.md`'s "A bug" starts from a fixture with no
  `grove/schemas/` and the session suggests `bugfix` itself; "A mistake in
  a schema" makes its mistakes in a `grove/schemas/bugfix/` of the tester's
  own and gets the same messages.
- `CLAUDE.md`'s layout table names `schemas/default/` and
  `schemas/bugfix/`.
