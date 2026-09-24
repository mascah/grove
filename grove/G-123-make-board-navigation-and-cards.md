---
id: "G-123"
type: work
title: "Make board navigation and cards quicker to read and move through"
status: review
created: "2026-09-24T01:23:27Z"
updated: "2026-09-24T04:36:21Z"
kind: feature
size: small
relates_to: ["G-043", "G-109", "G-124", "G-125"]
candidate: "ff15e25d4cf5c0b85d2c6ab618984dbf5cf726bd"
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

## Evidence

Implemented on `worktree-G-123`, base main `f81f7e9`, from G-123 revision
`sha256:2056133c…` (no plan: small work, and the proposed design in
Constraints was specific enough to build without one). Code is in `46f1322`,
with review fixes in `6e78245`. The candidate is the commit holding this
text.

Against each acceptance item:

1. **Stack cut.** In `openDetail` (`internal/tui/detail.go`), opening from a
   detail a record already on the path cuts the path back to that record.
   `o` from an attempt still opens a new layer, so Esc returns to the
   attempt. If a cut drops that layer, its return goes too. `settleFocus`
   now keeps the layer's depth right when a refresh drops a record below it.
   A breadcrumb is the first row of the detail header, such as
   `board › W-001 › attempts › W-001`, with the start cut when too long. It
   is left out in the compact header, below 16 rows.
   `TestReopeningCutsThePath` covers A → B → A (two deep, and one Esc to the
   board), the `o` layer and its return, a cut below it, and the refresh.
   The test fails with the cut disabled, and fails with the `settleFocus`
   fix removed.
2. **Column skip.** ←/→ and `h` `l` in `boardKey` (`model.go`) move to the
   next column with cards. Empty columns are still drawn, and past the last
   column with cards the focus stays. `TestColumnKeysSkipEmptyColumns`
   covers this at 120 and 80 columns: one → from Proposed lands on Done.
3. **Running card.** A card with a live attempt has a cyan border
   (`runningAccent`, ANSI 6, a colour no column uses), and `currentCards`
   leaves off `not on main` while the attempt is live.
   `TestRunningCardStandsOut` checks both, and checks that both revert when
   the attempt ends. An orphan is treated the same, with its text tag
   `● orphaned`. No spinner: it was optional, and only if the owner wants
   motion after seeing the colour. **Open: the owner's judgment in a
   terminal.**
4. **Hide the sidebar.** `w` in a detail, from 100 columns, hides the
   sidebar: `split()` then gives the content the full width and uses the
   one-pane layout, so no row holds sidebar text, and Tab still reaches the
   sidebar alone. `w` shows it again. The choice lasts for the session and
   does nothing below 100 columns. The hint shows only where the key acts.
   `TestHidingTheSidebar` covers this. The key is in the docs/board.md table.
   **Open: a real mouse selection** is for the owner to try, since the tests
   check only the rendered rows.
5. **Docs and checks.** docs/board.md now covers each changed behaviour:
   Columns and cards, Record detail, Attempts, and the keys table. At
   `6e78245`:
   - `go vet ./...`: clean.
   - `gofmt -l .`: empty.
   - `go run ./cmd/grove check`: `OK: 122 records`.
   - `go test -count=1 -timeout 120s ./...`: passed at `46f1322`. At
     `6e78245`, every package passed except `internal/tui`, where one run of
     `TestTerminal` failed in `attempt_lifecycle` (`never drew 'resolved'`).
     The screen showed `a read is already in progress`: the script's single
     `r` came while a re-read from before its answer commit was still
     running, under full-suite load.
   - `go test -count=1 -timeout 120s ./internal/tui` alone: passed three
     times in a row, in about 8.7 s each. Of that, `TestTerminal` takes
     8.4 s, as it does on main.
   - `python3 internal/tui/testdata/terminal.py BINARY`: all ten checks ok,
     at both `46f1322` and `6e78245`.

Review: [G-127](G-127-review-of-g-123-board-navigation.md), an independent
reviewer subagent, two rounds. Round 1 found one should-fix (the stale depth
of the `o` layer after a refresh) and four nits. Four were fixed in
`6e78245`, and the orphan border was kept for the owner. Round 2 found
nothing new.

Limits: the running-card colour, and a mouse selection with the sidebar
hidden, are unjudged by a person. The `attempt_lifecycle` script race is
not fixed here, since it is outside this work's scope.

## Next

**Handoff, 2026-09-24 (headless).** G-123 alone, on `worktree-G-123` in
`.claude/worktrees/worktree-G-123`, base main `f81f7e9`. Main is now
`25525cf`, which adds only a G-114 record edit; that touches none of these
files. No command is still running.

For the owner's judgment, run `go run ./cmd/grove` in this checkout, from a
terminal at least 100 columns wide, with an attempt running or a card in a
detail path:

- Open a card, open a linked record, then open the first again. The path
  should cut back, and the breadcrumb row should read well.
- → from Proposed should skip empty columns.
- Decide whether a cyan-bordered card with `● running` stands out enough,
  and whether an orphan should share that border.
- `w` in a detail, then select the content with the mouse.

Integration: `go run ./cmd/grove approve G-123 "VERDICT"` in this checkout,
then `go run ./cmd/grove integrate G-123` in the `main` checkout. Run
[G-124](G-124-keep-the-board-fresh-without-pre.md) and
[G-125](G-125-answer-a-blocking-question-from.md) one after another, since
they change the same package.
