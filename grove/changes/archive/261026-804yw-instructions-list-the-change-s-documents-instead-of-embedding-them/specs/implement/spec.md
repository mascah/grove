# Changes to implement

The maintained spec is [grove/specs/implement/spec.md](../../../../specs/implement/spec.md),
as it stands at ea7fbaf. Change 261026-ar7ec, which this change depends
on, also changes it.

## ADDED Requirements

### Requirement: List the change's documents

When `grove instructions implement <change-id>` or `grove instructions
archive <change-id>` is run, the system shall name each document of the
change by its path in the working tree, with the instruction to read every
one before acting, and print none of their text, and, for an approved
change that implement builds, name each document whose working-tree copy
differs from the target's, with the `git show` command that prints the
target's.

#### Scenario: A large change

- GIVEN approved change C whose proposal, change specs and tasks hold
  45 KB of text
- WHEN `grove instructions implement C` is run
- THEN the output names each of C's documents by path and holds none of
  their text, and the agent reads every one before building

#### Scenario: A document revised on the branch

- GIVEN approved change C, and the owner agreed, on this branch, to an edit
  of C's `tasks.md`
- WHEN `grove instructions implement C` is run with the owner
- THEN it names `tasks.md` as differing from the target, with
  `git show <commit>:<path>` for the approved text, and the agent builds
  the revision the owner agreed

#### Scenario: A document that differs, unattended

- GIVEN approved change C, and the working tree's copy of C's change spec
  differs from the target's
- WHEN `grove instructions implement C --unattended` is run
- THEN it names that change spec as differing, and the agent builds what
  `git show` prints for it

#### Scenario: Archiving

- GIVEN change C with every task done
- WHEN `grove instructions archive C` is run
- THEN it names each of C's documents, `implementation.md` included, by
  path, and holds none of their text

### Requirement: Implement leaves revising to propose

When `grove instructions implement <change-id>` is run, the system shall
print the template of `implementation.md` and of no other kind of
document, and the implement instructions shall have the agent run
`grove instructions propose <change-id>` and follow it before editing the
change's proposal, change specs or tasks.

#### Scenario: Templates

- WHEN `grove instructions implement C` is run
- THEN it prints the implementation template, and no proposal, change
  spec, tasks or `.grove.yaml` template

#### Scenario: An agreed revision

- GIVEN the owner agreed, while C is built, to an edit of its change spec
- WHEN the agent makes the edit
- THEN it has first run `grove instructions propose C` and follows what it
  prints
