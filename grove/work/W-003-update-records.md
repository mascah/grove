---
id: "W-003"
type: work
title: "Update record status and fields from the CLI"
status: active
kind: feature
priority: 3
size: medium
members: []
depends_on: []
relates_to: ["W-002", "D-003"]
created: "2026-09-19T15:36:19Z"
updated: "2026-09-19T18:55:41Z"
---

## Outcome

Change one record's status and fields from the CLI, preserving human-authored
Markdown and refusing stale updates. Grove's own records are the first data.
The owner approved this specification on 2026-09-19 and authorized
implementation; the [implementation plan](../../docs/plans/W-003-update.md)
maps acceptance to checks and records evidence.

## Why now

Finish the local record workflow before cross-branch coordination. W-002's
creation command is implemented; updating supplies the mutation boundary a
future UI can reuse. Priority 3 is unchanged. Medium reflects source-preserving
frontmatter edits and concurrency checks within one bounded command.

## Constraints

- Follow the accepted [record model](../../docs/record-model.md). Preserve ID,
  type, path, creation date, body bytes, and unrelated frontmatter bytes.
- Require Git for updates in this first version, as for creation. Plain-directory
  inspection remains supported. Do not change configuration, refs, or index.
- Apply one request to one record. No body editing, rename, move, delete,
  cross-branch routing, workflow enforcement, or generated reason/history prose.
- A valid status is an explicit assertion, including when reopening. The command
  does not establish acceptance, review, or integration.

## Design

### Command and revision interface

Proposed commands (not implemented yet):

```sh
grove show W-003 --json
grove update W-003 --expect sha256:HEX --set status=active --set priority=2
grove update W-003 --expect sha256:HEX --set 'depends_on=["W-002"]' --unset size
```

`show ID --json` returns one JSON object with exactly `id`, `path`, `revision`,
and `source`. Path is project-relative; source decodes to the exact UTF-8 file
bytes. Revision is `sha256:` followed by 64 lowercase hexadecimal digits of the
SHA-256 hash of those bytes, including BOM and line endings. Compute source and
revision from the same read buffer. Existing plain `show` stdout stays exact;
project/file context stays on stderr. Reads create no coordination state. W-004
will use this content-revision convention without needing update behavior.

`update ID` requires exactly one `--expect REVISION` and at least one `--set
FIELD=VALUE` or `--unset FIELD`. Accept repeated set/unset flags for different
fields. Reject any field mentioned twice, including set plus unset. Split a set
argument at the first `=`; preserve the rest as its value. Support both separate
option arguments and `--option=value`, existing `--project` placement, and help
without a project. Missing arguments, malformed revision syntax, or duplicate
options are usage errors (exit 2). Invalid field names, forbidden fields,
wrong-type fields, invalid values, missing IDs, and stale revisions are
operation errors (exit 1). No force/ignore-revision option.

Fields accepted by the update command:

| Records | Fields |
| --- | --- |
| All | `title`, `status`, `relates_to` |
| Work | `kind`, `priority`, `size`, `members`, `depends_on` |
| Questions | `blocks` |

String values are literal nonempty strings subject to the model's validation;
do not trim meaningful title whitespace. Priority uses decimal digits for an
integer 1 through 5. Relationship values are JSON arrays of strings; reject null,
wrong types, duplicate entries, invalid targets, self-links, and cycles through
existing schema/graph validation. Lists replace the whole ordered list. `[]`
sets an explicit empty list; `--unset` removes an optional field. Removing an
absent optional field is a no-op. Setting an absent field to an empty list adds
that field. Required fields cannot be unset. ID, type, created, and updated
cannot be set or unset through this command.

On success, stdout is one JSON object with exactly `id`, `path`, `revision`, and
`changed` (boolean), followed by a newline. The revision describes the resulting
bytes; a no-op returns the original revision. Diagnostics use stderr. This
result lets callers continue without hashing a later observation. Success
returns 0; failures return 1 except usage errors described above.

### Preservation and time

