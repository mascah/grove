---
id: "G-166"
type: question
title: "Which dependency-view layout should G-161's board implement?"
status: open
created: "2026-09-25T20:45:48Z"
updated: "2026-09-25T20:47:55Z"
blocks: ["G-161"]
relates_to: ["G-165"]
---

## Question

[G-161](G-161-dependency-view.md) says the owner judges concrete terminal
layouts before one is settled. Plan
[G-165](G-165-g-161-dependency-view-plan.md) draws both options at 120
columns from a synthetic ten-item unfinished backlog (chain, shared
prerequisite, convergence, unrelated work, a review candidate off `main`, an
abandoned prerequisite, a blocking question) and describes them at 80.
Answer two parts:

1. **Layout.** A, lanes: one column per layer, cards naming what each needs
   and unlocks, ←/→ scrolling by layer when the layers do not fit.
   B, layered list with a focus tree: one row per unfinished item, grouped by
   connected group and indented by layer, with the focused item's upstream
   and downstream trees beside it from 100 columns, or behind Tab below that.
2. **Handoff.** Proposed: Enter on the preview prints the ordered
   `grove context IDS` for the selection on stdout and exits, as selecting a
   version prints its workspace; the command rereads everything, so a stale
   preview cannot pass on stale facts. The alternative is no handoff: the
   preview only informs.

## Evidence

G-165's Observed section: this repository's largest group is 13 done items
over seven layers, so a lane layout needs four 80-column screens across it;
unfinished work here has almost no edges, so the drawings use a synthetic
backlog. The detail sidebar already lists direct `needs` and `needed by`, but
nothing shows more than one hop.

## Recommendation

B with the handoff as proposed. A list reads the same at 80 columns and at
any project size, and the tree answers "what does this need, and what does it
unlock" for the focused item without lines drawn across the screen. A shows
convergence at a glance while the backlog is small and shallow, and costs
horizontal scrolling as soon as it is not.

The noninteractive `grove deps` command, the shared model it prints, and the
shaping-guide change do not depend on this answer and are being built on
`worktree-G-161` meanwhile.

## Next

Open; blocks G-161's board steps (G-165 step 4). Answer here (for example
"B as drawn", or what changes), set `status=resolved`, and relaunch G-161;
the board is then built in the chosen layout.
