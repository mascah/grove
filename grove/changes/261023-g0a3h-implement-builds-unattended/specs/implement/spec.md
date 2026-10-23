# Changes to implement

The maintained spec is [grove/specs/implement/spec.md](../../../../specs/implement/spec.md).
"Verify before finishing" is added to it by change 261023-m2she, which this
change depends on; it is modified here as that change gives it.

## ADDED Requirements

### Requirement: Build unattended

When `/grove-implement <id> --unattended` is run for a change the target
approves or a cycle the target holds, the system shall print instructions
that have the agent build it as `/grove-implement <id>` would but ask no
one: a task an open question decides, or that cannot be built as the
change says, is marked `Status: blocked` with the question or the reason;
a file no task names is left as it is and named in that task's evidence;
the change's files other than `implementation.md` are not edited; and a
cycle is built whether or not it is marked `unattended` and whatever its
`needs_human` says.

#### Scenario: One change, by hand

- GIVEN the target approves change C, in no cycle, with tasks T-01 and
  T-02 that no open question decides
- WHEN `/grove-implement C --unattended` is run on the target branch
- THEN the agent creates a branch, builds T-01 and T-02 without asking
  anyone, records both as done in `implementation.md`, and archives C on
  that branch

#### Scenario: An open question decides a task

- GIVEN cycle K lists change C, whose open question decides T-02
- WHEN `/grove-implement K --unattended` is run and the agent reaches C's
  T-02
- THEN T-02 is marked `Status: blocked` with the question, nobody is
  asked, and C's proposal, change specs and tasks are unchanged

#### Scenario: Two members, one waiting on the other

- GIVEN cycle K lists A and B, B depending on A, neither archived
- WHEN `/grove-implement K --unattended` is run
- THEN A is built and archived, then B, on one branch, each by
  `grove instructions implement <change-id> --unattended`

#### Scenario: A cycle that needs a person, by hand

- GIVEN cycle K lists A and B, is not marked `unattended`, and
  `needs_human` names B
- WHEN `/grove-implement K --unattended` is run by hand
- THEN A and B are built asking no one; only `grove run` refuses K

#### Scenario: A change that waits

- GIVEN the target approves change B with `depends_on: [A]`, and the
  working tree does not hold A archived as finished
- WHEN `/grove-implement B --unattended` is run
- THEN nothing is built, and Grove says B waits on A

### Requirement: A member gets its cycle's plan

When `/grove-implement <change-id>` is run, with or without
`--unattended`, for a member of a cycle the target holds that is not
done, the system shall print that change's instructions with the cycle's
plan, naming the member's place in build order.

#### Scenario: The second member, unattended

- GIVEN cycle K lists A and B, B depending on A, and A is archived as
  finished on this branch
- WHEN `grove instructions implement B --unattended` is run
- THEN B's instructions are printed with K's plan, naming B as the second
  of two

#### Scenario: A member built with the owner

- GIVEN cycle K lists A and B, neither archived
- WHEN the owner runs `/grove-implement A`
- THEN A's instructions are printed with K's plan, naming A as the first
  of two

### Requirement: Implement's rules in both modes

Where `grove.yaml` holds lines under `rules.skills.implement`, the system
shall print them with the implement instructions, for a change and for a
cycle, with and without `--unattended`.

#### Scenario: A rule for both modes

- GIVEN `rules.skills.implement` holds "Commit after each task, naming its
  task ID."
- WHEN `grove instructions implement C` runs, and then
  `grove instructions implement C --unattended`
- THEN both print that line with the implement instructions

### Requirement: Report rules for apply

When the owner runs `grove check` and `grove.yaml` holds a
`rules.skills.apply` key, the system shall report that line, saying its
lines go under `rules.skills.implement`.

#### Scenario: A project that kept its apply rules

- GIVEN `grove.yaml` holds two lines under `rules.skills.apply`
- WHEN the owner runs `grove check`
- THEN it reports the `apply` line, pointing to `implement`

### Requirement: No apply skill

The system shall write no apply skill or command for any harness, `grove
init` shall remove the apply skill and command files an earlier Grove
wrote, and `grove instructions apply` shall print the usage and exit with
status 2.

