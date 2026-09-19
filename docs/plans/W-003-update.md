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

- [x] `repo` helpers and `create` write-lock participation; existing tests green.
- [x] `Revision`, `show --json`, exported parsing; JSON tests.
- [x] Editor with byte-exact fixtures.
- [x] `Apply` with lock, validation, publication, fault injection, races.
- [x] CLI `update`, workflow fixture, docs; `go test ./...`,
  `go test -race ./...`, `go vet ./...`, `gofmt -l .`, `go run ./cmd/grove check`.
- [x] Independent review; address blocking findings; reconcile records.

## Progress and evidence

Prepared against `8b23636`; W-003 set active on 2026-09-19. Implemented on
`worktree-W-003`: `56df181` (shared lock helpers, `new` takes the write lock),
`a0cf898` (`show --json`, exported parsing), `c4a906c` (`update`), `135cd4d`
(documentation), `14ed115` (review fixes).

- `go build`, `go vet ./...`, `gofmt -l .` (clean), `go test ./...`, and
  `go test -race ./...` pass for `internal/cli`, `internal/create`,
  `internal/project`, `internal/repo`, and `internal/update`.
- Editor fixtures (`edit_test.go`) assert exact output bytes for: same-line
  plain values with a tab after the colon and a retained inline comment; a
  double-quoted key with a single-quoted `''`-escaped Unicode value; a block
  sequence at the key's indentation and an indented one, replaced from the
  colon with the last item's inline comment retained; a multi-line flow list
  with a trailing comment; a literal block scalar whose content contains `#`;
  a folded scalar followed by a blank line; a multi-line plain scalar; unset
  removing whole lines and inline comments while standalone comments remain;
  appends before the closing delimiter with `updated` last; several changes in
  one pass; BOM plus CRLF with an indented mapping and a body without a final
  newline; a flow mapping with a Unicode key, a `]` inside a quoted item,
  first/last entry removal, and appends after the last value; refusals for
  tagged or anchored edited entries, a mapping value, and a missing closing
  delimiter, while an unrelated tagged entry still serves as a bound.
- `Apply` tests: every accepted field on work, question, and decision records
  with a title holding leading/trailing spaces, quotes, a colon, `#`, and CJK;
  optional removal and explicit empty lists; all lifecycle values including
  reopening without file or `created` changes; no-op requests preserving
  bytes, mtime, and a `0600` mode with a clock behind the record; stale no-ops
  refused with the current revision; absent versus empty list; twenty-four
  invalid requests refused without writes, including forbidden and wrong-type
  fields, malformed priorities and lists, duplicates, self-links, unresolved
  and non-work targets, and a status carrying a newline; dependency and
  membership cycles; an invalid neighbor; clock before `created` refused,
  equal-second changes sharing a stamp with distinct revisions; a missing
  `created` staying absent; a body-only edit invalidating the revision;
  neighbor, inventory, configuration, permission, and same-inode changes
  injected between validation and rename all refused with no temp file left;
  injected write, sync, close, compare, and rename failures leaving the
  original bytes; injected directory-sync and final-validation failures
  returning `*Failure` with the applied path and revision, permissions kept;
  a same-revision race with exactly one changed success; reciprocal
  dependency additions with exactly one success and a valid project; `new`
  and an `update` from a linked worktree both blocked while the write lock is
  held by another descriptor and proceeding after release, with the lock file
  retained; a non-Git directory refused.
- CLI tests: nineteen usage errors exit 2 without output; help works without a
  project; operation errors exit 1 with `grove:` diagnostics and unchanged
  record files; an output failure after publication reports the applied path
  and revision, and after a no-op reports that no change was needed; the
  workflow fixture creates W-002 with `new`, sets five fields at once, applies
  a no-op, marks it done while unsetting priority, reopens it, and passes
  `check` and `list` with the body skeleton and `created` intact; `show
  --json` hash/source pairs agree for BOM, CRLF, Unicode, and a missing final
  newline, with no Git state created.
- Real use in this worktree: `show W-003 --json`, a no-op `update` returned
  `changed: false` with the same revision, a status/priority change and its
  reversal left only `updated` different, a stale `--expect` was refused, and
  `check` reported 9 records. Dogfooding found the untouched-field guard
  comparing the `priority` pointer address, which refused every update of a
  record carrying a priority; fixed with a regression test before commit.
- Independent review (reviewer agent, `de3a37c..c4a906c`, in a throwaway copy):
  no blocking findings. Should-fix: `Edit` refused explicit `!!str` tags,
  unreferenced anchors, and explicit `? key` entries that the reader accepts;
  tags and anchors are now edited (a replaced value drops them) with six
  fixtures in `14ed115`, and explicit-key entries remain refused without
  writing as a recorded limitation. Nits taken: the flow-mapping bound now
  comes from scanning the mapping rather than the last `}` in the frontmatter;
  the duplicate-field usage message reads naturally; a fixture states that a
  standalone comment inside a removed list goes with the list while one after
  it stays; the untouched-field guard has a direct test. Nits recorded, not
  changed: acceptance item 1's field and invalid-request coverage lives in
  `internal/update` tests with a CLI subset; the overlap, key-position, bound,
  and plain-span guards have no mutation-detecting test. The reviewer also
  ran 9.6M fuzz executions over reader-accepted sources without a rejected,
  drifted, or body-altered result, 29 hostile titles round-tripping exactly,
  and confirmed the write lock is load-bearing by disabling it.
- Not verified: Windows (no `flock`), separate clones, explicit-key (`? key`)
  frontmatter, and a direct editor writing after the final comparison, which
  the specification states as an honest limit.
