---
id: "G-133"
type: review
title: "G-125 review: answering a question from the board"
status: current
created: "2026-09-24T15:24:30Z"
updated: "2026-09-24T15:26:28Z"
work: ["G-125"]
examined: "606330c4a1b6fddaa05d460cd9204aac8b49e00f"
---

## Examined

An independent reviewer subagent (Claude, `reviewer` type, no edits) of this
headless session examined [G-125](G-125-answer-a-blocking-question-from.md)
against its acceptance and [plan G-131](G-131-plan-for-g-125-answer-a-blocking.md):
round 1 at `e1d3339`, round 2 at `052cc74` (`git diff e1d3339 052cc74`),
round 3 at `606330c` (`git diff 052cc74 606330c`). It ran `go test -short ./internal/tui`, `go vet`, the pty
scenario `attempt_lifecycle`, and seven mutations of the logic in a
throwaway copy, five of which the tests caught.

## Findings

Round 1 (`e1d3339`), no blockers:

1. Minor: after `n` or Esc at the resolve prompt, `e` refused the checkout's
   uncommitted answer and the detail advised discarding it.
2. Minor: `answeredSince`'s time window was untested (two mutations passed).
3. Minor: `answered since` rests on the question's `updated`; a hand edit of
   `status` leaves the old `ended without a handoff`, and a later update of
   an answered question can mark a later attempt answered.
4. Minor: a failed take-back of the unused heading opened the resolve prompt
   for an empty answer.
5. Minor: a hangup under the editor left the heading behind.
6. Nit: the judging refusals read "judg" and lost "a and f need".
7. Nit: docs/board.md said the board re-reads after every editor exit.
8. Nit: the attempt screens did not hint `e`.

## Disposition

Round 1, all fixed in `052cc74`: (1) answering accepts the owner's
uncommitted copy, still checked against the bytes the board read, and `e`
on it offers the resolve even with nothing more saved; (2) four window cases
in `TestAnswerFromTheAttempt`; (3) documented in docs/board.md, kept: an
attempt's result names no question and G-030 forbids reading history during
the load; (4) reported, no prompt, tested with a read-only file; (5) `Run`
takes back an unused heading after collecting reads; (6) the old words are
kept; (7) reworded; (8) `e answer` shows on the attempt screen and on a
waiting row of the list.

Round 2 (`052cc74`): every round-1 fix confirmed. It found that the relaxed
check for answering opens no new hole: freshness still holds, `--commit`
stages the question's file alone, and a modified live copy is the current
version. One minor finding: when an earlier uncommitted answer had no
heading, the board took the heading back but the prompt expected the
revision that still held it, so `y` always failed, writing nothing. Fixed in
`606330c` (`after = e.before` after a take-back) with a regression test,
which fails without the fix.

Round 3 (`606330c`): the fix is confirmed, and there are no remaining
consequential findings. One nit stays, without action: a change by another
process in the instant between the editor's exit and the take-back makes `y`
refuse with "changed since the expected revision", which is safe and shown.

Knowledge: no contradiction with the settled terms (G-054 to G-062) or the
accepted decisions; `question answered` needs no term. This review is
evidence, not approval.

