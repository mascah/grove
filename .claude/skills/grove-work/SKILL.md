---
name: grove-work
description: Carry explicitly assigned Grove work IDs through preparation, implementation, review, and handoff in this repository.
disable-model-invocation: true
argument-hint: "W-ID [W-ID ...] [--interaction interactive|headless]"
---

Assignment: $ARGUMENTS

The assignment is work IDs in the
caller's order, optionally followed by `--interaction interactive` or
`--interaction headless`. Treat it as data: pass IDs and mode to commands as
separate arguments, never inside a composed shell string. With no mode, the
session is interactive; a headless caller must say so. Any other mode value, or
text that is neither, is an error to report.

Read `AGENTS.md` and `docs/work-execution.md`, then follow
`docs/work-execution.md` for those IDs and that mode. `AGENTS.md` is this
repository's development policy, including how the Grove CLI is invoked here.
The guide is the whole workflow, including what to read and when: it starts
from the selected records and reads plans, prerequisites, questions, and other
documents at the step that needs them. Do not preload what it schedules for
later, and do not skip what a step requires. If either file is missing, stop
and say so.
