**Implement builds an approved change and brings the maintained specs up to
date.** The change is agreed future work: build what it says, nothing more.

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

### 3. Record progress

Keep `implementation.md` current as you work; see "Documents" below.

### 4. Reconcile the specs

When every task is done, apply each change spec to its maintained spec,
`specs/<capability>/spec.md` under the document root. ADDED requirements become
requirements, MODIFIED ones replace the requirement of the same name, REMOVED
ones are deleted; the ADDED/MODIFIED/REMOVED headings never appear in a
maintained spec. Remove the known discrepancies the work resolved. Where there
is no maintained spec, create it from the spec template below. List each reconciled
spec in `implementation.md`. Then run Grove's `check` and fix what it reports.

### 5. Finish

Run the project's tests. Commit the code, `implementation.md` and the
reconciled specs together. Tell the owner what was done and that merging is
theirs.