Compare requested values by parsed meaning and list order, while distinguishing
an absent optional field from a present one. If all changes are no-ops, preserve
all file bytes, permissions, and dates; still enforce the expected revision.

Apply byte-range edits to changed frontmatter entries and `updated`. Preserve
all unrelated entries, comments, ordering, delimiters, BOM, existing line
endings, and the complete body after the closing delimiter. Support the forms
already accepted by the reader, including block/flow mappings and lists,
quoted keys/strings, multiline scalar values, and CRLF. Whole-mapping YAML
serialization does not meet preservation acceptance.

The changed value may use canonical YAML encoding. Preserve comments outside
that value's syntax span; comments inside a replaced list or multiline value
may be removed with that value. On unset, remove the entry's key/value and
syntactically required separator, plus its inline comment; retain standalone
comment lines and neighboring entries. In a flow mapping, adjust only the
separator/adjacent spacing needed to remove or append that entry. Append new
fields in request order at the mapping's end, with newly inserted `updated`
last. Use the opening delimiter's newline style for inserted lines. Editing
valid syntax must never fall back to rewriting unrelated entries; fail without
writing if the implementation cannot safely identify a span, and treat any
such refusal for an accepted fixture as an implementation gap.

For a changed record, preserve `created`, including its absence, and set
`updated` to current UTC time truncated to seconds. If that time precedes an
existing `created` or `updated`, refuse without writing and report clock/date
inconsistency. Do not invent a future timestamp to make validation pass. Multiple
changes within one second may share an updated timestamp; revisions remain
content-based. A no-op does not need a clock adjustment and succeeds even when
existing dates are ahead of the clock.

### Cooperating writers and publication

Introduce a shared advisory `flock` at `<git-common-dir>/grove/write.lock` for
record mutations across local worktrees. Keep D-003's allocator `grove/lock`
and `grove/next-ids` reservation contract unchanged. Never unlink a lock file;
process exit releases the held lock. A missing write lock is created on first
mutation; read commands never create it. Update does not initialize or advance
ID counters. The coarse repository-wide write lock is sufficient for this
first local-worktree audience.

Both `new` and `update` participate. `new` first reserves its ID under the
existing allocator lock and releases that lock, then takes the write lock,
reloads and validates the selected project, creates exclusively, and performs
its existing final validation. No code holds both locks at once. If the
configuration changed since the allocation input, or subsequent creation
fails, retain the consumed reservation and report the failure. Keep existing
creation output and no-overwrite behavior. Calls using an older Grove binary
or a direct editor do not participate in the new write-lock guarantee.

For update, acquire the write lock before loading the authoritative project
snapshot. Validate the whole selected project, find the unique target, check
`--expect`, and build the candidate in memory. Validate the project with that
candidate substituted before touching the target; this operation cannot repair
an already invalid project. Expose a shared validation path rather than copying
schema or graph rules into the writer.

For a change, prepare a temporary regular file beside the target with a name
that does not end in `.md`, write and sync its bytes, and preserve the target's
permission bits. Close it before publication. Immediately before replacement,
re-read configuration and the record inventory/content and compare them with
the snapshot used to validate the candidate; refuse a detected change. Verify
the destination still names the same regular file with the same permissions.
Replace through same-directory atomic rename and sync the parent directory.
Reload and validate the project while still holding the write lock, then return
the result. Clean temporary files on handled failures; an abrupt process kill
may leave a non-record temporary file but must release the lock.

A direct editor can still write after the last comparison. This is an honest
limit: cooperating current Grove commands are serialized, and observed external
changes are refused, but there is no atomic transaction against arbitrary file
or Git writers. No filesystem rollback can safely promise otherwise.

Before rename succeeds, invalid requests, stale input, and preparation failures
leave record bytes unchanged. Once rename succeeds, a directory-sync, final
validation, or output error returns failure and explicitly reports that the
update was applied (or was a no-op for output failure). Include the target path;
include the resulting revision when available. Never blindly restore old bytes
over a possible later edit. An internal result/error must retain publication
state so CLI output failure cannot erase this distinction. Diagnostics are
best-effort if stderr itself is unavailable. A failed directory sync reports
uncertain durability, not a successful durable write.

