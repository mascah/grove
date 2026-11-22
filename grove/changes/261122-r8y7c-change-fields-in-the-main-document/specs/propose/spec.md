# Changes to propose

The changes that change 261122-r8y7c proposes to the requirements of
[propose](../../../../specs/propose/spec.md), as it stands at 15816d7.
This is future work: nothing here describes the target branch until the
change is implemented and the maintained spec is reconciled.

## MODIFIED Requirements

### Requirement: Report a dependency that can never be met

When the owner runs `grove check`, the system shall report, for each open
change, each ID in its `depends_on` that names no change, open or
archived, that names a change archived as dropped, or that names an open
change which waits on this one in turn, directly or through other open
changes, at the `depends_on` line of the change's main document.

#### Scenario: No such change

- GIVEN open change B's `proposal.md` holds `depends_on: [260101-zzzzz]`
  among its fields, and no change, open or archived, has that ID
- WHEN the owner runs `grove check`
- THEN it reports B's `proposal.md` at its `depends_on` line, saying
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
- THEN its proposal's fields have no `depends_on`, and an open question
  names A as a change that may need to land first
