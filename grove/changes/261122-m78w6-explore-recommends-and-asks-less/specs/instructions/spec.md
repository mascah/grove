# Changes to instructions

<!-- The changes that change 261122-m78w6 proposes to the requirements of
`instructions`. This is future work: nothing here describes the target
branch until the change is implemented and the maintained spec is
reconciled. -->

Maintained spec: [instructions](../../../../specs/instructions/spec.md).

## ADDED Requirements

### Requirement: Ask the owner plainly

When Grove prints any skill's instructions, they shall tell the agent,
before the skill's own instructions, that when it asks the owner something
it puts the question first after at most a few lines of context, words it
so that it stands on its own, describes each option by what the owner
gains and gives up, says which it recommends and why, keeps the message to
about 150 words unless the owner asks for more, cites a file or commit only
where the owner asks or a document needs it, and names no change ID, task
ID, rule number or decision number without saying in words what it is.

#### Scenario: Every skill is told

- GIVEN any project
- WHEN `grove instructions explore`, `grove instructions propose`,
  `grove instructions plan` or `grove instructions prototype "<question>"`
  runs
- THEN the output, before the skill's own instructions, says to put the
  question first, describe options by what the owner gains and gives up,
  recommend one, and say in words what any ID or rule number named is

#### Scenario: An archived change named in a question

- GIVEN the ledger fixture with archived change 261015-k3m9q, which lets
  an older record keep its retired category on edit
- WHEN the agent asks the owner a question that depends on it
- THEN the question says what the change did, such as "the change that
  lets older records keep their category (261015-k3m9q)", not the ID alone
