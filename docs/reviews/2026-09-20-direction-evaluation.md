# First-days evaluation and research

Captured 2026-09-20. This report preserves observations and recommendation
rationale from the direction review; it is not another product specification.
The owner subsequently approved the recommendations in
[D-004](../../grove/decisions/D-004-interactive-adoption.md). The
[brief](../restart-brief.md) owns current direction and the
[roadmap](../plans/W-018-adoption-roadmap.md) owns the selected investment order.

## Scope and evidence limits

Inspected current Grove at main `c9904eae76be`, Bench at
`e39848aabcc9934477440f5b68b9c7afa649aacb`, and sibling skills at
`ec87bb2eb19a350634c4976cc5d215a57deb58de`. Read current code, records,
instructions and recent history; reviewed Bench shaping/debrief and nullsec
team definitions, sibling workflow sources, and existing restart investigations.
The archived application's architecture is historical; no service was accessed.

Retrieved skills/nullsec knowledge through the predecessor CLI. `grove status`
succeeded in skills. In nullsec it failed because it tried to create
`.git/grove-owner` outside the writable scope; no escalation or sibling repair
was attempted. `find` and `context` returned useful evidence, but deliberately
small context budgets reported required omissions. These calls did not establish
complete nullsec context or migration readiness. Inventory again at W-023.

The new CLI's `list` reported twelve Done work items and W-011 Proposed;
`check` passed with seventeen records before this reconciliation. These are
record observations and structural validation, not fresh execution of the full
Go suite or reacceptance of delivered features. No product code, agent run,
migration, visual trial or runtime recovery experiment was performed in this
evaluation. The installed Claude CLI reported `2.1.278`; only its version/help
were exercised. Sibling repositories were not intentionally modified.

## What the iterations contribute

| Iteration | Useful evidence to retain | Boundary |
| --- | --- | --- |
| Bench | Brief, domain atoms, linked knowledge, shaping/debrief, role mandates and exact-revision review receipts | Blueprint/export/run machinery is historical, not a required taxonomy or protocol |
| Archived Grove application | Separation of knowledge, work, execution and human review responsibilities | PostgreSQL/service/deployment authority was abandoned for this restart |
| Sibling skills | Research and proportional preparation, bounded review, checkpoint recovery, close/reconciliation | Installed executable and schema are different; no automatic compatibility |
| Current Go Grove | File contracts, allocation, safe updates, cross-worktree observations, staged context and terminal foundation | Authoring, review, adoption and human information architecture remain incomplete |

The evaluation's conclusion was to retain the foundation and prioritize a
complete useful loop. Building a board alone would leave the authoring and
review handoff incomplete; launching agents first would automate those gaps.

## Documentation and observed drift

- `docs/restart-brief.md` had 496 lines, combining direction, predecessor
  investigations, implementation history and stale next actions. At the observed
  main commit W-010 and W-012 were done and merged, while the brief still asked
  for W-012 integration and described W-010 as active.
- `docs/plans/W-010-W-011-agent-handoffs.md` had 518 lines. W-010 had shipped and
  its checklists were explicitly historical; W-011 was not started. Assigning
  them together again would risk repeating work or importing obsolete scope.
- Q-001 retained the earlier explicit-version presentation policy. The owner
  now wants a project-wide current view. This is a changed product choice,
  not evidence that the old implementation violated its accepted specification.
- `internal/tui/view.go` rendered history before metadata and raw record text.
  The owner now wants rendered work content first, with history and sources
  secondary. The source inspection establishes hierarchy, not visual quality.
- `internal/project/metadata.go` already allowed Abandoned; the execution guide
  already allowed small work without a separate plan. The missing pieces are
  transition authority, review handoff, integration meaning and presentation.

The previous brief remains recoverable from Git at `c9904ea:docs/restart-brief.md`.
Keep detailed historical reports linked from completed records; do not preload
them or copy the old brief into a second editable direction document.

## Knowledge and context

Bench's contract gave each kind of fact one owner, with terms, decisions,
questions and other atoms linked rather than restated. Its shaping and debrief
instructions distinguished consequential decisions from technical details and
carried evidence back into the owning knowledge. Those practices are useful
without requiring every type or a separate handoff-bundle repository.

