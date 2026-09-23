---
id: "G-120"
type: review
title: "G-109 Attempts redesign review"
status: current
created: "2026-09-23T20:34:58Z"
updated: "2026-09-23T20:35:13Z"
work: ["G-109"]
examined: "611f7ec"
---

## Examined

An independent reviewer subagent (read-only, headless attempt of
2026-09-23) reviewed [G-109](G-109-attempts-usability.md) in two rounds on
branch `worktree-G-109`. Round 1 covered `0315df9`, the implementation,
against its parent. Round 2 covered `611f7ec`, the fixes, against `0315df9`.
Both rounds checked the code against G-109's constraints,
[G-117](G-117-which-attempts-list-and-detail-l.md)'s answer and the section
"Revised after G-117's answer" in
[G-116](G-116-g-109-attempts-layouts-observed.md). The reviewer ran the
short tests and vet, probed widths 40, 41, 60, 99, 100 and 101 with `d` on
and off, and ran `ReadActivity` from the commit on all seven real
`events.jsonl` logs, comparing with counts from Python.

## Findings

Round 1 found no problem with terminal safety or widths, bounded reads, the
grouping rules apart from item 2, or the docs. It found:

1. Should-fix: Subagents was counted too high. Claude repeats `task_started`
   each time a background subagent is resumed, so G-108 showed ≥3 against
   1 spawned.
2. Should-fix: a work in review whose current candidate differed from the
   attempt's still read "judge candidate", with the attempt's commit.
3. Should-fix: Turns took the largest `num_turns` when it should sum them.
   Each result event of a resumed session counts only its own query, so
   G-107 showed ≥47 against a true 64.
4. Should-fix: Model read "not reported yet" when the window was cut, even
   though `result.json` held the model.
5. Nit: a failed read showed zeros as figures.
6. Nit: "spend known at the end" appeared for attempts that had ended
   without a result event.
7. Nit: every settled attempt got ✓.
8. Nit: the context window was the largest of the run's models.

Round 2 found all eight fixed and no regression. On every real log, Turns
and Subagents now equal Claude's own totals: G-046 122/1, G-107 74/2 and
64/3, G-108 63/1, G-105 39/1, G-109 32/0. It raised two optional nits:

9. Nit: a log small enough to fit the window, but with more than 200
   activity rows, still showed `≥` on every count.
10. Nit: the State line read "is review" instead of "is in review".

## Disposition

- Fixed in `611f7ec`, each with a regression test: 1 (distinct `task_id`),
  2 (W-108 in `TestAttemptStandings`), 3 (summed results and `Ended`, shown
  as `≈N` until a result ends the run), 4, 5 and 6
  (`TestAttemptScreenHonesty`), 7 and 8.
- Fixed in `feb6b42` after round 2, without a third review, both small and
  covered by tests:
  - 9: a separate `Activity.Dropped` flag. `TestReadActivityBounded` now
    asserts it is set without `Cut`.
  - 10: the wording, with its test string updated.
- No finding is open.
