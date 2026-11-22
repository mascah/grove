<!-- Each task is a section headed "## T-NN — title". IDs are never
renumbered or reused. -->

## T-01 — Read the target's schemas for claim and the project value

Result: `grove claim` and the project value find each change's
capabilities through the schemas of the target commit they read, whatever
the checkout's working tree holds, so an overlap on a capability is
refused and shown as runs and serve require.

Files:
- `internal/cli/project.go`
- `internal/cli/claim.go`
- `internal/cli/serve.go`
- `internal/cli/walkthrough_test.go`

Delegated: private helpers, and the test's name and layout.

Interfaces: none that callers see; no command, flag, message or JSON key
changes.

Acceptance:
- A test on the fixture, failing before the fix: a clone made before a
  project schema lands on the target, with A in flight and B sharing A's
  capability through that schema, is refused `grove claim B` naming the
  shared maintained spec and A's attempt (runs, "One capability,
  different requirements"), and its `grove changes --json` gives B
  `claimable` false with a verdict naming that spec (serve, "A change
  that overlaps one in flight").
- `grove instructions propose --schema <name>` for a schema drafted in
  the working tree and not yet on the target still prints that schema's
  templates.
- `grove changes --json` and `/api/project` stay equal but for `read`.
- `just check` passes.
