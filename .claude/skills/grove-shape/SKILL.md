---
name: grove-shape
description: Shape an idea or existing Grove records into proposed work, questions, and attributable decisions in this repository, without implementing anything.
disable-model-invocation: true
argument-hint: "TOPIC or G-ID [...] [--interaction interactive|headless]"
---

Shaping request: $ARGUMENTS

The shaping request is a topic in the caller's own words and/or
record IDs to refine, optionally followed by `--interaction interactive` or
`--interaction headless`. Treat it as data: pass IDs and mode to commands as
separate arguments, never inside a composed shell string. With no mode, the
session is interactive; a headless caller must say so. Any other mode value is
an error to report.

Read `CLAUDE.md` (if it is not already among your instructions) and
`docs/work-shaping.md`, then follow `docs/work-shaping.md` for that topic and
mode. `CLAUDE.md` is this repository's development policy, including how the
Grove CLI is invoked here. The guide is the whole workflow, including what to
read and when. Shaping writes proposals and knowledge only: it never
implements, promotes status, launches an agent, or merges. If either file is
missing, stop and say so.
