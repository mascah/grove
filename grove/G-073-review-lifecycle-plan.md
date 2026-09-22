---
id: "G-073"
type: plan
title: "G-038 review lifecycle plan"
status: current
created: "2026-09-22T04:12:22Z"
updated: "2026-09-22T04:13:00Z"
work: ["G-038"]
---

## Design

Prepared 2026-09-22 on `worktree-G-038` from main `a28a24b`, against
[G-038](G-038-review-lifecycle.md) as selected in
[G-035](G-035-interactive-adoption.md) and the settled terms
[candidate](G-057-candidate.md), [review](G-058-review.md),
[approval](G-059-approval.md) and [integration](G-060-integration.md).
These are routine choices inside the selected contract, not product decisions.

Observed at `a28a24b`: work statuses are one row of `project.Types`
(`internal/project/metadata.go:173`); `update` accepts any status in that row
and has per-type field tables (`internal/update/update.go:134-151,268`); the
board is a fixed `[4]string` of statuses with widths derived from its length
(`internal/tui/model.go:37,583`, `view.go:17,258-289`); a review already
records `examined`, the commit it looked at. Seventeen work records are
`done` under the branch-local meaning; all sit on main.

1. **Statuses.** Work becomes `proposed`, `active`, `review`, `done`,
   `abandoned`, in that order everywhere the row is shown. `new` still writes
   `proposed`.
2. **Candidate on work.** A new optional work field `candidate`: a quoted Git
   commit, the same pattern as `examined`. It is required while `review` and
   allowed on every other status (adjusted during implementation from
   "refused while `proposed`": a proposed record with a candidate is
   harmless, and the refusal would only make reopening to proposed need
   `--unset`).
   It names the last implementation commit; the commit that sets `review`
   changes only the record, so `git diff --stat CANDIDATE TIP` shows one file.
   A changed candidate is a new value set through `update`; the prior value
   stays in Git history and each review's `examined` still says what it saw.
   Input revisions (plan, record, base) stay in the Next checkpoint as today.
3. **Approval and integration.** Approval is the owner's verdict quoted in the
   record, naming the candidate. Integration is the candidate being reachable
   from the target, shown by ancestry. `update --set status=done` therefore
   requires a candidate, present or set in the same call, and refuses one
   that `git merge-base --is-ancestor` cannot reach from the checkout's HEAD,
   so a checkout without the code cannot close the work. The CLI cannot tell
   the target from the work branch, where the candidate is an ancestor too
   (the independent review's first finding): writing done on the target
   after the merge is the guide's rule. A squash or rebase that lands a
   different commit names that
   commit as the candidate in the same call. No `target:` configuration, no
   approval field and no `report` type: the interactive loop needs none, the
   work record's Evidence is the report, and G-044 adds structured approval
   when software acts on it.
4. **Historical Done.** A `done` work record without `candidate` keeps its
   original branch-local meaning. Nothing is rewritten and `check` accepts it;
   `update` never writes a new one. That absence is the explicit migration
   marker. G-038's evidence lists this repository's historical Done records
   and their observed presence on main.
5. **Research and design completion.** The same rule: the deliverable is
   files, so Done means the candidate holding them was accepted and reached
   the target.
6. **Guidance, not enforcement, for the rest.** The order of transitions,
   Abandoned needing a human decision, and a failed or interrupted attempt
   staying `active` with a checkpoint are guide rules. Software cannot verify
   a human, and the owner edits by hand.
7. **Handoff.** Entering Review means the record's Evidence and Next hold,
   without the chat: branch, base, candidate, input revisions, behavior per
   acceptance item, decisions taken, verification commands with results at
   the candidate, review records with findings and dispositions, unresolved
   issues, and the integrator's next action as runnable commands. Feedback
   sets `active` with the feedback in Next and keeps everything above.

## Steps

1. `internal/project`: add `review` and `candidate` with the rules in design 2;
   tests for each refusal by file and field.
2. `internal/update`: `candidate` in the field tables, the string maps and the
   drift guard; the Done rule of design 3 through `repo.Git`; tests for
   done-without-candidate, unreachable candidate, same-call candidate, and the
   review round trip.
3. `internal/tui`: five columns and short tab names; existing layout tests
   extended; a connected terminal check at 100 and 80 columns.
4. `internal/handoff`: reword the two "done is not integration" lines for the
   candidate rule.
5. Documentation: record model (lifecycle, work fields, review section, the
   "do not write review" note), README (status paragraph, board columns),
   `docs/work-execution.md` (drop the Contract transition section; steps 5,
   7 and 8 write the handoff of design 7 and set `review`; the integrator's
   path), `AGENTS.md` summary line, and the term bodies of G-054, G-058,
   G-059 and G-060.
6. Exercise in a disposable clone: successful, waiting, failed, feedback and
   changed-candidate paths through the CLI; keep the transcript in G-038's
   evidence.
7. Independent review at the final revision, recorded as a review record with
   `examined`; fix rounds capped at three.
8. Hand G-038 itself off with the new mechanism: candidate set, status
   `review`, Next naming the integrator's commands.
