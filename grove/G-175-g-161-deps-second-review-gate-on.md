---
id: "G-175"
type: review
title: "G-161 deps: second review gate on the board's dependency view"
status: current
created: "2026-09-25T21:33:16Z"
updated: "2026-09-25T21:33:37Z"
work: ["G-161"]
examined: "ffaca01"
---

## Examined

The board's dependency view of [G-161](G-161-dependency-view.md), plan
[G-165](G-165-g-161-dependency-view-plan.md) step 4 in layout B as
[G-166](G-166-g-161-dependency-layout.md) answered it, and the combined
candidate `05892a2..ffaca01` for regressions against the first gate
([G-168](G-168-g-161-deps-review.md)). Round 1 examined `476b5ac`; round 2
examined the fixes in `476b5ac..ffaca01`. Two fresh `grove-reviewer` agents
did the work, read-only. Each ran:

- `go build`, `go vet` and `gofmt -l`.
- `go test -short` over `internal/tui`, `internal/deps`, `internal/cli` and
  `internal/handoff`.
- `grove check`.
- The `dependencies` scenario of `internal/tui/testdata/terminal.py`.

Each also ran throwaway probes in a `git archive` copy, deleted afterwards:

- Every width from 40 to 160 and height from 10 to 40.
- Control sequences in titles and edges.
- This repository read through `versions.InspectContext`.

## Findings

Round 1:

1. High: the trees panicked below 60 columns at height 10.
2. Medium: indentation stopped at `min(layer, w/16)`, which drew G-040
   level with G-039, its prerequisite.
3. Medium: after a failed re-read the preview stayed open. It then read Git
   on the board.
4. Medium: divergent work took the first current state, not the one the
   board places its card by, so it could vanish from `g`. The plan's
   per-item `Compare` for the overview was not built, and no reason was
   recorded.
5. Medium: a tall tree was cut off and could not be scrolled.
6. Low to medium: the "press b" advice did nothing on this screen, and a
   stale selection could not be cleared.
7. Low: connected groups were ordered by size, not in `grove deps`' order.
8. Low: `A`, then `o`, from a record opened here loses the return to the
   dependency view. The one return slot predates this work.
9. Docs: the `dependencies` crumb was undocumented.

Round 2: rounds 1 to 7 and 9 were confirmed fixed and 8 documented, with
no regression. Two new low findings:

1. The heading counted a prerequisite missing from the records read as
   "done or abandoned".
2. `b`, `s`, Esc, Esc from the dependency view leaves the chooser stuck.
   Choosing then lands on the board. The `s` trap exists on main from the
   board's own `b`.

## Disposition

Round 1's 1 to 7 and 9 were fixed in `ffaca01`, each with a regression in
`TestDepsReviewRegressions`:

- 1: the shared height is at least one row.
- 2: every layer indents.
- 3: the preview closes on a failed read and reads only on its own screen.
- 4: `earliest` places both the card and the row. The trees list each
  diverging state with its own prerequisites; this replaces the per-item
  `Compare`, and the plan says so.
- 5: Tab gives the trees focus, and they scroll.
- 6: `b` works here and returns here, and `c` clears the selection.
- 7: groups keep `grove deps`' order.
- 9: `docs/board.md` documents the crumb.

Round 1's 8 is kept and documented in `docs/board.md` as a limit.

Round 2's 1 was fixed in `867a9d8`, with its regression. It was
self-checked, not independently re-reviewed. Round 2's 2 is left as the
existing chooser behaviour and reported as a limit.

Neither round judged the layouts; acceptance item 5 is the owner's.
