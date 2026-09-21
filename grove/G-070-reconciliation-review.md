---
id: "G-070"
type: review
title: "Independent review of the G-052 reconciliation"
status: current
created: "2026-09-21T22:37:40Z"
updated: "2026-09-21T22:37:40Z"
work: ["G-052"]
examined: "dd3a6f5"
relates_to: ["G-068", "G-069"]
---

## Examined

Two independent read-only reviews by reviewer agents that wrote none of the
work, plus one independent reading of every rewritten ID line during
preparation. Branch `worktree-W-029`, base main `70de539`.

- Preparation: all 860 lines where the rehearsal rewrote a typed ID, each read
  against its old context. About 50 lines named fixture, trial-clone or
  predecessor records, and four more were verbatim quotations; all became
  exceptions listed in [G-069](G-069-migration-map.md).
- Round 1: the content migration, `8bf691b` against `1854d8c`.
- Round 2, final: the combined revision `dd3a6f5`, which is this record's
  `examined`. The commit after it holds only the round-2 fixes below.

## Findings

Round 1 found no corruption, lost record, broken link or altered
timestamp or status: all 69 old/new pairs differed only by mechanical edits
and every relationship field matched the mapping. Consequential: (1) the
guides' example IDs had been hand-set to `G-012 G-014`, a plan and a work
record, so the guide's own `context` example failed; (2) the shaping evidence
review's `work` omitted the shaping work record, hiding it from that record's
context; (3) three ID ranges ("W-001 through W-005") became false once mapped,
because work IDs are no longer contiguous. Minor: the repairs review's `work`
omitted the middle of a filename range; two typed schema examples in the
record model and the `W-…`/`W-ID` authoring placeholders; two `convert`
rehearsal sentences made anachronistic; G-069 did not account for typed IDs
inside paths. Notes: the 22 converted documents carry no `created`, so `list`
puts them last; one anchor into the brief was already broken at the base.

Round 2 found no behaviour defect in the schema 1/2 removal and confirmed
round 1's fixes, acceptance 1, 2, 3 and 5, and the reintegration and
old-schema-branch behaviour in its own throwaway clone. Consequential: (1)
deleting the mixed-schema CLI test file left `convert` and `page` without any
test through `cli.Run`, shown by a surviving mutation of the changed argument
check; (2) the allocation floor scan lost its only nested-record case. Minor:
the board's new deleted-everywhere branch was untested; two example typed IDs
unaccounted for in G-069; fixtures and comments that still spell the old
layout. It also noted that board and history evidence on migrated records was
not yet in a record.

## Disposition

All consequential and minor findings are fixed: round 1 in `c08d213`, round 2
in the commit carrying this record (a trimmed `internal/cli/flexible_test.go`,
nested paths in `TestAllocateFloorsFromRefsAndWorktrees`, a deleted-everywhere
group in `TestBoardFollowsTypeNotIDOrPlacement` that a mutation of `isWork`
fails, one clause in G-069, two stale comments). Left as they are: `created`
on converted documents, because `convert` invents no dates and `update`
refuses the field, with the document dates kept in G-069; old-layout paths in
test fixtures, which recursive discovery makes legal; and the guides' "where
the schema has" hedges, which keep the guides portable. The round-2 fixes were
not reviewed again: they add tests and prose and change no shipped code. This
record is evidence, not approval.
