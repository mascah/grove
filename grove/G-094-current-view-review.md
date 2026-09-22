---
id: "G-094"
type: review
title: "Review of G-042 current view"
status: current
created: "2026-09-22T22:27:11Z"
updated: "2026-09-22T22:27:26Z"
work: ["G-042"]
examined: "96900e0"
---

## Examined

One independent reader, a Claude reviewer subagent of the implementing
session with no edit rights over the branch, reviewed branch `worktree-G-042`
against base `939d090` in three rounds, reading the code and running scratch
tests in a copy of the worktree.

- Round 1 examined `383e59d`: the projection, merge bases, caches, policy,
  TUI, and documentation, against [G-042](G-042-current-view.md) and
  [its plan](G-093-current-view-plan.md).
- Round 2 examined `e1606af`, the fixes for round 1.
- Round 3 examined `96900e0`, the fix for round 2. No defects remained.

`1c328cf` (a test made parallel) and `84115c4` (usage text) followed without
review; neither changes behavior.

## Findings

Round 1, all confirmed by running or by reading:

1. A cycle of older observations, from a revert carried across merges, left a
   record with no current state, so its card vanished.
2. A read failing during the projection, as a cancelled load's does,
   dereferenced nil and would end the board with a panic.
3. A branch's committed deletion row had no path, so its history read the
   whole branch's log and was labelled as a checkout's uncommitted deletion.
4. A current state that was only a checkout's uncommitted deletion lacked the
   `uncommitted` tag.
5. Two checkouts on one unborn HEAD gave a spurious "could not be ordered"
   note.

It also noted that the measurements were not yet recorded, and some stale
wording: the `Change` comment, the divergence text, and a progress sentence
in the brief. The versions package's short tests ran at 4.5 to 5.9 s,
against the 5 s budget.

Round 2: finding 1's fix covered only the case where every observation was
older. A cycle stayed hidden, with no note, while an unrelated branch that
could not be ordered against it stayed current.

Round 3: none. The reviewer checked that the loop trigger cannot miss a
sink cycle and that a cycle something else supersedes stays older, and ran
both from two invoking checkouts.

## Disposition

- 1 and its round-2 remainder: fixed in `e1606af` and `96900e0`. When a
  node's chain of reasons loops, the whole relation decides, and a node is
  current when everything newer than it, through any chain, is also older
  than it; a note says so. `TestCurrentViewCycle` includes the unrelated
  branch and fails at `e1606af`.
- 2: `e1606af`, `TestCurrentViewCancelled` (panicked before).
- 3 and 4: `e1606af`, with `TestCurrentViewBoard` now using the path-less
  row the projection produces, and an uncommitted deletion (W-005).
- 5: `e1606af`, `TestBaseOfUnborn`.
- Wording: fixed in `e1606af`. Measurements are in G-042's Evidence.
- Test budget: `1c328cf` runs the merge-base test in parallel; the package's
  short tests then took 4.0 to 4.3 s alone.
