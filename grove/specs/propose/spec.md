# Propose

## Purpose

`/grove-propose` turns a chosen direction into a change the owner can
review, writing a proposal, change specs and tasks that name their files
and interfaces; with `--unattended` it drafts a change from a request with
no one in the session, writing its round as open questions. This covers
how it settles what those documents say, what they must say, and what `grove check` reports on them, such as a named file
that moved since the change's base. A task's named files are the
repo-relative paths on its `Files:` line.

## Boundaries

This covers the documents `/grove-propose` writes and `grove check`'s
reports on a change. Building a change belongs to
`grove/specs/implement/spec.md`, and a
cycle's staleness to `grove/specs/cycles/spec.md`.

## Requirements

### Requirement: Tasks name files and interfaces

When a change's tasks are written, each task shall name, on a `Files:`
line, the repo-relative paths it creates, modifies or tests, one
backticked path per bullet or separated by commas, a path ending in `/`
naming a folder and every file under it, and, on an `Interfaces:` line,
the functions, types or commands it produces that other tasks or callers
use, with no code, commands or steps.

#### Scenario: A task that fixes a function

- GIVEN the owner decided, while proposing, that `find` in
  `ledger/records.py` raises `ValueError` for an unknown ID
- WHEN the agent writes the task
- THEN its `Files:` line names `ledger/records.py` and `tests/test_cli.py`,
  and its `Interfaces:` line says `find` raises `ValueError` with the
  message the requirement states

#### Scenario: A task that updates documentation

- GIVEN the owner wants the README's planned list updated in this change
- WHEN the agent writes the task
- THEN its `Files:` line names `README.md`, and the task has no steps

#### Scenario: One path per bullet

- GIVEN a task touching `ledger/records.py`, `tests/test_cli.py` and the
  fixture folder `tests/fixtures/`
- WHEN the agent writes the task
- THEN `Files:` is followed by three bullets, each one backticked path,
  the last `tests/fixtures/`, and `grove check` reports nothing for it

#### Scenario: A comma line still reads

- GIVEN a change written before this one, whose T-01 has
  `` Files: `ledger/records.py`, `tests/test_cli.py` ``
- WHEN the owner runs `grove check`
- THEN it reports nothing for the line, and Grove reads both files as
  T-01's

### Requirement: Report a named file that moved

When the owner runs `grove check`, the system shall report, once for each
change in no cycle that is not done, each file a task names, or a file
under a folder a task names, that differs on the target from the same file
at the change's `base`, naming each such file once and saying to revise
the change with propose.

#### Scenario: A file changed since base

- GIVEN change C, in no cycle, proposed at `base` b143991, with T-01 naming
  `internal/cli/run.go`, and the target has since merged a commit touching
  that file
- WHEN the owner runs `grove check`
- THEN it reports C's `tasks.md` once, saying `internal/cli/run.go`
  changed since base b143991 and to revise the change with propose

#### Scenario: A file under a named folder

- GIVEN change C, in no cycle, with T-01 naming `tests/fixtures/`, and the
  target has since merged a commit touching `tests/fixtures/records.json`
- WHEN the owner runs `grove check`
- THEN it reports C's `tasks.md`, saying `tests/fixtures/records.json`
  changed since base and to revise the change with propose

#### Scenario: Several files, several tasks

- GIVEN change C, in no cycle, with T-01 naming `ledger/cli.py` and
  `tests/test_cli.py` and T-02 naming `ledger/cli.py`, and the target has
  since changed both files
- WHEN the owner runs `grove check`
- THEN it reports one problem, on C's `tasks.md`, naming
  `ledger/cli.py` and `tests/test_cli.py` once each

#### Scenario: A member of a cycle

- GIVEN change C is a member of cycle K, not done, and T-01's named file
  changed on the target after K's `planned` commit
- WHEN the owner runs `grove check`
- THEN C is not reported against its `base`, and K is reported as stale
  (`grove/specs/cycles/spec.md`)

