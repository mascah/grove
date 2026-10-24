# Implementation: A cycle can depend on cycles

<!-- Progress on change 261023-2ej3w, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Read a cycle's depends_on and check it

Status: done

Built on main@9308040, as the owner agreed, though every named file moved
since base ab885af (#63–#67). `doc.Load` reads a cycle's optional
`depends_on` into `Meta.DependsOn`; `Cycle.WaitsOn` holds what it names that
is not a cycle done; `Collection.follows`, beside `dependencies`, reports a
`depends_on` naming no cycle and a loop. `grove changes` adds `, waits on
cycle <id>`; the project index adds `follows <id>` under a cycle.
`TestCycles` covers each acceptance line: `go test ./...` passed.

### T-02 — Run waits on the cycles a cycle depends on

Status: done

`runnable` refuses a cycle whose `WaitsOn` is not empty ("cycle K waits on
cycle L, which is not done"), so `eligible` passes it over.
`changedSincePlan` also leaves out a commit on the first-parent history
that touches the archived folder of a change listed by a cycle `cy` depends
on, directly or through others; bgtcb's cycle-file exception stays.
`TestFollows` covers each acceptance line, including that implement lists
the member of a waiting cycle. With the new exception disabled, four of its
checks fail. `TestReplan`, `TestMovedFiles`, `TestCycles` and `TestRun` pass
unchanged; `just check` passed.

### T-03 — Plan asks, and the docs say it

Status: done

The cycle template has a commented `# depends_on:` line; `requires/plan.md`
asks, for a new cycle, whether it follows a cycle not done, and its finish
names that obstacle. `schemas/default/skills/plan.md` needed no edit: the
rule says it. `docs/planning.md` gains "Follow another cycle", the stale
rule's exception and the two check problems; `docs/running.md` and
`docs/reference.md` name the field and the wait; `dev/walkthrough.md`
gains planning step 7 and a refusal in running step 6. `TestCycles`
checks the template line, the question and the obstacle; `just check`
passed. The owner agreed to add the words though `internal/instruct` was
already over rule 12's 7,000 (7,808 at main, 7,876 now; `just check` does
not count words until change 261023-ejarr). The hand-run `/grove-plan`
trial is not run: the owner agreed to mark the task done with its
scenarios unverified, to run before merging.

## Scenarios

cycles:

- "The cycle it follows is not done": `internal/cli/walkthrough_test.go` TestFollows, `grove run K` refused, passed
- "The cycle it follows is done": TestFollows "L's last change archived in a squash", `eligible` selects K, passed
- "Built with the owner": TestFollows, `instructions implement K` lists B, passed
- "The owner names one": unverified, the hand-run `/grove-plan` trial on the fixture (T-03) is not run
- "The owner names none": unverified, as above
- "No other cycle": unverified, as above; TestCycles checks only that the question is printed
- "Revised and re-planned in one squash merge": TestReplan, passed
- "Revised and re-planned in one merge commit": TestReplan, passed
- "Marked after a member changed": TestReplan, passed
- "Revision and re-plan as two commits on the target": TestReplan, passed
- "The cycle it follows lands in a squash merge": TestFollows "L's last change archived in a squash", passed
- "Stacked pull requests": TestFollows, C landed in a merge commit and D in a squash, passed
- "Through another cycle": TestFollows "M's change archived in a merge commit", passed
- "Any other commit still counts": TestFollows "a commit archiving nothing" and "an archive of a cycle K does not follow", passed
- "Nothing in the way": TestCycles checks the finish's order, passed; the agent's words unverified, a person judges them
- "A member needs a person": TestCycles checks the finish's text, passed; the agent's words unverified
- "A member waits on a change outside the cycle": unverified, how the agent says it is a person's to judge
- "The cycle waits on a cycle": TestCycles checks the finish's obstacle, passed; the agent's words unverified, the trial is not run
- "The runner needs something set up": unverified, how the agent says it is a person's to judge
- "One change in two cycles": TestCycles, passed
- "A dependency on no cycle": TestCycles, passed
- "Cycles that wait on each other": TestCycles, passed
- "One marked, one not": TestCycles, passed
- "A cycle that waits on a cycle": TestCycles, passed
- "The index names the cycle followed": TestCycles, through `instructions explore`, which prints the same index as plan, passed

runs:

- "Local target has commits of its own": TestRun, passed
- "A marked cycle of two changes, one waiting on the other": TestRun, passed
- "Claude Code listed as commands": TestRun, passed
- "Two marked cycles": TestRun, passed
- "The older marked cycle waits on the newer": TestFollows "M's change archived in a merge commit", `eligible` selects L, passed
- "A given cycle not marked": TestRun, passed
- "A given cycle while paused": unverified, no test gives a cycle while paused; TestRun tests the pause with none given
- "A project that does not list Claude Code": TestRun, passed
- "A checkout with no remote": TestRun, passed
- "The branch merged locally": TestRun, passed
- "A member needs a person": TestRun, passed
- "A stale plan": unverified, no test gives `grove run` a stale cycle; TestReplan and TestFollows test the same `runnable` through `eligible`
- "A cycle it follows not done": TestFollows, passed
- "No such cycle": TestRun, passed

## Reconciliation
