# Changes to implement

The changes that change 261108-jpg8n proposes to the requirements of
[implement](../../../../specs/implement/spec.md). This is future work.

## MODIFIED Requirements

### Requirement: Build a change, a draft or a cycle

When `/grove-implement <id>` is run, the system shall print instructions
for building, on a branch other than the target, the approved change
`<id>`, or each member of cycle `<id>` not yet archived in build order,
each listed with the members of the cycle it depends on and each member's
instructions naming its place in the cycle's plan; or, with
the owner, a change only the working tree holds, as a draft that merging
the branch approves with its code. The agent records each task's status
in `implementation.md`, builds no task already done again unless the
check before finishing finds its evidence no longer holds, and archives
the change when every task is done. The system shall print no
instructions, and say why, for a change that waits on one the working
tree has not archived as finished, or for a draft given `--unattended`;
and a cycle's instructions shall leave out, saying why, a member that
waits on a change neither being built nor archived as finished.

#### Scenario: An approved change

- GIVEN the target approves change C with tasks T-01 and T-02
- WHEN `/grove-implement C` is run on the target branch, with or without
  `--unattended`
- THEN the agent creates a branch, builds T-01 and T-02, records both as
  done in `implementation.md`, and archives C on that branch

#### Scenario: A draft

- GIVEN the owner proposed change D on branch `rename` and committed it,
  and the target does not hold D
- WHEN the owner runs `/grove-implement D` on `rename`
- THEN the instructions say D is a draft that merging `rename` approves
  with its code, and the agent builds and archives D on `rename`; with
  `--unattended`, Grove says the target holds no approved change D

#### Scenario: A change that waits

- GIVEN change B, approved or a draft, with `depends_on: [A]`, and the
  working tree does not hold A archived as finished
- WHEN `/grove-implement B` is run, with or without `--unattended`
- THEN nothing is built, and Grove says B waits on A

#### Scenario: A cycle

- GIVEN cycle K lists A and B, B depending on A, K not marked
  `unattended` and its `needs_human` naming B
- WHEN `/grove-implement K` is run, with or without `--unattended`
- THEN K's instructions list B `after A`, and A and then B are built on
  one branch, each by `grove instructions implement <change-id>`, B's
  instructions naming it the second of two in K's plan; only `grove claim`
  reads `unattended` and `needs_human`

#### Scenario: A member left out

- GIVEN cycle K lists A and C, and C depends on Z, which K does not list
  and the target has not archived as finished
- WHEN `/grove-implement K` is run
- THEN the instructions leave C out, saying it waits on Z, and A is built
