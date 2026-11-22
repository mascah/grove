# Instructions

## Purpose

The opening `grove instructions <skill>` prints is what every skill's
instructions say before the skill's own, whatever the skill and schema.

## Boundaries

This covers the opening every skill's instructions share. What each skill
requires and how it works belong to that skill's spec, such as `propose` or
`implement`; the project's `context` is in `config`.

## Requirements

### Requirement: Cite only what the project's readers can open

When Grove prints any skill's instructions, they shall tell the agent to
write Grove documents for whoever reads the project, citing only what those
readers can open and quoting or summarising what only the agent's machine
holds, such as an unpushed branch, a file outside the checkout or a private
project.

#### Scenario: Every skill is told

- GIVEN any project
- WHEN `grove instructions explore`, `grove instructions propose` or
  `grove instructions implement <change-id>` runs
- THEN the output, before the skill's own instructions, says to cite only
  what the project's readers can open and to quote or summarise what only
  this machine holds

#### Scenario: A note in a private project

- GIVEN the owner's notes in another private repository on their machine
  hold a finding the change relies on
- WHEN the agent proposes the change
- THEN the proposal states or quotes the finding and links no path outside
  the checkout

#### Scenario: A source the readers share

- GIVEN a public project whose public issue tracker has an issue on the
  behavior
- WHEN the agent proposes the change
- THEN the proposal may link that issue

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

## Known discrepancies

None recorded.

## Open questions

None.

## History

- [261028-rxsfa](../../changes/archive/261028-rxsfa-cite-what-readers-can-open/) Cite what the project's readers can open
- [261122-m78w6](../../changes/archive/261122-m78w6-explore-recommends-and-asks-less/) Explore recommends and asks less
