# Changes to implement

There is no maintained spec for `implement` at 7ce7b91. Every requirement
below is ADDED. "The owner" is the person in whose session the skill runs.

## ADDED Requirements

### Requirement: Implement an approved change

When the owner runs `/grove-implement <change-id>` for a change the target
approves and that waits on nothing the working tree has not archived as
finished, the system shall print instructions for building its tasks in
the owner's session, on a branch other than the target, recording each
task's status in `implementation.md` and archiving the change when every
task is done.

#### Scenario: A change with two tasks

- GIVEN the target approves change C with tasks T-01 and T-02, neither
  started
- WHEN the owner runs `/grove-implement C` on the target branch
- THEN the agent creates a branch, builds T-01 and T-02, records both as
  done in `implementation.md`, and archives C on that branch

#### Scenario: A change that is only a draft

- GIVEN change D exists in the working tree but not in the target
- WHEN the owner runs `/grove-implement D`
- THEN nothing is built, and Grove says the target holds no approved change
  D

#### Scenario: A change that waits

- GIVEN the target approves change B with `depends_on: [A]`, and the
  working tree does not hold A archived as finished
- WHEN the owner runs `/grove-implement B`
- THEN nothing is built, and Grove says B waits on A

### Requirement: Ask the owner instead of blocking

When a task is decided by an open question in the change, or cannot be
built as the change says, the implement instructions shall have the agent
ask the owner and wait for the answer, instead of marking the task
blocked.

#### Scenario: An open question decides a task

- GIVEN change C's proposal has an open question that decides how T-02 is
  built
- WHEN the agent reaches T-02
- THEN it asks the owner the question before building T-02

#### Scenario: The owner leaves it open

- GIVEN the agent asked the owner a question deciding T-02
- WHEN the owner says to leave it
- THEN T-02 is marked `Status: blocked` with the question under it, and the
  agent goes on with the other tasks

### Requirement: Revise the change on its branch

When the owner's answer or feedback changes what the change says, the
implement instructions shall have the agent say which of the change's
proposal, change specs and tasks it would edit and how, and, only once the
owner agrees, edit them on the implementation branch, run `grove check`
clean, and then build to the revised change.

#### Scenario: Feedback changes a requirement mid-build

- GIVEN the agent built T-01 of change C, and the owner, trying it, asks for
  different behavior than C's change spec states
- WHEN the owner agrees to the edit the agent describes
- THEN C's change spec on the branch states the new behavior, the code
  follows it, and `grove review` on the branch shows the edited change spec

#### Scenario: The owner does not agree

- GIVEN the agent proposed an edit to C's change spec
- WHEN the owner declines it
- THEN C's files are unchanged and the agent builds what C says

### Requirement: Apply asks no one

The apply instructions shall keep their current rule: a task an open
question decides, or that cannot be built as the change says, is marked
`Status: blocked` without asking anyone, and the change's files other than
`implementation.md` are not edited.

#### Scenario: The same open question under apply

- GIVEN change C's open question decides T-02
- WHEN C is built with `/grove-apply C`
- THEN T-02 is marked `Status: blocked` with the question, nobody is asked,
  and C's proposal, change specs and tasks are unchanged

## Observed current behavior

- At 7ce7b91, Grove's skills are explore, propose, apply and archive
  (`internal/instruct/skills.yaml`). There is no implement skill.
- At 7ce7b91, `internal/instruct/requires/apply.md` forbids editing the
  change's files other than `implementation.md` and settling an open
  question "even if you could ask"; `grove instructions apply <change-id>`
  refuses a change the target does not approve and one that waits
  (`internal/cli/instructions.go`, `instructionsApply`).
- At 7ce7b91, `docs/getting-started.md` step 5 builds a change with
  `/grove-apply <change-id>`.
