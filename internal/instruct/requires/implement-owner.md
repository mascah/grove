- The owner is in this session. Where an open question in the change decides
  a task, or a task cannot be built as the change says, ask the owner and
  wait; do not settle it yourself. If the owner says to leave it, mark the
  task `Status: blocked` with the question under it, and go on. Ask about a
  change the cycle's `needs_human` names as about any other, and tell the
  owner of a change Grove refuses as waiting.
- Make the choices and change the paths a task's `Delegated:` line names
  without asking. When a file neither its `Files:` nor its `Delegated:` line
  names needs a change, tell the owner what the task makes false in it;
  with their agreement add it to that task's `Files:` line and edit it,
  otherwise leave it.
- If, while building a task, you find something that makes the change's
  agreed approach wrong, even where the task could still be built as
  written, stop and build nothing more of it. Tell the owner what you found
  and the revision the change needs: which of its proposal, change specs
  and tasks you would edit, and how.
- Where a line below names a task's file as changed on the target since the
  change's base, tell the owner what changed before building that task, and
  ask whether to build on the file as it is now or revise the change with
  propose; never build on the file as it was.
- When a task is built, tell the owner how to try it and take their feedback
  before the next.
- When the owner's answer or feedback changes what the change says, tell
  them which of its proposal, change specs and tasks you would edit, and
  how, and build nothing depending on it yet. Edit them only once the owner
  agrees, and only as agreed; otherwise build what the change says. Edit no
  other file of the change than these and `implementation.md`, never its
  `.grove.yaml`, and never renumber or reuse a task ID. A done task whose
  result the revision changes is in progress again until rebuilt with its
  evidence brought up to date. Then run Grove's `check`, fix what it reports
  in the change, and build to the revision. If `check` reports a
  requirement of the change as rewritten or removed since its base, stop and
  tell the owner: the change must be revised with propose and merged again
  first. Merging this branch approves the revision with the code, as Grove's
  `review` shows.
