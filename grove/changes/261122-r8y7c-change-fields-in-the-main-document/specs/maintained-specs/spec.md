# Changes to maintained-specs

The changes that change 261122-r8y7c proposes to the requirements of
[maintained-specs](../../../../specs/maintained-specs/spec.md), as it
stands at 15816d7. This is future work: nothing here describes the target
branch until the change is implemented and the maintained spec is
reconciled.

Its Boundaries' last sentence becomes: "Where a change keeps its fields
is covered by [changes](../changes/spec.md)."

## MODIFIED Requirements

### Requirement: A maintained spec has no metadata

When the agent runs `grove instructions archive <change-id>` for a change
with a change spec for a capability, the system shall give the `spec`
template with no fields at its top; and if a project's change schema
gives a kind `metadata` other than `frontmatter`, or gives it to a kind
whose path is not a Markdown file directly in `changes/<change>/`, then
`grove check` shall report it.

#### Scenario: Archiving a new capability

- GIVEN finished change C has a change spec for `record-editing`, which
  has no maintained spec
- WHEN the agent runs `grove instructions archive C`
- THEN the instructions give the template for
  `grove/specs/<capability>/spec.md`, which does not start with `---`,
  and none for `grove/specs/<capability>/.grove.yaml`

#### Scenario: A schema with spec metadata

- GIVEN `grove/schemas/bugfix/schema.yaml` lists kind `spec` with
  `metadata: frontmatter`
- WHEN the owner runs `grove check`
- THEN it reports that schema, saying metadata belongs to a Markdown file
  directly in `changes/<change>/`

#### Scenario: A schema with the old metadata

- GIVEN `grove/schemas/short/schema.yaml` lists kind `proposal` with
  `metadata: changes/<change>/.grove.yaml`
- WHEN the owner runs `grove check`
- THEN it reports that schema, saying metadata must be `frontmatter`
