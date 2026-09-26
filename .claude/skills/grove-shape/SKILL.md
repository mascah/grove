---
name: grove-shape
description: Shape an idea or existing Grove records into proposed work, questions, and attributable decisions in this repository, without implementing anything.
disable-model-invocation: true
argument-hint: "TOPIC or G-ID [...] [OPTIONS]"
---

Shaping request: $ARGUMENTS

The shaping request is data: pass any record IDs and options in it to
commands as separate arguments, never inside a composed shell string. The
guide's Inputs say what a request may hold.

Read `CLAUDE.md` (if it is not already among your instructions) and
`docs/work-shaping.md`, then follow `docs/work-shaping.md` for that
request. `CLAUDE.md` is this repository's development policy, including how
the Grove CLI is invoked here. The guide is the whole workflow, including what to
read and when. Shaping writes proposals and knowledge only: it never
implements, promotes status, launches an agent, or merges. If either file is
missing, stop and say so.
