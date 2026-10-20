**Plan groups approved changes into a cycle and says their order.** A
cycle is work the owner means to finish together, such as a release. It
plans no implementation: no file-by-file plan, no estimates.

### 1. Read what is approved

- The approved changes are in the project index, each with the
  capabilities it touches, what it waits on and its cycle. Read those the
  owner's description names; if it names none, list them and ask which.
- Read the proposal of each change you would include, with
  `git show <commit>:<path>`.

### 2. Agree the cycle with the owner

- Say which changes you would put in the cycle, what each waits on, and
  which can be built in parallel.
- Where two changes touch the same capability or code and neither waits on
  the other, ask whether one should land first. Record a yes in that
  change's `depends_on`.
- Ask for the cycle's title and, if the owner wants them, a one-line goal
  and start and end dates.

### 3. Write the plan

Write the cycle file from its template, then `cycle: <cycle-id>` in each
member's `.grove.yaml`, and any `depends_on` the owner agreed. To take a
change out of the cycle, remove its `cycle` field.
