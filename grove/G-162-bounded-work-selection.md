---
id: "G-162"
type: work
title: "Execute an explicitly selected set of work with bounded progress and review"
status: proposed
created: "2026-09-25T20:35:28Z"
updated: "2026-09-25T20:38:16Z"
kind: feature
relates_to: ["G-045", "G-046", "G-044", "G-101", "G-054", "G-056", "G-057", "G-059", "G-060", "G-161"]
---

## Outcome

A person can explicitly assign several Grove work items, leave the terminal,
and return to attributable progress, remaining waits, and reviewable results
without an agent choosing extra scope or concealing dependency and review
boundaries.

Owner intent, conversation 2026-09-25: the predecessor's multi-item work
workflow enabled useful progress over roughly eight unattended hours, but
left too much sequencing responsibility with the lead agent. The owner wants
to understand and control what will execute next or together. Eight hours is
an example of unattended use, not a completion guarantee or an authorized
runtime/spend budget. The owner requested proposals, not execution.

## Constraints

### Observed evidence

At main `3f2b923`, [the work guide](../docs/work-execution.md) accepts explicit
multiple IDs, executes selected prerequisites first, and defaults to one
implementer at a time over shared interfaces. [Context](../internal/handoff/context.go)
already returns selected IDs, deterministic dependency order, input revisions,
and external prerequisites. However, [the CLI](../internal/cli/cli.go)
requires exactly one work ID for `run`, and [attempt
ownership](../internal/attempt/attempt.go) is organized around one work ID.

[Approval](../internal/update/review.go) and
[integration](../internal/integrate/integrate.go) judge one candidate and reject
subsequent changes outside that work record. Merely allowing more IDs in
`run` would not establish shared review, partial completion, or coherent
integration. [G-101](G-101-attempt-mechanism.md) remains the accepted process
ownership decision. The [roadmap](G-047-adoption-roadmap-plan.md) previously
identified explicit selection, budgets, bounded fan-out, stop conditions,
and partial completion as preparation needed for batches.

The predecessor was inspected as files only:
`../skills/cli/grove/work.py` (`batch`, `render_batch`) calculates selected
dependency order and external blockers; `../skills/skills/work/SKILL.md`
requires shared planning, per-item acceptance, and joint verification.
These are historical evidence, not entrypoints to invoke or authority to
restore its execution machinery. Those sibling paths are not portable product
documentation.

### Proposed design and scope

- Accept an explicit, bounded selection whose IDs, source revisions,
  execution checkout/base, dependency order, resource bounds, and review
  boundary are visible before launch and recorded durably. Reuse context's
  relationship interpretation; never select the rest of a backlog or add
  an unselected prerequisite automatically.
- Start with sequential implementation. Preserve independent review and
  per-item acceptance and evidence. The absence of a dependency is not a
  concurrency mandate. New parallel implementation, scheduling, automatic
  merge authority, provider expansion, and reboot recovery are outside the
  proposed initial scope.
- Define aggregate budget enforcement, any subordinate attempt budgets,
  retry/repair bounds, and partial completion before launch. The lead agent
  may plan inside the assignment but cannot expand it or change its review
  boundary. Do not promise progress from elapsed time or record size alone.
- Make progress inspectable after terminal exit: completed implementation,
  candidates awaiting judgment, active work, not-yet-started items, and
  blocked items with their next action. Plan how this appears through
  existing Attempts inspection; this does not require the graph to ship.
- Preserve duplicate-start protection for overlapping selections, source
  freshness, recoverable checkpoints, Stop, and owner-loss reporting.
  Resume must reconcile completed work and live owners before dispatching
  again. Define which unaffected selected items may continue after a wait
  or failure, and show that policy in the assignment.
- Keep prerequisite delivery relative to the actual execution base and the
  chosen review boundary. Persist missing human decisions as questions;
  unchanged waits do not consume repeated attempts. Handle dependencies
  outside the selection without acquiring them or silently changing base.

**Open human choice:** [G-163](G-163-selected-work-review-boundary.md) blocks
this proposal's execution contract. Whether B may use A's changes before
the owner reviews and integrates A determines checkout topology, candidate
ownership, feedback invalidation, and approval/integration behavior. Do not
resolve that choice indirectly in a plan or acceptance interpretation.

After the owner answers, reconcile the contract in the
[work guide](../docs/work-execution.md), [record model](../docs/record-model.md)
where its semantics change, and the owning command/board documentation.
Preserve the settled meanings of [work](G-054-work.md),
[attempt](G-056-attempt.md), [candidate](G-057-candidate.md),
[approval](G-059-approval.md), and [integration](G-060-integration.md).
Introduce durable batch representation only where the selected mechanics
need it; this proposal does not choose a new record type or revive old fields.

[G-161](G-161-dependency-view.md) provides a complementary discovery and
selection surface. Its suggested earlier implementation is investment order,
not a dependency: CLI selection can exercise this outcome independently.

## Acceptance

1. A person can inspect and launch a bounded explicit selection with its
   source/base, prerequisite order, outside blockers, resource bounds,
   and owner-selected review boundary. A changed input or incompatible
   selection produces a useful refusal before dependent work starts.
2. A chain, a branching selection, and a dependency path through unselected
   work exercise the selected policy. No unselected implementation or
   integration occurs. Later work starts only when its prerequisites meet
   the boundary resolved in G-163.
3. Every selected item retains its acceptance, evidence, exact candidate when
   offered for review, and concrete Next. Shared changes, feedback to an
   earlier item, dependent evidence becoming stale, and partial completion
   have a documented and demonstrated review/integration disposition.
   A successful provider exit alone does not mark any item done.
4. Reopening inspection does not launch again. Overlapping selections cannot
   duplicate ownership; Stop and owner loss preserve completed and partial
   work. Resume does not repeat completed units or continue from stale inputs.
5. Budget exhaustion, retry limits, a new human question, an external
   prerequisite, and failure of one member expose attributable waits/results.
   Any continuation of unaffected selected work follows the recorded policy
   and stays within the aggregate bounds. An unchanged wait does not retry.
6. Exercise lifecycle failures with fake providers first, then seek a
   separately bounded real-provider trial of selected work in a disposable
   project. The owner can return after terminal exit and understand what
   changed, what still waits, and what needs review without reading raw logs.
   Report untested recovery and timing limits explicitly.

## Next

Checkpoint, 2026-09-25 (headless, `--until plan`): G-163 is resolved
(option 1, shared implementation reviewed together). Plan
[G-185](G-185-g-162-selected-work-plan.md) at
`sha256:6584eb48aa591d5d16b0f399081434a9d883d6f840898a64cc3e62906cbb80ef`
(commit `dbb2f0d`, branch `worktree-G-162` from `main` `fe97300`) prepares
one attempt over an explicit selection, a shared candidate judged per member
and integrated as a group, feedback reopening the group, the
partial-completion and continuation policy, and launch preview and refusals.
Status left proposed; nothing implemented. Launching the continuation after
reading the plan approves that plan revision:

    /grove-work G-162

or `grove run G-162` from this branch's checkout. Its first step records the
G-163 answer as a decision. The real-provider trial (acceptance 6) is an
owner step with its own budget, not part of the attempt.
