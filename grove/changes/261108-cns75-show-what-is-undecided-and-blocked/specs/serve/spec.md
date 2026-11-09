# Changes to serve

The changes that change 261108-cns75 proposes to the requirements of
[serve](../../../../specs/serve/spec.md). This is future work.

## ADDED Requirements

### Requirement: What is undecided and blocked

When the page asks for `/api/project`, the system shall give each change
`open_questions` as `grove changes --json` does (changes spec, "Show what
is undecided and blocked"), and each change of a branch row the same,
read from the change's main document at that branch's tip or in the
working tree; and the page shall show, in the status column of a cycle's
member row, `open questions` for a change that has them and its tasks
blocked on the target, and, on a branch card, `open questions` beside the
change's task chips.

#### Scenario: A member with questions and a blocked task

- GIVEN cycle K on the target lists change C, whose proposal has an open
  question and whose `implementation.md` on the target gives T-02
  `Status: blocked`
- WHEN the owner opens the page
- THEN C's row under K shows `open questions` and `T-02 blocked`

#### Scenario: A draft with a question

- GIVEN branch `draft-notes` adds change D, which the target does not
  hold, and D's proposal has an open question
- WHEN the page asks for `/api/project`
- THEN the row naming `draft-notes` gives D `open_questions` with the
  question's text, and its card shows `open questions` beside D

#### Scenario: Blocked only on a branch

- GIVEN the target approves B with no `implementation.md`, and branch
  `build-b` holds B's `implementation.md` with T-02 `Status: blocked`
- WHEN the owner opens the page
- THEN `build-b`'s card shows T-02 blocked, and B's row under its cycle
  shows no blocked task
