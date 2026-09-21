---
id: "W-020"
type: work
title: "Hand implementation candidates into revision-bound human review"
status: proposed
created: "2026-09-21T00:54:14Z"
updated: "2026-09-21T01:03:38Z"
kind: feature
size: medium
priority: 1
depends_on: ["W-019"]
relates_to: ["D-004", "W-018", "W-010"]
---

## Outcome

An implementation returns a revision-bound candidate for human review instead
of leaving the owner to infer its disposition from Active/Done and chat prose.
Make the interactive loop work before adding a managed runner.

## Selected contract

Target Proposed → Active → Review → Done; Abandoned requires an explicit human
decision. For implementation, Done means accepted and integrated into the
configured target. Define equivalent completion for research/design deliverables.
Preparation, independent agent review and waiting are activities or additional
facts. A failed/interrupted terminal attempt does not imply Review readiness.

An explicit assignment authorizes bounded execution; a proposed record or
status edit alone does not. Prepare missing plans within that mandate, preserve
the small-work exception, and ask about consequential choices rather than routine
technical steps. Keep interactive/headless instructions shared.

## Acceptance

1. CLI/schema and shared execution guidance support Review and its transitions;
   the existing board can display the added status without requiring a redesign.
2. The handoff links outcome, changed behavior, decisions, acceptance evidence,
   test revision, independent review findings, unresolved issues and next action.
   A new session can inspect it without the original chat.
3. Review identifies candidate/input revisions. Changed content cannot silently
   inherit approval. Feedback that starts implementation returns work to Active
   and retains prior evidence.
4. Historical Done records retain their original meaning and evidence. Design
   and document an explicit migration; do not infer integration from status.
5. Manual local approval/integration can complete the loop with an honest
   recorded disposition. Automated merge/cleanup and the review TUI belong to
   W-026; managed process ownership belongs to W-027.
6. Exercise successful, waiting, failed, feedback and changed-candidate paths;
   reconcile the model, guide and CLI documentation. Owner usability judgment
   belongs to the real loop in W-021.

## Preparation and next

Depends on W-019's artifact contract. No implementation plan exists yet.
Prepare `docs/plans/W-020-review-lifecycle.md`, resolving representation of
approval/integration evidence and a migration for branch-local historical Done.
Keep lifecycle rules independent of any particular harness or subprocess.
