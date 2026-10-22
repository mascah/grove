# Changes to apply

There is no maintained spec for `apply` at 7ce7b91. Every requirement
below is ADDED. It assumes 261022-d1y6d, which keeps apply's rule that it
asks no one, is archived first.

## ADDED Requirements

### Requirement: Apply a cycle

When `/grove-apply <cycle-id>` is run for a cycle the target holds, the
system shall print instructions for building each member not yet archived,
in build order, each with its change's apply instructions, asking no one,
and leaving a member that waits on a change neither in the cycle nor
archived as finished as it is.

#### Scenario: Two members, one waiting on the other

- GIVEN cycle K lists A and B, B depending on A, neither archived
- WHEN `/grove-apply K` runs
- THEN A is built and archived, then B, on one branch

#### Scenario: A member blocked

- GIVEN cycle K lists A, and an open question in A decides A's T-02
- WHEN `/grove-apply K` runs
- THEN A's T-02 is marked `Status: blocked` with the question, nobody is
  asked, and A stays open

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
