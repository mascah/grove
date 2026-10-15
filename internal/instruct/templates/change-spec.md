# Changes to <capability>

<!-- The changes that change {{change}} proposes to the requirements of
`<capability>`. This is future work: nothing here describes the target branch
until the change is implemented and the maintained spec is reconciled.

Requirements have the same shape as in a maintained spec: one EARS sentence
(When / While / If ... then / Where / The system shall), then scenarios.
Delete the sections you do not use. Example:

    ## MODIFIED Requirements

    ### Requirement: Retain previously valid values

    When a user submits an otherwise valid edit that leaves an existing
    record's previously valid category unchanged, the system shall save the
    edit.

    #### Scenario: Update the note on a legacy record

    - GIVEN a record whose category was valid when created but is no longer
      permitted for new records
    - WHEN the user changes only its note and saves
    - THEN the note is saved and the category remains unchanged
-->

## ADDED Requirements

<!-- Rules that do not exist in the maintained spec. If `<capability>` has no
maintained spec yet, every proposed requirement goes here and today's behavior
goes under "Observed current behavior": nothing can be MODIFIED or REMOVED
until a maintained requirement exists.

Say "no maintained spec at {{base}}", not "no maintained spec exists in the
target branch": a later merge could make that false. -->

### Requirement: [short name of the rule]

When [event], the system shall [response].

#### Scenario: [a concrete case at the boundary]

- GIVEN [starting state]
- WHEN [action]
- THEN [observable result]

## MODIFIED Requirements

<!-- Only for requirements already in the maintained spec. Keep the
requirement's name and write the full intended requirement with all of its
scenarios, not just the words that change. -->

## REMOVED Requirements

<!-- Name each maintained requirement being removed and say in a sentence why. -->

## Observed current behavior

<!-- Evidence, not agreement: what the implementation does today where this
change touches it. Give each observation a source (file path and revision,
test, issue). Required when there is no maintained spec to compare against. -->
