# Changes to serve

The changes that change 261108-7grt3 proposes to the requirements of
[serve](../../../../specs/serve/spec.md). This is future work.

## MODIFIED Requirements

### Requirement: The project value

When the page asks for `/api/project`, the system shall read the target
branch's documents and answer with the one JSON object `grove changes
--json` prints (changes spec, "Print the project as data"), holding the
root, the target, its full commit, the time of the read, the maintained
specs, every change as `grove changes --json` gives it, each with the
cycle not done that lists it and whether it is approved and waits on
nothing, every cycle, the branches with work, the recorded attempts, the
history and the problems `grove check` would report on the target's
documents, as the requirements below say.

#### Scenario: An approved change in no cycle

- GIVEN the target approves change C, listed by no cycle, with
  `depends_on: [A]` and A archived as finished
- WHEN the page asks for `/api/project`
- THEN C appears with no cycle and as ready

#### Scenario: The same value as the command

- GIVEN a checkout of a Grove project with `grove serve` running
- WHEN the page asks for `/api/project` and the owner runs
  `grove changes --json` in the checkout
- THEN the two objects are equal but for the time of the read
