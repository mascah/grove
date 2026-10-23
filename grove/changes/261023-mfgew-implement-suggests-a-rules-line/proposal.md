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

`/grove-implement --unattended`, with nobody in the session, never edits
`grove.yaml` for this; its report may name a suggested line for the owner
to read.

This change is built on change 261023-g0a3h (B1), which folds apply into
`/grove-implement --unattended` and splits implement's requires text into
a shared file and a part for each mode.

The rules are in [specs/implement/spec.md](specs/implement/spec.md).

## Scope

### In scope

- The owner's part of implement's requires text, as B1 leaves it: the
  suggestion, and adding the agreed line to `grove.yaml` on the branch.
- The unattended part: a suggested line only in the report.
- `docs/reference.md` ("grove.yaml": `context` and `rules`),
  `docs/getting-started.md` (step 5) and `dev/walkthrough.md`.

### Out of scope

- Grove reading, checking or writing the suggestion: it is instruction
  text, and `grove check` already checks `grove.yaml`.
- Suggesting a line from explore, propose, plan or archive.
- A change to `grove/specs/config/spec.md`. It covers how Grove reads
  `context` and shows it in every skill's instructions; this change has
  the agent suggest a line at the end of building, and Grove reads
  nothing new, so the rules belong with implement.

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
  its report.** B1 makes unattended `/grove-implement --unattended`,
  which replaces the apply skill.

- **The agreed line is an exception beside the tasks' `Files:` lines,
  not added to a task.** The owner answered "yes" to this on 2026-10-23.
  Reason: the line is about future work, not this change.

## Open questions

None must be answered before review. None blocks a task.

## Related documents

- Requirements: [specs/implement/spec.md](specs/implement/spec.md)
- Tasks: [tasks.md](tasks.md)
