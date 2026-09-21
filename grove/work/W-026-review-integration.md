---
id: "W-026"
type: work
title: "Review candidates and integrate approved work locally"
status: proposed
created: "2026-09-21T00:54:16Z"
updated: "2026-09-21T01:03:40Z"
kind: feature
size: large
priority: 3
depends_on: ["W-020", "W-025"]
relates_to: ["D-004", "W-028"]
---

## Outcome

Review an implementation in Grove, record feedback or approval against its
candidate, and integrate approved work locally with an explicit result.

## Scope and bounds

Lead with outcome, changed behavior, consequential decisions, verification,
open findings and follow-ups; progressively disclose artifacts, changed files
and optional diffs. Reuse W-020's review evidence rather than creating another
editable final report. Provide equivalent CLI operations for supported actions.

Bind approval to the candidate and relevant inputs/target state. Revalidate
before integration; conflicts, changed candidates and moved targets have clear
outcomes. Preserve partial work and dirty worktrees. Only clean up branches or
worktrees after integration is proven and retained files/evidence are safe.
Local integration only; no remote PR/push/deployment requirement.

## Acceptance

1. The owner can understand and judge a completed interactive candidate without
   the originating chat, inspect evidence, and optionally read diffs.
2. Approval, rejection/feedback and integration results remain attributable;
   a changed candidate cannot silently inherit approval.
3. Feedback requesting implementation returns work to Active and preserves the
   earlier review. Before a runner exists, it produces an actionable interactive
   continuation; automatic relaunch is added in W-028.
4. Exercise local merge success, conflict/refusal, changed target, stale approval,
   uncommitted changes and cleanup refusal. Integration failure does not mark
   Done or discard work. Report approval and integration as separate facts.
5. The owner accepts the review hierarchy and connected terminal workflow.

## Preparation and next

No implementation plan exists yet. Prepare `docs/plans/W-026-review-integration.md`
from W-020 and W-025. Resolve approval representation, supported merge strategy,
post-merge status publication and safe cleanup policy before building commands.
Human review is the default; autonomous judging needs a later explicit policy.
