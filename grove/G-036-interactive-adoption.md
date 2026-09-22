---
id: "G-036"
type: work
title: "Complete the interactive Grove adoption milestone in nullsec"
status: done
created: "2026-09-21T00:54:13Z"
updated: "2026-09-22T21:09:52Z"
kind: release
size: large
priority: 1
members: ["G-025", "G-037", "G-065", "G-052", "G-038", "G-039", "G-040", "G-041"]
relates_to: ["G-035", "G-042", "G-043", "G-044", "G-045", "G-046", "G-064", "G-065"]
formerly: "W-018"
candidate: "977d98b"
---

## Outcome

Use this Grove for one real nullsec change from interactive shaping through
implementation, human review and local integration, without predecessor skills
or reconstructing a conversation. This is the first adoption milestone selected
by the owner in [G-035](G-035-interactive-adoption.md).

## Scope and coordination

Members own shaping, knowledge/artifacts, flexible storage and full reconciliation,
review handoff, native dogfooding, portable setup and the nullsec pilot. This
record owns their connected result. The storage direction was revised by the
owner on 2026-09-21 in [G-064](G-064-stable-knowledge.md).
Membership is not a prerequisite that blocks implementing children.
The [adoption roadmap](G-047-adoption-roadmap-plan.md) is the current
coordination plan, not a code-level implementation plan or batch assignment.

## Acceptance

1. A fresh session can find the brief, selected work, applicable knowledge,
   current preparation and next action through Grove-owned entrypoints.
2. One real Grove assignment exercises the native loop and records its limits.
3. A migration rehearsal preserves nullsec identity, links and important
   knowledge; the owner then selects live scope and the explicit cutover.
4. One real nullsec change has attributable shaping, implementation, independent
   review where required, human disposition, integration and updated knowledge.
5. The owner judges the loop useful enough to continue using. Missing judgment
   stays open; child completion or automated checks cannot substitute for it.

## Verdict

The owner closed this milestone on 2026-09-22, once every member was done:

> I've been using the loop in this project itself to know its already
> valuable and its basically the same process that the nullsec project was
> already using. We can mark these as accepted for now and I can always come
> back later with tweaks.

Against the acceptance, as closed:

1. Met by the members: `grove brief`, `grove context`, and the `grove-work`
   and `grove-shape` entrypoints (G-025, G-037, G-040).
2. Met by G-039, with its limits in
   [G-078](G-078-g-039-trial-evidence-for-the-int.md).
3. Met by G-041: rehearsal, owner-selected scope, and the cutover, merged in
   nullsec at `3eb2785`.
4. Not exercised in nullsec. No nullsec change had run through the loop when
   this closed. The owner accepted this repository's own use of the loop
   (G-076, G-079, G-041 and the other work that went through Review) in its
   place, since nullsec's earlier process was essentially the same.
5. The owner's verdict above: continue, with tweaks later.

The candidate is `977d98b`, the commit on `main` that closed the last member
(G-041). This record has no implementation of its own.

## Next

None for this record. Improvements left open by G-078 are still unowned until
assigned: `new review` leaves `Examined`, `Findings` and `Disposition`
headings the CLI could fill or the work guide could name; the headless shaping
call did not take the question path, so exercise it again after the guide
change; the headless work row is still unexercised. On the headless
divergence the owner said on 2026-09-22 that it concerns them and that an
eval suite for the guides should be considered in future; nothing is proposed
or assigned for it yet. The related work (G-042 to G-046) stays proposed and
is assigned on its own.
