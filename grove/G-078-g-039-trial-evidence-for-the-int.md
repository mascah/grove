---
id: "G-078"
type: review
title: "G-039 trial evidence for the interactive loop"
status: current
created: "2026-09-22T15:33:55Z"
updated: "2026-09-22T15:34:04Z"
work: ["G-039"]
examined: "840a78b"
---

## Examined

Trial evidence for [G-039](G-039-interactive-loop.md) under its plan
[G-075](G-075-trial-plan-for-the-interactive-l.md), written 2026-09-22 by
the fresh `/grove-work G-039` session that G-075 step 4 asked for. `examined`
is `840a78b`, the candidate of the real change
[G-076](G-076-filter-grove-list-by-status.md); the loop that produced it is
main `ff0f3e9..e8bcf03`. Harness: Claude Code 2.1.278, model Fable 5.1 for
every interactive session; the one headless call below ran on Claude Opus 5
(the `claude -p` default), go 1.26.2, git 2.55.0, macOS Darwin 25.6.0.

Every observation is labelled **real** (a fresh session on this repository),
**simulation** (a disposable clone reached by `--project`), or
**unexercised**. Nothing was rerun and no record was invented for the trial:
G-076 was a change the owner wanted.

## Findings

**Real trial, G-075 steps 1 to 3, all on main.** Elapsed from the proposal
commit to done: 22 minutes (09:03 to 09:25 local, commit timestamps).

1. Shape (`/grove-shape add a status filter to grove list`, interactive, main
   checkout): one proposed record, G-076, committed at `ff0f3e9`, with
   outcome, observed constraints at `5f07c94` with file and line, three
   offered scopes, the owner's choice recorded as theirs ("status only, and
   `list` without the flag keeps printing everything"), a design labelled
   proposed, and six acceptance items, the last the owner's judgment. Nothing
   assigned or implemented. The session found the brief, the `list` code and
   the record model's status vocabulary without being pointed at them (the
   record cites `internal/cli/cli.go:156` and the model's "Statuses"). One
   scope question was escalated; it was a genuine product choice.
2. Implement (`/grove-work G-076`, fresh session): branch `worktree-G-076`
   from main `ff0f3e9`; "no plan needed" with the reason in Evidence;
   `active` through the CLI; no escalation; feature `32b7b40`, independent
   review [G-077](G-077-g-076-independent-review-of-the.md) examined
   `32b7b40` with three findings fixed in `cd1067d`; evidence `840a78b`; the
   handoff commit `5d9a18a` changed the record alone. Acceptance items 1 to 5
   each have named evidence in the record.
3. Judge and integrate (owner): `git merge --ff-only worktree-G-076` and
   `update … status=done` at `e8bcf03`, four minutes after the handoff. The
   record and CLI carried enough to integrate, but the done commit changed
   `status` and `updated` only: the verdict is not quoted in the record, and
   its Next still says "Candidate awaits the owner's judgment". The owner,
   asked in the G-039 review session on 2026-09-22, said: "It went well. I
   did not need the chat to integrate it."
4. Continue (this session): `/grove-work G-039` in a fresh session found the
   checkpoint in G-039's Next and plan G-075 through `context`, rebased
   `worktree-G-039` from `5f07c94` onto `e8bcf03` (one commit, no conflict)
   so that G-076 and G-077 resolve as links, and carried on without the chat
   that wrote them.

**Simulations, in two clones of main `e8bcf03` under the session scratchpad,
with the CLI built once from that commit.**

5. Feedback round (clone `trial-clone-fb`, branch `sim-feedback` from
   `5d9a18a`): feedback appended to G-076's Next, `update --set
   status=active` (exit 0), a fix commit `98a3515`, `update --set
   status=review --set candidate=98a3515` (exit 0), and `git diff --stat
   98a3515 HEAD` touched the record alone. `context G-076` there listed
   G-077, whose `examined` is the first candidate `32b7b40`, so a second
   round shows the reviewer the gap to explain. In a branch at `ff0f3e9`,
   `update --set status=done --set candidate=840a78b` was refused with exit 1:
   "done means accepted and integrated: candidate 840a78b is not an ancestor
   of this checkout's HEAD; mark done where it was merged". On the branch
   holding the candidate the same write succeeded, as the model says it must.
6. Missing human decision, headless (clone `trial-clone`): `claude -p
   "/grove-shape hide finished work from grove list by default --interaction
   headless"` with write tools allowed, 92 seconds, exit 0. It created
   `worktree-shape-list-hide-finished` from `e8bcf03` in the clone's
   `.claude/worktrees/`, wrote one proposed work record, ran `check` (77
   records), committed `ab77fcd`, and returned the branch, commit, and that
   nothing was assigned or merged. It did **not** take the guide's
   missing-choice path: it scoped "finished" to work records itself, chose
   `--all` and a stderr note in a design labelled proposed, wrote acceptance
   that presumes them, and listed the three choices in Next as "left for
   you", reasoning that "none blocks the work, so each is in the record's
   Next rather than a question record". The record's ID, G-078, collides with
   this one: a clone has its own counter, as the shaping guide warns.
7. Question mechanics (clone `trial-clone-fb`, branch `sim-question`): `new
   question`, `update --set 'blocks=["G-036"]'`, commit `4b264b0`; `context
   G-036 --interaction headless` listed it as "question blocking G-036" and
   "open, blocks G-036"; `check` OK. The mechanics the headless path relies
   on work; only the harness's choice to use them was missing.

**Unexercised.** The headless work row (`claude -p "/grove-work …
--interaction headless"`), a rerun of the same headless shaping to confirm an
unchanged wait, and the review cap (three fix/review rounds) did not arise and
were not simulated.

**Knowledge (G-039 acceptance 6).** No new schema type or category was
needed: G-076's constraints hold the observed facts about `list`, the trial
plan and this evidence are plan and review records, and the terms G-057 to
G-060 already name candidate, review, approval and integration. `context`
staged the retrieval: the selected record and plan in full, G-076 and G-077
listed until this session read them. IDs and paths were stable throughout:
G-076 and G-077 kept their files from proposal to done, and G-039 and G-075
kept theirs across the rebase.

## Disposition

Against G-039's acceptance:

1. Met, real: item 1 above, with the harness discovering guides and code.
2. Met, real: item 2; proportionate preparation, one genuine escalation in
   shaping and none in implementation, checkpoints in Next, linked review.
3. Met with a gap, real: item 3; the owner integrated from files and the CLI
   without the chat, by their own account, but the verdict was not quoted
   and Next was not reconciled at done.
4. Met, real for continuation (item 4), simulation for feedback (5) and the
   missing-human case (6, 7), with the headless divergence recorded.
5. Met, this record and G-039's Evidence.
6. Met, real, as observed above.

Friction and the next improvement, each fed to its owner in the candidate
commit: the integration commands a handoff offers should include writing the
verdict, since the integrator ran them as given
(`docs/work-execution.md`, step 8); the headless shaping bound should say that
a scope or design choice the acceptance depends on is a missing choice even
when the proposal could be assigned without it (`docs/work-shaping.md`); and
`new review` leaves `Examined`, `Findings` and `Disposition` headings that
G-077 left empty while writing its own, which the CLI could fill or the guide
could name (G-036's Next, not changed here). On the headless divergence in
finding 6 the owner said, in the same session, that it concerns them, that
little can be done about it right now, and that an eval suite for the guides
should be considered in future; that is intent, recorded in G-036's Next,
not a proposal.
