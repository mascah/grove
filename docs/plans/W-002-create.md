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
  `for-each-ref --format=%(objectname) refs/heads refs/remotes refs/tags`,
  `grep -h -I -e ^id: <refs> -- <recordDir>` (Go validates each line), `worktree list --porcelain`.
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

- [x] Allocator tests: floor from refs and a linked worktree with live files,
  counter initialization message, counter-below-floor correction, reservation
  consumed by a failed creation, concurrent goroutines across two worktrees get
  distinct IDs, non-Git directory refuses. Confirm failure, implement, go green.
- [x] Creation and command tests: file name and frontmatter validate via
  `project.Load`, title quoting, explicit and derived slugs, existing file not
  overwritten, usage errors, read commands leave `<common>/grove` absent.
- [x] Verify: `go test ./...`, `go test -race ./...`, `go vet ./...`, `gofmt -l`;
  create a real record in this worktree and confirm `list`/`check`; confirm the
  kill-while-locked probe from D-003 against the built binary; independent
  review of the branch; reconcile W-002 and docs.

## Progress and evidence

Prepared against `dcbcdb7`; implemented on `worktree-W-002` (`7e15ccf`, lock
test `17cea9d`). Tests were written first and failed to compile before each step.

- `go build`, `go vet ./...`, `gofmt -l .` (clean), `go test ./...`, and
  `go test -race ./...` pass for `internal/cli`, `internal/create`, and
  `internal/project`.
- Allocator tests: floor 6 from committed W-003 on a branch plus live W-005 in
  a linked worktree, initialization message once, shared counter continues at
  7 from the worktree, Q starts at 1; a counter of 2 below committed W-004
  corrects to 5 with a message and persists `W 6`; twenty goroutines across two
  worktrees receive 2..21 without duplicates; a non-Git directory is refused
  with no state written.
- Creation tests: `W-002-title-with-quotes-more.md` with a title containing a
  colon and quotes validates through the W-001 reader with equal timestamps;
  question and decision defaults; an existing target path fails creation and
  the next attempt receives the following number; bad slugs and types create
  nothing; CLI usage errors exit 2; read commands leave `<common>/grove` absent;
  `new` prints the root-relative path and `check` validates the result.
- Lock test: a helper process holds the lock, the parent cannot acquire it for
  200 ms, the helper is killed with SIGKILL, and the parent acquires it at once.
- Real use: `go run ./cmd/grove new work "Update record status and fields from
  the CLI" --slug update-records` in this worktree reported initialization at 3
  from local refs and worktrees, wrote `grove/work/W-003-update-records.md`,
  and `check` reported 7 valid records; `.git/grove/next-ids` holds `W 4`.
- Independent review (opus reviewer agent, `dcbcdb7..17cea9d`): one blocking
  finding, the worktree walk discarded I/O errors so an unreadable record in
  another worktree silently lowered the floor and reissued its ID; fixed in
  `8d359e3` by propagating walk and read errors (a missing record folder stays
  normal) with a test that hides a `W-030` file. Should-fix items also taken:
  the ref-scan test now deletes the committed file so only the ref scan can
  find W-007; a persistence-failure test proves no ID and no file; `list` is
  asserted after creation; tags joined the ref scan; the Windows stub was
  dropped in favor of a compile error. The reviewer verified nested
  `grove.yaml` projects, thirteen title shapes, W-999 to W-1001 ordering, and
  twelve concurrent processes by probe. Remaining nits recorded: no directory
  fsync after rename, and an uninformative error if a repository has enough
  refs to exceed the argument limit.
- Not verified: Windows (does not compile) and separate clones.
