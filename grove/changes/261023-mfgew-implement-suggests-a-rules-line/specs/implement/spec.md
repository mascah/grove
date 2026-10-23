# Changes to implement

The maintained spec is [grove/specs/implement/spec.md](../../../../specs/implement/spec.md).
`/grove-implement --unattended` is added to it by change 261023-g0a3h,
which this change depends on.

## ADDED Requirements

### Requirement: Suggest a rules line

When the agent finishes implementing a change with the owner, and the
owner corrected it in that session in a way that would apply to future
work, the implement instructions shall have it suggest one line for
`context`, or for a key under `rules`, in `grove.yaml`, quoting the
correction, and add the line to `grove.yaml` on the implementation branch
only if the owner agrees; it suggests at most one line per change.

#### Scenario: The owner agrees

- GIVEN, while the agent built change C, the owner said "Name each test
  after the requirement it checks"
- WHEN the agent finishes C and the owner agrees to the line it suggests
  under `rules: skills: implement:`, quoting that correction
- THEN `grove.yaml` on the branch holds that line, and Grove's `check`
  reports no problem in it

#### Scenario: The owner declines

- GIVEN the agent suggested a line from the owner's correction
- WHEN the owner declines it
- THEN `grove.yaml` is unchanged

#### Scenario: Nothing for next time

- GIVEN the owner corrected nothing while the agent built C, or only
  something particular to C
- WHEN the agent finishes C
- THEN it suggests no line, and `grove.yaml` is unchanged

#### Scenario: Two corrections

- GIVEN the owner made two corrections while the agent built C, each of
  which would apply to future work
- WHEN the agent finishes C
- THEN it suggests one line, not two

### Requirement: Name a rules line unattended

When the agent, building a change with `/grove-implement <id>
--unattended`, found something that would apply to future work, the
unattended instructions shall let its report name at most one suggested
line per change for `context`, or for a key under `rules`, in
`grove.yaml`, with what it found, and shall have it add no such line to
`grove.yaml`.

#### Scenario: A test command nobody wrote down

- GIVEN cycle K lists change A, and A's tests run only with a command
  that `context` in `grove.yaml` does not name
- WHEN `/grove-implement K --unattended` builds A and finishes
- THEN its report names a `context` line giving that command, nobody is
  asked, and `grove.yaml` on the branch is unchanged
