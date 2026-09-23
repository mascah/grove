# Grove brief

Direction reconciled 2026-09-23 ([G-107](G-107-current-documentation.md)).
This is the single current source of product intent. The owner selected the
interactive adoption milestone in [G-035](G-035-interactive-adoption.md),
the storage and identity direction in [G-064](G-064-stable-knowledge.md),
and the next phase's audience on 2026-09-23 (below). Work records own
acceptance, progress and next actions; the
[record model](../docs/record-model.md) describes the implemented schema.

## Purpose

Grove is a local project workspace where a person shapes intent with agents,
delegates bounded work, and returns to a clear review of the result. It connects
project knowledge, work, and execution evidence so the next human or agent can
continue with **the right context, at the right time**.

The first audience is the owner working across local hobby repositories and
worktrees. The first adoption milestone, a complete interactive
shape → implement → review → integrate loop,
[G-036](G-036-interactive-adoption.md), was closed by the owner on 2026-09-22.
Its nullsec acceptance was not exercised in nullsec: the owner accepted this
repository's own use of the loop in its place. Nullsec's cutover to this
Grove ([G-041](G-041-nullsec-pilot.md)) was part of that milestone. Keyborg, selected in G-035
as the second adoption test, has no work record yet.

The next audience, selected by the owner in a shaping conversation on
2026-09-23, is a small external preview: initial use outside the owner's
projects. "Preview" names that audience only. It is not a selected release
channel, launch commitment, support policy or compatibility promise; release
scope, platforms and licensing remain later choices.

## Selected foundations

- Go CLI, ordinary Markdown with YAML frontmatter, `grove.yaml`, configurable
  record storage, and Git. Core inspection needs no running service.
- Stable sequential IDs in one neutral `G-` namespace, coordinated across
  linked worktrees; revision-checked mutations; ordinary edits preserve
  identity. Separate clones still require collision checks and reconciliation.
- Humans and agents use the same records. Deterministic validation, retrieval,
  mutation and lifecycle mechanics belong in software; judgment belongs in
  instructions and attributable human or delegated decisions.
- Branch-local editing with a combined project view. Actions bind to an exact
  checkout and revision; browsing should not make source selection the main job.
- Bare `grove` opens the TUI; explicit commands remain noninteractive. The TUI
  should be visually appealing, keyboard-friendly, and use the Charm ecosystem.
- Grove-owned shared workflows with thin Claude and Codex adapters, carried
  inside the binary. They replaced the sibling skills in this repository and
  in nullsec.

## Knowledge, work, and evidence

Keep their lifetimes and owners distinct:

| Information | Owner and purpose |
| --- | --- |
| Brief | Compact project purpose, constraints, selected direction and next investment |
| Terms | Domain vocabulary and relationships; no implementation diary |
| Decisions | Consequential choices, authority, alternatives and reconsideration conditions |
| Questions | Unresolved matters, affected work, evidence, and who can answer |
| Work | Outcome, bounds, acceptance, dependencies and next action |
| Artifacts | Plans, implementation reports and reviews linked to work and relevant revisions |
| Attempts | Particular executions, inputs, owner, progress, result and recovery state |

Retain terms and discoverable artifacts. Plans need not become an
independent ticket lifecycle. Small work may carry its preparation in the work
record; substantial work gets a concise linked plan. Designs, flows, systems,
research, releases, builds and team definitions can begin as general knowledge;
new operational types are introduced when software needs distinct behavior.

**Stable identity, stable placement, flexible content, multiple views.** One
configured Grove root has a flat creation default and recursive discovery;
folders do not determine validity or meaning. General knowledge needs no
predefined semantic category or ticket lifecycle. Known operational records
retain explicit contracts for the facts software acts on. Classification,
title and status changes preserve identity and path. Completed records stay
put; views bound everyday clutter. Per-type folder/prefix settings, automatic
filing and a general schema-extension engine are not selected. This repository
had one deliberate reconciliation into that layout
([G-052](G-052-migrate-knowledge.md), mapped by
[G-069](G-069-migration-map.md)); stable placement applies from then on.
Until a first release Grove keeps no backward compatibility: only the current
schema is read, and an old commit is inspected with the CLI it carries.

Retrieve context by activity: shaping starts with the brief and relevant
knowledge; preparation with selected work and affected interfaces; execution
with its mandate and current task; review with acceptance and the candidate;
resume with a checkpoint and changed inputs. Links enable discovery without
preloading everything. A context response is facts, never authorization.

## Workflow and authority

Work lifecycle: **Proposed → Active → Review → Done**, with Abandoned
available only through an explicit human decision for now. Preparation,
implementation and agent review are activities; waiting, failures and process
state are additional facts. A terminal attempt does not automatically enter Review.

