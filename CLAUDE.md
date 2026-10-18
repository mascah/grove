# Grove

Read `CONSTITUTION.md` first. It is the only file that says what is in scope
and how to work here; this file describes the code and never restates scope
(rule 20). `AGENTS.md` is a symlink to this file.

## What it is

A Go CLI. The `grove` command sets up a project, lists, checks, reviews and
archives Markdown specs and changes, and opens a reader. Skills that `grove
init` writes for Claude Code and Codex run `grove instructions <skill>` and
follow what it prints. What Grove does is specified by `docs/` (rule 8).

## Layout

| Path | What |
| --- | --- |
| `cmd/grove` | `main`, calling `cli.Main`. |
| `internal/cli` | Parses commands and wires the other packages together; the only package that writes to the terminal. One file per command. |
| `internal/instruct` | The embedded `skills.yaml` and `skills/` (Grove's skills and what it requires of each), `metadata/` (the `.grove.yaml` templates), the built-in change schema in `schemas/default/`, and building one skill's instructions as a value before printing. |
| `internal/doc` | Reads Grove's documents: `.grove.yaml`, diagnostics, IDs, the folder layout. |
| `internal/gitx` | The only place Git runs, by argv. `TreeFS` reads a commit's tree as an `fs.FS`. |
| `internal/tui` | The reader (Bubble Tea, Glamour). Knows nothing about Git or Grove's documents. |
| `docs/`, `README.md` | User documentation, and the specification. |
| `dev/walkthrough.md` | The manual test script: every docs step run on the fixture, in Claude Code and in Codex. |
| `dev/fixture.md`, `testdata/ledger/` | The ledger fixture, a small Python project. |
| `dev/reuse.md` | Code copied from earlier repositories, with path and revision. |

## Commands

```sh
just build                 # bin/grove
eval "$(just dev-path)"    # build, and put this checkout's bin/ first on PATH
just check                 # gofmt, vet, tests and the line budget; run before every commit
just fixture ~/tmp/ledger  # a fresh fixture repository to walk through
just install               # the owner's everyday grove, in ~/.local/bin
```

Guides, templates and schemas are embedded, so `just build` picks up an
edit to them. After changing the skill text or a harness row (both in
`internal/cli/init.go`), run `grove init` again in the fixture.

## Tests

`internal/cli/walkthrough_test.go` follows the walkthrough on real disposable
Git repositories built from the fixture (rule 17). `internal/tui/tui_test.go`
covers the reader.
