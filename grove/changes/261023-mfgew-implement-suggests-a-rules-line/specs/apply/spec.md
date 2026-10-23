# Changes to apply

The maintained spec is [grove/specs/apply/spec.md](../../../../specs/apply/spec.md).

## ADDED Requirements

### Requirement: Name a rules line, never write it

When the agent, building a member of a cycle, found something that would
apply to future work, the apply instructions shall let its report name at
most one suggested line per member for `context`, or for a key under
`rules`, in `grove.yaml`, with what it found, and shall have it add no
such line to `grove.yaml`.

#### Scenario: A test command nobody wrote down

- GIVEN cycle K lists change A, and A's tests run only with a command
  that `context` in `grove.yaml` does not name
- WHEN `/grove-apply K` builds A and finishes
- THEN its report names a `context` line giving that command, nobody is
  asked, and `grove.yaml` on the branch is unchanged
