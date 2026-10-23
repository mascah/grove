# Changes to config

The changes that change 261023-cnnhm proposes to the requirements of
`config`, the project's `grove.yaml`. There is no maintained spec at
2f8f44d; this adds one for the field it restores, and leaves the rest of
`grove.yaml` to `docs/reference.md`.

## ADDED Requirements

### Requirement: Give context to every skill

Where `grove.yaml` on the target holds a non-empty `context`, the system
shall print it, under a "Project context" heading before what Grove
requires of the skill, in the instructions of every skill.

#### Scenario: Applying a cycle

- GIVEN the target's `grove.yaml` holds `context: Run python3 -m unittest
  before saying a task is done.`
- WHEN `grove instructions apply <cycle-id> <change-id>` runs, as the
  session `grove run` starts does
- THEN its output has a "Project context" section with that line, before
  "What Grove requires of apply"

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

## Observed current behavior

- `configFields` in `internal/cli/check.go` at 2f8f44d omits `context`, so
  `grove check` reports `grove.yaml:<line>: unknown field "context"`.
- `Instructions.Markdown` in `internal/instruct/instruct.go` at 2f8f44d
  prints no project context; e3826d0 removed the "Project context"
  section, `instruct.Project.Context` and the `init` template's comment.
