# Reference

## Terminal commands

| Command | Does |
| --- | --- |
| `grove` | Opens the reader on the maintained specs and approved changes. |
| `grove init` | Writes `grove.yaml` if there is none, and the skills for each harness it lists. Safe to re-run: it rewrites the skills and leaves `grove.yaml` alone. Commits nothing. |
| `grove specs` | Lists the maintained specs on the target branch. |
| `grove changes` | Lists the approved changes on the target branch that are not archived. |
| `grove check` | Reports each problem in the Grove documents of this working tree, and each unknown field or `rules` key in `grove.yaml`, as `path:line: message`, and exits non-zero if there is one. |
| `grove review` | Shows what merging the current branch into the target would approve: each change or spec it adds or alters, and the diff. |
| `grove archive <change-id>` | Moves the change's folder to `changes/archive/`. Refuses unless `implementation.md` marks every task in `tasks.md` `Status: done`, or says `Status: dropped`. Commits nothing. |

Skills run these; you do not need to:

| Command | Prints the instructions to |
| --- | --- |
| `grove instructions explore` | Interview the user about an idea, with the project's specs and changes listed. |
| `grove instructions propose ["<description>" \| <change-id>]` | Write a new change from the conversation or the description, or revise a drafted one. |
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
  proposal:
    - List who is affected.
  tasks:
    - Every task names its test.
  apply:
    - Commit after each task, naming its task ID.
```

| Field | Default | Meaning |
| --- | --- | --- |
| `root` | `grove` | The folder that holds `specs/` and `changes/`. |
| `target` | `main` | The branch whose contents count as agreed. Merging into it approves. |
| `agents` | `[claude, codex]` | The harnesses `init` writes skills for. |
| `context` | empty | Facts about the project, given to every skill. Keep it to what an agent cannot read from the code. |
| `rules` | none | Extra lines for one kind of document, given when that document is written, or for one skill, given with its instructions. Kinds: `proposal`, `change-spec`, `tasks`, `implementation`, `spec`. Skills: `explore`, `propose`, `apply`, `archive`. |

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
is `default`, the software change [Concepts](concepts.md) describes. Its
files are:

```text
templates/   proposal.md, change-spec.md, tasks.md, implementation.md, spec.md
guides/      the same five, one per kind, and explore.md, propose.md,
             apply.md, archive.md, one per skill
```

To change one of them for your project, put your own at the same path under
`grove/schemas/default/`, such as `grove/schemas/default/templates/tasks.md`.
It replaces Grove's file whole; every file you do not put there stays
Grove's, including later improvements to it. A skill's instructions include
the guides and templates it uses, so `grove instructions propose "x"` shows
the proposal, change-spec and tasks templates to start from. Grove fills in
a template's `{{change}}` (the change's ID), `{{title}}`, `{{created}}` and
`{{base}}` (the target's commit) before printing it; put them back where you
copy one.

A schema cannot change Grove's skills or what Grove requires of each, which
its instructions give first; that a merge approves a change; a change's
`.grove.yaml`; or how progress is recorded: each task is a `T-NN` heading in
`tasks.md`, and its status a `Status:` line in `implementation.md`, which
`grove archive` reads. Keep those lines as Grove's implementation template
has them.

To add lines to Grove's instructions rather than replace a file, use `rules`
in `grove.yaml`.

## Files

```text
grove.yaml
grove/
  schemas/default/
    templates/, guides/
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
it holds a `.grove.yaml`; `grove check` reports one without its
`proposal.md` or `spec.md`. A maintained spec's `.grove.yaml` holds its `title`; a change's holds its
`title` and when it was `created`. Either may add a one-line `description`,
shown in lists. The folder name gives the ID: a spec's is the whole name, a
change's the start of it. `spec.md` and `proposal.md` start with the title as
their heading, and `grove check` reports one whose first heading is not `# `
and the title. A change's other files have no metadata: their folder names
their change. `schemas/` holds the project's own files for change schemas
(see "Change schemas").
