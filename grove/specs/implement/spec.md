# Implement

## Purpose

`/grove-implement` builds an approved change, a draft, or each change of a
cycle on a branch, asking the owner where the change does not settle a
task; with `--unattended` it asks no one and marks the task blocked
instead. "The owner" is the person in whose session the skill runs.

## Boundaries

This covers the implement instructions in both modes, `grove.yaml`'s
rules for implement, and how `grove review` names what an implementation
branch approves and builds. Starting a session with no one watching
belongs to `grove run` ([runs](../runs/spec.md)); when a cycle member
counts as changed since its plan, to [cycles](../cycles/spec.md);
revising a change off its implementation branch, to
[propose](../propose/spec.md). What makes a revision approved is the
merge into the target (`CONSTITUTION.md` rule 4).

## Requirements

### Requirement: Build a change, a draft or a cycle

When `/grove-implement <id>` is run, the system shall print instructions
for building, on a branch other than the target, the approved change
`<id>`, or each member of cycle `<id>` not yet archived in build order,
each member's instructions naming its place in the cycle's plan; or, with
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
- THEN A and then B are built on one branch, each by
  `grove instructions implement <change-id>` with the same mode, B's
  instructions naming it the second of two in K's plan; only `grove run`
  reads `unattended` and `needs_human`

#### Scenario: A member left out

- GIVEN cycle K lists A and C, and C depends on Z, which K does not list
  and the target has not archived as finished
- WHEN `/grove-implement K` is run
- THEN the instructions leave C out, saying it waits on Z, and A is built

### Requirement: Name the change's documents

When `grove instructions implement <change-id>` or `grove instructions
archive <change-id>` is run, the system shall name each document of the
change by its path, with the instruction to read every one first, print
none of their text and no template but `implementation.md`'s, and, for an
approved change, name each document whose working-tree copy differs from
the target's with the `git show` command that prints the approved text,
which the agent builds unless the owner agreed to the revision on this
branch.

#### Scenario: A large change

- GIVEN approved change C whose documents hold 45 KB of text
- WHEN `grove instructions implement C` is run
- THEN the output names each of C's documents by path, holds none of
  their text, and prints the implementation template and no other

#### Scenario: A document that differs

- GIVEN approved change C, whose `tasks.md` in the working tree differs
  from the target's
- WHEN `grove instructions implement C` is run
- THEN it names `tasks.md` as differing, with `git show <commit>:<path>`;
  the agent builds the revision if the owner agreed to it on this branch,
  and otherwise, as always with `--unattended`, what `git show` prints

### Requirement: Settle what the change does not

When a task is decided by an open question in the change, cannot be built
as the change says, or shows the change's approach is wrong, the implement
instructions shall have the agent, with the owner, stop that task and ask,
naming which of the proposal, change specs and tasks it would edit and
how, and edit them on the branch only once the owner agrees, by following
`grove instructions propose <change-id>` and then running `grove check`,
marking the task `Status: blocked` with the question if the owner leaves
it open; and, with `--unattended`, mark the task `Status: blocked` with
the question or reason, ask no one, and edit none of the change's files
but `implementation.md`. With the owner, they shall also have it tell the
owner how to try each task it builds and take their feedback before the
next, set a done task whose result an agreed edit changes back to in
progress until it is rebuilt, and stop when `grove check` reports a
requirement of the change as rewritten or removed since its base,
telling the owner the change must be revised with propose and merged
again first.

#### Scenario: An open question, with the owner

- GIVEN change C's open question decides how T-02 is built
- WHEN the agent reaches T-02
- THEN it asks the owner before building T-02; if the owner leaves it,
  T-02 is `Status: blocked` with the question and the agent goes on

#### Scenario: Feedback changes the change

- GIVEN the agent built T-01 of change C and told the owner how to try
  it, and the owner, trying it, asks for behavior C's change spec does
  not state
- WHEN the owner agrees to the edit the agent describes
- THEN the agent edits C's change spec on the branch by following
  `grove instructions propose C`, T-01 is in progress until rebuilt, and
  the agent builds to the edit; had the owner declined, C would be
  unchanged and built as it says

#### Scenario: The approach rests on something false

- GIVEN C's T-01 says `find` raises `ValueError` for an unknown ID, and
  the agent finds `report` relies on `find` returning nothing
- WHEN it makes the finding while building T-01
- THEN it stops T-01 and tells the owner what it found and which of C's
  documents it would edit and how

#### Scenario: Unattended

- GIVEN cycle K lists change C, whose open question decides T-02
- WHEN `/grove-implement K --unattended` reaches T-02
- THEN T-02 is `Status: blocked` with the question, nobody is asked, and
  C's proposal, change specs and tasks are unchanged

