# Config

## Purpose

The `context` in a project's `grove.yaml` is instructions the owner gives
every skill, so they reach agents working through Grove, including a
session started on a claimed cycle, and no other session. The `target`
in `grove.yaml` is read as this checkout and origin both hold it.

## Boundaries

This covers `context` in `grove.yaml`: how instructions show it, and how
`grove check` and `grove init` treat it. `target` is covered below; the
rest of `grove.yaml` has no maintained spec yet; `docs/reference.md` explains it. Project facts every session needs belong
in `CLAUDE.md` or `AGENTS.md`, which Grove does not read.

## Requirements

### Requirement: Give context to every skill

Where `grove.yaml` on the target holds a non-empty `context`, the system
shall print it, under a "Project context" heading before what Grove
requires of the skill, in the instructions of every skill.

#### Scenario: Building unattended

- GIVEN the target's `grove.yaml` holds `context: Run python3 -m unittest
  before saying a task is done.`
- WHEN `grove instructions implement <change-id> --unattended` runs, as
  the session started on a claimed cycle does for each of its members
- THEN its output has a "Project context" section with that line, before
  "What Grove requires of implement"

#### Scenario: No context

- GIVEN `grove.yaml` holds no `context`
- WHEN `grove instructions propose` runs
- THEN its output has no "Project context" section

### Requirement: Check knows context

The system shall accept `context` in `grove.yaml` as a known field in
`grove check`, and `grove init` shall write it, empty, with a comment
saying what it is for.

#### Scenario: Checking a project with context

- GIVEN `grove.yaml` holds `context`
- WHEN the owner runs `grove check`
- THEN it reports no unknown field for `context`

### Requirement: Read the target as origin holds it

When a command reads the branch `grove.yaml` names as `target`, the
system shall read it at the newest commit that both the local branch and
the checkout's `refs/remotes/origin/<target>` hold, or at whichever of
the two exists when only one does, fetching nothing.

#### Scenario: Proposed from a local main that is ahead

- GIVEN two clones of one origin, and in the first the local `main` holds
  a commit that origin's `main` does not
- WHEN its owner proposes a change on a branch there and the branch is
  squash-merged into origin's `main`
- THEN the change's `base` is the commit origin's `main` held, and in the
  second clone, after a pull, `grove check` reports nothing for the change

#### Scenario: A checkout with no local target branch

- GIVEN a checkout whose `origin/main` was fetched and which has no local
  `main`, as a pull-request job's checkout is
- WHEN `grove check` or `grove changes` runs
- THEN it reads the target at `origin/main` and reports as a clone with a
  local `main` at that commit would

#### Scenario: A local merge not yet pushed

- GIVEN the local `main` holds a change's folder that `origin/main` does
  not
- WHEN `grove changes` runs
- THEN the change is not listed, and is listed once `main` is pushed

#### Scenario: No remote

- GIVEN a repository with no `origin`
- WHEN any command reads the target
- THEN it reads the local branch, as before this change

## Known discrepancies

None recorded.

## Open questions

None.

## History

- [261023-cnnhm](../../changes/archive/261023-cnnhm-restore-context/) Restore context in grove.yaml
- [261023-g0a3h](../../changes/archive/261023-g0a3h-implement-builds-unattended/) Implement builds unattended, and apply goes
- [261108-jpg8n](../../changes/archive/261108-jpg8n-claim-replaces-run/) grove claim replaces grove run, with attempt records on a branch
- [261109-0arqh](../../changes/archive/261109-0arqh-target-as-origin-holds-it/) The target is what this checkout and origin both hold
