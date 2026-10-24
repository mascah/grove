**Propose a bug fix: "This is wrong; make it right."** It turns a bug into a
change the owner can review: a proposal that shows the bug, tasks that start
with a test reproducing it, and a change spec only when the owner agrees the
corrected behavior.

### 1. Reproduce it

Read the code and the maintained specs the bug touches, and reproduce it,
keeping the exact command and what it printed. If you cannot, ask the owner
how: propose no fix for a bug nobody has seen.

### 2. Find what says what should happen

Tell the owner which of these it is:

- A maintained requirement says what should happen, and the code breaks
  it: cite it in the proposal, and write no change spec.
- No maintained requirement covers it: ask in the round whether the
  corrected behavior should become agreed. If yes, write a change spec that
  ADDs it; if no, write none, and the fix stays observed behavior.
- The maintained spec itself is wrong: that is not a bug but a change to
  agreed behavior, proposed with the `default` schema.

### 3. Write the change

The first task adds a test that reproduces the bug and fails before the
fix. A simple fix is that one task: its `Files:` line names the test and
the file under fix, and its `Interfaces:` line what the fix changes that
callers see. No estimates.
