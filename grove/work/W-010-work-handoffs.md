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

That was the baseline before this work. The guide is now the workflow the
`grove-work` skill loads, and the saved prompts are spent; see Evidence.

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
shared guide and the repository's agent instructions. Both are implemented on
branch `worktree-W-010` (see Evidence for what has been verified in a harness).

Owner revision, 2026-09-20, superseding the conflicting prepared design:

- The repository-local adapters are intentional for this dogfooding phase.
  Plugin distribution, global installation, and portability to another adopted
  repository are not acceptance requirements here.
- The shared guide is Grove workflow only. This repository's development
  policy (restart history, Go verification, `go run ./cmd/grove`, branch
  conventions, coexistence with the predecessor) lives in `AGENTS.md`, once.
  The workflow never requires the predecessor's `grove:work`, `grove:close`, or
  `grove:shape` skills; the predecessor comparison is review evidence.
- Context is delivered in stages. `context` reads the configuration, the
  selected records, and explicit includes in full, and lists prerequisites,
  blocking questions, members, related and linked records, and links, with
  identity, status, revision, and how to retrieve them. The workflow includes
  the plan the record names when preparing or implementing, and must read the
  actual constraints of blocking questions and prerequisites before
  implementing a unit. A listing never stands in for required evidence, and no
  filename decides which artifact is current.
- The execution checkout is established and verified before the first write.
- W-011's shaping workflow should be able to reuse these boundaries; it is not
  part of this work.

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
revisions, prerequisites and their observed statuses, selected order, and the
explicit requested action; list linked plans and review evidence for retrieval
when the activity needs them (staged on 2026-09-20; first included in full). Never translate
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

Owner clarification, 2026-09-20: in 1 and 2, "retrieves" and "reproduced" mean
staged retrieval. The invocation delivers the guide and the selected records
with revisions at once, and the plan, prerequisites, questions, and evidence
through `--include` and `show` when the workflow's step needs them. Missing or
oversized requested sources still fail; nothing is truncated or summarized.

## Constraints and dependencies

Reuse current records and plain linked artifacts; no predecessor storage model,
workflow metadata, or mandatory service. Keep ../skills and ../nullsec unchanged.
Their instructions/own Grove CLI may supply read-only evidence during preparation;
do not invoke predecessor work/close behavior here. No hard dependency on the
Kanban implementation, and this record does not block W-009. Avoid shared-code
implementation concurrency with the repair branch.

## Evidence

Implemented on branch `worktree-W-010` from main `91edc0b`, in
`.claude/worktrees/W-010`, and revised there on 2026-09-20. On the owner's
instruction the branch was merged into `main` on 2026-09-20 with this work
still active: integration is not acceptance (establish it by Git ancestry, as
always). Not pushed. The [dogfooding evidence](../../docs/reviews/2026-09-19-W-010-dogfood.md)
holds the detail and keeps source inspection, tests, simulated runs, real
harness trials, and the owner's acceptance apart.

- `grove context WORK_ID...` (`internal/handoff`, `internal/cli/context.go`)
  and the README contract (`format_version` 2: selected records, configuration,
  and includes in full; everything else listed); the
  [work guide](../../docs/work-execution.md), workflow only, with staged
  reading and isolation before the first write; this repository's policy for
  it in `AGENTS.md`; thin `grove-work` adapters in `.claude/skills/` and
  `.agents/skills/`.
- Revision commits: `8ba9d89` (Markdown destinations decoded before URL
  interpretation), `10bfc51` (file identity for every source, checked before an
  alias is read or charged), `ff4deb3` (inclusion policy), `ded60de` and the
  commit carrying this text (guide, adapters, `AGENTS.md`, plan, evidence).
  Both defects were reproduced by failing tests before their fixes.
