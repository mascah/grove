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

## Scenarios

## Reconciliation
