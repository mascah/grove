- A cycle lists approved changes the owner means to build together. Its
  file is the only thing you write: edit no change's files, and no code.
- List under `changes` only changes the project index lists as approved.
  A change that is a draft, or archived, is not listed: tell the owner. A
  change is in at most one cycle that is not done; the project index lists
  those cycles.
- A new cycle: replace `<slug>` in its path with a few lowercase words of
  its title joined by hyphens. With one change named, that change is the
  only member, and its title is the cycle's unless the owner gives another.
- Read each member's proposal, change specs and tasks at the commit these
  instructions name. Where one calls for something an unattended session
  may not be able to do, such as a login or credentials the runner may
  lack, a manual step, or judging how something looks, tell the owner what
  and ask whether the change needs a person. Name it under `needs_human`,
  with the owner's reason in their words, only when the owner says it does;
  take it out only when they say it no longer does.
- Tell the owner of each member that waits on a change neither in the
  cycle nor archived as finished: an unattended run will not build the
  cycle until it is.
- Keep `created` as the template has it. Set `planned` to the commit the
  template shows: the target commit you read the members at. Change it
  only so.
- Never write or remove `unattended`. It hands the cycle to a `grove run`
  given no cycle and is the owner's to set, in a later merge of its own. If the cycle has
  it and the owner says a member needs a person, tell them `grove check`
  reports the mark and that they remove it.
- The plan is made when the owner merges the cycle file into the target
  branch; nothing you do makes it.
- Finish: run Grove's `check` and fix what it reports in the cycle file.
  Tell the owner the members in build order, which need a person, and which
  wait on what. Then, in this order: say `/grove-implement <cycle-id>`
  builds the cycle with them; say whether `grove run <cycle-id>` would
  build it; and say `unattended: true`, added to the cycle file in a merge
  of its own, hands it to a `grove run` given no cycle, such as a timer's.
  For each obstacle, say what removes it:
  - A member `needs_human` names: `grove run` refuses the cycle, given or
    marked. Revise the member with `/grove-propose <change-id>` so that its
    acceptance can be checked without a person, merge that, and re-plan
    the cycle with `/grove-plan <cycle-id>`.
  - A member waiting on a change neither in the cycle nor archived as
    finished: `grove run` builds nothing of the cycle until that change is
    archived as finished.
  - Something the runner's checkout must have, such as a login the owner
    said the runner has: name it.
  With no obstacle, say `grove run <cycle-id>` would build the cycle.
  Commit only when the owner or the project's rules ask. The
  next step is the owner's: commit, read what Grove's `review` shows, and
  merge it into the target branch.
