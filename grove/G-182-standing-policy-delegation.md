---
id: "G-182"
type: decision
title: "Delegate conflict resolution, approval and integration to a standing owner policy"
status: accepted
created: "2026-09-25T21:55:20Z"
updated: "2026-09-25T21:56:43Z"
relates_to: ["G-179", "G-180", "G-178", "G-177", "G-057", "G-058", "G-059", "G-060", "G-101", "G-142", "G-163"]
---

## Decision

On 2026-09-25 the owner answered [G-179](G-179-standing-policy-question.md)
with its option 4. Grove may act on a standing policy written in
`grove.yaml`, with no per-candidate human act, to start one bounded
resolution attempt when a candidate in review conflicts with the target, at
most once per target movement, and to approve and integrate a candidate
that meets the policy's written conditions after its independent review.
Each act is attributed in the record to the policy and its evidence, is
distinguishable from the owner's own verdict, and is reversible by an
ordinary revert. Anything the policy does not name waits for human judgment
as before. The policy's presence is the owner's standing instruction and
its absence means nothing automatic. The initial policy is narrow, proposed
in [G-180](G-180-policy-driven-integration.md) at the owner's request, and
the owner extends it.

The meanings of [candidate](G-057-candidate.md), [review](G-058-review.md),
[approval](G-059-approval.md) and [integration](G-060-integration.md) do
not change: approval was already the owner's or a delegate's acceptance of
one commit, and a review stays evidence, which the policy consumes.

## Alternatives

- **Nothing automatic** (G-179 option 1): unattended progress for
  implementation only; integration stays a morning chore.
- **Automatic resolution only** (option 2): delegates spend, not judgment.
  Not chosen as the initial state, though a policy without an approval
  section yields exactly it.
- **Delegated approval without automatic resolution** (option 3 alone):
  leaves conflicts to the manual feedback and relaunch loop that started
  this exploration.
- **The review agent decides by itself, with no written policy**: rejected.
  A review is evidence, never a verdict (G-058); a policy the owner can
  read is what makes each act attributable and extensible.

## Reconsideration

Reopen the conditions, not the delegation, when a delegated integration
lands a defect the written policy should have caught. Reopen the
delegation itself if the owner stops unattended use, if the conditions the
owner actually judges by cannot be expressed as deterministic checks plus
review evidence, or if the harness gains native gated merges that make
Grove's own policy redundant. A later decision supersedes this one if
[G-163](G-163-selected-work-review-boundary.md) selects a shared candidate
for a chain, whose review boundary this decision does not cover.
