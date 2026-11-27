<!-- Each task is a section headed "## T-NN — title". IDs are never
renumbered or reused. -->

## T-01 — `grove check` reports a rule or a spec a change grows past its limit

Result: for each open change, `grove check` reports a new rule over 60
words, a modified rule over 60 words that grew, and a maintained spec the
change leaves over 2,500 words and larger than at its `base`.

Files:
- `internal/cli/check.go`
- `internal/doc/requirements.go`
- `internal/doc/requirements_test.go`
- `internal/cli/walkthrough_test.go`
- `grove/specs/requirements/spec.md`

Interfaces: `doc.Requirement` gives the number of words in its rule, the
text before its first `#### Scenario:` heading; the three reports are
diagnostics of the change, printed by `grove check` as `path:line:
message`.

Constraints:
- The two limits are constants in Grove; `grove.yaml` gains no field.
- A requirement under REMOVED, an archived change, and a maintained spec
  that no open change has a block for are never reported.
- Each report is one sentence (`CONSTITUTION.md` rule 14).
- `grove claim` and `grove archive` refuse nothing they did not refuse
  before: `check` reports and blocks nothing (`CONSTITUTION.md` rule 10).
- The first sentence of the requirements spec's Purpose becomes what the
  proposal gives under "The text around the requirements".

Acceptance:
- On the fixture, an open change that ADDs a requirement with a rule of
  61 words makes `grove check` print the change spec's path, the
  requirement's line and name and the count, and exit non-zero; with a
  rule of 60 words and scenarios of any length it prints nothing for it
  (requirements, "Report a new rule over the limit").
- A MODIFIED requirement is reported when its rule is over 60 words and
  longer than at the change's `base`, and not when it is shorter, the
  same, or changed only in its scenarios (requirements, "Report a
  modified rule that grew").
- A change is reported for a maintained spec it leaves over 2,500 words
  and larger than at `base`, with both counts, and not when what it
  removes outweighs what it adds or the spec stays under; a new
  capability is counted from 0 (requirements, "Report a spec a change
  grows past the limit").
- A maintained spec holding a rule over 60 words that no open change
  names is not reported, and `grove check` in Grove's own repository,
  with this change open, prints `No problems in grove/.`
- A direct test in `internal/doc` covers the rule's word count: a rule of
  exactly 60 words, a rule with a fenced code block, and a requirement
  with no scenario.
- `just check` passes.

## T-02 — The change-spec guide, the reference and the walkthrough

Result: the propose instructions state both limits in the guide to a
change spec, the reference names the new reports and where to go when a
spec is over its limit, and the walkthrough runs one report.

Files:
- `internal/instruct/schemas/default/documents/change-spec/guide.md`
- `docs/reference.md`
- `dev/walkthrough.md`
- `internal/cli/walkthrough_test.go`

Acceptance:
- `grove instructions propose` prints, in the change spec's guide, that a
  rule keeps to 60 words, that a change leaves no spec over 2,500 words
  larger than it found it, and that `grove check` reports both
  (requirements, "Propose states the limits").
- `docs/reference.md` names, in what `grove check` reports, a rule or a
  spec that a change grows past its limit, linking the requirements spec
  and the guides to splitting and reconciling a maintained spec.
- `dev/walkthrough.md` has a step in which a change's rule over 60 words
  is reported and shortening it clears the report, and
  `internal/cli/walkthrough_test.go` follows it.
- `just check` passes, with the instruction text within its budget.
