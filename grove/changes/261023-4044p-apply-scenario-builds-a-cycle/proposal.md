# The apply scenario builds a cycle

## Problem

`grove/specs/implement/spec.md` at 5383ad8 records, under "Known
discrepancies", that the scenario under "Apply asks no one" builds C with
`/grove-apply C`, while apply has taken only a cycle ID since change
261022-1grb0 (`grove/specs/apply/spec.md` at 5383ad8, "Apply takes only a
cycle"; `internal/cli/instructions.go`, `instructionsApply`). The
discrepancy arose when 261022-1grb0 was revised on its implementation
branch to make apply cycle-only: its implement change spec
(`grove/changes/archive/261022-1grb0-cycles-for-unattended-runs/specs/implement/spec.md`)
only ADDED "Implement a cycle" and carried no MODIFIED "Apply asks no
one", so archiving could not fold the fix in.

## Intended outcome

The maintained implement spec's scenario under "Apply asks no one" builds C
as a member of a cycle, and the spec records no known discrepancy.

## Scope

### In scope

- The scenario under "Apply asks no one" in the implement spec.

### Out of scope

- The rule itself, which stays as it is.
- The apply spec, which already states that apply takes only a cycle.
- Code and documentation pages: none describes this scenario.

## Decisions

- **A docs-only follow-up change of its own.** The owner, 2026-10-23,
  offered a follow-up change, folding the fix into a larger change, or
  editing the maintained spec directly on a branch: "4 - 1." Reason: the
  fix is one requirement and waits on nothing.

## Open questions

None must be answered before review. None blocks a task.

## Related documents

- Requirements: [specs/implement/spec.md](specs/implement/spec.md)
  (maintained spec: `grove/specs/implement/spec.md`)
- Tasks: [tasks.md](tasks.md)