#### Scenario: A file this change's own branch touched

- GIVEN change C is being built on a branch, where T-01 edited its named
  file
- WHEN the owner runs `grove check` on that branch
- THEN C is not reported, since the file is unchanged on the target

### Requirement: Ask a round before writing

When the owner asks for a new change, the propose instructions shall have
the agent send, before writing any document, one message that states the
scope, offers a split and asks what must land first where either applies,
and asks at most five of the decisions it would otherwise assume, the most
consequential, numbered, each with a recommended answer and its reason
worded so that "yes" accepts it, those that change behavior, scope or
compatibility or lose data first, then those that fix architecture or
interfaces costly to reverse; and wait for the owner. Each answer is
recorded as a decision in the proposal, quoted; a decision the owner
declines to make is recorded as an open question naming the task it
blocks; and every other decision the agent would otherwise assume,
routine choices among them, is recorded as an open question with its
recommended answer, its reason and the task it decides.

#### Scenario: The owner accepts the round

- GIVEN the agent would otherwise assume whether an unknown ID exits 1 or
  2, and where the check goes
- WHEN it asks both with recommended answers and the owner says "yes"
- THEN the proposal records each as a decision quoting the owner's "yes"
  and the recommended answer, and no document lists an assumption

#### Scenario: The round is ordered, and routine choices are not asked

- GIVEN the owner asks for "Editing an unknown ID prints an error", and
  the agent would otherwise assume the exit status, whether `find` in
  `ledger/records.py` raises `ValueError` or `ledger/cli.py` checks the
  ID, and the test's name and place in `tests/test_cli.py`
- WHEN it sends the round
- THEN the message states the scope, asks the exit status first and where
  the check goes second, and does not ask the test's name and place; the
  proposal lists that as an open question with its recommended answer,
  naming T-01

#### Scenario: More than five decisions

- GIVEN the agent would otherwise assume seven decisions that change
  behavior or are costly to reverse
- WHEN it sends the round
- THEN it asks the five most consequential, and the proposal lists the
  other two as open questions, each with its recommended answer, its
  reason and the task it decides

#### Scenario: The owner declines one

- GIVEN the round asks where the check goes
- WHEN the owner says to leave that one open
- THEN the proposal lists it under open questions, naming the task it
  decides, and `/grove-implement` builds its best answer to it with a
  `Check:` line, or blocks the task where a wrong answer would be costly

#### Scenario: Nothing to assume

- GIVEN the conversation settled every decision the documents need
- WHEN the agent reaches the round
- THEN its message states the scope and says there is nothing to ask, and
  it writes the documents once the owner confirms the scope

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

### Requirement: Offer to split a change

When the agent confirms the scope of a new change whose idea is too big
for one change or has parts the owner could agree separately, the propose
instructions shall have it say so and ask; if the owner agrees, it writes
several changes on the same branch, each from its own
`grove instructions propose`, linked with `depends_on`.

#### Scenario: Two parts, one needing the other

- GIVEN the owner asks to propose `ledger delete ID` and an error instead
  of a traceback for an unknown ID, and `delete` finds its record the way
  `edit` does
- WHEN the agent confirms the scope
- THEN it says these are two changes the owner could agree separately,
  and, once the owner agrees, writes the unknown-ID change and the delete
  change on this branch, each from its own `grove instructions propose`,
  the delete change naming the other in `depends_on`

#### Scenario: The owner keeps it whole

- GIVEN the agent offered to split the same idea in two
- WHEN the owner says to keep it as one change
- THEN the agent writes one change

### Requirement: Delegate choices in a task

Where the owner agrees in the round that the implementer may make some
choices without asking, or change some paths beyond a task's `Files:`,
the propose instructions shall have the agent write them on an optional
`Delegated:` line of that task, and write no `Delegated:` line the owner
did not agree to; Grove shall read nothing on that line.

#### Scenario: The owner delegates test layout

