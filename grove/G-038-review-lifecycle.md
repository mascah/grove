---
id: "G-038"
type: work
title: "Hand implementation candidates into revision-bound human review"
status: active
created: "2026-09-21T00:54:14Z"
updated: "2026-09-22T04:13:00Z"
kind: feature
size: medium
priority: 1
depends_on: ["G-037", "G-065"]
relates_to: ["G-035", "G-036", "G-023", "G-064"]
formerly: "W-020"
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

Build on G-065's identity/storage contract, as selected in
[G-064](G-064-stable-knowledge.md). Review, approval and integration
remain explicit operational facts; generic knowledge does not imply authority.
Lifecycle transitions preserve IDs and paths, including Done. No completion
folder or record-filing operation is part of integration cleanup.

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
   G-044; managed process ownership belongs to G-045.
6. Exercise successful, waiting, failed, feedback and changed-candidate paths;
   reconcile the model, guide and CLI documentation. Owner usability judgment
   belongs to the real loop in G-039.

## Preparation and next

G-037 is delivered; wait for G-065's compatible foundation. G-052's migration
is preferred first but is not a technical dependency. No implementation plan
exists yet. Create a linked plan record through the supported CLI, resolving representation of
approval/integration evidence and a migration for branch-local historical Done.
Keep lifecycle rules independent of any particular harness or subprocess.
