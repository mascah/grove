---
id: "G-125"
type: work
title: "Answer a blocking question from the board"
status: active
created: "2026-09-24T01:23:28Z"
updated: "2026-09-24T15:07:23Z"
kind: feature
size: medium
relates_to: ["G-044", "G-046", "G-079", "G-114", "G-123"]
---

## Outcome

When an attempt stops on a question, the owner answers it without leaving
the board: open the question, edit it in their own editor, resolve and
commit it in the branch's checkout, and relaunch, from the attempt screen or
the question's detail.

Owner intent, 2026-09-23: "When an attempt stops with a blocking question,
can I just edit the file from the TUI and provide my response? and then
resolve it?"

## Constraints

Observed at main `28f5ddc`:

- The headless path of the [work guide](../docs/work-execution.md) (step 7)
  has the attempt write a question with `blocks`, checkpoint, commit and
  return the wait. The board shows the attempt as waiting on a question under
  Needs you ([docs/board.md](../docs/board.md)). Answering it today means a
  hand edit in the branch's worktree, `grove update ID --set status=resolved
  --commit` there, and then `R`.
- On `worktree-G-108`, G-118 and G-121 were written by attempts and resolved
  by the owner in commits `1d3ad82` and `e57a4c6`, each touching only the
  question file; G-121's answer landed under the agent's `## Next` heading.
  Both relaunches then succeeded, so the loop works and only the path is
  awkward.
- docs/board.md promises that selecting "never creates a worktree, edits a
  record, or starts an editor, shell, or agent; only `R` starts an agent".
  That contract changes with this work and the document is reconciled.
- Bubble Tea v2 has `tea.ExecProcess` (`exec.go`), which suspends the program
  for a child process; nothing in `internal/tui` uses it yet.
- A shown version belongs to a source. The workspace resolver binds a
  committed version to the one checkout holding its branch
  ([docs/commands.md](../docs/commands.md), `workspace`), and the review
  actions already run in the branch's checkout after a freshness check (G-044,
  [review.go](../internal/tui/review.go)).
- A question is `open` or `resolved`, and the answer stays in its body or
  links a decision ([record model](../docs/record-model.md)). `update
  --commit` exists (G-079).

Proposed design:

1. `e` on a question's detail, or on an attempt whose state is waiting on a
   question (which opens that question), resolves the checkout holding the
   shown version, refuses when there is none, it is ambiguous, or the file
   differs from what the board read, then runs `$VISUAL`, else `$EDITOR`,
   else `vi`, on the record file through `tea.ExecProcess`, and re-reads on
   return.
2. After the editor exits, a prompt offers to resolve and commit: `update
   --set status=resolved --commit` in that checkout with `--expect` the
   revision just read. Declining keeps the edit uncommitted and says so. The
   attempt row then leaves Needs you and `R` is offered.
3. Before opening the editor, append a `## Answer` heading when the body has
   none, so the answer has a home. Proposed; the implementer may drop it if
   it fights the record's own structure.
4. Bounded to question records, never code, and no text-editing widget in
   the board. The brief's "explicit commands remain noninteractive" holds:
   the editor is the owner's, and the board still only runs Grove's own
   operations.

Out of scope: editing any record type from the board, answering in an
in-board text box, resolving without an edit, relaunching automatically.

## Acceptance

1. With an attempt waiting on a question, one key from the attempt screen
   opens the question in the owner's editor in the branch's checkout; on
   save and quit the board offers to resolve and commit; accepting sets the
   question resolved, commits on that branch, and the attempt's work shows
   as launchable with `R`. Demonstrated once on a real attempt of Grove's own
   work and reported in the evidence.
2. Refusals are shown, never silent: no checkout holds the branch, the
   checkout is ambiguous, the record changed since the board read it, or the
   editor exited with an error. Nothing is written in those cases.
3. Tests in `internal/tui` with a fake exec and update cover the resolve
   path, the decline path and a stale revision. The terminal lifecycle check
   (`internal/tui/testdata/terminal.py`) covers suspending for the editor and
   resuming the board.
4. docs/board.md's contract and key table say what `e` does and where it
   writes. The checks in AGENTS.md pass.

## Next

Assign with `/grove-work G-125`. Medium: the checkout resolution and the
suspend-and-resume are the two unknowns, so a short plan record is worth
writing first. It changes the same package as
[G-123](G-123-make-board-navigation-and-cards.md) and
[G-124](G-124-keep-the-board-fresh-without-pre.md); run them one after
another.
