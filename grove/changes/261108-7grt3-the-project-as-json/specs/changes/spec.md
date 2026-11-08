# Changes to changes

The changes that change 261108-7grt3 proposes to the requirements of
[changes](../../../../specs/changes/spec.md). This is future work.

## MODIFIED Requirements

### Requirement: Print the project as data

When the owner runs `grove changes --json`, the system shall print one
JSON object, the project value `grove serve` answers `/api/project` with
(serve spec, "The project value"): `root`, `target`, `commit` (the
target's full commit), `read` (the time of the read), `specs` (the
maintained specs), `changes`: every change the target holds, archived
ones included, each with `id`, `path` (its main document under `root`),
`title`, `created`, `description`, `schema`, `depends_on`, `base`,
`state` (`approved`, `archived`, or `dropped` for a change archived
without being built), `waits_on` (what `depends_on` names that is not
archived as finished), `tasks`, each task with `id`, `title` and `status`
(its `Status:` line in `implementation.md`), `cycle` (the cycle not done
that lists it), `ready` (approved and waiting on nothing) and `files`
(its documents under `root`), leaving out a field that has no value; and
`cycles`, `branches`, `runs`, `history` and `problems` as the serve spec
gives them for `/api/project`.

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

#### Scenario: Cycles and attempts

- GIVEN the target holds cycle K, marked `unattended`, listing approved
  changes A and B, B with `depends_on: [A]`, and `origin/grove/runs`,
  fetched, holds a record naming K from host H
- WHEN the owner runs `grove changes --json`
- THEN `cycles` holds K with its members A then B and a verdict that it is
  claimed, B's `cycle` is K, and `runs` holds the record with H

#### Scenario: A branch with work

- GIVEN branch `build-B` from the target holds B's `implementation.md`
  with T-01 `Status: done`, and the target approves B
- WHEN the owner runs `grove changes --json`
- THEN `branches` holds a row naming `build-B` with B in progress and T-01
  done

#### Scenario: The same value as the page

- GIVEN `grove serve` running in the checkout
- WHEN the owner runs `grove changes --json` and asks the server for
  `/api/project`
- THEN the two objects are equal but for `read`

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
