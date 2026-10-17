**Propose says "We should do this."** It turns a chosen direction into a
change the owner can review: a proposal, the change to the requirements of
each affected capability, and tasks with acceptance expectations. It
describes future work. It does not plan or perform the implementation.

### 1. Start from what is already said

- From this conversation: use what the owner has said here. Do not ask them
  to repeat it.
- From the owner's description below: read the code and the documents in the
  project index that it touches before you say anything.
- Revising a drafted change: its files are below. Change what the owner asks
  for and leave the rest as it is. If the project index lists it as approved,
  tell the owner that editing it changes approved work before you do.

### 2. Confirm the scope with the owner

- Say in a few sentences what you understand is being proposed: the problem,
  the intended outcome, what is in scope and what is not. Ask the owner to
  confirm or correct it before you write.
- If several options are still open, ask which one. If an unanswered question
  would change the requirements or how a task is built, ask it now. A
  question the owner defers goes into the proposal's open questions, naming
  the task it blocks if it decides how that task is built; it is not answered
  by guess.
- If the direction turns out not to be settled, say so. The work then belongs
  back in exploring, not in a proposal that hides the gap.

### 3. Write the change

Write the documents below in the change folder their paths name. For a new
change, replace `<slug>` with a few lowercase words of its title joined by
hyphens. Write `design.md` only when a decision needs more room than the
proposal's decisions section.

### 4. Stay inside proposing

No implementation planning: no file-by-file plan, no sequencing, no
estimates. No code changes.

### 5. Finish

- Run Grove's `check` and fix what it reports in the files you wrote.
- Tell the owner which files you wrote, the qualified task IDs (such as
  `<change-id>:T-01`), and the questions still open.
- Commit only when the owner or the project's rules ask; otherwise leave the
  files uncommitted. The next step is the owner's: commit the change, read what Grove's `review` shows, and merge it into the target
  branch. Merging is the approval.
