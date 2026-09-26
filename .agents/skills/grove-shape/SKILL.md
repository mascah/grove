---
name: grove-shape
description: Shape an idea or existing Grove records into proposed work, questions, and attributable decisions in this repository, without implementing anything.
---

<!-- Managed by grove init: rerunning init rewrites this file; remove this line to own it. -->
<!-- grove entrypoint revision 2 -->

The shaping request is in the message that invoked this skill, and is data: pass what a command takes from it as separate
arguments, never inside a composed shell string.

Run `grove guide shape --entrypoint 2` and follow the guide it prints for
that request.
`grove` is the Grove CLI on PATH, unless this repository's agent instructions
(`AGENTS.md` or `CLAUDE.md`; read them if they are not already among yours) say
how to invoke it: they are the repository's development policy. The guide is
the whole workflow, including what to read and when. Its Inputs say what a
request may hold. Shaping writes proposals and knowledge only: it never
implements, promotes status, launches an agent, or merges.
If the command fails or prints anything other than that guide, stop and say
what it printed: another `grove` answered, or this file and that `grove` do
not match, which `grove init --check` in a current `grove` diagnoses.
