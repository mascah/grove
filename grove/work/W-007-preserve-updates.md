---
id: "W-007"
type: work
title: "Preserve accepted frontmatter edits and reject changed configuration"
status: proposed
created: "2026-09-19T20:13:54Z"
updated: "2026-09-19T20:21:16Z"
kind: fix
priority: 2
size: medium
relates_to: ["W-003", "W-002"]
---

## Outcome

Valid field-update requests preserve human-authored bytes outside changed
values, and mutations refuse configuration changes observed during preparation.
Proposed repair of the existing W-003 contract. [Review R3–R5](../../docs/reviews/2026-09-19-integrated-cli.md)
reproduced comment loss, multi-unset/explicit-key refusals, and a missed
configuration change. [Implementation plan](../../docs/plans/W-007-preserve-updates.md).

## Constraints

Keep W-003's request syntax, revisions, no-op semantics, timestamps, permission
preservation, shared write lock, publication order, and applied-failure reporting.
No whole-frontmatter serialization, schema narrowing, body edits, cross-branch
writes, or force option. Do not change allocator reservation/lock ordering.

## Design

Locate key, colon, value, comments, and separators separately. For a value that
begins on a later line, retain key-line and standalone comments outside its
syntax span; retain needed indentation/newlines so canonical replacement YAML
parses. Support explicit scalar keys already admitted by the reader.
Plan flow separators across the complete change set, assigning each comma once;
handle adjacent/trailing removals and insertion of absent `updated` together.
Unset removes that entry's inline comment as specified while retaining unrelated
standalone comments. Keep byte-span editing and fail safely for truly ambiguous
spans; accepted fixtures must succeed.

Compare `Project.Config` byte-for-byte in update's final snapshot comparison.
Creation must compare the allocation input's configuration bytes against the
under-lock reload, before publication. Configuration comments/formatting count
as observed changes. A failed creation still consumes its reservation. The review also reproduced the creation behavior with a loaded input followed
by a configuration-only change; retain that regression fixture.

## Acceptance

1. Exact-byte tests preserve a comment in `title: # comment` before a later-line
   scalar and standalone comments before values; cover LF/CRLF, BOM, quoted
   keys/values, lists, and unchanged body bytes.
2. Removing adjacent optional entries at a flow mapping's end succeeds in either
   request order, with/without an existing `updated`, trailing commas, and
   multiline separators. Unset drops its inline comment and retains neighbors.
3. Explicit `? status` / `: proposed` updates and explicit optional-key removal
   succeed with exact unrelated-byte preservation; tagged/anchored accepted
   forms and existing block/flow fixtures remain green.
4. Valid configuration byte changes injected before update's comparison or
   between creation allocation and publication refuse without record publication;
   invalid configurations and changed record roots still refuse. Creation's
   reservation remains advanced; subsequent operations work normally.
5. Existing no-op/stale/clock, concurrent writers, graph validation, permissions,
   pre/post-publication error and cleanup tests pass, plus full suite, race, vet,
   formatting, Grove `check`, and independent preservation/concurrency review.

## Dependencies and handoff

W-003 is integrated. No new product choice or unfinished dependency. Implement
serially after W-006 for ownership coordination; this ordering is not a semantic
`depends_on`. Use one Fable agent and an isolated checkout. Keep review findings
and branch completion distinct from integration. The reader's accepted syntax
is evidence; lack of an old fixture does not authorize narrowing it.

## Next

Use the [combined repair handoff](../../docs/prompts/W-006-W-008-implementation.txt)
after W-006 in the same assignment's isolated worktree.

Fable: reproduce the four update failures in the linked plan on the current
base; implement byte-span/separator repairs and exact configuration guards;
return focused commits, evidence, and independent review. Do not expand the
mutation surface or merge the branch.
