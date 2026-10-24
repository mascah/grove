## T-01 — Propose re-plans the cycle of a revised member

Result: `/grove-propose <change-id>` behaves as
[specs/propose/spec.md](specs/propose/spec.md) "Re-plan the cycle of a
revised member" says, and the plan instructions' finish names that way
past a member that needs a person, as
[specs/cycles/spec.md](specs/cycles/spec.md) "Say whether the cycle can
run unattended" says.

Files: `internal/cli/instructions.go`, `internal/instruct/instruct.go`,
`internal/instruct/requires/propose.md`,
`internal/instruct/requires/plan.md`, `internal/cli/walkthrough_test.go`

Interfaces: `instruct.Revise` also takes the cycle not done that lists
the change, as the working tree holds it, with its file's content, or
none; given one, it prints the file under "## The cycle" and the cycle
document from `metadata/cycle.yaml` with `planned` filled with the
instructions' target commit, as `Replan` does.

Acceptance:
- In `internal/cli/walkthrough_test.go`, on the fixture:
  `instructions propose <member>` for a member of a cycle not done prints
  the cycle's file, `planned: <main's commit>` and the rule never to write
  or remove `unattended`; for a change in no cycle it prints no cycle;
  `instructions plan` prints the shorter finish for a member that needs a
  person.
- On the ledger fixture, in Claude Code: `/grove-propose <budgets-id>` on
  a branch, asked to make budgets checkable without a person, revises
  budgets and the 0.4 cycle file, asks whether budgets still needs a
  person, takes it out of `needs_human` when told so, sets `planned` to
  main's commit and leaves `unattended` as it was; `grove check` is clean.

## T-02 — The stale check counts a member changed with its plan

Result: `grove check` and `grove run` treat a cycle as stale as
[specs/cycles/spec.md](specs/cycles/spec.md) "Count a member changed with
its plan" and [specs/runs/spec.md](specs/runs/spec.md) "Retry" say.

Files: `internal/cli/check.go`, `internal/cli/walkthrough_test.go`

Interfaces: `changedSincePlan`, used by `stalePlans` and by `runnable` in
`internal/cli/run.go`, leaves out a member every commit of which on the
target's first-parent history since `planned`, touching its folder or
named files, also edits the cycle file; its messages keep their words.

Constraints: Git runs through `gitx.Run`, as the check does today; within
the rule 12 budget.

Acceptance:
- In `internal/cli/walkthrough_test.go`, on the fixture: a member revised
  and its cycle re-planned in one squash commit leaves `grove check` clean
  and `grove run` builds the cycle; the same as two commits merged with
  `--no-ff` does too; the same as two commits fast-forwarded onto main is
  stale in both; a member changed on main and the cycle marked
  `unattended` in a later commit is still reported stale.
- After an attempt that paused on a blocked member, one squash merge of
  its revision and re-plan lets `grove run --retry` build the cycle.
- The existing stale-plan tests pass unchanged.

## T-03 — Document the shorter way past a pause

Result: the documentation describes the re-plan propose does and the
stale rule, as the change specs say (rule 8).

Files: `docs/running.md`, `docs/planning.md`, `docs/coordinating.md`,
`docs/reference.md`, `dev/walkthrough.md`

Interfaces: none.

Acceptance:
- `docs/running.md` "The pause" says: revise the change with
  `/grove-propose <change-id>`, merge, run `grove run --retry`, with no
  separate re-plan merge, and that the revision and the re-plan must
  reach main as one commit, squash or merge.
- `docs/planning.md` "Re-plan a cycle" says revising a member with
  `/grove-propose` re-plans its cycle on the same branch; its definition
  of a stale plan names the exception; the quoted finish for a member
  that needs a person names one merge.
- `docs/coordinating.md` says, where it names re-planning, that revising a
  member re-plans its cycle; `docs/reference.md`'s `grove instructions
  propose <change-id>` row says it re-plans a cycle the change is in.
- `dev/walkthrough.md`: planning step 1's finish names the shorter way;
  in the unattended-run section's step 6 ("Run a cycle you name"),
  `/grove-propose <export-id>`
  makes export checkable without a person and re-plans its cycle on the
  same branch, one squash merge, and `grove run --retry
  <export-cycle-id>` then starts a session; in planning step 5, two
  commits fast-forwarded onto main leave the cycle stale.
