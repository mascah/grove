# Changes to implement

The changes that change 261122-6ndwb proposes to the requirements of
[implement](../../../../specs/implement/spec.md), as it stands at 07d319f.
This is future work: nothing here describes the target branch until the
change is implemented and the maintained spec is reconciled.

## ADDED Requirements

### Requirement: A requirement rewritten on the target

When a change MODIFIES or REMOVES a requirement that `grove check`, run
in the working tree, would report as rewritten or removed since the
change's `base` ([propose](https://github.com/mascah/grove/blob/main/grove/specs/propose/spec.md), "Report a
stale requirement"), the system shall name that capability and
requirement in the implement instructions, saying which of the two
happened and since which commit, and they shall have the agent, before
building any task, revise the change as "Settle what the change does not"
says, so that the requirement holds what the maintained spec now says
together with what this change alters, and write a `Check:` line under a
task that builds the requirement, naming it and what was kept from each
text, asking no one.

#### Scenario: Rewritten by a change that landed first

- GIVEN approved change B, in no cycle, proposed at `base` b143991,
  MODIFIES "Keep a category" of `record-editing` to add "unless the user
  clears it", and change A, which B does not depend on, rewrote that
  requirement to keep the tags too and is archived on the target
- WHEN `/grove-implement B` is run on a branch from the target
- THEN the instructions name `record-editing`'s "Keep a category" as
  rewritten since base b143991; before building any task the agent
  revises B by following `grove instructions propose B`, so that B's
  "Keep a category" keeps the tags and adds "unless the user clears it";
  a `Check:` line under the task that builds it names the requirement and
  what was kept from each text; and once B is archived the maintained
  spec's "Keep a category" still keeps the tags

#### Scenario: Removed by a change that landed first

- GIVEN approved change B MODIFIES "Keep a category", and a change B does
  not depend on removed that requirement from the maintained spec and is
  archived on the target
- WHEN `/grove-implement B` is run
- THEN the instructions name "Keep a category" as removed since B's
  `base`, and the agent revises B before building any task, with a
  `Check:` line saying what it did with the requirement and why

#### Scenario: Nothing rewritten

- GIVEN approved change C MODIFIES "Keep a category", whose text in the
  maintained spec is what it was at C's `base`
- WHEN `/grove-implement C` is run
- THEN the instructions name no requirement as rewritten or removed

#### Scenario: Rewritten by a change it depends on

- GIVEN B holds `depends_on: [A]`, both MODIFY "Keep a category", and the
  target holds A archived as finished
- WHEN `/grove-implement B` is run
- THEN the instructions name no requirement as rewritten

#### Scenario: Two members of one cycle

- GIVEN cycle K lists A and B, neither depending on the other, both
  MODIFY "Keep a category", and A is built and archived on the cycle's
  branch
- WHEN `grove instructions implement --runner B` is run on that branch
- THEN B's instructions name "Keep a category" as rewritten since B's
  `base`, and B is revised on the branch before any of its tasks is built

## Observed current behavior

At 07d319f, `grove instructions implement <change-id>` adds one line to
what it requires when a task's file changed on the target, "Changed on
main since base ...", and none for a requirement
(`internal/cli/instructions.go`, `instructionsChange`). `grove check`
reports the requirement from the working tree (`internal/cli/check.go`,
`stale`), and stays quiet once the fold has put the change's own text in
the maintained spec.
