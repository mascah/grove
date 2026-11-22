<!-- Each task is a section headed "## T-NN — title". IDs are never
renumbered or reused. -->

## T-01 — The implement instructions name a rewritten requirement

Result: `grove instructions implement` names each requirement the change
modifies or removes that was rewritten or removed since its `base`, and
tells the agent to reconcile the change before building.

Files:
- `internal/cli/instructions.go`
- `internal/cli/check.go`
- `internal/instruct/requires/implement.md`
- `internal/cli/walkthrough_test.go`

Delegated: private helpers, test names and layout, and the exact wording
of the line the instructions print.

Constraints:
- `grove check` reports a stale requirement in the words it uses now.
- The instructions for a change with no such requirement are what they
  are now.

Acceptance:
- For a change in no cycle whose MODIFIED requirement another archived
  change rewrote after its `base`, `grove instructions implement
  <change-id>` names the capability and the requirement as rewritten
  since that base, and says to revise the change before building any
  task, keeping what the maintained spec now says and what the change
  alters, and to write a `Check:` line (implement, "Rewritten by a change
  that landed first").
- For a requirement another change removed, the line says removed
  (implement, "Removed by a change that landed first").
- For a change whose requirements are as they were at its `base`, or
  were rewritten only by a change it depends on, the instructions name
  none (implement, "Nothing rewritten" and "Rewritten by a change it
  depends on").
- With `--runner`, for a member of a cycle whose requirement an earlier
  member rewrote when it was archived on the same branch, the
  instructions name it (implement, "Two members of one cycle").
- `just check` passes.

## T-02 — Docs and walkthrough

Result: the documentation says what implement does with a requirement
rewritten on the target, and the walkthrough runs it on the fixture.

Files:
- `docs/coordinating.md`
- `dev/walkthrough.md`

Acceptance:
- `docs/coordinating.md`, where it shows `grove check` reporting a
  requirement rewritten since base, says that `/grove-implement` names
  the same requirement and that the agent revises the change on its
  branch before building, the merge approving the revision with the code.
- `dev/walkthrough.md`, at the step where `grove check` reports a
  requirement rewritten since base, has a fresh session's implement
  instructions name that requirement.
