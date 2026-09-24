---
id: "G-128"
type: review
title: "G-114 independent review of the knowledge rules"
status: current
created: "2026-09-24T04:42:56Z"
updated: "2026-09-24T04:43:14Z"
work: ["G-114"]
examined: "ebfcb7b"
---

## Examined

The combined diff of [G-114](G-114-capture-and-reuse-terms-question.md)
against its base, `git diff 25525cf HEAD -- docs grove/G-056-attempt.md`:
both guides, `docs/record-model.md`, `docs/commands.md` and
[G-056](G-056-attempt.md). The headless `/grove-work G-114` session asked a
separate reviewer subagent (Claude Opus 5.5, read-only, no edits), which
did not write the change. It checked the record's In and Out of scope, the
anchors, G-051 and the settled terms, and the new knowledge check.

Rounds:

- Round 1, before `bf60f18`, by an earlier attempt of this assignment that
  was interrupted before it recorded its findings. Its fixes are `bf60f18`
  and `70ecd0c`. Its findings are known only through those commits.
- Round 2 examined `70ecd0c` and found three consequential and five minor
  findings, fixed in `54cb163`.
- Round 3 examined `54cb163`, found two minor findings and gave the verdict
  "ready". Both were fixed in `ebfcb7b`, which is the `examined` commit. The
  reviewer did not reread those two sentences.

## Findings

Round 2, at `70ecd0c`:

1. Consequential. The resolution rule would record G-118's preference as an
   accepted decision. The Out of scope list forbids that. The rule did not
   separate an answer that decides from one that prefers or defers.
2. Consequential. The shaping and work guides explained why a note in Next
   is not enough by saying it "blocks nothing". That leaned toward the
   "block" answer on the headless bound (G-118, G-122), which the text must
   leave open.
3. Consequential. The work guide said "a concept" where shaping says
   "domain vocabulary". Nearly every candidate introduces some concept, so
   the review check would fire on every review and push agents toward
   writing too many terms.
4. Minor. "Read the terms and decisions that touch it" did not say how to
   find the unlinked ones.
5. Minor. The interactive bullet said "steps 1 and 2", which leaves out the
   commit.
6. Minor. The knowledge check sat between self-review sentences, so it was
   unclear whether a self-check runs it.
7. Minor. The record-model sentence was awkward.
8. Minor. The guides digest in G-114's Next was stale.

Round 3, at `54cb163`: the reviewer judged findings 1 to 7 resolved. It
accepted the author's reply to finding 2: the work guide keeps `blocks` for
an interactive question, because a missing decision there stops the unit,
and the G-118 bound concerns only headless shaping proposals. It raised two
new minor findings:

9. Minor. "A part left open stays in the question" leaves an open choice
   inside a resolved question, where nothing shows it as open.
10. Minor. In the interactive bullet, the thing created was "a choice"
    rather than a question.

Knowledge check on the candidate: no domain concept without a term, and no
contradiction of a settled term or of
[G-051](G-051-typed-knowledge-records.md). The only conflict with an open
choice was finding 2. The threshold is guide text that is cheap to reverse,
so by its own rule it needs no decision record.

## Disposition

- Findings 1 to 7 were fixed in `54cb163`.
- Finding 8 is fixed in G-114's Evidence, which records the final guides
  digest.
- Findings 9 and 10 were fixed in `ebfcb7b` with the reviewer's suggested
  wording.

That keeps the fixes within three rounds for this gate. Nothing is open.

Finding 9 describes G-118 as it stands: the owner's preference on the
headless bound, left open, sits only in a resolved question's body and in
G-122. This review does not create a question for it. That belongs to
whoever acts on G-118 or G-122, and G-114 keeps the bound out of scope.
