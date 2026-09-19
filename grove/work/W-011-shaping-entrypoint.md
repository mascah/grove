---
id: "W-011"
type: work
title: "Shape project work through reusable agent instructions"
status: proposed
created: "2026-09-19T21:20:41Z"
updated: "2026-09-19T21:21:07Z"
kind: feature
priority: 2
size: small
relates_to: ["W-010", "W-009"]
---

## Outcome

Let the owner discuss Grove with Claude or Codex and have the agent maintain
useful proposed work, questions, and decisions through the restart's files and
CLI, without bespoke workflow instructions on every session. Supply a thin
repository-owned entrypoint over one shared authoring guide. This is proposed
dogfooding work, not an installed skill or an authorized scheduler.

## Why now and evidence

On 2026-09-19 the owner described two proposal sources: interactive planning
with Claude/Codex, and future scheduled or triggered headless research. The
[predecessor and Bench review](../../docs/reviews/2026-09-19-shaping-and-runner-evidence.md)
distinguishes conversational exploration, outcome-oriented shaping, and process
supervision. W-010 owns reusable implementation handoffs; this record owns
authoring and shaping. Neither must wait for a board or runner to be useful.

## Constraints

Use current work/question/decision types and lifecycles. No new schema, agent
processes, timers, board mutation, claim subsystem, or automatic merge. Preserve
concurrent implementation worktrees. Follow this repository's authority and
record model, not predecessor storage or close/archive commands.

Common guidance owns behavior; any Claude/Codex entrypoints are thin references.
Inspect harness discovery locally before choosing locations. A read-only context
command from W-010 can be reused when available, but is not a prerequisite for
an initial guide that uses existing list/show/check/new/update commands.

## Proposed behavior

- Begin from the brief, current records, actual code/evidence, and existing
  worktrees. Find the owning record before creating another; do not infer missing
  work from the current branch alone or turn Done into integration proof.
- Explore open ideas conversationally; distinguish the owner's settled choices
  from suggestions. Shape a selected outcome by investigating technical unknowns,
  naming constraints, testable acceptance, necessary design, and a concrete Next.
- Create IDs with `go run ./cmd/grove new`. Use revision-checked `update` for
  supported fields; edit bodies and linked plans as ordinary source, preserving
  concurrent edits. Do not pretend the CLI supports body updates or arbitrary
  metadata. Validate and check links before a handoff.
- Put proposed work in proposed status. Persist unresolved human questions and
  affected work links; only assert accepted decisions with attributable authority.
  Keep brief direction, work requirements, and plans in their respective owners.
- Report the exact checkout/branch containing new proposals and how to find
  them in W-009's source-scoped board. A proposal not present in the selected
  checkout cannot be promised a card in that checkout's Proposed column.
- Describe headless adaptation in the shared guide: explicit research mandate,
  allowed writes and bounds, evidence-backed proposals, durable questions and
  waiting conditions. This slice does not run that mode, choose a publication
  policy for it, or implement schedules. Never promote new proposals into
  implementation merely because a research process completed.

## Acceptance

1. A short entrypoint in each supported interactive harness loads one shared
   authoring guide and current project context without predecessor initialization.
2. A real Grove planning session produces or refines a bounded work proposal,
   an unresolved question, and an attributable decision only where one exists;
   it does not manufacture records merely to satisfy a checklist. Retain evidence
   and the concrete next action. Exercise missing cases in disposable fixtures.
3. Probes cover an existing overlapping proposal, an unanswered product choice,
   a stale revision, work owned on another branch, and a headless caller that
   must return a wait. No unsupported statuses/metadata or duplicate IDs appear.
4. Records validate and links resolve; the selected checkout and cross-branch
   visibility are reported accurately. No code, implementation status, sibling
   repository, or concurrent checkout is changed by the planning session.
5. Claude and Codex use the same responsibilities without duplicated editable
   policy. If one harness cannot be exercised, state that limit rather than
   claiming verified support.

## Next

Prepare a small implementation plan linked from this body: select repository
entrypoint locations, write the shared authoring instructions from the reviewed
responsibilities, and dogfood on the next requirements conversation. Reuse W-010's
context interface if implemented; do not block on it or expand into a runner.
Keep scheduled publication policy and launching in separately shaped work.