- Checks on the final code revision, uncached: `go test -count=1 ./...`,
  `go test -race -count=1 ./...`, `go vet ./...`, `gofmt -l .`,
  `go mod tidy -diff`, `go run ./cmd/grove check`, and a relative-link and
  anchor check of the changed documents all pass. `FuzzResolve` ran 4.7 million
  inputs cleanly with Markdown-escaped traversal seeds.
- Measured context: `W-012` 99229 to 2719 source bytes at the start; `W-009`
  210956 to 17383, 32579 with its plan; `W-010 W-011` 193078 to 20694,
  47932 with the shared plan (like for like, before this revision's text
  enlarged this record and the plan). The evidence has composition and what each
  stage reads.
- Independent review of `079a7da..ded60de` by a separate reviewer agent that
  edited nothing: ten guard mutations each failed a test, fuzzing and
  hand-built destinations found no escape or double decoding, output was
  deterministic, and documentation matched code. One consequential finding, a
  contradiction in the guide between redirecting a large result in step 1 and
  writing nothing before step 3, is fixed; so are the minor ones (an `included`
  row with no matching source path under an alias, link and record matching
  that disagreed, an unnamed changed record, stale failure advice), with
  regressions where code changed. The follow-up fixes and the final documents
  were checked by the implementer only.
- Real harness trials in disposable clones (Claude Code 2.1.278, codex-cli
  0.155.1): both discover the adapters from an explicit headless invocation,
  and both carried a small fixture assignment through staged reads, a worktree
  before the first write, a resolved blocking question read in full and
  honoured, and a handoff with `main` untouched.
- Acceptance 1–3: met by the command, tests, and the real W-006–W-008, W-009,
  and W-010/W-011 contexts, under the owner's staged-retrieval clarification.
  4, 6, 7: met for the supported path by the guide, the two simulated headless
  runs against the first guide (durable question and wait, unacquired external
  blocker, missing-plan preparation, serial batch, resume from a checkpoint),
  and the two real trials; interruption mid-implementation was not exercised,
  and the wait path was not rerun against the revised guide. 8: both
  invocations are documented and both were exercised headlessly on a fixture;
  the interactive invocations were not. 5: **open**. A fixture written by the
  workflow's author is not a real prepared assignment, and the owner has not
  judged it.

## Next

Owner: start a fresh interactive Claude Code session in the main checkout,
where the skill and `context` now are, and run

```text
/grove-work W-012
```

W-012 is the dogfooding target the owner chose on 2026-09-20: it is real Go
work in the board, lists three related records to retrieve when needed, has no
plan yet, and its Next holds an undecided owner choice (whether lineage
replaces a card's version list or sits beside it). The run should establish
what no trial has: that an interactive session discovers the skill; that the
workflow asks that one question before planning instead of choosing; that a
real assignment then goes from its ID through a prepared plan, code, the Go
checks, and review to a handoff in its own worktree, created before the first
write, without a composed prompt; and that the staged reads arrived when
needed. Interrupt it once midway and rerun the same invocation to exercise
resume. Then say whether this replaces asking for a prompt, and record what
happened in the dogfooding evidence before marking this done. Optionally try
`$grove-work` in the Codex TUI.

An earlier draft of this Next proposed `/grove-work W-011` from the W-010
worktree, stacked on the then-unmerged branch. The owner questioned it and it
was dropped: it would have stacked work on an unaccepted branch, was chosen
partly because it exercised a rule this revision had just written, and is
documentation-only work that mirrors the guide under test. The guide's rule
for records that exist only on another branch therefore stays untested beyond
reading.

[W-011](W-011-shaping-entrypoint.md) (`grove-shape`, `docs/work-shaping.md`) is
not started; it is ordinary later work from `main` and can reuse the adapter shape, the interaction-mode convention,
the workflow/policy split, staged reading, and isolation before the first
write. The shared ID counter stands at `W 18` after a fixture mistake recorded
in the evidence; the gap is historical and is not to be reset. Plugin
distribution, nullsec migration, an attachment schema, and supervised
launching stay separate, later investments.
