# Changes to requirements

The changes that change 261127-hvq1r proposes to the requirements of
[requirements](../../../../specs/requirements/spec.md), as it stands at
3cfd2ef. This is future work: nothing here describes the target branch
until the change is implemented and the maintained spec is reconciled.

The Purpose's first sentence gains the new reports, as the proposal gives
it under "The text around the requirements".

## ADDED Requirements

### Requirement: Report a new rule over the limit

When the owner runs `grove check`, the system shall report each
requirement an open change ADDs whose rule, its text before its first
scenario, is over 60 words as `wc -w` counts them, naming its file, line,
name and word count.

#### Scenario: A rule of 61 words

- GIVEN an open change ADDs requirement "Export a ledger" with a rule of
  61 words and one scenario
- WHEN the owner runs `grove check`
- THEN it reports the change spec at the requirement's line, saying
  "Export a ledger" has a rule of 61 words against a limit of 60, and
  exits non-zero

#### Scenario: A rule of 60 words with long scenarios

- GIVEN an open change ADDs a requirement with a rule of 60 words and
  scenarios of 400 words
- WHEN the owner runs `grove check`
- THEN the requirement is not reported

#### Scenario: A long rule no open change names

- GIVEN `grove/specs/record-editing/spec.md` holds "Keep a category" with
  a rule of 90 words, and no open change has a block for it
- WHEN the owner runs `grove check`
- THEN "Keep a category" is not reported

### Requirement: Report a modified rule that grew

When the owner runs `grove check`, the system shall report each
requirement an open change MODIFIES whose rule is over 60 words and
longer than the rule the change's `base` holds for it, naming its file,
line, name and both word counts.

#### Scenario: A long rule made longer

- GIVEN the maintained spec at an open change's `base` holds "Keep a
  category" with a rule of 80 words, and the change MODIFIES it with a
  rule of 95
- WHEN the owner runs `grove check`
- THEN it reports the change spec at the requirement's line, saying the
  rule of "Keep a category" grew from 80 words to 95 against a limit of
  60

#### Scenario: A long rule made shorter

- GIVEN the rule is 80 words at `base` and 70 in the change
- WHEN the owner runs `grove check`
- THEN the requirement is not reported

#### Scenario: A short rule taken over the limit

- GIVEN the rule is 50 words at `base` and 65 in the change
- WHEN the owner runs `grove check`
- THEN the requirement is reported

#### Scenario: Only the scenarios change

- GIVEN the rule is 80 words at `base`, and the change MODIFIES the
  requirement with the same rule and one scenario more
- WHEN the owner runs `grove check`
- THEN the requirement is not reported

### Requirement: Report a spec a change grows past the limit

When the owner runs `grove check`, the system shall report each open
change whose requirements, folded into a maintained spec as its `base`
holds it, leave that spec over 2,500 words and larger than at `base`,
naming the change spec, the capability and both word counts.

#### Scenario: Adding to a spec already over

- GIVEN `grove/specs/record-editing/spec.md` holds 3,000 words at an open
  change's `base`, and the change ADDs a requirement of 200 words to it
- WHEN the owner runs `grove check`
- THEN it reports the change's spec for `record-editing`, saying the
  change takes `record-editing` from 3,000 words to 3,200 against a limit
  of 2,500

#### Scenario: Adding and removing

- GIVEN that spec holds 3,000 words at `base`, and the change ADDs a
  requirement of 200 words and REMOVEs one of 300
- WHEN the owner runs `grove check`
- THEN the change is not reported

#### Scenario: Taking a spec over

- GIVEN that spec holds 2,400 words at `base`, and the change ADDs a
  requirement of 200 words
- WHEN the owner runs `grove check`
- THEN the change is reported, from 2,400 words to 2,600

#### Scenario: A spec that stays under

- GIVEN that spec holds 1,000 words at `base`, and the change ADDs a
  requirement of 200 words
- WHEN the owner runs `grove check`
- THEN the change is not reported

#### Scenario: A new capability

- GIVEN no maintained spec for `exports` at an open change's `base`, and
  the change ADDs requirements of 2,600 words to it
- WHEN the owner runs `grove check`
- THEN the change is reported, from 0 words to 2,600

### Requirement: Propose states the limits

When the agent runs `grove instructions propose`, the system shall print,
in the guide to a change spec, that a rule keeps to 60 words, that a
change leaves no spec over 2,500 words larger than it found it, and that
`grove check` reports both.

#### Scenario: The propose instructions

- GIVEN any project
- WHEN the agent runs `grove instructions propose`
- THEN the change spec's guide gives both limits and says `grove check`
  reports a change that exceeds either
