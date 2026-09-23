---
id: "G-107"
type: work
title: "Reconcile current documentation for a small external preview"
status: proposed
created: "2026-09-23T15:51:37Z"
updated: "2026-09-23T16:07:55Z"
relates_to: ["G-036", "G-078", "G-108", "G-109", "G-110"]
---

## Outcome

Give fresh agent sessions and prospective users an accurate, navigable account
of Grove as it works today, without carrying the completed adoption push as
current instructions.

Owner intent, shaping conversation 2026-09-23: Grove is a "minimum functional
product," useful for its own development, and nullsec has cut over. The owner
requested documentation reconciliation to avoid stale conversation context and
then selected "Prepare for a small external preview" as the next-phase audience.
That audience is selected; this work and the investment sequence below remain
proposals, not assignments or approval of a release.
"Preview" was the assistant's shorthand for initial use outside the owner's
projects, not a selected release channel or compatibility promise; the owner
asked for that clarification in the same conversation.

## Scope and constraints

Observed at main `f27444e43e581f7cc01856c459a514b57d62a406`:

- [AGENTS.md](../AGENTS.md) still opens with restart/read-only framing and says
  the runner contract remains open, alongside accumulated implementation history.
- [The brief](brief.md) describes shipped features as future work, the installed
  CLI as the predecessor, and branch-local Done as still permitted. Its adoption
  sequence has completed. [G-036](G-036-interactive-adoption.md) records the
  owner's closure, including the explicit substitution of Grove dogfooding for
  a real nullsec change; do not erase that evidence distinction.
- [README.md](../README.md) calls adoption the next milestone and contains a
  "Resume this conversation" prompt. Its roughly 6,000 words mix onboarding,
  command reference, and history. [The record model](../docs/record-model.md)
  mixes the implemented contract with starter-era design narrative.
- `docs/prompts/` retains old assignments naming removed paths and typed IDs;
  some prescribe race-suite practices superseded by repository policy.
- [The work guide](../docs/work-execution.md) and
  [shaping guide](../docs/work-shaping.md) are the shared workflow sources,
  embedded through [guides.go](../guides.go). Repository adapters load those
  sources; generated adapters load the binary's guides.

Proposed approach: inventory current instruction and documentation entrypoints,
then give each fact one editable owner. AGENTS.md owns durable repository policy;
the brief owns compact purpose, constraints and selected direction; README owns
the product introduction and routes to onboarding/reference; the record model
owns the current schema; shared guides own workflow. Keep adapters thin.

Reconcile documents outside `grove/`, plus the brief and narrowly identified
passages that still present completed adoption coordination as current. Preserve
dated reviews, decision authority, stable record identity/paths, and historical
acceptance limits. Retire spent prompts from the live instruction surface while
retaining provenance in Git or a linked historical record where useful.

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

Proposed follow-on order: a small behavioral eval baseline before substantial
skill changes; Attempts usability improvements and specific dogfooding fixes;
then preview packaging and onboarding. The eval and UI work can be independently
scoped. [G-078](G-078-g-039-trial-evidence-for-the-int.md) supplies a historical
headless missing-choice failure to revisit, not a claim that it still reproduces.

## Acceptance

1. An inventory accounts for current documentation, agent entrypoints, and spent
   prompts, with a clear retained, relocated, reconciled, or retired disposition.
   Current entrypoints no longer present completed adoption as pending or the
   predecessor as installed. Claims match the inspected revision and commands.
2. Product purpose and the owner's preview audience are discoverable in the
   brief. Next actions remain in work records; no second editable roadmap or
   status account is introduced. Historical verdicts remain attributable.
3. A fresh session can identify how to shape and execute assigned work, retrieve
   context in stages, and find current contracts without being directed through
   the old adoption roadmap. Record what was actually exercised, the harness
   and revision, and any limits; launch no session without its required mandate.
4. A reader unfamiliar with this repository can find what Grove does, its
   current capabilities and limits, installation/setup instructions, and the
   shape-to-review workflow without knowing old work IDs. The owner judges the
   proposed navigation useful for a small external preview; automated checks
   cannot supply that judgment.
5. Local links and command examples resolve or are explicitly labelled
   historical. `go run ./cmd/grove check` passes. Shared guide ownership and thin
   adapter behavior remain consistent, and no safeguarded behavior or authority
   boundary changes merely to shorten text. Report context reduction as evidence,
   not as a substitute for correctness.

## Next

The owner can review and assign this proposal. Before execution, commit the
agreed proposal so an isolated worktree based on main contains it, then invoke
`$grove-work G-107`. Preparation should inventory the surfaces and propose their
final ownership/navigation before editing. If fresh-session validation needs a
separate harness launch, include that bounded validation in the assignment.

Related proposals now own the [behavioral eval baseline](G-108-workflow-evals.md),
[Attempts usability](G-109-attempts-usability.md), and
[distribution preparation](G-110-external-preview.md). They remain independently
assignable. Actual release scope and compatibility expectations are later owner
choices, not prerequisites for this documentation cleanup.
