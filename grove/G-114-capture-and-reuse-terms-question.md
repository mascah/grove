---
id: "G-114"
type: work
title: "Capture and reuse terms, questions and decisions across shaping, work and review"
status: active
created: "2026-09-23T19:44:56Z"
updated: "2026-09-24T04:31:44Z"
relates_to: ["G-037", "G-051", "G-056", "G-107", "G-108", "G-118", "G-122", "G-125"]
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

Refined 2026-09-24 in an interactive shaping session after
[G-108](G-108-workflow-evals.md) closed: the sequencing this record waited on
is met, and the questions written since narrow the diagnosis (below).

## Constraints

Observed at main `f81f7e9`, first observed at `a54513f`; both guides are
byte-identical between the two and `grove version` prints guides digest
`3f5487904c61`, the digest G-108's baseline ran at:

- The [shaping guide](../docs/work-shaping.md) says settled vocabulary belongs
  in a term record and when a question or decision may be written. The
  [work guide](../docs/work-execution.md) names decision records only under
  rejection and terms nowhere; its interactive missing-decision path asks in
  chat and relies on the checkpoint's pending judgments, while only the
  headless path persists a question. Review (step 6) examines acceptance and
  evidence, not knowledge.
- Every term (G-054 to G-062) was created in one batch under G-037 on
  2026-09-21 and none since; no work record after G-098 links a term, other
  than this one. No decision has been recorded since G-101.
- Questions are now created when a headless attempt meets the choice.
  [G-117](G-117-which-attempts-list-and-detail-l.md),
  [G-118](G-118-what-mandate-should-the-g-108-pa.md) and
  [G-121](G-121-how-do-the-g-108-eval-runs-log-i.md) were each written by a
  headless work attempt following the work guide's missing-decision path,
  with options, evidence, recommendation and `blocks`, and each was resolved
  by the owner writing the answer into the body. On the eval fixture,
  [G-122](G-122-g-108-baseline-runs-the-missing.md) reports 5 of 5 headless
  shaping runs persisting a blocking question and 5 of 5 companion runs
  asking nothing. The gap is after resolution and in interactive sessions,
  not in the question mechanics: none of the three resolved questions links
  a decision or term; G-118's answer carries a product preference about the
  headless bound that no decision record holds and G-122 reports as still
  open; the choice the G-108 session took about synced account content in
  the clean configuration lives only in G-108's Evidence; and the owner's
  interactive choices in G-108's Constraints sit in the work record, which
  the shaping guide allows for routine ones and which nothing distinguishes
  from consequential ones.
