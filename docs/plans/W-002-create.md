# W-002 record creation implementation plan

Goal: deliver `grove new` per [W-002](../../grove/work/W-002-create-records.md)
using the accepted [D-003](../../grove/decisions/D-003-allocator-mechanism.md)
mechanism. Spec: [record model](../record-model.md#identity-and-dates).

Execution: single agent, sequential allocator → creation → command integration.
Reason: one new package with one owner; the command depends on the allocator,
and W-001's reader is reused unchanged. Dispatching implementers and reviewers
for roughly three hundred lines would cost more than one independent review of
the finished diff, which this plan schedules before close.
Delegation: one independent review of the whole branch before close.
Runtime: interactive session in worktree `.claude/worktrees/W-002`, branch
`worktree-W-002`, base `dcbcdb7`. Git 2.50 and Go 1.26.2 available.
Reassess: if the allocator needs a second lock primitive or Windows support.

## Implementation decisions

- `internal/create`: `Allocate(root, recordDir, prefix string, report io.Writer)`
  and `New(p *project.Project, kind, title, slug string, now time.Time, report io.Writer)`.
  `project.Project` gains `RecordDir` (relative) so creation reuses its config.
- Git access through `os/exec` only in `create`; the reader stays Git-free.
  Commands: `rev-parse --path-format=absolute --git-common-dir --show-prefix`,
  `for-each-ref --format=%(objectname) refs/heads refs/remotes`,
  `grep -h -I -E <pattern> <refs> -- <prefix><recordDir>`, `worktree list --porcelain`.
- State under `<common>/grove/`: `lock` (flock) and `next-ids` (`W 3` lines),
  written via temp file, fsync, rename while locked.
- Floor scan on every allocation: ref grep plus a walk of every worktree's live
  `<recordDir>` for `.md` files matching `^id:\s*["']?W-(\d+)["']?\s*$`. Body
  lines that look like IDs only raise the floor; gaps are acceptable.
- Frontmatter: id, type, title (Go-quoted, YAML-compatible), status default
  (`proposed`, `open`, `proposed`), `created`/`updated` equal. Body skeleton per type.
- Slug: explicit `--slug` must match `^[a-z0-9-]+$`; otherwise derive from the
  title per the model (32 chars, fallback `record`).
- After `O_EXCL` creation, reload the project; any diagnostic fails the command
  without deleting the file. stdout: root-relative path. stderr: allocation notes.
- Outside Git: exit 1, no file, no state. `list`/`show`/`check` never create state.

## Steps and verification

- [ ] Allocator tests: floor from refs and a linked worktree with live files,
  counter initialization message, counter-below-floor correction, reservation
  consumed by a failed creation, concurrent goroutines across two worktrees get
  distinct IDs, non-Git directory refuses. Confirm failure, implement, go green.
- [ ] Creation and command tests: file name and frontmatter validate via
  `project.Load`, title quoting, explicit and derived slugs, existing file not
  overwritten, usage errors, read commands leave `<common>/grove` absent.
- [ ] Verify: `go test ./...`, `go test -race ./...`, `go vet ./...`, `gofmt -l`;
  create a real record in this worktree and confirm `list`/`check`; confirm the
  kill-while-locked probe from D-003 against the built binary; independent
  review of the branch; reconcile W-002 and docs.

## Progress and evidence

Prepared against `dcbcdb7`.
