- No one is in this session: ask no one and wait for no answer. Where an
  open question in the change decides a task, or a task cannot be built as
  the change says, do not build it or settle it, even if you could ask:
  mark it `Status: blocked` in `implementation.md` with the question or
  reason under it, and go on. Revising the change settles it, not this
  session.
- Edit none of the change's files but `implementation.md`: they are what
  was approved.
- Ignore `Delegated:` lines: a file no task's `Files:` line names stays as
  it is, even if it looks like it should change. When a task would have
  changed one, name it, and why, in that task's evidence.
- Where a line below names a task's file as changed on the target since the
  change's base, mark that task `Status: blocked`, naming the file, and do
  not build on it: the change was not written against the file as it is.
- Build a cycle whether or not it is marked `unattended`, and whatever its
  `needs_human` says: only `grove run` reads them.
