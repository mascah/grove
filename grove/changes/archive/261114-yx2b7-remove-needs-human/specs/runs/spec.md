# Changes to runs

The changes that change 261114-yx2b7 proposes to the requirements of
[runs](../../../../specs/runs/spec.md), as it stands at 999ab56. This is
future work: nothing here describes the target branch until the change is
implemented and the maintained spec is reconciled.

The Purpose's "a buildable cycle" loses "has no member `needs_human`
names,".

## MODIFIED Requirements

### Requirement: Refuse a given cycle that cannot be built

If the operator gives `grove claim` a cycle ID that is not a buildable
cycle on the fetched target, then the system shall record nothing, create
no branch, print nothing on standard output, exit non-zero and say in one
sentence why: there is no such cycle, it is done, it waits on a cycle its
`depends_on` names that is not done, its plan is stale, or a member waits
on a change neither in the cycle nor archived as finished.

#### Scenario: A needs_human left over

- GIVEN `origin`'s target holds cycle K, buildable, whose file still holds
  `needs_human` naming member C
- WHEN the operator runs `grove claim K`
- THEN an attempt is recorded for K, and K's ID is printed

#### Scenario: A stale plan

- GIVEN `origin`'s target holds cycle K, and member A was revised on the
  target after K's `planned` commit
- WHEN the operator runs `grove claim K`
- THEN nothing is recorded, and Grove says K's plan is stale

#### Scenario: A cycle it follows not done

- GIVEN `origin`'s target holds cycle K, whose `depends_on` names cycle
  L, and L's change C is not archived
- WHEN the operator runs `grove claim K`
- THEN nothing is recorded, no branch is created, and Grove says K waits
  on cycle L, which is not done

#### Scenario: No such cycle

- GIVEN `origin`'s target holds no cycle with ID Z
- WHEN the operator runs `grove claim Z`
- THEN nothing is recorded, standard output is empty, the exit status is
  non-zero, and Grove says there is no cycle Z
