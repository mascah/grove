# Working on Grove

Read `docs/restart-brief.md` first. This repository is a fresh product restart,
with a Go CLI for read-only inspection, configuration, and operational records.
The brief distinguishes selected
direction, observed evidence, and proposed design; preserve those distinctions.

- Keep the product useful through ordinary local files and a CLI without a
  required running service. The core record model in `docs/record-model.md`
  is accepted: Markdown with YAML frontmatter. Use Go for the first CLI.
  The starter file defaults use short sequential IDs, replacing the random-ID
  trial. Future allocation coordinates across local worktrees through Git's
  common metadata directory; reading records requires no allocator state.
  W-001 implements the reader, W-002 record creation with shared allocation,
  and W-003 field updates with content revisions and a shared write lock.
  W-009 implements the read-only terminal board that bare `grove` opens, with
  Bubble Tea v2 in `internal/tui`; keep explicit subcommands noninteractive,
  and the board's text escaping and explicit version selection intact. The
  runner contract remains open.
- Keep deterministic validation and state changes in software where useful;
  do not assume software can replace judgment instructions or prove acceptance.
- Record settled choices and the concrete next action in the brief while it
  remains small. Do not create a second editable account of the same direction.
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
  (`go test ./...`), and `go vet ./...`; use race tests for relevant changes.
