## T-01 — Delivery rules, and implement --runner

Result: `grove.yaml`'s `rules.delivery` prints in its skill's instructions
unless `grove instructions implement` is given `--runner`, which prints
Grove's paragraph in its place, and a cycle names each member's command
with `--runner`, as the implement change spec's requirements say.

Files:
- `internal/instruct/requires/implement-runner.md`
- `internal/instruct/requires/implement-cycle.md`
- `internal/instruct/instruct.go`
- `internal/cli/cli.go`
- `internal/cli/instructions.go`
- `internal/cli/check.go`
- `internal/cli/walkthrough_test.go`
- `docs/reference.md`
- `docs/running.md`
- `dev/walkthrough.md`
- `grove/specs/runs/spec.md`

Delegated: test layout and private helpers

Interfaces:
- `grove instructions implement [--runner] <change-id> | <cycle-id>`.
- `instruct.Rules` gains `Delivery`, keyed by skill.
- `instruct.Implement` and `instruct.Cycle` take whether the caller
  delivers.

Constraints:
- The paragraph is Grove's own text, under `requires/`, so that no project
  file replaces it (rule 10).
- Without `--runner`, and with no `rules.delivery`, the instructions are
  unchanged except the cycle's member commands.

Acceptance:
- `just check` passes, within the instruction budget.
- `grove instructions implement C` prints the lines under
  `rules.delivery.implement` under their own heading after the skill's
  rules.
- `grove instructions implement --runner C` prints the lines under
  `rules.skills.implement`, leaves out those under
  `rules.delivery.implement`, and prints Grove's paragraph.
- The paragraph says nothing about committing.
- `grove instructions implement --runner K`, for a cycle, does the same
  and names `instructions implement --runner <change-id>` for each
  member. Without the flag, a cycle's member commands also carry
  `--runner`.
- `grove instructions plan --runner` prints the usage and exits 2.
- `grove check` reports an unknown `rules.delivery` key as it does for
  `rules.skills`, and accepts the six skill keys.
- `docs/reference.md`'s `rules` row names `delivery`.
- `docs/running.md`'s example puts its push-and-PR line under
  `rules.delivery.implement` and says a caller that delivers the work
  itself passes `--runner`.
- `dev/walkthrough.md` has a step showing both.
- At archive, the runs spec's Boundaries line says the agent pushes and
  opens pull requests only as `rules.delivery` tells it, and not at all
  under `--runner`.
- The walkthrough test checks each of these.
