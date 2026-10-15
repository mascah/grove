**Explore is a conversation about what the owner wants and what is true
today.** It is an interview, not a document: write nothing unless the owner
asks. No change, several options or a narrower prerequisite is a complete
result.

### 1. Interview before you conclude anything

- Take the owner's description as written, from the message that asked for
  this. If there is none, ask what they want to explore.
- Read whatever in the code and the project index bears on it.
- Ask about intent first, unless the description already answers it: what
  should become possible, for whom, and why now. Then boundaries (what must
  not change, what is out of scope) and what they already know or have tried.
- One focused question at a time; offer likely answers as choices.
- Never ask what the repository can tell you. Look it up, then ask what it
  means or whether it is intended.
- Your first reply is a question or a short account of what you will look up,
  not a plan, a spec or a list of requirements.

### 2. Keep evidence and agreement apart

- What the code does is observed behavior: cite its source (file path and
  commit). Whether it is intended is the owner's to say.
- Where documents or code disagree, ask which is intended.
- A question is settled only when the owner has answered it. Your own
  inference stays open; say so.

### 3. Do not narrow early

Do not push toward a proposal; write no proposal, change spec, tasks or
requirements until the owner asks. Do not settle on the first plausible
option: the owner sees the options and open questions before anything is
recommended. Do not change code.

### 4. Optional: one small baseline spec

A baseline records the behavior the owner says is intended today, so a change
is measured against agreed rules. When the conversation reaches rules no
maintained spec holds, offer it once as a choice, with that reason in a line.
Write it only if the owner says yes, as the spec under "Documents" below,
covering only the area this conversation touches. It becomes agreed when the
owner merges this branch. Then run Grove's `check` and fix what it reports.

### 5. Finish

- Say what you found, what the owner decided (quoted), what is still open, and
  which result this reached: no change, options to choose from, a
  prerequisite, or ready to propose.
- The owner keeps anything by asking you to write it down, or by proposing
  from this conversation. When they ask for a proposal, run Grove's
  `instructions propose` with no argument and follow what it prints, here.
- Leave anything you wrote uncommitted: the owner commits.
