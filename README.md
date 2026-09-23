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
| `run`, `attempts`, `attempt`, `stop` | Start one headless agent attempt of a work item that outlives the terminal; list, inspect and stop attempts | [Attempts](docs/commands.md#attempts) |
| `versions` | Show each record's versions across branches and worktrees, current or older | [Versions](docs/commands.md#versions) |
| `workspace` | Print the checkout holding a selected version | [Workspace](docs/commands.md#workspace) |
| `context` | Assemble staged context for selected work | [Context](docs/commands.md#context) |
| `init` | Set up Grove in a Git checkout | [Init](docs/commands.md#init) |
| `guide`, `version` | Print a workflow guide; name this build | [Version and guide](docs/commands.md#version-and-guide) |

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
grove version     # must print "grove v…"; a usage error means another grove answered
grove init        # at the checkout's top, or: grove --project /absolute/path init
grove check
```

Then `/grove-shape TOPIC` develops the brief and proposes work, and
`/grove-work G-001` carries it out (`$grove-shape` and `$grove-work` in
Codex). [Init](docs/commands.md#init) says what `init` writes and how to
point Codex at the right binary.

## Where each subject lives

| Subject | Owner |
| --- | --- |
| Command usage | `grove --help` |
| Command behaviour | [Command reference](docs/commands.md), [record model](docs/record-model.md) |
| The board | [docs/board.md](docs/board.md) |
| Record types, fields, statuses, validation and lifecycle | [Record model](docs/record-model.md) |
| Executing assigned work (`grove-work`) | [Work guide](docs/work-execution.md), loaded by the [Claude](.claude/skills/grove-work/SKILL.md) and [Codex](.agents/skills/grove-work/SKILL.md) adapters; [G-032](grove/G-032-dogfood-review.md) records which invocations were exercised |
| Shaping ideas into proposed work (`grove-shape`) | [Shaping guide](docs/work-shaping.md), loaded by the [Claude](.claude/skills/grove-shape/SKILL.md) and [Codex](.agents/skills/grove-shape/SKILL.md) adapters; [G-050](grove/G-050-shaping-review.md) records what was exercised |
| Why Grove exists and its selected direction | [The brief](grove/brief.md) |
| Progress and the next action on any work | That work's record: `grove list`, or `/` on the board |
| Developing Grove, and its constraints | [AGENTS.md](AGENTS.md) |
| An old typed ID such as `W-001` | [G-069](grove/G-069-migration-map.md) |

## Develop Grove

[AGENTS.md](AGENTS.md) holds the development and verification policy, the
repository's hooks, CI and tooling.
