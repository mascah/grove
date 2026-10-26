# Changes to propose

The maintained spec is [grove/specs/propose/spec.md](../../../../specs/propose/spec.md),
as it stands at ea7fbaf. Change 261026-ar7ec, which this change depends
on, also changes it.

## ADDED Requirements

### Requirement: List a revised change's documents

When `grove instructions propose <change-id>` is run for a change the
working tree holds, the system shall name each of the change's documents
by path, with the instruction to read every one before editing, and print
none of their text.

#### Scenario: Revising a large change

- GIVEN drafted change D whose documents hold 31 KB of text
- WHEN `grove instructions propose D` is run
- THEN the output names each of D's documents by path and holds none of
  their text, and still prints the templates of the documents propose
  writes
