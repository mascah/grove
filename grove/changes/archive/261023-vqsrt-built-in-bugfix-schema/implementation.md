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

Status: done

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
- Trial in Claude Code 2.1.310 on a fresh ledger fixture with no
  `grove/schemas/`, this build first on PATH, headless over three turns of
  one session: `/grove-propose Editing a missing record crashes` ran
  `grove instructions propose "..."` without `--schema`, reproduced the
  `StopIteration` traceback, and asked "this is behavior that's wrong, so
  I suggest the `bugfix` change schema instead of `default`. Use
  `bugfix`?", writing nothing. Confirmed, it ran `grove instructions
  propose --schema bugfix "..."`, asked its round, and on "Yes." wrote a
  change with `schema: bugfix`, the proposal's Observed, Expected,
  Reproduction, The requirement it breaks and Open questions, a change
  spec ADDing the requirement since no maintained spec covers it and the
  owner agreed, and T-01 whose acceptance says its test fails before the
  fix. `grove check` was clean.

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
- Recheck before finishing: T-01 to T-03's evidence holds against the code
  as it now is. `just check` passes (3732 Go lines, 1356 test lines);
  `grove check` prints "No problems in grove/.".

## Scenarios

- schemas "A project with no schemas of its own": `walkthrough_test.go` `TestSchemas`, passed.
- schemas "Proposing with bugfix": `TestSchemas`, passed (`schema: bugfix` in the `.grove.yaml` template).
- schemas "A change with schema bugfix": `TestSchemas`, passed.
- schemas "A schema that does not exist": `TestSchemas`, passed.
- schemas "The bugfix proposal template": `TestSchemas`, passed (the five headings, and the project's `default` tasks template reaching it, so the tasks template is `default`'s).
- schemas "A fix that restores a requirement": `TestSchemas`, passed ("Write no change spec" and the failing first test in the guide).
- schemas "Behavior no requirement gives": the trial above, passed (asked in the round, ADDed on "yes"); the guide text is the reference's example, word for word.
- schemas "The project's own default file": `TestSchemas`, passed.
- schemas "The project's own bugfix file": `TestSchemas`, passed.
- schemas "A project that copied the example": by hand on a fresh fixture, passed (see T-01).
- schemas "A mistake in the project's bugfix": `TestSchemas`, passed.
- schemas "Proposing a bug fix": `TestSchemas`, passed; the trial above, passed.
- schemas "Exploring a bug": `TestSchemas`, passed.

## Reconciliation

- schemas: created specs/schemas/spec.md and its .grove.yaml, with the change's requirements; its "Observed current behavior" is resolved and not carried over.
