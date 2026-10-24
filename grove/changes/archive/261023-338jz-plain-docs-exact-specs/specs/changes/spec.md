# Changes to changes

No maintained spec at 0c6f419. `changes` is how Grove lists the changes
the target holds, prints them as data, and moves a change to the archive.
How `grove changes` shows a change's cycle belongs to
`grove/specs/cycles/spec.md`, "Show cycles"; checking the fold before
archiving belongs to change 261023-bfwp5's `requirements` capability.

## ADDED Requirements

### Requirement: List approved changes

When the owner runs `grove changes`, the system shall print one line for
each change the target holds that is not archived, giving its ID, its
title and `approved`, followed, when its `depends_on` names changes the
target does not hold archived as finished, by `waits on` and their IDs.

#### Scenario: A change that waits

- GIVEN the target approves changes A and B, B with `depends_on: [A]`,
  neither archived
- WHEN the owner runs `grove changes`
- THEN A's line ends `approved`, and B's line ends `approved  waits on A`

#### Scenario: Its dependency archived

- GIVEN the target holds A archived as finished and approves B with
  `depends_on: [A]`
- WHEN the owner runs `grove changes`
- THEN A has no line, and B's line names nothing it waits on

#### Scenario: A draft

- GIVEN change D is in the working tree and not in the target
- WHEN the owner runs `grove changes`
- THEN D has no line

### Requirement: Print the project as data

When the owner runs `grove changes --json`, the system shall print one
JSON object holding `root`, `target`, `commit` (the target's full commit)
and `changes`: every change the target holds, archived ones included, each
with `id`, `path` (its main document under `root`), `title`, `created`,
`description`, `schema`, `depends_on`, `base`, `state` (`approved`,
`archived`, or `dropped` for a change archived without being built),
`waits_on` (what `depends_on` names that is not archived as finished),
and `tasks`, each task with `id`, `title` and `status` (its `Status:` line
in `implementation.md`), leaving out a field that has no value.

#### Scenario: One approved, one archived

- GIVEN the target approves B, with `depends_on: [A]` and one task, and
  holds A archived with its one task `Status: done`
- WHEN the owner runs `grove changes --json`
- THEN `changes` holds B with `state` `approved` and `waits_on` absent,
  and A with `state` `archived`, its `path` under `changes/archive/`, and
  its task's `status` `done`

#### Scenario: A dropped change

- GIVEN the target holds change C archived with `implementation.md` saying
  `Status: dropped`
- WHEN the owner runs `grove changes --json`
- THEN C's `state` is `dropped`

#### Scenario: No implementation yet

- GIVEN the target approves C, whose `implementation.md` is not on the
  target, and C has no `description`
- WHEN the owner runs `grove changes --json`
- THEN C's tasks have no `status`, and C has no `description` key

### Requirement: Archive only a finished or dropped change

When the agent or the owner runs `grove archive <change-id>`, the system
shall refuse, naming the first task that is not done, unless the change's
`implementation.md` gives every task in `tasks.md` `Status: done` or says
`Status: dropped`; otherwise it shall move the change's folder to
`changes/archive/` under the root and commit nothing.

#### Scenario: A task not done

- GIVEN change C's `tasks.md` has T-01 and T-02, and its
  `implementation.md` gives T-01 `Status: done` and T-02 `Status: blocked`
- WHEN `grove archive C` runs
- THEN C's folder does not move, and Grove says T-02 is not done in
  `implementation.md`

#### Scenario: Every task done

- GIVEN every task of C is `Status: done` in its `implementation.md`
- WHEN `grove archive C` runs
- THEN C's folder is under `changes/archive/`, Grove says so and that
  nothing is committed, and `git status` shows the move uncommitted

#### Scenario: A dropped change

- GIVEN C's `implementation.md` says `Status: dropped`, and no task is
  done
- WHEN `grove archive C` runs
- THEN C's folder moves to `changes/archive/`

## Observed current behavior

- `grove changes` prints `<id>  <title>  approved`, then `  waits on` and
  the IDs a change waits on, then its cycle, for each change on the target
  not archived; it also prints a change's problems under its line
  (`internal/cli/changes.go`, `listChanges`, at 0c6f419).
- `grove changes --json` prints `root`, `target`, `commit` and the open
  changes followed by the archived ones; the keys and their omission when
  empty come from the JSON tags on `Meta`, `Document`, `Change` and `Task`
  (`internal/doc/doc.go` at 0c6f419). A change's `state` is `approved`
  outside `changes/archive/`, and `archived` or `dropped` inside it.
- `grove archive` refuses with "<task> is not done in implementation.md;
  archive a change when every task is done or the file says Status:
  dropped", else renames the folder into `changes/archive/` and prints
  "Archived <id> to <root>/changes/archive/<folder>/; nothing is
  committed." (`internal/cli/instructions.go`, `archive`, at 0c6f419).
- `docs/reference.md` and `docs/coordinating.md`, "The project as data",
  state these rules at 0c6f419; no maintained spec does.
