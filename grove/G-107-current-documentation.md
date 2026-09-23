---
id: "G-107"
type: work
title: "Reconcile current documentation and give each fact one owner"
status: active
created: "2026-09-23T15:51:37Z"
updated: "2026-09-23T16:49:34Z"
relates_to: ["G-036", "G-108", "G-109", "G-110"]
---

## Outcome

Give fresh agent sessions and readers of this repository an accurate,
navigable account of Grove as it works today, with each fact owned by one
editable document, and without carrying the completed adoption push as current
instructions.

Owner intent, shaping conversation 2026-09-23: Grove is a "minimum functional
product," useful for its own development, and nullsec has cut over. The owner
requested documentation reconciliation to avoid stale conversation context and
selected "Prepare for a small external preview" as the next-phase audience.
"Preview" was the assistant's shorthand for initial use outside the owner's
projects, not a selected release channel or compatibility promise; the owner
asked for that clarification in the same conversation. On 2026-09-23 the owner
narrowed this work, in a review of the first draft, to reconciliation and
ownership: newcomer onboarding and any preview-readiness judgment belong to
[G-110](G-110-external-preview.md), and agent behavioral trials to
[G-108](G-108-workflow-evals.md).

## Scope and constraints

Observed at main `4b2a01c6d301557b82ad989aabb80d826ef38a55`:

- [AGENTS.md](../AGENTS.md) opens with restart/read-only framing, says "The
  runner contract remains open" although [G-101](G-101-attempt-mechanism.md)
  is accepted and `run` is merged, and carries accumulated implementation
  history beside its policy.
- [The brief](brief.md) describes shipped features as future work, the installed
  CLI as the predecessor, and branch-local Done as still permitted. Its
  "Observed state" is dated at `c9904ea`, and its "Suggested sequence" points at
  the adoption roadmap that [G-036](G-036-interactive-adoption.md) closed; G-036
  records the owner's closure, including the explicit substitution of Grove
  dogfooding for a real nullsec change. Do not erase that evidence distinction.
- [README.md](../README.md) calls adoption the next milestone and ends with a
  "Resume this conversation" prompt. Its roughly 6,000 words mix onboarding,
  command reference, and history. [The record model](../docs/record-model.md)
  mixes the implemented contract with starter-era design narrative.
- `docs/prompts/` retains three spent assignments naming removed paths and typed
  IDs; two prescribe full race suites, which repository policy now forbids.
- [The work guide](../docs/work-execution.md) and
  [shaping guide](../docs/work-shaping.md) are the shared workflow sources,
  embedded through [guides.go](../guides.go). Repository adapters load those
  sources; generated adapters load the binary's guides.

Proposed approach: inventory current instruction and documentation entrypoints,
then give each fact one editable owner. AGENTS.md owns durable repository policy;
the brief owns compact purpose, constraints and selected direction; README owns
the product introduction and routes to reference; the record model owns the
current schema; shared guides own workflow. Keep adapters thin. The inventory
is a plan record linked by `work`, written before any document is edited.

In scope: documents outside `grove/`, and the brief. Records, including done
work and dated reviews, are history and are not edited: their identity, paths,
acceptance limits and verdicts stay as written. Retire spent prompts from the
live instruction surface; Git keeps their provenance.

Retain consequential constraints, including repository-safe Git subprocesses,
shared allocation/write locking, source freshness, terminal escaping, lifecycle
authority, staged retrieval, and verification policy. Removing the narrative of
how a constraint arose must not remove the constraint. Do not infer a need to
load every implementation record into every session.

This is documentation reconciliation. Behavior changes to the CLI, skills or
context assembly, an eval runner, Attempts redesign, release automation, Pages
publication, GitHub settings, and sibling-repository changes are separate work.
Flag newly found behavioral contradictions for follow-up rather than silently
changing the workflow contract. No release compatibility promise is selected.

## Acceptance

1. A plan record inventories current documentation, agent entrypoints, and
   spent prompts, giving each a retained, relocated, reconciled, or retired
   disposition. Current entrypoints no longer present completed adoption as
   pending, the predecessor as installed, or the runner contract as open.
   Claims match the inspected revision and commands.
2. The brief states product purpose and the owner's selected preview audience,
   and its sequence section names [G-108](G-108-workflow-evals.md),
   [G-109](G-109-attempts-usability.md) and
   [G-110](G-110-external-preview.md) as the proposals of that phase, their
   order labelled proposed until the owner selects one. Next actions remain in
   work records; no second editable roadmap or status account is introduced.
   Historical verdicts remain attributable.
3. A fresh session given only AGENTS.md can say how to shape work, how to
   execute assigned work, and how to retrieve context in stages, without being
   routed through the old adoption roadmap. This is a reading check, not a
   behavioral eval: one bounded interactive session, its harness and revision
   recorded, launched only with the assignment's mandate.
4. The README's opening says what Grove is, what this build does, and where a
   reader goes next, without requiring old work IDs. Whether that serves an
   external newcomer is G-110's judgment, not this work's.
5. Local links and command examples resolve or are explicitly labelled
   historical. `go run ./cmd/grove check` passes. Shared guide ownership and thin
   adapter behavior remain consistent, and no safeguarded behavior or authority
   boundary changes merely to shorten text. Report context reduction as evidence,
   not as a substitute for correctness.

## Next

The owner can assign this proposal. Before execution, commit the agreed
proposal so an isolated worktree based on main contains it, then invoke
`$grove-work G-107`. Preparation writes the inventory plan and proposes final
ownership and navigation before editing. The eval baseline, Attempts usability
and distribution preparation are independently assignable; none is a
prerequisite for this cleanup.
