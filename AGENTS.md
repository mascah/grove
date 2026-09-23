# Working on Grove

Grove is a Go CLI and terminal board over Markdown records with YAML
frontmatter, kept in ordinary files and Git, with no required running
service. This file is this repository's development policy. Every other fact
has one owner; read it when the task needs it, not up front:

- [`grove/brief.md`](grove/brief.md): purpose, constraints and selected
  direction, with selected direction, observed evidence and proposed design
  kept apart. Read it for anything that shapes direction or is not bounded by
  a work record.
- [`README.md`](README.md): what Grove is and where each subject is owned.
- [`docs/commands.md`](docs/commands.md) and [`docs/board.md`](docs/board.md):
  command and board behaviour beyond `grove --help`.
- [`docs/record-model.md`](docs/record-model.md): configuration, record schema,
  validation and lifecycle.
- [`docs/work-execution.md`](docs/work-execution.md) and
  [`docs/work-shaping.md`](docs/work-shaping.md): the shared workflows.
- Work records under `grove/`: each outcome, its acceptance, evidence and
  Next. [G-069](grove/G-069-migration-map.md) maps old typed IDs and paths.

## Invoking the CLI

- Here, every `grove …` is `go run ./cmd/grove …` in the selected checkout;
  the installed `~/.local/bin/grove` is a build from the commit
  `grove version` prints and can lag it.
- Create records only with `go run ./cmd/grove new`, never by hand-numbering,
  and change status or fields only with `go run ./cmd/grove update`.
- `go run` reports every failure as exit 1 and the command's code as
  `exit status N` on stderr; to tell usage (2) from failure (1), build once
  with `go build -o <temp path> ./cmd/grove` and run that.

## Shaping, executing and retrieving context

- Shape with `/grove-shape TOPIC` (Claude) or `$grove-shape TOPIC` (Codex);
  shaping never assigns, implements or merges.
- Execute assigned IDs with `/grove-work G-030` or `$grove-work G-030`;
  `grove run G-030` (board `R`) starts `/grove-work G-030 --interaction
  headless` as a Grove-owned attempt that outlives the terminal.
- Retrieve context in stages from `grove context IDs`, which reads the
  selected records in full, lists the rest, and is facts, not authorization;
  a listing is not a reading.
- The guides write commands as `grove …`, and this file is their policy
  here; the adapters in `.claude/skills/` and `.agents/skills/` only load
  this checkout's guide files and stay unmanaged on purpose.
- Do not invoke the predecessor's `grove:work`, `grove:close` or
  `grove:shape` skills or its close/archive commands.
- Work branches are `worktree-G-030`, or `worktree-G-030-G-031` for several
  IDs, and a headless proposal branch is `worktree-shape-SLUG`, each in a
  linked worktree under `.claude/worktrees/`.
- The base is `main` only when it holds the selected records; otherwise the
  work guide says what to do.
- Do not merge or push unless the assignment says so.
- Plans and reviews are records: `go run ./cmd/grove new plan` or
  `new review`, then `update` to set `work` and a review's `examined`.
- When a contract changes, reconcile the document that owns it; change the
  brief only when the direction it selects changes, never to report
  progress, and never keep a second editable account of that direction.
- Progress and the next action belong in a record's Next, never the brief.
- Fixtures that create records run in a disposable clone reached by an
  explicit absolute `--project` path, never a `cd` that can fail, since `new`
  in a worktree of this repository advances the shared ID counter.

## Constraints

Each rule's reasons live in the record named, and the rule applies whatever
that record's status.

- `schema_version: 3` is the only schema, with no backward compatibility
  before the first release; inspect an old commit with the CLI in that
  commit.
- Never hand-author or renumber an ID, and never move or rename a record
  file; keep `convert` for documents outside the record root (G-064).
- `new` allocates from one counter in Git's common directory, shared by every
  linked worktree, and `new` and `update` serialize through a shared write
  lock (G-007, G-009).
- Work runs `proposed`, `active`, `review`, `done`, and an implementation
  ends in `review` with its `candidate` commit (G-038).
- `done` is written on `main` after the merge, never on the work branch, and
  `update` refuses a candidate HEAD lacks and a checkout off the target
  (G-038).
- `approve`, `feedback` and `integrate` (board `a`, `f`, `i`) record the
  verdict and merge (G-044).
- Never backfill a candidate on a `done` record that has none (G-038).
- Every Git process Grove or its tests start goes through `repo.Command`,
  never a bare `exec.Command("git", …)`, and a hook that runs tests scrubs
  `GIT_DIR` and the other repository variables as well (G-089).
- Read every branch through one `git cat-file` process scoped to what the
  project loader reads, merge bases included, never a process per branch
  (G-031, G-042).
- Read a record's Git history only while its card is open, as a read any
  key may cancel, never during the board load (G-030).
- Bare `grove` opens the board (Bubble Tea v2, `internal/tui`), and explicit
  subcommands stay noninteractive.
- Keep the board's text escaping and exact source targeting, with freshness
  checks before acting on a selected version (G-017, G-011).
- Escape a rendered body before glamour and filter it to glamour's own
  styles after; keep both layers (G-043).
- The current view derives from Git ancestry in
  `internal/versions/current.go`, and `b` still chooses one checkout's own
  board (G-042).
- `../skills/` and `../nullsec/` are evidence and potential compatibility
  targets, not automatically in an implementation's write scope; follow
  their own instructions.
- Read nullsec's records with the installed `grove` from nullsec's checkout
  (`grove brief`, `grove list`, `grove context G-NNN`), and name the
  repository when a `G-` ID could be either's.
- Read `../skills/`, the uninstalled predecessor, as files; G-041 records how
  it was removed and how to restore it.
- The archived application's service authority, architecture, credentials,
  deployment procedures and backlog are historical: do not revive them or
  copy private local data into this repository.
- Keep deterministic validation and state changes in software where useful,
  and do not assume software can replace judgment or prove acceptance.

## Changes and verification

- Use focused Conventional Commits, preserve unrelated work, and isolate
  concurrent implementation in separate worktrees.
- Verify claims against actual results; documentation-only changes need link
  and consistency checks.
- For Go changes iterate with `go test -short ./<package>`, then run
  `go vet ./...`, `gofmt -l .`, `go run ./cmd/grove check`, and one
  `go test -count=1 -timeout 120s ./...` as final evidence.
- Rerun a failed package, not the suite.
- Reproduce a Linux-only failure with `docker run --rm -v "$PWD":/src -w /src
  -e GOFLAGS=-buildvcs=false golang:1.26 go test ./...`.
- Never `-p 1`, and never `-race` across the suite on macOS, where the race
  runtime hangs in the forked child before `exec`.
- Run `-race` only per package, only for a concurrency change, with
  `-timeout 120s`; a hang in `syscall.forkExec` is that toolchain bug, not
  evidence, and kill any `*.test` process a timeout leaves behind.
- No package over five seconds, and no test that builds, sleeps, or waits on
  a shim without a `-short` skip and a comment saying why.
- TUI work also needs terminal lifecycle and connected-workflow checks.
