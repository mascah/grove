# Implementation: A built-in bugfix schema

<!-- Progress on change 261023-vqsrt, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Embed the bugfix schema

Status: done

- `internal/instruct/schemas/bugfix/`: `schema.yaml` (its description, with
  a header comment), `skills/propose.md`, `documents/proposal/template.md`
  and `documents/proposal/guide.md`, copied word for word from
  `docs/reference.md`'s "Example: a bug schema" at 57af2e3. The example
  already had the sections Decisions names (Observed, Expected,
  Reproduction, The requirement it breaks, Open questions), the first task
  failing before the fix, and no change spec for a fix that restores a
  requirement, so nothing needed trimming beyond what it holds.
- `internal/instruct/instruct.go`: `SchemaNames` returns `default`,
  `bugfix`, then the project's other folders. A schema's files and
  `schema.yaml` come from its places, first first: the project's folder
  for it, Grove's built-in one, then the project's `default` and Grove's
  (`places`, used by `load`, `file` and `Check`); `load` reads each place's
  `schema.yaml` field by field, so the built-in `bugfix` names only its
  description and keeps `default`'s kinds and `writes`. `Check` counts only
  the project's own folders as places whose files are used.
- `internal/cli/walkthrough_test.go`: `TestSchemas` starts with no
  `grove/schemas/`: propose lists `default` then `bugfix` with its
  description; `--schema bugfix` prints the built-in propose guide (a first
  task failing before the fix, "Write no change spec") and the template's
  five sections, with the project's `default` tasks template reaching it;
  a project's `bugfix/documents/proposal/template.md` replaces the built-in
  one while the built-in guide stays; `grove check` is clean with a change
  of `schema: bugfix`. The existing checks of a mistake in
  `grove/schemas/bugfix/` and of an unknown schema keep their messages.
- On a fresh ledger fixture, with `grove/schemas/bugfix/` copied from the
  example in `docs/reference.md` at 57af2e3 (its propose guide marked to
  tell it apart), `grove check` printed "No problems in grove/." and
  `grove instructions propose --schema bugfix "x"` printed the project's
  guide.
- `grove check` on Grove's own repository prints "No problems in grove/.".

### T-02 — Explore and propose suggest bugfix

Status: in progress

- `internal/instruct/requires/explore.md`: when the owner describes
  behavior that is wrong, suggest the `bugfix` schema listed below, used
  only once the owner confirms; the proposal then adds `--schema bugfix`.
- `internal/instruct/requires/propose.md`: for a new change whose
  instructions do not follow `bugfix`, suggest it when the owner describes
  wrong behavior, and run `instructions propose --schema <name>` once the
  owner confirms it or names another listed kind. The old condition, "if
  the project's change schemas are listed", is gone: with `bugfix` built in
  they always are.
- `internal/cli/walkthrough_test.go`: `TestSchemas` checks that `grove
  instructions explore` and `grove instructions propose "Editing a missing
  record crashes"`, on the fixture with no `grove/schemas/`, list both
  schemas and say to suggest `bugfix` and use it once the owner confirms.

### T-03 — Document the built-in

Status: done

- `docs/reference.md`, "Change schemas": Grove has `bugfix` and `default`;
  "A schema of your own" keeps a short example (a schema whose propose
  writes no change specs), says the instructions always list the schemas,
  Grove's first, and that the agent suggests `bugfix` for wrong behavior;
  "Example: a bug schema" is now "Grove's bugfix", describing the built-in,
  the four places a `bugfix` file comes from, and replacing one of its files
  in `grove/schemas/bugfix/`. `check`'s "none of the three places" is "none
  of its places".
- `docs/concepts.md`, "Change schemas" and the glossary: Grove has
  `default` and `bugfix`, and the agent suggests `bugfix` for a bug.
- `dev/walkthrough.md`: "A bug" starts from a fixture with no
  `grove/schemas/`, and the session suggests `bugfix` itself; "A mistake in
  a schema" writes a `grove/schemas/bugfix/` of the tester's own with the
  same two mistakes and gets the same messages, and is where review's
  `change schema` label is seen.
- `CLAUDE.md`'s layout table names `schemas/default/` and
  `schemas/bugfix/`.
