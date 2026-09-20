---
id: "W-010"
type: work
title: "Prepare reusable work instructions and execution handoffs"
status: active
created: "2026-09-19T20:44:28Z"
updated: "2026-09-20T05:23:53Z"
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

Deliver a usable `grove-work` skill as part of Grove's agent interface. Its
workflow teaches the agent how and when to use the CLI through preparation,
implementation, verification, review, and evidence reconciliation. Context
generation alone does not satisfy this outcome. Interactive and headless callers
must consume the same workflow, with explicit handling of human availability.

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

## Predecessor behavior reviewed

The owner explicitly requested review of the actual
`../skills/skills/work/SKILL.md`. The
[2026-09-19 predecessor review](../../docs/reviews/2026-09-19-predecessor-work.md)
inspected that skill, its supporting instructions, instruction-generator code,
and headless adapter at sibling revision `ec87bb2`. It is required design input
for this work, not authority to invoke predecessor workflows here.

The old skill prepared missing plans, coordinated implementers and reviewers,
bounded retries, recovered interrupted attempts, verified connected behavior,
and reconciled evidence before an integration handoff. Its `grove launch`
printed instructions; frozen export/reconcile and a `claude -p` process adapter
were separate facilities. The current shared guide implements only part of that
instructional responsibility and must not be presented as equivalent already.

## Proposed first slice

Package the execution guide as a thin repository-owned `/work`-equivalent
entrypoint for the chosen agent harness, backed by a read-only CLI prompt/context
operation accepting one or several explicit work IDs in one selected checkout.
The CLI assembles facts; the instructions supply implementation/review judgment.
Keep one owner for the guide: an adapter should reference it, not fork its rules.
Thin describes the harness adapter, not reduced workflow responsibilities. The
shared guide owns the substantive behavior, including which CLI operations to
use and when to stop, wait, resume, or hand off.
The [shared implementation plan](../../docs/plans/W-010-W-011-agent-handoffs.md)
now proposes `grove context` plus repository-local `grove-work` adapters for
Claude and Codex. The command assembles facts; adapters explicitly load the
shared guide and required project instructions. These interfaces are prepared
design, not implemented commands or verified harness behavior.

Use the review's responsibility mapping to select the supported execution path:
plan preparation, implementation/review ownership, bounded retry/stop conditions,
checkpoint recovery, and evidence reconciliation must be explicit. Do not copy
the predecessor's model assignments or controller API by assumption. Keep
unimplemented controller/runner behavior visibly deferred. The instruction
entrypoint may prepare missing plans within the authorized outcome; distinguish
that preparation task from a frozen, implementation-ready assignment. A missing
required source is an error, while an absent plan must produce an explicit
preparation step rather than silent omission or a false readiness claim.

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

The owner described two proposal sources and a later board action on 2026-09-19:
interactive shaping, unattended research/planning, then an option to start
headless implementation when activating selected work. The companion
[W-011](W-011-shaping-entrypoint.md) owns authoring instructions; this record
continues to own the implementation assignment. Do not turn it into a combined
shaping, scheduler, TUI mutation, and runner project.

The reusable instructions need an explicit interaction mode. An interactive
caller can ask for a consequential missing decision; a headless caller must
persist the question/affected work and return a concrete waiting condition,
continuing only independent authorized work. Include outcome, allowed writes,
verification/review obligations, and completion/wait/stop conditions in the
assignment. A provider-specific session ID or launch command is not the mandate.
Mode selection must reach the shared guide explicitly; neither a short skill
invocation nor a headless process should fall back to interactive assumptions.
The intended headless entrypoint invokes this same skill or explicitly loads
its shared guide, without maintaining a second editable work prompt.

The same prepared assignment can later feed a Kanban Implement control and
`claude -p`. That action must separately define launch authorization, fresh
source/workspace binding, input revisions, attempt identity, duplicate-start
handling, logs/results, exit/failure states, cancellation, interruption recovery,
and reconciliation. A saved prompt or background PID is not a supervised run.
Do not smuggle these policies into the first board or add a run schema here.

Settled lifecycle requirement: closing the TUI leaves running sessions working;
reopening reconnects to them, while Stop is a separate action. The
[shaping and runner evidence review](../../docs/reviews/2026-09-19-shaping-and-runner-evidence.md)
examines Bench's stream capture and tmux host as possible sources of mechanisms.
No runtime choice is settled. Raw events, provider completion, verified work
acceptance, and integration must remain distinct. Status updates alone must not
trigger a process; the later activation UI should offer manual activation or
explicitly starting implementation.

## Acceptance

1. A single short invocation for selected IDs retrieves the reusable guide and
   correct current records/plans/evidence, without composing a bespoke prompt.
2. W-006–W-008 serial and W-009 prerequisite-gated handoffs can be reproduced from
   their owning artifacts. Changes to a record produce changed revision/context;
   no duplicated specification or stale copied acceptance silently wins.
3. Tests cover missing/invalid records and links, explicit order, multi-ID scope,
   changed source revisions and partial/oversized context. CLI prompt/context
   preparation performs no filesystem/Git/project-state writes and launches
   nothing; agent-authored plans belong to the separately authorized execution.
4. The entrypoint obeys the restart schema/CLI and preserves the distinction
   between implementation completion, review, human judgment, and integration.
5. Dogfood the entrypoint on a real prepared work assignment and retain the
   generated handoff and observed shortcomings before designing automatic launch.
6. Document reuse/adaptation/deferral of the reviewed predecessor responsibilities.
   Verify missing-plan preparation, external blockers, serial batch handoffs,
   interruption/resume, and pending human judgment for the supported path.
   Do not claim full `/work` equivalence from prompt-generation tests alone.
7. The assignment states whether human interaction is available. A simulated
   headless missing-decision case produces a durable question/wait handoff,
   without invented answers, automatic retries, or a provider launch. Interactive
   and headless instructions share outcome/constraints/acceptance sources.
8. Document the interactive skill invocation and intended `claude -p` entrypoint,
   including how interaction mode reaches the guide. Verify that both paths load
   the same workflow and CLI guidance. Distinguish adapter/source inspection and
   simulated behavior from an actual headless harness trial; leave any unrun
   trial explicit. This slice does not require implementing or launching a runner.

## Constraints and dependencies

Reuse current records and plain linked artifacts; no predecessor storage model,
workflow metadata, or mandatory service. Keep ../skills and ../nullsec unchanged.
Their instructions/own Grove CLI may supply read-only evidence during preparation;
do not invoke predecessor work/close behavior here. No hard dependency on the
Kanban implementation, and this record does not block W-009. Avoid shared-code
implementation concurrency with the repair branch.

## Next

Implement with W-011 using the linked shared plan and
[Fable handoff](../../docs/prompts/W-010-W-011-implementation.txt). The owner agreed
to this bounded dogfooding investment on 2026-09-19. Verify repair integration
and coordinate shared CLI/docs ownership with W-009; prefer following its current
handoff without making it a semantic dependency. The prepared path is native
serial execution with independent review, explicit headless waits, and prose
checkpoints; the old controller/runtime machinery remains deferred. Dogfood the
entrypoint before shaping a manually launched supervised run. Use existing saved
prompts until the new interfaces are actually implemented and verified.