## Acceptance

1. CLI fixtures exercise setting every allowed field on its permitted types,
   multiple changes together, optional removal, explicit empty lists, literal
   title punctuation/Unicode, and all invalid requests described above.
2. Status updates support all valid lifecycle values and reopening without body
   entries, renaming, or creation-date changes. No acceptance/review policy is
   inferred from status, relationships, headings, or prose.
3. Same-value updates, same ordered lists, and removing absent fields preserve
   full-file hashes and dates. Absent versus explicit empty is tested. Stale
   expectations still refuse a request that otherwise would be a no-op.
4. Hash/source pairs agree exactly for plain and JSON inspection, including BOM,
   CRLF, Unicode, and a body without its final newline. Existing inspection
   behavior and its zero allocator/lock side effects remain intact.
5. Byte-preservation fixtures cover accepted block/flow forms, quoted keys,
   multiline strings, inline and standalone comments, changed-list comments,
   unrelated field formatting, insertion, and clearing. Semantic validity alone
   is insufficient; assert unchanged byte ranges as well as final values.
6. Invalid candidate relationships, duplicate/self links, membership cycles,
   dependency cycles, and already invalid projects fail without record changes.
   Missing creation dates stay absent; clock rollback and same-second changes
   obey the stated time contract.
7. Detect stale callers after body-only changes with unchanged timestamps, and
   detected file/configuration/inventory/permission changes during preparation.
   Two processes updating one record with the same revision produce exactly one
   changed success; the other refuses without losing the winning change.
8. Concurrent updates to different records serialize full-project validation:
   reciprocal dependency additions cannot both succeed and create a cycle.
   Concurrent new/update operations participate in the same publication lock.
   Include linked-worktree lock serialization and killing a lock holder; W-002's
   allocation, gaps, recovery, and no-overwrite tests continue to pass.
9. Inject temporary-write, sync, close, rename, directory-sync, final-validation,
   and output failures. Verify original bytes before publication, applied-state
   reporting after publication, permissions, handled-error cleanup, and no
   destructive rollback. Do not claim rejection for an undetectable external
   edit after the final comparison.
10. A CLI fixture creates a record, updates multiple fields, marks it done,
    reopens it, and checks the project. This proves the command workflow rather
    than the truth of the record's prose acceptance.

## Handoff boundary

Code inspected at `bea8e92` (product code unchanged from `ee69c42`). Expected
interfaces are `internal/cli`, shared parsing/graph validation in
`internal/project`, a new update package, and publication coordination reused
by `internal/create`. Keep source-editing logic separate from filesystem writes
so preservation checks do not need concurrency fixtures.

Execution recommendation: one Fable agent in an isolated worktree, W-003 only.
It owns these overlapping interfaces through verification and review. Runtime
is a user-started Fable session; no Grove agent runner is assumed. Reassess scope
if preserving accepted YAML forms requires a substantially larger parser.

After owner review of this specification, prepare the implementation plan with
acceptance-to-check mapping, then implement in the selected worktree. Required
checks are relevant fixtures, `go test ./...`, `go test -race ./...`,
`go vet ./...`, and gofmt. Obtain independent review focused on source
preservation, concurrent writers, and publication failure reporting. Reconcile
README, model, brief, and this work record with actual evidence on closure.
Implementation completion is in its branch context; integration remains a
separate step. Do not start W-004/W-005 as part of this assignment.

## Next

Implement per the [plan](../../docs/plans/W-003-update.md) on
`worktree-W-003`, then independent review and closure. While Fable implements W-003, finish W-004's output/source
selector contract and W-005's consuming contract. Implement those sequentially
after reviewing the actual shared interfaces; W-004's read-only behavior has no
hard dependency on updates, while W-005 requires W-004.
