# Changes to changes

The changes that change 261108-cns75 proposes to the requirements of
[changes](../../../../specs/changes/spec.md). This is future work.

## ADDED Requirements

### Requirement: Show what is undecided and blocked

When the owner runs `grove changes`, the system shall end the line of a
change whose main document's "Open questions" section holds anything
once HTML comments are removed, other than "None" with an optional
period in any case, with `open questions`, and the line of a change whose
`implementation.md` on the target gives tasks `Status: blocked` with
their IDs and `blocked`, each after what the line says today and two
spaces; and when the owner runs `grove changes --json`, it shall give
each change that has open questions `open_questions`, the section's text
without its HTML comments.

#### Scenario: Questions and a blocked task

- GIVEN the target approves change C, whose proposal's "Open questions"
  says "Should IDs be six digits? Blocks T-02.", and whose
  `implementation.md` gives T-01 `Status: done` and T-02 and T-04
  `Status: blocked`
- WHEN the owner runs `grove changes`
- THEN C's line ends `approved  open questions  T-02, T-04 blocked`
- WHEN the owner runs `grove changes --json`
- THEN C's `open_questions` is "Should IDs be six digits? Blocks T-02."

#### Scenario: None

- GIVEN the target approves change D, whose proposal's "Open questions"
  holds the template's HTML comment and "None."
- WHEN the owner runs `grove changes` and `grove changes --json`
- THEN D's line ends `approved`, and D has no `open_questions` key

#### Scenario: Blocked only on a branch

- GIVEN the target approves change B with no `implementation.md`, and
  branch `build-b` holds B's `implementation.md` with T-02
  `Status: blocked`
- WHEN the owner runs `grove changes` on the target's checkout
- THEN B's line names no blocked task

#### Scenario: Waiting and in a cycle

- GIVEN the target approves C, with `depends_on: [A]` and A not archived,
  listed by cycle K, and C's proposal has an open question
- WHEN the owner runs `grove changes`
- THEN C's line ends `waits on A  in cycle K  open questions`
