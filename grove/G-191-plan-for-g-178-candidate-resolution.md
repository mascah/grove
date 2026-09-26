---
id: "G-191"
type: plan
title: "Plan for G-178: resolve a conflicting candidate through one bounded attempt"
status: current
created: "2026-09-26T01:18:08Z"
updated: "2026-09-26T01:19:13Z"
work: ["G-178"]
---

## Design

Plan for [G-178](G-178-candidate-target-update.md), written headless from
main `e812672` (G-177 integrated), on `worktree-G-178`.

**Name.** `grove resolve ID` and the board's `m` on a record in review. It
matches G-180's `policy.resolve` and G-182's "resolution attempt". It avoids
`reconcile` and `refresh`, which the record rules out, and `update`, which is
taken. The board's existing internal "resolve" of an answered question keeps
its name. The new prompt's kind is `conflict`.

**One function for both triggers.** `attempt.Resolve(req, shown, now,
report)` in `internal/attempt/resolve.go` holds the refusals and the
mandate, so G-180's policy calls the same function. It does the following,
in order, and every refusal comes before anything is written:

1. One work ID, no `--until`, `--branch` or `--worktree`: the operation
   chooses where it runs.
2. The project names a target. `versions.Inspect` finds the one committed
   branch, not the target, that holds the record in review, and that
   branch's single live checkout. Feedback commits there, as `f` does.
3. `versions.PredictContext` against `refs/heads/TARGET` with the
   candidate. Anything but `conflict` is refused with the fact's own text.
   When the caller passes the fact it showed, `shown`, a different
   candidate, or a target commit other than the one read now, is refused
   ("changed since the fact was computed; look again").
4. The group sharing the candidate (G-188) is the selection, the given ID
   first. No member may have a running or orphaned attempt. The launch
   defaults (`Defaulted` with the branch's `run:`) must supply a budget and
   a mode. After review round 1, what `Start` would refuse is asked here,
   before the feedback, of the branch's checkout with the group active: a
   wait, the entrypoints, and the provider and its `--version`.
5. `update.Feedback` in the branch's checkout, with `Mandate(fact,
   target, candidate)` as its text. The text names the target commit in
   full and the conflicting files. It says to merge that commit rather
   than rebase or merge a later tip, resolve those files, rerun
   verification, and hand off the merge as the new candidate with the
   previous candidate, the merged commit and the resolved files in
   Evidence. It says to change nothing else and to stop with a checkpoint
   naming any choice the record does not settle.
6. `Start` with `Root` set to the branch's checkout, and `Branch` and
   `Worktree` set to it (a reuse), with the group's IDs. If it fails after
   the feedback, the error says the feedback stands and gives the
   `grove run … --branch … --worktree …` that launches it. Start's own
   launch lock and running check remain the duplicate guard. A second
   concurrent `resolve` finds the record already active and is refused.

The assignment stays `/grove-work IDS --interaction headless`. The mandate
travels in the record, which is where the guide reads a checkpoint and
feedback, so no new assignment input is needed.

**Guide.** Step 5 gains "A target that moved": merge the named commit,
never rebase; resolve only the named files; verify; hand off; and scope
the review to the resolution (`git show --cc MERGE`, the previous
candidate and its reviews). A choice the record does not settle is a
missing human decision. The judging section adds `grove resolve` as the
disposition for a conflict. `integrate`'s refusal names it as the next
action.

**Judging the resolution.** `versions.Changes` gains `Resolution {merge,
target, previous, files}`:

- `merge` is the latest first-parent merge on the candidate not on the
  target (`rev-list --first-parent --merges --parents -n1 CAND ^TARGET`)
  whose second parent the target contains.
- `files` is the conflicts that merging the merge's two parents again
  (`merge-tree`, in objects only) reports. Each file also says whether the
  result is one side's content. Adjusted after review round 1: `diff-tree
  --cc` missed a conflict settled by taking one side, which drops the other
  side's change. It also listed files that Git had merged by itself.
- `previous` is the `candidate` the record held at the merge's first
  parent, since feedback keeps it.

The Review block prints a row for it, and the Changes list marks resolved
files. It is read on demand with the rest of the changes, never in the
board load. The earlier candidate and every review's `examined` stay
ancestors, because the resolution is a merge.

**Docs.** `docs/commands.md` (`resolve`), `docs/board.md` (`m`, the
resolution row), `docs/work-execution.md`, and `grove --help`.

## Steps

- [x] `attempt.Resolve` and `Mandate`, with refusal tests.
- [x] Fake-provider lifecycle tests: a clean resolution, one that needs a
      choice, a target that moves during the attempt, and a Stop.
- [x] `grove resolve` in the CLI and in the usage text; `integrate`'s
      refusal names it.
- [x] `versions.Changes.Resolution`, with a test on a real repository.
- [x] The board: the `m` prompt through the launch line, the resolution
      row, and the resolved-file marks, with tests.
- [x] The guide, `commands.md` and `board.md`; then verification,
      independent review and the handoff.

Out of scope, as the record says: automatic triggering (G-180), and any
change to what candidate, approval or integration mean. A real-provider
trial is bounded separately and is not run here.
