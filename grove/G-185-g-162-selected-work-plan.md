---
id: "G-185"
type: plan
title: "G-162 selected work: one attempt, shared candidate, group judgment"
status: current
created: "2026-09-25T23:08:19Z"
updated: "2026-09-25T23:08:22Z"
work: ["G-162"]
---

## Inputs

Prepared on 2026-09-25, headless, bounded at the plan, on branch
`worktree-G-162` from `main` `fe97300`, for
[G-162](G-162-bounded-work-selection.md) at `sha256:bf15491b…`. Read in
full: G-162; [G-163](G-163-selected-work-review-boundary.md) (resolved,
answer "option 1": explicit shared implementation, reviewed together); the
settled terms [Work](G-054-work.md), [Attempt](G-056-attempt.md),
[Candidate](G-057-candidate.md), [Approval](G-059-approval.md) and
[Integration](G-060-integration.md); decision
[G-101](G-101-attempt-mechanism.md); `internal/attempt/attempt.go` (Start,
Own, reconcile, List, Stop), `internal/update/review.go`,
`internal/integrate/integrate.go`, `internal/deps/deps.go` (G-161's
Order, Preview, Deliver), `versions.Others`, the Attempts section of
`docs/commands.md`, and the work guide.

## Observed at `fe97300`

- The work guide already runs several IDs in one assignment, in `context`
  order, one implementer at a time, in one execution checkout, and
  `/grove-work` accepts several IDs. The gap is Grove's side: `run` takes one
  ID, an attempt records one `work`, one record state and one revision, and
  duplicate-start refusal matches attempt directories by one ID prefix.
- `deps.Preview` (G-161) already gives a selection's order, the outside
  prerequisites (never added), blocking questions and delivery relative to a
  checkout's HEAD. It is the one interpretation the launch must reuse.
- `approve` and `integrate` refuse any file other than the record's changed
  after the candidate (`versions.Others`). With several records handed off
  on one branch, each other record's own status commit trips that check, and
  integrating one record merges the whole branch, so its siblings' code with
  it, unapproved.
- `feedback` reopens one record and keeps its candidate.
- `grove.yaml` here defaults `run.budget: 50`; the provider enforces
  `--max-budget-usd` over the whole process, subagents included. Grove
  never retries a process.

## Design

The owner chose option 1 (G-163). Everything below follows from it, and none
of it adds a record type, a field, or a status.

### The assignment: one attempt over an explicit selection

- `grove run ID... [flags]` launches one Grove-owned process,
  `claude -p "/grove-work ID... [--until plan] --interaction headless"`,
  under the existing owner: one worktree, one budget, one Stop, one
  owner-loss path (G-101 unchanged). Several IDs never mean several
  processes: implementation is sequential inside the one session, and the
  absence of an edge is not a concurrency mandate.
- The IDs are passed as given; `deps.Order` orders them, and the attempt
  records both. The branch defaults to `worktree-` plus the IDs joined with
  `-` as given (`worktree-G-030-G-031`), as the repository's naming already
  says; `--branch` and `--worktree` still override. The attempt id keeps
  the first ID as its prefix, and `attempts ID` lists every attempt whose
  selection contains ID.
