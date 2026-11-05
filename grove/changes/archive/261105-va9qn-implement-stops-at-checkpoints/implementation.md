# Implementation: Implement stops at checkpoints and hands off in the branch

<!-- Progress on change 261105-va9qn, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Implement stops at checkpoints and hands off in the branch

Status: done

- `internal/instruct/requires/implement.md` (both modes): what a
  checkpoint is; build a guess (`Status: to check`) or, where a wrong guess
  is costly, nothing (`Status: blocked`); a `Check:` line; leave a task
  that builds on a checkpoint unbuilt, naming it; build every other task;
  take an answer from the session, the words that started it or an
  `Answer:` line, writing `Answer: "<their words>"` and building on it;
  the report lists checkpoints first.
- `implement-owner.md`: "tell the owner how to try it and take their
  feedback before the next" is gone; the agent does not stop between
  tasks, then stops once with every checkpoint. A task that cannot be
  built is stopped at once and asked about; the open-question clause moved
  to the checkpoint rule. The stop for a wrong approach stays.
- `implement-unattended.md`: the open-question clause moved to the
  checkpoint rule; it leaves checkpoints for a later session and
  finishes; an answer that changes the change leaves the task blocked,
  naming the revision.
- `implement-cycle.md`: the report lists every change's checkpoints first.
- `implementation/guide.md` and `template.md`: `Status: to check`,
  `Check:` and `Answer:`.
- `internal/cli/walkthrough_test.go`: `TestImplementInstructions` checks
  the new text in both modes and for a cycle, and that "tell the owner how
  to try it" and "an open question in the change decides" are gone;
  `TestArchive` shows `grove archive` refusing a change whose T-01 is
  `Status: to check`, as for `in progress`. No Go code changed: status is
  read as free text.
- `docs/concepts.md` (the loop's step 5, glossary "Blocked" and a new
  "Checkpoint"), `docs/getting-started.md` (implement) and
  `docs/coordinating.md` (a later session picks up the checkpoints, with
  the owner or unattended; an answer changing the change still needs
  `/grove-propose`).
- `dev/walkthrough.md`: new step "A checkpoint, then its answer"; the
  "Implement" step no longer expects "tells me how to try it".
- Would have changed `docs/guides/reconcile-a-maintained-spec.md`, which
  no `Files:` line names: it cites "telling the owner how to try each
  task" as an example of past instruction text, which still reads true as
  history. Left as it is.
- `just check`: passed (4031 Go lines, 1707 test lines, 8426 words).

## Scenarios

- implement "A screen to look at, unattended": internal/cli/walkthrough_test.go TestImplementInstructions checks the unattended text; an agent doing it is the walkthrough's "A checkpoint, then its answer", not run here: unverified
- implement "An open question, with the owner": TestImplementInstructions checks the owner text; agent behavior not run here: unverified
- implement "No checkpoint": TestImplementInstructions checks "Do not stop between tasks." and that the per-task stop is gone; agent behavior not run here: unverified
- implement "A task that builds on a checkpoint": TestImplementInstructions checks the text; agent behavior not run here: unverified
- implement "Confirmed in a short session": TestArchive shows `to check` is refused until done; the session is the walkthrough step, not run here: unverified
- implement "An answer in the words that start an unattended session": TestImplementInstructions checks the shared answer text; agent behavior not run here: unverified
- implement "An answer that changes the change, unattended": TestImplementInstructions checks "naming the revision the change needs"; agent behavior not run here: unverified
- implement "Feedback changes the change": TestImplementInstructions checks the owner's propose text, passed; agent behavior not run here: unverified
- implement "The approach rests on something false": TestImplementInstructions checks "stop that task and build nothing more of it", passed; agent behavior not run here: unverified
- implement "Unattended": TestImplementInstructions checks the unattended text, passed; agent behavior not run here: unverified

## Reconciliation

- implement: applied to grove/specs/implement/spec.md