- GIVEN the round recommends that the implementer may lay out tests and
  name private helpers as it likes, and may add `tests/test_records.py`,
  and the owner says "yes"
- WHEN the agent writes T-01
- THEN T-01 has a `Delegated:` line naming those choices and
  `tests/test_records.py`, and its `Files:` line does not name
  `tests/test_records.py`

#### Scenario: Nothing delegated

- GIVEN the round asked nothing about delegation, or the owner declined it
- WHEN the agent writes the tasks
- THEN no task has a `Delegated:` line

#### Scenario: Check ignores the line

- GIVEN a change whose T-01 has a `Delegated:` line right after its
  `Files:` line, and the target changed `tests/test_records.py` after the
  change's `base`
- WHEN the owner runs `grove check`
- THEN it reports nothing about T-01: its `Files:` line reads as before,
  and the delegated path is not watched

### Requirement: Re-plan the cycle of a revised member

When the owner runs `/grove-propose <change-id>` for a change that a cycle
not yet done lists, the propose instructions shall print that cycle's file
and have the agent, on the same branch as the revision, set the cycle's
`planned` to the target commit the instructions name, change nothing else
in it, and never write or remove `unattended`.

#### Scenario: A blocked member answered

- GIVEN cycle K, marked `unattended`, lists change B, left blocked by an
  open question
- WHEN the owner runs `/grove-propose B` and answers the question
- THEN the agent revises B, sets K's `planned` to the target commit the
  instructions name, asks nothing about whether B needs a person, and
  leaves `unattended: true` as it is

#### Scenario: A change in no cycle

- GIVEN change D, which no cycle not yet done lists
- WHEN the owner runs `/grove-propose D`
- THEN the instructions print no cycle file, and the agent edits no cycle

### Requirement: Decisions go in the proposal

When the agent runs `grove instructions propose`, the system shall print
instructions that put every decision, however much room it needs, in the
proposal's decisions section, and that name no `design.md`.

#### Scenario: A new change

- GIVEN the default schema with none of the project's own files
- WHEN the agent runs `grove instructions propose "Delete a record"`
- THEN the instructions, the propose guide and the proposal template among
  them, contain no `design.md`

#### Scenario: A decision that needs room

- GIVEN the owner's answer to a decision needs several paragraphs
- WHEN the agent writes the change from the propose instructions
- THEN the decision and its reasons are in `proposal.md` under
  Decisions, and the change folder has no `design.md`

### Requirement: Report a dependency that can never be met

When the owner runs `grove check`, the system shall report, for each open
change, each ID in its `depends_on` that names no change, open or
archived, that names a change archived as dropped, or that names an open
change which waits on this one in turn, directly or through other open
changes, at the `depends_on` line of the change's `.grove.yaml`.

#### Scenario: No such change

- GIVEN open change B's `.grove.yaml` holds `depends_on: [260101-zzzzz]`,
  and no change, open or archived, has that ID
- WHEN the owner runs `grove check`
- THEN it reports B's `.grove.yaml` at its `depends_on` line, saying
  260101-zzzzz is no change, open or archived, and exits non-zero

#### Scenario: A dropped change

- GIVEN open change B depends on A, and A is archived with
  `implementation.md` saying `Status: dropped`
- WHEN the owner runs `grove check`
- THEN it reports B's `depends_on`, saying A was dropped, so B could never
  be applied

#### Scenario: A loop

- GIVEN open change A depends on B, and open change B depends on A
- WHEN the owner runs `grove check`
- THEN it reports the `depends_on` of each, saying the other waits on it in
  turn

#### Scenario: A dependency archived as finished

- GIVEN open change B depends on A, and A is archived with every task
  `Status: done`
- WHEN the owner runs `grove check`
- THEN B's `depends_on` is not reported

### Requirement: Report a stale requirement

