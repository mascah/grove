# Requirements

## Purpose

`grove check` wants "shall" and a scenario in each requirement, and `grove
archive` refuses a fold that does not hold the change's requirements word
for word. A requirement is one `### Requirement: <name>` block of a
maintained spec or a change spec: a rule that says "shall", with at least
one scenario. Archive folds a change's ADDED, MODIFIED and REMOVED
requirements into the maintained specs; "where this branch left the
target" is the merge base of the current branch and the target.

## Boundaries

This covers how `grove check` reads the shape of requirements and how
`grove archive` compares a fold with what was agreed. Sections of a
maintained spec other than its requirements, such as Purpose and Known
discrepancies, are not compared. Whether the code meets the scenarios is
not checked here. The agent, not Grove, does the fold.

## Requirements

### Requirement: Each requirement says shall and has a scenario

When the owner runs `grove check`, the system shall report each
requirement not under REMOVED, in a maintained spec or in a change spec of
an open change, whose text has no "shall" or no `#### Scenario:` heading,
naming its file, line and name.

#### Scenario: A requirement with no scenario

- GIVEN `grove/specs/record-editing/spec.md` holds requirement "Keep a
  category", which says "the system shall" and has no `#### Scenario:`
- WHEN the owner runs `grove check`
- THEN it reports that file at the requirement's line, saying "Keep a
  category" has no scenario, and exits non-zero

#### Scenario: A requirement with no shall

- GIVEN an open change's change spec ADDs a requirement that reads "Users
  can delete a record" and has a scenario
- WHEN the owner runs `grove check`
- THEN it reports the requirement, saying it has no "shall"

#### Scenario: A removed requirement

- GIVEN a change spec REMOVEs a requirement with a sentence saying why,
  and no scenario
- WHEN the owner runs `grove check`
- THEN the requirement is not reported

#### Scenario: An archived change

- GIVEN an archived change's change spec holds a requirement with no
  scenario
- WHEN the owner runs `grove check`
- THEN it is not reported

### Requirement: No change headings in a maintained spec

When the owner runs `grove check`, the system shall report each `## ADDED
Requirements`, `## MODIFIED Requirements` or `## REMOVED Requirements`
heading in a maintained spec, naming its file and line.

#### Scenario: A heading left by a fold

- GIVEN `grove/specs/record-editing/spec.md` holds `## ADDED Requirements`
  above a requirement
- WHEN the owner runs `grove check`
- THEN it reports that line, saying the heading belongs in a change spec

### Requirement: Archive refuses a fold that does not match

When the owner runs `grove archive <change-id>` for a finished change, the
system shall compare each maintained spec the change has a change spec
for with what it should hold: its requirements where this branch left the
target, then what each change archived on this branch since then gives
them, then what this change gives them; and if a requirement is missing,
is present though removed, or has other text, whitespace aside, it shall
print each such requirement as `path:line: message`, move nothing, and
exit non-zero.

#### Scenario: Requirements copied word for word

- GIVEN finished change C MODIFIES "Keep a category" and ADDs "Delete a
  record" in `record-editing`, and the maintained spec holds both as C
  gives them
- WHEN the owner runs `grove archive C`
- THEN C moves to `changes/archive/`

#### Scenario: A modified requirement reworded

- GIVEN the maintained spec holds "Keep a category" with one word other
  than C gives it
- WHEN the owner runs `grove archive C`
- THEN it reports the maintained spec at that requirement, saying it is
  not as C gives it, and C stays where it is

#### Scenario: A removed requirement left in

- GIVEN C REMOVEs "Old rule" and the maintained spec still holds it
- WHEN the owner runs `grove archive C`
- THEN it reports "Old rule" as still there though C removes it, and C
  stays where it is

#### Scenario: An unrelated requirement edited

- GIVEN the maintained spec's "Show a record", which no change archived on
  this branch names, differs from its text where the branch left the
  target
- WHEN the owner runs `grove archive C`
- THEN it reports "Show a record" as changed though no change names it,
  and C stays where it is

#### Scenario: Lines reflowed

- GIVEN the maintained spec holds C's requirements with the same words
  wrapped at other places
- WHEN the owner runs `grove archive C`
- THEN C moves to `changes/archive/`

#### Scenario: The second change of a cycle

- GIVEN change A, archived earlier on this branch, MODIFIED "Keep a
  category", and C ADDs "Delete a record" in the same capability
- WHEN the owner runs `grove archive C` with the spec holding both
- THEN C moves to `changes/archive/`

#### Scenario: The target moved

- GIVEN the target rewrote "Show a record" after this branch left it
- WHEN the owner runs `grove archive C` on the branch, whose spec has the
  text from where it left
- THEN "Show a record" is not reported

#### Scenario: A dropped change

- GIVEN C's `implementation.md` says `Status: dropped`
- WHEN the owner runs `grove archive C`
- THEN no spec is compared, and C moves to `changes/archive/`

### Requirement: Archive copies requirements word for word

When the agent runs `grove instructions archive <change-id>`, the system
shall print instructions to copy each ADDED and MODIFIED requirement into
its maintained spec word for word, to delete each REMOVED one, and, if
`grove archive` refuses, to fix the spec as its messages say and run it
again.

#### Scenario: The archive instructions

- GIVEN finished change C
- WHEN the agent runs `grove instructions archive C`
- THEN the instructions say to copy requirements word for word and to run
  `grove archive` again after fixing what it reports

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

## Known discrepancies

None recorded.

## Open questions

None.

## History

- [261023-bfwp5](../../changes/archive/261023-bfwp5-check-requirements-and-the-fold/) Check requirements and the archive fold
- [261023-xkehc](../../changes/archive/261023-xkehc-specs-list-the-changes-that-shaped-them/) Maintained specs list the changes that shaped them
