# Changes to cycles

The changes that change 261118-40s1a proposes to the requirements of
[cycles](../../../../specs/cycles/spec.md), as it stands at f7e0140. This
is future work: nothing here describes the target branch until the change
is implemented and the maintained spec is reconciled.

The Boundaries gain, after "the project index show them", that the
project index naming a requirement two approved changes both change is
covered here too.

## ADDED Requirements

### Requirement: Name a requirement two approved changes both change

When Grove prints a skill's instructions and two approved changes, neither
archived, each hold in their change specs an ADDED, MODIFIED or REMOVED
block for the same requirement of the same capability, and neither
depends on the other, directly or through other changes, the project
index shall name, under each of the two, the other change, the capability
and the requirement.

#### Scenario: Two changes modify one requirement

- GIVEN the target approves changes A and B, each with a change spec for
  `record-editing` that holds a MODIFIED block for "Keep a category", and
  neither names the other in `depends_on`
- WHEN the owner runs `grove instructions plan "0.5 release"`
- THEN the project index names, under A, change B, `record-editing` and
  "Keep a category", and under B, change A with the same

#### Scenario: One removes what the other modifies

- GIVEN A's change spec for `record-editing` holds a REMOVED block for
  "Keep a category" and B's a MODIFIED block for it
- WHEN the owner runs `grove instructions propose`
- THEN the project index names each under the other

#### Scenario: One depends on the other

- GIVEN A and B both modify "Keep a category", and B holds
  `depends_on: [A]`
- WHEN the owner runs `grove instructions plan "0.5 release"`
- THEN the project index names neither under the other

#### Scenario: Through another change

- GIVEN A and C both modify "Keep a category", C depends on B, and B
  depends on A
- WHEN the owner runs `grove instructions plan "0.5 release"`
- THEN the project index names neither under the other

#### Scenario: The same name in two capabilities

- GIVEN A modifies "Keep a category" of `record-editing`, and B adds a
  requirement "Keep a category" to `reports`
- WHEN the owner runs `grove instructions plan "0.5 release"`
- THEN the project index names neither under the other

#### Scenario: One of the two is archived

- GIVEN A and B both modified "Keep a category", and the target holds A
  archived
- WHEN the owner runs `grove instructions plan "0.5 release"`
- THEN the project index names no change under B

### Requirement: Plan tells the owner of members that share a requirement

While agreeing a cycle's members with the owner, when the project index
names a requirement two of them both change, the plan instructions shall
have the agent tell the owner which two changes and which requirement,
and that the one built second is stale on it once the first is archived,
ask whether both stay in this cycle or one goes to a later cycle, and
list under `changes` only what the owner then agrees to.

#### Scenario: The owner moves one to a later cycle

- GIVEN the project index names "Keep a category" of `record-editing`
  under A and B
- WHEN the owner runs `/grove-plan 0.5 release`, the agent proposes A and
  B, and the owner, told of the pair, keeps A and leaves B for later
- THEN the new cycle file lists A and not B

#### Scenario: The owner keeps both

- GIVEN that index
- WHEN the owner, told of the pair, keeps both
- THEN the new cycle file lists A and B, `grove check` reports nothing for
  the cycle, and `grove claim <cycle-id>` claims it when nothing else
  stops it

#### Scenario: No pair among the members

- GIVEN the project index names the pair A and B, and the agent proposes
  A and C for the cycle
- WHEN the owner agrees to A and C
- THEN the agent says nothing about requirements two members share
