# Working on Grove

This repository is a fresh product restart, with a Go CLI for read-only
inspection, configuration, and operational records. `grove/brief.md` is
the current source of intent; it distinguishes selected direction, observed
evidence, and proposed design, and those distinctions must be preserved. Read
it first for anything that shapes direction or is not bounded by a work record.
Assigned work IDs start from their records instead (see Assigned work) and
read the brief when the record, a product question, or reconciliation needs it.

- Keep the product useful through ordinary local files and a CLI without a
  required running service. The core record model in `docs/record-model.md`
  is accepted: Markdown with YAML frontmatter. Use Go for the first CLI.
  The starter file defaults use short sequential IDs, replacing the random-ID
  trial. Allocation coordinates across local worktrees through Git's common
  metadata directory; reading records requires no allocator state.
  G-003 implements the reader, G-007 record creation with shared allocation,
  and G-009 field updates with content revisions and a shared write lock.
  G-017 implements the read-only terminal board that bare `grove` opens, with
  Bubble Tea v2 in `internal/tui`; keep explicit subcommands noninteractive,
  and the board's text escaping and exact source targeting intact. G-035 selects
  a future project-wide current view; G-042 owns that projection and may replace
  mandatory version picking in ordinary browsing, while retaining explicit source
  inspection and freshness checks. The current board is still checkout-scoped. G-031
  reads every branch through one `git cat-file` process, scoped to what the
  project loader reads; do not add a Git process per branch. G-030 reads a
  record's Git history only while its card is open, as a read any key may
  cancel; do not read history during the board load. The runner contract
  remains open.
- Keep deterministic validation and state changes in software where useful;
  do not assume software can replace judgment instructions or prove acceptance.
- Record settled choices in the brief while it remains small. Progress and the
  concrete next action belong to each work record's Next, never the brief. Do
  not create a second editable account of the same direction.
  G-036 owns the selected interactive adoption milestone; its linked roadmap is
  investment order, not a batch assignment. G-037 added term, plan, and review
  records (G-051) and the `brief:` key that `grove brief` reads. G-064 selected
  stable identity/placement and flexible knowledge, which G-065 implemented:
  neutral `G-` IDs from one counter, `page` records, flat creation, recursive
  discovery, and `update --set type=`. G-035's target lifecycle is pending:
  G-038 owns that migration. G-052 reconciled this repository on 2026-09-21:
  every record, former `docs/plans` and `docs/reviews` document included, has
  a neutral ID flat under `grove/`, the brief is `grove/brief.md`, and
  [G-069](grove/G-069-migration-map.md) maps each old ID and path to its
  counterpart; use it to follow a reference in an old commit or message.
  G-052 also deleted schemas 1 and 2 with their typed IDs, type folders and
  counters: `schema_version: 3` is the only schema, Grove keeps no backward
  compatibility before its first release, and an old commit is inspected with
  the CLI in that commit. Ordinary operations keep IDs and paths stable: never
  hand-author or renumber an ID, and keep `convert` for documents outside the
  record root.
- `../skills/` and `../nullsec/` are evidence and potential compatibility targets,
  not automatically part of an implementation's write scope. Follow their
  instructions; retrieve their Grove knowledge through `grove status`,
  `grove find`, and `grove context` from the relevant repository.
- The installed `grove` command currently belongs to the sibling skills project.
  Use `go run ./cmd/grove` for this restart's `grove.yaml` and `grove/` records;
  create records with `go run ./cmd/grove new`, never by hand-numbering, and
  change status or fields with `go run ./cmd/grove update`.
  Do not use the predecessor's initialization or validation commands here.
- The archived application's service authority, architecture, credentials,
  deployment procedures, and backlog are historical. Do not revive them as
  requirements for this project or copy private local data into this repository.
- Use focused Conventional Commits. Preserve unrelated work and isolate
  concurrent implementation in separate worktrees.
- Verify claims against actual results. Documentation-only changes need link
  and consistency checks. For Go changes iterate with `go test -short
  ./<package>`, then run `go vet ./...`, `gofmt -l .`, `go run ./cmd/grove
  check`, and one `go test -count=1 -timeout 120s ./...` as final evidence.
  Rerun a failed package, not the suite. Never `-p 1`, and never `-race`
  across the suite on macOS: the race runtime hangs in the forked child before
  `exec` there, and the test timeout leaves it behind at full CPU. Run `-race`
  only per package, only for a concurrency change, with `-timeout 120s`, and
  treat a hang in `syscall.forkExec` as that toolchain bug, not evidence; kill
  any `*.test` process a timeout leaves behind. Budget: no package over five
  seconds, and no test that builds, sleeps, or waits on a shim without a
  `-short` skip and a comment saying why. TUI work also needs terminal
  lifecycle and connected-workflow checks.

## Assigned work

To carry out assigned work IDs, follow `docs/work-execution.md`; the
`grove-work` skill adapters in `.claude/skills/` and `.agents/skills/` load it
and are kept in this repository on purpose while the workflow is dogfooded.
That guide is the workflow and writes commands as `grove …`. This section is
this repository's policy for it, kept here so that neither the guide nor the
adapters repeat it:

- Every `grove …` in the guide is `go run ./cmd/grove …` in the selected
  checkout, because the installed `grove` is the predecessor's (above).
  `go run` reports every failure as exit 1 and prints the command's own code as `exit status N` on
  stderr; to tell a usage error (2) from a failure (1), build once with
  `go build -o <temp path> ./cmd/grove` and run that.
- Do not invoke the predecessor's `grove:work`, `grove:close`, or
  `grove:shape` skills or its close/archive commands for work here.
  Comparisons with the predecessor are history in review records, not
  required reading.
- `context` is facts, not authorization, and must stay read-only. It reads the
  selected records in full and lists the rest; read plans, prerequisites,
  questions, `docs/record-model.md`, and the brief when the guide's step needs
  them, not up front.
- Work branches are `worktree-G-030`, or `worktree-G-030-G-031` for several
  IDs, in a linked worktree under `.claude/worktrees/`. The default base is
  `main` only when it holds the selected records; the guide says what to do
  when it does not. Do not merge or push unless the assignment says so.
- New plans and review evidence are records: `go run ./cmd/grove new plan`
  or `new review`, then `update` to set `work` (and a review's `examined`
  commit). Reconcile the README or
  `docs/record-model.md` when their contract changes, and the brief only when
  the work changes the direction it selects, not to report progress.
- Fixtures that create records belong in a disposable clone reached by an
  explicit absolute `--project` path, never a `cd` that can fail: `new` in a
  worktree of this repository advances the shared ID counter.

## Shaping work

To turn an idea or existing records into proposed work, follow
`docs/work-shaping.md`; the `grove-shape` adapters beside `grove-work` load it.
The Assigned work policy above applies unchanged: `grove …` is
`go run ./cmd/grove …`, predecessor skills such as `grove:shape` are not used
here, and fixtures that create records use a disposable clone. A headless
proposal branch is `worktree-shape-SLUG` in a linked worktree under
`.claude/worktrees/`. Shaping never assigns, implements, or merges.
