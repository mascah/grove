- Write the documents below in the change folder their paths name, with its
  `.grove.yaml` from its template. For a new change, replace `<slug>` with a
  few lowercase words of its title joined by hyphens.
- A new change: if these instructions do not follow `bugfix` and the owner
  describes behavior that is wrong, suggest the `bugfix` change schema
  listed below. Once the owner confirms it, or says the work is of another
  listed kind than these instructions follow, run Grove's `instructions
  propose --schema <name>` with the owner's description in quotes, if there
  is one, and follow what it prints instead.
- A new change: ask the owner whether it needs another open change, approved
  or drafted, to land first. Name any in the project index that touches the
  same capabilities or code. Write each change the owner names in
  `depends_on`; leave the field out when there is none. Keep `base` as the
  template has it: the target commit the change is proposed from.
- Revising a drafted change: its files are below. Change what the owner asks
  for and leave the rest as it is. If the project index lists it as
  approved, tell the owner that editing it changes approved work before you
  do.
- Revising a change that a cycle not done lists: the cycle's file is below,
  and you re-plan it on this branch, so one merge brings the revision and
  the plan. Ask the owner whether the change needs a person, as planning
  does: name it under `needs_human`, with their reason in their words, only
  when they say it does, and take it out only when they say it no longer
  does. Set `planned` to the commit the cycle template below shows, and
  change nothing else in the file; never write or remove `unattended`. Tell
  the owner the revision and the re-plan must reach the target as one
  commit, a squash or a merge, or the plan is stale.
- Revising a change that Grove's `check` reports as stale on a
  requirement, or for a named file changed since its base: bring that
  MODIFIED or REMOVED requirement, or the tasks naming that file, up to
  date with the target, keeping what this change alters, then set `base`
  to the commit the `.grove.yaml` template below shows. Change `base` only
  so; never clear the report by hand.
- A new change: while confirming its scope, if the idea is too big for one
  change or has parts the owner could agree separately, say so and ask
  whether to split it. If the owner agrees, write each part as its own
  change on this branch, each from its own run of Grove's `instructions
  propose`, the later naming the earlier in `depends_on`; otherwise write
  one.
- A new change: once the owner has confirmed its scope, and before you
  write any document, ask one round: every decision you would otherwise
  assume, numbered, each with your recommended answer worded so that "yes"
  accepts it. Order it by consequence: first what changes behavior, scope
  or compatibility or loses data, then architecture or interfaces costly
  to reverse, and last the routine choices inside the agreed approach, as
  recommended defaults on one line, not a list, that a single "yes"
  accepts. Leave none out. Wait for the owner. If there is nothing to assume, say so and
  write. Each answer is a decision in the proposal, quoting the owner.
- Write a task's `Delegated:` line, naming choices the implementer may make
  without asking and paths it may change beyond `Files:`, only where the
  owner agreed to it in the round; otherwise write none.
- A question the owner defers is not answered by guess. It goes into the
  change's open questions, naming the task it blocks when it decides how
  that task is built. Write nothing you inferred as an assumption.
- When this conversation holds a prototype's findings, write them into the
  proposal, citing the prototype branch and the commit that holds the
  experiment. Link no file outside the change and the target branch: the
  prototype branch is never merged.
- Change no code.
- Before you tell the owner the change is written, check that each
  requirement has an outcome a user or caller can observe, each ADDED or
  MODIFIED requirement is covered by a task's acceptance, the proposal's
  scope, the change specs and the tasks do not contradict each other, and
  no task does what the scope leaves out. Fix what you can and tell the
  owner what you changed; ask the owner about what needs a decision, and
  change nothing it decides until they answer.
- Finish: run Grove's `check` and fix what it reports in the files you
  wrote. Tell the owner which files you wrote, the qualified task IDs (such
  as `<change-id>:T-01`) and the questions still open. The next step is the
  owner's: commit the change and read what Grove's `review` shows. For a
  change the project index does not list as approved, offer both next steps
  and choose neither for them: merge it into the target branch to approve
  it, or run `/grove-implement <change-id>` on this branch so one merge
  approves it with its code. For a change it lists as approved, the next
  step is merging the revision; `/grove-implement` builds it as the target
  holds it.
