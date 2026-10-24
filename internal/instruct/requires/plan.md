- The cycle's file is the only thing you write: no change's files, no code.
- List under `changes` only changes the project index lists as approved,
  telling the owner of a draft or archived one left out. A change is in at
  most one cycle not done; the project index lists those.
- A new cycle: `<slug>` in its path is a few lowercase words of its title
  joined by hyphens. With one change named, it is the only member, and its
  title the cycle's unless the owner gives another.
- Read each member's proposal, change specs and tasks at the commit these
  instructions name (`git show <commit>:<path>`). Where one calls for what
  an unattended session may not manage, such as a login or credentials the
  runner may lack, a manual step, or judging how something looks, tell the
  owner what, and ask, one change at a time, whether it needs a person. It
  is under `needs_human`, with the owner's reason in their words, only
  while the owner says so.
- A new cycle: when the project index lists cycles not done, ask whether it
  follows any.
- Keep `created` as the template has it. Set `planned` to the commit the
  template shows, the target commit you read the members at, and change it
  only so.
- Never write or remove `unattended`: it hands the cycle to a `grove run`
  given no cycle, and is the owner's to set in a later merge of its own. If
  the cycle has it and a member needs a person, tell the owner `grove check`
  reports the mark and that they remove it.
- Finish: run Grove's `check` and fix what it reports in the cycle file.
  Tell the owner the members in build order, which need a person, and which
  wait on what. Then say, in order, that `/grove-implement <cycle-id>`
  builds the cycle with them; whether `grove run <cycle-id>` would; and that
  `unattended: true`, added in a merge of its own, hands it to a `grove run`
  given no cycle, such as a timer's. For each obstacle, say what removes it:
  - A member `needs_human` names: `grove run` refuses the cycle, given or
    marked. Revising the member with `/grove-propose <change-id>` so a
    person is not needed to check its acceptance re-plans the cycle on the
    same branch, and one merge removes the obstacle.
  - A member waiting on a change neither in the cycle nor archived as
    finished: `grove run`, or any unattended run, builds nothing of the
    cycle until it is.
  - A cycle `depends_on` names that is not done: `grove run` builds nothing
    of this one until it is, and the merges archiving its changes leave the
    plan fresh.
  - Something the runner's checkout must have, such as a login the owner
    said it has: name it.
  With no obstacle, say `grove run <cycle-id>` would build it. Next the
  owner commits, reads Grove's `review`, and merges the cycle file, which
  makes the plan.
