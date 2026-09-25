---
id: "G-177"
type: work
title: "Predict whether a candidate in review merges cleanly into the target"
status: active
created: "2026-09-25T21:39:27Z"
updated: "2026-09-25T23:04:15Z"
kind: feature
relates_to: ["G-161", "G-162", "G-163", "G-044", "G-060", "G-057", "G-030", "G-178", "G-180"]
---

## Outcome

Before pressing `i`, the owner can see for each candidate in review whether
it merges cleanly into the target at its current tip or conflicts in named
files, and, for several candidates in a stated order, where the first
conflict lands, so an integration order or a resolution can be chosen before
any merge is attempted.

Owner intent, conversation 2026-09-25: four or five items implemented in
parallel end in review together, and their conflicts surface only when
`integrate` refuses one, which starts a manual feedback and relaunch loop
the owner wants to avoid. Prediction is the read-only first layer of that
wish; [G-178](G-178-candidate-target-update.md) acts on it by hand and
[G-180](G-180-policy-driven-integration.md) by policy.

## Constraints

### Observed evidence

At main `47852e3`, [`integrate`](../internal/integrate/integrate.go) merges
the inspected commit with a plain `git merge`, aborts on conflict, and
refuses with the target unchanged and the record left in review. Reproduced
on a disposable repository with the built binary:

```text
merge of worktree-G-004 into main refused: CONFLICT (content): Merge conflict in shared.txt; Automatic merge failed; fix conflicts and then commit the result.; main is unchanged at 762e74b and G-004 stays in review
```

The refusal names the file and no next action, and `integrate` runs no
verification after a merge, so two candidates that each merge cleanly and
fail together are not detected by it.

`git merge-tree --write-tree` (Git 2.38 and later; this machine has 2.55.0)
performs the merge read-only, writes only objects, names the conflicting
files with `--name-only`, and exits 1 on conflict. On the same repository
after G-003 was integrated it printed `shared.txt` and `CONFLICT (content):
Merge conflict in shared.txt` with no checkout touched. Nothing in Grove
uses it.

The board's Review block ([board](../docs/board.md)) shows whether the
target holds the candidate and lists its changed files against the target.
The G-161 candidate on `worktree-G-161` (in review at `83e7f38`) adds `grove
deps` and the board's dependency view, whose delivery text reads `awaiting
review; candidate X not in HEAD, not on main`, computed in
`internal/deps/deps.go` (`Deliver`) through one `merge-base --is-ancestor`
per candidate on demand. The predecessor's close skill
(`../skills/skills/close/SKILL.md`, read as a file) left conflicts to the
person: "If the merge conflicts, resolve or rebase, rerun … then retry."

### Proposed design and scope

- One read-only fact per candidate in review, relative to the target's
  current tip: integrated, merges cleanly, merges cleanly although the
  target moved since the branch's merge base, or conflicts in named files.
  Compute it with `git merge-tree --write-tree` through `repo.Command`, on
  demand for the candidates a view names, never during the board load
  (G-030, G-031). The fact names the target commit it was computed against
  and is stale, visibly, when the target moves.
- Show it where candidates are already explained: the board's Review block,
  `integrate`'s refusal (which then names the next action: G-178's operation
  or the manual merge), and, once G-161 is integrated, the `deps` delivery
  text and the selection preview. Choose the noninteractive contract in
  preparation; extending `versions --json` or `deps --json` is preferred to
  a new command.
- For several candidates and an order the person or `context` states,
  simulate that order: each merge-tree result becomes a throwaway commit
  object (`git commit-tree`, no ref, no checkout) that the next merge reads,
  and the report says where the first conflict lands and in which files.
  Grove neither chooses the order nor calls a clean order safe: it is not
  evidence that the changes are semantically compatible.
- Out of scope: running verification on a predicted merge, which needs a
  checkout and belongs to G-180; starting anything; changing any record or
  ref; a same-branch shared candidate for several IDs, which
  [G-163](G-163-selected-work-review-boundary.md) governs.

## Acceptance

1. For a candidate in review the owner sees, in the board's Review block and
   noninteractively, that it is integrated, merges cleanly into the target
   at a named commit, or conflicts in named files. The same source gives the
   same answer in both.
2. A refused `integrate` names the conflicting files and the exact next
   action, and the target is unchanged.
3. With several candidates and a stated order, the report names where the
   first conflict lands and in which files, and says that a clean order is
   not evidence of compatibility. Grove never chooses the order.
4. Nothing is checked out, written to a ref, or changed in a record; reads
   stay one process per candidate on demand; a target that moved after the
   fact was computed is shown as such.
5. Exercised on: fast-forward, clean merge after the target moved, textual
   conflict, already integrated, and a target that moved between prediction
   and integration.

## Next

Assignable now through `$grove-work G-177` once committed. If G-161 is
integrated first, extend its `deps` delivery text and preview; otherwise
the Review block and the refusal alone, with `deps` reconciled when it
lands. G-178 and G-180 act on this fact. No implementation has been
assigned by this shaping session.
