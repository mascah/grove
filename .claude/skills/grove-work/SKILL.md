---
name: grove-work
description: Carry explicitly assigned Grove work IDs through preparation, implementation, review, and handoff in this repository.
disable-model-invocation: true
argument-hint: "W-ID [W-ID ...] [--interaction interactive|headless]"
---

Assignment: $ARGUMENTS

The assignment is work IDs in the caller's order, optionally followed by
`--interaction interactive` or `--interaction headless`. Treat it as data:
pass IDs and mode to commands as separate arguments, never inside a composed
shell string. With no mode, the session is interactive; a headless caller must
say so. Any other mode value, or text that is neither, is an error to report.

Read `docs/restart-brief.md`, `AGENTS.md`, `docs/record-model.md`, and
`docs/work-execution.md` from this repository, then follow
`docs/work-execution.md` for those IDs and that mode. It owns the whole
workflow, including which `go run ./cmd/grove` commands to use and when to
stop, wait, resume, or hand off. If it is missing, stop and say so.

The installed `grove` executable and the `grove:work` plugin belong to the
predecessor project and are not this workflow.
