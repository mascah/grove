# Working on Grove

This repository is a fresh product restart, with a Go CLI for read-only
inspection, configuration, and operational records. `docs/restart-brief.md` is
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
  W-001 implements the reader, W-002 record creation with shared allocation,
  and W-003 field updates with content revisions and a shared write lock.
  W-009 implements the read-only terminal board that bare `grove` opens, with
  Bubble Tea v2 in `internal/tui`; keep explicit subcommands noninteractive,
  and the board's text escaping and exact source targeting intact. D-004 selects
  a future project-wide current view; W-024 owns that projection and may replace
  mandatory version picking in ordinary browsing, while retaining explicit source
  inspection and freshness checks. The current board is still checkout-scoped. W-013
  reads every branch through one `git cat-file` process, scoped to what the
  project loader reads; do not add a Git process per branch. W-012 reads a
  record's Git history only while its card is open, as a read any key may
  cancel; do not read history during the board load. The runner contract
  remains open.
- Keep deterministic validation and state changes in software where useful;
  do not assume software can replace judgment instructions or prove acceptance.
- Record settled choices in the brief while it remains small. Progress and the
  concrete next action belong to each work record's Next, never the brief. Do
  not create a second editable account of the same direction.
  W-018 owns the selected interactive adoption milestone; its linked roadmap is
  investment order, not a batch assignment. W-019 implements schema 2: term,
  plan, and review records (D-005) and the `brief:` key that `grove brief`
  reads. D-006 selects stable identity/placement and flexible knowledge: W-030
  owns neutral IDs, general pages and flat creation with recursive discovery;
  it is not implemented schema yet. D-004's target lifecycle is also pending:
  W-020 owns that migration. W-029 follows W-030 to reconcile all existing Grove
  records, the brief and legacy `docs/plans`/`docs/reviews` content into one flat
  layout with neutral IDs, a durable old-to-new mapping and repaired references.
  This explicit one-time conversion is separate from normal operations keeping
  IDs and paths stable. Until support ships, use the current record model and
  execution guide; do not hand-author future IDs or general pages.
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
  and consistency checks. For Go changes run the relevant tests, the full suite
  (`go test ./...`), `go vet ./...`, `gofmt -l .`, and
  `go run ./cmd/grove check`; use race tests (`go test -race ./...`) for
  relevant changes and uncached runs (`-count=1`) for final evidence. TUI work
  also needs terminal lifecycle and connected-workflow checks.

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
  Comparisons with the predecessor are history in `docs/reviews/`, not
  required reading.
- `context` is facts, not authorization, and must stay read-only. It reads the
  selected records in full and lists the rest; read plans, prerequisites,
  questions, `docs/record-model.md`, and the brief when the guide's step needs
  them, not up front.
- Work branches are `worktree-W-012`, or `worktree-W-012-W-014` for several
  IDs, in a linked worktree under `.claude/worktrees/`. The default base is
  `main` only when it holds the selected records; the guide says what to do
  when it does not. Do not merge or push unless the assignment says so.
- New plans and review evidence are records: `go run ./cmd/grove new plan`
  or `new review`, then `update` to set `work` (and a review's `examined`
  commit). Those written before W-019 stay in `docs/plans/` and
  `docs/reviews/`, linked from the owning record, until W-029 migrates them. Reconcile the README or
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
