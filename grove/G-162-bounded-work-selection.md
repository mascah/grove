---
id: "G-162"
type: work
title: "Execute an explicitly selected set of work with bounded progress and review"
status: review
created: "2026-09-25T20:35:28Z"
updated: "2026-09-26T00:03:07Z"
kind: feature
relates_to: ["G-045", "G-046", "G-044", "G-101", "G-054", "G-056", "G-057", "G-059", "G-060", "G-161", "G-188"]
candidate: "2fb756b66f1f72c96be38db1d287fc1b12997e5d"
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

## Evidence

Implemented 2026-09-25 by a headless `/grove-work G-162` session on
`worktree-G-162`, from main `fe97300` (plan checkpoint `06015a9`), merged
with main `38f82aa` in `a11dd94`. Started from this record at
`sha256:70baf910…` and plan [G-185](G-185-g-162-selected-work-plan.md) at
`sha256:6584eb48…` (`dbb2f0d`), whose launch approved it. The plan's
"Adjustments during implementation" records each bounded change since.
G-163's answer is decision
[G-188](G-188-selected-work-shared-candidate.md) (`e25925e`).

Commits: `4d1507f` selection in `internal/attempt` and `run`; `5f44b78`
group approve, feedback and integrate; `8073de4` board; `99a143c` docs;
`be48722` and `15774e4` review fixes; `a11dd94` the merge of main.

Against each acceptance item:

1. `grove run ID... --dry-run` prints, writing nothing, the IDs as given,
   the order `deps` gives, each member's status, revision and whether it
   can start or waits and why, the outside prerequisites with their delivery
   at the base, the base, worktree, bounds, review boundary, continuation
   policy and a digest; `--expect DIGEST` refuses a changed assignment,
   printing what it would run now. Refusals before any write: a member not
   work, missing, not proposed or active here or on the branch, uncommitted,
   selected twice, a digest mismatch, an overlapping live attempt, a
   reopened group launched in part, and no member able to start, with waits
   from both the launching checkout and a reused branch. Tests:
   `TestSelectionChainAndBranch`, `TestSelectionRefusals`,
   `TestSelectionWaitsBothPlaces`, `TestSelectionReopenedGroupRunsTogether`,
   `TestRunDryRun`, `TestAttemptCommandsUsage`.
2. A chain, a branch and a path through unselected work: order, waits
   ("needs G-003, which is proposed and not selected"), outside items never
   made members (`TestSelectionThroughUnselectedWork`). Later members start
   only when nothing they need waits; the agent's side is the work guide's
   new selection paragraph.
3. Each member keeps its record, acceptance, evidence and Next. A shared
   candidate: records on a branch with the same `candidate` are one group;
   `approve` per member ignores the group's record commits; `feedback` on
   one reopens all, each committed alone; `integrate` of any merges the
   group only when every member is approved, marks each done alone, and
   refuses to carry unfinished, unapproved work's candidate
   (`TestGroupIntegratesOnlyWhenEveryMemberIsApproved`,
   `TestGroupFeedbackReopensEveryMember`,
   `TestGroupRefusesToCarryAReopenedSibling`). Partial completion and
   waiting members' checkpoints are in the work guide's step 8. A provider
   exit never marks anything done: `attempt` prints each member as the
   branch holds it (`TestMemberResults`, `TestRunSelection`).
4. Inspection starts nothing. Overlap refusal is by membership; `attempts
   ID` lists every selection including ID; Stop and owner loss reconcile
   every member through the same `reconcile` (`TestMemberResults` calls it
   as Stop does; `TestStopAfterReconnect` and `TestOwnerLost` still pass).
   Resume reuses the branch and rechecks the members there; not repeating
   completed units is the guide's rule for the agent.
5. Budget is the provider's one `--max-budget-usd` over the whole
   selection; Grove never retries; per-member waits and states are in
   `attempt` and the board's attempt screen (`TestSelectionAttempt`); an
   unchanged wait is refused at launch.
6. Fake-provider lifecycle: `TestRunSelection` runs a two-member selection
   through the real owner. The real-provider trial is the owner's step
   below. Untested: recovery after a reboot (not selected), a real
   provider's handling of several members, timing of long selections, and
   a Stop or owner loss mid-selection through a live process.

Verification at `15774e4` (the tree the candidate holds, apart from this
record and G-190): `go vet ./...` clean, `gofmt -l .` empty, `go run
./cmd/grove check` OK (181 records), `go test -count=1 -timeout 120s
./...` all pass, `python3 internal/tui/testdata/terminal.py` all pass on a
binary built there. The attempt package runs about 6 s uncached against
5.2 s on main before, 1.2 s under `-short`.

Review: [G-190](G-190-g-162-review-1.md), three fresh `grove-reviewer`
rounds at `99a143c`, `be48722` and `15774e4`; two blocking findings fixed
with regressions and re-reviewed. Open after the last round allowed: two
should-fix refusal messages (`run`'s suggestion for a partial reopened group
lacks `--branch` and other members; `integrate`'s carry refusal gives the
wrong advice for a sibling approved on a stacked branch) and three notes.
Also open for the owner: whether selection and group need term records.

## Next

Checkpoint, 2026-09-25 (headless): in review on `worktree-G-162`, the
candidate the status commit names. To judge it, read the Evidence above and
G-190, then from this checkout
(`.claude/worktrees/worktree-G-162`):

    grove approve G-162 "VERDICT"      # or: grove feedback G-162 "TEXT"

and in main's checkout:

    grove integrate G-162 --cleanup

Acceptance 6's real-provider trial is the owner's step with its own budget,
in a disposable project, with a `grove` built from the integrated main:

    t=$(mktemp -d) && git -C "$t" init -q -b main && grove --project "$t" init
    grove --project "$t" new work "Add greeting.txt saying hello"
    grove --project "$t" new work "Add a second line to greeting.txt"
    grove --project "$t" new work "Describe greeting.txt in README.md"
    # set depends_on so each needs the one before, write each Acceptance,
    # then commit: git -C "$t" add -A && git -C "$t" commit -qm init
    grove --project "$t" run ID1 ID2 ID3 --budget 5 --permission-mode auto --dry-run
    grove --project "$t" run ID1 ID2 ID3 --budget 5 --permission-mode auto --expect DIGEST
    grove --project "$t" attempts    # later, after the terminal is gone:
    grove --project "$t" attempt ATTEMPT
