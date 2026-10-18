# Reference

## Terminal commands

| Command | Does |
| --- | --- |
| `grove` | Opens the reader on the maintained specs and approved changes. |
| `grove init` | Writes `grove.yaml` if there is none, and the skills for each harness it lists. Safe to re-run: it rewrites the skills and leaves `grove.yaml` alone. Commits nothing. |
| `grove specs` | Lists the maintained specs on the target branch. |
| `grove changes` | Lists the approved changes on the target branch that are not archived. |
| `grove check` | Reports each problem in the Grove documents of this working tree as `path:line: message`, and exits non-zero if there is one. |
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

`grove init` writes one skill per workflow for each harness in `grove.yaml`.

| Workflow | Claude Code | Codex |
| --- | --- | --- |
| explore | `/grove-explore` | `$grove-explore` |
| propose | `/grove-propose` | `$grove-propose` |
| apply | `/grove-apply` | `$grove-apply` |
| archive | `/grove-archive` | `$grove-archive` |

Claude Code reads them from `.claude/skills/grove-<workflow>/SKILL.md`, Codex
from `.agents/skills/grove-<workflow>/SKILL.md`. Commit them with your
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
| `rules` | none | Extra lines for one kind of document, given when that document is written, or for one workflow, given with its instructions. Kinds: `proposal`, `change-spec`, `tasks`, `implementation`, `spec`. Workflows: `explore`, `propose`, `apply`, `archive`. |

`context` and `rules` add to Grove's own instructions; they cannot remove or
replace them. Short rules work better than long ones.

Grove's skills commit nothing unless you ask in the conversation or a rule
asks, as `apply` does above.

## Files

```text
grove.yaml
grove/
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

A maintained spec's `.grove.yaml` holds its `title`; a change's holds its
`title` and when it was `created`. Either may add a one-line `description`,
shown in lists. The folder name gives the ID: a spec's is the whole name, a
change's the start of it. `spec.md` and `proposal.md` start with the title as
their heading. A change's other files have no metadata: their folder names
their change.
