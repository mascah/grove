- A change describes future work. It is approved only when the owner merges
  it into the target branch; nothing you do approves it.
- Write the documents below in the change folder their paths name, with its
  `.grove.yaml` from its template. For a new change, replace `<slug>` with a
  few lowercase words of its title joined by hyphens.
- A new change: if the project's change schemas are listed below and the
  owner says the work is of another kind than these instructions follow,
  confirm it with them, run Grove's `instructions propose --schema <name>`
  with the owner's description in quotes, if there is one, and follow what
  it prints instead.
- A new change: ask the owner whether it needs another open change, approved
  or drafted, to land first. Name any in the project index that touches the
  same capabilities or code. Write each change the owner names in
  `depends_on`; leave the field out when there is none. Keep `base` as the
  template has it: the target commit the change is proposed from.
- A new change: when the project index lists open cycles, ask the owner
  which, if any, the change joins. Write the one named as `cycle`; leave the
  field out otherwise.
- Revising a drafted change: its files are below. Change what the owner asks
  for and leave the rest as it is. If the project index lists it as
  approved, tell the owner that editing it changes approved work before you
  do.
- A question the owner defers is not answered by guess. It goes into the
  change's open questions, naming the task it blocks when it decides how
  that task is built.
- Change no code.
- Finish: run Grove's `check` and fix what it reports in the files you
  wrote. Tell the owner which files you wrote, the qualified task IDs (such
  as `<change-id>:T-01`) and the questions still open. Commit only when the
  owner or the project's rules ask. The next step is the owner's: commit the
  change, read what Grove's `review` shows, and merge it into the target
  branch.
