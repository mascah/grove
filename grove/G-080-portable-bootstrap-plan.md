---
id: "G-080"
type: plan
title: "G-040 portable bootstrap plan"
status: current
created: "2026-09-22T16:05:22Z"
updated: "2026-09-22T16:06:06Z"
work: ["G-040"]
---

## Design

Prepared 2026-09-22 on `worktree-G-040` from main `ccdc92d`, against
[G-040](G-040-portable-bootstrap.md) at `sha256:5d171634…` under
[G-064](G-064-stable-knowledge.md)'s bounds. The record delegates packaging,
update behaviour and the distribution mechanism to preparation; the choices
below are routine technical ones inside the outcome, with no product choice
left open.

Observed at `ccdc92d`: the two guides are ordinary files, `docs/work-execution.md`
and `docs/work-shaping.md`, and this repository's four adapters under
`.claude/skills/` and `.agents/skills/` tell the agent to read those files and
`AGENTS.md`; the guides link into `../grove/G-…` records and the adapters, so
printed elsewhere those links are dead. `cli.Run` loads the project before
dispatching every command but the board and help, so a command that needs no
`grove.yaml` must branch before `project.Load`. `go build` stamps the main
module's pseudo-version and `vcs.revision`/`vcs.modified` into the binary,
`go run` stamps `(devel)` only, and `go install pkg@rev` stamps the
pseudo-version. The installed `grove` on this machine is the predecessor
(Python; `init`, `lint`, `status`, …), which refuses `version` and `guide` as
unknown commands. nullsec and the skills project use `grove.toml`, not
`grove.yaml`, so a target with predecessor knowledge has no `grove.yaml` to
conflict with.

1. **One owner for the workflow: the binary.** A root package embeds both
   guide files (`//go:embed docs/work-execution.md docs/work-shaping.md`), so
   the files stay the editable owner here and every built Grove carries the
   guides of its own commit. `grove guide work` and `grove guide shape` print
   them, needing no project. The workflow version is therefore the executable
   version, which `grove version` prints from Go's build info: the module
   version, then the VCS revision and `modified` when stamped. Nothing is
   copied into a target repository as an editable guide.
2. **Portable guide text.** The guides drop their relative links into this
   repository's records and adapters, keeping the sentences as plain
   references (Grove's own G-035, G-038, G-032, G-050), and the "any agent
   without skills" invocation rows read the guide through `grove guide`. This
   repository's own adapters keep reading the files, since `go run` builds
   from them; the generated adapters read the binary. Both load the one
   owner.
3. **`grove init`** sets up the directory named by `--project` (default the
   working directory), which must be the top of a Git checkout: `new` needs
   Git, and a nested `grove.yaml` inside a checkout would shadow discovery.
   It plans every path first and writes only when nothing conflicts, then
   prints one line per path, `created`, `unchanged`, `updated` or `kept`,
   and `conflict PATH: reason` with exit 1 and no write at all otherwise.
   - `grove.yaml`: created with `schema_version: 3`, `records: grove`,
     `brief: grove/brief.md`; kept when it exists and loads as schema 3 (its
     own `records` and `brief` then name the paths below); a conflict when it
     exists and does not load, quoting the diagnostics.
   - The record root: created as a directory; a conflict when a file is there.
   - The brief: created as a placeholder that says it is not written yet and
     states no intent, naming the shaping skill as the way to develop it; kept
     whenever a regular file exists. Never rewritten.
   - Six managed entrypoints: `.claude/skills/grove-work/SKILL.md`,
     `.claude/skills/grove-shape/SKILL.md`, `.agents/skills/grove-work/SKILL.md`,
     `.agents/skills/grove-shape/SKILL.md` and the two Codex
     `agents/openai.yaml` policies. Each carries a marker line, "Managed by
     grove init". A file with the marker is compared with the template:
     `unchanged` or `updated`, which is the managed update. A file without the
     marker is the user's and is `kept`, with the note that deleting it lets
     init write the managed version. A directory, symlink or unreadable
     file at a managed path is a conflict.
   - `AGENTS.md` and `CLAUDE.md` are never read or written: the harnesses load
     them themselves and the generated adapters defer to them for how the CLI
     is invoked. No `--force`, no removal, no rename.
4. **Generated adapters** are the thin form of this repository's: the same
   assignment-as-data paragraph, then "run `grove guide work` and follow it",
   with `grove` meaning the Grove on `PATH` unless the repository's agent
   instructions say otherwise, and an instruction to stop when the command
   fails or prints something else, since that means another `grove` answered.
   Claude adapters keep `disable-model-invocation: true` and the argument
   hint; Codex adapters keep `allow_implicit_invocation: false`.
5. **Distribution and upgrade.** A built binary from a named commit:
   `go build -o DIR/grove ./cmd/grove` from a checkout, or
   `go install github.com/mascah/grove/cmd/grove@COMMIT` into `GOBIN`. No
   release tags, installer, or self-update: upgrading is building again, which
   changes the guides with the code, and rerunning `init` where the adapter
   templates changed. Coexistence with the predecessor: build to a directory
   that is first on `PATH` only in the shells that use the new Grove, never
   over the installed one; `grove version` tells them apart because the
   predecessor rejects the command. The README documents this.
6. **Boundary.** The portable package is the binary and what `init` writes.
   `go run ./cmd/grove`, worktree names, the Go tests, `AGENTS.md` and this
   repository's records stay here; the generated adapters and the printed
   guides name none of them.

## Steps

All done on 2026-09-22; evidence and the handoff are in
[G-040](G-040-portable-bootstrap.md) and [G-082](G-082-portable-bootstrap-review.md).

1. [x] Plan committed with the record revisions above (`73b2044`).
2. [x] Root package `grove` with the embedded guides; `guide` and `version`
   commands in `internal/cli`, dispatched before the project loads; usage text
   (`7ddbf52`; `version` gained a guide digest in `66fcccc`).
3. [x] `init` in `internal/cli/init.go`: templates, plan-then-write, the report
   lines; tests in `init_test.go` for creation, `check` passing afterwards, a
   second run reporting `unchanged`, a marked edit reported `updated`, an
   unmarked file `kept`, an invalid `grove.yaml` and a nested directory
   refused with nothing written (`7ddbf52`; symlinked parents and the default
   path added by the review, `66fcccc`, `3113280`).
4. [x] Guide link edits (design item 2); README section "Adopt Grove in
   another repository"; record model's configuration section names `init`;
   AGENTS.md names the new commands (`7ddbf52`, `731674e`, `eb32209`).
5. [x] Evidence in a disposable Git repository under the session scratchpad,
   with a binary built from a clone of the candidate first on `PATH`: init,
   check, rerun, managed update, conflict; `claude -p` and `codex exec` fresh
   sessions invoking the generated `grove-shape` and `grove-work` skills
   headless, reported as observed behaviour apart from file checks (G-082).
6. [x] `go test -short`, `go vet ./...`, `gofmt -l .`, `grove check`,
   `go test -count=1 -timeout 120s ./...`; independent review of the combined
   diff, two rounds; review record G-082; handoff into Review.
