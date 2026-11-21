<!-- Each task is a section headed "## T-NN — title". IDs are never
renumbered or reused. -->

## T-01 — The index names a shared requirement, and plan tells the owner

Result: the project index names, under each of two approved changes that
change the same requirement and do not depend on each other, the other
change, the capability and the requirement, and the plan instructions
have the agent tell the owner of such a pair among a cycle's members.

Files:
- `internal/doc/doc.go`
- `internal/doc/doc_test.go`
- `internal/instruct/instruct.go`
- `internal/instruct/requires/plan.md`
- `internal/cli/walkthrough_test.go`

Delegated: test layout and private helpers.

Interfaces: `doc.Change.Shares` lists, for an open change, each
requirement another open change also holds a block for: that change's ID,
the capability and the requirement's name, sorted. `doc.Load` fills it.
The index prints one line for each, suggested as `also changed by
<change-id>: <capability> "<requirement>"` (the proposal's open question).

Constraints:
- A pair is found from the ADDED, MODIFIED and REMOVED blocks of the
  change specs alone, by capability and requirement name; the requirement
  texts are not compared.
- A change that depends on the other, directly or through other changes,
  open or archived, is no pair with it.
- `grove check`, `grove changes`, `grove claim` and `grove serve` print
  and decide what they did (`CONSTITUTION.md` rule 10).
- `just check` prints at most 80 instruction words more than it printed
  when 261118-emse7 was archived.

Acceptance:
- With the fixture's two changes that each modify the same requirement of
  `record-editing` approved and neither depending on the other, `grove
  instructions plan "0.5 release"` and `grove instructions propose` each
  print, under each of the two in the project index, the other's ID, the
  capability and the requirement (cycles, "Name a requirement two
  approved changes both change").
- With `depends_on` naming the other added to one of them, or with one
  archived, the index prints no such line.
- The plan instructions say to tell the owner of each such pair among the
  members proposed, that the one built second is stale on the requirement
  once the first is archived, and to ask whether both stay (cycles, "Plan
  tells the owner of members that share a requirement").
- `grove check` reports nothing more than before for a cycle that lists
  both, and `grove claim` claims it.
- A direct test in `internal/doc` covers a pair, a pair ordered by
  `depends_on` through a third change, and the same name in two
  capabilities.
- `just check` passes.

## T-02 — Docs and walkthrough

Result: the documentation says the project index names a requirement two
changes both change and what plan does with it, and the walkthrough runs
it on the fixture.

Files:
- `docs/planning.md`
- `docs/coordinating.md`
- `dev/walkthrough.md`

Acceptance:
- `docs/planning.md`, under "Plan a cycle", says the agent tells you when
  two changes it proposes change the same requirement, what that means
  for the one built second, and that the choice is yours.
- `docs/coordinating.md`, where it describes two changes rewriting the
  same requirement, says the project index names the pair before either
  is built.
- `dev/walkthrough.md` has a step printing the plan instructions with the
  fixture's two changes that modify one requirement, showing the line
  under each, and its note on what the project index lists for an
  approved change says it also names such a pair.
