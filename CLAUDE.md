# Working on Grove

Grove is a Go CLI and terminal board over Markdown records with YAML
frontmatter, kept in ordinary files and Git, with no service to run. This
repository is two things at once, and this file is its policy for both:

- **Developing Grove**: the source under `cmd/` and `internal/`, its tests,
  and the documents the binary ships. Those rules say what the code must keep
  doing and how a change is verified.
- **Using Grove here**: this repository tracks its own work as records under
  `grove/`, with the CLI built from this checkout, following the guides it
  ships. Those rules say how records and work move here.

An adopting project gets only the second kind of rule, through `grove init`
and its own instructions file, and never sees this one. `AGENTS.md` is a
symlink to this file so Codex reads the same policy: edit `CLAUDE.md`.

Every other fact has one owner. Read it when the task needs it, not up front:

- [`grove/brief.md`](grove/brief.md): purpose, constraints and selected
  direction, keeping selected, observed and proposed apart.
- [`README.md`](README.md): what Grove is and where each subject is owned.
- [`docs/commands.md`](docs/commands.md) and [`docs/board.md`](docs/board.md):
  command and board behaviour beyond `grove --help`.
- [`docs/record-model.md`](docs/record-model.md): configuration, schema,
  validation and lifecycle (`grove guide model`).
- [`docs/work-execution.md`](docs/work-execution.md) and
  [`docs/work-shaping.md`](docs/work-shaping.md): the two workflows
  (`grove guide work`, `grove guide shape`).
- Records under `grove/`: each outcome, its acceptance, evidence and Next.
  [G-069](grove/G-069-migration-map.md) maps old typed IDs and paths.

Each rule below names the record or document that holds its reasons and
applies whatever that record's status; a rule naming none is this file's own.

## The CLI here

- `grove …`, in the guides and below, means `go run ./cmd/grove …` in the
  checkout you are working in, so a worktree runs its own code. The installed
  `~/.local/bin/grove` is a build of the commit `grove version` prints,
  refreshed by `just install` (the post-merge hook runs it), and can lag.
- `go run` turns every failure into exit 1 and prints the real code as
  `exit status N` on stderr; to tell usage (2) from failure (1), build once
  with `go build -o <temp path> ./cmd/grove` and run that.

## Developing Grove

### What the code must keep doing

- `schema_version: 3` is the only schema, with no backward compatibility
  before the first release; inspect an old commit with the CLI in that
  commit (G-065, G-052).
- `new` allocates from one counter in Git's common directory, shared by every
  linked worktree, and `new` and `update` serialize through a shared write
  lock (G-007, G-009).
- Every Git process Grove or its tests start goes through `repo.Command`,
  never a bare `exec.Command("git", …)`, and a hook that runs tests scrubs
  `GIT_DIR` and the other repository variables as well (G-089).
- Read every branch through one `git cat-file` process scoped to what the
  project loader reads, merge bases included, never a process per branch
  (G-031, G-042).
- The current view derives from Git ancestry in
  `internal/versions/current.go`, and `b` still chooses one checkout's own
  board (G-042).
- Bare `grove` opens the board (`internal/tui`), and explicit subcommands
  stay noninteractive (G-017).
- Read a record's Git history only while its card is open, as a read any
  key may cancel, never during the board load (G-030).
- Keep the board's text escaping and exact source targeting, with freshness
  checks before acting on a selected version (G-017, G-011).
- Escape a rendered body before glamour and filter it to glamour's own
  styles after; keep both layers (G-043).
- A document the binary ships (`grove guide work|shape|model`, the reviewer
  definition) links only within itself or to `https://`, names another
  shipped document by its command, never by path, never links a record or
  names one beyond its own example IDs, and never names Grove's repository
  or the predecessor (G-146, G-151).

### Changes and verification

- Focused Conventional Commits; preserve unrelated work; isolate concurrent
  implementation in separate worktrees.
- Verify claims against actual results; a documentation-only change needs
  link and consistency checks.
