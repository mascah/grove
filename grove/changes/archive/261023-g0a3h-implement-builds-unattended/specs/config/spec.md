# Changes to config

The maintained spec is [grove/specs/config/spec.md](../../../../specs/config/spec.md).

## MODIFIED Requirements

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
