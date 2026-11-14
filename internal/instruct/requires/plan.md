- A cycle lists approved changes the owner means to build together. Write
  its file and no code, and edit no change's files except as the next
  line says.
- A new cycle: before writing its file, run Grove's `check`. For each
  member it reports stale on a requirement, or for a named file changed
  since its base, tell the owner this edits approved work, bring the
  member up to date on this branch as Grove's
  `instructions propose <change-id>` says, then go on with the plan.
- List under `changes` only changes the project index lists as approved.
  A change that is a draft, or archived, is not listed: tell the owner. A
  change is in at most one cycle that is not done; the project index lists
  those cycles.
- A new cycle: write its title yourself from its members, without asking.
  With one change named, that change is the only member, and its title is
  the cycle's unless the owner gives another. Replace `<slug>` in its path
  with a few lowercase words of the title joined by hyphens.
- Read each member's proposal, change specs and tasks at the commit these
  instructions name.
- Tell the owner of each member that waits on a change neither in the
  cycle nor archived as finished: `grove claim` will not claim the
  cycle until it is.
- A new cycle: when the project index lists cycles not done, ask the owner
  whether it follows any of them.
- Keep `created` as the template has it. Set `planned` to the commit the
  template shows: the target commit you read the members at. Change it
  only so.
- Never write or remove `unattended`. It hands the cycle to a `grove claim`
  given no cycle and is the owner's to set, in a later merge of its own.
- The plan is made when the owner merges the cycle file into the target
  branch; nothing you do makes it.
- Finish: run Grove's `check` and fix what it reports in the cycle file.
  Tell the owner the members in build order and which wait on what.
  Then, in this order: say `/grove-implement <cycle-id>` builds the cycle
  with them; say whether `grove claim <cycle-id>` would claim it; and
  say `unattended: true`, added to the cycle file in a merge of its own,
  hands it to a `grove claim` given no cycle, such as a timer's. For each
  obstacle, say what removes it:
  - A member waiting on a change neither in the cycle nor archived as
    finished: `grove claim` will not claim the cycle until that change is
    archived as finished.
  - A cycle `depends_on` names that is not done: `grove claim` will
    not claim the cycle until it is, and the merges archiving its
    changes leave the plan fresh.
  - Something the runner's checkout must have, such as a login a member
    calls for: name it.
  With no obstacle, say `grove claim <cycle-id>` would claim the cycle. The
  next step is the owner's: commit, read what Grove's `review` shows, and
  merge it into the target branch.
