# Reference

## Terminal commands

| Command | Does |
| --- | --- |
| `grove` | Opens the reader on the maintained specs and approved changes. |
| `grove init` | Writes `grove.yaml` if there is none, and Grove's skills or commands for each harness it lists. Safe to re-run: see "Harnesses". Commits nothing. |
| `grove specs` | Lists the maintained specs on the target branch. |
| `grove changes` | Lists the approved changes on the target branch that are not archived, what each waits on, and the cycle not done that holds it, with `unattended` when it is marked. |
| `grove changes --json` | Prints the target branch's changes, archived ones included, as JSON (see [Coordinating changes](coordinating.md), "The project as data"). |
| `grove check` | Reports each problem in the Grove documents of this working tree, each unknown field or `rules` key in `grove.yaml`, and each mistake in a change schema, as `path:line: message`, and exits non-zero if there is one. Among the problems: a `depends_on` naming no change, open or archived, or a dropped change, changes that wait on each other in a loop, the cycle problems [Planning cycles](planning.md) lists, a requirement a change modifies or removes whose text in the working tree's maintained spec is neither what it was at the change's `base` nor what the change or a change it depends on gives it, a task's `Files:` line it cannot read, and, for a change in no cycle that is not done, a file a task names that changed on the target since the change's `base` ([Coordinating changes](coordinating.md), "When the plan is stale"). |
| `grove review` | Shows what merging the current branch into the target would approve: each change, maintained spec, cycle or change schema file it adds or alters, and the diff. |
| `grove archive <change-id>` | Moves the change's folder to `changes/archive/`. Refuses unless `implementation.md` marks every task in `tasks.md` `Status: done`, or says `Status: dropped`. Commits nothing. |
| `grove run` | In a clean checkout, fetches the target from `origin` if there is one, else uses the local target, selects the cycle there with the earliest `created` that is marked `unattended` and can be built unattended, and starts one Claude Code session on the apply skill for it on a new branch, `grove/run-<attempt>`. Keeps a record and the output in `.git/grove/`, then starts nothing until every change of that cycle is archived on the target. See [Running apply unattended](running.md). |
| `grove run <cycle-id>` | As `grove run`, but builds the cycle given, whether or not it is marked `unattended`, and says why when it cannot be built. See [Running apply unattended](running.md), "Run a cycle you name". |
| `grove run --retry` | Starts a new attempt even though the latest one's cycle is not archived. |
| `grove run --timeout <duration>` | Stops the session after `<duration>`, such as `30m` or `1h30m`, rather than 2 hours. Goes with `--retry`. |
| `grove version` | Prints the version of Grove, such as `grove 0.1.0`. A build from a checkout prints a Go pseudo-version naming its commit, and `dev` when Go recorded none. |

Skills run these; you do not need to:

| Command | Prints the instructions to |
| --- | --- |
| `grove instructions explore` | Interview the user about an idea, with the project's specs and changes listed. |
| `grove instructions prototype "<question>"` | Answer the question, or one from the conversation, with a throwaway experiment on a `prototype/<slug>` branch that is never merged, writing no Grove document; a proposal cites its findings. |
| `grove instructions propose [--schema <name>] ["<description>"]` | Write a new change from the conversation or the description, with the project's change schema or the one named. |
| `grove instructions propose <change-id>` | Revise a drafted change, with the schema it has. |
| `grove instructions plan ["<description>"]` | Write a new cycle of approved changes from the conversation or the description, asking which need a person, editing no change (see [Planning cycles](planning.md)). |
| `grove instructions plan <change-id>` | Write a new cycle whose only change is that approved change. |
| `grove instructions plan <cycle-id>` | Revise the cycle as the working tree holds it, setting `planned` to the target's commit. |
| `grove instructions implement <change-id>` | Build an approved change with the owner, asking them what the change leaves open and, with their agreement, revising it on the branch; record progress, and archive it when every task is done. Refuses a change that waits on one this working tree does not hold archived as finished. |
| `grove instructions implement <cycle-id>` | Build each change of the cycle this working tree has not archived with the owner, in build order, one at a time, each as `grove instructions implement <change-id>` does. |
| `grove instructions apply <cycle-id>` | Build each change of the cycle this working tree has not archived without asking anyone, in build order, one at a time, each with its own instructions (see [Coordinating changes](coordinating.md), "Applying a cycle"). Refuses a change's ID, and no ID. |
| `grove instructions apply <cycle-id> <change-id>` | Build that change of the cycle without asking anyone, marking a task an open question decides blocked, record progress, and archive it when every task is done; with the cycle's changes in build order. Refuses a change the cycle does not list, and one that waits on a change this working tree does not hold archived as finished. |
| `grove instructions archive <change-id>` | Fold the change's specs into the maintained specs, then run `grove archive`. |

