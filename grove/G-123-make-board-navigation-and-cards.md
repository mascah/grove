---
id: "G-123"
type: work
title: "Make board navigation and cards quicker to read and move through"
status: active
created: "2026-09-24T01:23:27Z"
updated: "2026-09-24T04:24:10Z"
kind: feature
size: small
relates_to: ["G-043", "G-109", "G-124", "G-125"]
---

## Outcome

The board and a record's detail take fewer keystrokes to move through and
say at a glance which work is running: opening a record that is already on
the path returns to it instead of nesting deeper, a breadcrumb shows the
path, the column keys skip empty columns, a card with a live attempt is told
apart by more than a tag, and the sidebar can be hidden to read or copy the
content.

Owner intent, from notes taken while using the board on 2026-09-23 and
confirmed in the shaping conversation of the same day.

## Constraints

Observed at main `28f5ddc`:

- The detail keeps the open records as a list of IDs. `openDetail`
  ([detail.go](../internal/tui/detail.go)) appends every time and
  `leaveDetail` pops one, so A → B → A is three deep and Esc walks every
  step. A record opened with `o` from an attempt returns there through
  `workDepth` and `workBack`, which any change to the stack must keep.
- ←/→ and h/l step one visible column at a time
  ([model.go](../internal/tui/model.go), `boardKey`), and `visible()` includes
  empty columns, so with nothing active or in review, Done is three presses
  from Proposed.
- `attemptTag` ([attempts.go](../internal/tui/attempts.go)) yields `● running`
  or `● orphaned`, joined to the state tag with a space in `cards()`, so a
  running card reads `● running not on main`. An attempt runs on a branch,
  so the second tag is nearly always true while the first shows. Only focus
  changes a card's border ([view.go](../internal/tui/view.go), `cardBox`).
- From 100 columns the detail splits content and sidebar at 11/20 with no way
  to hide the sidebar (`detailBody`). No mouse mode is enabled, so the
  terminal's own selection works, and it takes both panes because they share
  rows. Linked records are already grouped by role in the order `roles`
  gives; `related` mixes types.
- The board redraws on a 2 s tick while an attempt runs;
  `charm.land/bubbles/v2`, which has a spinner, is not a dependency.

In scope, as proposed design:

1. Opening an ID already on the stack cuts the stack back to it. A breadcrumb
   row in the detail header, such as `board › G-108 › G-115`, clipped from
   the left when long. The `o` return path still works.
2. ←/→ and h/l skip empty columns, which are still drawn.
3. A card with a live attempt gets a distinct border colour or marker, and
   `not on main` is not shown beside `● running`. A spinner is optional:
   only if the owner wants motion after seeing the colour, and only with the
   redraw cost measured.
4. One key hides and shows the sidebar on a wide terminal, and the content
   takes the width. Narrow behaviour (one pane, Tab cycles) is unchanged.

Out of scope: regrouping linked records by type (the owner leaned to keeping
roles in the shaping conversation), mouse support, and any change to what the
board reads.

## Acceptance

1. From the board, Enter on A, then on B in A's sidebar, then on A in B's
   sidebar shows A with a path two deep, and one Esc returns to the board. A
   test in `internal/tui` covers the cut and the `o` return.
2. With Active and Review empty, one → from Proposed lands on Done. A test
   covers it.
3. In an actual terminal the owner can tell a running card at a glance, and
   no card shows `not on main` beside `● running`. The owner's judgment in a
   terminal is the check, as the brief asks for visual work.
4. On a wide terminal one key hides the sidebar, the content fills the width,
   and a mouse selection of the content takes no sidebar text. The key is in
   the table in [docs/board.md](../docs/board.md).
5. docs/board.md describes each changed behaviour. `go run ./cmd/grove check`,
   `go vet ./...`, `gofmt -l .`, one full `go test -count=1 -timeout 120s
   ./...` and the terminal lifecycle check pass.

## Next

Assign with `/grove-work G-123`. Small, judged visually.
[G-124](G-124-keep-the-board-fresh-without-pre.md) and
[G-125](G-125-answer-a-blocking-question-from.md) change the same package;
run them one after another, not as parallel attempts.
