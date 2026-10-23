# Implementation: Plan's verdict, runs without a remote, streamed output

## Tasks

### T-01 — Plan ends with a verdict

Status: done

- `internal/instruct/requires/plan.md`: the finish has the agent say
  whether `grove run` would build the cycle once marked, and what removes
  each obstacle (a member needing a person, one waiting outside the cycle,
  something the runner's checkout must have), or how to add the mark.
- `docs/planning.md` shows the verdict for a cycle with a member that
  needs a person and for one with nothing in the way; `dev/walkthrough.md`
  planning steps 1 and 3 run them.
- `TestCycles` checks `grove instructions plan` prints the requirement.
  `just check` passes.
- The ledger fixture run of `/grove-plan` in Claude Code was not done: the
  owner, 2026-10-23, "I'll skip the fixture check, it's getting too time
  consuming to setup".

### T-02 — Run without a remote

Status: done

- `internal/cli/run.go`: `grove run` fetches and checks the fast-forward
  only when `git remote` lists `origin`; otherwise the local target's
  commit is the fetched target.
- `docs/running.md` has "Without a remote", saying the owner merges the
  attempt's branch; "The remote" and the ordinary run's step 1 say
  `origin` is optional. `docs/reference.md`'s `grove run` row no longer
  requires `origin`. `dev/walkthrough.md` running step 6 runs it.
- `TestRun` ends with the checkout's `origin` removed: a local commit of
  its own is no refusal, the record holds the local target's commit, the
  project is paused, and merging the attempt's branch locally ends the
  pause. The `origin` cases above it are unchanged. Without the code
  change it fails with Git's "'origin' does not appear to be a git
  repository". `just check` passes.

### T-03 — Show the session as it runs

Status: not started

## Reconciliation
