---
id: "G-099"
type: review
title: "G-044 review-integration review"
status: current
created: "2026-09-23T01:12:55Z"
updated: "2026-09-23T01:16:37Z"
work: ["G-044"]
examined: "956da77"
---

## Examined

Branch `worktree-G-044`, base `main` c4b3aea. An independent reviewer (a
subagent given [G-044](G-044-review-integration.md), the plan
[G-098](G-098-g-044-review-integration-plan.md), the record model's lifecycle
section, the guide's judging section and the diff, with no write access)
examined `git diff c4b3aea..2e0aca3` in round 1 and the fix commit 956da77
in round 2, which is the `examined` field. It ran `go vet`, `gofmt -l` and
`go test -short` per package, and non-short `go test` on tui, cli and
integrate. The candidate differs from 956da77 by a0b2f2b (a test for
changes and diffs under a prefix, and the result screen's hint while the
re-read is pending: the two round-2 notes) and by this record and G-044's
evidence.

## Findings

1. **Blocking: `integrate` checked one commit and merged another.** The
   approval and tip checks ran on the commit `versions.Inspect` read, but
   the merge was `git merge BRANCH`, which takes whatever the name resolves
   to at that moment: a commit pushed to the branch between the check and
   the merge would be merged unapproved and marked done, and a tag of the
   branch's name wins the name.
2. **Should-fix: the remedy after a failed done write did nothing.** The
   record was already written and staged, so the printed `update … --commit`
   saw no change and committed nothing.
3. **Should-fix: `--cleanup` deleted ignored files** such as `.env`, which
   `git worktree remove` removes without `--force`.
4. **Should-fix: a project under a subdirectory broke approve, integrate
   and the diff view**: three copies of the "other files changed" loop
   compared top-relative `git diff` paths with project-relative record
   paths, and the diff read used a cwd-relative pathspec.
5. **Should-fix: diff lines showed tabs as `\t`.**
6. **Should-fix: the result screen claimed a re-read that had not happened**,
   Esc could return to a stale detail, and the pseudo-terminal scenario
   raced on it (one failure seen in eight runs under load).
7. **Nit: Esc from a detail left a changes or diff read running.**
8. **Nit: Ctrl-C during an action** waits for it and then reports only the
   interruption.
9. **Nit: `update --set approved=` records an approval by hand** without a
   verdict or the tip check.
10. **Should-fix (docs): the board was still described as read-only** in
    the README, AGENTS.md, the record model and two package comments.
11. **Nit: the approval fact did not name the approving commit.**

No finding: every Git process goes through `repo`; one read at a time and
nothing read during the board load; the review block, changes, diff, prompt
and result rows are escaped before any style and the hostile OSC sequence
stays inert; the appended paragraph preserves LF, CRLF and no-final-newline
files under the write lock and revision check; every refusal comes before
the merge and a conflict is aborted with the target unchanged; `done` is
refused off the target; the docs name the commands, keys and facts as
implemented.

## Disposition

1. Fixed in 956da77: the merge is of the inspected commit with `-m "Merge
   branch 'NAME'"`; a test tags `main` with the branch's name and still gets
   the branch's tip fast-forwarded and done.
2. Fixed: the message says to commit the staged record with `git commit --
   PATH`, or to rerun `update` if it was not written, which `update`'s own
   error states.
3. Fixed: `!!` entries of `git status --porcelain --ignored` keep the
   worktree with the files named; tested.
4. Fixed: `versions.Others` places the record under the prefix from
   `repo.IdentifyContext` and serves `update`, `integrate` and `Changes`;
   `DiffContext` uses `:(top,literal)`; the integrate fixture takes a prefix
   and a0b2f2b covers `Changes`, `Others` and `DiffContext` under one.
5. Fixed: tabs become four spaces before escaping; the fixture has a tab.
6. Fixed: the result screen says "Re-reading the board…" until the re-read
   lands, Esc waits for it (the hint says so since a0b2f2b), and the
   pseudo-terminal scenario waits for "The board has been re-read." before
   Esc.
7. Fixed: Esc from a detail stops any pending read but the inspect and an
   action; tested.
8. Recorded as a limit in G-098; the action's commits are in Git.
9. Recorded as an adjustment in G-098: `check` holds `approved` to the
   candidate and `integrate` re-checks the tip.
10. Fixed in the five places.
11. Fixed: "approved on branch X at TIP (Verdict …)".

Round 2 verdict: **clear with notes**, the notes being the two items
a0b2f2b closes.
