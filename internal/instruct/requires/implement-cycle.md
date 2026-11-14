- Work in this working tree on a branch other than the target; create one if
  you are on the target. Ask no one and wait for no one.
- Build each change below in a fresh context of your own, such as a
  subagent, telling it to run Grove's `instructions implement <change-id>`
  and follow what it prints, but neither archive nor commit, and to report
  each task's status, its `Check:` lines and blocked tasks. Start at once
  every change whose `after` changes are archived and whose tasks'
  `Files:`, read from its `tasks.md`, name no file that a change being
  built names; start each other change once that holds.
- Only you archive and commit, one returned change at a time: a change
  whose tasks are all done, by Grove's `instructions archive <change-id>`.
  Where two changes built together changed the same file, run the
  project's tests once both have returned; if they fail, archive the one
  started first and build the other again on its own before archiving it.
- Where you cannot start a fresh context, build each change yourself in
  the order given, one at a time, by the Grove command the list names,
  following what it prints to the end: it archives the change when every
  task is done.
- A change Grove refuses as waiting stays as it is.
- Finish: report for the reviewer first every change's `Check:` lines and
  blocked tasks with their questions, then, for each change, whether it
  was archived and what is left, and each change left out and why.
