---
id: "G-198"
type: review
title: "Review of G-180 policy-driven integration"
status: current
created: "2026-09-26T03:45:33Z"
updated: "2026-09-26T03:45:52Z"
work: ["G-180"]
examined: "8eb741f1e5c82be5b1b3b679243f6ca57eb4bb8b"
---

## Examined

[G-180](G-180-policy-driven-integration.md) on `worktree-G-180`: the
change from main `fd7744e` to `8eb741f`, against the record's acceptance,
plan [G-196](G-196-plan-for-g-180-policy-integration.md), decision
[G-182](G-182-standing-policy-delegation.md) and the repository's
instructions. Three rounds, each by a fresh `grove-reviewer` agent,
read-only; each ran the tests it needed and exercised the built binary or
the package tests on disposable repositories.

- Round 1 examined `2edcb1a`.
- Round 2 examined `6aba90e`.
- Round 3 examined `8eb741f`, the code of the candidate. The candidate
  adds only this record and G-180's evidence.

## Findings

Round 1 (`2edcb1a`), six findings:

1. **Consequential.** A sweep with several clean candidates integrated only
   the first: each later one was verified against the target as planned,
   approved, and then refused by `integrate`'s `Expect`, leaving a delegated
   verdict whose verification did not cover the moved target. Fixed in
   `6aba90e`: `approve` predicts again and verifies against the current
   target; `TestSweepIntegratesSeveralCandidatesInOneSweep` (the reviewer
   confirmed it fails without the fix).
2. **Consequential.** An older clean review outvoted a newer current review
   with open findings. Fixed: every current review covering the candidate
   must close `Open findings: none`; `TestSweepHeedsEveryReviewOfTheCandidate`.
3. **Consequential.** A delegated resolution took `run:` defaults from the
   candidate's own `grove.yaml`, so the change under judgment chose the
   permission mode of an unattended spend. Fixed in `6aba90e` for the mode
   and budget, and in `8eb741f` for model and effort (a second defaulting in
   `prepare`); the resolve test gives the branch `bypassPermissions`, a
   model and an effort and asserts none reach the launch.
4. `--dry-run` planned a resolution that `resolve` would refuse for want of
   a permission mode, reserving its budget. Fixed: the plan waits and names
   the mode.
5. Acceptance 5's real-provider trial was not run. Disputed and agreed in
   rounds 2 and 3: the record bounds it separately, before use in this
   repository; it stays open for the owner.
6. The record model still said Review "awaits human judgment" beside the
   new delegate sentence. Reworded; round 2 found the rewording contradicted
   the settled term [G-058](G-058-review.md), and it now reads "awaits
   human judgment: the owner's own, or given in advance as a standing
   `policy:`".

Round 2 (`6aba90e`): findings 1, 2, 4 and 6 resolved; finding 3 partly
(model and effort still leaked through `prepare`); one knowledge finding
(the Review wording against G-058). Both fixed in `8eb741f`.

Round 3 (`8eb741f`): both round-2 findings resolved, the first verified by
reverting the fix in a copy and watching the test fail. No new finding of
consequence. Notes for the owner, not findings: `verify` commands have no
timeout; files a candidate adds under `.claude/` or `.agents/` still reach
its resolution attempt, which runs in the branch's checkout, since `never`
guards approval, not resolution; a few lines of `docs/record-model.md` are
over 80 columns.

Knowledge: the term [G-197](G-197-policy.md) (proposed) defines the
policy and a delegated act. The work contradicts no settled term or
accepted decision: approval by a delegate is inside
[G-059](G-059-approval.md), and a review stays evidence the policy reads.

## Disposition

Every consequential finding was fixed and re-reviewed. Three rounds, the
cap, were used; the last found nothing consequential. This review is
evidence, not approval.

Open findings: none
