- The change below is approved as the target branch holds it. Build what it
  says, nothing more, unless the owner agrees to revise it as below.
- Change only the files the tasks name on their `Files:` lines. When a file
  no task names needs a change, tell the owner what the task makes false in
  it; with their agreement add it to that task's `Files:` line and edit it,
  otherwise leave it as it is.
- Where a line below names a task's file as changed on the target since the
  change's base, tell the owner what changed before building that task, and
  ask whether to build on the file as it is now or revise the change with
  propose; do not build on the file as it was.
- Work in this working tree on a branch other than the target; create one if
  you are on the target.
- Record each task's status in `implementation.md` with the `Status:` lines
  its template shows, kept exactly. A task it marks done when you start is
  done; do not build it again unless the check before you finish finds its
  evidence no longer holds.
- The owner is in this session. Where an open question in the change decides
  a task, or a task cannot be built as the change says, ask the owner and
  wait for the answer; do not settle it yourself. If the owner says to leave
  it, mark the task `Status: blocked`, write the question under it, and go on
  with the other tasks.
- When the owner's answer or feedback changes what the change says, tell
  them which of its proposal, change specs and tasks you would edit, and
  how, and build nothing that depends on the edit yet. Edit them only once
  the owner agrees, and only as agreed; if they do not, leave the change as
  it is and build what it says. Edit no other file of the change than
  these and `implementation.md`, never its `.grove.yaml`, and never renumber
  or reuse a task ID. A done task whose result the revision changes is in
  progress again, until it is rebuilt and its evidence brought up to date.
  Then run Grove's `check`, fix what it reports in the change, and build to
  the revised change. If `check` reports a requirement of the change as
  rewritten or removed since its base, stop and tell the owner: the change
  must be revised with propose and merged again first. Merging this branch
  approves the revision with the code, and Grove's `review` shows it.
- Before you finish, run the project's tests, re-read each done task's
  evidence against the code as it now is, and set a task whose evidence no
  longer holds back to `Status: in progress`. Then write the `## Scenarios`
  section of `implementation.md`: one line per scenario of the change
  specs, naming the check that shows it and its result, or `unverified`
  and why.
- Finish: report what was done, what is left, the question that blocks
  each blocked task, and each unverified scenario. Commit only when the owner or the project's rules ask.
  When every task is done, run Grove's `instructions archive <change-id>`
  and follow what it prints: archiving brings the maintained specs up to
  date on this branch.