For implementation, Done means accepted and integrated into the configured
target ([G-038](G-038-review-lifecycle.md)); a `done` record without a
candidate predates that meaning and keeps its historical evidence.
Research/design work needs a completion condition suitable to its deliverable.

An assignment or Implement action authorizes bounded execution; merely creating
a proposal or editing status does not. Preparation investigates technical
unknowns and surfaces consequential human decisions before implementation.
Routine technical choices proceed within the mandate. Separate plan-file or
stage-by-stage human approval is not a universal gate.

Independent agent review examines a candidate against acceptance and evidence.
Human review presents changed behavior, decisions, outstanding issues and checks
before optional diffs. Approval is tied to the candidate being integrated;
changed candidates require reconsideration. Feedback that starts another
implementation attempt returns work to Active and preserves prior reviews.

Interactive and headless callers share the workflow, with explicit human
availability, authority and resource bounds. Missing human decisions become
durable questions; unchanged waits do not trigger repeated work. Unattended
research publishes proposals on an isolated reviewable branch, with supporting
knowledge identified for selective integration, and cannot authorize its own
implementation or merge.

## Current view and TUI

The normal view is project-wide and independent of the invoking checkout.
Present current work first, collapse identical observations, place demonstrably
superseded states in history, and label unintegrated and uncommitted changes.
Genuine divergence stays visible. Git ancestry supports this view; timestamps,
status rankings, or the newest branch tip do not define authority.

This supersedes G-002's explicit-versions-first presentation as the default.
Exact source inspection and fresh workspace binding remain available.
Ancestry here means merge bases: a copy is superseded when, since it and
another split, only the other changed the record. The integration target is
optional configuration (`target: main` in `grove.yaml`), chosen by the owner
on 2026-09-22. It labels progress as not yet on the target and never decides
which state is current.

Board columns follow the lifecycle, with bounded recent Done items,
searchable older work, and Abandoned hidden by default. Detail leads with
rendered work Markdown, with artifacts, workflow indicators and a timeline;
source observations are secondary. Review leads with a concise result and
progressively exposes evidence, files and diffs. Design these views together;
visual quality requires the owner's judgment in an actual terminal.

## Adoption and execution

`grove init` bootstraps another repository: valid configuration, a
placeholder brief, the record root and managed entrypoints, preserving existing
instructions. A shaping session develops the brief. Preserve provenance and
repair current references; do not retain parallel editable layouts. Product
and workflow instructions keep their functional homes. Live changes to a
sibling repository need a separately assigned scope, and an explicit cutover
preserves identity, knowledge and evidence with one editable authority for
each fact.

One bounded implementation runs independently of the TUI as a Grove-owned
`claude -p` process that continues after the terminal closes, reconnects
without duplication and stops explicitly
([G-101](G-101-attempt-mechanism.md) records the choice and when to
reconsider it). No machine-reboot guarantee is selected. Roles start with
shaping/preparation, implementation and independent review; route model
strength by uncertainty and consequence and retain actual configuration per
attempt. Bound batches and fan-out before adding schedules.

Reports and reviews should support evidence-backed process learning. A bounded
debrief proposes changes to knowledge, guidance or mechanical checks; it does
not silently rewrite acceptance or turn every observation into permanent policy.

## Observed state

At main `768efab` (2026-09-23), the Go CLI inspects, creates, updates and
converts records; shows versions across branches and worktrees and resolves
their checkouts; assembles staged context; bootstraps a repository with
`init` and prints its embedded guides; records approval and feedback and
integrates an approved candidate; and runs, lists and stops headless attempts.
Bare `grove` opens the board on the current view, with record detail, search,
the review actions and attempts. Terms, plans, reviews and pages are records.
This is a dated observation, not a progress log: work records own what has
changed since.

The installed `grove` is a build of this CLI from a named commit and can lag
this checkout; use `go run ./cmd/grove` here. The predecessor is uninstalled.
The archived FastAPI/PostgreSQL application
(the sibling checkout `grove-archive-2026-09-18/`, last commit `be40e46`) is
historical, with no service,
credential, deployment or backlog authority over this project.

## Suggested sequence

The preview phase has three proposals, each independently assignable and none
a prerequisite of another: behavioral evaluations of context and workflows
([G-108](G-108-workflow-evals.md)), usable Attempts screens
([G-109](G-109-attempts-usability.md)), and distribution preparation
([G-110](G-110-external-preview.md)). Their order is proposed: the owner has
not selected one. Completing them is not by itself evidence that a preview is
ready; the owner judges that.

The brief does not track progress or a current next action. Each work
record's Next owns that; change this section only when the selected sequence
itself changes. The adoption [roadmap](G-047-adoption-roadmap-plan.md) and
[evaluation](G-048-direction-evaluation-review.md) are the closed milestone's
history.
