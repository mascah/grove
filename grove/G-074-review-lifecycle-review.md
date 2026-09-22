---
id: "G-074"
type: review
title: "G-038 review lifecycle review"
status: current
created: "2026-09-22T04:33:20Z"
updated: "2026-09-22T04:33:40Z"
work: ["G-038"]
examined: "f9f6c0d"
---

## Examined

Independent review of [G-038](G-038-review-lifecycle.md) against its
acceptance and the plan [G-073](G-073-review-lifecycle-plan.md), by a
reviewer session with no write access to the checkout, on 2026-09-22.
Round 1 examined `f9f6c0d` on `worktree-G-038` (diff from main `a28a24b`);
round 2 examined the fix commit `12b7752`. The reviewer mutation-tested both
new rules in a scratch copy: forcing `integrated` to return nil fails
`TestUpdateDoneMeansAnIntegratedCandidate`, and removing the
required-in-review branch fails `TestReviewLifecycleAndCandidate`.

## Findings

Round 1, consequential:

1. `update` refused `done` only where HEAD lacked the candidate, so on the
   work branch, where the candidate is an ancestor too, `done` was accepted;
   the record model, guide, README, AGENTS.md, the plan and the code comment
   all said it was refused there. Reproduced in a disposable project.
2. `context` printed "done without a candidate is not integration" above
   Requirements lines that carried no candidate, so a reader could not tell
   the two apart.
3. The G-038 record had no Evidence and a stale Next, so acceptance 2 was not
   met as the branch stood.
4. The plan said `candidate` is refused while `proposed`; the code allows it.
5. The README said "a `review` record" where it meant a work record in
   Review status, the distinction the term G-058 draws in the same commit.

Round 1, minor:

6. A reopened record could be closed again on its old candidate; only the
   guide's `git diff --stat CANDIDATE TIP` check catches it, and the record
   model's not-enforced list did not say so.
7. `--unset candidate` on a done record gave the "set candidate=COMMIT" advice.
8. An unknown commit and a genuine non-ancestor shared one message.
9. No test proved a letters-only candidate stays a quoted string.
10. The narrow tab names were a hand-sized array indexed over the statuses.
11. The plan promised a connected terminal check at 100 and 80 columns; the
    diff widened the connected test to 160 instead, and five columns at 100
    cells truncate titles sooner.
12. The brief still says Done "will" mean integrated and that G-038 "must"
    migrate it.

Acceptance as the reviewer read it at `f9f6c0d`: 1, 4 and 5 met; 3 met in
substance with the enforcement gap in finding 6 stated as guidance; 2 not met
(finding 3); 6 partly, pending the evidence transcript.

## Disposition

Fixed in `12b7752`: finding 1 by rewording the six places to what the check
enforces, that a checkout without the code cannot close the work, with the
target rule named as the guide's and AGENTS.md's, plus a test that the
on-branch write is accepted; 2 by a `candidate` field on `Requirement`,
rendered as `done, candidate COMMIT` or `done, no candidate`; 4 and 5 by
text; 6 by a line in the not-enforced list; 7 and 8 by separate messages
(`a record that stays done cannot lose its candidate`; Git's exit 1 is reported as
not an ancestor, anything else as could not be checked), each with a test;
9 by a quoting test; 10 by an array sized from the status count.

Finding 3 is the evidence and Next written into G-038 with the handoff.
Finding 11: the 100- and 80-column checks were run on a pseudo-terminal and
are in G-038's evidence; the truncation is reported there for the owner, and
G-043 owns the visual redesign. Finding 12 is left for the owner: AGENTS.md
reserves the brief for direction changes, not progress.

Round 2, on `12b7752`, consequential: (13) the round-1 fix for finding 2
printed `done, candidate` for a prerequisite with no candidate, and no test
rendered a Requirements line; (14) the record model said "a done record's
candidate cannot be removed", but reopening with `--unset candidate` in the
same call is allowed, so the enforced rule is that a record that stays done
cannot lose it. Minor: (15) closing from review while unsetting the candidate
got the removal message instead of the advice to set one; (16) the JSON
comment on `Requirement.Candidate` said "merged" though a review
prerequisite's candidate is only offered; (17) a Git warning beside exit 1
(a branch named like a hex commit) would get the "could not be checked"
wording, cosmetic and still a refusal; (18) the tab-name array and its
comment; (19) a test scoping hack. Round 2 confirmed every round-1 fix
closed and the fixture sound.

Fixed in the following commit: 13 by an explicit branch with
`TestRequirementsNameADoneCandidateOrItsAbsence`; 14 by the sentence; 15 by
guarding on the prior status, with a test; 16 and 18 by the comments; 19 by
a plain block. 17 left as cosmetic. ROUND3_PLACEHOLDER
