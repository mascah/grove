# Reference

## Terminal commands

| Command | Does |
| --- | --- |
| `grove` | Opens the reader on the maintained specs and approved changes. |
| `grove init` | Writes `grove.yaml` if there is none, and the skills for each harness it lists. Safe to re-run: it rewrites the skills and leaves `grove.yaml` alone. Commits nothing. |
| `grove specs` | Lists the maintained specs on the target branch. |
| `grove changes` | Lists the approved changes on the target branch that are not archived. |
| `grove check` | Reports each problem in the Grove documents of this working tree, each unknown field or `rules` key in `grove.yaml`, and each mistake in a change schema, as `path:line: message`, and exits non-zero if there is one. |
| `grove review` | Shows what merging the current branch into the target would approve: each change, maintained spec or change schema file it adds or alters, and the diff. |
| `grove archive <change-id>` | Moves the change's folder to `changes/archive/`. Refuses unless `implementation.md` marks every task in `tasks.md` `Status: done`, or says `Status: dropped`. Commits nothing. |

Skills run these; you do not need to:

| Command | Prints the instructions to |
| --- | --- |
| `grove instructions explore` | Interview the user about an idea, with the project's specs and changes listed. |
| `grove instructions propose [--schema <name>] ["<description>"]` | Write a new change from the conversation or the description, with the project's change schema or the one named. |
| `grove instructions propose <change-id>` | Revise a drafted change, with the schema it has. |
| `grove instructions apply <change-id>` | Build an approved change and record progress. |
| `grove instructions archive <change-id>` | Fold the change's specs into the maintained specs, then run `grove archive`. |

## Skills

`grove init` writes one skill for each of Grove's four skills, for each harness in `grove.yaml`.

| Skill | Claude Code | Codex |
| --- | --- | --- |
| explore | `/grove-explore` | `$grove-explore` |
| propose | `/grove-propose` | `$grove-propose` |
| apply | `/grove-apply` | `$grove-apply` |
| archive | `/grove-archive` | `$grove-archive` |

Claude Code reads them from `.claude/skills/grove-<skill>/SKILL.md`, Codex
from `.agents/skills/grove-<skill>/SKILL.md`. Commit them with your
project.

## grove.yaml

```yaml
root: grove
target: main
agents: [claude, codex]
context: |
  Go 1.24 service; tests with go test ./...
  Every public API stays backwards compatible.
rules:
  documents:
    proposal:
      - List who is affected.
    tasks:
      - Every task names its test.
  skills:
    apply:
      - Commit after each task, naming its task ID.
```

| Field | Default | Meaning |
| --- | --- | --- |
| `root` | `grove` | The folder that holds `specs/` and `changes/`. |
| `target` | `main` | The branch whose contents count as agreed. Merging into it approves. |
| `agents` | `[claude, codex]` | The harnesses `init` writes skills for. |
| `schema` | `default` | The change schema a new change has unless the owner names another while proposing. See "Change schemas". |
| `context` | empty | Facts about the project, given to every skill. Keep it to what an agent cannot read from the code. |
| `rules` | none | Extra lines for Grove's instructions. Under `documents`, for one kind of document, given when that document is written: `proposal`, `change-spec`, `tasks`, `implementation`, `spec`, and those of the project's own schemas. Under `skills`, for one skill, given with its instructions: `explore`, `propose`, `apply`, `archive`. |

`grove check` reports a field or `rules` key that Grove does not know; other
commands ignore it.

`context` and `rules` add to Grove's own instructions; they cannot remove or
replace them. Short rules work better than long ones.

Grove's skills commit nothing unless you ask in the conversation or a rule
asks, as `apply` does above.

## Change schemas

A change schema says what a change consists of and how each skill goes: its
kinds of document, each with a path, a template and a guide, and for each
skill its guide and the kinds of document it writes. Grove's built-in schema
is `default`, the software change [Concepts](concepts.md) describes:

