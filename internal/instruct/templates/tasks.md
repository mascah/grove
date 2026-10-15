<!-- Each task is a section headed "## T-NN — title". One task is usually
enough. A task says what result is wanted and how to tell it was achieved; it
refers to requirements instead of repeating them, and leaves ordering and
execution plans out.

IDs are allocated in order and never renumbered or reused. This change's first
task is selected as {{change}}:T-01. Example:

    ## T-01 — Edit a record using its previously valid values

    Result: existing records can be edited without forcing the new-record rule.

    Constraints: new records must still meet the current rule.

    Acceptance:
    - An older record retaining its prior values saves successfully.
    - A newly created record violating the current rule is rejected.
    - The affected behavior and its checks are identified for later reconciliation.

"Constraints:" is optional. A finished task is not a checked box: do not add
status here. -->

## T-01 — [title]

Result: [what is true once this task is done, in one sentence]

Acceptance:
- [a check someone else could carry out]
