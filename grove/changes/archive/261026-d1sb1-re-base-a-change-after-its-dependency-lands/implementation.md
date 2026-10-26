# Implementation: Re-base a change after its dependency lands

## Tasks

### T-01 — Revising with propose re-bases a change

Status: done

Evidence:
- `internal/cli/check.go`: the moved-file message reads "T-01's x changed
  on main since base b; revise the change with propose, which also
  re-bases it after a dependency lands"; `grove check` still exits
  non-zero.
- `internal/instruct/requires/propose.md`: revising a change reported for
  a named file changed since its base brings the tasks naming it up to
  date with the target and sets `base`, as for a stale requirement.
- `internal/instruct/requires/implement-owner.md`: when a change in
  `depends_on` brought the file, the agent recommends propose, which
  re-bases the change. The unattended text is unchanged, so a moved file
  still blocks the task.
- `docs/coordinating.md` and `dev/walkthrough.md` (step 5 of
  "Coordinating changes") quote the new message and say to re-base with
  propose once a dependency lands.
- `internal/cli/walkthrough_test.go`: `TestMovedFiles` expects the new
  message, with a non-zero status; `TestImplementInstructions` expects
  the owner's recommendation.
- The on-fixture acceptance with agent sessions (`/grove-propose B`
  re-basing, `/grove-implement B` recommending propose) was not run: see
  Scenarios.
- `just check` passed: 4020 Go lines, 1681 test lines, 7997 instruction
  words.

## Scenarios

- implement "Built after its dependency": unverified; it needs an agent session with the owner (TestImplementInstructions checks the owner's instructions recommend propose when a dependency brought the file)
- propose "The message": `internal/cli/walkthrough_test.go` TestMovedFiles, passed
- propose "A dependency landed": unverified; it needs a propose session on the fixture (the instruction text that sets `base` for a moved file is in `requires/propose.md`)

## Reconciliation

- implement: applied to grove/specs/implement/spec.md
- propose: applied to grove/specs/propose/spec.md
