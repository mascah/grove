# A cycle can depend on cycles

## Problem

Cycles planned together go stale as each earlier one lands. Nearly every
change's tasks name the same busy files, such as `docs/reference.md`,
`dev/walkthrough.md` and `internal/cli/walkthrough_test.go`. At ab885af,
`changedSincePlan` (`internal/cli/check.go`) counts any difference in a
member's folder or named files between the cycle's `planned` commit and
the target, whatever commit made it.

On 2026-10-23 the four cycles of that day's review were planned in one
merge (65173bb, #56). Once cycle A's two changes merged (f749c46, #57;
aa245f1, #58), `grove run 261023-ynzcd` refused cycle B: "cycle
261023-ynzcd's plan is stale, since 261023-g0a3h and ... changed on main
after it". None of B's change folders had moved; only files B's tasks name
had, in the merges that built and archived cycle A. B was re-planned by
hand (ab885af, #60). Cycles C and D, planned the same day, will go stale
the same way as B lands.

Nothing in a cycle file says one cycle follows another, so Grove cannot
tell such an expected merge from a change the plan did not see. At
ab885af, `doc.Load` (`internal/doc/doc.go`) reads a cycle's `title`,
`created`, `changes`, `planned`, `goal`, `needs_human` and `unattended`,
and reports any other field.

## Intended outcome

A cycle file may name, in `depends_on`, the cycles it follows, as a
change's `.grove.yaml` names the changes it waits on. `grove run` builds
no cycle while a cycle it depends on is not done, given or marked, and
says which. A change to a member's folder, or to a file its tasks name,
does not make the plan stale when it arrived in a commit on the target's
first-parent history that archived a change of a cycle this one depends
on, directly or through others; any other change still does. So a cycle
planned behind another stays buildable as that one lands, by squash
merges, merge commits or stacked pull requests.

`grove check` reports a `depends_on` naming no cycle, and cycles that
depend on each other in a loop. `grove changes` says a change's cycle
waits on a cycle that is not done. `/grove-plan` asks, when other cycles
are not done, whether the new cycle follows any of them, and writes those
the owner names. `/grove-implement <cycle-id>` ignores the field.

The rules are in [specs/cycles/spec.md](specs/cycles/spec.md) and
[specs/runs/spec.md](specs/runs/spec.md).

## Scope

### In scope

- The optional `depends_on` field of a cycle file: reading it, the
  template line, and `grove check`'s reports on it.
- `grove run`'s selection of a marked cycle and its refusal of a given
  one.
- The stale-plan check shared by `grove check` and `grove run`.
- `grove changes`' line for a change whose cycle waits on a cycle, and
  the instructions' project index naming the cycles a cycle follows.
- The plan instructions: the question for a new cycle, and the finish
  naming a cycle it waits on as an obstacle, as the maintained "Say
  whether the cycle can run unattended" already asks for each obstacle.
- `docs/planning.md`, `docs/running.md`, `docs/reference.md` and
  `dev/walkthrough.md`.

### Out of scope

- `/grove-implement <cycle-id>`: it builds the cycle as it would without
  `depends_on`. Each member's own `depends_on` still applies, and the
  owner orders attended work.
- Waiting through another cycle when building: `grove run` waits only on
  the cycles a cycle's `depends_on` names. Only the stale check follows
  them further.
- `grove check` reporting a marked cycle that waits on a cycle not done:
  it is waiting, as a change waiting on its `depends_on` is.
- `/grove-plan <cycle-id>` asks no new question: revising, the agent
  changes `depends_on`, like every other field, only as the owner asks.
- A member revised, or a file its tasks name changed, in a commit that
  archived no change of a cycle this one depends on: it still makes the
  plan stale, as change 261023-bgtcb leaves the check.
- Change 261023-bgtcb's exception, a member changed in the same commit as
  an edit to the cycle file: this change starts from its text and keeps
  it.

This change needs the constitution amendment of PR #61, open at this
writing, merged first: at ab885af, rule 10 has `run` refuse any cycle
whose members or named files differ from what its plan checked.

## Decisions

The owner chose the field's name ("We use depends_on for changes so we
should maybe keep it consistent. Let's do it before C and D") and agreed
to every answer below in one round ("Yes"), 2026-10-23.

1. **An optional `depends_on: [<cycle-id>, ...]` in a cycle file, read as
   a change's is;** the cycle template gets a commented line for it.
   Reason: one word for "finish these first" across changes and cycles,
   and `doc.Load` already reads cycles with required and optional field
   lists.
2. **`grove run`, given a cycle or selecting a marked one, builds no cycle
   while a cycle it depends on is not done, saying which in one
   sentence.** Reason: a cycle built before the one it follows would build
   on code that is not there yet; rule 14 keeps the message to one
   sentence.
3. **A change to a member's folder, or to a file its tasks name, does not
   count when it arrived in a commit on the target's first-parent history
   that archived a change of a cycle this one depends on, directly or
   through others; anything else still counts.** Squash merges, merge
   commits and stacked pull requests, one commit per archived change, all
   work. Reason: the plan expected exactly those merges; reading first
   parents treats a merge commit as one commit, as a squash is; following
   cycles through others covers a chain of cycles planned together.
4. **`grove check` reports a cycle's `depends_on` naming no cycle, and
   cycles that depend on each other in a loop,** as it does for changes'
   `depends_on`. Reason: either would leave a cycle that never runs.
5. **`grove changes` adds "waits on cycle X" to a change whose cycle waits
   on a cycle not done; `/grove-implement <cycle-id>` ignores the field.**
   Reason: the owner sees why a run skips the change; attended work the
   owner orders themselves, and members' own `depends_on` still apply.
6. **`/grove-plan` asks, when other cycles are not done, whether this
   cycle follows any of them, and writes those the owner names.** Reason:
   the agent infers no order; only the owner's answer goes in the file.
7. **The specs: cycles gets the field and its check reports, and the
   staleness rule changes; runs' selection and given-cycle refusal
   change.** Docs: `docs/planning.md`, `docs/running.md`,
   `docs/reference.md` and `dev/walkthrough.md`, with tests in
   `internal/cli/walkthrough_test.go`. Reason: the documentation is the
   specification (rule 8).

The owner answered the proposal's two open questions on PR #62 ("yes do
both", 2026-10-23):

8. **A cycle `depends_on` naming no cycle is reported by `grove check`,
   and `grove run` treats the cycle as waiting,** as a change's
   `depends_on` naming no change keeps it waiting at ab885af
   (`Change.WaitsOn`, `internal/doc/doc.go`). Reason: one rule for both
   kinds of `depends_on`, and a typo stops the run instead of skipping
   the wait.
9. **The instructions' project index shows a cycle's `depends_on`,** as
   "follows <cycle-id>" under the cycle. Reason: `/grove-plan` and the
   other skills see the order the owner set without opening each cycle
   file.

## Open questions

None must be answered before review. None blocks a task.

## Related documents

- Requirements: [specs/cycles/spec.md](specs/cycles/spec.md),
  [specs/runs/spec.md](specs/runs/spec.md)
- Tasks: [tasks.md](tasks.md)
