---
id: "W-007"
type: work
title: "Preserve accepted frontmatter edits and reject changed configuration"
status: done
created: "2026-09-19T20:13:54Z"
updated: "2026-09-19T22:07:50Z"
kind: fix
priority: 2
size: medium
relates_to: ["W-003", "W-002"]
---

## Outcome

Valid field-update requests preserve human-authored bytes outside changed
values, and mutations refuse configuration changes observed during preparation.
Repairs the existing W-003 contract (implemented; see Evidence). [Review R3–R5](../../docs/reviews/2026-09-19-integrated-cli.md)
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

## Evidence

Implemented 2026-09-19 on branch `worktree-W-006-W-008` (base `2d6de36`), not
yet integrated into main. Code: `2cc7814` (comments, explicit keys, flow
separators), `dac27fe` (configuration bytes in update and creation), `dedac73`
(review fix). The [plan](../../docs/plans/W-007-preserve-updates.md#implementation-notes-2026-09-19)
records the bounded adjustments. Request syntax, revisions, no-op and clock
behaviour, permissions, lock order, publication steps, and applied-failure
reporting are untouched; editing stays byte-span based.

Reproduced before repair on the unfixed branch: `title: # retain` lost its
comment; removing the final two flow entries failed with `frontmatter:
overlapping edits` in either order; `? status` failed with `cannot locate the
key's colon`; a comment-only `grove.yaml` change at the `compare` step
published, and `create.New` created a record after the same change.

Acceptance: (1) `TestEditKeepsCommentsBeforeLaterLineValues` asserts exact
bytes for key-line and standalone comments, LF and BOM+CRLF, quoted keys and
values, lists at and beyond the key's indentation, a flow mapping, and
unchanged bodies. (2) `TestEditFlowRemovals` asserts exact bytes for the
final two in both orders, absent and existing `updated`, first/middle/last,
every entry, trailing commas, set plus unset, and multi-line mappings with
inline and standalone comments in LF and CRLF;
`TestEditFlowRemovalsExhaustive` checks every subset of four entries, with
and without appends, over six layouts for valid YAML, exact remaining values,
and surviving standalone comments. (3) `TestEditExplicitKeys` covers set and
unset of explicit keys in block and flow form with comments, tags, CRLF, and
multi-line values; the tagged and anchored fixtures stay green with two more
forms. (4) `TestUpdateDetectsChangesDuringPreparation` (comment-only
configuration) and `TestNewRefusesWhenConfigurationBytesChangedAfterLoad`
refuse without publication, leave no temporary file, and keep the
reservation consumed (the next record is W-003).
`TestUpdatePreservesAcceptedForms` repeats the review's reproducers through
`Apply`. (5) Existing W-003 tests pass unchanged.

Independent review (separate reviewer agent, commits `2cc7814` and `dac27fe`,
own export): no P1/P2; base failures reproduced; 46 adversarial `Edit` probes
produced correct bytes or a safe refusal, none a silent wrong result. Fixed:
the always-true key guard inherited from W-003 (`dedac73`). Accepted as is: a
stray space before an appended comma in one multi-line layout (valid YAML),
and removal of a same-line comment that follows a removed last entry's
trailing comma.

Limits: a flow separator on a later line than its entry's value refuses
rather than edits; the existing-`updated` and CRLF flow cases are asserted at
`Edit` level, not again through `Apply`; arbitrary editors writing after the
final comparison remain outside the guarantee. Suite results are in the
[combined repair evidence](../../docs/reviews/2026-09-19-repairs-W-006-W-008.md).

## Next

Integrate branch `worktree-W-006-W-008` into main as a separate, explicit
step; this record being done asserts completion on that branch only.
