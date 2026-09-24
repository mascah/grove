---
id: "G-124"
type: work
title: "Keep the board fresh without pressing r"
status: review
created: "2026-09-24T01:23:28Z"
updated: "2026-09-24T14:57:55Z"
kind: feature
size: small
relates_to: ["G-046", "G-109", "G-123"]
candidate: "3d4597b02229ae36c623846959e6f49ff048d469"
approved: "3d4597b02229ae36c623846959e6f49ff048d469"
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

## Evidence

Implemented on `worktree-G-124`, base main `2491eda`, from G-124 revision
`sha256:72e2d5d7…`. There is no plan: the work is small, and the proposed
design in Constraints was specific enough to build from. The code is in
`a6807a4`, the review fixes in `70509fb`, and a docs rewording in `eba3197`.
The candidate is the commit holding this text.

Against each acceptance item:

1. **Focus.** `Model.View` sets `ReportFocus`. On `tea.FocusMsg`, the board
   calls `refresh()`, which is what `r` calls, unless:
   - the session is ending,
   - `busy()` (an inspect, resolve or action is under way), or
   - `pinned()` (a detail left at a timeline commit or a diff, which a
     re-read would close; finding 1 of [G-130](G-130-review-of-g-124-board-re-reads-o.md)).

   Blur does nothing. The header now reads `read 2 branches, 2 checkouts at
   14:05:06`, from `readAt`, set when a read lands.

   `TestFocusRereadsTheBoard` sends `FocusMsg` and checks that a read
   starts and the header time shows. It also checks that no read starts
   while one is pending, or with `asOf` or `diff` set, and that blur starts
   nothing. The test fails with the focus handling disabled.

   In a real pty, the new `focus_rereads` scenario in
   `internal/tui/testdata/terminal.py` checks four things:
   - The board emits `\x1b[?1004h`.
   - Focus-in (`\x1b[I`) after a new branch is made starts a read, and `s`
     lists the new branch.
   - `\x1b[?1004l` comes last on exit.
   - The terminal modes are restored.

   It fails against the base build (`2491eda`).
2. **Moved tip.** `versions.TipsContext` lists `refs/heads` through the
   same `listBranches` (one `for-each-ref` via `repo.GitContext`) that
   Inspect starts from. `Backend.Tips` is read beside the attempts, only
   while `m.running()`. `gotAttempts` re-reads the board when the tips differ
   from the committed sources of the last read, but not while an inspect,
   resolve or action is under way, and not while pinned. It never restarts
   an inspection under way, so a slow read cannot starve, and the next poll
   compares again with what that read found.

   `TestMovedTipRereadsTheBoard` covers:
   - unchanged tips, which re-read nothing;
   - a moved tip, which re-reads once;
   - a pinned detail, which waits;
   - an inspection under way, which is not restarted;
   - the finished read, which then matches.

   It fails with the `moved` term removed, and fails without the `pending
   != "inspect"` guard. In the pty `attempt_lifecycle` scenario, a status
   committed on the running attempt's branch moves the card to `Active 1`
   with no key pressed. That also fails against the base build.
3. **Nothing when idle.** With no attempt live, no tick is scheduled, and
   `Tips` is never called. The same test asserts both (`ticking` false, 0
   listings). Under `-short`, `internal/tui` runs in 0.5 s. The pty test was
   already skipped under `-short` (G-071).
4. **Docs.** [docs/board.md](../docs/board.md) says when the board re-reads
   by itself: on focus with its exceptions, on an attempt's end or a moved
   tip, and nothing else without a key. It also says that the header shows
   the read time.

Verification at `5b63579`, whose code matches the candidate. Every command
passed:
- `gofmt -l .`: clean.
- `go vet ./...`: clean.
- `go run ./cmd/grove check`: OK, 125 records.
- `go test -count=1 -timeout 120s ./...`: every package ok, including
  `TestTerminal`, whose `focus_rereads` and `attempt_lifecycle` both pass.

The load average during that run was about 7. Under that load,
`internal/versions` took 4.4–5.7 s even under `-short`; it took 5.6 s on the
base build too, and this change adds no test there.

Review: [G-130](G-130-review-of-g-124-board-re-reads-o.md), by an
independent reviewer subagent over two rounds, examined `70509fb`. It found
nothing blocking. Its dispositions:
- The should-fix (a re-read closed an open diff) is fixed.
- Two nits are fixed.
- Focus clearing the versions cursor is kept, as `r` does.
- The round-2 nit (the pin holds beneath the attempts screen) is kept as
  behaviour and documented in `eba3197`, which was self-checked and not
  re-reviewed.

Limits:
- Only a terminal that reports focus re-reads on focus. The owner's
  Ghostty is one, but it was not driven here: the pty sends the escape
  sequence directly.
- An uncommitted edit still needs `r`.

## Next

**Handoff, 2026-09-24 (headless).** G-124 alone, on `worktree-G-124` in
`.claude/worktrees/worktree-G-124`, base main `2491eda`. No command is
still running.

For the owner's judgment, run `go run ./cmd/grove` in this checkout from
Ghostty:
- Switch to another window, commit something on a branch, and come back.
  The header's read time should change with no key pressed.
- With an attempt running, its record should move to Active while it runs.

Integration: `go run ./cmd/grove approve G-124 "VERDICT"` in this checkout,
then `go run ./cmd/grove integrate G-124` in the `main` checkout.

Verdict on candidate 3d4597b, 2026-09-24: working as expected
