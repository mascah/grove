# Requirements

## Purpose

`grove check` wants "shall" and a scenario in each requirement and
reports a change that grows a rule or a spec past its limit, and `grove
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

### Requirement: Report a new rule over the limit

When the owner runs `grove check`, the system shall report each
requirement an open change ADDs whose rule, its text before its first
scenario, is over 60 words as `wc -w` counts them, naming its file, line,
name and word count.

#### Scenario: A rule of 61 words

- GIVEN an open change ADDs requirement "Export a ledger" with a rule of
  61 words and one scenario
- WHEN the owner runs `grove check`
- THEN it reports the change spec at the requirement's line, saying
  "Export a ledger" has a rule of 61 words against a limit of 60, and
  exits non-zero

#### Scenario: A rule of 60 words with long scenarios

- GIVEN an open change ADDs a requirement with a rule of 60 words and
  scenarios of 400 words
- WHEN the owner runs `grove check`
- THEN the requirement is not reported

#### Scenario: A long rule no open change names

- GIVEN `grove/specs/record-editing/spec.md` holds "Keep a category" with
  a rule of 90 words, and no open change has a block for it
- WHEN the owner runs `grove check`
- THEN "Keep a category" is not reported

### Requirement: Report a modified rule that grew

When the owner runs `grove check`, the system shall report each
requirement an open change MODIFIES whose rule is over 60 words and
longer than the rule the change's `base` holds for it, naming its file,
line, name and both word counts.

#### Scenario: A long rule made longer

- GIVEN the maintained spec at an open change's `base` holds "Keep a
  category" with a rule of 80 words, and the change MODIFIES it with a
  rule of 95
- WHEN the owner runs `grove check`
- THEN it reports the change spec at the requirement's line, saying the
  rule of "Keep a category" grew from 80 words to 95 against a limit of
  60

#### Scenario: A long rule made shorter

- GIVEN the rule is 80 words at `base` and 70 in the change
- WHEN the owner runs `grove check`
- THEN the requirement is not reported

#### Scenario: A short rule taken over the limit

- GIVEN the rule is 50 words at `base` and 65 in the change
- WHEN the owner runs `grove check`
- THEN the requirement is reported

#### Scenario: Only the scenarios change

- GIVEN the rule is 80 words at `base`, and the change MODIFIES the
  requirement with the same rule and one scenario more
- WHEN the owner runs `grove check`
- THEN the requirement is not reported

### Requirement: Report a spec a change grows past the limit

When the owner runs `grove check`, the system shall report each open
change whose requirements, folded into a maintained spec as its `base`
holds it, leave that spec over 2,500 words and larger than at `base`,
naming the change spec, the capability and both word counts.

#### Scenario: Adding to a spec already over

- GIVEN `grove/specs/record-editing/spec.md` holds 3,000 words at an open
  change's `base`, and the change ADDs a requirement of 200 words to it
- WHEN the owner runs `grove check`
- THEN it reports the change's spec for `record-editing`, saying the
  change takes `record-editing` from 3,000 words to 3,200 against a limit
  of 2,500

#### Scenario: Adding and removing

- GIVEN that spec holds 3,000 words at `base`, and the change ADDs a
  requirement of 200 words and REMOVEs one of 300
- WHEN the owner runs `grove check`
- THEN the change is not reported

#### Scenario: Taking a spec over

- GIVEN that spec holds 2,400 words at `base`, and the change ADDs a
  requirement of 200 words
- WHEN the owner runs `grove check`
- THEN the change is reported, from 2,400 words to 2,600

#### Scenario: A spec that stays under

- GIVEN that spec holds 1,000 words at `base`, and the change ADDs a
  requirement of 200 words
- WHEN the owner runs `grove check`
- THEN the change is not reported

#### Scenario: A new capability

- GIVEN no maintained spec for `exports` at an open change's `base`, and
  the change ADDs requirements of 2,600 words to it
- WHEN the owner runs `grove check`
- THEN the change is reported, from 0 words to 2,600

### Requirement: Propose states the limits

When the agent runs `grove instructions propose`, the system shall print,
in the guide to a change spec, that a rule keeps to 60 words, that a
change leaves no spec over 2,500 words larger than it found it, and that
`grove check` reports both.

#### Scenario: The propose instructions

- GIVEN any project
- WHEN the agent runs `grove instructions propose`
- THEN the change spec's guide gives both limits and says `grove check`
  reports a change that exceeds either

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
- [261127-hvq1r](../../changes/archive/261127-hvq1r-check-limits-the-size-of-rules-and-specs/) Check limits the size of rules and specs
