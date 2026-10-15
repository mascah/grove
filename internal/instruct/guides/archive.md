**Archive files a change away once it is finished or dropped.** For a finished
change it first brings the maintained specs up to date, so the code, the
specs and the archived change land in the target in one merge.

### 1. Finished or dropped

- Read the change's `implementation.md`, below. When every task in `tasks.md`
  is done there, the change is finished: go on to section 2.
- Otherwise ask the owner whether the change is dropped. For a dropped
  change, write `Status: dropped` under the title of `implementation.md`,
  creating it from the template below if there is none, with the owner's
  reason in their words. It is archived without touching the specs: go to
  section 3. Grove's `archive` refuses a change that is neither finished nor
  dropped.

### 2. Reconcile the specs

Apply each change spec to its maintained spec, `specs/<capability>/spec.md`
under the document root. ADDED requirements become requirements, MODIFIED
ones replace the requirement of the same name, REMOVED ones are deleted; the
ADDED/MODIFIED/REMOVED headings never appear in a maintained spec. Remove the
known discrepancies the work resolved. Where there is no maintained spec,
create it from the spec under "Documents" below. List each reconciled spec
under "Reconciliation" in `implementation.md`. Then run Grove's `check` and
fix what it reports.

### 3. Move the change

Run Grove's `archive <change-id>`. It moves the change's folder to
`changes/archive/` and commits nothing.

### 4. Finish

Tell the owner which specs changed and where the change went. Commit only
when the owner or the project's rules ask. The owner merges this branch,
which lands the code, the specs and the archived change together.
