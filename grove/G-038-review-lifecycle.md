---
id: "G-038"
type: work
title: "Hand implementation candidates into revision-bound human review"
status: active
created: "2026-09-21T00:54:14Z"
updated: "2026-09-22T04:13:00Z"
kind: feature
size: medium
priority: 1
depends_on: ["G-037", "G-065"]
relates_to: ["G-035", "G-036", "G-023", "G-064"]
formerly: "W-020"
---

## Outcome

An implementation returns a revision-bound candidate for human review instead
of leaving the owner to infer its disposition from Active/Done and chat prose.
Make the interactive loop work before adding a managed runner.

## Selected contract

Target Proposed → Active → Review → Done; Abandoned requires an explicit human
decision. For implementation, Done means accepted and integrated into the
configured target. Define equivalent completion for research/design deliverables.
Preparation, independent agent review and waiting are activities or additional
facts. A failed/interrupted terminal attempt does not imply Review readiness.

An explicit assignment authorizes bounded execution; a proposed record or
status edit alone does not. Prepare missing plans within that mandate, preserve
the small-work exception, and ask about consequential choices rather than routine
technical steps. Keep interactive/headless instructions shared.

Build on G-065's identity/storage contract, as selected in
[G-064](G-064-stable-knowledge.md). Review, approval and integration
remain explicit operational facts; generic knowledge does not imply authority.
Lifecycle transitions preserve IDs and paths, including Done. No completion
folder or record-filing operation is part of integration cleanup.

## Acceptance

1. CLI/schema and shared execution guidance support Review and its transitions;
   the existing board can display the added status without requiring a redesign.
2. The handoff links outcome, changed behavior, decisions, acceptance evidence,
   test revision, independent review findings, unresolved issues and next action.
   A new session can inspect it without the original chat.
3. Review identifies candidate/input revisions. Changed content cannot silently
   inherit approval. Feedback that starts implementation returns work to Active
   and retains prior evidence.
4. Historical Done records retain their original meaning and evidence. Design
   and document an explicit migration; do not infer integration from status.
5. Manual local approval/integration can complete the loop with an honest
   recorded disposition. Automated merge/cleanup and the review TUI belong to
   G-044; managed process ownership belongs to G-045.
6. Exercise successful, waiting, failed, feedback and changed-candidate paths;
   reconcile the model, guide and CLI documentation. Owner usability judgment
   belongs to the real loop in G-039.

## Evidence

Branch `worktree-G-038` from main `a28a24b`, in `.claude/worktrees/G-038`;
inputs were this record at `sha256:a39d1331…` and the plan
[G-073](G-073-review-lifecycle-plan.md), prepared in `190669c` with one
adjustment recorded in it. Commits: `6c15442` code, `f9f6c0d` docs,
`12b7752` and `4109c95` fixes from the independent review
[G-074](G-074-review-lifecycle-review.md) (three rounds, every consequential
finding fixed), `87e66ab` the test that round 3 found missing, fixed after
the review cap and therefore self-checked by mutation rather than
independently reviewed; then the evidence commit, which is the candidate
and changes only records. G-074's dispositions list what was left to the
owner.

Verified at `87e66ab`, the last code change, uncached: `gofmt -l .` and
`go vet ./...` clean, `go run ./cmd/grove check` OK, and
`go test -count=1 -timeout 120s ./...` ok in every package (see Next for the
timing note). The board was opened on a pseudo-terminal at 100×24 and 80×24:
it read 3 branches and 3 checkouts, rendered `Proposed (9) Active (1)
Review (0) Done (17) Abandoned (0)` (five tabs at 80), and exited 0 on `q`.

1. Work statuses are `proposed`, `active`, `review`, `done`, `abandoned`
   (`internal/project/metadata.go`); `update` takes `review` and the new
   `candidate` field through the same tables; the board shows five columns
   from the same list with no other change (`internal/tui`). Visible cost:
   five columns at 100 cells leave 19 per card, so titles truncate sooner
   than with four; G-043 owns the visual redesign. The guide has a Lifecycle
   section, the Review handoff in step 8, and the integrator's path.
2. The handoff is the record's Evidence and Next as step 8 lists them; this
   record is the first, and `context G-038` prints it in full. The change is
   in `docs/record-model.md` (Work lifecycle), `docs/work-execution.md`,
   README, AGENTS.md and the terms G-054, G-058, G-059 and G-060.
3. `candidate` is a quoted commit, required in `review` and changeable
   through `update`; a review's `examined` is compared to it by the reader.
   Approval is of the candidate, and `done` is refused where HEAD lacks it,
   so changed content cannot inherit approval through the CLI; a further
   commit after the handoff is a new candidate (guide). Feedback sets
   `active` and removes nothing (`TestUpdateLifecycleAndReopening`).
4. A `done` record without `candidate` keeps its meaning: the reader accepts
   it, its fields stay editable, `context` shows it as `done, no candidate`,
   and `update` never writes a new one
   (`TestUpdateDoneMeansAnIntegratedCandidate`). That rule and the record
   model section are the migration. This repository's seventeen historical
   Done records (G-003, G-007, G-009, G-010, G-011, G-014, G-015, G-016,
   G-017, G-023, G-025, G-030, G-031, G-037, G-052, G-065, G-071) are all
   on main at `a28a24b` and were left untouched, with no backfilled candidate.
5. Manual approval and integration: merge, then `update --set status=done`
   in the target's checkout, where the check holds it to the merged code. The
   CLI cannot tell the target from the work branch, where the candidate is an
   ancestor too, so "done on the target, never on the branch" is the guide's
   and AGENTS.md's rule (G-074 finding 1). No automation, no `target`
   configuration, no approval field, no `report` type.
6. Exercised with binaries built from `f9f6c0d` and again from `4109c95`
   in a disposable clone, with the same outcomes, transcripts kept outside
   the repository: success (new → active → implement
   → `review` refused without a candidate, accepted with one → `git diff
   --stat CANDIDATE HEAD` shows one file → `versions` shows `review` on the
   branch and `proposed` on main → done on main refused before the merge,
   accepted after); feedback (review → active keeps the candidate → new
   commit → new candidate → review, both candidates in Git history);
   changed candidate (a commit after the handoff shows as a two-file diff,
   `--set candidate=` replaces it, `candidate=main` refused by form, an
   unknown commit refused at done); failed attempt (stays `active`);
   waiting (an open question `blocks` it and `context` lists it);
   historical done (valid, renamed, refused an unreachable candidate).
   Owner usability is G-039's.

## Next

In Review: the candidate is the evidence commit named in `candidate`, and
the branch tip adds only this status change. To judge it, from a checkout
of `worktree-G-038`:

```sh
go run ./cmd/grove context G-038          # this record and G-073, G-074 listed
git diff --stat "$(git log -1 --format=%H)"~1  # the handoff commit touches one file
git diff a28a24b..HEAD --stat             # 25 files
go test -count=1 -timeout 120s ./...
```

Approve: `git merge --ff-only worktree-G-038` on `main`, then there
`go run ./cmd/grove update G-038 --expect REVISION --set status=done`,
quoting the verdict here, and commit. Feedback: `--set status=active` on the
branch with the feedback here. Open for the owner: the brief's "Done will
mean" and "G-038 must migrate" sentences are now dated (G-074 finding 12);
AGENTS.md reserves the brief for direction changes, so they were left.
Timing: `internal/versions` took 7.1 s in the full run, above the 5 s
budget, from the Git process ceiling G-071 recorded, not from this change.
G-039 exercises this loop on real work; G-044 adds the approval and
integration actions.