| Kind | Path, under the root | Written by |
| --- | --- | --- |
| `proposal` | `changes/<change>/proposal.md`, the change's main document | propose |
| `change-spec` | `changes/<change>/specs/<capability>/spec.md` | propose |
| `tasks` | `changes/<change>/tasks.md` | propose |
| `implementation` | `changes/<change>/implementation.md` | apply, archive |
| `spec` | `specs/<capability>/spec.md`, a maintained spec | archive |

A schema's folder holds `schema.yaml`, a folder for each kind of document,
`documents/<kind>/` with its `template.md` and `guide.md`, and the guide for
each skill, `skills/<skill>.md`.

A schema cannot change Grove's skills or what Grove requires of each, which
its instructions give first; that a merge approves a change; a change's
`.grove.yaml`; or how progress is recorded: each task is a `T-NN` heading in
`tasks.md`, and its status a `Status:` line in `implementation.md`, which
`grove archive` reads. So every schema has `tasks` and `implementation` at
those paths, keeps the `Status:` lines of Grove's implementation template,
and has one kind with `metadata: changes/<change>/.grove.yaml`, the change's
main document. Its other kinds are in the change folder, except `spec`.

To add lines to Grove's instructions rather than replace a file, use `rules`
in `grove.yaml`.

### Your own files for default

