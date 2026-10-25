# Maintained specs list the changes that shaped them

## Problem

A maintained spec says what is agreed, not why. The reasons are in the
proposals of the changes that shaped it, kept in `changes/archive/`, but
nothing in the spec points to them. At 0c6f419, the maintained spec
template ends with Known discrepancies and Open questions
(`internal/instruct/schemas/default/documents/spec/template.md`), and the
archive guide folds requirements into the spec and leaves no trace of the
change (`internal/instruct/schemas/default/skills/archive.md`). To find
the decision behind a rule, a reader has to search the archive for change
specs that name the capability. In the Codex review of the 2026-10-23
plan, this was the gap: a later agent sees a constraint in a maintained
spec without finding the decision behind it.

## Intended outcome

Every maintained spec ends with a `## History` section: one line for each
change folded into it, in the order they were archived, with the change's
ID linked to its archived folder and its title. Archiving a finished
change adds its line to each maintained spec it folds into. History is
not a requirement, so the fold check from 261023-bfwp5 does not compare
it. Grove's own seven maintained specs get their History from the changes
already archived. The rules are in
[specs/requirements/spec.md](specs/requirements/spec.md).

## Scope

### In scope

- A closing `## History` section in the maintained spec template, and a
  line about it in the spec guide.
- The archive instructions adding one History line per maintained spec a
  finished change folds into.
- `docs/concepts.md` ("Why archive folds a change into the specs"),
  `docs/reference.md`, `dev/walkthrough.md` and
  `internal/cli/walkthrough_test.go`.
- History for Grove's own maintained specs in `grove/specs/`, from
  `grove/changes/archive/`.

### Out of scope

- `grove archive` or `grove check` reading, writing or checking History:
  the agent writes it, as it does the fold.
- A dropped change: it folds into no spec, so it adds no line.
- Other projects' maintained specs: they get History as their changes are
  archived.

## Decisions

The owner agreed to the direction on 2026-10-23 ("yes to all") and to
every answer below in one round ("yes to C & D").

- **The maintained spec template gets a closing `## History` section;
  archiving a finished change adds one line per change to each maintained
  spec it folds into, with the change ID, its title and a link to its
  archived folder.** Reason: the Codex review found that a later agent
  sees a constraint without finding the decision behind it; the archived
  proposals hold the reasoning.
- **History is not a requirement, and the fold check from 261023-bfwp5,
  which compares only requirements, ignores it.** Reason: a History line
  added or missing changes nothing agreed.
- **Grove's own seven maintained specs are backfilled from the archived
  changes in `grove/changes/archive/` that touched them,** each archived
  change's `specs/<capability>/` folders saying which. Other projects are
  not touched.

## Open questions

None must be answered before review. None blocks a task.

## Related documents

- Requirements: [specs/requirements/spec.md](specs/requirements/spec.md)
- Tasks: [tasks.md](tasks.md)
