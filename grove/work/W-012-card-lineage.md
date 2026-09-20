---
id: "W-012"
type: work
title: "Show a work item's lineage from Git history in its card"
status: active
created: "2026-09-20T04:36:57Z"
updated: "2026-09-20T14:55:55Z"
kind: feature
priority: 2
size: small
relates_to: ["W-009", "W-013", "Q-001"]
---

## Outcome

Opening a card shows how that work item changed over time: when it was first
proposed, each commit that touched its record, and the status the record held
at that commit. The owner asked for this on 2026-09-19 after the first W-009
demo: the list of versions across branches did not say what had happened to an
item, and lineage was what they expected to find useful. Lineage is selected
direction; everything below is proposed design.

Git already holds the lineage, so nothing new is stored. For W-006 on
2026-09-19, `git log --follow` on its file plus the `status:` line at each
commit gave:

```text
e147e43 09-19 16:07  done      docs: close W-006 with evidence...
fd20223 09-19 14:48  proposed  docs: add reusable execution handoffs...
400e366 09-19 14:32  proposed  docs: review integrated CLI and shape...
```

Proposed: a History section in the card's detail pane for the focused version,
newest first, read from that version's branch or checkout HEAD. Read it when a
card is opened, not for every record during the board load, so the board's load
time (W-013) does not grow with history. Uncommitted live edits appear as a
first "uncommitted" row using the change classification the board already has.

## Constraints

Reads only, through Git, with the board's cancellation and one-read-at-a-time
rule. Commit subjects and author text are file-like input: escape them as the
board escapes record text. Follow renames of the record's file. History on one
branch says nothing about integration into another; do not imply it. No new
record fields and no stored history.

## Acceptance

1. A fixture record created, edited, renamed, and moved through statuses across
   commits shows each commit with its date, subject, and the status at that
   commit, newest first, for the branch of the focused version.
2. Two branches with different histories for one ID show different lineages
   depending on the focused version.
3. Control characters in a commit subject are displayed escaped.
4. Opening a card cancels cleanly and never blocks keys; the board load makes
   no history reads.
5. Targeted, full, and race suites, vet, and formatting pass; the owner judges
   the view in a demo.

## Next

On 2026-09-20 the owner selected "beside, history first": the card screen and
its version list stay, and History leads the details pane, following the
focused row (the board's checkout while the ID header has focus). That is
selected direction; the read and the model are in the
[plan](../../docs/plans/W-012-card-lineage.md).

Checkpoint: `/grove-work W-012`, branch `worktree-W-012` in
`.claude/worktrees/W-012`, base `a0fd23a`. Plan committed; implementation
follows its tasks in order.
