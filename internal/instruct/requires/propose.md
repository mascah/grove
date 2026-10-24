- Write the documents below in the change folder their paths name, with its
  `.grove.yaml` from its template; for a new change, `<slug>` is a few
  lowercase words of its title joined by hyphens. Change no code.
- A new change: if these instructions do not follow `bugfix` and the owner
  describes wrong behavior, suggest the `bugfix` schema listed below. Once
  the owner confirms it, or names another listed kind, run Grove's
  `instructions propose --schema <name>` with the owner's description in
  quotes, if any, and follow that instead.
- A new change: ask whether another open change, approved or drafted, must
  land first, naming any in the project index on the same capabilities or
  code. Write each the owner names in `depends_on`, or leave the field out.
  Keep `base` as the template has it: the target commit proposed from.
- A new change: while confirming scope, if the idea is too big for one
  change or has parts the owner could agree separately, ask whether to
  split it. If so, write each part as its own change on this branch, from
  its own run of `instructions propose`, the later naming the earlier in
  `depends_on`.
- A new change: once the owner confirms the scope, and before writing, ask
  one round: every decision you would otherwise assume, such as which
  option, requirements and scenarios, where each task goes in the code,
  what is out of scope and which files outside the code change, numbered,
  each with your recommended answer and why, worded so "yes" accepts it.
  Order it by consequence: first what changes behavior, scope or
  compatibility or loses data, such as an exit status; then architecture or
  interfaces costly to reverse, such as which function holds a check; last,
  on one line of defaults a single "yes" accepts, routine choices such as a
  test's name and place. Leave none out, and wait. With nothing to assume,
  say so and write. Each answer is a decision in the proposal, quoting the
  owner.
- You may recommend delegating small choices, such as test layout, private
  helpers or a new test file. Write a task's `Delegated:` line, naming those
  choices and paths beyond `Files:`, only where the owner agreed.
- A question the owner defers goes into the open questions, naming the task
  it blocks when it decides how that task is built; never guess it. Write
  nothing you inferred as an assumption.
- Revising a drafted change: its files are below. Change what the owner asks
  and leave the rest; if the project index lists it as approved, first tell
  the owner this changes approved work.
- Revising a change in a cycle not done: its file is below, and you re-plan
  it here, so one merge brings both. Ask whether the change needs a person,
  as planning does: it is under `needs_human`, with the owner's reason in
  their words, only while they say so. Set `planned` to the commit the cycle
  template shows and change nothing else; never write or remove
  `unattended`. Tell the owner both must reach the target as one commit, a
  squash or a merge, or the plan is stale.
- Revising a change `check` reports as stale on a requirement: bring that
  MODIFIED or REMOVED requirement up to date with the target's maintained
  spec, keeping what this change alters, then set `base` to the commit the
  `.grove.yaml` template shows. Change `base` only so, never to clear the
  report.
- Write a prototype's findings from this conversation into the proposal,
  citing its branch and commit, and link no file outside the change and the
  target branch.
- Before telling the owner the change is written, read its documents side by
  side: each requirement has an outcome a user or caller can observe; each
  ADDED or MODIFIED one has a task whose acceptance shows it, or needs one;
  nothing contradicts, such as a task naming a file the scope puts out, or a
  scenario its requirement does not cover. Fix what you can and say what you
  changed; where which side is wrong is unsettled, ask, and change nothing it
  decides until answered.
- Finish: run Grove's `check` and fix what it reports in your files. Tell the
  owner the files you wrote, the qualified task IDs (such as
  `<change-id>:T-01`) and the questions still open. Next the owner commits
  and reads Grove's `review`. For a change the project index does not list as
  approved, offer both next steps and choose neither: merge it to approve it,
  or run `/grove-implement <change-id>` on this branch so one merge approves
  it with its code. For an approved one, the next step is merging the
  revision; `/grove-implement` builds it as the target holds it.
