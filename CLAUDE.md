# Grove

`CONSTITUTION.md` is the only file that says what is in scope. This file
describes the code and how agents work here. `CONTRIBUTING.md` says how to
build, test, commit and release. `AGENTS.md` is a symlink to this file.

## What it is

A Go CLI. The `grove` command sets up a project, lists, checks, reviews and
archives Markdown specs and changes, and opens a reader. Skills and commands
that `grove init` writes for each harness run `grove instructions <skill>` and
follow what it prints. What Grove does is specified by the maintained specs in `grove/specs/`, explained by `docs/` (rule 8).

## Layout

| Path | What |
| --- | --- |
| `cmd/grove` | `main`, calling `cli.Main`. |
| `internal/cli` | Parses commands and wires the other packages together; the only package that writes to the terminal. One file per command. |
| `internal/instruct` | The embedded `skills.yaml` and `requires/` (Grove's skills and what it requires of each), `metadata/` (the `.grove.yaml` templates), the built-in change schemas in `schemas/default/` and `schemas/bugfix/`, and building one skill's instructions as a value before printing. |
| `internal/doc` | Reads Grove's documents: `.grove.yaml`, diagnostics, IDs, the folder layout. |
| `internal/gitx` | The only place Git runs, by argv. `TreeFS` reads a commit's tree as an `fs.FS`. |
| `internal/tui` | The reader (Bubble Tea, Glamour). Knows nothing about Git or Grove's documents. |
| `docs/`, `README.md` | User documentation, explaining the specs in `grove/specs/`. |
| `grove/` | Grove's own specs, changes and cycles: the specification (rule 8). |
| `mkdocs.yml`, `.github/`, `.goreleaser.yaml`, `install.sh` | The docs site, CI, the release pipeline and the install script. |
| `dev/walkthrough.md` | The manual test script: every docs step run on the fixture. |
| `dev/fixture.md`, `testdata/ledger/` | The ledger fixture, a small Python project. |
| `dev/reuse.md` | Code copied from earlier repositories, with path and revision. |

## How agents work here

- Read `CONSTITUTION.md` first in every session and every sub-agent brief;
  quote its rule number for anything called fixed or out of scope.
- No fact is written in two of `CONSTITUTION.md`, `CLAUDE.md` and
  `CONTRIBUTING.md`; one links to another instead. Scope lives only in the
  constitution.
