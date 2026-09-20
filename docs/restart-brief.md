# Grove restart brief

Captured: 2026-09-18.
Status: product direction and starter file defaults selected; the Go CLI
reads, validates, creates, and updates the operational records, shows their
versions across local branches, and locates a selected version's checkout.
A read-only terminal board over those operations opens from bare `grove`; the
owner accepted it as a starting point on 2026-09-19 and it is integrated.
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
  2026-09-18. The W-009 branch uses Bubble Tea v2 for the terminal board.
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
  owner subsequently selected the Kanban board as the first terminal experience
  after reliability repairs, with version choices inside each card rather than
  a separate picker as the starting screen.
  [W-009](../grove/work/W-009-terminal-picker.md) designs that board-first TUI:
  one checkout supplies column statuses, while each card exposes explicit
  cross-branch versions. The owner selected bare `grove` (no arguments) to open
  the TUI, with the board as its first screen; there will be no `grove board`
  subcommand. Worktree creation and agent launching remain later
  investments. The board interaction is implemented and integrated
  ([evidence](reviews/2026-09-19-board-W-009.md)). After a first demo the owner
  had identical versions folded and the wording changed, then accepted it as
  good enough for the moment; boundaries, colour, and lineage remain wanted.
- Eventual agent execution from that workspace, with `claude -p` as the concrete
  first-provider idea. Exact invocation and lifecycle behavior need validation.
- Lifecycle choice settled on 2026-09-19: running agent sessions continue when
  the TUI closes; reopening reconnects, and Stop is a separate action. This does
  not promise automatic recovery after machine shutdown. The run owner must be
  independent of the viewing TUI; tmux is a candidate, not a selected dependency.
- Near-term dogfooding need identified on 2026-09-19: repeated requests for
  implementation prompts expose the missing restart-native `/work` equivalent.
  The [shared execution guide](work-execution.md) and saved handoff prompts now
  provide reusable instructions. [W-010](../grove/work/W-010-work-handoffs.md)
  proposes a thin work entrypoint and deterministic context/prompt preparation
  before supervised launching. A future Kanban Implement action can reuse that
  prepared assignment; launch, cancellation, recovery, and result reconciliation
  still need their own contract. This does not add agent launching to W-009.
  The [actual predecessor `/work` review](reviews/2026-09-19-predecessor-work.md)
  is required input to W-010: its preparation, review, recovery, and closure
  responsibilities extend beyond prompt generation. Adapt or explicitly defer
  them; the current guide is only a partial baseline.
- Emerging authoring/execution workflow, described by the owner on 2026-09-19:
  interactive Claude/Codex sessions create and refine proposed work, questions,
  and decisions through the CLI; bounded scheduled/triggered headless research
  can also propose work. Moving proposed work to Active would offer launching a
  headless implementation session. Detailed launch policy is still proposed.
  Publication choice settled on 2026-09-19: unattended research publishes its
  proposals on a separate reviewable branch; selected proposals are integrated
  afterward. It does not write directly into the planning checkout or merge its
  own proposals automatically. Until integration, W-009 exposes branch-only work
  under Elsewhere (first named Other sources), or in Proposed when viewing that research checkout.
  A status edit alone does not authorize a launch or prove that
  one is running. [W-011](../grove/work/W-011-shaping-entrypoint.md) proposes the
  immediate interactive shaping guide; W-010 owns the reusable work assignment.
  Scheduling and run supervision remain later work, outside W-009. Inspect the
  [predecessor/Bench evidence](reviews/2026-09-19-shaping-and-runner-evidence.md)
  before designing those mechanisms. Streamed activity is not work acceptance.
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
precisely. Optional agent execution needs a run owner that outlives the viewing
UI; neither ordinary record access nor opening the board requires that owner.

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
were not tested by that early probe. W-004/W-005 subsequently implemented the
read-only foundations; the review below qualifies their current reliability.

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
   recovery, and reconciliation. TUI exit now explicitly leaves runs working;
   define ownership and reconnection independently of the UI process.
   Reuse existing export/run/reconcile ideas only after inspecting their fit.
3. **How should the first TUI feel in use?** The owner selected the terminal
   experience and requested early Kanban value. W-009 proposes the checkout
   board and explicit version view, implemented with Bubble Tea v2.0.9. The
   owner's first answer (2026-09-19, recorded in W-009): an acceptable start;
   a list of versions across branches does not say what happened to an item,
   and lineage over time is what they want from a card (W-012). Agent
   execution is outside the first interaction.

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

**Next action:** decide with the owner whether a card opens to lineage or
keeps the version list beside it, then plan
[W-012](../grove/work/W-012-card-lineage.md). The board
([W-009](../grove/work/W-009-terminal-picker.md)) and its load fix
([W-013](../grove/work/W-013-load-scaling.md)) are done and integrated; run
the board with `go run ./cmd/grove`. The CLI repairs below were integrated
first (main `acfc905`, established by Git ancestry). The board's
[evidence](reviews/2026-09-19-board-W-009.md) separates what automated checks
and two independent reviews established from the owner's judgment, which
W-009 records. The [2026-09-19 review](reviews/2026-09-19-integrated-cli.md) verified W-003
integration at `b20d2b0` and W-004/W-005 integration at `5041ae1`, and
demonstrated wrong-repository routing through a symlinked project prefix,
success for a checkout deleted during inspection, comment loss and
valid-update refusals, an ignored configuration change, and newline-path
parsing that wrote a lock outside the actual Git directory.

