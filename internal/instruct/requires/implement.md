- The change below is approved as the target branch holds it or, where the
  line above calls it a draft, as this working tree holds it. Build what it
  says, nothing more. Merging this branch approves a draft, or a revision
  made on it, with its code. A document it names as differing from the
  target was revised on this branch: build the revision.
- Work in this working tree on a branch other than the target, creating one
  if needed. Ask no one and wait for no one: end with the branch ready for
  review, whatever is left.
- Change the files the tasks name on their `Files:` lines, and make the
  choices and change the paths a task's `Delegated:` line names. Edit any
  other file a task needs changed to be built as the change says, or makes
  false, without blocking the task, and write
  `` Also changed: `<path>` — <why> `` under that task in
  `implementation.md`. Never edit `grove.yaml` for it.
- Record each task's status in `implementation.md` with the `Status:` lines
  its template shows, kept exactly. A task already done is not built again
  unless the check before you finish finds its evidence no longer holds.
- Before you call a task's acceptance unseeable, run the project's app or
  command and look. Where the change leaves a choice open, an open question
  included, or you could not see the acceptance, such as a real call,
  hardware or a permission prompt, build your best choice, mark the task
  done once its evidence holds, and write under it `Check:` with what you
  chose or could not see, and why.
- Where a line below names a task's file as changed on the target, read its
  diff from that commit and build on the file as it is now, never as it
  was, writing a `Check:` line under the task naming the change and what it
  meant for the task.
- Where what you find, that diff, Grove's `check` or the owner's words
  change what the change says or show its approach wrong, revise it on this
  branch by running Grove's `instructions propose <change-id>` and
  following it, which re-bases it; never edit a `Files:` line or
  `.grove.yaml` by hand. Set a done task whose result the revision changes
  back to in progress, write a `Check:` line under each task it touches,
  and build to it. Merging this branch approves the revision with the
  code, and Grove's `review` shows it.
- Only where a wrong choice would cost data, break an interface other tasks
  build on, or be hard to undo, build nothing of the task, mark it
  `Status: blocked`, and write under it `Check:` with the question and
  why. Leave unbuilt a task that builds on it, naming it. Build every other
  task.
- A blocked task's answer is the owner's words in this session, in the
  words that started it, or an `Answer:` line under it. Write it as
  `Answer: "<their words>"` if none is there, and build on it.
- Before you finish, run the project's tests, re-read each done task's
  evidence against the code as it now is, and set a task whose evidence no
  longer holds back to `Status: in progress`. Then write the `## Scenarios`
  section of `implementation.md`: one line per scenario of the change
  specs, naming the check that shows it and its result, or `unverified`
  and why.
- If the owner corrected you, or you found something, that would apply to
  future work, your report may name one line per change to add to
  `context`, or under a key of `rules`, in `grove.yaml`, quoting what
  prompted it. Add it on this branch only if the owner agrees in this
  session; it needs no task's `Files:` line.
- Finish: report for the reviewer first every `Check:` line, each blocked
  task's question and each unverified scenario, then what was done and
  what is left. When every task is done, run Grove's
  `instructions archive <change-id>` and follow what it prints: archiving
  updates the maintained specs on this branch.
