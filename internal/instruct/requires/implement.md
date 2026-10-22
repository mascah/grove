- The change below is approved as the target branch holds it. Build what it
  says, nothing more, unless the owner agrees to revise it as below.
- Work in this working tree on a branch other than the target; create one if
  you are on the target.
- Record each task's status in `implementation.md` with the `Status:` lines
  its template shows, kept exactly. A task it already marks done is done; do
  not build it again.
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
  or reuse a task ID. Then run Grove's `check`, fix what it reports in the
  change, and build to the revised change. Merging this branch approves
  the revision with the code, and Grove's `review` shows it.
- Finish: report what was done, what is left, and the question that blocks
  each blocked task. Commit only when the owner or the project's rules ask.
  When every task is done, run Grove's `instructions archive <change-id>`
  and follow what it prints: archiving brings the maintained specs up to
  date on this branch.
