- No one is in this session: ask no one and wait for no one's answer. Where
  an open question in the change decides a task, or a task cannot be built
  as the change says, do not build it and do not settle it yourself, even if
  you could ask: mark it `Status: blocked` in `implementation.md`, write the
  question or the reason under it, and go on with the other tasks. It is
  settled by revising the change, not here.
- Do not edit the change's files other than `implementation.md`: they are
  what was approved.
- Make the choices and change the paths a task's `Delegated:` line names.
  Beyond them and the `Files:` lines, edit any other file a task needs
  changed to be built as the change says, or makes false, without blocking
  the task, and write `` Also changed: `<path>` — <why> `` under that task
  in `implementation.md`. Never edit `grove.yaml` for it.
- Where a line below names a task's file as changed on the target,
  mark that task `Status: blocked`, naming the file, and do not build on
  it: the change was not written against the file as it is.
- Build a cycle whether or not it is marked `unattended`, and whatever its
  `needs_human` says: only `grove run` reads them.
- If you found something that would apply to future work, your report may
  name one line per change to add to `context`, or under a key of `rules`,
  in `grove.yaml`, with what you found. Do not add it to `grove.yaml`.