- **Bounds.** The one `--max-budget-usd` is the aggregate; there are no
  subordinate attempt budgets (the provider's subagents share it). Grove
  starts one process and never retries; the guide's three fix/review rounds
  per gate are the repair bound; an unchanged wait is refused at launch, so
  it cannot consume a repeated attempt.
- **Preview.** `run ID... --dry-run` loads, checks and prints the assignment
  without writing or starting anything: IDs as given, order, each member's
  status and revision, the launching checkout and base commit (or the
  reused branch's tip), branch and worktree, budget, permission mode, model,
  effort and bound, the outside prerequisites with their delivery relative to
  that base, open questions, which members can start and which wait and why,
  the review boundary and the continuation policy (below), and an assignment
  digest: sha256 over the ordered IDs with their revisions, the base, and the
  resolved launch options. `run … --expect DIGEST` refuses a launch whose
  digest differs, naming what changed; `attempt.json` records the digest
  and every member revision.
- **Launch refusals**, all before any write, as today: a member that is not
  work, not in this checkout, or not proposed/active (here or on a reused
  branch, where review means "judge it first"); a member record with
  uncommitted changes; a selection `deps.Order` refuses; a digest mismatch;
  a running or orphaned attempt whose selection shares any member
  (overlapping selections cannot both own a record); and a selection where
  **no member can start**. A member waits when an open question blocks it,
  when an outside prerequisite is definitely undelivered at the base
  (proposed, active, review or abandoned, or done with a candidate the base
  lacks; done without one is reported as unrecorded delivery, not refused),
  or when a selected prerequisite waits. Waiting members do not refuse the
  launch while another member can start; they are shown in the preview and
  in the launch output.

### Execution inside the attempt (work guide)

- Members run in the context order. A member whose inputs wait is not
  started. A new question, an external blocker, or a failure of one member
  stops that member and every selected member that needs it, directly or
  through other work; members that need none of them continue. Budget
  exhaustion or Stop ends everything, with what is committed kept.
- Each member keeps its own acceptance, evidence and Next; the guide's
  per-member plan, question and checkpoint rules are unchanged.
- **Review boundary.** The independent review runs on each plan's named
  boundaries and once on the combined result, with each member's evidence
  kept separate.
- **Handoff.** When every member the attempt started is complete, the
  complete members enter review together with **one shared candidate**: the
  combined commit. One commit then sets `status=review` and `candidate=C` on
  all of them, and touches only their record files. A member not started
  (waiting) stays as it was, with a checkpoint naming its wait.
- **Partial completion.** If a started member is incomplete (failed, out of
  budget, stopped, or at the review cap) with changes outside the record
  root on the branch, nothing on that branch enters review: the complete
  members stay active with a checkpoint "complete at commit X, review waits
  on G-…", since integrating the branch would carry the incomplete code.
  Records, plans and questions of a waiting member are record-root files
  and do not hold the branch back. Relaunching the same selection resumes on
  the same branch and does not redo members whose checkpoint the branch
  confirms.

### Judging and integrating a shared candidate

- **Group.** The records on a branch whose `candidate` is the same commit
  are one group. Nothing else defines it.
- **Approve** stays per record (per-item acceptance): `approve` of a member
  binds that member's verdict to C. Its "nothing but the record changed"
  check excludes every group member's record file, not only its own
  (`versions.Others` takes several record paths).
- **Integrate** of any member integrates the group, because merging C merges
  all of it: it refuses, naming them, unless every group member on the
  branch is approved in review; then one merge, and `done` written and
  committed for each member alone, each reported as its own fact. The board's
  `i` gets this for free.
- **Feedback** on any member reopens the group: that member gets the owner's
  text, and every other member is set active with its approval dropped and a
  line "Reopened with G-…'s feedback on candidate C" appended, each committed
  alone. A changed candidate would void their approvals anyway (Approval);
  this makes their stale evidence visible instead of leaving them in review
  on a candidate that can no longer be integrated alone. The board's
  feedback prompt names the group.
- A single-member group behaves exactly as today.

### Progress after the terminal is gone

- `attempt.json` gains the selection (IDs as given, order, per-member record
  path and revision, digest); `result.json` gains the state of every member
  as the worktree holds it (status, candidate, revision, uncommitted) and
  the open questions blocking each. Attempt files from before read as a
  selection of one.
- `grove attempt ATTEMPT` prints one line per member: awaiting judgment
  (review, candidate), active (checkpoint in its Next), not started,
  waiting on question G-… or prerequisite G-…, plus the existing per-launch
  inputs-changed check for each member. `attempts` lists the selection. The
  board's attempt list and detail show the selection and the per-member
  lines; the board's `R` stays one work's launch (CLI selection exercises
  this outcome; G-161 owns the board preview).
- Nothing in inspection starts a process; reopening it launches nothing.

## Steps

1. Record the G-163 answer as an accepted decision attributed to the owner
   (`grove new decision`), linked from G-163 and named in G-162's
   `relates_to`, with the consequences above that the owner approves by
   launching this plan's continuation.
2. `internal/attempt`: several IDs in `Request`, `Launch` and `Result`;
   `deps` for order, waits and delivery; the digest; `--dry-run` and
   `--expect`; overlap refusal by membership; per-member reconcile. Fake
   provider tests: chain A→B→C, branch A→{B,C}, a path through unselected
   work (ordered, never implemented or listed as selected), overlap refusal,
   all-waiting refusal, digest mismatch, reuse of a branch with a member in
   review, Stop and orphan with a multi-member result, an old single-work
   attempt file.
3. `internal/cli`: `run ID...`, usage, `attempts`/`attempt` output.
4. `internal/update` and `internal/integrate`: group-aware approve,
   feedback and integrate, with `versions.Others` taking several paths.
   Tests in disposable repositories: shared candidate approved one by one;
   integrate refused with one member unapproved and target unchanged;
   integrate of the group writes done for each; feedback on the first member
   reopens the rest; a one-member group unchanged.
5. `internal/tui`: attempt list and detail for a selection, feedback prompt
   naming the group; `terminal.py` run.
6. Documents: work guide (several members: continuation, shared-candidate
   handoff, partial completion, group judgment), record model (candidate
   shared by a group, approval per member, integration and feedback per
   group), `docs/commands.md` (run, attempts, approve, feedback, integrate),
   `docs/board.md` where the attempt view changes. Shipped-document link
   rules apply (G-146, G-151).
7. Verification per CLAUDE.md; independent review of the combined candidate
   through `grove-reviewer`; handoff into review.
8. **Owner step, not in the attempt:** acceptance 6's real-provider trial
   runs in a disposable project with its own explicit budget, e.g. three
   small chained work items launched with `grove run A B C --budget 5`,
   inspected afterwards with `attempt`. The implementation hands this off as
   a runnable command and reports recovery and timing limits as untested;
   it does not spend on a provider itself.

## Limits

No parallel implementation, scheduling, automatic merge, reboot recovery or
new provider. Budget is the provider's enforcement, so an overrun inside one
turn is the provider's. Waits are detected at launch and by the agent;
Grove does not interrupt a running attempt when a record changes on the
target, and reports it afterwards per member.
