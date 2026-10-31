# Changes to cycles

The maintained spec is [grove/specs/cycles/spec.md](../../../../specs/cycles/spec.md)
(at a5b5674).

The Purpose's definition of stale becomes: a cycle is stale when the
folder of a member not yet archived, or a file one of its tasks names,
differs on the target from the same path at the cycle's `planned` commit,
unless every commit on the target's first-parent history since then that
changed it also edited the cycle file, archived a change the cycle itself
lists, or archived a change that a cycle this one depends on, directly or
through other cycles, lists.

## MODIFIED Requirements

### Requirement: Count a member changed with its plan

When a commit on the target's first-parent history after a cycle's
`planned` commit changes the folder of a member not yet archived, or a
file one of its tasks names, the system shall treat the cycle as stale, in
`grove check` and in `grove run`, unless that commit also edits the cycle
file, archives a change the cycle itself lists, or archives a change that
a cycle this one depends on, directly or through other cycles, lists.

#### Scenario: Revised and re-planned in one squash merge

- GIVEN cycle K, marked `unattended`, lists B, and on a branch B's
  proposal is revised and K's `planned` set to the target's commit
- WHEN the branch is squash-merged into the target and the owner runs
  `grove check`
- THEN it reports nothing for K, and `grove run` builds K

#### Scenario: Revised and re-planned in one merge commit

- GIVEN a branch with one commit revising member B and another setting
  cycle K's `planned`
- WHEN the branch is merged into the target with a merge commit
- THEN K is not stale

#### Scenario: Marked after a member changed

- GIVEN member B's proposal was revised in one merge into the target, and
  cycle K was marked `unattended` in a later merge
- WHEN the owner runs `grove check`
- THEN it reports K, saying B changed since the plan and to re-plan K with
  `/grove-plan K`

#### Scenario: Revision and re-plan as two commits on the target

- GIVEN the revision of member B and the edit to cycle K's file reached
  the target as two commits, by rebase or fast-forward
- WHEN the operator runs `grove run K`
- THEN no session is launched, and Grove says K's plan is stale

#### Scenario: A member lands in its own squash merge

- GIVEN cycle K, marked `unattended`, lists A and B, each with a task
  naming `ledger/cli.py`, and a squash merge into the target archives A,
  edits `ledger/cli.py` and leaves K's file as it is
- WHEN the owner runs `grove check`
- THEN it reports nothing for K, `grove run` builds K, and
  `grove instructions implement B` names no file changed since K's plan

#### Scenario: Then a commit archiving nothing

- GIVEN that squash merge, and a later commit on the target that archives
  no change and edits `ledger/cli.py`
- WHEN the owner runs `grove check`
- THEN it reports K, saying B changed since the plan and to re-plan K with
  `/grove-plan K`

#### Scenario: The cycle it follows lands in a squash merge

- GIVEN cycle K, marked `unattended`, depends on cycle L, K's member B
  has a task naming `docs/reference.md`, and a squash merge into the
  target archives L's only change and edits `docs/reference.md`
- WHEN the owner runs `grove check`
- THEN it reports nothing for K, and `grove run` builds K

#### Scenario: Stacked pull requests

- GIVEN cycle K depends on cycle L, whose changes C and D were built as
  stacked pull requests, and the target took them as two commits on its
  first-parent history, each archiving one of them and editing a file K's
  member B names
- WHEN the operator runs `grove run K`
- THEN K's plan is not stale

#### Scenario: Through another cycle

- GIVEN cycle K depends on cycle L, L depends on cycle M, and a merge
  commit into the target archives M's change E and edits a file K's
  member B names
- WHEN the owner runs `grove check`
- THEN K is not stale

#### Scenario: Any other commit still counts

- GIVEN cycle K, marked `unattended`, depends on cycle L, and a commit on
  the target archives no change of L, nor of any cycle L depends on, and
  edits a file K's member B names
- WHEN the owner runs `grove check`
- THEN it reports K, saying B changed since the plan and to re-plan K with
  `/grove-plan K`
