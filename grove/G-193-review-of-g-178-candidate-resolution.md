---
id: "G-193"
type: review
title: "Review of G-178 candidate resolution"
status: current
created: "2026-09-26T02:40:12Z"
updated: "2026-09-26T02:40:33Z"
work: ["G-178"]
examined: "7cb7ceae8ee0f55b253a8a124670012777680ed0"
---

## Examined

[G-178](G-178-candidate-target-update.md) on `worktree-G-178`: the change
from main `e812672` to `7cb7cea`, against the record's acceptance, plan
[G-191](G-191-plan-for-g-178-candidate-resolution.md) and the repository's
instructions. Three rounds, each by a fresh `grove-reviewer` agent,
read-only. Each built and ran what it needed on disposable repositories.

- Round 1 examined `14b31f6`.
- Round 2 examined `b52448f`.
- Round 3 examined `7cb7cea`, the code of the candidate. The candidate adds
  only this record and G-178's evidence.

## Findings

Round 1 (`14b31f6`):

1. **Consequential.** The resolved files came from `diff-tree --cc`. That
   hid a conflict settled by taking one side, which drops the other side's
   change, and it listed files that Git had merged by itself. **Fixed in
   `b52448f`:** the files are now the conflicts that merging the merge's
   parents again reports, each saying which side it kept.
2. **Minor.** A merge of unrelated history into the branch failed the whole
   changes read. **Fixed** with `merge-base --is-ancestor`.
3. **Minor.** A wait, a missing skill, an incompatible entrypoint or a
   failing `--version` was refused only after the feedback was written.
   **Fixed:** `Resolve` asks all of them before it writes, and
   `entrypoints` and `provider` are factored out of `Start` unchanged.
4. **Minor, knowledge.** "resolve" has two senses, and no term defined this
   one. **Fixed:** term [G-192](G-192-resolution.md), proposed.
5. **Minor.** Nothing tested a group. **Fixed:** `TestResolveAGroup`.
6. **Note for G-180.** `Resolve` takes no policy attribution. Left to
   G-180.
7. **Cosmetic.** The report line was empty when Git named no file.
   **Fixed.**

Round 2 (`b52448f`):

1. **Consequential.** Under a project prefix, every resolved file was
   labelled "took the target's side". The cause is that `git merge-tree`
   names files relative to the current directory, which also affected
   G-177's prediction. **Fixed in `7cb7cea`:** `resolveCommits` reads the
   prefix in its `rev-parse` and `predict` joins it, and the side diffs
   use `:(top,literal)` pathspecs.
2. **Minor.** Renames were paired by the user's `diff.renames`.
   **Fixed:** `--no-renames`, and a stricter side rule.
3. **Minor.** A branch whose history began apart from the target failed the
   changes read. **Fixed:** it names no resolution.
4. **Documentation.** The guide sent the reviewer to `--cc`. **Fixed:**
   `--remerge-diff`.
5. **Cosmetic.** A line in `commands.md` was not wrapped. **Fixed.**

Round 3 (`7cb7cea`): nothing consequential. Every round-2 fix was
confirmed by running it, and every G-177 caller of the changed
prediction was checked. Two minor findings and one note are left open:

1. A user's `diff.relative=true` makes the side check find nothing under a
   prefix. Every file is then unlabelled rather than mislabelled. The
   older `numstat` and `Others` diffs share this exposure. The fix is
   `--no-relative`.
2. The test suite does not cover a branch whose history began apart and
   then merged the target. The reviewer verified it with a throwaway test.
3. Note: any prediction failure on the merge's parents, such as Git before
   2.38, drops the resolution row without a reason.

Knowledge: term G-192 covers the new concept. The work contradicts no
settled term or accepted decision: candidate, approval and integration keep
their meanings, and the operation is the resolution attempt that G-182
names.

## Disposition

Every consequential finding was fixed and re-reviewed. The review cap of
three rounds was reached, and the last round found nothing consequential.
The three open minor items are listed in G-178's limits for the owner.
This review is evidence, not approval.
