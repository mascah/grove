# W-018 adoption roadmap

Selected 2026-09-20; implementation has not started. The
[brief](../restart-brief.md) owns direction and
[D-004](../../grove/decisions/D-004-interactive-adoption.md) records authority.
[W-018](../../grove/work/W-018-interactive-adoption.md) owns the first milestone's
acceptance. This is its coordination plan, not a batch assignment or a detailed
implementation plan. Records below own their scope, acceptance and next action.

## How to continue without this conversation

1. Start with the chosen work record, repository instructions and the shared
   execution guide through `/grove-work W-ID` (Claude) or `$grove-work W-ID`
   (Codex). In this repository every command is `go run ./cmd/grove …`.
2. Inspect its current context, branch/worktree, prerequisites and linked current
   plan. Prepare a concise implementation plan when needed. The future plan
   paths named in records are destinations, not claims that those files exist.
3. Read only the relevant research section when its uncertainty matters. This
   roadmap and historical reviews are not mandatory context for every task.
4. Complete the selected assignment and its acceptance, preserve evidence and
   next action, and obtain its explicit integration disposition. Do not expand
   into all roadmap entries automatically.

Work remains Proposed until individually assigned and started. Approval of this
roadmap does not assign sibling migrations, background execution or automatic
merges; those actions belong to their explicitly selected work and mandate. Existing assignment authority still applies within a unit;
routine technical decisions do not require repeated permission.

## First adoption milestone

Preferred order is top to bottom. `depends_on` in the records expresses actual
prerequisites; order here also reflects product priorities.

| Work | Deliverable | Completion evidence |
| --- | --- | --- |
| [W-011](../../grove/work/W-011-shaping-entrypoint.md) | Native interactive shaping and shared authoring guidance | Real requirements session and fresh harness discovery; no fabricated knowledge records |
| [W-019](../../grove/work/W-019-knowledge-artifacts.md) | Domain terms and discoverable linked plans/reports/reviews | Compatible schema and staged retrieval, preserved identity/links |
| [W-020](../../grove/work/W-020-review-lifecycle.md) | Review lifecycle and durable candidate handoff | Revision-bound evidence, explicit completion migration, manual review/integration path |
| [W-021](../../grove/work/W-021-interactive-loop.md) | Complete loop on real Grove work | Fresh-session continuation, independent review, owner verdict and integration evidence |
| [W-022](../../grove/work/W-022-portable-bootstrap.md) | Minimal setup and versioned portable workflows | Disposable project adoption and observed Claude/Codex entrypoints |
| [W-023](../../grove/work/W-023-nullsec-pilot.md) | One real nullsec change using Grove | Rehearsal, authorized live cutover, full loop and owner's continue/revise verdict |

W-011 can ship with today's schema. It teaches where richer knowledge is not yet
supported rather than inventing metadata. W-019/W-020 extend the contract and
update the same workflow owners. W-021 checks their connected result before
W-022 carries it into another repository. W-023 completes W-018; member status
alone cannot establish the milestone's human acceptance.

For the initial loop, plans stay proportional and human review can read CLI/files.
Design board/detail/review together early, but do not make the full TUI redesign
a prerequisite for the real hobby-project change.

## Following investments

| Work | Deliverable | Boundary |
| --- | --- | --- |
| [W-024](../../grove/work/W-024-current-view.md) | Project-wide current projection | Resolve concrete ambiguous histories; preserve explicit sources and batched reads |
| [W-025](../../grove/work/W-025-board-detail.md) | Polished board, detail, artifacts, timeline and list/search | Owner-reviewed visuals and connected terminal checks |
| [W-026](../../grove/work/W-026-review-integration.md) | Progressive review plus local approval/integration actions | Candidate-bound approval, conflict/refusal and safe cleanup |
| [W-027](../../grove/work/W-027-durable-attempt.md) | One bounded durable implementation attempt | Fake-process failure probes, then a bounded actual-provider trial |
| [W-028](../../grove/work/W-028-managed-runs.md) | TUI launch, runs overview, reconnect, stop and feedback continuation | UI observes the independent owner and reuses the same review contract |