Those defects are repaired on branch `worktree-W-006-W-008` (base `2d6de36`):
[W-006](../grove/work/W-006-workspace-provenance.md) verifies live project
ownership through every prefix component and re-checks the selected target
before returning a workspace; [W-007](../grove/work/W-007-preserve-updates.md)
preserves comments, edits explicit keys, plans flow separators for the whole
request, and compares exact configuration bytes in update and creation;
[W-008](../grove/work/W-008-git-paths.md) reads every Git path on its own and
shares one NUL-delimited worktree inventory with the allocator. Each failure
was reproduced on the base before its fix. The
[repair evidence](reviews/2026-09-19-repairs-W-006-W-008.md) holds commits,
suite results, independent review dispositions, and remaining limits. These
repair existing contracts; they add no product decision. Done records assert
completion on that branch: establish integration by Git ancestry, not by
status, and keep it a separate, explicit step. Preserve the retained
W-004/W-005 worktree; it has no commits missing from main.

The [W-009 handoff](prompts/W-009-implementation.txt) and the
[W-006–W-008 handoff](prompts/W-006-W-008-implementation.txt) are both spent.
Both reference the shared execution guide; the work records and plans remain
the specifications. W-010/W-011 below make future assignments need only work
IDs; they are not a prerequisite that delays the board.

Requirements work continued on main while the repairs ran in their separate
worktree; keep such work outside that checkout until it is integrated. Refine
W-010's shared interactive/headless assignment and W-011's authoring guide.
Apply the selected separate-branch publication policy to unattended research
handoffs; identify the proposals and supporting questions/decisions needed for
selective integration. W-009's existing source-scoped columns and Elsewhere
shelf remain unchanged.

The owner agreed to W-010/W-011 as the next bounded dogfooding investment:
shape records, assign selected work through short entrypoints, and recover the
next action from repository evidence. The
[shared plan](plans/W-010-W-011-agent-handoffs.md) and
[Fable handoff](prompts/W-010-W-011-implementation.txt) specify context assembly
and thin Claude/Codex adapters. Their interface details are prepared design,
not implemented capabilities. Prefer following the owner's W-009 handoff to
avoid shared CLI/docs edits; W-009 is not a semantic dependency. After proving
this loop, shape one manually launched supervised headless attempt before
schedules or board launch controls.

The selected next user outcome is a terminal Kanban board
with explicit version selection inside each card. On 2026-09-19 the owner
preferred the board over a standalone version picker. Version selection supports
choosing the correct workspace when branches differ; the board supplies the
everyday overview. Proposed policy: columns
show one labelled live checkout's statuses, card details group all versions by ID,
and an Elsewhere shelf exposes work absent from that checkout. Workspace
selection remains explicit, incomplete results visible, and stale selections
require refresh. This brings board value into the first TUI; it does not select
one authoritative status across branches. [W-009](../grove/work/W-009-terminal-picker.md) owns
the proposed screen, keyboard/output contract, acceptance, and linked plan.
Bubble Tea v2.0.9 is used for the board only; adopting it
more widely follows the owner's judgment of the board. After the owner's first
demo on 2026-09-19 the card folds identical versions into one row and the board
says branch and checkout; the next card outcome the owner asked for is a work
item's lineage from Git history ([W-012](../grove/work/W-012-card-lineage.md)).
A full load was measured at 27.7 s with 300 branches and is now 0.35 s
([W-013](../grove/work/W-013-load-scaling.md)). Keep worktree creation
as a separate design with
explicit destination and failure policies; agent launches, claims, and mutation
through an interactive view come later. Q-001's accepted policy stays unchanged.

The implemented CLI remains useful for dogfooding: `list`, `show`, `check`,
`new`, `update`, `versions`, and `workspace`, subject to the limits recorded
in the review and the repair evidence. Use `go run ./cmd/grove new` for every new record and `update` for field
changes; keep the installed sibling CLI untouched. Separate-clone and import
conflicts remain outside the local allocator guarantee.

## Reset and scope record

The old local application was archived intact as `../grove-archive-2026-09-18/`,
including Git history and ignored local data, from clean main at `be40e46`.
The fresh `../grove/` began with this brief, README, agent guidance, and a
minimal ignore file. Its first implementation is now the local Go inspection CLI.

The old application's PostgreSQL/API architecture, deployed Mini service,
planning-authority cutover, and old roadmap are historical. This local reset
does not alter any external service, database, credential, remote repository,
or installed CLI. The sibling skills and nullsec projects remain unchanged.
No runner has been launched. The board uses Bubble Tea v2, the only TUI
framework in the module. The new
Markdown/frontmatter format and file defaults are selected; `grove.yaml` and
operational records exist, and the inspection, creation, update, versions,
and workspace CLI is implemented and verified locally. It has not replaced
the installed sibling executable.
