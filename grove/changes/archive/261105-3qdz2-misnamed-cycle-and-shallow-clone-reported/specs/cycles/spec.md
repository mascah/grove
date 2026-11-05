# Changes to cycles

The changes that change 261105-3qdz2 proposes to the requirements of
[cycles](../../../../specs/cycles/spec.md), as it stands at 8a8d064.

## ADDED Requirements

### Requirement: Report a cycle file named without an ID

When the owner runs `grove check` and a file under `<root>/cycles/` has a
name that does not start with a cycle ID, the system shall report the
file, saying its name must start with a cycle ID, and exit non-zero, and
`grove changes` shall name no cycle for the changes that file lists.

#### Scenario: A letter IDs leave out

- GIVEN `grove/cycles/261105-proto-unknown-id.yaml` lists approved change
  C, and `o` is a letter cycle IDs leave out
- WHEN the owner runs `grove check` and `grove changes`
- THEN `grove check` reports the file, saying the file name must start
  with a cycle ID, and exits non-zero, and `grove changes` shows C with
  no "in cycle"
