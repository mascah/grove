**Plan says "Build these together."** It puts approved changes in a cycle
and finds out, with the owner, which of them need a person. It plans no
implementation: no file-by-file plan, no estimates.

### 1. Read what is approved

- The approved changes are in the project index, each with the
  capabilities it touches and what it waits on. Read those the owner's
  description or this conversation names; if none is named, list them and
  ask which.
- Read each one with `git show <commit>:<path>`.

### 2. Agree the cycle with the owner

- Say which changes you would put in the cycle, in build order, and what
  each waits on.
- Ask about each change that may need a person, one at a time, saying what
  in it does.
- Ask for the cycle's title and, if the owner wants one, a one-line goal.

### 3. Write the cycle file

Write it from its template. Revising, change what the owner asks for and
leave the rest as it is.
