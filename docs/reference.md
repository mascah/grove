# Reference

## Terminal commands

| Command | Does |
| --- | --- |
| `grove` | Opens the reader on the maintained specs and approved changes. |
| `grove init` | Writes `grove.yaml` if there is none, and Grove's skills or commands for each harness it lists. Safe to re-run: see "Harnesses". Commits nothing. ([init spec](https://github.com/mascah/grove/blob/main/grove/specs/init/spec.md)) |
| `grove specs` | Lists the maintained specs on the target branch. |
| `grove changes` | Lists the approved changes on the target branch that are not archived, what each waits on, and the cycle not done that holds it, with `unattended` when it is marked and each cycle not done it follows. ([changes spec](https://github.com/mascah/grove/blob/main/grove/specs/changes/spec.md), [cycles spec](https://github.com/mascah/grove/blob/main/grove/specs/cycles/spec.md)) |
| `grove changes --json` | Prints the target branch's changes, archived ones included, as JSON (see [Coordinating changes](coordinating.md), "The project as data"). ([changes spec](https://github.com/mascah/grove/blob/main/grove/specs/changes/spec.md)) |
| `grove check` | Reports each problem in this working tree's Grove documents, in `grove.yaml` and in the project's change schemas, as `path:line: message`, and exits non-zero if there is one. Among them: a requirement with no "shall" or no scenario ([requirements spec](https://github.com/mascah/grove/blob/main/grove/specs/requirements/spec.md)); a `depends_on` that can never be met, and a requirement or a file that changed under a change since its `base` ([propose spec](https://github.com/mascah/grove/blob/main/grove/specs/propose/spec.md); [Coordinating changes](coordinating.md)); and a cycle that is malformed or that `grove run` would skip ([cycles spec](https://github.com/mascah/grove/blob/main/grove/specs/cycles/spec.md)). |
| `grove review` | Shows what merging the current branch into the target would approve: each change it brings, as future work or with its code, each Grove document it adds, changes or removes, and the diff. While the branch holds a change open, or archived where the target does not, it lists before the documents each file the branch changes outside the document root: `planned <change-id>:T-NN` when a task's `Files:` names it or a folder above it, `agreed <change-id>:T-NN: <reason>` when a task's `Also changed:` line in `implementation.md` names it, or `named nowhere`. ([implement spec](https://github.com/mascah/grove/blob/main/grove/specs/implement/spec.md), "Review names a change approved with its code", "Review shows planned against actual") |
| `grove archive <change-id>` | Moves the change's folder to `changes/archive/`, once every task in `tasks.md` is `Status: done` in `implementation.md`, or the file says `Status: dropped`. For a finished change it first checks that the maintained specs hold the change's requirements word for word. Commits nothing. ([changes spec](https://github.com/mascah/grove/blob/main/grove/specs/changes/spec.md), [requirements spec](https://github.com/mascah/grove/blob/main/grove/specs/requirements/spec.md)) |
| `grove run` | In a clean checkout, takes the oldest cycle on the target that is marked `unattended` and can be built, and starts one Claude Code session that builds it on a new branch. Then it starts nothing until that cycle is archived on the target. See [Running unattended](running.md). ([runs spec](https://github.com/mascah/grove/blob/main/grove/specs/runs/spec.md)) |
| `grove run <cycle-id>` | As `grove run`, but builds the cycle given, marked or not, and says why when it cannot. See [Running unattended](running.md), "Run a cycle you name". |
| `grove run --retry` | Starts a new attempt even though the latest one's cycle is not archived. |
| `grove run --timeout <duration>` | Stops the session after `<duration>`, such as `30m` or `1h30m`, rather than 2 hours. Goes with `--retry`. |
| `grove version` | Prints the version of Grove, such as `grove 0.1.0`. A build from a checkout prints a Go pseudo-version naming its commit, and `dev` when Go recorded none. |

Skills run these; you do not need to:

| Command | Prints the instructions to |
| --- | --- |
| `grove instructions explore` | Interview the user about an idea, with the project's specs and changes listed. |
| `grove instructions prototype "<question>"` | Answer the question with a throwaway experiment on a `prototype/<slug>` branch that is never merged, writing no Grove document ([prototype spec](https://github.com/mascah/grove/blob/main/grove/specs/prototype/spec.md)). |
| `grove instructions propose [--schema <name>] ["<description>"]` | Write a new change from the conversation or the description, after one round of questions, with the project's change schema or the one named; check it, then offer to merge it now or build it on this branch first ([propose spec](https://github.com/mascah/grove/blob/main/grove/specs/propose/spec.md)). |
| `grove instructions propose <change-id>` | Revise a drafted change; for a change in a cycle not done, re-plan that cycle on the same branch too. |
| `grove instructions plan ["<description>"]` | Write a new cycle of approved changes, asking which need a person, editing no change (see [Planning cycles](planning.md); [cycles spec](https://github.com/mascah/grove/blob/main/grove/specs/cycles/spec.md)). |
| `grove instructions plan <change-id>` | Write a new cycle whose only change is that approved change. |
| `grove instructions plan <cycle-id>` | Revise the cycle as the working tree holds it, setting `planned` to the target's commit. |
| `grove instructions implement <change-id>` | Build a change with the owner, asking them what it leaves open and, with their agreement, revising it on the branch; record progress, and archive it when every task is done. The change may be approved, or a draft this working tree holds. Refuses a change that waits on one not archived as finished ([implement spec](https://github.com/mascah/grove/blob/main/grove/specs/implement/spec.md)). |
| `grove instructions implement <change-id> --unattended` | Build an approved change as above but ask no one: what it cannot settle is marked blocked, and the change itself is not edited. |
| `grove instructions implement <cycle-id>` | Build each change of the cycle this working tree has not archived with the owner, in build order, one at a time, each as `grove instructions implement <change-id>` does. |
| `grove instructions implement <cycle-id> --unattended` | Build each change of the cycle the same way without asking anyone, marked `unattended` or not (see [Coordinating changes](coordinating.md), "Building a cycle unattended"). |
| `grove instructions archive <change-id>` | Fold the change's specs into the maintained specs word for word, adding the change's line to each one's History, then run `grove archive`, fixing the specs and running it again if it refuses ([requirements spec](https://github.com/mascah/grove/blob/main/grove/specs/requirements/spec.md)). |

The instructions of `implement`, `archive` and `propose <change-id>` name
each of the change's documents by path, for the agent to read with its own
file tool, and print none of their text, so they do not grow with the
change. For an approved change, `implement` also names each document whose
copy in the working tree differs from the target's, with the
`git show <commit>:<path>` command that prints the approved text.

## Harnesses

Grove has six skills: explore, prototype, propose, plan, implement and
archive. `grove init` writes a file for each of them, for each harness
under `agents` in `grove.yaml`. It writes them as skills, or as commands
where the harness has commands. It also removes the files of skills Grove
no longer has, such as `grove-apply`.

How you type implement, and the other skills alike:

| Harness | `agents` key | As skills | As commands |
| --- | --- | --- | --- |
| Claude Code | `claude` | `/grove-implement` | `/grove:implement` |
| Codex | `codex` | `$grove-implement` | none |
| Cursor | `cursor` | `/grove-implement` | none |
| OpenCode | `opencode` | not typed; ask for `grove-implement` by name | `/grove-implement` |
| Pi | `pi` | `/skill:grove-implement` | `/grove-implement` |
| Hermes | `hermes` | `/grove-implement` | none |

Type the words a skill takes after its name, such as
`/grove-implement 261020-m4k8t --unattended`.

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

The first `grove init` runs where there is no `grove.yaml`. It chooses the
harnesses from what the repository holds, each as skills:

- `.claude/` or `CLAUDE.md` gives `claude`.
- `.codex/` or `AGENTS.md` gives `codex`.
- `.cursor/` gives `cursor`.
- `.opencode/` or `opencode.json` gives `opencode`.
- `.pi/` gives `pi`.
- `.hermes/` gives `hermes`.

With none of these, it chooses `claude` and `codex`. It writes them under
`agents`, and says which it chose.

Later runs follow `agents` and nothing else. Each run writes every file of
the harnesses listed, overwriting edits to them. It also removes Grove's
files that `agents` no longer gives: those of a harness taken out, of a
form changed, or of a skill Grove no longer has. It looks only in the
places in the table above. There it removes only `grove-` folders and
files, and files in `.claude/commands/grove/`. Its message names each one
it removed. The
[init spec](https://github.com/mascah/grove/blob/main/grove/specs/init/spec.md) has the exact rules.

## grove.yaml

```yaml
root: grove
target: main
agents:
  claude: commands
  codex: skills
context: |
  Run go test ./... before saying a task is done.
  Every public API stays backwards compatible.
rules:
  documents:
    proposal:
      - List who is affected.
    tasks:
      - Every task names its test.
  skills:
    implement:
      - Commit after each task, naming its task ID.
```

| Field | Default | Meaning |
| --- | --- | --- |
| `root` | `grove` | The folder that holds `specs/`, `changes/` and `cycles/`. |
| `target` | `main` | The branch whose contents count as agreed. Merging into it approves. |
| `agents` | `claude` and `codex`, as `skills` | Each harness `init` writes Grove's skills for, and whether as `skills` or `commands` (see "Harnesses"). The first `init` chooses them. |
| `schema` | `default` | The change schema a new change has unless the owner names another while proposing. See "Change schemas". |
| `context` | empty | Instructions given to every skill, so they reach agents working through Grove and no other session; project facts every session needs stay in `CLAUDE.md` or `AGENTS.md`. |
| `rules` | none | Extra lines for Grove's instructions. Under `documents`, for one kind of document, given when that document is written: `proposal`, `change-spec`, `tasks`, `implementation`, `spec`, and those of the project's own schemas. Under `skills`, for one skill, given with its instructions: `explore`, `prototype`, `propose`, `plan`, `implement`, `archive`; lines under `implement` are given with the owner and with `--unattended`. |

`grove check` reports a field or a `rules` key that Grove does not know.
Other commands ignore it. A `rules.skills.apply` key is from before apply
became `implement --unattended`; `grove check` says its lines go under
`rules.skills.implement`. `grove check` also reports a harness under
`agents` that Grove does not know, or a form that harness does not have.
`init` stops at either.

`context` and `rules` add to Grove's own instructions; they cannot remove
or replace them. Short rules work better than long ones. When you correct
`/grove-implement` in a way that would apply to future work, it may suggest
one line for them when it finishes, and adds it only if you agree.
`--unattended` only names a line in its report. The
[config spec](https://github.com/mascah/grove/blob/main/grove/specs/config/spec.md) has the exact rules for `context`.

Grove's skills commit nothing unless you ask in the conversation or a rule
asks, as `implement` does above.

## Change schemas

A change schema says what a change consists of and how each skill goes.
It lists kinds of document, each with a path, a template and a guide. For
each skill, it gives the skill's guide and the kinds of document it
writes. Grove has two built-in
schemas ([schemas spec](https://github.com/mascah/grove/blob/main/grove/specs/schemas/spec.md)): `bugfix`, for behavior that is wrong (see "Grove's bugfix"), and
`default`, the software change [Concepts](concepts.md) describes:

| Kind | Path, under the root | Written by |
| --- | --- | --- |
| `proposal` | `changes/<change>/proposal.md`, the change's main document | propose, implement |
| `change-spec` | `changes/<change>/specs/<capability>/spec.md` | propose, implement |
| `tasks` | `changes/<change>/tasks.md` | propose, implement |
| `implementation` | `changes/<change>/implementation.md` | implement, archive |
| `spec` | `specs/<capability>/spec.md`, a maintained spec | archive |

A schema's folder holds `schema.yaml`, a folder for each kind of document,
`documents/<kind>/` with its `template.md` and `guide.md`, and the guide for
each skill, `skills/<skill>.md`.

Some things a schema cannot change:

- Grove's skills, and what Grove requires of each. A skill's instructions
  give those first, such as propose asking one round of every decision it
  would otherwise assume before it writes.
- That a merge approves a change.
- A change's `.grove.yaml`.
- How progress is recorded. Each task is a `T-NN` heading in `tasks.md`.
  It names the files it touches under `Files:`, one backticked path per
  bullet or, as before, separated by commas on one line; a path ending in
  `/` names a folder and every file under it. It names the interfaces it
  produces on an `Interfaces:` line. An optional `Delegated:` line names
  choices and paths the implementer may take without asking the owner;
  Grove does not read it. A task's status is a `Status:` line in
  `implementation.md`, which `grove archive` reads.

So every schema has `tasks` and `implementation` at those paths, and keeps
the `Status:` lines of Grove's implementation template. Every schema has
one kind with `metadata: changes/<change>/.grove.yaml`: the change's main
document. No other kind has `metadata`. Every other kind is in the change
folder, except `spec`.

To add lines to Grove's instructions rather than replace a file, use `rules`
in `grove.yaml`.

### Your own files for default

To change one of `default`'s files, put your own at the same path under
`grove/schemas/default/`, such as `grove/schemas/default/documents/tasks/template.md`.
It replaces Grove's file whole. Every file you do not put there stays
Grove's, including later improvements to it. A skill's instructions include
the guides and templates it uses, so `grove instructions propose "x"` shows
the proposal, change-spec and tasks templates to start from.

Before printing a template, Grove fills in these values. Put them back
where you copy one.

- `{{change}}`: the change's ID.
- `{{title}}`: its title, once it has one.
- `{{created}}`: when it was created.
- `{{schema}}`: its schema.
- `{{base}}`: the target's commit, shortened to seven characters.
- `{{commit}}`: the target's full commit.

### A schema of your own

A project adds a schema in `grove/schemas/<name>/`, with a `schema.yaml`
that names only what differs from `default`, such as a schema for changes
too small for change specs:

```yaml
description: A small change with no change specs.
skills:
  propose:
    writes: [proposal, tasks]
```

| Field | Meaning |
| --- | --- |
| `description` | One line saying which work the schema is for, shown when the agent offers the project's schemas. |
| `documents` | The kinds of document, as a list of `kind`, `path`, `template`, `guide` and, for the main document, `metadata`. Given, it replaces `default`'s list whole. |
| `skills` | For a skill, such as `propose`, its `guide` and the kinds it `writes`. A field or skill left out is `default`'s. |

Each file the schema needs comes from `grove/schemas/<name>/`, else from
`grove/schemas/default/`, else from Grove's `default`. So a file of your
own for `default` is also your other schemas' file, unless they have
their own.

A change has a schema, which its `.grove.yaml` names as `schema`. It keeps
that schema: revising, implementing and archiving the change use it. A
change without one has `default`. To change a change's schema, edit its
`.grove.yaml`.

A new change has `grove.yaml`'s `schema`, unless the owner says which kind
of work it is while exploring or proposing. For that, the instructions
list the project's schemas, Grove's two first. The agent suggests
`bugfix` when the owner describes behavior that is wrong. Once the owner
confirms a schema, the agent asks for `grove instructions propose
--schema <name>`.

Schemas are shared by copying their folder.

`grove check` reports these mistakes in a schema:

- an unknown field in its `schema.yaml`;
- a kind a skill writes that the schema does not list;
- a file it names that is in none of its places;
- a file in its folder that it does not use, often a misspelled path;
- `tasks` or `implementation` missing or moved;
- a `schema` in `grove.yaml` or a change's `.grove.yaml` that names no
  schema.

Grove builds no skill's instructions from a schema with one of these
mistakes, except the unknown field and the unused file.

### Grove's bugfix

`bugfix` is for fixing behavior that is wrong. Every project has it, with
no file to write. It keeps `default`'s kinds and each skill's `writes`, and
has its own description, propose guide, and proposal template and guide:

- **Proposal.** Its sections are Observed, Expected, Reproduction, The
  requirement it breaks and Open questions.
- **Propose guide.** Reproduce the bug first, keeping the command and what
  it printed. If a maintained requirement says what should happen, cite it
  and write no change spec: what is agreed does not change. If none covers
  the behavior, ask in the round whether the corrected behavior should
  become agreed, and write a change spec that ADDs it only if the owner
  says yes. If the maintained spec itself is wrong, it is not a bug but a
  `default` change. The first task adds a test that reproduces the bug and
  fails before the fix.

`grove instructions propose --schema bugfix "x"` shows its files as a
session gets them. Each file comes from the first of these places that
has it, and each field of its `schema.yaml` from the first that gives it:

1. `grove/schemas/bugfix/`
2. Grove's `bugfix`
3. `grove/schemas/default/`
4. Grove's `default`

So to replace one of its files, put yours at the same path under
`grove/schemas/bugfix/`, such as
`grove/schemas/bugfix/documents/proposal/template.md`. The rest stay
Grove's. A file of your own for `default`, such as a tasks template,
reaches `bugfix` too, unless `grove/schemas/bugfix/` has its own.

## Files

```text
grove.yaml
grove/
  schemas/<name>/
    schema.yaml
    documents/<kind>/template.md, guide.md
    skills/<skill>.md
  specs/<capability>/
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

A folder under `specs/` is a maintained spec when it holds `spec.md`
([maintained-specs spec](https://github.com/mascah/grove/blob/main/grove/specs/maintained-specs/spec.md)).
Its ID is the folder's name. Its title is the file's first `# ` heading;
`grove check` reports a spec with none. Its description, in each skill's
project index, is the first sentence of its `## Purpose`. A spec has no
`.grove.yaml`. `grove check` reports one left in a spec folder as no
longer read, until it is deleted.

A folder under `changes/` is a change when it holds a `.grove.yaml`. Its
ID is the start of the folder's name. `grove check` reports a change
without its main document: its schema's, which is `proposal.md` in
`default`. The main document starts with the change's title as its
heading, and `grove check` reports one whose first heading is not `# `
and the title. A change's other files have no metadata: their folder
names their change.

A change's `.grove.yaml` holds:

- `title`;
- `created`, when it was created;
- `schema`;
- `base`, the full commit of the target it was proposed from;
- optionally `description`, one line shown in lists;
- optionally `depends_on`, the IDs of the changes it waits on (see
  [Coordinating changes](coordinating.md)).

Propose writes it from Grove's template.

A cycle file's name starts with its ID. It holds:

- `title`;
- `created`, when it was created;
- optionally `goal`, one line;
- `changes`, the IDs of the approved changes in it;
- optionally `depends_on`, the IDs of the cycles it follows;
- `planned`, the full commit of the target its plan read;
- optionally `needs_human`, each change the owner said needs a person,
  with their reason;
- optionally `unattended: true`, which hands it to `grove run` (see
  [Planning cycles](planning.md)).

The change folder above is `default`'s. `schemas/` holds the project's own
change schemas (see "Change schemas").
