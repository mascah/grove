---
id: "G-190"
type: review
title: "Final review of G-162 in three rounds: selection, shared candidate, group judgment"
status: current
created: "2026-09-25T23:44:31Z"
updated: "2026-09-26T00:03:01Z"
work: ["G-162"]
examined: "15774e4"
---

## Examined

The final combined review gate of [G-162](G-162-bounded-work-selection.md)
(plan [G-185](G-185-g-162-selected-work-plan.md) step 7), on
`worktree-G-162` from main `fe97300` (merged with main `38f82aa` in
`a11dd94`), by three fresh `grove-reviewer` agents dispatched in turn from
the headless `/grove-work G-162` session of 2026-09-25. Each was read-only
on the checkout and probed in disposable copies under the temp directory:

- Round 1 examined `06015a9..99a143c`, the whole implementation.
- Round 2 examined the fixes `99a143c..be48722` and the branch's readiness.
- Round 3 examined the fixes `be48722..15774e4`, the merge `a11dd94`, and
  `main...15774e4`. `examined` is `15774e4`; the commit after it adds only
  this record and G-162's evidence.

Each ran `go build`, `go vet`, `gofmt -l`, `grove check` and the package
tests; round 3 ran `go test -count=1 -timeout 120s ./...` (all pass) and
`terminal.py` (all pass). No provider ran.

## Findings

Round 1 (at `99a143c`):

1. **Blocking, fixed in `be48722`.** A reused branch's waits replaced the
   launching checkout's, so an open question on main was ignored where a
   branch existed. Waits from both places now combine;
   `TestSelectionWaitsBothPlaces`.
2. **Should-fix, fixed.** The guide let a waiting member's checkpoint land
   after the shared candidate, which `approve` then refuses. The guide now
   commits it before the candidate.
3. **Should-fix, changed.** One ID is a selection of one, so `run` refused
   work with an undelivered prerequisite even under `--until plan`. Under
   the bound only questions stop a member now; without it the refusal
   stands, documented (plan adjustments).
4. **Should-fix, fixed.** The launch printed no selection; it prints the
   preview's lines now. Untested.
5. **Should-fix, fixed.** The board's running tag and latest attempt keyed
   on the first ID; both use membership now, tested.
6. **Note, kept.** A digest mismatch prints the current assignment, not a
   difference (plan adjustments).
7. **Note, kept.** An unreadable outside candidate counts as undelivered.
8. **Note, kept.** The launch decides delivery against the base itself;
   deps' preview is reused for order, outside items and questions.
9. **Note, fixed.** A failure marking one group member done now names the
   rest. Untested.
10. **Note, fixed.** The board's feedback prompt named members not in
    review; it names only those in review, tested.
11. **Note, fixed.** One unreadable attempt blocked every launch; it now
    fails only listings that could include it, tested.
12. **Note, kept.** The attempt package runs about 6 s uncached (5.2 s on
    main before; 1.2 s under `-short`).

Round 2 (at `be48722`):

1. **Blocking, fixed in `15774e4`.** After group feedback, continuing one
   member alone let `integrate` merge its sibling's code whose approval
   feedback had withdrawn. `integrate` now refuses a merge carrying the
   candidate of unfinished, unapproved work the target lacks
   (`TestGroupRefusesToCarryAReopenedSibling`), and `run` refuses a
   selection leaving out a member still sharing the candidate on that
   branch (`TestSelectionReopenedGroupRunsTogether`).
2. **Should-fix, fixed.** The `--until plan` exemption was missing from the
   guide the agent follows; step 2 now says an external blocker stops
   implementation, not preparation.
3. **Should-fix, fixed.** The branch conflicted with main after G-177;
   merged in `a11dd94`, keeping main's merge prediction and the group's
   record paths.
4. **Should-fix, fixed here.** This record was empty.
5. **Note, fixed.** Notes and wait clauses repeated on a reused branch.
6. **Note, kept.** The board's attempt list and `e` know only the handed-off
   or first member's questions; the attempt screen lists each member's.
7. **Note, mitigated.** The default branch follows the IDs' order; feedback
   output and docs give `--branch`.
8. **Note.** Fixes 4 and 9 of round 1 remain untested.
9. **Note, open for the owner.** No term record defines selection or group;
   decision G-188 and the record model do, and "group" also names the
   board's attempt groups and version groups.
10. **Note, fixed at the handoff.** G-162's Next was stale.

Round 3 (at `15774e4`), nothing blocking:

1. **Should-fix, open.** `run`'s refusal of a partial reopened group
   suggests `grove run SELECTED SIBLING` without `--branch` and names only
   the first left-out member, so following it starts a subset on a new
   branch and strands the others. `integrate`'s carry check still refuses
   the unsafe merge. Fix: name every left-out member and add
   `--branch BRANCH`.
2. **Should-fix, open.** Where a sibling was approved on another branch
   stacked below, `integrate`'s carry refusal says to hand it off and judge
   it, though it is approved; integrating it first works. Fix: say "or
   integrate it first".
3. **Note, open.** A sibling candidate Git cannot read makes `integrate`
   fail with Git's own error, naming no record; only hand edits make one.
4. **Note, open.** The board's `i` prompt leaves out a same-candidate
   record not in review, which `integrate` then refuses, explaining why.
5. **Note, open.** The CLI's `grove run IDS --branch BRANCH` feedback hint
   is untested.

## Disposition

Rounds 1 and 2 were fixed with regressions where the finding was
consequential, then re-reviewed. Round 3 was the last round the work guide
allows at this gate. Its two should-fix findings are refusal messages whose
refusals hold, so the candidate is handed to the owner with them open
rather than changed without review. The owner's `feedback` can ask for them.
