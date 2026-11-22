# Changes to changes

The changes that change 261122-7vht2 proposes to the requirements of
[changes](../../../../specs/changes/spec.md), as it stands at 07d319f.
This is future work: nothing here describes the target branch until the
change is implemented and the maintained spec is reconciled.

"In flight", "a loose change" and "overlap" are as the
[runs](../../../../specs/runs/spec.md) spec defines them.

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
that lists it), `ready` (approved and waiting on nothing), for a loose
change and no other `claimable` (whether `grove claim` given its ID
would record an attempt now) and `verdict` (one sentence saying that it
would, or why not, as the serve spec's "Readiness of each change in no
cycle" gives it), and `files` (its documents under `root`), leaving out
a field that has no value; and `cycles`, `branches`, `runs`, `history`
and `problems` as the serve spec gives them for `/api/project`.

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
  in flight, B's `cycle` is K, B has no `claimable` key, and `runs` holds
  the record with H and `in_flight` true

#### Scenario: A change in no cycle

- GIVEN the target approves C, which no cycle lists and which waits on
  nothing, and no change is in flight
- WHEN the owner runs `grove changes --json`
- THEN C's `ready` is true, its `claimable` is true and its `verdict`
  says `grove claim <C's ID>` would claim it

#### Scenario: A change in flight

- GIVEN `origin/grove/runs`, fetched, holds a record naming loose change
  C from host H, and the target does not hold C archived
- WHEN the owner runs `grove changes --json`
- THEN C's `ready` is true, its `claimable` is false and its `verdict`
  says it is in flight, naming the attempt and H

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
