---
id: "G-039"
type: work
title: "Prove the complete interactive workflow on real Grove work"
status: active
created: "2026-09-21T00:54:14Z"
updated: "2026-09-22T14:57:26Z"
kind: investigation
size: small
priority: 1
depends_on: ["G-025", "G-037", "G-038"]
relates_to: ["G-035", "G-036", "G-064", "G-065"]
formerly: "W-021"
---

## Outcome

Prove Grove's complete interactive shape → prepare → implement → independent
review → human review → integrate loop on actual Grove development.

## Scope and evidence

Use delivered G-025/G-037/G-065/G-038 behavior on a newly selected, bounded real change.
Choose a small real assignment with the owner after the prerequisites land.
Do not count a prerequisite's earlier implementation as a trial of behavior that
was not available then, rerun completed work, or invent demonstration records.
This is a connected-workflow trial, not another implementation of those features.

## Acceptance

1. Fresh interactive harness invocation discovers the guides and appropriate
   context, shapes useful work and records actual knowledge changes.
2. Implementation prepares proportionally, escalates only genuine human
   decisions, retains checkpoints and finishes with linked review evidence.
3. The owner can review the result from files/CLI without reconstructing chat,
   request a follow-up or accept it, and complete explicit local integration.
4. Observe fresh-session continuation. Exercise missing-human and feedback
   cases in disposable fixtures if they do not arise naturally.
5. Record invocation, harness/version, source and candidate revisions, actual
   outcomes, owner feedback, observed friction and the next improvement. Separate
   real trials, simulations and unsupported/unexercised harness behavior.
6. Observe useful general knowledge captured without a new schema type, staged
   retrieval of it when relevant, and stable IDs/paths through work completion.
   Do not manufacture new categories or knowledge solely to satisfy the trial.

## Preparation

The prerequisites are delivered on main at `5f07c94`: G-038's candidate
`fae1e4c` is an ancestor, and G-025, G-037 and G-065 show as the shaping guide
with both skill adapters, the `brief:` key with plan and review records, and
neutral IDs with `page` records. The owner selected the real change on
2026-09-22: a status filter for `grove list`, observed missing at `5f07c94`.
The trial plan is [G-075](G-075-trial-plan-for-the-interactive-l.md).
The user must supply the actual usability verdict. This record does not itself
authorize merging any trial branch: obtain that work's explicit disposition.
Feed lessons into their owning guide or work record and the adoption milestone.

## Evidence

Trial run 2026-09-22 under plan G-075 on branch `worktree-G-039` (worktree
`.claude/worktrees/G-039`), rebased from main `5f07c94` onto `e8bcf03` once
the real change was integrated there; started from this record at revision
`43830baa…` and G-075 at `8dedc7c6…`, neither changed by the rebase. The
evidence per acceptance item, with each observation labelled real,
simulation or unexercised, is the review record
[G-078](G-078-g-039-trial-evidence-for-the-int.md), examined `840a78b`, the
candidate of the real change [G-076](G-076-filter-grove-list-by-status.md)
whose loop is main `ff0f3e9..e8bcf03`. In short: items 1, 2, 4 (continuation)
and 6 were met on real sessions; item 3 was met with one gap, the owner's
verdict not quoted at done; item 4's feedback and missing-human cases were
simulated in disposable clones, where the headless shaping call diverged from
the guide by settling product choices itself; item 5 is G-078 and this
section. The candidate is the commit that adds this section.

Decisions: the rebase, so that G-076 and G-077 resolve as links from G-078;
the headless simulation ran `claude -p` with write tools allowed on the
clone, which used Claude Opus 5 rather than the interactive model. Lessons
went to their owners in the same commit: `docs/work-execution.md` step 8 (the
handoff's integration commands include the verdict edit),
`docs/work-shaping.md` headless bounds (a choice the acceptance depends on is
a missing choice), and G-036's Next (open improvements).

Verification at the candidate: `go run ./cmd/grove check` → `OK: 78
records`; every link written here and in G-078 resolves to a file in this
checkout; no Go source changed, so `go vet`, `gofmt` and the suite are
unaffected. Review: G-078 is this session's evidence, not an independent
review; the record and plan ask for none, so that omission is reported here
as open rather than blocking.

Limits: the owner's verbatim verdict on G-076 and whether chat was needed for
it are not on file; the headless work row and the review cap were not
exercised; the clones were left under the session scratchpad and are not
part of this branch.

## Next

Candidate awaits the owner's judgment. Two facts only the owner has, to write
into G-078's finding 3 when judging: the G-076 verdict verbatim, and whether
chat was needed to integrate it. Demo, from any checkout:

```sh
cd .claude/worktrees/G-039
go run ./cmd/grove context G-039 --include grove/G-078-g-039-trial-evidence-for-the-int.md
go run ./cmd/grove show G-078
git diff --stat main
```

Integrate on approval, from the main checkout, quoting the verdict here
before setting done:

```sh
git merge --ff-only worktree-G-039
go run ./cmd/grove show G-039 --json    # take the revision
go run ./cmd/grove update G-039 --expect REVISION --set status=done
git commit -am "docs(G-039): mark done on the owner's acceptance and merge"
```

Feedback instead: write it here and set `status=active` on the branch.
