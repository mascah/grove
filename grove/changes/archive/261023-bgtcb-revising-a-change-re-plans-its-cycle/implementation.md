# Implementation: Revising a change re-plans its cycle

<!-- Progress on change 261023-bgtcb, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Propose re-plans the cycle of a revised member

Status: done

- `internal/cli/instructions.go`: `instructionsPropose`, revising a
  change, finds the cycle not done in the working tree that lists it and
  passes it, with its file, to `instruct.Revise`.
- `internal/instruct/instruct.go`: `Revise` takes that cycle or none;
  given one, the task line says it re-plans the cycle, the cycle document
  is written from `metadata/cycle.yaml` with `planned` the instructions'
  target commit and `created` the cycle's, and the file is printed under
  "## The cycle", as `Replan` does.
- `internal/instruct/requires/propose.md`: revising a member, re-plan its
  cycle on this branch: ask whether the change needs a person as planning
  does, set `planned` from the template, change nothing else, never write
  or remove `unattended`, and tell the owner the revision and re-plan must
  reach the target as one commit.
- `internal/instruct/requires/plan.md`: the finish for a member that needs
  a person names revising it with `/grove-propose`, which re-plans the
  cycle on the same branch, and one merge.
- `internal/cli/walkthrough_test.go`: `TestCycles` checks that `instructions
  propose <member>` of a marked cycle prints the re-planning task line, the
  cycle document with `planned: <main's commit>`, the rule never to write
  or remove `unattended` and the cycle's file; that a change in no cycle
  gets no cycle; and that `instructions plan` names the shorter way.
- Trial in Claude Code 2.1.310 on a fresh ledger fixture with this build
  first on PATH, headless over three turns: an approved "Monthly budgets in
  report" in a cycle whose `needs_human` names it for judging how the
  report reads. `/grove-propose <budgets-id>` on a branch printed
  "re-planning cycle <id>", asked unprompted whether the `needs_human`
  entry should stay, and, asked to make the acceptance checkable by a test
  and told budgets no longer needs a person, revised the proposal, change
  spec and tasks, removed the entry, set `planned` to main's head, added no
  `unattended`, and told the owner the revision and re-plan must land as
  one commit. `grove check` was clean on the branch, and on main after a
  squash merge, the cycle not stale.

### T-02 — The stale check counts a member changed with its plan

Status: done

- `internal/cli/check.go`: `changedSincePlan`, which `stalePlans` and
  `runnable` in `internal/cli/run.go` share, reads the target's
  first-parent history since `planned` once per cycle with `git log
  --first-parent --diff-merges=first-parent --name-only`, so a merge commit
  counts as one change against its first parent, as a squash does; a
  member is stale when some commit there changed its folder or a named
  file without also editing the cycle file. Its messages are unchanged.
- `internal/cli/walkthrough_test.go`: `TestReplan`, on the fixture, with a
  marked cycle of one member: the revision and a re-plan setting `planned`
  to main's commit, squash-merged, leave `grove check` clean and
  `eligible`, which `grove run` and `grove run --retry` select by, picks
  the cycle; as two commits merged with `--no-ff`, the same; as two
  commits fast-forwarded onto main, `grove check` reports the cycle stale
  with its old message and `eligible` picks nothing; a member revised in a
  merge and the cycle marked in a later commit is reported stale. The
  existing stale-plan tests in `TestCycles`, `TestMovedFiles` and
  `TestRun` pass unchanged.
- The pause after an attempt is not replayed with a stub session here;
  `--retry` starts any attempt `eligible` or `runnable` allows, which the
  test checks after the squash merge.

### T-03 — Document the shorter way past a pause

Status: done

- `docs/running.md`, "The pause": answer the question with `/grove-propose
  <change-id>`, which re-plans the cycle on the same branch, merge, and run
  `grove run --retry`, with no separate re-plan merge; the revision and the
  re-plan reach main as one commit, squash or merge.
- `docs/planning.md`: the quoted finish for a member that needs a person
  names one merge; the definition of a stale plan names the exception, read
  on main's first-parent history; "Re-plan a cycle" says revising a member
  with `/grove-propose` re-plans its cycle, and that two commits by rebase
  or fast-forward leave it stale.
- `docs/coordinating.md`: where it names re-planning a cycle, revising a
  member re-plans it. `docs/reference.md`: the `grove instructions propose
  <change-id>` row says it re-plans a cycle the change is in.
- `dev/walkthrough.md`: planning step 1's finish names the shorter way;
  planning step 5 adds two commits fast-forwarded onto main leaving the
  cycle stale; running step 6 revises export with `/grove-propose
  <export-id>`, re-planning its cycle on the same branch, one squash merge,
  and `grove run --retry <export-cycle-id>` then starts a session.
- Recheck before finishing: T-01 to T-03's evidence holds against the code
  as it now is. `just check` passes (3767 Go lines, 1441 test lines);
  `grove check` prints "No problems in grove/.".

## Scenarios

- propose "A blocked member answered": unverified as a blocked member; the trial revised a member that needed a person, with the same instructions (`TestCycles` checks them for a marked cycle, leaving `unattended` as it is).
- propose "A member that no longer needs a person": the trial above, passed (entry removed, `planned` set, `grove check` clean after a squash merge); `TestReplan` shows `grove run` would then select the cycle.
- propose "A change in no cycle": `walkthrough_test.go` `TestCycles`, passed.
- cycles "Revised and re-planned in one squash merge": `TestReplan`, passed.
- cycles "Revised and re-planned in one merge commit": `TestReplan`, passed.
- cycles "Marked after a member changed": `TestReplan`, passed.
- cycles "Revision and re-plan as two commits on the target": `TestReplan`, passed (check reports it stale and `eligible` selects nothing; `grove run K` refuses through the same `runnable`).
- cycles "Nothing in the way", "A member waits on a change outside the cycle", "The runner needs something set up": unverified here, unchanged by this change; the finish text is printed, as `TestCycles` checks.
- cycles "A member needs a person": `TestCycles`, passed, for the finish naming the re-plan by `/grove-propose` and one merge; an agent's finish was not tried.
- runs "A blocked change revised": `TestReplan`, passed, for a squash merge of a revision and re-plan leaving the cycle eligible; the full pause and `--retry` with a stub session was not replayed.

## Reconciliation

- propose: applied to specs/propose/spec.md
- cycles: applied to specs/cycles/spec.md; its Purpose's definition of stale names the exception, as the change spec says.
- runs: applied to specs/runs/spec.md
