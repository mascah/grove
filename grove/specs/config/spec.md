# Config

## Purpose

The `context` in a project's `grove.yaml` is instructions the owner gives
every skill, so they reach agents working through Grove, including the
session `grove run` starts, and no other session.

## Boundaries

This covers `context` in `grove.yaml`: how instructions show it, and how
`grove check` and `grove init` treat it. The rest of `grove.yaml` has no
maintained spec yet; `docs/reference.md` explains it. Project facts every session needs belong
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
  the session `grove run` starts does for each member of its cycle
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

## Known discrepancies

None recorded.

## Open questions

None.
