# Implementation: Plan's verdict, runs without a remote, streamed output

## Tasks

### T-01 — Plan ends with a verdict

Status: in progress

- `internal/instruct/requires/plan.md`: the finish has the agent say
  whether `grove run` would build the cycle once marked, and what removes
  each obstacle (a member needing a person, one waiting outside the cycle,
  something the runner's checkout must have), or how to add the mark.
- `docs/planning.md` shows the verdict for a cycle with a member that
  needs a person and for one with nothing in the way; `dev/walkthrough.md`
  planning steps 1 and 3 run them.
- `TestCycles` checks `grove instructions plan` prints the requirement.
  `just check` passes.
- Left: the ledger fixture run of `/grove-plan` in Claude Code.

### T-02 — Run without a remote

Status: not started

### T-03 — Show the session as it runs

Status: not started

## Reconciliation