## Harnesses

`grove init` writes a file for each of Grove's seven skills, explore,
prototype, propose, plan, implement, apply and archive, for each harness under `agents` in
`grove.yaml`: as a skill, or as a command where the harness has commands.
How you type apply, the others alike:

| Harness | `agents` key | As skills | As commands |
| --- | --- | --- | --- |
| Claude Code | `claude` | `/grove-apply` | `/grove:apply` |
| Codex | `codex` | `$grove-apply` | none |
| Cursor | `cursor` | `/grove-apply` | none |
| OpenCode | `opencode` | not typed; ask for `grove-apply` by name | `/grove-apply` |
| Pi | `pi` | `/skill:grove-apply` | `/grove-apply` |
| Hermes | `hermes` | `/grove-apply` | none |

Type the words a skill takes after its name, such as
`/grove-apply 261020-m4k8t`.

Where the files go, `<skill>` being the skill's name:

| Harness | Skills | Commands |
| --- | --- | --- |
| `claude` | `.claude/skills/grove-<skill>/SKILL.md` | `.claude/commands/grove/<skill>.md` |
| `codex`, `cursor`, `hermes` | `.agents/skills/grove-<skill>/SKILL.md` | none |
| `opencode` | `.agents/skills/grove-<skill>/SKILL.md` | `.opencode/commands/grove-<skill>.md` |
| `pi` | `.agents/skills/grove-<skill>/SKILL.md` | `.pi/prompts/grove-<skill>.md` |

Harnesses that share `.agents/skills/` share one copy of each skill. Commit
the files with your project.

Two harnesses read a project's files only once you trust it: in Hermes, run
`hermes skills trust` in the repository; in Pi, trust the project when it
asks, or with `/trust`.

The first `grove init`, with no `grove.yaml`, chooses the harnesses from
what the repository holds, each as skills: `.claude/` or `CLAUDE.md` gives
`claude`, `.codex/` or `AGENTS.md` gives `codex`, `.cursor/` gives
`cursor`, `.opencode/` or `opencode.json` gives `opencode`, `.pi/` gives
`pi`, and `.hermes/` gives `hermes`. With none of these, it chooses `claude`
and `codex`. It writes them under `agents`, and says which it chose.

Later runs follow `agents` and nothing else. Each run writes every file the
harnesses listed have, overwriting edits to them, and removes Grove's files
that `agents` no longer gives: those of a harness taken out, of a form
changed, or of a skill Grove no longer has. It looks only in the places in
the table above, and removes only `grove-` folders and files, and files in
`.claude/commands/grove/`. Its message names each one it removed.

## grove.yaml

