# W-003 record update implementation plan

Goal: deliver `show ID --json` and `update ID` per
[W-003](../../grove/work/W-003-update-records.md), the canonical specification.
Schema: [record model](../record-model.md#on-disk-contract).

Execution: single agent, sequential: shared Git/lock helpers → revision and
JSON inspection → source-preserving frontmatter editor → locked publication →
command integration → creation participation. Reason: one new package whose
parts depend on each other; the editor is the only design risk and is testable
without concurrency fixtures. Delegation: one independent review of the whole
branch before close, focused on source preservation, concurrent writers, and
publication failure reporting. Runtime: interactive session in worktree
`.claude/worktrees/W-003`, branch `worktree-W-003`, base `8b23636`. Git 2.50,
Go 1.26.2. Reassess: if preserving accepted YAML forms needs more than a
bounded span scanner over yaml.v3 node positions.

## Implementation decisions

- `internal/repo`: `CommonDir(root)` (one `git rev-parse` call) and the `flock`
  helper moved from `create`, so `create` and `update` share them without a
  dependency cycle. Write lock: `<common>/grove/write.lock`, never unlinked.
- `internal/project`: `Revision(source)` returns `sha256:` + 64 hex digits;
  `ParseRecord` and `Validate` are exported so the writer reuses the reader's
  schema and graph rules instead of copying them.
- `internal/update/edit.go`: `Edit(source, changes)` locates the frontmatter
  block on raw bytes (BOM, CRLF retained), parses it with yaml.v3 for entry
  order and rune-based key/value start positions, then scans each edited
  value's end by style (plain, single/double quoted, literal/folded, block or
  flow sequence). The next key's start bounds every span; an unexpected byte
  at a computed position refuses the edit rather than guessing. Same-line
  values are replaced in place, keeping an inline comment; a value that starts
  on a later line is replaced from after the colon. Unset removes whole lines
  in block mappings and one entry plus one separator in flow mappings. New
  fields append before the closing delimiter (block) or after the last value
  (flow), request order, `updated` last, using the opening delimiter's newline.
  Tagged or anchored edited entries are refused as unsupported spans.
- Canonical encodings: `%q` for strings (matches `new`), decimal for priority,
  `["A", "B"]` or `[]` for lists, quoted UTC seconds for `updated`.
- `internal/update/update.go`: `Apply(root, Request, now, fault)`: take the
  write lock, load and validate, find the record, compare `--expect` with the
  content revision, parse values, drop no-op fields by parsed meaning, refuse
  a clock behind existing dates, edit, re-parse the candidate, check unchanged
  fields are byte-for-byte equal in meaning, validate the graph with the
  candidate substituted, then publish: temp file beside the target (name not
  ending in `.md`), write, sync, chmod, close; re-load and compare
  configuration, inventory, and sources with the snapshot; verify the target
  is the same regular file with the same permissions; rename; sync the
  directory; reload and validate. `Result{Applied}` and a typed error retain
  publication state so the CLI reports an applied update on later failures.
  `fault(step)` is a test hook for injected failures and external writes.
- `create.New`: reserve under the allocator lock, release, take the write lock,
  reload, refuse a changed record root, create exclusively, validate.
- CLI: `show ID --json`; `update ID --expect REV --set F=V... --unset F...`.
  Duplicate fields or options, malformed revisions, and missing arguments exit
  2; everything else exits 1. Output failure after publication reports the
  applied path and revision on stderr.

## Acceptance to checks

| Acceptance | Check |
| --- | --- |
| 1 fields, multiple changes, removal, empty lists, Unicode, invalid requests | `TestUpdateFields*`, `TestUpdateRejects*` in `internal/update`; CLI usage tests |
| 2 lifecycle values and reopening | `TestUpdateLifecycle` |
| 3 no-op preservation, absent vs empty, stale no-op | `TestUpdateNoOp*` |
| 4 hash/source agreement, BOM/CRLF/no final newline, read side effects | `TestShowJSON*` in `internal/cli` |
| 5 byte preservation forms | `TestEdit*` table in `edit_test.go` asserting exact bytes |
| 6 invalid candidates and already invalid projects, missing created, clock | `TestUpdateRejects*`, `TestUpdateClock*` |
| 7 stale after body edit, changes during preparation, same-revision race | `TestUpdateStale*`, `TestUpdateDetectsExternal*`, `TestUpdateSameRevisionRace` |
| 8 cross-record serialization, new/update lock, worktrees, killed holder | `TestUpdateReciprocalDependencies`, `TestNewAndUpdateShareLock`, `internal/repo` lock test, existing `create` tests |
| 9 injected failures | `TestUpdateFault*` |
| 10 workflow fixture | `TestUpdateWorkflow` in `internal/cli` |

## Steps and verification

- [ ] `repo` helpers and `create` write-lock participation; existing tests green.
- [ ] `Revision`, `show --json`, exported parsing; JSON tests.
- [ ] Editor with byte-exact fixtures.
- [ ] `Apply` with lock, validation, publication, fault injection, races.
- [ ] CLI `update`, workflow fixture, docs; `go test ./...`,
  `go test -race ./...`, `go vet ./...`, `gofmt -l .`, `go run ./cmd/grove check`.
- [ ] Independent review; address blocking findings; reconcile records.

## Progress and evidence

Prepared against `8b23636`; W-003 set active on 2026-09-19.
