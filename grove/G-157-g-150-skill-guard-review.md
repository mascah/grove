---
id: "G-157"
type: review
title: "G-150 skill-guard review"
status: current
created: "2026-09-25T19:35:52Z"
updated: "2026-09-25T19:36:07Z"
work: ["G-150"]
examined: "b322825"
---

## Examined

[G-150](G-150-launch-attempts-only-where-the-w.md) on `worktree-G-150`,
base `28aaf95`, two independent rounds by fresh `grove-reviewer` subagents
on 2026-09-25: round 1 at `cd3a1c5` (the guard, docs and tests), round 2
at `b322825` (the fix commit `cd3a1c5..b322825` against round 1's
findings), the `examined` commit. Both ran fixtures in disposable `mktemp`
directories with a built binary; neither edited the checkout.

## Findings

1. Medium, round 1: the on-disk skill check ran before the branch's own
   review and blocking-question refusals, so a branch whose candidate is
   in review but lacks the skill was told to commit `init`'s files to it
   instead of to judge the candidate. Reproduced by the reviewer.
2. Low to medium, round 1: the board's `R` collects `Start`'s report
   facts and shows them after `Start` returns, so the reviewer warning is
   drawn after the owner and provider have started; `grove run` prints it
   before the owner starts. Neither path pauses, so the warning cannot
   prevent spend in either.
3. Low, round 1: `internal/attempt` is over five seconds, which predates
   the change (the reviewer measured 4.93 s and 5.48 s at `28aaf95`,
   5.18 s and 5.80 s at `cd3a1c5`).
4. Low, round 1: the new-branch check read symbolic `HEAD` while the
   branch is created from the `head` SHA read earlier.
5. Low, round 1: the README's `git add grove.yaml grove .claude .agents`
   assumes the checkout's top.

Round 2 found nothing consequential and no regression. It confirmed that
the `TestRefusals` regression for 1 fails when the check is moved back
(mutation run in a `git archive b322825` copy), and repeated the knowledge
check: no new concept, no conflict with G-101, G-134 or the Attempt term.

## Disposition

1, 4 and 5 fixed in `b322825` and verified in round 2. 2 is left open for
the owner: the record's proposed design has `R` show the messages "as it
shows the other launch messages", and drawing them before the owner starts
would need facts streamed into the board mid-launch, a wider change. It
limits acceptance 2's "before the provider starts" to `grove run`. 3 is
reported as an inherited limit, not caused here.
