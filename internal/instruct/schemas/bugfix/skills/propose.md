**Propose a bug fix: "This is wrong; make it right."** It turns a bug into a
change the owner can review: a proposal that shows the bug, tasks that start
with a test reproducing it, and a change spec only when the owner agrees the
corrected behavior. It describes future work. It does not plan or perform
the implementation.

### 1. Reproduce it

- Read the code and the maintained specs in the project index that the bug
  touches.
- Reproduce it, and keep the exact command and what it printed. If you
  cannot, ask the owner how; do not propose a fix for a bug nobody has seen.

### 2. Find what says what should happen

Tell the owner which of these it is:

- A maintained requirement says what should happen, and the code breaks it.
  Cite the requirement in the proposal. Write no change spec: what is agreed
  does not change.
- No maintained requirement covers the behavior. Ask, in the round before
  writing, whether the corrected behavior should become agreed. If yes, write a change spec that
  ADDs it; archiving folds it into the maintained spec. If no, write none,
  and the fix stays observed behavior.
- The maintained spec itself is wrong. That is not a bug: tell the owner it
  changes agreed behavior and belongs to proposing with the `default`
  schema.

### 3. Write the change

- The first task adds a test that reproduces the bug and fails before the
  fix. A simple fix is that one task: its `Files:` line names the test and
  the file under fix, and its `Interfaces:` line what the fix changes that
  callers see.
- No code, commands or steps, no sequencing, no estimates.
