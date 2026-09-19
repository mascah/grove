---
id: "W-003"
type: work
title: "Update record status and fields from the CLI"
status: proposed
kind: feature
priority: 3
size: medium
members: []
depends_on: []
relates_to: ["W-002", "D-003"]
created: "2026-09-19T15:36:19Z"
updated: "2026-09-19T17:50:16Z"
---

## Outcome

Change a record's status and planning fields from the CLI so agents and the
future workspace can move work without hand-editing frontmatter. Grove's own
records are the first data.

## Why now

Created by `grove new` as W-002's first real record. Priority 3: creation and
inspection are enough to dogfood; this becomes urgent when a UI or run needs
to change state. Medium reflects revision checks and frontmatter rewriting
that preserve the human-authored body byte for byte.

## Constraints

- Preserve `created`; set `updated` only when content changes. Compare actual
  content, not timestamps, before writing; refuse when the file changed
  underneath the command.
- Rewrite only the changed frontmatter fields; leave the body and unrelated
  fields untouched. Validate the whole project after the write.
- Respect the lifecycle values and planning vocabulary in the record model.
  Defer renames, moves, deletes, and cross-branch edits.

## Acceptance

- Set status and each optional work field, including clearing one, with the
  file otherwise identical.
- A concurrent edit between read and write is detected and refused.
- Fixture tests cover valid transitions, invalid values, unchanged no-ops that
  leave `updated` alone, and unchanged file hashes on failure.

## Design proposal for handoff

Shaping inspected main at `ee69c42`, where W-002 is already present. This
section is a proposal for review, not an accepted extension of the record
model or an instruction to begin implementation.

### Smallest useful outcome

One command changes one record in the explicitly selected checkout. It can
apply several field changes together, allowing an agent to set status and
planning metadata without intermediate partially applied requests. Keep the
existing file path and identity. Cross-branch selection belongs to
[Q-001](../questions/Q-001-branch-versions.md).

Recommended field boundary: `title` and `status` on every record;
`kind`, `priority`, `size`, `members`, and `depends_on` on work;
`blocks` on questions; `relates_to` on every type. Optional fields can be
removed. Relationship updates replace an ordered list; incremental add/remove
operations can follow if actual usage warrants them. `id`, `type`, `created`,
and `updated` are not user-settable through this command.

Allow any status value valid for the record type, including reopening. Do not
require a reason or generate body entries. The accepted model gives status its
meaning; software validates its value without asserting that acceptance,
review, or integration happened. A mandatory reason would introduce workflow
policy and body editing into this first field-update operation.

### Revision and write boundary

Recommend both a caller-supplied content revision and a fresh comparison at
write time. They address separate cases: an agent acting on an old observation,
and a file changing while the command prepares an update. Derive the revision
from all original file bytes, including frontmatter, body, BOM, and line endings;
do not derive it from `updated`. Expose the revision with the exact source it
describes through the read interface. Keep existing `show` stdout byte-exact.
The final command syntax and revision transport still need to be specified.

Serialize cooperating Grove writers across validation and publication. Re-read
the project after acquiring the write lock and validate the proposed replacement
against the complete record set before changing the target. Check relationships
and cycles against the candidate, rather than writing invalid input and relying
on rollback. Creation and updates need an explicit coordination contract;
the existing allocator lock only covers ID reservation.

Preserve body bytes and unrelated frontmatter bytes, including comments and
formatting. Replacing the entire YAML mapping through a serializer does not
satisfy that requirement. Inspect source spans and fixture coverage for block
and flow mappings, multiline values, comments, BOM, and CRLF before choosing
the editing technique. Define how comments attached to a changed or removed
field behave. Do not silently narrow the reader's accepted YAML forms.

A semantic no-op leaves all bytes and timestamps unchanged. A changed record
preserves `created` and updates `updated`; the implementation contract must
also specify behavior when the clock is earlier than an existing date.

Publish a prepared replacement atomically, preserving file permissions. Reject
invalid input, stale revisions, and preparation failures without changing record
bytes. Distinguish these from errors after publication: a failed final project
check or output write must report that the update took effect. Do not blindly
restore an old file over a later edit. The original acceptance statement about
unchanged hashes on failure needs this explicit boundary before implementation.

An advisory lock coordinates participating commands; it does not exclude a
direct editor. Define detected-change refusal and the remaining final
comparison/publication race honestly. Do not promise a transactional project
snapshot against arbitrary external writers.

### Acceptance cases to make executable

Alongside the existing acceptance, the final specification should cover:

- Updating several fields together; clearing optional fields; wrong-type fields;
  invalid lifecycle values; missing targets; self-links; and both cycle types.
- Reopening each record type without adding generated narrative or changing its
  identity, path, or creation date.
- No-op comparisons by field meaning, including an already absent optional
  field, with unchanged full-file hashes.
- A stale caller revision caused only by a body edit or unchanged timestamp;
  an edit during preparation; and two cooperating writers using one revision.
- Byte preservation across the accepted YAML forms, comments, Unicode, BOM,
  CRLF, and a body without a final newline.
- Write, replacement, and post-publication failures, with assertions appropriate
  to which side of publication failed; preserved permissions and temporary-file
  cleanup.
- A CLI fixture that creates a record, updates it, marks it done, reopens it,
  and validates the resulting project. This proves the command workflow, not
  fulfillment of the record's prose acceptance.

### Handoff readiness

Recommend one Fable implementation session in an isolated worktree after this
contract is settled. CLI parsing, project validation, preservation, and write
coordination share interfaces and should be prepared together. W-002 is present
in the inspected base; Q-001 does not block this checkout-local operation.

Before handing off, settle command/revision syntax, supported non-Git mutation
behavior, coordination with `new`, changed-field comment handling, and the exact
failure guarantees above. Then prepare ordered implementation steps against the
actual code and map checks to acceptance. Use the repository's full Go suite,
race suite, and vet, plus independent review of preservation and concurrency.
No Fable runtime or adapter has been invoked or demonstrated in this shaping.

## Next

Review the proposed update boundary, then finish the command and write contract
and prepare W-003 for Fable. The owner selected cross-branch coordination as
the following experience: [W-004](W-004-record-versions.md) proposes version
inspection and [W-005](W-005-record-workspace.md) workspace location. Their
read-only behavior does not depend on updates; coordinate shared CLI and
project-loader interfaces when preparing execution.
