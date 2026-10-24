# Changes to propose

The maintained spec is [grove/specs/propose/spec.md](../../../../specs/propose/spec.md).

## ADDED Requirements

### Requirement: Decisions go in the proposal

When the agent runs `grove instructions propose`, the system shall print
instructions that put every decision, however much room it needs, in the
proposal's decisions section, and that name no `design.md`.

#### Scenario: A new change

- GIVEN the default schema with none of the project's own files
- WHEN the agent runs `grove instructions propose "Delete a record"`
- THEN the instructions, the propose guide and the proposal template among
  them, contain no `design.md`

#### Scenario: A decision that needs room

- GIVEN the owner's answer to a decision needs several paragraphs
- WHEN the agent writes the change from the propose instructions
- THEN the decision and its reasons are in `proposal.md` under
  Decisions, and the change folder has no `design.md`

## Observed current behavior

- The propose guide says to write `design.md` only when a decision needs
  more room than the proposal's decisions section
  (`internal/instruct/schemas/default/skills/propose.md` at 0c6f419).
- The proposal template's Decisions comment ends "Use design.md only when
  one needs more room."
  (`internal/instruct/schemas/default/documents/proposal/template.md` at
  0c6f419).
- No schema lists `design.md` as a kind, so no instructions give it a
  template or guide (`internal/instruct/schemas/default/schema.yaml` at
  0c6f419).
