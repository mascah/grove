# Implementation: Cycles planned for unattended runs

<!-- Progress on change 261022-1grb0, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Plan cycles

Status: done

The second acceptance check, `/grove-plan` in Claude Code on the fixture,
is not yet run: the owner, 2026-10-22, "Please mark them done and archive
and I will do the test after and then give feedback if needed."

- `internal/doc`: cycle files `cycles/*.yaml` load as `Cycle`, with
  `title`, `changes`, `planned` required and `goal`, `needs_human`,
  `unattended` optional; `Done` when every member is archived. Diagnostics
  for a member that is no change, a member another cycle not done lists
  first, and a marked cycle whose `needs_human` names a member.
- `internal/instruct`: the `plan` skill (`skills.yaml`, `requires/plan.md`,
  the default schema's `skills/plan.md`, `metadata/cycle.yaml`), `Plan`
  and `Replan`, and cycles not done in the project index.
- `internal/cli`: `grove instructions plan [<description>|<change-id>|<cycle-id>]`;
  `grove check` reports a marked cycle whose member's folder differs on the
  target from `planned` (reading the target only when such a cycle
  exists); `grove changes` names each change's cycle and its mark;
  `grove review` labels `cycles/` files `cycle`.
- Docs: new `docs/planning.md` (in `mkdocs.yml` and `docs/index.md`);
  `docs/reference.md` lists the skill, `grove instructions plan`, the cycle
  file's fields and `rules.skills.plan`; `dev/walkthrough.md` gains
  "Planning cycles".
- Checks: `TestCycles` covers plan instructions (new, one change, revise),
  review, changes, the index, check on a marked cycle with `needs_human`,
  a marked stale cycle, a malformed cycle file and one change in two open
  cycles. `just check` passes: 3105 Go lines, 968 test lines.
- Revised with the owner's agreement (proposal decisions, cycles and runs
  specs, tasks): a cycle has `created`, from Grove's template, required by
  `grove check`; cycles are ordered by it. `TestCycles` checks the
  template's `created` and a cycle file without it.

### T-02 — Run, apply and implement take cycles

Status: done

The second acceptance check, with a remote and Claude Code on the fixture,
is not yet run, as for T-01, at the owner's word. Revised before building, with the
owner's agreement (proposal decisions, apply spec, tasks): apply builds a
cycle one change at a time through `instructions apply <cycle-id>
<change-id>`, which prints the cycle's plan with the change.

- `internal/cli/instructions.go`: `instructions apply <cycle-id>` prints
  the cycle's loop and file; `instructions apply <cycle-id> <change-id>`
  prints that member's apply instructions with the cycle's members in
  build order and its place among them, refusing a change outside the
  cycle and one waiting on a change not archived; a change ID alone or no
  ID is refused, pointing to `/grove-implement`. `instructions implement
  <cycle-id>` prints a loop over `instructions implement <change-id>`.
- `internal/instruct`: `Cycle` replaces `ApplyAll`; `InCycle`;
  `requires/apply-cycle.md` (was `apply-all.md`) and
  `requires/implement-cycle.md`; apply's and implement's skill
  descriptions.
- `internal/cli/run.go`: `grove run` selects the cycle with the earliest
  `created` that is marked, not done, names no member in `needs_human`,
  has no waiting member and is not stale; launches
  `claude -p "/grove-apply <cycle-id>"`; the record gains `cycle`, its
  `changes` being the cycle's in build order, so the pause is per cycle.
- Docs: `docs/running.md`, `docs/coordinating.md` ("Applying a cycle"),
  `docs/reference.md`, `docs/index.md`, `README.md`,
  `docs/guides/run-apply-on-a-timer.md`; also `docs/concepts.md` and
  `docs/getting-started.md`, which said apply takes a change ID or none.
  `docs/guides/run-grove-in-ci.md` types no `/grove-apply` and describes no
  runner, so it is unchanged. `dev/walkthrough.md`: coordinating step 2
  uses implement, planning step 7 applies a cycle, the running steps and
  other paths use cycles.
- Checks: `TestApplyInstructions` (refusals, the cycle loop, one member
  with the plan, a change outside the cycle, implement of a cycle),
  `TestCoordinating` (a waiting member refused, order, a member's place
  after its dependency is archived, a member left out), `TestRun` (changes
  in no cycle, an unmarked cycle, a stale marked cycle, two marked cycles
  the same day selecting the earlier `created` over the lower ID, the
  pause per cycle with another marked cycle, the next cycle after the
  merge). `just check` passes: 3213 Go lines, 1034 test lines.

## Reconciliation

- apply: applied to specs/apply/spec.md, a new maintained spec
- cycles: applied to specs/cycles/spec.md, a new maintained spec
- implement: applied to specs/implement/spec.md; its Purpose and
  Boundaries name cycles, and "Known discrepancies" records that the
  scenario under "Apply asks no one" types `/grove-apply C`, which apply
  now refuses
- runs: applied to specs/runs/spec.md; its Purpose defines an eligible
  cycle in place of eligible changes
