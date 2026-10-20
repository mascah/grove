- A cycle groups approved changes. Membership is each change's `cycle`
  field and order is each change's `depends_on`, both in its `.grove.yaml`;
  the cycle file holds only the cycle's own details. A change is in at most
  one cycle.
- Write only the cycle file below and the `cycle` and `depends_on` fields
  of approved changes' `.grove.yaml`. Editing an approved change's
  `.grove.yaml` changes approved work: tell the owner before you do. Change
  nothing else in a change, and no code.
- A new cycle: replace `<slug>` in its path with a few lowercase words of
  its title joined by hyphens.
- Tell the owner about each change in the cycle that waits on a change not
  archived as finished that is in another cycle or in none, and ask whether
  it belongs in this cycle too.
- The plan is accepted when the owner merges it into the target branch;
  nothing you do accepts it.
- Finish: run Grove's `check` and fix what it reports in the files you
  wrote. Tell the owner the order: which changes can be applied now, which
  can be built in parallel because neither waits on the other, and which
  wait on which. Commit only when the owner or the project's rules ask. The
  next step is the owner's: commit, read what Grove's `review` shows, and
  merge it into the target branch.