- Iterate with `go test -short ./<package>`, then run `go vet ./...`,
  `gofmt -l .`, `go run ./cmd/grove check`, and one
  `go test -count=1 -timeout 120s ./...` as final evidence. Rerun a failed
  package, not the suite.
- Reproduce a Linux-only failure with `docker run --rm -v "$PWD":/src -w /src
  -e GOFLAGS=-buildvcs=false golang:1.26 go test ./...`.
- Never `-p 1`, and never `-race` across the suite on macOS, where the race
  runtime hangs in the forked child before `exec` (G-081). Run `-race` per
  package, only for a concurrency change, with `-timeout 120s`; a hang in
  `syscall.forkExec` is that bug, not evidence, and kill any `*.test`
  process a timeout leaves behind.
- No package over five seconds, and no test that builds, sleeps, or waits on
  a shim without a `-short` skip and a comment saying why (G-071).
- TUI work also needs the terminal lifecycle and connected-workflow checks:
  `python3 internal/tui/testdata/terminal.py BINARY` (Unix; skipped without
  `python3` and under `-short`).
- Run `lefthook install` once per clone: pre-commit formats staged Go files
  and runs `go vet` and `go mod tidy -diff`; there is no pre-push hook.
  GitHub Actions runs the same checks on Ubuntu and macOS as a signal, not a
  gate (G-081).
- `just clean-merged` removes local branches merged into `main` and their
  clean worktrees, after asking.

## Using Grove here

### Records

- Create records only with `go run ./cmd/grove new`, and change status or
  fields with `update`. Never hand-author or renumber an ID, and never move
  or rename a record file; `convert` is for documents outside the record
  root (G-064).
- Plans and reviews are records: `new plan` or `new review`, then `update`
  to set `work` and a review's `examined`.
- Work runs `proposed`, `active`, `review`, `done`, and an implementation
  ends in `review` with its `candidate` commit. `done` is written on `main`
  after the merge, never on the work branch; `update` refuses a candidate
  HEAD lacks and a checkout off the target `grove.yaml` names. Never
  backfill a candidate on a `done` record that has none (G-038).
- `approve`, `feedback` and `integrate` (board `a`, `f`, `i`) record the
  verdict and merge (G-044).
- Progress and the next action belong in a record's Next, never the brief.
  Change the brief only when the direction it selects changes, and never
  keep a second editable account of that direction. When a contract
  changes, reconcile the document that owns it.

### Shaping and executing work

- Shape with `/grove-shape TOPIC` (Claude) or `$grove-shape TOPIC` (Codex);
  shaping never assigns, implements or merges. Execute assigned IDs with
  `/grove-work G-030` or `$grove-work G-030`; `grove run G-030` (board `R`)
  runs the same headless as a Grove-owned attempt that outlives the terminal.
- The adapters in `.claude/skills/` and `.agents/skills/` load this file and
  this checkout's guide files, not the binary's copies, and carry no init
  marker on purpose.
- `grove context IDs` reads the selected records in full and lists the rest;
  a listing is not a reading, and `context` writes nothing and authorizes
  nothing. The guides say what to read at each step.
- Work branches are `worktree-G-030`, or `worktree-G-030-G-031` for several
  IDs, and a headless proposal branch is `worktree-shape-SLUG`, each in a
  linked worktree under `.claude/worktrees/`. The base is `main` only when it
  holds the selected records; otherwise the work guide says what to do.
- Do not merge or push unless the assignment says so.
- Fixtures that create records run in a disposable clone reached by an
  explicit absolute `--project` path, never a `cd` that can fail: `new` in
  any worktree of this repository advances the shared ID counter.

## Outside this repository

- `../skills/` (the uninstalled predecessor, read as files) and `../nullsec/`
  (a pilot adopter, read with its installed `grove` from its own checkout)
  are evidence, never in an implementation's write scope. Follow their own
  instructions, name the repository when a `G-` ID could be either's, and
  never invoke the predecessor's `grove:*` skills or close/archive commands
  (G-041).
- The archived application is history: do not revive its service,
  architecture, credentials, deployment or backlog, and copy no private local
  data into this repository (the brief).
