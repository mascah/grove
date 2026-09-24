---
id: "G-131"
type: plan
title: "Plan for G-125: answer a blocking question from the board"
status: current
created: "2026-09-24T15:06:48Z"
updated: "2026-09-24T15:07:18Z"
work: ["G-125"]
---

## Design

Base: `main` at `c38d914`, branch `worktree-G-125`, G-125 at revision
`sha256:5db1849f…`. One implementer; `internal/tui` only, plus
`docs/board.md`.

- **Checkout.** `judgeRoot` in `review.go` already finds the one clean
  checkout on the branch of a shown version. It becomes `checkoutOf(g, v,
  doing)`, shared by `a`/`f` and `e`, returning the live version it found and
  refusing, besides its present cases, when more than one checkout is on the
  branch (ambiguous). Messages keep their words for judging.
- **Freshness.** Before the editor, `e` reads the file in that checkout and
  refuses when its revision is not the one the board read; nothing is written
  then.
- **Answer heading.** When the body has no `## Answer` line, `e` appends one
  before the editor, and restores the original bytes if the editor leaves the
  file exactly as it was given, so a quit without an answer writes nothing.
  Kept: G-121's answer landed under the agent's `## Next`.
- **Suspend and resume.** A new `Backend.Edit(path, done)` returns the
  `tea.Cmd`; `Run` sets it to `tea.ExecProcess` of `$VISUAL`, else `$EDITOR`,
  else `vi`, with the terminal's own input and screen files as the editor's
  stdin, stdout and stderr (the program's output is the `watched` wrapper,
  which is not a terminal to the child). Tests fake it by writing the file.
- **Resolve.** On return the file is read again: an editor error or no change
  says so and writes nothing more; otherwise a `resolve` prompt (y/n) opens,
  and the board re-reads. `y` runs a new `Backend.Answer(root, id, expect)`:
  `update.Apply` with `status=resolved`, `Commit`, and `Expect` the revision
  read after the editor, so the commit holds the answer and the status
  together. `n` or Esc leaves the edit uncommitted and says where.
- **From an attempt.** `e` on the attempt screen or the attempts list, for an
  attempt waiting on a question, opens that question's detail as `o` opens
  work (Esc returns to the attempt) and starts the same edit.
- **After the answer.** The latest attempt of work whose blocking question was
  created before the attempt ended and resolved after it gets the outcome
  `question answered`, and stays under Needs you as `question answered: R
  again`, like `feedback given: R again`: the relaunch is still the owner's
  to make. This is a deliberate reading of the proposed design's "leaves
  Needs you"; acceptance 1 asks only that the work show as launchable with
  `R`, which its detail does.

## Steps

1. `review.go`: `checkoutOf` with the ambiguity refusal; `judgeRoot` through it.
2. `answer.go`: `e` (freshness, heading, edit), the edited message, the
   `resolve` prompt and action; `Backend.Edit` and `Backend.Answer`; `Run`
   and `Live` wiring; key routing on the detail, attempt and attempts screens;
   the question detail's header row and the footer hint.
3. `attempts.go`: the `answered` outcome and standing; `e` in the question
   standing's Next.
4. Tests in `internal/tui`: resolve, decline, stale revision, editor error,
   ambiguous and missing checkout, from the attempt screen.
5. `terminal.py`: a scenario that suspends the board for a fake editor
   (`VISUAL` a script writing the answer) and resumes it, then declines.
6. `docs/board.md`: the selection contract sentence, a section on answering
   a question, the key table.
7. The checks in AGENTS.md; a demonstration on a real attempt of Grove's own
   work, or its absence reported.