- [G-056](G-056-attempt.md) still says Grove has no attempt record and that
  an attempt is visible only as a branch, a worktree and a checkpoint;
  [attempt.go](../internal/attempt/attempt.go) has written attempt files
  since G-045, and [docs/commands.md](../docs/commands.md#attempts) is where
  attempts are documented. The definition mixed an implementation observation
  into a settled meaning, and that part went stale unnoticed.
- `grove context` lists only what the selected work names in `depends_on`,
  `blocks`, `relates_to` or `members`, or links from its body
  ([context.go](../internal/handoff/context.go)). A term or decision nothing
  links is invisible to the next agent. G-122 adds that the eval pair cannot
  see the listing at all, since the fixture's only knowledge is its brief.
- The [record model](../docs/record-model.md) already gives every type it
  needs: terms with meaning and relationships, questions with `blocks`,
  decisions with authority, alternatives and reconsideration conditions, and
  `relates_to` on any record. A resolved question keeps its answer in its
  body or links the durable decision.
  [G-125](G-125-answer-a-blocking-question-from.md) proposes answering a
  question from the board under an `## Answer` heading; the resolver's step
  this record adds must fit that path, and G-125 must not presume this
  record's wording.
- The owner's open choice about the headless bound, whether a shippable
  proposal with a non-blocking question is preferred to blocking (G-118
  item 5, G-122 Disposition), edits the same "Missing human choice"
  paragraph this record's shaping change touches. This record does not make
  that choice and its text must fit either answer: it says when a question
  is persisted and what it carries and links, never whether the proposal
  blocks. If the owner makes the choice, it is a separate guide edit paired
  with a change to the eval check `question-blocks-proposal`.
- G-108 is done at `60c9537`. Its case pair reruns for about $3 and ten
  minutes and is the regression check for a guide edit; its companion case
  is where a guide change that makes agents write more questions or terms
  would show as over-asking. The end-to-end knowledge sequence in G-108's
  Next is the eval of this procedure and is follow-on work, not this
  record's.

In scope, all as guide text unless stated: before introducing or changing a
concept, read the terms and decisions that touch it and name conflicts and
synonyms; capture a definition in a term record when it settles, as meaning,
relationships, boundaries and misleading alternatives, with implementation
state and progress kept elsewhere; persist a consequential open question
before any wait or handoff in either interaction mode, with the choice,
evidence, recommendation, who can answer and what it blocks; on resolution,
keep the question, and when the answer is consequential by the threshold
below record it as a decision the question links, so an answer does not live
only in a resolved question's body; record a
decision only under a selective threshold (proposed: reversal cost, reasoning
a future reader would otherwise lack, and real alternatives), short, with
authority attributed as the shaping guide already requires; add to review
whether the candidate introduces a concept, contradicts a settled term,
depends on an unanswered choice, or implements a consequential decision
without its rationale, reported as findings for the author to reconcile; and
require work to link the terms and decisions that govern it so `context`
lists them. Repairing G-056 is the first exercised instance of the term rule.

Out of scope: `context` supplying terms or decisions automatically, a
generated glossary or index, a new record type or field, any parallel
glossary file or separately numbered decision tree, the headless bound's
block-or-surface choice, and recording G-118's preference as a decision
before the owner makes it. Reconsider the retrieval change only if evidence
shows links present and unread rather than absent.

## Acceptance

1. Both guides carry the procedure once, in the step where each activity
   meets it, and the review step names the knowledge check; the adapters are
   unchanged, `grove guide shape` and `grove guide work` print the change,
   and the Evidence records the new guides digest `grove version` prints, so
   a G-108 rerun is comparable with its baseline at `3f5487904c61`.
2. G-056 states the meaning without implementation state, and the observation
   it dropped lives where attempts are documented, with a link.
3. One real assignment after the change is walked through the procedure and
   its outcome reported honestly: which terms, questions and decisions it
   read, created or linked, and which rule it found unclear or unnecessary.
   Absence of a new record is a valid result when nothing settled.
4. `grove check` passes and every new link resolves. No product change beyond
   guide text and record bodies.

## Next

Checkpoint, 2026-09-24, headless `/grove-work G-114` on `worktree-G-114`,
based on main `25525cf` (record revision `f4b81aa`). No plan needed: the
In scope list names each rule and both guides, and the change is guide text
and two record bodies. Done: guide edits `d2379f0` (guides digest
`dff386aa5bb2`), G-056 repair `08a30b6`, `relates_to` gains G-051, which
governs term records. Pending: an independent review of the diff, then the
G-108 pair rerun at the reviewed guides, then Evidence and Review.

Owner decisions, 2026-09-24, in the shaping session that refined this
record:

- The G-108 pair reruns after the guide edit, as this record's own
  mandate: `python3 evals/run.py run --runs 5 --budget 5 --model
  claude-opus-5-5 --permission-mode auto --config-dir
  ~/.cache/grove-evals/claude --out DIR`, with `DIR` outside every checkout,
  about $3 and a $50 cap. Read the companion case for over-asking and the
  missing-choice case against G-122, and report both in this record's
  Evidence; G-108's Next stays untouched.
- G-125 walks acceptance 3, since it changes how a question is answered and
  so meets the resolution rule directly. Assign it after this record lands.

The headless bound's block-or-surface choice (G-118, G-122) stays open and
is not this record's; G-108's Next lists the knowledge-sequence eval case as
the follow-on once this lands.
