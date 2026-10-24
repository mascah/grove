# Changes to runs

The maintained spec is [grove/specs/runs/spec.md](../../../../specs/runs/spec.md)
(at ab885af). "Start one attempt" starts from the text change 261023-g0a3h
gives it; this change depends on change 261023-bgtcb, which depends on
that change. Only its new scenario differs. "Refuse a given cycle that
cannot be built" starts from the maintained text.

The Purpose's definition of "a buildable cycle" gains one condition: it
waits on no cycle, that is, the fetched target holds done every cycle its
`depends_on` names (`grove/specs/cycles/spec.md`). "An eligible cycle"
stays a buildable cycle that holds `unattended: true`.

## MODIFIED Requirements

### Requirement: Start one attempt

When the operator runs `grove run` in a clean checkout that is not paused,
and either gives a cycle ID that is a buildable cycle on the fetched target
or gives none and the fetched target holds an eligible cycle, the system
shall, when the checkout has a remote named `origin`, fetch the target from
it and fast-forward the local target branch to it; then create the
attempt's branch from the fetched target, record the attempt for the given
cycle, or else for the oldest eligible cycle, the one whose `created` is
earliest, launch one Claude Code session on the implement skill with that
cycle's ID and `--unattended`, typed `/grove-implement <cycle-id>
--unattended` where the project lists Claude Code as skills and
`/grove:implement <cycle-id> --unattended` where it lists it as commands,
and wait for the session to end.

#### Scenario: Local target has commits of its own

- GIVEN the checkout's local target branch holds a commit `origin` does not
- WHEN the operator runs `grove run`
- THEN no branch is created and no session is launched, and Grove says the
  local target cannot be fast-forwarded

#### Scenario: A marked cycle of two changes, one waiting on the other

- GIVEN `origin`'s target holds cycle K, marked `unattended`, listing
  approved changes A and B, B depending on A, neither archived
- WHEN the operator runs `grove run`
- THEN the checkout is on a new branch from the fetched target, one Claude
  Code session is launched on `/grove-implement K --unattended`, and the
  attempt's record names K, A then B, and the target commit

#### Scenario: Claude Code listed as commands

- GIVEN the fetched target's `grove.yaml` has `agents` with `claude:
  commands`, and eligible cycle K
- WHEN the operator runs `grove run`
- THEN one session is launched on `/grove:implement K --unattended`

#### Scenario: Two marked cycles

- GIVEN `origin`'s target holds eligible cycles K and L planned the same
  day, K's `created` the earlier
- WHEN the operator runs `grove run`
- THEN the attempt selects K only

#### Scenario: The older marked cycle waits on the newer

- GIVEN `origin`'s target holds cycles K and L, both marked `unattended`,
  K's `created` the earlier, K's `depends_on` naming L, and L not done
- WHEN the operator runs `grove run`
- THEN the attempt selects L only

#### Scenario: A given cycle not marked

- GIVEN `origin`'s target holds buildable cycle K without `unattended`,
  and eligible cycle L
- WHEN the operator runs `grove run K`
- THEN one session is launched on the implement skill with K's ID and
  `--unattended`, and the attempt's record names K

#### Scenario: A given cycle while paused

- GIVEN the latest attempt selected cycle L, whose members are not
  archived on the target, and K is buildable
- WHEN the operator runs `grove run K`
- THEN no session is launched and Grove says the project is paused; with
  `--retry`, an attempt starts and selects K

#### Scenario: A project that does not list Claude Code

- GIVEN the fetched target's `grove.yaml` has `agents` without `claude`
- WHEN the operator runs `grove run`
- THEN no session is launched, and Grove says `run` needs Claude Code
  listed under `agents`

#### Scenario: A checkout with no remote

- GIVEN the checkout has no remote named `origin`, and its local target
  branch holds eligible cycle K
- WHEN the operator runs `grove run`
- THEN nothing is fetched, the checkout is on a new branch from the local
  target, one session is launched on the implement skill with K's ID and
  `--unattended`, and the attempt's record holds the local target's commit

#### Scenario: The branch merged locally

- GIVEN the latest attempt, in a checkout with no remote, built cycle K on
  its branch, and the owner merged that branch into the local target
- WHEN the operator runs `grove run`
- THEN the project is not paused

### Requirement: Refuse a given cycle that cannot be built

If the operator gives `grove run` a cycle ID that is not a buildable cycle
on the fetched target, then the system shall launch no session, create no
branch, and say in one sentence why: there is no such cycle, it is done, a
member is under `needs_human`, it waits on a cycle its `depends_on` names
that is not done, its plan is stale, or a member waits on a change neither
in the cycle nor archived as finished.

#### Scenario: A member needs a person

- GIVEN `origin`'s target holds cycle K whose `needs_human` names member C
- WHEN the operator runs `grove run K`
- THEN no branch is created and no session is launched, and Grove says C
  needs a person

#### Scenario: A stale plan

- GIVEN `origin`'s target holds cycle K, and member A was revised on the
  target after K's `planned` commit
- WHEN the operator runs `grove run K`
- THEN no session is launched, and Grove says K's plan is stale

#### Scenario: A cycle it follows not done

- GIVEN `origin`'s target holds cycle K, whose `depends_on` names cycle
  L, and L's change C is not archived
- WHEN the operator runs `grove run K`
- THEN no branch is created and no session is launched, and Grove says K
  waits on cycle L, which is not done

#### Scenario: No such cycle

- GIVEN `origin`'s target holds no cycle with ID Z
- WHEN the operator runs `grove run Z`
- THEN no session is launched, and Grove says there is no cycle Z