To change one of `default`'s files, put your own at the same path under
`grove/schemas/default/`, such as `grove/schemas/default/documents/tasks/template.md`.
It replaces Grove's file whole; every file you do not put there stays
Grove's, including later improvements to it. A skill's instructions include
the guides and templates it uses, so `grove instructions propose "x"` shows
the proposal, change-spec and tasks templates to start from. Grove fills in
a template's `{{change}}` (the change's ID), `{{title}}`, `{{created}}`,
`{{schema}}` and `{{base}}` (the target's commit) before printing it; put
them back where you copy one.

### A schema of your own

A project adds a schema in `grove/schemas/<name>/`, with a `schema.yaml`
that names only what differs from `default`:

| Field | Meaning |
| --- | --- |
| `description` | One line saying which work the schema is for, shown when the agent offers the project's schemas. |
| `documents` | The kinds of document, as a list of `kind`, `path`, `template`, `guide` and, for the main document, `metadata`. Given, it replaces `default`'s list whole. |
| `skills` | For a skill, such as `propose`, its `guide` and the kinds it `writes`. A field or skill left out is `default`'s. |

Each file the schema needs comes from `grove/schemas/<name>/`, else from
`grove/schemas/default/`, else from Grove's `default`. A file of your own
for `default` is so your other schemas' too, unless they have their own.

A change has a schema, which its `.grove.yaml` names as `schema`, and keeps
it: revising, applying and archiving the change use it. A change without
one has `default`. A new change has `grove.yaml`'s `schema`, unless the owner
says which kind of work it is while exploring or proposing: when the project
has schemas other than `default`, the instructions list them, and the agent
asks for `grove instructions propose --schema <name>`. To change a change's
schema, edit its `.grove.yaml`. Schemas are shared by copying their folder.

`grove check` reports a mistake in a schema: an unknown field in its
`schema.yaml`, a kind a skill writes that the schema does not list, a file
it names that is in none of the three places, a file in its folder that it
does not use (often a misspelled path), `tasks` or `implementation` missing
or moved, and a `schema` in `grove.yaml` or a change's `.grove.yaml` that
names no schema. A skill's instructions are not built from a schema with
any of these but the unknown field and the unused file.

### Example: a bug schema

This schema is for fixing behavior that is wrong. It keeps `default`'s kinds
and each skill's `writes`, and replaces the propose guide, the proposal's
template and its guide; everything else is `default`'s.

`grove/schemas/bugfix/schema.yaml`:

```yaml
description: A fix for behavior that is wrong, whether or not a maintained spec describes it.
```

`grove/schemas/bugfix/skills/propose.md`:

```markdown
**Propose a bug fix: "This is wrong; make it right."** It turns a bug into a
change the owner can review: a proposal that shows the bug, tasks that start
with a test reproducing it, and a change spec only when the owner agrees the
corrected behavior. It describes future work. It does not plan or perform
the implementation.

### 1. Reproduce it

- Read the code and the maintained specs in the project index that the bug
  touches.
- Reproduce it, and keep the exact command and what it printed. If you
  cannot, ask the owner how; do not propose a fix for a bug nobody has seen.

### 2. Find what says what should happen

Tell the owner which of these it is:

- A maintained requirement says what should happen, and the code breaks it.
  Cite the requirement in the proposal. Write no change spec: what is agreed
  does not change.
- No maintained requirement covers the behavior. Ask the owner whether the
  corrected behavior should become agreed. If yes, write a change spec that
  ADDs it; archiving folds it into the maintained spec. If no, write none,
  and the fix stays observed behavior.
- The maintained spec itself is wrong. That is not a bug: tell the owner it
  changes agreed behavior and belongs to proposing with the `default`
  schema.

### 3. Write the change

- The first task adds a test that reproduces the bug and fails before the
  fix. A simple fix is that one task.
- No implementation planning: no file-by-file plan, no sequencing, no
  estimates.
```

`grove/schemas/bugfix/documents/proposal/template.md`:

```markdown
# <title>

<!-- A bug fix describes future work. It becomes approved future work when
it is merged into the target branch, which says nothing about
implementation. -->

## Observed

<!-- What happens today: the command, what it printed, and the source
(file path and revision). -->

## Expected

<!-- What should happen instead, and who says so. -->

## Reproduction

<!-- Steps someone else can follow from a clean checkout. -->

## The requirement it breaks

<!-- The maintained requirement, by spec and name; or that none covers this
behavior, and what the owner decided about agreeing one. -->

## Open questions

<!-- Say which must be answered before review and which can wait. Name the
task a question blocks when it decides how that task is built. -->

## Related documents

- Tasks: [tasks.md](tasks.md)
- Requirements: [the change spec, if the owner agreed one]
```

`grove/schemas/bugfix/documents/proposal/guide.md`:

```markdown
`proposal.md` shows the bug: what is observed, what is expected, how to
reproduce it, and the requirement it breaks.

- Title: write it in `.grove.yaml` and repeat it as the proposal's first
  heading. Name the corrected behavior, such as "Editing a missing record
  prints an error".
- Observed: cite each observation by its source instead of restating it as
  fact.
- Expected: a maintained requirement, cited, or the owner's words, quoted.
- The requirement it breaks: the one cited, or that none covers the
  behavior and whether the owner agreed to one.
- Open questions: name the task a question blocks when it decides how that
  task is built.
```

## Files

```text
grove.yaml
grove/
  schemas/<name>/
    schema.yaml
    documents/<kind>/template.md, guide.md
    skills/<skill>.md
  specs/<capability>/
    .grove.yaml
    spec.md
  changes/<id>-<slug>/
    .grove.yaml
    proposal.md
    specs/<capability>/spec.md
    tasks.md
    implementation.md
  changes/archive/<id>-<slug>/
```

A folder under `changes/` or `specs/` is a change or a maintained spec when
it holds a `.grove.yaml`; `grove check` reports one without its main
document: `spec.md`, or the change's schema's, `proposal.md` in `default`.
A maintained spec's `.grove.yaml` holds its `title`; a change's holds its
`title`, when it was `created`, and its `schema`. Either may add a one-line
`description`, shown in lists. The folder name gives the ID: a spec's is the
whole name, a change's the start of it. A main document starts with the
title as its heading, and `grove check` reports one whose first heading is
not `# ` and the title. A change's other files have no metadata: their
folder names their change. The change folder above is `default`'s;
`schemas/` holds the project's own change schemas (see "Change schemas").