### Requirement: Change only the files the tasks name

The implement instructions shall have the agent change only the files the
tasks' `Files:` lines name and, with the owner, also make the choices and
change the paths a task's `Delegated:` line names without asking; for any
other file a task makes false, they shall have it tell the owner and,
only with their agreement, edit it and write `` Also changed: `<path>` —
<reason> `` under that task in `implementation.md`; with `--unattended`,
ignore `Delegated:` lines, write no `Also changed:` line, and name such a
file, and why, in the task's evidence.

#### Scenario: A delegated path

- GIVEN T-01's `Files:` line names `tests/test_cli.py` and its
  `Delegated:` line names `tests/test_records.py`
- WHEN the agent puts `find`'s tests in a new `tests/test_records.py`
- THEN it creates the file without asking the owner

#### Scenario: A file nothing names

- GIVEN the agent, building T-01, finds `README.md`, which no line of T-01
  names, states what T-01 makes false
- WHEN it tells the owner and they agree
- THEN it edits `README.md` and writes `` Also changed: `README.md` ``
  under T-01 with the reason; had they declined, nothing is written

#### Scenario: Unattended

- GIVEN the same T-01, built by `/grove-implement K --unattended`
- WHEN the agent builds it
- THEN neither `tests/test_records.py` nor `README.md` is touched, T-01's
  evidence names both and why, and there is no `Also changed:` line

### Requirement: A task's file changed on the target

When a task names a file, or a folder holding one, that the target changed
since the change's `base`, or, for a member of a cycle not done, since
the cycle's `planned` commit as the cycles spec's "Count a member changed
with its plan" counts it, the system shall name
that task and file in the implement instructions, and they shall have the
agent, with the owner, tell them what changed before building the task
and ask whether to build on the file as it is now or revise the change
with propose, recommending propose, which re-bases the change, when a
change in its `depends_on` brought the file; and, with `--unattended`,
mark the task `Status: blocked` naming the file and not build it.

#### Scenario: Changed since base

- GIVEN approved change C, in no cycle, whose T-02 names
  `internal/cli/run.go`, which the target changed after C's `base`
- WHEN `/grove-implement C` is run
- THEN the instructions name T-02 and `internal/cli/run.go` as changed
  since base; with the owner, the agent asks before building T-02; with
  `--unattended`, T-01 is built and T-02 is `Status: blocked` naming the
  file

#### Scenario: Its dependency landed

- GIVEN change B depends on A, and the merge that archived A on the
  target changed `ledger/export.py`, which B's T-01 names
- WHEN the owner runs `/grove-implement B` and the agent reaches T-01
- THEN it tells the owner what changed and recommends revising B with
  `/grove-propose B`, or else building on the file as it is now

#### Scenario: A member of a cycle

- GIVEN change C is a member of cycle K, and the target changed a file
  C's T-01 names after K's `planned` commit
- WHEN `/grove-implement C` is run
- THEN the instructions name T-01 and the file as changed since K's plan,
  not since C's `base`; once K is re-planned, or when the change came by
  the merge that archived a cycle K depends on, they name nothing

### Requirement: Verify before finishing

When the agent has built every task it can of a change, the implement
instructions shall have it run the project's tests, re-read each done
task's evidence against the code as it now is and set a task whose
evidence no longer holds back to `Status: in progress`, and write a
`## Scenarios` section in `implementation.md` with one line per scenario
of the change specs naming the check that shows it and its result, or
`unverified` and why; its report, and the archive instructions before the
fold, shall name each unverified scenario, refusing nothing and, with
`--unattended`, asking no one.

#### Scenario: Every scenario has a test

- GIVEN a change whose two scenarios are shown by passing tests
- WHEN the agent finishes it
- THEN `## Scenarios` names each with its test and "passed"

#### Scenario: A done task later broken

- GIVEN T-01 was done, and building T-02 broke the test T-01's evidence
  relies on
- WHEN the agent rechecks before finishing
- THEN T-01 is `Status: in progress` again until it is rebuilt

#### Scenario: A scenario only a person can judge

- GIVEN a scenario says a list "reads clearly"
- WHEN the agent finishes, and later archives, the change
- THEN that line says `unverified` and why, the report names it, and
  archive names it to the owner before folding

### Requirement: Rules for implement

Where `grove.yaml` holds lines under `rules.skills.implement`, the system
shall print them with the implement instructions in both modes; when the
agent finishes a change with the owner, who corrected it in a way that
would apply to future work, the instructions shall have it suggest at most
one line for `context` or a key under `rules`, quoting the correction, and
add it on the branch only if the owner agrees, naming it in no task's
`Files:` line; and, with `--unattended`, let its report name at most one
such line and add none.