#### Scenario: Init after upgrading

- GIVEN a project whose `.claude/skills/grove-apply/SKILL.md` an earlier
  Grove wrote
- WHEN the owner runs `grove init`
- THEN that folder is removed, and init names `grove-apply` among what it
  removed

#### Scenario: The old command

- GIVEN cycle K on the target
- WHEN `grove instructions apply K` is run
- THEN no instructions are printed, and Grove prints the usage

## MODIFIED Requirements

### Requirement: Implement an approved change

When the owner runs `/grove-implement <change-id>`, without
`--unattended`, for a change the target approves and that waits on
nothing the working tree has not archived as finished, the system shall
print instructions for building its tasks in the owner's session, on a
branch other than the target, recording each task's status in
`implementation.md` and archiving the change when every task is done.

#### Scenario: A change with two tasks

- GIVEN the target approves change C with tasks T-01 and T-02, neither
  started
- WHEN the owner runs `/grove-implement C` on the target branch
- THEN the agent creates a branch, builds T-01 and T-02, records both as
  done in `implementation.md`, and archives C on that branch

#### Scenario: A change that is only a draft

- GIVEN change D exists in the working tree but not in the target
- WHEN the owner runs `/grove-implement D`, with or without
  `--unattended`
- THEN nothing is built, and Grove says the target holds no approved change
  D

#### Scenario: A change that waits

- GIVEN the target approves change B with `depends_on: [A]`, and the
  working tree does not hold A archived as finished
- WHEN the owner runs `/grove-implement B`
- THEN nothing is built, and Grove says B waits on A

### Requirement: Ask the owner instead of blocking

When, without `--unattended`, a task is decided by an open question in the
change, or cannot be built as the change says, the implement instructions
shall have the agent ask the owner and wait for the answer, instead of
marking the task blocked.

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

### Requirement: Implement a cycle

When the owner runs `/grove-implement <cycle-id>`, without
`--unattended`, for a cycle the target holds, the system shall print
instructions for building each member not yet archived, in build order,
each as `/grove-implement <change-id>` builds it, in the owner's session.

#### Scenario: A cycle the owner builds

- GIVEN cycle K lists A and B, is not marked `unattended`, and
  `needs_human` names B
- WHEN the owner runs `/grove-implement K`
- THEN A and then B are built in the owner's session, the agent asking the
  owner where `/grove-implement` would

### Requirement: Verify before finishing

When the agent has built every task it can of a change, with the owner or
with `--unattended`, the implement instructions shall have it, before it
finishes, run the project's tests, re-read each done task's evidence
against the code as it now is and set a task whose evidence no longer
holds back to `Status: in progress`, and write a `## Scenarios` section in
`implementation.md` with one line per scenario of the change specs naming
the check that shows it and its result, or `unverified` and why; its
finishing report names each unverified scenario, and with `--unattended`
it asks no one.

#### Scenario: Every scenario has a test

- GIVEN a change whose change spec has two scenarios, both shown by tests
  that pass
- WHEN the agent finishes implementing it
- THEN `implementation.md`'s `## Scenarios` names each scenario with its
  test and "passed"

#### Scenario: A done task later broken

- GIVEN T-01 was marked done, and building T-02 changed the code T-01's
  evidence relies on so its test fails
- WHEN the agent rechecks before finishing
- THEN T-01 is `Status: in progress` again until it is rebuilt

#### Scenario: A scenario only a person can judge

- GIVEN a scenario says a list "reads clearly"
- WHEN the agent writes the scenario list
- THEN that line says `unverified` and why, and the finishing report
  names it

#### Scenario: A scenario no test shows, unattended

- GIVEN a member, built with `--unattended`, whose change spec has a
  scenario no test checks
- WHEN the agent finishes the member
- THEN that line says `unverified` and why, the report names it, and the
  agent goes on without asking

## REMOVED Requirements

### Requirement: Apply asks no one

Apply is removed; its rule is now "Build unattended", which
`/grove-implement <id> --unattended` follows.
