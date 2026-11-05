# Implementation: Plan brings a stale member up to date

<!-- Progress on change 261105-gmpa4, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Plan brings a stale member up to date

Status: done

- `internal/instruct/requires/plan.md`: the first bullet no longer says
  the cycle file is the only thing written; a new bullet says that for a
  new cycle, before writing its file, the agent runs `check` and, for each
  member it reports stale on a requirement or for a named file changed
  since its base, tells the owner this edits approved work, brings it up
  to date on this branch as `instructions propose <change-id>` says, then
  goes on with the plan. It points to propose rather than restating it
  (+59 words, 8107 of 9000).
- `internal/cli/walkthrough_test.go`: `TestCycles` checks the new text and
  that "the only thing you write" is gone.
- `docs/planning.md` "Plan a cycle": says so, and that one merge then
  brings the revision and the plan.
- `dev/walkthrough.md` "Planning cycles" step 1 first moves
  `ledger/report.py` on main so budgets is stale on its base, then expects
  the session to bring T-01 up to date and set `base`, `git status` to show
  budgets' files beside the cycle file, and `grove check` clean. Step 2's
  `grove review` now expects budgets' revised files beside the cycle file,
  since the plan branch now carries them.
- `just check`: passed (4031 Go lines, 1701 test lines, 8107 words).

## Scenarios

- cycles "Two approved changes": internal/cli/walkthrough_test.go TestCycles (plan instructions, and the cycle file merged), passed; an agent writing the file is the walkthrough's "Plan a cycle", not run here: unverified
- cycles "A draft change named": internal/cli/walkthrough_test.go TestCycles checks the plan instructions; an agent refusing the draft is not run here: unverified
- cycles "A member stale on its base": internal/cli/walkthrough_test.go TestCycles checks the plan instructions say it; an agent doing it is the walkthrough's "Plan a cycle" step 1, not run here: unverified

## Reconciliation

- cycles: applied to grove/specs/cycles/spec.md