```yaml
root: grove
target: main
agents:
  claude: commands
  codex: skills
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
| `root` | `grove` | The folder that holds `specs/`, `changes/` and `cycles/`. |
| `target` | `main` | The branch whose contents count as agreed. Merging into it approves. |
| `agents` | `claude` and `codex`, as `skills` | Each harness `init` writes Grove's skills for, and whether as `skills` or `commands` (see "Harnesses"). The first `init` chooses them. |
| `schema` | `default` | The change schema a new change has unless the owner names another while proposing. See "Change schemas". |
| `rules` | none | Extra lines for Grove's instructions. Under `documents`, for one kind of document, given when that document is written: `proposal`, `change-spec`, `tasks`, `implementation`, `spec`, and those of the project's own schemas. Under `skills`, for one skill, given with its instructions: `explore`, `prototype`, `propose`, `plan`, `implement`, `apply`, `archive`. |

`grove check` reports a field or `rules` key that Grove does not know; other
commands ignore it. It also reports a harness under `agents` that Grove has
no row for, or a form that harness does not have, and `init` stops at
either.

`rules` adds to Grove's own instructions; it cannot remove or
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
| `proposal` | `changes/<change>/proposal.md`, the change's main document | propose, implement |
| `change-spec` | `changes/<change>/specs/<capability>/spec.md` | propose, implement |
| `tasks` | `changes/<change>/tasks.md` | propose, implement |
| `implementation` | `changes/<change>/implementation.md` | implement, apply, archive |
| `spec` | `specs/<capability>/spec.md`, a maintained spec | archive |

A schema's folder holds `schema.yaml`, a folder for each kind of document,
`documents/<kind>/` with its `template.md` and `guide.md`, and the guide for
each skill, `skills/<skill>.md`.

A schema cannot change Grove's skills or what Grove requires of each, which
its instructions give first, such as propose asking one round of every
decision it would otherwise assume before it writes; that a merge approves a change; a change's
`.grove.yaml`; or how progress is recorded: each task is a `T-NN` heading in
`tasks.md`, naming the files it touches on a `Files:` line and the
interfaces it produces on an `Interfaces:` line, and its status a `Status:`
line in `implementation.md`, which `grove archive` reads. So every schema has `tasks` and `implementation` at
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
it: revising, implementing, applying and archiving the change use it. A change without
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
- No maintained requirement covers the behavior. Ask, in the round before
  writing, whether the corrected behavior should become agreed. If yes, write a change spec that
  ADDs it; archiving folds it into the maintained spec. If no, write none,
  and the fix stays observed behavior.
- The maintained spec itself is wrong. That is not a bug: tell the owner it
  changes agreed behavior and belongs to proposing with the `default`
  schema.

### 3. Write the change

- The first task adds a test that reproduces the bug and fails before the
  fix. A simple fix is that one task: its `Files:` line names the test and
  the file under fix, and its `Interfaces:` line what the fix changes that
  callers see.
- No code, commands or steps, no sequencing, no estimates.
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
  cycles/<id>-<slug>.yaml
```

A folder under `changes/` or `specs/` is a change or a maintained spec when
it holds a `.grove.yaml`; `grove check` reports one without its main
document: `spec.md`, or the change's schema's, `proposal.md` in `default`.
A maintained spec's `.grove.yaml` holds its `title`; a change's holds its
`title`, when it was `created`, its `schema`, and `base`, the full commit of
the target it was proposed from, which propose writes from Grove's template.
Either may add a one-line `description`, shown in lists. A change may add
`depends_on`, the list of change IDs it waits on (see
[Coordinating changes](coordinating.md)). The folder name gives the ID: a spec's is the
whole name, a change's the start of it. A main document starts with the
title as its heading, and `grove check` reports one whose first heading is
not `# ` and the title. A change's other files have no metadata: their
folder names their change. A cycle file's name starts with its ID, and it
holds its `title`, when it was `created`, an optional one-line `goal`, `changes`, the list of
approved change IDs in it, `planned`, the full commit of the target its
plan read, an optional `needs_human`, each change the owner said needs a
person with their reason, and an optional `unattended: true`, which hands
it to `grove run` (see [Planning cycles](planning.md)). The change folder above is `default`'s;
`schemas/` holds the project's own change schemas (see "Change schemas").
