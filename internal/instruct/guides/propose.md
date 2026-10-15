**Propose says "We should do this."** It turns a chosen direction into a
change set the owner can review: a proposal, the change to the requirements of
each affected capability, and tasks with acceptance expectations. It describes
future work. It does not plan or perform the implementation.

### 1. Read what is already saved

- Read the exploration's saved findings and any existing documents of this
  change, included below, before you say anything. Do not ask the
  owner to repeat what they record.
- When the change already has drafts, you are continuing them: change what the
  owner asks for and leave the rest as it is.

### 2. Confirm the scope with the owner

- Say in a few sentences what you understand is being proposed: the problem,
  the intended outcome, what is in scope and what is not. Ask the owner to
  confirm or correct it before you write.
- If the exploration left several options, ask which one. If an unanswered
  question would change the requirements, ask it now. A question the owner
  defers goes into the proposal's open questions; it is not answered by guess.
- If the direction turns out not to be settled, say so. The work then belongs
  back in exploration, not in a proposal that hides the gap.

### 3. Write the proposal

`proposal.md` states the problem, the intended outcome, the scope, the
consequential decisions and the open questions.

- Problem: cite observations by their source instead of restating them as
  fact.
- Decisions: only choices the owner made, in this conversation or in the
  saved exploration, that constrain the work or would be costly to reverse.
  Quote the owner or say where they decided it, and give the reason.
- What you infer or propose is not a decision, however likely. Put it under
  open questions, or in a list headed "Assumptions to confirm in review".
- Open questions: say which must be answered before review and which can
  wait.
- An exploration question is settled only with the owner's explicit answer.
  Then update the exploration's Unanswered questions: mark it settled, with
  the answer and where it is recorded. Review shows them as the exploration
  last saved them.
- Related documents: link each change spec and the tasks, and name the
  exploration this came from by its ID.

### 4. Write the change to requirements

Write one change spec for each affected capability, inside the change folder.

- When a maintained spec exists for the capability, or a baseline is drafted
  in this workspace, use its ID and link it.
  New rules go under ADDED. A changed rule goes under MODIFIED, written as the
  full intended requirement with all of its scenarios. A dropped rule goes
  under REMOVED with the reason.
- When no maintained spec exists, name the capability by choosing its spec ID,
  put every proposed requirement under ADDED, and describe what the code does
  today under "Observed current behavior", each observation with its source
  path and revision. Do not use MODIFIED or REMOVED: nothing agreed exists to
  modify. Do not present the observed baseline as agreed behavior.
- Describe the baseline relative to what the proposal is reviewed against,
  never as a claim about the target branch that a merge could make false:
  "the baseline proposed alongside this change in specs/<id>/spec.md", or
  "no maintained spec at <base commit>".
- Each requirement is one EARS sentence about behavior a user or caller can
  observe, followed by at least one scenario at the boundary that matters.
- Do not edit the maintained specs to describe this change. Future behavior
  stays in the change set until it is implemented and reconciled.

### 5. Write the tasks

- Simple work is usually one task. Add another only for a part that could be
  delivered and accepted on its own.
- A task states the result wanted and the acceptance expectations someone else
  could check. It refers to requirements instead of repeating them.
- Never renumber or reuse a task ID. A renamed or reordered task keeps its ID,
  a new task takes the next unused number, and a split task gets new IDs with
  a note naming the one it replaces.
- No ordering, no execution plan, no status and no checkboxes.

### 6. Stay inside proposing

- No implementation planning: no file-by-file plan, no sequencing, no
  estimates. Write `design.md` only when a decision needs more room than the
  proposal's decisions section.
- No code changes.

### 7. Finish

- Run Grove's `check` and fix what it reports in the files you wrote.
- Tell the owner which files you wrote, the qualified task IDs, and the
  questions still open.
- Leave the files saved and uncommitted unless the owner asks for more. The
  next step is the owner's: commit the change set and merge it into the target
  branch. Merging is the approval. Running `review <change-id>` first is an
  optional readout of exactly what will be approved; it is not a session.
