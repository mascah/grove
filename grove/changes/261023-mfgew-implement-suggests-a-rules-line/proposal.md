# Implement suggests a rules line

## Problem

A correction the owner gives the agent while implementing is lost when
the session ends. At 0c6f419, `context` and `rules` in `grove.yaml` are
how the owner tells every later session what to do (`docs/reference.md`,
"grove.yaml"), but nothing brings a correction there: implement's
instructions end with a report and archive
(`internal/instruct/requires/implement.md`, "Finish"), and change only the
files the tasks name. So the owner corrects the next agent the same way,
or writes the rule by hand from memory.

AI-DLC ends a run by asking "Anything to add for next time?". The owner's
notes park the same idea: "Could a skill propose a rules line at the end
of a run, with the owner deciding?"

## Intended outcome

When `/grove-implement` finishes a change, and the owner corrected the
agent in a way that would apply to future work, the agent suggests one
line for `context` or `rules` in `grove.yaml`, quoting the correction. It
writes the line on the implementation branch only if the owner agrees, so
the merge brings it in with the code. It suggests at most one per change,
and nothing when there was no such correction.

Apply, with nobody in the session, never edits `grove.yaml` for this; its
report may name a suggested line for the owner to read.

The rules are in [specs/implement/spec.md](specs/implement/spec.md) and
[specs/apply/spec.md](specs/apply/spec.md).

## Scope

### In scope

- Implement's requires text: the suggestion, and adding the agreed line
  to `grove.yaml` on the branch.
- Apply's requires text: a suggested line only in the report. Change
  261023-g0a3h (B1) folds apply into `/grove-implement --unattended`, and
  carries this rule into that mode.
- `docs/reference.md` ("grove.yaml": `context` and `rules`),
  `docs/getting-started.md` (step 5) and `dev/walkthrough.md`.

### Out of scope

- Grove reading, checking or writing the suggestion: it is instruction
  text, and `grove check` already checks `grove.yaml`.
- Suggesting a line from explore, propose, plan or archive.
- A change to `grove/specs/config/spec.md`. It covers how Grove reads
  `context` and shows it in every skill's instructions; this change has
  the agent suggest a line at the end of building, and Grove reads
  nothing new, so the rules belong with implement and apply.

## Decisions

The owner agreed to the direction on 2026-10-23 ("yes to all") and to
every answer below in one round ("yes to C & D").

- **At the end of implementing with the owner, if the owner corrected
  the agent in a way that would apply to future work, the agent suggests
  one line for `context` or `rules` in `grove.yaml`, quoting the
  correction, and writes it on this branch only if the owner agrees. At
  most one suggestion per change.** Reason: AI-DLC's learning loop
  ("Anything to add for next time?"), parked in the owner's notes as
  "Could a skill propose a rules line at the end of a run, with the owner
  deciding?".
- **Unattended never edits `grove.yaml`, but may name a suggested line in
  its report.** Today that is the apply skill; B1 makes it
  `/grove-implement --unattended`.

## Open questions

None must be answered before review. Two can wait:

- Blocks T-01: implement's rule "Change only the files the tasks name"
  covers `grove.yaml`. Is the agreed line an exception beside the tasks,
  or added to a task's `Files:` line as a file no task names? Recommended:
  an exception, since the line is about future work, not this change.
- Blocks T-02 only if B1 (261023-g0a3h) is archived first: then apply's
  rule goes to implement's unattended mode instead, in the files B1
  leaves. Recommended: yes.

## Related documents

- Requirements: [specs/implement/spec.md](specs/implement/spec.md),
  [specs/apply/spec.md](specs/apply/spec.md)
- Tasks: [tasks.md](tasks.md)
