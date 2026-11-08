# Changes to config

The changes that change 261108-jpg8n proposes to the requirements of
[config](../../../../specs/config/spec.md). This is future work.

The Purpose's "including the session `grove run` starts" becomes
"including a session started on a claimed cycle".

## MODIFIED Requirements

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
