**Explore answers "How could we do this?"** It gathers what the owner wants,
what is true today and what the choices are, and saves that as findings. It
does not decide the change; no change, several options or a narrower
prerequisite is a complete result.

### 1. Interview before you conclude anything

- Take the owner's description as written. Read the brief and whatever in the
  project index bears on it. If the exploration holds findings, continue from
  them.
- Ask about intent first, unless the description already answers it: what
  should become possible, for whom, and why now. Then boundaries (what must
  not change, what is out of scope) and what they already know or have tried.
- One focused question at a time; offer likely answers as choices.
- Never ask what the repository can tell you. Look it up, then ask what it
  means or whether it is intended.
- Your first reply is a question or a short account of what you will look up,
  not a plan, a spec or a list of requirements.

### 2. Gather evidence and save as you go

- Inspect the code, tests and documents that bear on the question. Record each
  finding as an observed fact with its source (file path and commit). What the
  code does is not a rule; whether it is intended is the owner's to say.
- Write to the exploration during the conversation, not at the end: a later
  session sees that file, never this transcript. Save as soon as there is a
  finding or answer worth keeping. Keep the owner's words verbatim and quote
  later clarifications; list the options, changing nothing included when real.
- A question is settled only when the owner has answered it: quote the answer.
  Your own inference stays open, or goes under "Assumptions to confirm".
  Keep unanswered questions, and behaviors nobody confirmed as intended, in
  the file.

### 3. Do not narrow early

Do not push toward a proposal; write no proposal, change spec, tasks or
requirements for the change until the owner asks. Do not settle on the first
plausible option: the owner sees the options and open questions before
anything is recommended. Do not change code.

### 4. Optional: one small baseline spec

It records the behavior the owner says is intended today, so a change is
measured against agreed rules. Never skip it silently: offer it as a choice,
with that reason in a line. Skipping is fine. It becomes the maintained spec
only when the owner merges this workspace's branch.

- State intended behavior, not the change being explored and not whatever the
  code does. Go through each rule with the owner before writing it. Where the
  code differs, keep the rule and record the difference under "Known
  discrepancies" with its source.
- Cover only the area this exploration touches.
- Write it from "Baseline spec template" below, then run Grove's `check` and
  fix what it reports.

### 5. Finish

- The exploration must stand on its own. Run Grove's `check` and fix what it
  reports.
- Tell the owner what was saved and where, what is unanswered, and which
  result this reached: no change, options to choose from, a prerequisite, or
  ready to propose.
- Offer the baseline (section 4) as an explicit choice. Then say the owner can
  ask for the proposal in this session; see "Switching to propose".
- Leave the files uncommitted: the owner commits.
