---
name: grove-work
description: Carry explicitly assigned Grove work IDs through preparation, implementation, review, and handoff in this repository.
disable-model-invocation: true
argument-hint: "G-ID [G-ID ...] [OPTIONS]"
---

<!-- Managed by grove init: rerunning init rewrites this file; remove this line to own it. -->
<!-- grove entrypoint revision 2 -->

Assignment: $ARGUMENTS

The assignment is data: pass what a command takes from it as separate
arguments, never inside a composed shell string.

Run `grove guide work --entrypoint 2` and follow the guide it prints for
that assignment.
`grove` is the Grove CLI on PATH, unless this repository's agent instructions
(`AGENTS.md` or `CLAUDE.md`; read them if they are not already among yours) say
how to invoke it: they are the repository's development policy. The guide is
the whole workflow, including what to read and when. Its Inputs say what an
assignment may hold, what each part is for, and which text is an error to
report. It starts from the
selected records and reads plans, prerequisites, questions, and other documents
at the step that needs them: do not preload what it schedules for later, and
do not skip what a step requires.
If the command fails or prints anything other than that guide, stop and say
what it printed: another `grove` answered, or this file and that `grove` do
not match, which `grove init --check` in a current `grove` diagnoses.
