# Changes to implement

The maintained spec is [grove/specs/implement/spec.md](../../../../specs/implement/spec.md).
"Implement an approved change" is modified by change 261023-g0a3h, which
this change depends on; it is modified here as that change gives it, with
only its draft scenario dropped, since "Implement a draft with the owner"
now covers a draft.

## ADDED Requirements

### Requirement: Implement a draft with the owner

When the owner runs `/grove-implement <change-id>`, without
`--unattended`, for a change this working tree holds and the target does
not, and that waits on nothing the working tree has not archived as
finished, the system shall print instructions for building it as the
working tree holds it, in the owner's session, that say the change is a
draft, that merging the branch into the target approves it with its
code, and that its proposal, change specs and tasks are edited only with
the owner's agreement, as when revising a change on its branch.

#### Scenario: Drafted and built on one branch

- GIVEN the owner proposed change D on branch `rename` and committed it,
  and the target does not hold D
- WHEN the owner runs `/grove-implement D` on `rename`
- THEN Grove's instructions say D is a draft that merging `rename`
  approves with its code, and the agent builds D's tasks, records them as
  done in `implementation.md`, and archives D on `rename`

#### Scenario: Unattended refuses a draft

- GIVEN change D exists in the working tree but not in the target
- WHEN `/grove-implement D --unattended` is run
- THEN nothing is built, and Grove says the target holds no approved
  change D

#### Scenario: Building shows the draft is wrong

- GIVEN the agent, building draft D's T-01, finds that D's change spec
  states behavior the owner, trying it, does not want
- WHEN the agent says which of D's files it would edit and how
- THEN it edits them only once the owner agrees, runs `grove check`
  clean, and builds to the edited change; had the owner declined, D's
  files would be unchanged

#### Scenario: A large draft

- GIVEN draft D has eight tasks across four capabilities
- WHEN the owner runs `/grove-implement D`
- THEN Grove prints the instructions for D as for a draft of one task,
  with nothing about its size

#### Scenario: A draft that waits

- GIVEN draft D has `depends_on: [A]`, and the working tree does not
  hold A archived as finished
- WHEN the owner runs `/grove-implement D`
- THEN nothing is built, and Grove says D waits on A

### Requirement: Review names a change approved with its code

When the current branch holds a change archived as finished that the
target holds neither open nor archived, `grove review` shall say that
merging the branch into the target approves that change with its code,
instead of as future work.

#### Scenario: Drafted, built and archived on one branch

- GIVEN branch `rename` holds change D archived as finished, and the
  target does not hold D
- WHEN the owner runs `grove review` on `rename`
- THEN its first line is "Merging rename into main approves change D
  with its code."

#### Scenario: A draft not yet built

- GIVEN branch `rename` holds change D open, and the target does not hold
  D
- WHEN the owner runs `grove review`
- THEN it says merging `rename` approves change D as future work and does
  not implement it

#### Scenario: An approved change archived on its branch

- GIVEN the target holds change C open, and branch `c-impl` holds C
  archived as finished
- WHEN the owner runs `grove review` on `c-impl`
- THEN it names no change approved by the merge, and lists C's archived
  folder and the maintained specs it changes

## MODIFIED Requirements

### Requirement: Implement an approved change

When the owner runs `/grove-implement <change-id>`, without
`--unattended`, for a change the target approves and that waits on
nothing the working tree has not archived as finished, the system shall
print instructions for building its tasks in the owner's session, on a
branch other than the target, recording each task's status in
`implementation.md` and archiving the change when every task is done.

#### Scenario: A change with two tasks

- GIVEN the target approves change C with tasks T-01 and T-02, neither
  started
- WHEN the owner runs `/grove-implement C` on the target branch
- THEN the agent creates a branch, builds T-01 and T-02, records both as
  done in `implementation.md`, and archives C on that branch

#### Scenario: A change that waits

- GIVEN the target approves change B with `depends_on: [A]`, and the
  working tree does not hold A archived as finished
- WHEN the owner runs `/grove-implement B`
- THEN nothing is built, and Grove says B waits on A
