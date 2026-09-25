---
id: "G-180"
type: work
title: "Resolve, approve and integrate candidates under an explicit owner policy"
status: proposed
created: "2026-09-25T21:39:28Z"
updated: "2026-09-25T21:41:23Z"
kind: feature
depends_on: ["G-177", "G-178"]
relates_to: ["G-044", "G-058", "G-059", "G-060", "G-101", "G-134", "G-140", "G-142", "G-162", "G-163", "G-179"]
---

## Outcome

Where the owner has written a standing policy, Grove starts a resolution
attempt for a candidate that conflicts with the target, and approves and
integrates a candidate that meets the policy's conditions after its
independent review, with each act attributed to the policy and its
evidence, so a set of independently implemented items can go from review to
the target while the owner is away, and the owner judges only what the
policy leaves to them.

Owner intent, conversation 2026-09-25, as [G-179](G-179-standing-policy-question.md)
records it. The policy's content is the owner's; this record does not
choose it, and G-179 blocks its execution contract.

## Constraints

### Observed evidence

At main `47852e3`, `approved` is a commit that must equal `candidate`, and
the verdict is a body paragraph (record model, work lifecycle); nothing in
the schema distinguishes the owner's verdict from a delegated one, and
`check` accepts any equal value, which [G-044](G-044-review-integration.md)
notes as the hand-set escape that skips the verdict and tip checks. The
independent reviewer (`.claude/agents/grove-reviewer.md`,
[G-134](G-134-bound-an-attempt-at-its-plan-and.md)) returns findings and
limits and is "never approval"; the [work guide](../docs/work-execution.md)
dispatches it per gate and caps fix rounds at three. `integrate` merges,
then writes done, and runs no verification; approval needs the branch's
clean checkout and integration the target's. Each attempt records its
model, effort, reviewer definition hash and cost
([G-134](G-134-bound-an-attempt-at-its-plan-and.md),
[G-140](G-140-default-an-attempt-s-budget-mode.md)), and nothing resident runs between
attempts ([G-101](G-101-attempt-mechanism.md)). G-179's evidence lists the
authority and spend constraints.

### Proposed design and scope

- A policy in `grove.yaml`, shape chosen in preparation: what it may do
  (resolve, approve, integrate, in whichever combination G-179 selects) and
  its conditions, deterministic wherever possible: verification commands
  that must pass on the merged result, path scope or size of the change, no
  open question blocking the work, a clean merge, the reviewer reporting
  nothing consequential and no knowledge finding, at most one resolution per
  target movement, and an aggregate budget.
- Attribution: a delegated approval writes `approved` and a verdict
  paragraph naming the policy, the review record and the attempt, and is
  distinguishable from the owner's own verdict wherever the record is
  shown; the merge and `done` are the ordinary `integrate` acts. Introduce a
  field only if prose attribution proves insufficient.
- Verify before the target moves: make the merge in a temporary worktree,
  run the policy's commands there, and only then merge on the target. A
  failure leaves the target unchanged and the candidate in review with the
  reason, exactly as a refused `integrate` does.
- Trigger: the owner process of a finishing attempt, or an explicit sweep
  command the owner or an external scheduler runs. Grove starts no resident
  service (G-101); an unchanged wait does not retry.
- Everything the policy does not name waits for the owner as today.
  Reconcile the brief's lifecycle paragraph and the [approval](G-059-approval.md)
  term if the answer changes their meaning; the existing reviewer
  definition suffices initially, with its findings consumed by the policy.
- Out of scope: choosing the policy's conditions for the owner; a second
  reviewer definition; executing several selected items
  ([G-162](G-162-bounded-work-selection.md)); a shared candidate for a chain
  ([G-163](G-163-selected-work-review-boundary.md)).

## Acceptance

1. With no policy, nothing changes. With one, every automatic act is visible
   in advance, as what would happen to each candidate and why, and
   attributable afterwards, with the policy, evidence, attempt and cost in
   the record and on the board.
2. A candidate outside the policy waits for the owner. One inside it is
   integrated only after the merged result passed the policy's verification;
   a failure leaves the target and the record unchanged with the reason.
3. Automatic resolution runs at most once per target movement within the
   aggregate budget; a repeated conflict waits with the reason.
4. A delegated verdict is distinguishable from the owner's in the record and
   on the board, and reversible: the record names the merge to revert.
5. Exercised with fake providers on clean, conflicting, failing-verification
   and out-of-policy candidates, then a separately bounded real-provider
   trial on a disposable project before use in this repository.

## Next

The owner answers G-179. Then refine this scope to the selected option and
prepare the policy shape, attribution and the verified merge before
assignment through `$grove-work G-180`. Depends on G-177 and G-178. No
policy, spend, merge or implementation is authorized here.
