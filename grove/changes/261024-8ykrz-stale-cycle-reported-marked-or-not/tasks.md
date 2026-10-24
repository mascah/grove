## T-01 — Check reports a stale cycle whether or not it is marked

Result: `grove check` reports a stale cycle that is not marked `unattended`
as the scenarios "A member of a cycle" in `grove/specs/propose/spec.md`
and `grove/specs/implement/spec.md` say, in the wording the proposal's
decisions 1 and 2 give.

Files: `internal/cli/check.go`, `internal/cli/walkthrough_test.go`,
`docs/planning.md`, `dev/walkthrough.md`

Delegated: test layout and private helper names

Interfaces: `grove check` reports a new problem, `planned: <change-id>
changed on <target> since this plan; re-plan it with /grove-plan
<cycle-id>`, on a cycle not marked `unattended` (breaking, decision 3).

Constraints: a marked cycle's messages, and what `grove run` builds, stay
as they are; a member of a cycle not done is still not reported against
its `base`.

Acceptance:
- A test in `TestMovedFiles` reproduces the bug and fails before the fix:
  a member of an unmarked cycle whose named file changed on main after the
  plan gives check's unmarked line and exit status 1; the marked case
  still gives today's line.
- An unmarked cycle whose `planned` is no commit is reported as a marked
  one is.
- `docs/planning.md` describes and quotes the stale report for any cycle,
  not only under "Mark it unattended" (rule 8).
- `dev/walkthrough.md`'s Planning cycles step 5 touches a file budgets
  names, shows the 0.4 cycle's unmarked line, and re-plans it.

## T-02 — Implement names a member's files moved since the plan

Result: `grove instructions implement <change-id>` behaves as
[specs/implement/spec.md](specs/implement/spec.md) "Tell the owner of a
file moved since the plan" says, with the line of the proposal's
decision 8.

Files: `internal/cli/instructions.go`, `internal/cli/check.go`,
`internal/instruct/requires/implement-owner.md`,
`internal/instruct/requires/implement-unattended.md`,
`internal/cli/walkthrough_test.go`, `docs/coordinating.md`,
`dev/walkthrough.md`

Delegated: test layout and private helper names

Interfaces: a member's implement instructions carry ``- Changed on
<target> since cycle <cycle-id>'s plan <sha7>: <task>'s `<file>`.``;
check, `grove run` and implement count the commits since a plan with one
function (decision 13).

Constraints: `grove instructions implement <cycle-id>` prints nothing
new; the requires text grows by no words (decision 9).

Acceptance:
- `TestMovedFiles` checks that a member's moved file is named against the
  plan, with and without `--unattended`, and still not against `base`.
- `TestFollows` checks that the merge archiving the followed cycle's
  change, which edits a file the member names, leaves that member's
  instructions naming nothing; a re-planned cycle names nothing either.
- The implement requires text, with the owner and unattended, applies its
  rule for a changed file to that line as to the line against `base`.
- `docs/coordinating.md` says a task of a change in a cycle, naming a file
  that changed since the cycle's plan, is named as one changed since
  `base` is (rule 8).
- `dev/walkthrough.md`'s Planning cycles step 5 shows
  `grove instructions implement <budgets-id>` naming T-01's file since the
  cycle's plan before the re-plan.
