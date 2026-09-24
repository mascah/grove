---
id: "G-124"
type: work
title: "Keep the board fresh without pressing r"
status: proposed
created: "2026-09-24T01:23:28Z"
updated: "2026-09-24T01:25:33Z"
kind: feature
size: small
relates_to: ["G-046", "G-109", "G-123"]
---

## Outcome

The board shows the latest state without the owner pressing `r`: it re-reads
when the terminal window regains focus and when a running attempt's branch
moves, so work an attempt sets active appears in Active while the attempt is
still running.

Owner intent, 2026-09-23: "I typically always hit `r` every time I tab back
to the window to make sure I'm looking at the latest."

## Constraints

Observed at main `28f5ddc`:

- `r` re-reads everything (`Model.Update` in
  [model.go](../internal/tui/model.go)). While any attempt is live the board
  polls the attempt files every 2 s (`pollEvery` in
  [attempts.go](../internal/tui/attempts.go)) and re-reads the board only when
  one ends (`gotAttempts`). A status the attempt commits mid-run is not seen
  until `r`: on 2026-09-24 the G-108 attempt set its record active on
  `worktree-G-108` while running, and the card stayed in Proposed until then.
- Bubble Tea v2 delivers `tea.FocusMsg` and `tea.BlurMsg` when the view sets
  `ReportFocus` (bubbletea/v2 `focus.go`, `screen.go`). `Model.View` in
  [view.go](../internal/tui/view.go) does not set it. The owner's terminal,
  Ghostty, reports focus; a terminal that does not sends nothing, and the
  board is unchanged there.
- A board load reads every branch through one `git cat-file` (G-031, G-042)
  and history and changes only while a detail is open (G-030). A re-read on
  every tick would repeat the load every 2 s.
- The last read's result holds the commit each committed source was read at
  (`versions` prints it as `Source: committed refs/heads/… HASH`).

Proposed design:

1. Set `ReportFocus`. On `FocusMsg` do what `r` does when no read is
   pending, otherwise nothing. Blur does nothing.
2. On each attempt tick, list the branch tips through one
   `git for-each-ref refs/heads` via `repo.Command` and compare them with the
   tips of the last read; re-read when one moved. Only while an attempt is
   live, so the board starts no process of its own otherwise. Live worktree
   files are not watched: an uncommitted edit still needs `r`.
3. `r` stays as it is.

Out of scope: file watching (a new dependency), a periodic refresh with no
attempt running, and re-reading a detail's timeline on focus.

## Acceptance

1. Switching away from the terminal and back re-reads the board, and the
   header's read time shows it. A test sends `FocusMsg` and sees a read start,
   and sees none while a read is pending.
2. With an attempt running, a commit on its branch that changes its record's
   status moves the card within one poll interval with no key pressed. A test
   with a fake tips reader covers a moved tip (re-read) and an unchanged one
   (no re-read).
3. With no attempt live the board schedules no tick and starts no process by
   itself. No package exceeds five seconds.
4. [docs/board.md](../docs/board.md) says when the board re-reads on its own.
   The checks in AGENTS.md pass.

## Next

Assign with `/grove-work G-124`. Small. It changes the same package as
[G-123](G-123-make-board-navigation-and-cards.md); run them one after
another.
