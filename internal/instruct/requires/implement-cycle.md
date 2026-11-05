- Work in this working tree on a branch other than the target; create one if
  you are on the target.
- Build each change below in a fresh context of your own, such as a
  subagent, telling it to run Grove's
  `instructions implement <change-id> --unattended` and follow what it
  prints, but neither archive nor commit, and to report each task's status
  and `Check:` lines. Start at once every change whose `after` changes are
  archived and whose tasks' `Files:`, read from its `tasks.md`, name no file
  that a change being built names; start each other change once that holds.
- Only you archive and commit, one returned change at a time: a change
  whose tasks are all done, by Grove's `instructions archive <change-id>`.
  Where two changes built together changed the same file, run the
  project's tests once both have returned; if they fail, archive the one
  started first and build the other again on its own before archiving it.
- With the owner, put the checkpoints each returned change holds to the
  owner at once and build that change again in a fresh context, giving it
  their words; a revision their words need goes through them and Grove's
  `instructions propose <change-id>` here. Unattended, leave checkpoints
  for your report.
- Where you cannot start a fresh context, build each change yourself in
  the order given, one at a time, by the Grove command the list names,
  following what it prints to the end: it archives the change when every
  task is done.
- A change Grove refuses as waiting stays as it is.
- Finish: report the checkpoints of every change first, then, for each
  change, whether it was archived, what is left and the question that
  blocks each blocked task, and each change left out and why.
