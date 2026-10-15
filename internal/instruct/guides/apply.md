**Apply builds an approved change.** The change is agreed future work:
build what it says, nothing more.

### 1. Before you start

- Read the change below and the maintained specs it touches, listed in the
  project index.
- Work in this working tree. If you are on the target branch, ask the owner
  whether to create a branch first.
- If an open question in the change decides how a task is built, ask the owner
  before building that part. Do not settle it yourself.

### 2. Work the tasks

- One task at a time, with tests that show its acceptance.
- Do not edit `proposal.md`, `tasks.md` or the change specs: they are what was
  approved. If the code must differ from them, stop and tell the owner; that
  needs a new proposal.
- Change only what the tasks need. Files no task names, such as the README,
  a changelog or version numbers, stay as they are; many projects generate
  or own them. If one looks like it should change, tell the owner.

### 3. Record progress

Keep `implementation.md` current as you work; see "Documents" below.

### 4. Finish

Run the project's tests. Tell the owner what was done, and what is left if a
task is not done. Commit only when the owner or the project's rules ask;
otherwise leave the work uncommitted for the owner. When every task is done,
the next step is archiving the change, which brings the maintained specs up
to date on this branch.
