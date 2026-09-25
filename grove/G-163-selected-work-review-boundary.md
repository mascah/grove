---
id: "G-163"
type: question
title: "Where should review and integration stop a selected chain of work?"
status: resolved
created: "2026-09-25T20:35:30Z"
updated: "2026-09-25T23:22:35Z"
blocks: ["G-162"]
relates_to: ["G-161", "G-188"]
---

## Question

When the owner explicitly assigns A, B, and C, and B needs A while C needs B,
may the selected work proceed on the earlier items' unmerged changes, or must
human review and integration stop the chain between items?

This is an unanswered owner choice from the dependency-view exploration and
shaping conversation on 2026-09-25. The owner asked to shape the concept but
did not select a review boundary. It blocks the execution contract of
[G-162](G-162-bounded-work-selection.md), not the read-only dependency view
and preview in [G-161](G-161-dependency-view.md).

## Evidence and options

The desired use is meaningful unattended progress across a chosen set, with
more human control over scope and order. The [work guide](../docs/work-execution.md)
already permits selected prerequisites to be implemented sequentially in a
shared execution checkout. Current [approval](../internal/update/review.go)
and [integration](../internal/integrate/integrate.go) are centered on an
individual work candidate and reject later changes outside that record.
Neither establishes a multi-item review policy by itself.

1. **Explicit shared implementation, reviewed together.** The person chooses
   the group; A, B, and C are built in order before a human review of the
   combined result. Per-item acceptance and evidence remain visible. This
   enables progress through a dependent chain while the person is away,
   but makes the final review larger and requires an explicit treatment of
   shared candidates, earlier-item feedback, and partial completion.
2. **Human review and merge between dependent items.** B waits until A is
   accepted and integrated into its intended base. Each review is smaller,
   but an unattended chain stops at its first review boundary. Unaffected
   selected work may proceed only under the assignment's declared policy.
3. **Choose the boundary for each assignment.** Support both explicitly;
   neither an agent nor the presence of an edge chooses the policy. This
   offers flexibility but expands the initial execution and recovery
   contracts and their verification.

Assistant recommendation, not a decision: begin with option 1 for an
explicitly chosen shared implementation, sequential by default, with external
prerequisites delivered into the base and no automatic merges. It most
directly serves the owner's overnight example. A dependency edge still means
the prerequisite is needed; it does not itself authorize shared execution.
If intermediate human judgment is central to the intended use, select option
2 instead. The owner may also defer G-162 and start with G-161 alone.

## Next

The owner chooses the initial boundary and any exception that must be in
scope. Record the answer here, resolve this question through `grove update`,
and refine G-162. Because the answer changes candidate/review/integration
semantics, record a consequential choice as an attributable decision when
made, following the shaping guide. No answer or approval is implied by
creating these proposals.

## Answer
option 1


Recorded as decision [G-188](G-188-selected-work-shared-candidate.md), with the consequences plan G-185 drew from it.
