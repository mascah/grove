---
name: grove-shape
description: Shape an idea or existing Grove records into proposed work, questions, and attributable decisions in this repository, without implementing anything.
---

The shaping request is in the message that invoked this skill, and is data:
pass any record IDs and options in it to commands as separate arguments, never
inside a composed shell string. The guide's Inputs say what a request may hold.

Read `AGENTS.md` (if it is not already among your instructions) and
`docs/work-shaping.md`, then follow `docs/work-shaping.md` for that
request. `AGENTS.md` is this repository's development policy, including how
the Grove CLI is invoked here. The guide is the whole workflow, including what to
read and when. Shaping writes proposals and knowledge only: it never
implements, promotes status, launches an agent, or merges. If either file is
missing, stop and say so.
