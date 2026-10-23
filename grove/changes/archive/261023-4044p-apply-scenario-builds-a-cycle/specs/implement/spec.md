# Changes to implement

The changes that change 261023-4044p proposes to the requirements of
`implement` (`grove/specs/implement/spec.md` at 5383ad8). Archiving this
change also removes the entry under "Known discrepancies" that this
requirement resolves.

## MODIFIED Requirements

### Requirement: Apply asks no one

The apply instructions shall keep their rule: a task an open question
decides, or that cannot be built as the change says, is marked
`Status: blocked` without asking anyone, and the change's files other than
`implementation.md` are not edited.

#### Scenario: The same open question under apply

- GIVEN cycle K lists change C, whose open question decides T-02
- WHEN K is built with `/grove-apply K` and the agent reaches C's T-02
- THEN T-02 is marked `Status: blocked` with the question, nobody is asked,
  and C's proposal, change specs and tasks are unchanged