When the owner runs `grove check`, the system shall report each
requirement an open change MODIFIES or REMOVES whose text in the working
tree's maintained spec is neither its text at the change's `base`, nor
the text this change gives it, nor the text a change it depends on gives
it, naming the change spec's file and line, saying whether the
requirement was rewritten or removed since `base`, and saying to revise
the change with propose.

#### Scenario: Rewritten by another change

- GIVEN change S, proposed at `base` 4f1c2a9, MODIFIES "Keep a category"
  of `record-editing`, and change T, which S does not depend on, rewrote
  that requirement in the maintained spec and was archived
- WHEN the owner runs `grove check`
- THEN it reports S's change spec at the requirement's line, saying "Keep
  a category" was rewritten since base 4f1c2a9, and to revise the change
  with propose

#### Scenario: Removed by another change

- GIVEN change S MODIFIES "Keep a category", and another change it does
  not depend on removed that requirement from the maintained spec
- WHEN the owner runs `grove check`
- THEN it reports that "Keep a category" was removed since S's `base`

#### Scenario: Rewritten by a change it depends on

- GIVEN change S depends on T, and T's archive gave "Keep a category" the
  text T's change spec gives it
- WHEN the owner runs `grove check`
- THEN S is not reported

#### Scenario: Already reconciled

- GIVEN the maintained spec gives "Keep a category" the text S's own
  change spec gives it, as archiving S leaves it
- WHEN the owner runs `grove check`
- THEN S is not reported

### Requirement: A change that records existing behavior has no tasks

When a change only records behavior the code already has, the propose
instructions shall have the agent write it with no tasks, and `grove
check` and `grove archive` shall accept it with none.

#### Scenario: Recording what edit does

- GIVEN the owner asks for a change that records what `ledger edit` does
  today, and agrees each rule, with no code to change
- WHEN the agent writes the change
- THEN its change spec holds the rules under ADDED and what the code does
  under "Observed current behavior", its `tasks.md` holds no `T-NN`
  task, and `grove check` reports no problem in it

#### Scenario: Implemented with no tasks

- GIVEN change R, drafted on branch `record-edit` with no tasks, and the
  target does not hold R
- WHEN the owner runs `/grove-implement R` on `record-edit`
- THEN the agent builds nothing, folds R's change spec into
  `grove/specs/record-editing/spec.md`, and `grove archive R` moves R to
  `changes/archive/`, since no task waits on `implementation.md`; one
  merge of `record-edit` approves and archives R

#### Scenario: The owner wants the code changed too

- GIVEN, while recording what `ledger edit` does, the owner decides that
  `ledger edit 99 --note x` should print an error instead of a traceback
- WHEN the agent writes the change
- THEN fixing it is a task with its `Files:` line, or the owner leaves
  the traceback under "Known discrepancies" for a later change; the
  change is not written as one that only records behavior

### Requirement: Report a malformed Files line

If a task's `Files:` line holds anything but repo-relative paths in
backticks, one per bullet or separated by commas, then `grove check` shall
report it at its line in `tasks.md`, naming the task and both forms and
saying that a folder ends in `/`, and Grove shall read no file from that
line.

#### Scenario: A path outside the repository

- GIVEN T-01's `Files:` line names `` `../ledger/cli.py` ``
- WHEN the owner runs `grove check`
- THEN it reports the line, saying T-01's `Files:` line should list
  repo-relative paths in backticks, one per bullet or separated by commas,
  and that a folder ends in `/`, and exits non-zero

#### Scenario: A folder

- GIVEN T-01's `Files:` line names `` `tests/fixtures/` ``
- WHEN the owner runs `grove check`
- THEN it reports nothing for the line

### Requirement: List a revised change's documents

When `grove instructions propose <change-id>` is run for a change the
working tree holds, the system shall name each of the change's documents
by path, with the instruction to read every one before editing, and print
none of their text.

#### Scenario: Revising a large change

- GIVEN drafted change D whose documents hold 31 KB of text
- WHEN `grove instructions propose D` is run
- THEN the output names each of D's documents by path and holds none of
  their text, and still prints the templates of the documents propose
  writes

