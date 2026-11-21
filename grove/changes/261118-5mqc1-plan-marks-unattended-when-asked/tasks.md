<!-- Each task is a section headed "## T-NN — title". IDs are never
renumbered or reused. -->

## T-01 — Plan writes the mark when the owner asks

Result: the plan instructions have the agent write or remove `unattended`
only where the owner asks, and say at the finish what merging a marked
cycle file does.

Files:
- `internal/instruct/requires/plan.md`
- `internal/cli/walkthrough_test.go`

Delegated: test layout and private helpers.

Constraints:
- The propose instructions still say never to write or remove
  `unattended` (`internal/instruct/requires/propose.md` unchanged).
- `just check` prints at most 60 instruction words more than it printed
  when 261118-40s1a was archived.

Acceptance:
- `grove instructions plan "0.4 release"` no longer says "Never write or
  remove `unattended`"; it says to write it only where the owner's
  description or the conversation asks, never to offer it, to write it
  even with an obstacle, and, for a cycle being revised, to remove it
  only where the owner asks and keep it otherwise (cycles, "Mark a cycle
  unattended").
- The finish in those instructions says, for a cycle file holding
  `unattended: true`, that merging it hands the cycle to a `grove claim`
  given no cycle, and otherwise that `unattended: true`, added in a merge
  of its own, does, still after naming `/grove-implement <cycle-id>` and
  whether `grove claim <cycle-id>` would claim it (cycles, "Say whether
  the cycle can run unattended").
- A cycle file merged with `unattended: true` in the merge that plans it
  is claimed by a `grove claim` given no cycle (cycles, "The owner asks
  for an unattended cycle").
- `grove instructions propose <member-id>` for a member of a cycle not
  done still says never to write or remove `unattended`.
- `just check` passes.

## T-02 — Docs and walkthrough

Result: the documentation says plan writes the mark when you ask for an
unattended cycle, and the walkthrough runs it on the fixture.

Files:
- `docs/planning.md`
- `dev/walkthrough.md`

Acceptance:
- `docs/planning.md`, under "Plan a cycle", says the agent writes
  `unattended: true` when you ask it for an unattended cycle and what its
  verdict then says; under "Mark it unattended", it says the mark is
  added in a merge of its own when you did not ask plan for it, that only
  plan writes or removes it and only when you ask, and that re-planning
  keeps it otherwise.
- `dev/walkthrough.md`, "Running unattended" step 4, plans the budgets
  cycle with `/grove-plan <budgets-id>`, saying in the session that it
  should run unattended; the cycle file holds `unattended: true`, and one
  merge and push replaces the separate mark. Its steps where the owner
  asks nothing still show a cycle file with no `unattended`.
