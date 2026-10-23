## T-01 — Plan ends with a verdict

Result: `/grove-plan` finishes as [specs/cycles/spec.md](specs/cycles/spec.md)
"Say whether the cycle can run unattended" says, and the documentation
specifies it.

Constraints: the plan instructions stay data under `internal/instruct`
(rule 10); no new `grove check` report.

Acceptance:
- `docs/planning.md` shows the verdict for a cycle with nothing in the way
  and for one with a member that needs a person, with what removes it;
  `dev/walkthrough.md` runs those steps (rule 8).
- On the ledger fixture, in Claude Code: `/grove-plan` on two approved
  changes, one of which the owner says needs a person, ends by saying
  `grove run` would skip the cycle if marked and what removes that.
- `internal/cli/walkthrough_test.go` checks that `grove instructions plan`
  prints the requirement; `just check` passes.

## T-02 — Run without a remote

Result: `grove run` behaves as [specs/runs/spec.md](specs/runs/spec.md)
"Start one attempt" says in a checkout with no remote named `origin`, and
the documentation specifies it.

Constraints: within the rule 12 budget; the behavior with `origin` is
unchanged.

Acceptance:
- `docs/running.md` has a step for running in a checkout with no remote,
  saying the owner merges the attempt's branch into the target themselves;
  `docs/reference.md`'s `grove run` row no longer requires `origin`;
  `dev/walkthrough.md` runs the step.
- `internal/cli/walkthrough_test.go` covers a checkout with no remote
  starting an attempt with the local target's commit in its record, and
  the pause ending when the branch is merged locally, with the `origin`
  cases unchanged; `just check` passes.

## T-03 — Show the session as it runs

Result: `grove run` behaves as [specs/runs/spec.md](specs/runs/spec.md)
"Show the session as it runs" says, and the documentation specifies it.

Constraints: the standard library only (rule 16); messages are one
sentence (rule 14).

Acceptance:
- `docs/running.md` describes the lines printed while the session runs and
  the JSON Lines log, and no longer says a stopped session's output is
  usually empty; `docs/guides/run-apply-on-a-timer.md` says what the
  launchd log receives; `dev/walkthrough.md` runs the step.
- On the ledger fixture with Claude Code: `grove run` on a marked cycle
  prints lines while the session runs, and the attempt's log holds the
  session's events as JSON Lines.
- `internal/cli/walkthrough_test.go`, with a stand-in `claude` that emits
  stream-json, checks the lines on standard output and the log's content;
  `just check` passes.
