---
id: "G-042"
type: work
title: "Derive a project-wide current view of work"
status: active
created: "2026-09-21T00:54:15Z"
updated: "2026-09-22T21:40:04Z"
kind: feature
size: medium
priority: 3
depends_on: ["G-037", "G-038"]
relates_to: ["G-035", "G-002", "G-010", "G-011", "G-030", "G-031", "G-064", "G-065"]
formerly: "W-024"
---

## Outcome

Opening Grove from any linked checkout presents the same useful project-wide
current work view, with source observations and genuine divergence accessible
without making every user choose a branch first.

## Selected direction

Collapse identical observations; move demonstrably superseded states into
history; expose unintegrated progress and label live uncommitted changes.
Keep genuine competing changes visible. Never use the largest status, latest
timestamp or newest branch tip as authority. A revert is a real change, not
automatically an older state. Exact source routing and stale-selection checks
remain intact; view selection does not merge content or grant write authority.

Use G-065's stable identities and recursive discovery. Neither type prefixes
nor a completed/history folder establish which record is current or integrated.
"History" here is a view of evidence, not a filesystem move. Include supported
old/new schema sources and legacy/new IDs in the projection fixtures.

## Acceptance

1. Specify and fixture-test a deterministic projection for old branches,
   integrated work, unmerged progress, dirty/deleted records, branch-only work,
   genuine divergence, reverts and missing/invalid sources.
2. The policy distinguishes record history from unrelated branch-tip commits
   and committed observations from live overlays; ambiguity is explicit.
3. Equivalent observations from different invoking checkouts produce the same
   default view. Users can inspect the chosen evidence and other sources.
4. For the observed G-030 history, an old Proposed copy on G-023 does not
   obscure the Done/integrated result on main. Do not hard-code these branches.
5. Preserve read-only behavior, incomplete-result diagnostics, safe action
   targeting and G-031's batched Git reads. Avoid per-branch Git processes;
   any added history work must retain acceptable measured board load behavior.
6. Update G-002's implementation note and current documentation. Full visual
   redesign is G-043; projection should be independently testable and explainable.

## Preparation and next

The [plan](G-093-current-view-plan.md) holds the projection, the owner's
2026-09-22 decisions (uncommitted edits count, labelled; one card per
divergence in its earliest status; no target), and the steps. Implementation
is on branch `worktree-G-042` from `939d090`.
