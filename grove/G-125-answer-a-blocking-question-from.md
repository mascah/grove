---
id: "G-125"
type: work
title: "Answer a blocking question from the board"
status: done
created: "2026-09-24T01:23:28Z"
updated: "2026-09-24T21:03:16Z"
kind: feature
size: medium
relates_to: ["G-044", "G-046", "G-079", "G-114", "G-123"]
candidate: "6833eb5c0cb1c19ef6d80f04a3ceba6feedf210d"
approved: "6833eb5c0cb1c19ef6d80f04a3ceba6feedf210d"
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
  question file; G-121's answer landed under the agent's `## Evidence

Branch `worktree-G-125`, base `main` at `c38d914`, started from this record
at `sha256:5db1849f…` and [plan G-131](G-131-plan-for-g-125-answer-a-blocking.md)
at `sha256:ce915f85…` (`dc3c981`). Implementation `e1d3339`, review fixes
`052cc74` and `606330c`; the candidate is the commit holding this evidence.
Headless attempt `G-125.20260924T150358Z`.

What changed (`internal/tui`, `docs/board.md`):

- `e` on an open question's detail, or on the attempts list or one attempt
  whose outcome is waiting on a question (which opens that question's
  detail above the attempt, as `o` opens work), answers it (`answer.go`).
  The checkout is found by `checkoutOf` (`review.go`), now shared with `a`
  and `f`: the one valid checkout on the branch of the shown version,
  refusing none or several. For judging its copy must match HEAD, as
  before; for answering an uncommitted copy is the owner's answer so far.
  The file on disk must have the revision the board read.
- The board appends `## Answer` when the body has none, then suspends
  through `tea.ExecProcess`: `sh -c '$VISUAL "$@"'` (else `$EDITOR`, else
  `vi`, as Git runs it) on the file, with the terminal files themselves as
  stdio, since the program's output is the `watched` wrapper. On exit the
  file is read again. An unused heading is taken back, also by `Run` when a
  hangup ended the session under the editor. An editor error, or no change,
  says so and writes nothing more; a failed take-back is reported with no
  prompt. After an edit, or on an earlier uncommitted one, the board
  re-reads and asks `Resolve ID and commit it with your answer on branch B?
  y/n`. `y` runs `update.Apply` with `status=resolved`, `Commit`, and
  `Expect` the revision after the editor, so one commit of the question's
  file alone holds the answer and the status. The result names `R` on each
  work the question blocks. `n` or Esc leaves the edit uncommitted and says
  where; `e` reopens it.
- The latest attempt of work blocked by a question that is now resolved, not
  asked after the attempt ended and last `updated` from that second on, is
  `answered since`: `question answered: R again`, still under Needs you like
  `feedback given: R again`, since the relaunch is the owner's to make. This
  is a deliberate reading of the design's "leaves Needs you"; acceptance 1
  asks only for `R` to be offered, which the work's detail does.
- The detail of an open question names where `e` writes, or why it cannot.
  The footers hint `e`. docs/board.md adds *Answering a question*, its key
  row, the outcome and the Needs you case, and changes the selection
  contract to "only `e` starts the owner's editor, on an open question".

Against acceptance:

1. From the attempt screen, one key opens the question in the owner's editor
   in the branch's checkout. After the editor quits, the board offers to
   resolve and commit. `y` sets it resolved and commits on that branch, and
   the attempt shows `question answered: R again`; `o`, then `R`, launches.
   All of this is exercised by `terminal.py`'s `attempt_lifecycle`, which
   runs the real binary in a pseudo-terminal against a fake provider and a
   fake `VISUAL`. That scenario checks:
   - the editor got a terminal in canonical mode, on the branch's copy;
   - the board came back on the alternate screen;
   - one commit `docs(G-002): set status=resolved` holds only the question,
     with its answer;
   - the next `R` continues on that branch to a candidate.
   **Not done: the demonstration on a real attempt of Grove's own work.** A
   headless session cannot answer a real question, and none is open (`grove
   attempts`, `grove list` at `606330c`). It is left to the owner; see Next.
2. Refusals are shown and write nothing: no checkout, two checkouts, a file
   changed since the read, a resolved question, a record that is not a
   question, an editor error (`TestAnswerRefusals`,
   `TestAnswerDeclinedUnsavedOrFailed`).
3. `internal/tui/answer_test.go` covers the following, with a fake `Edit`
   and `Answer`:
   - the resolve path;
   - the decline path, then reopening the answer;
   - a stale revision, both before the editor and refused at resolve;
   - the take-back, and a failed take-back;
   - the path from the attempt;
   - the answered-since window.

   `terminal.py` covers suspending for the editor and resuming.
4. docs/board.md is reconciled as above. At `606330c` these all passed:
   - `go vet ./...`;
   - `gofmt -l .` (empty);
   - `go run ./cmd/grove check` (`OK: 128 records`);
   - `go test -count=1 -timeout 120s ./...` (every package ok, `internal/tui`
     16.1s with `TestTerminal`);
   - `python3 internal/tui/testdata/terminal.py BIN`: 11 of 11 ok, run by
     `TestTerminal` in that suite and directly at `052cc74`.

Review: [G-133](G-133-g-125-review-answering-a-questio.md), an independent
subagent reviewer over three rounds. It found no blockers. Every finding was
fixed, or kept and documented with its reason.

Limits: `answered since` rests on `updated`, as docs/board.md says; `e`
starts only on Unix (`sh`), as Grove does.

## Next

In review with the candidate this record names. The integrator's actions:

1. Demonstrate acceptance 1 on a real attempt when one next stops on a
   question: in the board, `A`, Enter on the waiting attempt, `e`, write the
   answer, quit the editor, `y`; the attempt shows `question answered: R
   again`, and `o` then `R` relaunches. Report it here, or approve with it
   noted as still owed.
2. `go run ./cmd/grove approve G-125 "VERDICT"` in this checkout
   (`.claude/worktrees/worktree-G-125`), then `go run ./cmd/grove integrate
   G-125` in the `main` checkout; or `go run ./cmd/grove feedback G-125
   "TEXT"` here.

Verdict on candidate 6833eb5, 2026-09-24: lgtm
