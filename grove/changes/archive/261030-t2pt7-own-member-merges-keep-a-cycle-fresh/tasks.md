## T-01 — Count a cycle's own archived members as expected

Result: a commit that archives a change the cycle itself lists no longer
makes the cycle's plan stale, in `grove check`, `grove run` and the
implement instructions' "Changed on main since cycle <id>'s plan" line,
as "Count a member changed with its plan" in
[specs/cycles/spec.md](specs/cycles/spec.md) gives it.

Files:
- `internal/cli/check.go`
- `internal/cli/walkthrough_test.go`
- `docs/planning.md`
- `dev/walkthrough.md`

Delegated: test layout and private helper names

Interfaces: `sincePlan` keeps its signature; it also skips each commit
that archives a change `cy` lists, reading `cy` as the caller passes it
rather than the target's copy, and its comment says so.

Constraints: `CONSTITUTION.md` rule 10, amended to name the merges that
archived the changes of the cycle itself, is merged into the target
first.

Acceptance:
- `TestFollows` gives cycle L's member d a task naming `ledger/cli.py`;
  after c alone is landed, editing that file, `grove check` is clean and
  `grove run` selects L. This fails at a5b5674.
- A test covers "Then a commit archiving nothing": a later commit that
  archives nothing and edits `ledger/cli.py` makes L stale, and
  `grove instructions implement d` names the file.
- The other scenarios of "Count a member changed with its plan" still
  pass as their tests in `internal/cli/walkthrough_test.go` have them.
- `docs/planning.md` "Re-plan a cycle" says a commit that archives a
  change of this cycle does not count, and that the cycle's changes can
  then reach main one pull request at a time.
- `dev/walkthrough.md`'s "Building a cycle unattended" step has the
  paragraph decision 6 of the proposal describes, and its `grove check`
  is clean as written.
- `just check` passes.
