---
id: "G-183"
type: plan
title: "Plan for G-177: merge prediction for candidates in review"
status: current
created: "2026-09-25T23:03:54Z"
updated: "2026-09-25T23:04:15Z"
work: ["G-177"]
---

Plan for [G-177](G-177-merge-prediction.md), prepared headless in
`worktree-G-177` from main `6b14141`, where G-161 is integrated, so `deps`
and the board's selection preview are in scope. Single implementer; one
independent review on the final revision.

## Design

- **One source.** `internal/versions/merge.go` owns the fact: `Merge`
  (`target` commit it was computed against, `commit`, `outcome` one of
  `integrated`, `fast-forward`, `clean`, `conflict`, and `conflicts`, the
  files) and `Text(targetName)`, the one wording every surface prints.
  `PredictContext(ctx, root, target, commits)` resolves the target to a
  commit once, then for each commit in order: merge base equal to the commit
  is `integrated`, equal to the target is `fast-forward`, otherwise
  `git merge-tree --write-tree --name-only --no-messages -z` gives `clean`
  or, on exit 1, `conflict` with the named files. With several commits, a
  clean result becomes a throwaway `git commit-tree` object (no ref, no
  checkout) that the next merge reads, and the sequence stops after the
  first conflict. Everything goes through `repo.GitContext`.
- **Board Review block.** `ChangesContext` resolves the target to a commit
  first, reads every fact against that commit, derives `OnTarget` from the
  merge base (dropping one `merge-base --is-ancestor`), and fills
  `Changes.Merge` through the same classifier: one `rev-parse` plus a
  `merge-tree` only when neither ancestry answers. The block shows
  `Merge.Text` in place of `on/not on TARGET`, and, when the board's own
  reading of the target tip differs from the fact's commit, says the target
  moved and that `r` re-reads. Still on demand, only while the card is open.
- **Noninteractive contract.** `deps` (CLI and the board's preview) extends
  `Deliver` with a prediction function: each item in review whose candidate
  is not on the target gets `merge` in `--json` and `Text` appended to its
  delivery. A selection with two or more such candidates also gets
  `merge_order` in `--json` and a note naming, in the preview's order (the
  person's, constrained only by `depends_on`), each outcome up to the first
  conflict and its files, and saying Grove chose no order and a clean order
  is not evidence of compatibility. No new command; `versions --json` is
  unchanged.
- **Integrate.** Before merging, predict the approved commit against HEAD.
  A conflict refuses before anything changes, naming the target commit, the
  files and the next action: bring the target into the branch in its
  checkout (`git merge TARGET`, resolve, commit, which is a new candidate to
  hand to review), or `grove feedback ID "…"` there to return it to an
  implementer. The `git merge` refusal path stays as the fallback for what
  prediction cannot see (a dirty or untracked file in the way).
- **Not in scope**, per the record: verification of a predicted merge,
  starting anything, writing refs or records, same-branch shared
  candidates.

## Steps

1. `versions/merge.go` with tests over a fixture repository: fast-forward,
   clean after the target moved, textual conflict, integrated, a sequence
   whose second merge conflicts only because of the first, and no ref or
   checkout change after any of them.
2. `ChangesContext` on the target commit with `Merge`; board Review block
   text and the moved-target note; TUI tests with the fake backend.
3. `deps.Deliver` with prediction, `merge` and `merge_order` JSON, text
   note; CLI tests; board preview wired through a `Predict` backend entry.
4. `integrate` pre-merge prediction and refusal; test that a target moved
   between prediction and integration is refused with the files and the
   target unchanged.
5. Docs: `docs/board.md` Review block and preview, `docs/commands.md`
   `deps` and `integrate`, and the integrate sentence of
   `docs/work-execution.md`.
6. Verification per `CLAUDE.md`, independent review, handoff.
