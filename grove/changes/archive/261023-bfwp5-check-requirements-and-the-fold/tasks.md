## T-01 — Check each requirement's shape

Result: `grove check` reports requirements and maintained specs as
[specs/requirements/spec.md](specs/requirements/spec.md) "Each
requirement says shall and has a scenario" and "No change headings in a
maintained spec" say.

Files: `internal/doc/requirements.go`, `internal/doc/requirements_test.go`,
`internal/cli/check.go`

Interfaces: `doc.Requirement` reports its own shape problems, used by
`check` for maintained specs and open change specs.

Acceptance:
- On the ledger fixture, a maintained requirement with its scenario
  deleted, and one with "shall" reworded, are each reported at their line;
  a REMOVED requirement and an archived change are not.
- A `## ADDED Requirements` heading pasted into a maintained spec is
  reported.
- `grove check` on Grove's own repository stays clean.

## T-02 — Archive refuses a fold that does not match

Result: `grove archive` compares the fold as
[specs/requirements/spec.md](specs/requirements/spec.md) "Archive refuses
a fold that does not match" says, from a file of its own.

Files: `internal/cli/archive.go`, `internal/cli/instructions.go`,
`internal/cli/walkthrough_test.go`

Interfaces: none beyond `grove archive`'s messages; `archive` moves from
`instructions.go` to `archive.go` unchanged in what it takes.

Constraints: the comparison reads the specs where this branch left the
target with Git's merge base, as `gitx` reads trees today.

Acceptance:
- In `walkthrough_test.go`, on the fixture: a fold copied word for word
  archives; a reworded MODIFIED, a REMOVED left in and an unrelated
  requirement edited each refuse with their `path:line: message`; reflowed
  lines archive; the second change of a cycle archives after the first on
  the same branch; a dropped change archives with no comparison.
- The existing walkthrough tests that archive a change fold its specs
  first, as the archive guide says.

## T-03 — Document the checks

Result: the documentation and the archive instructions describe both
checks, as [specs/requirements/spec.md](specs/requirements/spec.md)
"Archive copies requirements word for word" says.

Files: `internal/instruct/requires/archive.md`,
`internal/instruct/schemas/default/skills/archive.md`, `docs/reference.md`,
`docs/concepts.md`, `dev/walkthrough.md`

Interfaces: none.

Acceptance:
- `docs/reference.md` lists both new `check` problems and archive's new
  refusal; `docs/concepts.md`'s list of what Grove requires at each step
  names it (rule 8).
- `dev/walkthrough.md` has a step where archive refuses a reworded
  requirement and accepts it once copied.
- `grove instructions archive` prints the word-for-word rule.