W-024 depends on the knowledge/lifecycle contracts, not on the act of migrating
nullsec. W-027 needs portable assignment and review contracts, not the TUI merge
screen; its later place is investment order. These distinctions allow deliberate
reordering without manufacturing dependencies or silently widening an assignment.

## Preparation questions, at their point of use

These are design tasks, not unanswered human preferences that block W-011 today.
Investigate routine choices; create a question through the CLI only when an
actual consequential human decision remains, linked to the work it blocks.

| Boundary | Questions to resolve | Owner |
| --- | --- | --- |
| Knowledge | Term representation; artifact identity, multi-work links, schema version and migration; brief location | W-019 |
| Review lifecycle | Candidate/input identity; disposition and integration receipts; historical Done migration; research/design completion | W-020 |
| Portability | Workflow packaging, managed adapter updates, executable coexistence and minimal initialization | W-022 |
| Live adoption | Current nullsec scope, retained/archive mapping, collision handling, rollback and explicit cutover | W-023 |
| Current projection | Target branch, record-level ancestry, dirty overlays, deletion/reverts, divergent-card placement | W-024 |
| Visual experience | Concrete layouts, narrow terminals, bounded Done defaults, timeline navigation, compatible Charm modules | W-025 |
| Local integration | Supported merge strategy, moved target, stale approval, status publication and cleanup failures | W-026 |
| Runtime | Native Claude background versus owned process; role/model controls, budgets, event protocol and owner-loss recovery | W-027 |
| Run experience | Separate run screen versus detail section, activity summarization, feedback relaunch | W-028 |

## Research and source map

The [direction evaluation](../reviews/2026-09-20-direction-evaluation.md) holds
the repository observations and external references behind this sequence.
Read the relevant subsection at preparation, not every predecessor skill.

- Shaping, terms and context: evaluation's **Knowledge and context**;
  prior [shaping/runner research](../reviews/2026-09-19-shaping-and-runner-evidence.md).
- Preparation, review and roles: evaluation's **Workflow and roles**;
  prior [predecessor work review](../reviews/2026-09-19-predecessor-work.md).
- Projection and TUI: evaluation's **Current state and presentation**, Q-001's
  historical answer, and the current versions/history implementation.
- Process ownership: evaluation's **Future runtime boundary**, then current
  provider documentation and installed capabilities. Old adapter code is
  evidence, not a supported implementation to copy unexamined.

## Later candidates, not implementation-ready assignments

- **Learning/debrief:** mine linked plans, reports and reviews for repeated
  failures or missing context; retain supporting examples and propose bounded
  changes to knowledge, guidance or checks. Collect useful evidence now; measure
  whether a change improves subsequent work before making universal policy.
- **Multiple work items:** explicit bounded selection, dependency/ownership
  analysis, sequential execution by default, per-item acceptance and shared
  integration evidence. Define budgets, maximum fan-out, stop conditions and
  partial completion before enabling unattended batches.
- **Research/dream loop and scheduling:** one research mandate, no automatic
  expansion into implementation; separate reviewable branch, evidence-backed
  proposals, durable human questions, deduplication and unchanged-wait handling.
  Scheduling is another caller of the proven run contract.
- **Dependency graph:** visualize actual dependency edges and distinguish them
  from membership or preferred order. Use real graph size/readability evidence
  before choosing a terminal rendering approach.
- **Additional knowledge types:** designs, flows, systems, research, releases,
  builds and team definitions are candidates. Introduce only when a distinct
  owner and retrieval need emerge; do not recreate Bench's entire taxonomy.
- **Keyborg adoption:** a second test of packaging and migration with its pinned
  Bench contract; shape it from the nullsec results and current project needs.
- **Autonomous judging, other providers and recovery after reboot:** separate
  policy/capability decisions. Interactive Claude/Codex support does not establish
  managed-provider parity or unattended acceptance authority.

## Completion and reconsideration

Use W-021 and W-023 evidence to adjust this order and the owning records.
Changes to product direction belong in the brief with an attributable decision.
Current progress belongs in work records, not another status table here. When
the first milestone completes, record the owner's verdict in W-018 and select
the next bounded investment from the observed friction.