### Requirement: Say that revising re-bases

When `grove check` reports a change's named files as changed on the target
since the change's `base`, the system shall say, in the same sentence,
that revising the change with propose also re-bases it after a dependency
lands.

#### Scenario: The message

- GIVEN change C, in no cycle, proposed at `base` b143991, whose T-01
  names `ledger/cli.py`, which the target has since changed
- WHEN the owner runs `grove check`
- THEN it reports "ledger/cli.py changed on main since base b143991;
  revise the change with propose, which also re-bases it after a
  dependency lands"

### Requirement: Re-base a change whose named file moved

When the owner revises, with `/grove-propose <change-id>`, a change that
`grove check` reports for a named file changed since its `base`, the
propose instructions shall have the agent bring the change up to date with
the file as the target holds it and set `base` to the target's commit.

#### Scenario: A dependency landed

- GIVEN change B depends on change A, and the merge that archived A
  created `ledger/export.py`, which B's T-01 names
- WHEN the owner runs `/grove-propose B`
- THEN B's `base` is the target's commit, B's tasks hold for
  `ledger/export.py` as the target has it, and `grove check` reports
  nothing for B

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

## Known discrepancies

None recorded.

## Open questions

None.

## History

- [261023-6j5bw](../../changes/archive/261023-6j5bw-tasks-name-their-files/) Tasks name their files
- [261023-fca8x](../../changes/archive/261023-fca8x-propose-asks-a-round/) Propose asks a round
- [261023-m2she](../../changes/archive/261023-m2she-verify-before-done/) Verify before done
- [261023-g0a3h](../../changes/archive/261023-g0a3h-implement-builds-unattended/) Implement builds unattended, and apply goes
- [261023-6yafg](../../changes/archive/261023-6yafg-one-merge-for-small-changes/) One merge for small changes
- [261023-yqfpj](../../changes/archive/261023-yqfpj-propose-asks-what-matters/) Propose asks what matters
- [261023-bgtcb](../../changes/archive/261023-bgtcb-revising-a-change-re-plans-its-cycle/) Revising a change re-plans its cycle
- [261023-d6r79](../../changes/archive/261023-d6r79-specs-without-grove-yaml/) Specs without .grove.yaml, and no design.md
- [261023-338jz](../../changes/archive/261023-338jz-plain-docs-exact-specs/) Plain docs, exact specs
- [261023-psd1b](../../changes/archive/261023-psd1b-guide-to-specs-of-existing-code/) A guide to writing specs of existing code
- [261026-ar7ec](../../changes/archive/261026-ar7ec-files-lists-one-path-per-line-and-review-shows-planned-against-actual/) Files lists one path per line, and review shows planned against actual
- [261026-804yw](../../changes/archive/261026-804yw-instructions-list-the-change-s-documents-instead-of-embedding-them/) Instructions list the change's documents instead of embedding them
- [261026-d1sb1](../../changes/archive/261026-d1sb1-re-base-a-change-after-its-dependency-lands/) Re-base a change after its dependency lands
- [261108-jpg8n](../../changes/archive/261108-jpg8n-claim-replaces-run/) grove claim replaces grove run, with attempt records on a branch
- [261114-40gd5](../../changes/archive/261114-40gd5-check-names-each-moved-file-once/) Check names each moved file once per change
- [261114-yx2b7](../../changes/archive/261114-yx2b7-remove-needs-human/) Remove needs_human
- [261114-zpgx9](../../changes/archive/261114-zpgx9-implement-ends-at-a-branch-ready-for-review/) Implement ends at a branch ready for review
- [261115-z5kj4](../../changes/archive/261115-z5kj4-propose-drafts-with-no-one-in-the-session/) Propose drafts a change with no one in the session
- [261122-m78w6](../../changes/archive/261122-m78w6-explore-recommends-and-asks-less/) Explore recommends and asks less
