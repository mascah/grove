- The owner is in this session. Where an open question in the change decides
  a task, or a task cannot be built as the change says, ask the owner and
  wait for the answer; do not settle it yourself. If the owner says to leave
  it, mark the task `Status: blocked`, write the question under it, and go on
  with the other tasks. Ask about a change the cycle's `needs_human` names
  as about any other, and tell the owner of a change Grove refuses as
  waiting.
- Make the choices and change the paths a task's `Delegated:` line names
  without asking the owner. When a file neither its `Files:` nor its
  `Delegated:` line names needs a change, tell the owner what the task
  makes false in it; with their agreement edit it and write
  `` Also changed: `<path>` — <reason> `` under that task in
  `implementation.md`, otherwise leave it as it is.
- If, while building a task, you find something that makes the change's
  agreed approach wrong, even where the task could still be built as
  written, stop that task and build nothing more of it. Tell the owner
  what you found and the revision the change needs, naming which of its
  proposal, change specs and tasks you would edit and how.
- Where a line below names a task's file as changed on the target, tell
  the owner what changed before building that task, and ask whether to
  build on the file as it is now or revise the change with propose; do not
  build on the file as it was.
- When a task is built, tell the owner how to try it and take their feedback
  before the next.
- When the owner's answer or feedback changes what the change says, tell
  them which of its proposal, change specs and tasks you would edit, and
  how, and build nothing that depends on the edit yet. Only once the owner
  agrees, make the edit as agreed by running Grove's
  `instructions propose <change-id>` and following it, never editing a
  `Files:` line or `.grove.yaml` by hand; if they do not agree, build what
  the change says. A done task whose result the revision changes is in
  progress again until rebuilt and its evidence updated.
  Then run Grove's `check`, fix what it reports in the change, and build to
  the revised change. If `check` reports a requirement of the change as
  rewritten or removed since its base, stop and tell the owner: the change
  must be revised with propose and merged again first. Merging this branch
  approves the revision with the code, and Grove's `review` shows it.
- When you finish a change, if the owner corrected you in this session in a
  way that would apply to future work, suggest one line to add to `context`,
  or under a key of `rules`, in `grove.yaml`, from one correction, quoting
  it: at most one per change, never several merged, and none for a
  correction particular to this change. Add it on this branch only if the
  owner agrees; it needs no task's `Files:` line.
