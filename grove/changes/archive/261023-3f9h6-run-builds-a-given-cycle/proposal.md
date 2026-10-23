# grove run builds a cycle it is given

## Problem

The owner's usage feedback of 2026-10-23, from session
`8617ec36-2e5d-47ba-828b-7a0d218200a8`. After the owner merged a plan, the
agent closed with: "The cycle isn't marked unattended, so grove run won't
pick it up. To hand it to grove run, add unattended: true to the cycle file
in a separate PR. Otherwise, /grove-apply 261023-58dcm builds it in a
session you start." The owner: "It seems like we don't have a way to build
a cycle interactively. Right now they can only be run unattended."

- `/grove-implement <cycle-id>` already builds a cycle with the owner
  (`docs/coordinating.md` at 2f8f44d, "Applying a cycle"), but the plan's
  finish in `internal/instruct/requires/plan.md` at 2f8f44d names only
  `grove run` and the `unattended` mark, so the owner was not told.
- `grove run` takes no cycle (`internal/cli/cli.go` at 2f8f44d refuses any
  argument) and selects only a cycle marked `unattended`
  (`eligible` in `internal/cli/run.go` at 2f8f44d). Handing a planned
  cycle to `grove run` once takes a second pull request that adds the
  mark, which also hands it to every later timer-driven run.

`CONSTITUTION.md` rules 2, 4 and 10, as amended in #32, now let `run`
build a cycle it is given.

## Intended outcome

- `grove run <cycle-id>` builds that cycle whether or not it is marked
  `unattended`, under every other gate that applies to a marked cycle, and
  says in one sentence why when it cannot.
- Plain `grove run`, as a timer runs it, is unchanged.
- The plan's finish names `/grove-implement <cycle-id>` first, then
  `grove run <cycle-id>`, then the mark.

## Scope

### In scope

- The `grove run` cycle argument, its refusals, and its attempt record and
  pause.
- The plan instructions' finish.
- `docs/running.md`, `docs/planning.md`, `docs/reference.md` and
  `dev/walkthrough.md` steps for each.

### Out of scope

- An interactive `grove run`, one launching `claude` without `-p`;
  `/grove-implement <cycle-id>` is that.
- How `grove run` displays the session: change 261023-yz5bk.
- Harnesses other than Claude Code for `run`.
- `grove check`'s reports for marked cycles, which stay as they are.

## Decisions

- **Naming the cycle is the consent.** The owner, 2026-10-23, offered a
  flag that ignores the mark, flipping `grove run` to interactive with
  `--apply` for marked cycles, or `grove run <cycle-id>`: "Yes" to the last,
  and merged it into the constitution as #32. Reason: a flag that turns
  off one of Grove's own gates is a switch on a gate (rule 10); an
  interactive `grove run` is `/grove-implement`.
- **Usage is `grove run [--retry] [--timeout <duration>] [<cycle-id>]`.**
  The owner, 2026-10-23: "yes to all".
- **A given cycle that cannot be built is refused in one sentence naming
  why**, with nothing started: an unknown ID, a done cycle, a member under
  `needs_human`, a stale plan, or a member waiting on a change neither in
  the cycle nor archived as finished. The owner, 2026-10-23: "yes to all".
  Reason: `grove check` reports these only for a marked cycle, so "Nothing
  to apply" would leave the operator without a reason.
- **A given cycle respects the pause and needs `--retry`; its attempt is
  recorded and pauses the project like any other.** The owner, 2026-10-23:
  "yes to all".
- **The plan's finish names `/grove-implement <cycle-id>`, then
  `grove run <cycle-id>`, then `unattended: true` merged on its own for a
  `grove run` given no cycle.** With a member that needs a person, it says
  `grove run <cycle-id>` would refuse too and offers `/grove-implement`.
  The owner, 2026-10-23: "plan skill's closing message names
  /grove-implement <cycle-id> first", and "yes to all".
- **`docs/running.md` gets a section "Run a cycle you name", and the
  walkthrough a step for it.** The owner, 2026-10-23: "yes to all".
- **No `depends_on`.** The owner, 2026-10-23: "yes to all". This change and
  261023-yz5bk both edit `internal/cli/run.go` and `docs/running.md`.

## Open questions

None.

## Related documents

- Requirements: [specs/runs/spec.md](specs/runs/spec.md) (maintained spec:
  `grove/specs/runs/spec.md`), [specs/cycles/spec.md](specs/cycles/spec.md)
  (maintained spec: `grove/specs/cycles/spec.md`)
- Tasks: [tasks.md](tasks.md)
