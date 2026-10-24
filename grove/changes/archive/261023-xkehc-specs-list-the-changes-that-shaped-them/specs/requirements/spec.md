# Changes to requirements

No maintained spec for `requirements` at 0c6f419; change 261023-bfwp5,
which this change depends on, proposes it.

## ADDED Requirements

### Requirement: Archive adds a History line

When the agent runs `grove instructions archive <change-id>` for a
finished change, the system shall print instructions to add one line at
the end of the `## History` section of each maintained spec the change
has a change spec for, adding the section at the end of the spec if it
has none: `- [<change-id>](../../changes/archive/<folder>/) <title>`,
where `<folder>` is the change's folder name and `<title>` its title.

#### Scenario: The archive instructions

- GIVEN finished change 261015-k3m9q, "Allow editing older records", with
  a change spec for `record-editing`
- WHEN the agent runs `grove instructions archive 261015-k3m9q`
- THEN the instructions say to add a line to the History of
  `grove/specs/record-editing/spec.md` linking
  `../../changes/archive/261015-k3m9q-allow-editing-older-records/`

#### Scenario: A new maintained spec

- GIVEN no maintained spec for `record-editing`
- WHEN the agent runs `grove instructions archive 261015-k3m9q`
- THEN the maintained spec template it prints ends with `## History`

#### Scenario: A second change

- GIVEN `grove/specs/record-editing/spec.md` has a History line for
  261015-k3m9q, and finished change C has a change spec for
  `record-editing`
- WHEN the agent archives C
- THEN C's line follows 261015-k3m9q's

#### Scenario: A dropped change

- GIVEN change C's `implementation.md` says `Status: dropped`
- WHEN the agent archives C
- THEN no maintained spec gains a line

### Requirement: Archive does not compare History

When the owner runs `grove archive <change-id>` for a finished change,
the system shall leave each maintained spec's `## History` section out of
the fold it compares.

#### Scenario: A line added

- GIVEN finished change C's requirements are in the maintained spec word
  for word, and its History has a line for C
- WHEN the owner runs `grove archive C`
- THEN C moves to `changes/archive/`

#### Scenario: No line added

- GIVEN finished change C's requirements are in the maintained spec word
  for word, and its History has no line for C
- WHEN the owner runs `grove archive C`
- THEN C moves to `changes/archive/`

## Observed current behavior

- The maintained spec template ends with Known discrepancies and Open
  questions and has no History section
  (`internal/instruct/schemas/default/documents/spec/template.md` at
  0c6f419).
- The archive guide has the agent apply ADDED, MODIFIED and REMOVED to
  the maintained spec and list each reconciled spec in
  `implementation.md`; the maintained spec records nothing of the change
  (`internal/instruct/schemas/default/skills/archive.md` at 0c6f419).
- None of Grove's seven maintained specs names a change; the changes that
  shaped them are found only by the `specs/<capability>/` folders of
  `grove/changes/archive/` (`grove/specs/*/spec.md` at 0c6f419).
- A requirement block ends at the next `## ` heading, so a section after
  the requirements is not part of one (`internal/doc/requirements.go` at
  0c6f419).
