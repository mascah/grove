# Changes to apply

There is no maintained spec for `apply` at 7ce7b91. Every requirement
below is ADDED. It assumes 261022-d1y6d, which keeps apply's rule that it
asks no one, is archived first.

## ADDED Requirements

### Requirement: Apply a cycle

When `/grove-apply <cycle-id>` is run for a cycle the target holds, the
system shall print the cycle's plan and instructions for building each
member not yet archived, in build order and one at a time, each by running
`grove instructions apply <cycle-id> <change-id>` and following what it
prints to the end, asking no one, and leaving a member that waits on a
change neither in the cycle nor archived as finished as it is.

#### Scenario: Two members, one waiting on the other

- GIVEN cycle K lists A and B, B depending on A, neither archived
- WHEN `/grove-apply K` runs
- THEN A is built and archived, then B, on one branch

#### Scenario: A member blocked

- GIVEN cycle K lists A, and an open question in A decides A's T-02
- WHEN `/grove-apply K` runs
- THEN A's T-02 is marked `Status: blocked` with the question, nobody is
  asked, and A stays open

### Requirement: Apply one member of a cycle

When `grove instructions apply <cycle-id> <change-id>` is run for a member
of a cycle the target holds, the system shall print that change's apply
instructions with the cycle's plan, naming the member's place in build
order.

#### Scenario: The second member

- GIVEN cycle K lists A and B, B depending on A, and A is archived as
  finished on this branch
- WHEN `grove instructions apply K B` is run
- THEN B's apply instructions are printed with K's plan, naming B as the
  second of two

### Requirement: Refuse a member that cannot be applied

If `grove instructions apply <cycle-id> <change-id>` is run for a change
the cycle does not list, or for a member that waits on a change this
working tree has not archived as finished, then the system shall print no
instructions and say why.

#### Scenario: A change outside the cycle

- GIVEN cycle K lists only A, and the target approves C
- WHEN `grove instructions apply K C` is run
- THEN no instructions are printed, and Grove says C is not in K

#### Scenario: A member whose dependency was blocked

- GIVEN cycle K lists A and B, B depending on A, and a task of A is blocked
  so A is not archived
- WHEN `grove instructions apply K B` is run
- THEN no instructions are printed, and Grove says B waits on A

### Requirement: Apply takes only a cycle

If `/grove-apply` is run with a change ID or with no ID, then the system
shall build nothing and say apply takes a cycle ID and that the owner
builds a change with `/grove-implement`.

#### Scenario: A change ID

- GIVEN the target approves change C
- WHEN `/grove-apply C` is run
- THEN nothing is built, and Grove points to `/grove-implement C`

## Observed current behavior

- At 7ce7b91 `grove instructions apply <change-id>` prints instructions for
  one approved change, and `grove instructions apply` with no ID for every
  approved change the working tree has not archived, in build order
  (`internal/cli/instructions.go`, `instructionsApply` and
  `instructionsApplyAll`; `docs/coordinating.md`, "Applying every approved
  change").
- At 7ce7b91 `docs/guides/run-grove-in-ci.md` and `docs/index.md` show
  `/grove-apply` with no ID.
