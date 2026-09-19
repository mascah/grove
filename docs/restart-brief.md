# Grove restart brief

Captured: 2026-09-18.
Status: product direction and starter file defaults selected; the Go CLI
reads, validates, creates, and updates the operational records, shows their
versions across local branches, and locates a selected version's checkout.
This document is the restart's current source of intent.

## Intent and authority

The owner proposed a Jira/Linear-like project tool implemented as a CLI, with
project state persisted locally as Markdown or another inspectable file format.
Humans and agents would coordinate through those records instead of requiring
an external tracker. Desired features include linked items, dependencies,
assignees, attachments such as an agent's implementation review, and project
configuration, potentially YAML. A TUI and/or optional local web server would
make the same files convenient for people to use.

The conversation connected this idea to the work concept in `../skills/` and
its use in `../nullsec/docs/grove/`. The owner then imagined a TUI showing work,
questions, research, and related knowledge, with actions that launch
`claude -p` sessions to advance work. This exposed a possible harness role for
Grove in addition to project tracking.

The owner explicitly endorsed the direction and requested a fresh start from
the first `../grove/` application, plus durable documentation so the idea could
survive conversation or usage limits. They see Grove becoming its own CLI
project, with the skills in `../skills/` potentially changing substantially or
eventually disappearing. That possibility is not a decision to delete the
skills now or to replace all reasoning guidance with code.

This is authority to capture and develop the restart direction. The detailed
design suggestions below are not an approved implementation specification.

## Product direction

Grove should give a person and their agents a shared, inspectable account of
what the project wants, what is known, what is being worked on, what happened,
and what remains unresolved.

The initial audience is the owner working with coding agents across local
repositories and worktrees. Broader team and cross-machine collaboration are
possibilities, not established first-release requirements.

Selected direction:

- A standalone Grove CLI project, with durable local project records.
- Go for the first CLI implementation, explicitly selected by the owner on
  2026-09-18. A TUI framework has not been selected.
- Core use without a required hosted server or persistent daemon.
- Human and agent access to the same underlying project state.
- Linked work, dependencies, assignees, and attached evidence/artifacts.
- Project configuration in `grove.yaml`; the [record model](record-model.md)
  owns the accepted starter defaults.
- An interactive workspace that can expose work, questions, research, and
  related project knowledge. A TUI is the currently exciting direction; a local
  web UI is also possible, without committing to two initial interfaces.
- Working direction: branch-local project records with a combined cross-branch
  view. The owner accepts editing a record in its branch context for now,
  provided Grove makes reaching that context low-friction. Opening the existing
  worktree as the editing/execution context is accepted interaction direction;
  detailed source and routing contracts remain to be designed.
- Next experience selected on 2026-09-19: coordinate work across local branches
  by seeing record versions and opening the right workspace. This guides work
  after safe record updates. The owner also selected grouping by record ID,
  showing each branch's status, and requiring explicit version selection before
  opening a workspace; [Q-001](../grove/questions/Q-001-branch-versions.md) retains
  that answer. W-004 and W-005 implemented the CLI foundations on
  2026-09-19: `versions` and `workspace`, with the selector contract in
  their shared [coordination plan](plans/W-004-W-005-coordination.md). The
  interactive Open workspace action, creating a worktree for a branch
  without one, a checkout-local board, and agent launching are later
  investments.
- Eventual agent execution from that workspace, with `claude -p` as the concrete
  first-provider idea. Exact invocation and lifecycle behavior need validation.
- A clean implementation start, informed by the working skills and nullsec
  experience rather than constrained by the first application's architecture.

