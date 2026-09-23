---
id: "G-044"
type: work
title: "Review candidates and integrate approved work locally"
status: active
created: "2026-09-21T00:54:16Z"
updated: "2026-09-23T00:25:14Z"
kind: feature
size: large
priority: 3
depends_on: ["G-038", "G-043"]
relates_to: ["G-035", "G-046", "G-064"]
formerly: "W-026"
---

## Outcome

Review an implementation in Grove, record feedback or approval against its
candidate, and integrate approved work locally with an explicit result.

## Scope and bounds

Lead with outcome, changed behavior, consequential decisions, verification,
open findings and follow-ups; progressively disclose artifacts, changed files
and optional diffs. Reuse G-038's review evidence rather than creating another
editable final report. Provide equivalent CLI operations for supported actions.

Bind approval to the candidate and relevant inputs/target state. Revalidate
before integration; conflicts, changed candidates and moved targets have clear
outcomes. Preserve partial work and dirty worktrees. Only clean up branches or
worktrees after integration is proven and retained files/evidence are safe.
Local integration only; no remote PR/push/deployment requirement.

Cleanup means safe branch/worktree cleanup, not moving completed records.
G-064's stable identity and placement apply to integration and feedback too.

## Acceptance

1. The owner can understand and judge a completed interactive candidate without
   the originating chat, inspect evidence, and optionally read diffs.
2. Approval, rejection/feedback and integration results remain attributable;
   a changed candidate cannot silently inherit approval.
3. Feedback requesting implementation returns work to Active and preserves the
   earlier review. Before a runner exists, it produces an actionable interactive
   continuation; automatic relaunch is added in G-046.
4. Exercise local merge success, conflict/refusal, changed target, stale approval,
   uncommitted changes and cleanup refusal. Integration failure does not mark
   Done or discard work. Report approval and integration as separate facts.
5. The owner accepts the review hierarchy and connected terminal workflow.

## Preparation and next

No implementation plan exists yet. Create a linked plan through the supported
CLI from G-038 and G-043. Resolve approval representation, supported merge strategy,
post-merge status publication and safe cleanup policy before building commands.
Human review is the default; autonomous judging needs a later explicit policy.
