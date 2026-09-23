# Working on Grove

Grove is a Go CLI and terminal board over Markdown records with YAML
frontmatter, kept in ordinary files and Git, with no required running
service. This file is this repository's development policy. Other facts have
one owner each; read them when the task needs them, not up front:

- [`grove/brief.md`](grove/brief.md): purpose, constraints and selected
  direction. It keeps selected direction, observed evidence and proposed
  design apart; preserve those distinctions. Read it for anything that shapes
  direction or is not bounded by a work record.
- [`README.md`](README.md): what the CLI and board do, command by command.
- [`docs/record-model.md`](docs/record-model.md): configuration, record schema,
  validation and lifecycle rules.
- [`docs/work-execution.md`](docs/work-execution.md) and
  [`docs/work-shaping.md`](docs/work-shaping.md): the shared workflows.
- Work records under `grove/`: each outcome, its acceptance, evidence and
  Next. Progress and the concrete next action belong to a record's Next,
  never the brief. [G-069](grove/G-069-migration-map.md) maps old typed IDs
  and paths, for a reference in an old commit or message.

## Invoking the CLI

The installed `grove` (`~/.local/bin/grove`) is a build of this CLI from a
named commit, which `grove version` prints; it can lag this checkout. Here,
every `grove …` is `go run ./cmd/grove …` in the selected checkout. Create
records only with `go run ./cmd/grove new`, never by hand-numbering, and
change status or fields with `go run ./cmd/grove update`. `go run` reports
every failure as exit 1 and prints the command's own code as `exit status N`
on stderr; to tell a usage error (2) from a failure (1), build once with
`go build -o <temp path> ./cmd/grove` and run that.

## Shaping, executing and retrieving context

- **Shape** an idea or existing records into proposed work with
  `/grove-shape TOPIC` (Claude) or `$grove-shape TOPIC` (Codex), which load
  `docs/work-shaping.md`. Shaping never assigns, implements or merges. A
  headless proposal branch is `worktree-shape-SLUG` in a linked worktree
  under `.claude/worktrees/`.
- **Execute** assigned work IDs with `/grove-work G-030` or `$grove-work
  G-030`, which load `docs/work-execution.md`. `grove run G-030` (or `R` on
  the board) starts `/grove-work G-030 --interaction headless` as a
  Grove-owned attempt that outlives the terminal.
- **Retrieve context in stages.** Start from `grove context IDs`: it reads the
  selected records in full and lists the rest. Read plans, prerequisites,
  questions, the record model and the brief when the guide's step needs them.
  `context` is facts, not authorization, and stays read-only. A listing is
  not a reading.

The adapters in `.claude/skills/` and `.agents/skills/` only load the guides
and are kept unmanaged here on purpose: they read the guide files this
checkout's `go run` builds from, where `grove init` elsewhere writes managed
adapters that read the binary's copy (`grove guide`). The guides are workflow
and write commands as `grove …`; this file is the policy for them here, so
neither the guides nor the adapters repeat it:

- Do not invoke the predecessor's `grove:work`, `grove:close` or
  `grove:shape` skills or its close/archive commands.
- Work branches are `worktree-G-030`, or `worktree-G-030-G-031` for several
  IDs, in a linked worktree under `.claude/worktrees/`. The default base is
  `main` only when it holds the selected records; the guide says what to do
  when it does not. Do not merge or push unless the assignment says so.
- New plans and review evidence are records: `go run ./cmd/grove new plan`
  or `new review`, then `update` to set `work` (and a review's `examined`
  commit). Reconcile the README or `docs/record-model.md` when their contract
  changes, and the brief only when the work changes the direction it selects,
  not to report progress. Record settled choices in the brief; do not create
  a second editable account of the same direction.
- Fixtures that create records belong in a disposable clone reached by an
  explicit absolute `--project` path, never a `cd` that can fail: `new` in a
  worktree of this repository advances the shared ID counter.

## Constraints

Each rule's detail lives in the record named; the rule applies whatever the
record's status.