#### Scenario: A rule for both modes

- GIVEN `rules.skills.implement` holds "Commit after each task, naming its
  task ID."
- WHEN `grove instructions implement C` is run, with or without
  `--unattended`
- THEN the instructions print that line

#### Scenario: The owner agrees

- GIVEN the owner said "Name each test after the requirement it checks"
  while C was built, and made one other correction that would apply to
  future work
- WHEN the agent finishes C and the owner agrees to the line it suggests
- THEN `grove.yaml` on the branch holds that one line, and no task's
  `Files:` line names `grove.yaml`; had the owner declined, or corrected
  only something particular to C, `grove.yaml` is unchanged

#### Scenario: Unattended

- GIVEN A's tests run only with a command `context` does not name
- WHEN `/grove-implement K --unattended` builds A and finishes
- THEN its report names a `context` line giving that command, and
  `grove.yaml` is unchanged

### Requirement: Review names what the branch approves and builds

When the current branch holds a change the target does not hold, `grove
review` shall say whether merging the branch approves it as future work
or, once it is archived as finished on the branch, with its code; and when
the branch holds a change open, or archived where the target does not
hold it archived, `grove review` shall list each file outside the
document root that the branch changes against the target, as planned,
naming the task whose `Files:` line names it or a folder above it, as
agreed, naming the task and the reason its `Also changed:` line in
`implementation.md` gives, or as named nowhere.

#### Scenario: A draft, before and after building

- GIVEN branch `rename` holds change D, which the target does not hold
- WHEN the owner runs `grove review` on `rename` with D open, and again
  once D is archived as finished there
- THEN, with D open, it says merging `rename` approves change D as future
  work and does not implement it; with D archived, its first line is
  "Merging rename into main approves change D with its code."

#### Scenario: Planned, agreed and named nowhere

- GIVEN branch `c-impl` holds approved change C, open or archived, whose
  T-01 names `ledger/records.py` and `tests/fixtures/` on its `Files:`
  line, whose `implementation.md` has `` Also changed: `README.md` — T-01
  makes its usage example false `` under T-01, and the branch changes
  `ledger/records.py`, `tests/fixtures/records.json`, `README.md` and
  `Cargo.lock`
- WHEN the owner runs `grove review` on `c-impl`
- THEN it names no change approved by the merge, and lists
  `ledger/records.py` and `tests/fixtures/records.json` as planned by
  C:T-01, `README.md` as agreed by C:T-01 with that reason, and
  `Cargo.lock` as named nowhere

## Known discrepancies

None recorded.

## History

- [261022-d1y6d](../../changes/archive/261022-d1y6d-implement-with-the-owner/) Implement a change with the owner
- [261022-1grb0](../../changes/archive/261022-1grb0-cycles-for-unattended-runs/) Cycles planned for unattended runs
- [261023-4044p](../../changes/archive/261023-4044p-apply-scenario-builds-a-cycle/) The apply scenario builds a cycle
- [261023-6j5bw](../../changes/archive/261023-6j5bw-tasks-name-their-files/) Tasks name their files
- [261023-m2she](../../changes/archive/261023-m2she-verify-before-done/) Verify before done
- [261023-g0a3h](../../changes/archive/261023-g0a3h-implement-builds-unattended/) Implement builds unattended, and apply goes
- [261023-6yafg](../../changes/archive/261023-6yafg-one-merge-for-small-changes/) One merge for small changes
- [261023-yqfpj](../../changes/archive/261023-yqfpj-propose-asks-what-matters/) Propose asks what matters
- [261024-8ykrz](../../changes/archive/261024-8ykrz-stale-cycle-reported-marked-or-not/) A stale cycle is reported whether or not it is marked
- [261023-mfgew](../../changes/archive/261023-mfgew-implement-suggests-a-rules-line/) Implement suggests a rules line
- [261026-ar7ec](../../changes/archive/261026-ar7ec-files-lists-one-path-per-line-and-review-shows-planned-against-actual/) Files lists one path per line, and review shows planned against actual
- [261026-804yw](../../changes/archive/261026-804yw-instructions-list-the-change-s-documents-instead-of-embedding-them/) Instructions list the change's documents instead of embedding them
- [261026-d1sb1](../../changes/archive/261026-d1sb1-re-base-a-change-after-its-dependency-lands/) Re-base a change after its dependency lands
- [261031-kkzva](../../changes/archive/261031-kkzva-reconcile-the-implement-spec-with-current-behavior/) Reconcile the implement spec with current behavior
