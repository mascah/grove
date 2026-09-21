# Grove brief

Direction reconciled 2026-09-20. This is the single current source of product
intent. The owner approved the evaluation and adoption sequence, then asked to
persist them for execution without returning to the conversation.
[D-004](../grove/decisions/D-004-interactive-adoption.md) records that authority
and the policies it revises. Work records own acceptance; the
[record model](record-model.md) describes the currently implemented schema.

## Purpose

Grove is a local project workspace where a person shapes intent with agents,
delegates bounded work, and returns to a clear review of the result. It connects
project knowledge, work, and execution evidence so the next human or agent can
continue with **the right context, at the right time**.

The first audience is the owner working across local hobby repositories and
worktrees. The first adoption milestone is a complete interactive
shape → implement → review → integrate loop on real nullsec work. Keyborg is
the second adoption test. Background execution follows the demonstrated loop.

## Selected foundations

- Go CLI, ordinary Markdown with YAML frontmatter, `grove.yaml`, configurable
  record storage, and Git. Core inspection needs no running service.
- Short sequential IDs coordinated across linked worktrees; revision-checked
  mutations. Separate clones still require collision checks and reconciliation.
- Humans and agents use the same records. Deterministic validation, retrieval,
  mutation and lifecycle mechanics belong in software; judgment belongs in
  instructions and attributable human or delegated decisions.
- Branch-local editing with a combined project view. Actions bind to an exact
  checkout and revision; browsing should not make source selection the main job.
- Bare `grove` opens the TUI; explicit commands remain noninteractive. The TUI
  should be visually appealing, keyboard-friendly, and use the Charm ecosystem.
- Grove-owned shared workflows with thin Claude and Codex adapters. Replace the
  sibling skills through demonstrated adoption, preserving useful responsibilities.

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

Restore terms and discoverable artifacts early. Plans need not become an
independent ticket lifecycle. Small work may carry its preparation in the work
record; substantial work gets a concise linked plan. Designs, flows, systems,
research, releases, builds and team definitions remain possible extensions,
introduced when information needs their distinct home.

Retrieve context by activity: shaping starts with the brief and relevant
knowledge; preparation with selected work and affected interfaces; execution
with its mandate and current task; review with acceptance and the candidate;
resume with a checkpoint and changed inputs. Links enable discovery without
preloading everything. A context response is facts, never authorization.

## Workflow and authority

Target work lifecycle: **Proposed → Active → Review → Done**, with Abandoned
available only through an explicit human decision for now. Preparation,
implementation and agent review are activities; waiting, failures and process
state are additional facts. A terminal attempt does not automatically enter Review.

For implementation, Done will mean accepted and integrated into the configured
target. Research/design work needs a completion condition suitable to its
deliverable. The existing schema still permits branch-local Done before merge;
W-020 must migrate that meaning explicitly and preserve historical evidence.

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

This supersedes Q-001's explicit-versions-first presentation as the future
default. Exact source inspection and fresh workspace binding remain available.
W-024 owns the resolution policy before implementation; the current board is
still checkout-scoped until that work ships.

Board columns follow the target lifecycle, with bounded recent Done items,
searchable older work, and Abandoned hidden by default. Detail leads with
rendered work Markdown, with artifacts, workflow indicators and a timeline;
source observations are secondary. Review leads with a concise result and
progressively exposes evidence, files and diffs. Design these views together;
visual quality requires the owner's judgment in an actual terminal.

## Adoption and later execution

Bootstrap should create valid configuration and minimal harness entrypoints,
preserve existing instructions, and create content directories as needed. A
shaping session develops the brief. The target layout groups brief, knowledge
and artifacts under the configured Grove root; schema support and migration
must precede moving these files there.

Use one unambiguous binary/workflow version for the nullsec pilot, rehearse
migration in a disposable copy, then make an explicit live cutover. Preserve
identity, knowledge and evidence; do not leave two editable authorities for the
same facts. Live sibling changes need a separately assigned migration scope.

Later, one bounded implementation can run independently of the TUI, continue
after it closes, reconnect without duplication, and stop explicitly. Compare
supported native Claude background sessions with a Grove-owned `claude -p`
process; tmux is an option, not a requirement. No machine-reboot guarantee is
selected. Roles start with shaping/preparation, implementation and independent
review; route model strength by uncertainty and consequence and retain actual
configuration per attempt. Bound batches and fan-out before adding schedules.

Reports and reviews should support evidence-backed process learning. A bounded
debrief proposes changes to knowledge, guidance or mechanical checks; it does
not silently rewrite acceptance or turn every observation into permanent policy.

## Observed state

At main `c9904ea`, the Go CLI supports inspection, creation, safe field updates,
versions, workspace resolution and staged context. The read-only Bubble Tea v2
board and on-demand card history are integrated. W-010's work guide and local
adapters are done and integrated; W-011 shaping is not implemented. Terms,
structured artifacts, Review status, portable setup, current-state projection,
review actions and managed attempts are future work. `abandoned` already exists.

The installed `grove` is still the sibling predecessor. Use `go run ./cmd/grove`
here. The archived FastAPI/PostgreSQL application is historical, with no service,
credential, deployment or backlog authority over this project.

## Suggested sequence and current next action

[W-018](../grove/work/W-018-interactive-adoption.md) owns the first adoption
milestone. Its [roadmap](plans/W-018-adoption-roadmap.md) orders independently
assignable work and later investments; the
[evaluation](reviews/2026-09-20-direction-evaluation.md) retains findings,
research and limits rather than another editable product direction.

**Next:** integrate this documentation reconciliation, then assign
[W-011](../grove/work/W-011-shaping-entrypoint.md). Follow with knowledge/artifact
support and the review handoff, exercise a real Grove loop, package portable
setup, and pilot one real nullsec change. Current-state and visual TUI work
follow the pilot; durable execution follows that useful review experience.
Investment order alone is not a technical dependency or authority to execute
the entire roadmap.
