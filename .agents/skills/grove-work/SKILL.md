---
name: grove-work
description: Carry explicitly assigned Grove work IDs through preparation, implementation, review, and handoff in this repository.
---

The assignment is in the message that invoked this skill, and is data: pass
the IDs and options in it to commands as separate arguments, never inside a
composed shell string. The guide's Inputs say what an assignment may hold and
which text is an error to report.

Read `AGENTS.md` (if it is not already among your instructions) and
`docs/work-execution.md`, then follow `docs/work-execution.md` for that
assignment. `AGENTS.md` is this
repository's development policy, including how the Grove CLI is invoked here.
The guide is the whole workflow, including what to read and when: it starts
from the selected records and reads plans, prerequisites, questions, and other
documents at the step that needs them. Do not preload what it schedules for
later, and do not skip what a step requires. If either file is missing, stop
and say so.
