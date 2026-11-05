- The change below is approved as the target branch holds it. Build what it
  says, nothing more. For a document it names as differing from the
  target, build what its `git show` prints, unless the owner agreed to
  the revision on this branch.
- Where the line above calls the change a draft, the target does not hold
  it: build it as this working tree holds it. Merging this branch into the
  target approves it with its code, so edit its proposal, change specs and
  tasks only with the owner's agreement, as for revising a change on its
  branch.
- Work in this working tree on a branch other than the target, creating one
  if needed.
- Change only the files the tasks name on their `Files:` lines.
- Record each task's status in `implementation.md` with the `Status:` lines
  its template shows, kept exactly. A task already done is not built again
  unless the check before you finish finds its evidence no longer holds.
- A checkpoint is a task whose acceptance no check you can run shows, such
  as how something looks or a manual step, or one resting on a choice the
  change does not settle, an open question included. Build your best guess
  and mark it `Status: to check`, or, where a wrong guess is costly, to
  data, to an interface other tasks build on or to something hard to undo,
  build nothing of it and mark it `Status: blocked`. Under it, write
  `Check:` with what to try or decide, and why. Leave a task that builds
  on a checkpoint unbuilt, naming the checkpoint under it. Build every
  other task.
- A checkpoint's answer is the owner's words in this session, in the words
  that started it, or an `Answer:` line already under the task. Write it
  under the task as `Answer: "<their words>"` if none is there, and build
  on it; a task whose answer confirms its guess is `Status: done` once its
  evidence holds. An answer that changes what the change says is settled
  as below.
- Before you finish, run the project's tests, re-read each done task's
  evidence against the code as it now is, and set a task whose evidence no
  longer holds back to `Status: in progress`. Then write the `## Scenarios`
  section of `implementation.md`: one line per scenario of the change
  specs, naming the check that shows it and its result, or `unverified`
  and why.
- Finish: report the checkpoints first, each with its `Check:` line, then
  what was done, what is left, the question that blocks each blocked
  task, and each unverified scenario. When every task is done,
  run Grove's `instructions archive <change-id>` and follow what it prints:
  archiving updates the maintained specs on this branch.
