# Changes to propose

The changes that change 261124-mhdgj proposes to the requirements of
[propose](../../../../specs/propose/spec.md), as it stands at acf076d.
This is future work: nothing here describes the target branch until the
change is implemented and the maintained spec is reconciled.

## ADDED Requirements

### Requirement: Write a change spec against a change it depends on

Where a change depends on an open change whose change spec adds or
modifies a requirement of the same capability, the change-spec
instructions shall have the agent write that requirement under MODIFIED,
from the text the change it depends on gives it, naming that change, even
where the capability has no maintained spec at the change's `base`.

#### Scenario: A requirement only the dependency adds

- GIVEN open change A ADDS "Show a record" to `record-view`, which has no
  maintained spec at `base` 4f1c2a9, and change B depends on A and alters
  that requirement
- WHEN the agent writes B's change spec for `record-view`
- THEN it holds "Show a record" under MODIFIED, starting from the text A
  gives it, says it starts from change A, and `grove check` does not
  report it

#### Scenario: A requirement the dependency modifies

- GIVEN open change A MODIFIES "Keep a category" of `record-editing`, and
  change B depends on A and alters that requirement again
- WHEN the agent writes B's change spec for `record-editing`
- THEN it holds "Keep a category" under MODIFIED, starting from A's text
  rather than the maintained spec's at `base`, and names change A

## MODIFIED Requirements

### Requirement: Offer to split a change

When the agent confirms the scope of a new change whose idea is too big
for one change or has parts the owner could agree separately, the propose
instructions shall have it say so and ask; if the owner agrees, it writes
every agreed part in that session as its own change on the same branch,
each from its own `grove instructions propose`, linked with `depends_on`,
each later part written against the requirements the earlier parts leave,
and defers no part because an earlier one rewrites the same requirements.

#### Scenario: Two parts, one needing the other

- GIVEN the owner asks to propose `ledger delete ID` and an error instead
  of a traceback for an unknown ID, and `delete` finds its record the way
  `edit` does
- WHEN the agent confirms the scope
- THEN it says these are two changes the owner could agree separately,
  and, once the owner agrees, writes the unknown-ID change and the delete
  change on this branch, each from its own `grove instructions propose`,
  the delete change naming the other in `depends_on`

#### Scenario: Parts that rewrite the same requirement

- GIVEN the owner agrees to split a redesign into three changes, each
  rewriting "Show a record" of `record-view` in turn
- WHEN the agent writes the changes
- THEN it writes all three on this branch in that session, the second
  naming the first in `depends_on` and modifying "Show a record" from the
  text the first gives it, the third likewise from the second, and records
  no part as to be proposed once an earlier one is built

#### Scenario: The owner keeps it whole

- GIVEN the agent offered to split the same idea in two
- WHEN the owner says to keep it as one change
- THEN the agent writes one change