Matt Pocock's current `grill-with-docs` entrypoint delegates to grilling and
domain modeling. Its modeling instructions test terminology against scenarios
and code, create files as needed, keep the glossary free of implementation
details and reserve ADRs for consequential trade-offs. These support Grove's
term/decision separation; they do not determine its storage schema.
[Entrypoint](https://github.com/mattpocock/skills/blob/main/skills/engineering/grill-with-docs/SKILL.md),
[modeling source](https://raw.githubusercontent.com/mattpocock/skills/main/skills/engineering/domain-modeling/SKILL.md).

Recommended retrieval discipline: shape with brief/relevant knowledge; prepare
with selected outcomes and interfaces; execute with task-local references;
review with acceptance/candidate/evidence; resume with checkpoints and changed
inputs. The existing read-only `context` command's explicit inclusion is a
foundation. A future term/artifact model must not silently preload the wiki.

The current reader rejects Markdown outside its three type folders. A proposed
`grove/brief.md`, terms folder or artifacts tree therefore requires schema work
before file relocation. The present docs/plans and docs/reviews stay valid homes
until the supported migration is delivered.

## Workflow and roles

The current work guide already prepares missing plans, permits small-work
exceptions, persists unanswered human questions, and caps review fix rounds.
It distinguishes implementation, review, owner acceptance and integration in
prose. W-020 makes the human-review handoff explicit while retaining this useful
authority model. A proposal alone is not an implementation assignment.

Bench's nullsec planner/reviewer definitions included mandate, expected output,
model class, write restrictions and escalation conditions. The reviewer required
exact contract/candidate identity and did not edit the implementation. Retain
these concepts with a small role set; evaluate extra reviewers at consequential
boundaries instead of mandating a roster for every task.

Model-routing recommendation, to validate through actual work: use strong
reasoning for uncertain shaping, consequential preparation and difficult review;
use a capable implementation model for bounded coding; use faster models for
well-bounded retrieval or mechanical work when measured quality permits. Record
the actual model and configuration. These are hypotheses, not comparative
performance measurements or fixed provider/model-version assignments.

Claude supports named agent configurations and per-agent model, tools and
permission controls. Map Grove mandates onto demonstrated harness capabilities;
do not assume identical controls in every adapter. The first adoption supports
interactive Claude/Codex; managed execution initially investigates Claude.
[Claude subagent documentation](https://code.claude.com/docs/en/sub-agents).

## Current state and presentation

At the observed main revision, `versions W-012` showed five rows: Done on main
and W-012 committed tips, Proposed on W-010's committed tip, and unchanged Done
live copies in main and W-012. The owner's earlier six-row example included
another live checkout no longer in that inventory. The older branch state is
real evidence but should not dominate the ordinary current-work experience.

`internal/versions.Version` is a source observation plus record/content revision.
A current projection can consume those observations without replacing their
inspection contract. Git ancestry can establish reachability, but unrelated
branch-tip commits, dirty files, reverts and competing edits require explicit
record-level policy. Time alone supplies no authority. W-024 must investigate
these cases rather than promising an always-unambiguous latest version.
[Git ancestry reference](https://git-scm.com/docs/git-merge-base).

Proposed visual implementation ingredients reviewed: Bubble Tea for lifecycle,
Lip Gloss for layout/styles, Bubbles for components, Glamour for Markdown and
Huh where forms improve setup/edits. Charm Log is a logger, not a provider event
protocol or durable attempt store. Check compatible module versions during
implementation and preserve untrusted-text escaping and on-demand history.
[Bubble Tea](https://github.com/charmbracelet/bubbletea),
[Lip Gloss](https://github.com/charmbracelet/lipgloss),
[Bubbles](https://github.com/charmbracelet/bubbles),
[Glamour](https://github.com/charmbracelet/glamour),
[Huh](https://github.com/charmbracelet/huh),
[Log](https://github.com/charmbracelet/log).

## Future runtime boundary

The observed local Claude help supports `--agent`, `--output-format stream-json`,
partial-message output and native background session commands. The streaming
flag is not `--stream-json-output`. Official documentation treats native
`--background` and print mode as different invocation paths; their interfaces
must be compared, not combined blindly.
[Claude CLI reference](https://code.claude.com/docs/en/cli-reference),
[programmatic use](https://code.claude.com/docs/en/headless).

Compare native session management with a Grove-owned process against a concrete
contract: identity, ownership independent of the TUI, bounded resources,
duplicate-start protection, durable events/result, interruption, reconnect and
Stop. A tmux session supplies neither acceptance nor result reconciliation by
itself. Retain raw logs separately from stable report/review artifacts.

The earlier [runner research](2026-09-19-shaping-and-runner-evidence.md) lists
failure probes and predecessor mechanisms; the
[work review](2026-09-19-predecessor-work.md) records visible adapter weaknesses.
Read those at W-027 preparation, then recheck current capabilities. No mechanism
or unattended reliability was demonstrated by this evaluation.

## Documentation reconciliation verification

The reconciliation branch is `worktree-direction-reconciliation`, based on
`c9904ea`. All new records were created with `go run ./cmd/grove new`; decision
status, priorities and relationships were set with revision-checked `update`.
No Go source or runtime behavior changed.

- `go run ./cmd/grove check` passed with 29 records. Read-only `context` calls
  for W-011, W-018, the adoption sequence and later TUI/runtime selections
  succeeded; ordering follows actual prerequisites, not roadmap priority.
- A local Markdown check verified 94 outgoing links and 21 incoming links to
  changed documents, including fragments. `git diff --check` was clean.
- An independent read-only agent reviewed the diff and new documents for
  lost intent, contract contradictions, dependency traps and misleading scope.
  It found one ownership gap: W-022 expected brief-location support that W-019
  did not explicitly require. W-019 now owns brief discovery/location and
  preservation of one authority and existing links. The reviewer found no other
  consequential issues; external sources and real harness behavior were outside
  that review. The correction received a focused follow-up review.
- Full Go tests and fresh harness/TUI trials were not run for this documentation
  change. The native and external workflow trials remain future acceptance.

This verification concerns the durable direction and assignments, not delivery
of their proposed capabilities. The main checkout's unrelated `lefthook.yml`
was preserved; no merge, push or sibling migration was performed.
