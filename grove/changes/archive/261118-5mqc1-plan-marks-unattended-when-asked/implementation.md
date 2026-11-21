# Implementation: Plan marks a cycle unattended when the owner asks

<!-- Progress on change 261118-5mqc1, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Plan writes the mark when the owner asks

Status: done

`internal/instruct/requires/plan.md`: the "Never write or remove
`unattended`" line is replaced by one that has the agent write the mark only
where the owner's description or the conversation asks for an unattended
cycle, even with an obstacle, never offering it, and, revising, remove it only
where the owner asks and keep it otherwise. The finish now says, where the
cycle file holds `unattended: true`, that merging it hands the cycle to a
`grove claim` given no cycle, or else that `unattended: true`, added in a merge
of its own, does, still after `/grove-implement <cycle-id>` and whether
`grove claim <cycle-id>` would claim it.

`internal/cli/walkthrough_test.go` `TestCycles` checks the new lines, that
"Never write or remove `unattended`" is gone, the finish's order and both
branches of its last sentence; its revise step still checks that
`grove instructions propose <member-id>` says "never write or remove
`unattended`" (`propose.md` unchanged). `TestClaim` already claims cycle `l`,
committed with `unattended: true` in the commit that plans it, with a
`grove claim` given no cycle.

`just check` passes; instruction words 8199 → 8227 (+28, limit 60).

Check: 261118-emse7 and 261118-40s1a changed `plan.md` and the test after
the change's dc8d335 reading; the cycle's plan brought `base` to 7a33629, and
the text the acceptance points at was still there, so the edit is on the file
as main holds it.

### T-02 — Docs and walkthrough

Status: done

`docs/planning.md`, "Plan a cycle": asked for an unattended cycle, the agent
writes `unattended: true`, even with an obstacle, never offering it, and its
verdict says merging the file hands the cycle to a `grove claim` given no
cycle; the verdict paragraph gives both cases. "Mark it unattended": the
separate merge is for a cycle you did not ask plan to mark; only plan writes
or removes the line, only when asked, and re-planning keeps it otherwise.

`dev/walkthrough.md`, "Running unattended" step 4: the budgets cycle is
planned with `/grove-plan <budgets-id>`, saying it should run unattended; the
file holds `unattended: true` and one merge and push replaces the mark.
Planning steps 1 and 3 and coordinating's 0.5 cycle, where the owner asks
nothing, still show no `unattended` and a mark in a merge of its own.
`TestClaim` follows that step: cycle `l` is committed marked in the commit
that plans it and claimed by `grove claim` given no cycle.

Check: the docs' example request is "/grove-plan an unattended cycle for
delete and budgets"; the change names no wording. The walkthrough's
session step was not run with a real agent here; only the instructions it
prints were read (`bin/grove instructions plan "0.4 release"`).

## Scenarios

The plan scenarios are an agent's behavior; the tests check the instructions
that produce it, and no live `/grove-plan` session was run.

- cycles "Nothing in the way": internal/cli/walkthrough_test.go TestCycles (finish order and "or else say `unattended: true`, added in a merge of its own, does."), passed
- cycles "Marked as the owner asked": TestCycles (finish's "where the cycle file holds `unattended: true`, say merging it hands the cycle to a `grove claim` given no cycle"), passed
- cycles "A member waits on a change outside the cycle": unchanged instruction text; TestCycles checks the obstacle lines, passed
- cycles "The cycle waits on a cycle": TestCycles ("`grove claim` will not claim the cycle until it is"), passed
- cycles "The runner needs something set up": unverified; agent behavior, instruction text unchanged and not asserted by a test
- cycles "Planned but not marked": TestClaim (nothing to claim after cycle k is planned unmarked), passed
- cycles "The owner asks for an unattended cycle": TestCycles (instruction line) and TestClaim (cycle l, marked in the commit that plans it, claimed with no cycle given), passed; the agent writing the file is unverified without a live session
- cycles "The owner asks nothing about it": TestCycles ("never offer it"), passed; agent behavior unverified
- cycles "An obstacle does not stop the mark": TestCycles ("even with an obstacle below"), passed
- cycles "Re-plan keeps the mark": TestCycles ("keep it otherwise"), passed
- cycles "The owner asks to remove the mark": TestCycles ("remove it only where the owner asks"), passed

## Reconciliation

- cycles: applied to specs/cycles/spec.md
