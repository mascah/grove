# Grove

A local project workspace for humans and agents: work, questions, decisions,
plans, reviews and knowledge as Markdown files with YAML frontmatter in your
Git repository, read and changed through a CLI and a terminal board, with no
service to run. Records carry sequential IDs shared by every linked worktree;
Grove reads each record's versions across local branches and worktrees and
decides which are current by Git ancestry. Two workflows, shaping ideas into
proposed work and executing assigned work into a candidate for review, travel
inside the binary for agents to follow.

## Run it

Requires Go 1.26.8 or later (`go.mod`). From this repository:

```sh
go run ./cmd/grove              # the terminal board; needs a terminal
go run ./cmd/grove --help       # every command's usage
go run ./cmd/grove list
go run ./cmd/grove show G-003
go run ./cmd/grove check
```

## Commands

`grove --help` has each command's usage. Every command takes
`--project DIR`; without it, Grove searches upward for `grove.yaml`.

| Command | What it does | Details |
| --- | --- | --- |
| (none) | Opens the terminal board | [Board](docs/board.md) |
| `list`, `show`, `brief`, `check` | List, print and validate records and the brief, reading only | [Record model](docs/record-model.md#reading-and-writing-records) |
| `new`, `update` | Create a record with the next shared ID; change its frontmatter | [Record model](docs/record-model.md#identity-and-dates) |
| `convert` | Make a record from a Markdown document outside the record root | [Record model](docs/record-model.md#identity-and-placement-apart-from-classification) |
| `approve`, `feedback`, `integrate` | Judge a candidate in review; merge an approved one and mark it done | [Work lifecycle](docs/record-model.md#work-lifecycle) |
| `resolve` | Give a candidate that conflicts with the target to one attempt that merges the target and resolves it | [Resolving a conflict](docs/commands.md#resolving-a-conflict) |
| `sweep` | Resolve, approve and integrate candidates in review under the owner's standing `policy:`, each act attributed to it | [Sweep](docs/commands.md#sweep) |
| `run`, `attempts`, `attempt`, `stop` | Start one headless agent attempt of a work item that outlives the terminal; list, inspect and stop attempts | [Attempts](docs/commands.md#attempts) |
| `versions` | Show each record's versions across branches and worktrees, current or older | [Versions](docs/commands.md#versions) |
| `workspace` | Print the checkout holding a selected version | [Workspace](docs/commands.md#workspace) |
| `context` | Assemble staged context for selected work | [Context](docs/commands.md#context) |
| `deps` | Show how unfinished work depends on other work, or preview a selection's order | [Dependencies](docs/commands.md#dependencies) |
| `init` | Set up Grove in a Git checkout | [Init](docs/commands.md#init) |
| `guide`, `version` | Print a workflow or review guide or the record model; name this build | [Version and guide](docs/commands.md#version-and-guide) |

## The board

`grove` with no command opens a board of work by status in its current state
across every branch and checkout. From it you open a record's detail,
history and versions, search every record, judge and integrate a candidate,
and start and stop attempts. [The board](docs/board.md) describes each view
and key.

## Use Grove in another repository

Build one binary from a named commit, put it on `PATH`, and check which
`grove` answers before `init`:

```sh
go build -o "$HOME/.local/bin/grove" ./cmd/grove   # from a clone at that commit
GOBIN="$HOME/.local/bin" go install github.com/mascah/grove/cmd/grove@COMMIT  # or from the module; pushed commits only
grove version     # must print "grove … content sha256:…"; a usage error means another grove answered
grove init        # at the checkout's top, or: grove --project /absolute/path init
grove check
git add grove.yaml grove .claude .agents && git commit -m "chore: set up grove"   # at the top
```

Commit what `init` wrote before the first attempt: an attempt's worktree
holds only committed files, and `grove run` refuses one without the
`grove-work` skill. Then `/grove-shape TOPIC` develops the brief and proposes
work, and `/grove-work G-001` carries it out (`$grove-shape` and `$grove-work`
in Codex), or `grove run G-001 --budget USD --permission-mode MODE` starts it
as a headless attempt ([Attempts](docs/commands.md#attempts)).
[Init](docs/commands.md#init) says what `init` writes and how to point Codex
at the right binary. To upgrade, install the new `grove`, then run
`grove init --check`, `grove init` and commit what changed
([Entrypoint revisions](docs/commands.md#entrypoint-revisions)).

## Where each subject lives

| Subject | Owner |
| --- | --- |
| Command usage | `grove --help` |
| Command behaviour | [Command reference](docs/commands.md), [record model](docs/record-model.md) |
| The board | [docs/board.md](docs/board.md) |
| Record types, fields, statuses, validation and lifecycle | [Record model](docs/record-model.md), printed by `grove guide model` |
| Executing assigned work (`grove-work`) | [Work guide](docs/work-execution.md), loaded by the [Claude](.claude/skills/grove-work/SKILL.md) and [Codex](.agents/skills/grove-work/SKILL.md) adapters; [G-032](grove/G-032-dogfood-review.md) records which invocations were exercised |
| Shaping ideas into proposed work (`grove-shape`) | [Shaping guide](docs/work-shaping.md), loaded by the [Claude](.claude/skills/grove-shape/SKILL.md) and [Codex](.agents/skills/grove-shape/SKILL.md) adapters; [G-050](grove/G-050-shaping-review.md) records what was exercised |
| How well the guides steer an agent, measured (behavioral evaluations) | [evals/README.md](evals/README.md); [G-108](grove/G-108-workflow-evals.md) owns the outcome |
| Why Grove exists and its selected direction | [The brief](grove/brief.md) |
| Progress and the next action on any work | That work's record: `grove list`, or `/` on the board |
| Developing Grove, and its constraints | [CLAUDE.md](CLAUDE.md) |
| An old typed ID such as `W-001` | [G-069](grove/G-069-migration-map.md) |

## Develop Grove

[CLAUDE.md](CLAUDE.md) holds the development and verification policy, the
repository's hooks, CI and tooling, and keeps developing Grove apart from
using Grove to track this repository's own work. `AGENTS.md` is a symlink to
it, so Codex reads the same file.
