# W-019 plan: terms, plans, reviews, and a discoverable brief

Owner: [W-019](../../grove/work/W-019-knowledge-artifacts.md). Representation is
selected in [D-005](../../grove/decisions/D-005-typed-knowledge-records.md); this
plan settles the items that record left to preparation. Base: main `d9fc2a5`,
branch `worktree-W-019`.

## Observed

- A type is spread over small tables: loader folders
  (`internal/project/project.go:67`), `validID` and lifecycles
  (`metadata.go:162-165,231`), `create.kinds` and `idLine`
  (`create.go:26-36`), the counter reader's `WQD` check (`create.go:198`),
  update's `allowed`/`lists`/`strs`/`fields` (`update.go:130-136,248`), and
  `idForm`/`compareIDs` in `internal/versions`.
- The board shows only `Type == "work"` (`internal/tui/model.go:615`); new
  types never become cards. `list`, `show`, `versions` are type-agnostic.
- `context` reads selected work and lists related records by role
  (`internal/handoff/context.go:236-262`); selection refuses non-work.
- Committed sources give the loader only `grove.yaml` and the record root
  (`internal/versions/tree.go:248-285`), one `git cat-file` process (W-013).
- The shared counter file is read strictly: an older binary calls a `T` line
  corrupt and refuses `new`.

## Design

1. **One type table** in `internal/project`: name, prefix, folder, statuses
   (first is `new`'s initial status), and the schema version that introduces
   it. The loader, `validID`, create, update and versions read it instead of
   their own copies. Body skeletons stay in `create`.

   | Type | Prefix | Folder | Statuses | Extra fields |
   | --- | --- | --- | --- | --- |
   | term | `T-` | `terms/` | `proposed`, `settled` | none |
   | plan | `P-` | `plans/` | `current`, `superseded` | `work` |
   | review | `R-` | `reviews/` | `current`, `superseded` | `work`, `examined` |

   `work` is an optional list of work IDs, validated like `depends_on`
   targets (must resolve, must be work, no duplicates). It is optional because
   `new` takes no fields; `update --set 'work=[…]'` sets it. `examined` is an
   optional quoted Git commit (7 to 40 lowercase hex): what the review looked
   at. Staleness is a comparison a reader makes, not stored state; W-020 owns
   approval and candidate semantics and may extend the review type. Two terms
   with the same title (case-insensitive) are an error.
2. **Schema 2.** The CLI accepts `schema_version` 1 and 2. Version 1 keeps
   exactly today's rules, so existing projects stay readable and nothing is
   rewritten. The new types and the `brief` key exist only in 2; migrating is
   the deliberate one-line edit of `grove.yaml`, documented in the record
   model. An older binary refuses 2 with its existing "unsupported version"
   message instead of reporting new folders as misplaced records.
3. **Brief.** Optional `brief: PATH` in `grove.yaml`: a clean project-relative
   `.md` path without `..`, anywhere in the project including directly under
   the record root, but not inside a type folder. The loader exempts that one
   path from the misplaced-Markdown rule. `LoadFS` validates the value's form
   only, so committed trees need no extra reads; live `Load` (and so `check`
   and every live command) requires a regular non-symlink file there. New
   `grove brief` prints it like `show`: bytes to stdout, path to stderr;
   `--json` gives `{path, revision, source}`. `context` never includes it;
   `--include` still can. This repository sets `brief: docs/restart-brief.md`;
   W-029 owns moving it.
4. **Context.** `context W-NNN` lists each plan and review whose `work` names
   a selected ID, with role `plan for W-NNN` / `review of W-NNN`, not included.
   The scope notice says so. No structural change, so `format_version` stays 2.
5. **Terms.** Create work, preparation, attempt, review, approval,
   integration, source and revision with `grove new term`, written from the
   brief, record model and guides, status `proposed` until the owner judges
   them. Meaning and relationships only.

## Steps

1. Type table, schema gating, new fields and validation, brief key, with
   loader and metadata tests (schema 1 refuses the new folders and key; schema
   2 accepts them; bad `work`/`examined`/brief values name file and field).
2. create, update, versions, CLI: table-driven prefixes and counters, update's
   field lists (`examined` always quoted), `grove brief`, help text; tests in
   disposable clones via `--project`.
3. `context` listing and notice, with tests for a shared plan.
4. Reconcile `docs/record-model.md`, README, `AGENTS.md`, both guides: new
   plans and reviews are `grove new plan|review` records from here on;
   existing `docs/plans` and `docs/reviews` stay until W-029.
5. Last, because it touches shared state: `grove.yaml` to schema 2 with
   `brief`, the eight terms, and this work's review as the first `R-` record.
6. Full verification per `AGENTS.md`, then one independent review of the
   final combined revision; fix rounds capped at three.

## Known limits

- From step 5 until this branch is merged, `grove new` run by older code in
  another checkout of this repository fails on the shared counter file's `T`
  or `R` line, and older code reports this branch as unsupported schema 2.
  Removing the counter file is safe (it reinitializes from a scan) but loses
  unwritten reservations. Merging resolves both.
- A brief outside the record root is not existence-checked on committed
  sources in `versions` or the board.
- This plan itself stays in `docs/plans/`: it predates the support. W-029
  migrates it with the rest.
