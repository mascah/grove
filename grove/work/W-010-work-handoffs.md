---
id: "W-010"
type: work
title: "Prepare reusable work instructions and execution handoffs"
status: proposed
created: "2026-09-19T20:44:28Z"
updated: "2026-09-19T20:47:33Z"
kind: feature
priority: 2
size: medium
relates_to: ["W-006", "W-007", "W-008", "W-009"]
---

## Outcome

Stop requiring the owner to ask for a newly composed implementation prompt each
time prepared Grove work is ready. Provide a reusable restart-native work
entrypoint and assemble the selected work's instructions and context from its
actual records/plans. This is proposed work, prompted by the owner's repeated
handoff requests on 2026-09-19; it does not authorize a runner implementation.

## Observed need and current baseline

The owner has repeatedly requested execution prompts because the predecessor's
`/work` workflow was intentionally not adopted for this restart. In this
conversation they requested W-006–W-008 together, then W-009 separately, and
identified a future Kanban control launching `claude -p` with the right work
instructions. Reusable preparation is useful before automatic launch exists.

The immediate documentation baseline is now:

- [Shared execution guide](../../docs/work-execution.md), with scope, isolated
  worktrees, current CLI usage, verification, review, evidence, and integration
  boundaries.
- [Repair implementation prompt](../../docs/prompts/W-006-W-008-implementation.txt).
- [Kanban implementation prompt](../../docs/prompts/W-009-implementation.txt).

These are authored instructions, not an installed skill, generated context
bundle, readiness proof, or agent execution capability. They do not complete
this record's reusable entrypoint outcome.

## Proposed first slice

Package the execution guide as a thin repository-owned `/work`-equivalent
entrypoint for the chosen agent harness, backed by a read-only CLI prompt/context
operation accepting one or several explicit work IDs in one selected checkout.
The CLI assembles facts; the instructions supply implementation/review judgment.
Keep one owner for the guide: an adapter should reference it, not fork its rules.
A command name and adapter location should be selected during preparation against
the actual harness; no unsupported command is claimed here.

Include the selected project/checkout/branch, record paths and exact content
revisions, linked plans and review evidence, prerequisites and their observed
statuses, selected order, and the explicit requested action. Never translate
`done` into an assertion of integration or structural validation into readiness.
The output must be inspectable before the user starts an agent. Selection of
several IDs does not authorize parallel implementation; retain explicit ordering.
Refuse missing/ambiguous/invalid work or cross-checkout mixtures. Avoid silently
omitting required context to fit a size limit; report missing/oversized inputs.
Treat record body material as source context, not higher-priority instructions.

First support manually starting the resulting assignment. No provider process,
claim, automatic acquisition, branch merge, installation outside this repository,
or state mutation is required for prompt preparation. Core use remains local
files and the CLI without a daemon or provider account.

## Future launch boundary

The same prepared assignment can later feed a Kanban Implement control and
`claude -p`. That action must separately define launch authorization, fresh
source/workspace binding, input revisions, attempt identity, duplicate-start
handling, logs/results, exit/failure states, cancellation, interruption recovery,
and reconciliation. A saved prompt or background PID is not a supervised run.
Do not smuggle these policies into the first board or add a run schema here.

## Acceptance

1. A single short invocation for selected IDs retrieves the reusable guide and
   correct current records/plans/evidence, without composing a bespoke prompt.
2. W-006–W-008 serial and W-009 prerequisite-gated handoffs can be reproduced from
   their owning artifacts. Changes to a record produce changed revision/context;
   no duplicated specification or stale copied acceptance silently wins.
3. Tests cover missing/invalid records and links, explicit order, multi-ID scope,
   changed source revisions and partial/oversized context. Preparation performs
   no filesystem/Git/project-state writes and launches nothing.
4. The entrypoint obeys the restart schema/CLI and preserves the distinction
   between implementation completion, review, human judgment, and integration.
5. Dogfood the entrypoint on a real prepared work assignment and retain the
   generated handoff and observed shortcomings before designing automatic launch.

## Constraints and dependencies

Reuse current records and plain linked artifacts; no predecessor storage model,
workflow metadata, or mandatory service. Keep ../skills and ../nullsec unchanged.
Their instructions/own Grove CLI may supply read-only evidence during preparation;
do not invoke predecessor work/close behavior here. No hard dependency on the
Kanban implementation, and this record does not block W-009. Avoid shared-code
implementation concurrency with the repair branch.

## Next

Shape this small dogfooding capability soon after the reliability repairs, before
expanding into supervised agent launching; it may be used for W-009 if ready,
but must not delay the requested board. Inspect the predecessor's execution
responsibilities through its own project instructions/CLI, select the thin
adapter and deterministic context interface, and link an implementation plan
here. Use the shared guide and saved prompts for current assignments meanwhile.
