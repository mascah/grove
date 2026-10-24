- The change below is approved as the target branch holds it. Build what it
  says, nothing more.
- Where the line above calls the change a draft, build it as this working
  tree holds it: merging this branch approves it with its code, so edit its
  proposal, change specs and tasks only with the owner's agreement, as for
  revising a change on its branch.
- Work on a branch other than the target, creating one if you are on it.
- Change only the files the tasks name on their `Files:` lines.
- Record each task's status in `implementation.md` with the `Status:` lines
  its template shows, kept exactly. A task it marks done when you start is
  done; build it again only if the check before you finish finds its
  evidence no longer holds.
- Before you finish, run the project's tests, re-read each done task's
  evidence against the code as it now is, and set a task whose evidence no
  longer holds back to `Status: in progress`. Then write the `## Scenarios`
  section of `implementation.md`: one line per scenario of the change
  specs, naming the check that shows it and its result, or `unverified`
  and why.
- Finish: report what was done, what is left, the question blocking each
  blocked task, and each unverified scenario. When every task is done, run
  Grove's `instructions archive <change-id>` and follow what it prints:
  archiving brings the maintained specs up to date on this branch.
