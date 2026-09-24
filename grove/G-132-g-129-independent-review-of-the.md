---
id: "G-132"
type: review
title: "G-129 independent review of the headless long-command rule"
status: current
created: "2026-09-24T15:09:34Z"
updated: "2026-09-24T15:09:49Z"
work: ["G-129"]
examined: "4e8cd92"
---

## Examined

The diff of [G-129](G-129-headless-attempts-run-long-comma.md) against main
`c38d914`, which is `docs/work-execution.md` and G-129's own status change.
The headless `/grove-work G-129` session asked a separate reviewer subagent
(Claude Opus 5.5, read-only, no edits). The reviewer did not write the
change. It checked the text against G-129's Outcome and scope, against the
rest of the work guide, and against the settled terms (G-056 in particular)
and G-101. It also checked the anchors by applying GitHub's slug rules by
hand.

Rounds:

- Round 1 examined `88312db`: two consequential findings, four minor ones
  and one nit. They were fixed in `9ad05c5`.
- Round 2 examined `9ad05c5`: every fix held, no new contradiction was
  found, and two nits were raised. Both were fixed in `4e8cd92`, which is
  the `examined` commit. The reviewer did not reread those two changes; the
  author checked them against the nits.

## Findings

Round 1, at `88312db`:

1. Consequential. The checkpoint did not say who can run the command. A
   headless rerun would repeat the same roughly $2 wait that the record set
   out to stop.
2. Consequential. The text said nothing about a foreground command that
   hits the timeout. The eval pair runs close to the ten-minute tool
   maximum, so a timeout is the likely failure.
3. Minor. The sentence above the new rule, "If a command must outlive the
   session, the handoff names its owner, handle…", still seemed to allow a
   headless session to detach a job. Nothing said which of the two rules
   wins.
4. Minor. "Longer than one tool call in the foreground" contradicts itself
   if read literally.
5. Minor. The text did not say that the work stays active or that the
   checkpoint is committed. Step 8's outcome list had no wait on a command.
   Step 8's sentence about a background agent differs in scope from the new
   rule, so it is not a duplicate (nit).
6. Minor. G-129's Next was stale once the record was set active, and
   acceptance 1 and 2 still needed their evidence.
7. Nit. The Inputs link points to all of step 5, not to its last paragraph.

Round 2, at `9ad05c5`:

1. Nit. "Note its partial output" did not say where.
2. Nit. In the Inputs sentence, ", at the end of step 5" could be read as
   applying to both links.

Knowledge check, both rounds: the change introduces no new domain concept.
It does not contradict G-056 or G-101, and it does not depend on the open
choice of whether the attempt owner should report an abandoned job.

## Disposition

- Round 1, findings 1 to 4 and 7: fixed in `9ad05c5`.
- Round 1, finding 5: fixed in `9ad05c5`. The checkpoint links step 7,
  "the work stays active" was added, and step 8 now says "question, blocker
  or command". The step 8 nit was left as it is, by choice.
- Round 1, finding 6: resolved in G-129's Evidence and Next at handoff.
- Round 2, nits 1 and 2: fixed in `4e8cd92`. The checkpoint names "the path
  of any partial output", and the Inputs link reads "(end of step 5)".
