# Changes to propose

The changes that change 261115-z5kj4 proposes to the requirements of
[`propose`](../../../../specs/propose/spec.md). This is future work:
nothing here describes the target branch until the change is implemented
and the maintained spec is reconciled.

## ADDED Requirements

### Requirement: Draft a change with no one in the session

When the agent runs `grove instructions propose --unattended` with a
request, the flag given as its own argument or as the first word of the
quoted description, the system shall print propose instructions that have
the agent confirm no scope, ask no round and wait for no one, write the
change from the request on a branch other than the target, and end with
the branch ready for review.

#### Scenario: A request from a tracker issue

- GIVEN the ledger fixture with Grove set up, on `main`, and a session
  with no one in it
- WHEN it runs `/grove-propose --unattended Editing an unknown ID prints
  an error and exits 2`
- THEN the agent creates a branch, writes the change's proposal, change
  spec and tasks there, asks nothing, and its last message names the
  files, the qualified task IDs and the open questions

#### Scenario: The flag inside the quotes

- GIVEN a harness whose skill runs `grove instructions propose "<args>"`
- WHEN the agent runs `grove instructions propose "--unattended Editing
  an unknown ID prints an error"`
- THEN it prints what `grove instructions propose --unattended "Editing
  an unknown ID prints an error"` prints, and the request it prints
  holds no `--unattended`

#### Scenario: A schema the caller names

- GIVEN the same fixture
- WHEN the agent runs `grove instructions propose --unattended --schema
  bugfix "Editing a missing record crashes"`
- THEN the instructions follow `bugfix` and have the agent ask no one

#### Scenario: Without the flag

- GIVEN the same fixture
- WHEN the agent runs `grove instructions propose "Editing an unknown ID
  prints an error"`
- THEN the instructions have the agent confirm the scope, ask one round
  and wait for the owner

### Requirement: Write the round as open questions

When the agent writes a change with no one in the session, the propose
instructions shall have it write each decision it would have asked in the
round as an open question, in the round's order, with its recommended
answer, the reason and the task it decides, and record under Decisions
only what the request itself says, quoted.

#### Scenario: The request settles one decision

- GIVEN the request "Editing an unknown ID prints an error and exits 2",
  and the agent would otherwise assume whether `find` in
  `ledger/records.py` raises `ValueError` or `ledger/cli.py` checks the ID
- WHEN it writes the change with no one in the session
- THEN Decisions quotes the request for the exit status, an open question
  recommends that `find` raise `ValueError`, with its reason, and names
  T-01, and no document lists an assumption

#### Scenario: The request settles nothing more

- GIVEN the request "Add ledger delete ID"
- WHEN the agent writes the change with no one in the session
- THEN Decisions says that none was made beyond the request, and every
  choice the documents need is an open question with a recommended answer

### Requirement: Write no change for a request with no outcome

If, with no one in the session, the request does not say what outcome is
wanted, then the propose instructions shall have the agent write no
document and say what the request lacks.

#### Scenario: A request to look into something

- GIVEN the request "Look into the report command"
- WHEN the agent follows the instructions with no one in the session
- THEN it writes no change folder, and its last message says the request
  names no outcome and what would make it one

### Requirement: One request, one change

When the agent writes a change with no one in the session, the propose
instructions shall have it write one change, in the schema the
instructions follow, with no `depends_on` the request does not name, and
name as an open question a split it would have offered, an open change
it would have asked about, or another schema it would have suggested.

#### Scenario: Two parts

- GIVEN the request "Add ledger delete ID, and an error for an unknown ID"
- WHEN the agent writes the change with no one in the session
- THEN it writes one change holding both, and an open question says the
  two could be agreed separately

#### Scenario: An open change touches the same code

- GIVEN approved change A, whose tasks name `ledger/cli.py`, and a
  request whose change also names it and does not mention A
- WHEN the agent writes the change with no one in the session
- THEN its `.grove.yaml` has no `depends_on`, and an open question names A
  as a change that may need to land first

### Requirement: Refuse an unattended propose it cannot write

If `grove instructions propose --unattended` is given a change's ID, or
no request, then the system shall say in one sentence that an unattended
propose writes a new change from a request, print no instructions and
exit non-zero.

#### Scenario: A change's ID

- GIVEN drafted change D
- WHEN the agent runs `grove instructions propose --unattended D`
- THEN Grove says an unattended propose writes a new change from a
  request, prints no instructions and exits non-zero

#### Scenario: No request

- GIVEN the ledger fixture
- WHEN the agent runs `grove instructions propose --unattended`
- THEN Grove says the same, prints no instructions and exits non-zero

## MODIFIED Requirements

### Requirement: Check the change before presenting it

When the agent has written a change's documents, the propose instructions
shall have it check, before it tells the owner the change is written, that
each requirement has an outcome a user or caller can observe, each ADDED
or MODIFIED requirement is covered by a task's acceptance, the proposal's
scope, the change specs and the tasks do not contradict each other, and
no task does what the scope leaves out; fix what it can, tell the owner
what it changed, and ask the owner about what needs a decision or, with
no one in the session, write it as an open question.

#### Scenario: A requirement no task covers

- GIVEN the agent wrote a change whose change spec ADDs "Delete a record"
  and whose only task's acceptance covers editing
- WHEN it checks the change before presenting it
- THEN it adds the missing acceptance or task, and tells the owner it did

#### Scenario: Scope and a task disagree

- GIVEN the proposal puts the README out of scope and a task names
  `README.md`
- WHEN the agent checks the change
- THEN it asks the owner which holds, and changes nothing until they say

#### Scenario: Scope and a task disagree with no one to ask

- GIVEN the same change, written with no one in the session
- WHEN the agent checks the change
- THEN it changes neither, and an open question says the two disagree,
  names the task and must be answered before review

### Requirement: Offer one merge or two

When the agent finishes writing a change the target does not hold, the
propose instructions shall have it offer the owner both next steps:
merge the proposal now to approve it, or run `/grove-implement
<change-id>` on this branch so that one merge approves the change with
its code; shall have it choose neither for the owner; and, with no one in
the session, shall have it offer no `/grove-implement` on the branch and
say that the draft waits for the owner to read and merge it.

#### Scenario: A new change on its branch

- GIVEN the agent wrote change D on branch `rename`, and `grove check`
  reports no problem
- WHEN it finishes
- THEN it names D's files, task IDs and open questions, and offers to
  merge `rename` to approve D, or to build D with `/grove-implement D` on
  `rename` so one merge approves it with its code

#### Scenario: Revising an approved change

- GIVEN the target holds change C, and the owner revises it with
  `/grove-propose C`
- WHEN the agent finishes
- THEN it does not offer to implement C on this branch, since
  `/grove-implement C` builds C as the target holds it

#### Scenario: A draft nobody attended

- GIVEN the agent wrote change D on branch `unknown-id` with no one in
  the session, and `grove check` reports no problem
- WHEN it finishes
- THEN it names D's files, task IDs and open questions, says that merging
  `unknown-id` approves D, and names no `/grove-implement D` on
  `unknown-id`
