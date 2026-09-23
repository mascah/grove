---
id: "G-114"
type: work
title: "Capture and reuse terms, questions and decisions across shaping, work and review"
status: proposed
created: "2026-09-23T19:44:56Z"
updated: "2026-09-23T19:45:35Z"
relates_to: ["G-037", "G-056", "G-107", "G-108"]
---

## Outcome

Shaping, work and review share one knowledge-capture procedure, so that a
term settled in conversation, a consequential choice left open in execution,
and the answer that resolves it each reach the record that owns it and can be
retrieved by the next session instead of contradicted.

Owner intent, 2026-09-23: the owner is concerned that Grove does not reliably
save terms and questions during shaping, work or review, and favours the
pattern of Matt Pocock's grilling and domain-modeling skills, where
definitions are challenged against scenarios and code and recorded when they
settle, and consequential trade-offs become selective architecture decision
records. An external source-based assessment the owner read the same day
found the storage adequate and the discipline missing; this record follows
its diagnosis. Its formats and thresholds are proposed design.

## Constraints

Observed at main `a54513f`:

- The [shaping guide](../docs/work-shaping.md) says settled vocabulary belongs
  in a term record and when a question or decision may be written. The
  [work guide](../docs/work-execution.md) names decision records only under
  rejection and terms nowhere; its interactive missing-decision path asks in
  chat and relies on the checkpoint's pending judgments, while only the
  headless path persists a question. Review (step 6) examines acceptance and
  evidence, not knowledge.
- Every term (G-054 to G-062) was created in one batch under G-037 on
  2026-09-21 and none since. One question has ever existed (G-002). No
  decision has been recorded since G-101. No work record after G-098 links a
  term. Several attempts, reviews and the G-107 reconciliation ran in that
  time.
- [G-056](G-056-attempt.md) still says Grove has no attempt record and that
  an attempt is visible only as a branch, a worktree and a checkpoint;
  [attempt.go](../internal/attempt/attempt.go) has written attempt files
  since G-045. The definition mixed an implementation observation into a
  settled meaning, and that part went stale unnoticed.
- `grove context` lists only what the selected work names in `depends_on`,
  `blocks`, `relates_to` or `members`, or links from its body
  ([context.go](../internal/handoff/context.go)). A term or decision nothing
  links is invisible to the next agent.
- The [record model](../docs/record-model.md) already gives every type it
  needs: terms with meaning and relationships, questions with `blocks`,
  decisions with authority, alternatives and reconsideration conditions, and
  `relates_to` on any record.

In scope, all as guide text unless stated: before introducing or changing a
concept, read the terms and decisions that touch it and name conflicts and
synonyms; capture a definition in a term record when it settles, as meaning,
relationships, boundaries and misleading alternatives, with implementation
state and progress kept elsewhere; persist a consequential open question
before any wait or handoff in either interaction mode, with the choice,
evidence, recommendation, who can answer and what it blocks, and keep it
after resolution linked to the term or decision that answered it; record a
decision only under a selective threshold (proposed: reversal cost, reasoning
a future reader would otherwise lack, and real alternatives), short, with
authority attributed as the shaping guide already requires; add to review
whether the candidate introduces a concept, contradicts a settled term,
depends on an unanswered choice, or implements a consequential decision
without its rationale, reported as findings for the author to reconcile; and
require work to link the terms and decisions that govern it so `context`
lists them. Repairing G-056 is the first exercised instance of the term rule.

Out of scope: `context` supplying terms or decisions automatically, a
generated glossary or index, a new record type or field, and any parallel
glossary file or separately numbered decision tree. Reconsider the retrieval
change only if evidence shows links present and unread rather than absent.

## Acceptance

1. Both guides carry the procedure once, in the step where each activity
   meets it, and the review step names the knowledge check; the adapters are
   unchanged and `grove guide shape` and `grove guide work` print the change.
2. G-056 states the meaning without implementation state, and the observation
   it dropped lives where attempts are documented, with a link.
3. One real assignment after the change is walked through the procedure and
   its outcome reported honestly: which terms, questions and decisions it
   read, created or linked, and which rule it found unclear or unnecessary.
   Absence of a new record is a valid result when nothing settled.
4. `grove check` passes and every new link resolves. No product change beyond
   guide text and record bodies.

## Next

The owner can commit this proposal and assign it with `$grove-work G-114`.
The guide edits change the digest [G-108](G-108-workflow-evals.md) pins:
land them after G-108's baseline has run, or accept that its first rerun then
differs by two levers. G-108's Next lists the end-to-end knowledge sequence
as a follow-on evaluation case for this procedure.