Starting scope accepted on 2026-09-18: use work, questions, and decisions as the
minimum to begin dogfooding. Defer implementing structured attachments such as
reviews and reports; ordinary prose and links can support early development.
The owner subsequently accepted the [starter record model](record-model.md):
Markdown with YAML frontmatter, the four required fields, optional typed
relationships, explicit lifecycles, and the initial validation boundary.
On 2026-09-19, the owner selected configurable root-level `grove/` storage and
accepted the starter file defaults. Actual browsing then exposed excessively
long filenames. The owner accepted replacing the random-ID/timestamp naming
trial with short type-prefixed sequential IDs, allocated under a shared lock in
Git's common metadata directory across local branches and worktrees. Separate
clones remain outside that coordination guarantee. [D-002](../grove/decisions/D-002-sequential-ids.md)
records the revision; the [record model](record-model.md#on-disk-contract) owns
the current schema and allocation requirements. Operational records now use
names such as `W-001-inspect-records.md`, with dates retained in frontmatter.
W-002 implements creation and the shared allocator accepted in
[D-003](../grove/decisions/D-003-allocator-mechanism.md).

The owner then proposed carrying richer work metadata from nullsec W-032 into
the starter model: members, dependencies, priority, size, and kind, to support
the UI and project functionality, then accepted the optional fields and their
semantic distinctions. The [record model](record-model.md) owns those meanings.
The initial allowed values were accepted with the file defaults on 2026-09-19;
this does not adopt the predecessor's preparation or execution policies.

“Stateless CLI” means commands need no resident application process. Project
records and recoverable coordination necessarily have state. File-backed,
local-first, and no required daemon express the intended properties more
precisely. Whether an optional running UI supervises launched jobs is open.

## Why this is worth exploring

Files keep the working record available to editors, agents, version control,
and future interfaces. The interface can expose both the intended work and the
reasoning/evidence surrounding it. Agent execution can then operate against an
explicit mandate and return a result that the project retains.

Backlog.md is a close existing project: Markdown tasks, CLI, terminal board,
local browser, dependencies, and configuration. Its documented spec/plan/review
workflow overlaps with Grove's existing lifecycle, which the owner identified
as a reason integration may be less neutral than it first appears. The owner
found it inspiring and expressed interest in investing more in Grove's CLI.
No hands-on comparison established a missing Backlog.md capability.

The product hypothesis is that Grove's connected project knowledge, practical
human/agent coordination, and evidence-aware execution can justify building
further. A Markdown board alone is not an established differentiator.

Sources inspected in this conversation:

- [Backlog.md](https://github.com/MrLesk/Backlog.md)
- [git-bug](https://github.com/git-bug/git-bug), another offline/distributed
  tracker using Git storage.

## Evidence from the working Grove workflow

The sibling `skills` repository already implements a Grove CLI with `status`,
`context`, `find`, `batch`, `close`, `export`, `run`, and `reconcile` commands.
Work distinguishes capability scope, release membership, dependencies,
preparation, acceptance, and evidence. Skills provide exploration, shaping,
implementation/review, closure, and knowledge maintenance guidance.

The CLI context retrieved for skills W-013 and the status output expose concrete
pressure points. These are recorded observations and proposed fixes, not claims
that those fixes have shipped:

- W-013: main cannot see ownership across linked worktrees. Proposed shared
  local claims should prevent duplicate acquisition and distinguish ownership
  from checkout-local progress or completion awaiting integration.
- W-014: completed code was merged twice while its work page remained proposed.
  Proposed handoff validation checks committed work closure and evidence.
- W-010/W-011: an overnight nullsec run exposed repeated idle notifications,
  duplicated reporting, uncollected background commands, and missing compact
  recovery checkpoints. These motivate identified attempts and durable results.
- W-012: passing component checks did not establish correct connected player
  interactions or the owner's usability judgment.
- W-015: concurrent planning has shared-checkout hazards; isolating it also
  raises visibility and independently allocated work-ID problems.

Nullsec provides real work with release membership, blockers, design artifacts,
research/spikes, implementation, verification, and human feedback. Its W-032
release was partially closed, with W-037 blocked on W-036, when status was read.
It is useful evidence for the proposed workspace, not a migration target yet.

Read evidence through each project's CLI and follow its AGENTS.md. Useful
starting points are `grove context --work W-013 --phase shape` in skills and
`grove status` in nullsec. Repository revisions observed during capture were
skills `82d4b16` and nullsec `49d2dee`; these projects continue to evolve.

## Proposed responsibility boundary

The CLI could own identity, relationships, validation, safe mutations, ownership,
context queries, evidence references, and result reconciliation. The interactive
UI would use those same operations. An optional execution adapter could run an
agent against recorded work and return an inspectable result.

Skills could become thinner instructions for exploration, judgment, planning,
implementation, and review. Evaluate each responsibility before moving it:
deterministic checks belong naturally in code; product judgment may still need
instructions and actual human input. Full skill elimination is an open outcome.

Avoid a second ticket that duplicates an existing Grove work record. A future
compatibility or migration design should give each fact one canonical owner.
Whether the existing schema is adopted, adapted, or replaced remains open.

## First experience to design

Proposed vertical slice, not yet an approved implementation scope:

1. Open one local project and see its work, blockers, questions, and evidence.
2. Select one prepared work item and inspect the outcome and constraints.
3. Launch one bounded agent attempt from the workspace, or perform the same
   operation from the CLI. Record which work and input revision it consumes.
4. Observe ownership and progress from another local session/worktree.
5. Interrupt the supervising interface or agent and reopen the workspace.
   Recover an honest account without silently launching duplicate work.
6. Inspect a returned result and a linked review tied to the candidate revision.
7. Distinguish implementation completion, required human judgment, and
   integration. Preserve the evidence needed by the next session.

Possible UI actions are Investigate, Prepare, Implement, and Review. They would
request work with a completion condition; a visual status move cannot itself
prove execution or acceptance. Actual command names and workflow flexibility
are undecided.

The first demo may be smaller than this full sequence. The sequence is a way to
test the architecture, particularly ownership and recovery, before expanding the
feature list. Tracking and inspection should remain useful without running AI.

## Consequential design questions

### Cross-branch view and branch-context editing, 2026-09-18

The owner wants a project-wide view and, after discussing Backlog.md's
branch-local storage and combined board, accepts its branch-context editing
restriction for now. The condition is low user effort: Grove should fit context
selection into the workflow. This establishes the working direction without
adopting Backlog.md's schema, version-selection rules, or exact interaction.

The owner correctly identified a checkout conflict: a branch already checked out
in a linked worktree cannot normally also be checked out in main's working
directory. [Git switch](https://git-scm.com/docs/git-switch) documents this guard.

Accepted interaction direction (owner follow-up: "that feels acceptable"):
a card offers Open workspace. Grove identifies the selected
record's branch and locates its existing checkout with
`git worktree list --porcelain -z`. The UI's editing context and commands target
that checkout explicitly; this does not require switching main. If no checkout
exists, a workspace-opening action could prepare a linked worktree. Ambiguous,
missing, or changed branch/worktree identities need explicit handling, without
forcing a second checkout or guessing which version of a work item to edit.

Proposed safeguards: revalidate the target branch and record revision before a
mutation, preserve dirty files/index state, and coordinate writes with an active
agent. A worktree's existence does not imply exclusive ownership or permission
to change an agent's assignment. Show whether a card reflects committed branch
data or live working files. Keep the selected version's source visible. The
owner subsequently selected explicit versions with no authoritative aggregate
status; Q-001 retains that answer. Source-discovery details and claim/handoff
semantics remain open.

A temporary Git 2.50.1 probe reproduced the duplicate-checkout refusal, located
the existing feature worktree through porcelain output, and edited its record.
Main's current branch, copy of the record, and unrelated staged changes remained
unchanged. The fixture was removed. This proves basic routing feasibility only;
concurrent writers, branch-change races, workspace provisioning, and recovery
were not tested. No product implementation exists yet.

Earlier alternative retained for reconsideration: a separate `grove-project`
branch with a shared working directory was also proven locally feasible in a
temporary repo. Main and two worktrees could resolve it through Git's common
directory and see board changes with independent history. The current working
direction instead explores branch-local records and cross-branch visibility.
Reconsider shared storage if navigating/editing branch-specific records proves
too disruptive. [Git worktree facilities](https://git-scm.com/docs/git-worktree).

### Remaining questions

1. **How does a selected version reach its workspace?**
   [Q-001](../grove/questions/Q-001-branch-versions.md) resolves presentation:
   group by ID, show each branch's status, and select a version explicitly.
   W-004 and W-005 deliver the CLI contracts: a selector per observed
   version and a `workspace` command that resolves it to an existing
   checkout or refuses with the reason. How an interactive view presents
   those rows, when a missing checkout is created, and claims and run-record
   lifetimes remain future design.
2. **What does one run promise?** Define its input, actor/attempt identity,
   working directory, output, failure/wait state, cancellation, interruption
   recovery, and reconciliation. Establish what happens when the TUI exits.
   Reuse existing export/run/reconcile ideas only after inspecting their fit.
3. **What is the next useful interactive interface?** Go inspection commands now
   read real Grove data. An interactive view could validate navigation; one
   executable action would test the harness boundary. Choose a TUI stack against
   that bounded slice when it becomes the next investment.

Related details to resolve at the appropriate boundary:

- Safe concurrent Markdown/frontmatter edits, schema versions, direct-editor
  behavior, and rebuildable indexes.
- Shared-counter initialization/recovery, clone/import collision handling,
  dependency and membership semantics, stale writes and conflict handling.
- Responsibility/assignee versus the particular session holding a claim.
- Attachments and review provenance: author, reviewed revision, result, and
  what makes a previous review stale.
- Configurable workflow versus Grove's useful defaults; avoid encoding a
  universal sequence before exercising investigations and implementation.
- Local linked-worktree coordination versus independent clones/machines.
  Disconnected clones cannot promise globally exclusive claims.
- How existing skills projects adopt the new executable without a name/path
  collision or two competing project records.

## Suggested sequence and current next action

1. Core record model and file defaults accepted: [the model](record-model.md)
   defines the schema and links to Grove's operational records. Keep work and
   question details in those records, schema details in the model, and direction
   in this brief.
2. Walk one work record through creation, implementation on a branch, review,
   integration, and reopening. Keep work, an execution attempt, and review
   evidence distinct. Resolve identity, version selection, and completion meaning.
3. The Go CLI now lists, shows, validates, creates, and updates the agreed
   records. Use actual Grove development records for dogfooding;
   keep the restart brief as the
   direction owner until an explicit migration avoids duplicate ownership.
4. Add the combined board and Open workspace interaction over those operations.
   Exercise concurrent edits and worktree changes. Extend to one agent run and
   interruption recovery after the record workflow is useful on its own.
5. Evaluate against a real nullsec scenario in a fixture/copy. Migrate skill
   responsibilities incrementally after the replacement demonstrates useful
   behavior. The existing skills may assist development without dictating the
   new product schema.

**Next action:** shape the interactive Open workspace experience over the
integrated CLI. [W-004](../grove/work/W-004-record-versions.md) and
[W-005](../grove/work/W-005-record-workspace.md) were integrated into main at
`5041ae1` on 2026-09-19: `versions` shows each record's committed version on
every local branch and live version in every worktree with a selector per
version, and `workspace --source` resolves a selected version to its existing
checkout or refuses with the reason, per the
[coordination plan](plans/W-004-W-005-coordination.md). The next investment
is presenting those version rows, explicit selection, and creating a linked
worktree for a branch that has no checkout, which the CLI deliberately
refuses today; choose the interface stack against that bounded slice. Keep
implementation serial while CLI and project-loader ownership overlaps.
Creation is complete in [W-002](../grove/work/W-002-create-records.md): use
`go run ./cmd/grove new` for every new record; the allocator per
[D-003](../grove/decisions/D-003-allocator-mechanism.md) shares one counter
across linked worktrees, and separate-clone and import conflicts stay out of
scope. Inspection is complete in
[W-001](../grove/work/W-001-inspect-records.md): use `list`, `show ID`, or
`check`. Keep the installed sibling CLI untouched; TUI and execution remain
later work.

## Reset and scope record

The old local application was archived intact as `../grove-archive-2026-09-18/`,
including Git history and ignored local data, from clean main at `be40e46`.
The fresh `../grove/` began with this brief, README, agent guidance, and a
minimal ignore file. Its first implementation is now the local Go inspection CLI.

The old application's PostgreSQL/API architecture, deployed Mini service,
planning-authority cutover, and old roadmap are historical. This local reset
does not alter any external service, database, credential, remote repository,
or installed CLI. The sibling skills and nullsec projects remain unchanged.
No runner has been launched and no TUI framework has been selected. The new
Markdown/frontmatter format and file defaults are selected; `grove.yaml` and
operational records exist, and the inspection, creation, update, versions,
and workspace CLI is implemented and verified locally. It has not replaced
the installed sibling executable.
