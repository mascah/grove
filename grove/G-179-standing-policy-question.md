---
id: "G-179"
type: question
title: "What may a standing owner policy do to a candidate without a per-item human act?"
status: open
created: "2026-09-25T21:39:27Z"
updated: "2026-09-25T21:41:23Z"
blocks: ["G-180"]
relates_to: ["G-058", "G-059", "G-060", "G-142", "G-162", "G-163", "G-177", "G-178"]
---

## Question

When a candidate in review conflicts with the target, or when its
independent review finds nothing consequential, what may Grove do on the
owner's standing instruction alone, with no per-candidate human act: start
a resolution attempt, record approval, integrate? Under what conditions?

Owner intent, conversation 2026-09-25: the owner leans toward automation.
A predicted conflict should be resolvable by an agent, manually or
automatically, and the owner is starting to consider a review after each
attempt that evaluates whether a candidate is safe to merge automatically
or truly needs human approval. Nothing was decided. The answer blocks
[G-180](G-180-policy-driven-integration.md), not the read-only
[G-177](G-177-merge-prediction.md) or the human-triggered
[G-178](G-178-candidate-target-update.md).

## Evidence and options

- The brief's selected foundations put judgment "in instructions and
  attributable human or delegated decisions". [Approval](G-059-approval.md)
  is "a person with authority, or someone they delegated to, accepting one
  specific candidate", and a [review](G-058-review.md) is evidence, never a
  verdict by itself. Delegation is already inside the vocabulary; what no
  record selects is a standing delegation, and the record model's work
  lifecycle says Review status means a candidate awaits human judgment.
  The repository policy says not to assume software can replace judgment or
  prove acceptance.
- [G-045](G-045-durable-attempt.md) and
  [G-162](G-162-bounded-work-selection.md) keep automatic merge out of their
  scope, and [G-163](G-163-selected-work-review-boundary.md)'s
  recommendation excludes automatic merges within a chain. Those are scope
  bounds, not a decision against it.
- Spend is explicit today: every attempt takes a budget
  ([G-140](G-140-default-an-attempt-s-budget-mode.md)), and the owner treated spend
  authorization as an explicit act ([G-141](G-141-never-run-gpt-6-astra-unless-the.md),
  [G-142](G-142-keep-a-blank-mandate-answer-from.md)). An automatic
  resolution attempt spends without a per-item keypress.
- `integrate` runs nothing after the merge. A delegated integration that
  verified nothing could land a semantic conflict on the target. An
  integrated mistake is recoverable by revert, but `done` would have been
  written and the record would need to say so.

Options:

1. **Nothing automatic.** G-177 and G-178 only; every launch, approval and
   integration stays a human act. Simplest; the overnight goal is served
   for implementation but not for integration.
2. **Automatic resolution only.** A predicted or refused conflict starts one
   bounded G-178 attempt under a policy budget, at most once per target
   movement; approval and integration stay human. Spend is delegated,
   judgment is not.
3. **Delegated approval and integration under a written policy.** After the
   independent review, a candidate that meets conditions the owner wrote in
   `grove.yaml` is approved with a verdict attributed to the policy and the
   reviewer's evidence, and integrated; everything else waits for the owner
   as today. Deterministic conditions (verification passed on the merged
   result, changed paths or size within a scope, no open question blocking,
   clean merge) are checked in software; the reviewer's "nothing
   consequential" and "no knowledge finding" are evidence the policy
   consumes, not a verdict the reviewer gives.
4. **Both 2 and 3.**

Assistant recommendation, not a decision: option 4 with a narrow initial
policy the owner writes, because the brief already allows delegated
decisions and every delegated act stays attributable and reversible, while
human approval remains the default for anything the policy does not name.
If the owner would rather see G-177 and G-178 in use first, choose 2 now
and reopen this question for 3.

## Next

The owner answers. Record the answer here and resolve this question with
`grove update`; record a consequential choice as a decision attributed to
the owner, then refine G-180's scope. No answer is implied by creating it.
