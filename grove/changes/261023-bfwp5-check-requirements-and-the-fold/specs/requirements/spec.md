# Changes to requirements

## ADDED Requirements

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

## Observed current behavior

- `grove check` reports no requirement for missing "shall" or scenario
  (`internal/cli/check.go` at 57af2e3).
- `grove archive` refuses only a change with a task not `Status: done`
  that is not dropped, then moves the folder
  (`internal/cli/instructions.go`, `archive`, at 57af2e3).
- The archive guide has the agent apply ADDED, MODIFIED and REMOVED to the
  maintained spec and run `grove check`; nothing compares the result with
  the change (`internal/instruct/schemas/default/skills/archive.md` at
  57af2e3).
