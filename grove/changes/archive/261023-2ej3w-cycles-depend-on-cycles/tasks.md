## T-01 — Read a cycle's depends_on and check it

Result: a cycle file may hold `depends_on`, and `grove check`,
`grove changes` and the instructions' project index treat it as
[specs/cycles/spec.md](specs/cycles/spec.md) "Report a malformed cycle"
and "Show cycles" say.

Files: `internal/doc/doc.go`, `internal/cli/changes.go`,
`internal/instruct/instruct.go`, `internal/cli/walkthrough_test.go`

Interfaces: `doc.Load` reads a cycle's optional `depends_on` into
`Meta.DependsOn`, as for a change; `doc.Cycle` gains `WaitsOn`, the
cycles its `depends_on` names that the collection does not hold done; a
check beside `dependencies` reports, for each cycle not done, a
`depends_on` naming no cycle and one that waits on it in turn.

Acceptance:
- In `internal/cli/walkthrough_test.go` (`TestCycles`), on the fixture: a
  cycle with `depends_on` naming another cycle checks clean; one naming
  261099-zzzzz is reported on its `depends_on` line; two cycles naming
  each other are each reported; the unknown-field message lists
  `depends_on` among the fields a cycle file takes.
- `grove changes` adds `waits on cycle <id>` to a change whose cycle
  waits on a cycle not done, and nothing once that cycle is done.
- In `internal/cli/walkthrough_test.go` (`TestCycles`): the project index
  any skill's instructions print lists `follows <cycle-id>` under a cycle
  not done whose `depends_on` names that cycle.

## T-02 — Run waits on the cycles a cycle depends on

Result: `grove run` and the stale-plan check behave as
[specs/runs/spec.md](specs/runs/spec.md) "Start one attempt" and "Refuse
a given cycle that cannot be built", and
[specs/cycles/spec.md](specs/cycles/spec.md) "A cycle can depend on
cycles" and "Count a member changed with its plan", say.

Files: `internal/cli/run.go`, `internal/cli/check.go`,
`internal/cli/walkthrough_test.go`

Interfaces: `runnable` refuses a cycle whose `WaitsOn` is not empty, in
one sentence naming the cycle it waits on, so `eligible` passes it over;
`changedSincePlan` also leaves out a commit on the target's first-parent
history since `planned` that archives a change listed by a cycle `cy`
depends on, directly or through other cycles, as `at` holds them; its
messages keep their words.

Constraints: Git runs through `gitx.Run`, as the check does today; the
exception change 261023-bgtcb adds, a commit that also edits the cycle
file, stays; within the rule 12 budget.

Acceptance:
- In `internal/cli/walkthrough_test.go` (`TestRun`), on the fixture:
  `grove run K` for a cycle K whose `depends_on` names a cycle not done
  launches nothing and says so in one sentence; with K the older of two
  marked cycles, `grove run` selects the one it depends on.
- In `internal/cli/walkthrough_test.go`, on the fixture: with K
  depending on L, a squash commit archiving L's change and editing a file
  K's member names leaves `grove check` clean and `grove run` builds K;
  the same as a `--no-ff` merge does too; two commits each archiving one
  of L's changes do too; a commit archiving a change of a cycle L depends
  on does too; a commit archiving nothing of L that edits that file makes
  K stale in both, as does an archive of a cycle K does not depend on.
- `/grove-implement <cycle-id>` instructions for a cycle waiting on a
  cycle not done list its members as without `depends_on`.
- The existing stale-plan tests, and change 261023-bgtcb's, pass
  unchanged.

## T-03 — Plan asks, and the docs say it

Result: the plan instructions ask and finish as
[specs/cycles/spec.md](specs/cycles/spec.md) "Ask which cycles a new cycle
follows" and "Say whether the cycle can run unattended" say, and the
documentation describes the field, the wait and the stale rule (rule 8).

Files: `internal/instruct/requires/plan.md`,
`internal/instruct/schemas/default/skills/plan.md`,
`internal/instruct/metadata/cycle.yaml`, `docs/planning.md`,
`docs/running.md`, `docs/reference.md`, `dev/walkthrough.md`,
`internal/cli/walkthrough_test.go`

Interfaces: the cycle template holds a commented `# depends_on:` line.

Constraints: within the rule 12 word budget for `internal/instruct`.

Acceptance:
- In `internal/cli/walkthrough_test.go` (`TestCycles`): `instructions
  plan "<description>"` prints the commented `depends_on` line, the
  question for when other cycles are not done, and the finish's obstacle
  for a cycle that waits on a cycle.
- On the ledger fixture, in Claude Code: with the show-saved cycle not
  done, `/grove-plan 0.5` asks whether the new cycle follows it, writes
  `depends_on: [<show-saved-cycle-id>]` when told so, and ends saying
  `grove run` would not build it until that cycle is done; with no other
  cycle not done it asks nothing about cycles.
- `docs/planning.md` shows `depends_on` in the cycle file, says plan asks
  it, names the stale rule's exception for the cycles a cycle depends on,
  and lists the two new check problems; `docs/running.md` names the wait
  in an ordinary run's selection and among "Run a cycle you name"'s
  refusals; `docs/reference.md` names `depends_on` among a cycle file's
  fields and the `grove run` and `grove changes` rows say the wait.
- `dev/walkthrough.md`: in planning cycles, a step plans a second cycle
  that follows another, checks the two new problems, and lands the first
  so the second stays fresh while a commit archiving nothing of it makes
  it stale; in running unattended, `grove run <cycle-id>` refuses a cycle
  that waits on a cycle not done.