- **Identity.** `schema_version: 3` is the only schema, and Grove keeps no
  backward compatibility before its first release; inspect an old commit with
  the CLI in that commit. Never hand-author or renumber an ID, and do not
  move or rename record files: IDs and paths stay stable (G-064). `new` allocates
  from one counter in Git's common directory, shared by every linked worktree,
  and `new` and `update` serialize through a shared write lock (G-007, G-009).
  Keep `convert` for documents outside the record root.
- **Lifecycle.** Work runs `proposed`, `active`, `review`, `done`. An
  implementation ends in `review` with its `candidate` commit; `done` is
  written on `main` after the merge, never on the work branch, and `update`
  refuses a candidate HEAD does not contain and a checkout off the configured
  target. `approve`, `feedback` and `integrate` (board `a`, `f`, `i`) record
  the verdict and merge (G-038, G-044). A `done` record without a candidate
  predates that meaning; do not backfill one.
- **Git processes.** Every Git process Grove or its tests start goes through
  `repo.Command`, never a bare `exec.Command("git", …)`: it drops `GIT_DIR`
  and the other variables that name a repository, because Git exports
  `GIT_DIR` to hooks and a child that inherits it acts on the real repository
  (G-089: test fixtures run by a pre-push hook from a linked worktree
  committed into this repository). A hook that runs tests must scrub them as
  well.
- **Reading branches.** Every branch is read through one `git cat-file`
  process, scoped to what the project loader reads, merge bases included; do
  not add a Git process per branch (G-031, G-042). A record's Git history is
  read only while its card is open, as a read any key may cancel, never during
  the board load (G-030).
- **The board.** Bare `grove` opens it (Bubble Tea v2, `internal/tui`);
  explicit subcommands stay noninteractive. Keep its text escaping and exact
  source targeting intact, with freshness checks before acting on a selected
  version (G-017, G-011). A rendered body is escaped before glamour and
  filtered to glamour's own styles after it; keep both layers (G-043). The
  current view derives from Git ancestry in `internal/versions/current.go`
  (G-042), and `b` still chooses one checkout's own board.
- **Evidence and authority.** `../skills/` and `../nullsec/` are evidence and
  potential compatibility targets, not automatically part of an
  implementation's write scope; follow their own instructions. Nullsec runs on this Grove: read its records with the installed
  `grove` from nullsec's checkout (`grove brief`, `grove list`, `grove
  context G-NNN`); its `G-` numbers are its own, so name the repository when an
  ID could be either. `../skills/` keeps the uninstalled predecessor's files;
  read them as files. G-041 records how the predecessor was removed and how to
  restore it. The archived application's service authority, architecture,
  credentials, deployment procedures and backlog are historical: do not revive
  them or copy private local data into this repository.
- Keep deterministic validation and state changes in software where useful;
  do not assume software can replace judgment instructions or prove acceptance.

## Changes and verification

- Use focused Conventional Commits. Preserve unrelated work and isolate
  concurrent implementation in separate worktrees.
- Verify claims against actual results. Documentation-only changes need link
  and consistency checks. For Go changes iterate with `go test -short
  ./<package>`, then run `go vet ./...`, `gofmt -l .`, `go run ./cmd/grove
  check`, and one `go test -count=1 -timeout 120s ./...` as final evidence.
  Rerun a failed package, not the suite. Reproduce a Linux-only failure with
  `docker run --rm -v "$PWD":/src -w /src -e GOFLAGS=-buildvcs=false golang:1.26 go test ./...`.
- Never `-p 1`, and never `-race` across the suite on macOS: the race runtime
  hangs in the forked child before `exec` there, and the test timeout leaves
  it behind at full CPU. Run `-race` only per package, only for a concurrency
  change, with `-timeout 120s`, and treat a hang in `syscall.forkExec` as that
  toolchain bug, not evidence; kill any `*.test` process a timeout leaves
  behind.
- Budget: no package over five seconds, and no test that builds, sleeps, or
  waits on a shim without a `-short` skip and a comment saying why. TUI work
  also needs terminal lifecycle and connected-workflow checks.
