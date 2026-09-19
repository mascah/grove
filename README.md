# Grove

A local project workspace for humans and agents, built around a CLI and durable files.

**Status: the first Go CLI lists, shows, validates, creates, and updates local project records.**

Start with [the restart brief](docs/restart-brief.md) and
[the accepted record model](docs/record-model.md). The brief records the selected
direction, the reasoning from the existing Grove skills and nullsec workflow,
open design questions, and the next useful experiment. [grove.yaml](grove.yaml)
configures the record tree. The completed first implementation is
[CLI inspection](grove/work/W-001-inspect-records.md).
The installed `grove` still belongs to the sibling skills project; it does not
read this new format.

## Use the CLI

Requires Go 1.26 or later. Run from this repository:

```sh
go run ./cmd/grove list
go run ./cmd/grove show W-001
go run ./cmd/grove check
go run ./cmd/grove new work "Title of the work" --slug short-name
go run ./cmd/grove show W-001 --json
go run ./cmd/grove update W-001 --expect sha256:HEX --set status=active --unset size
go run ./cmd/grove --project /path/to/project check
```

The first three commands read live files without modifying them; `new` adds
one file and prints its path. `list` shows ID, type,
status, and title; `show` prints the exact Markdown source, or with `--json`
one object holding the path, a `sha256:` content revision, and the source;
`check` validates metadata and relationships. `update` changes frontmatter
fields of one record when its file still hashes to `--expect`, keeps every
other byte of the file, sets `updated`, and prints the resulting revision.
Lists are JSON arrays such as `'["W-001"]'`; `--unset` removes an optional
field. Project/file context and errors go to stderr, so
stdout can be redirected. Exit codes are 0 for success, 1 for inspection/output
errors, and 2 for invalid command usage.

Without `--project`, discovery searches upward for `grove.yaml` and stops at
the current Git checkout boundary. Plain directories also work. Any invalid
record makes the command fail; no partial list or record is printed.

`new` and `update` require a Git checkout. `new` takes the next `W-`, `Q-`,
or `D-` number from a counter under the repository's common Git directory,
shared by every linked worktree, and floors it by the highest ID on any local
ref or worktree. Never number new records by hand. Both commands serialize
through a write lock in that same directory; `update` refuses a stale
`--expect`, an invalid project, or any change it observes while preparing the
write, and reports when a failure happened after the file was replaced.

Use `go test ./...`, `go test -race ./...`, and `go vet ./...` for verification.
Cross-branch views, TUI, and agent execution remain future work.

The intended experience combines linked work, questions, research, project
knowledge, and evidence. A CLI serves agents and humans; a TUI can make the
project visible and eventually launch agent sessions. An optional local web UI
can operate on the same records. Core project operations require no hosted
service or persistent daemon.

## Local restart

On 2026-09-18, the first Grove application checkout was preserved intact at
`../grove-archive-2026-09-18/` and this fresh repository replaced `../grove/`.
The archived checkout retains its Git history and ignored local data. Its last
commit was `be40e46`. Existing deployments and external data were not changed.
The archive's instructions and service-based planning authority belong to that
old application; they do not govern this restart.

`../skills/` remains the working Grove skill suite and CLI. `../nullsec/` remains
the real project providing workflow evidence. Neither is migrated by this reset.

## Resume this conversation

Give a new agent this prompt:

> Read AGENTS.md, docs/restart-brief.md, and docs/record-model.md. Continue the file-backed
> Grove CLI and interactive workspace described there. The old application was
> archived. The Go list/show/check/new/update CLI and its records now work
> locally; W-001, W-002, and W-003 record verification. Follow the brief's
> next action. Treat the brief's
> remaining proposals as proposals. Inspect
> the sibling skills and nullsec projects through their Grove CLI when evidence
> is needed. Preserve this direction and update the brief as choices settle.
