# Grove

A local project workspace for humans and agents, built around a CLI and durable files.

**Status: restart brief captured; new implementation has not started.**

Start with [the restart brief](docs/restart-brief.md). It records the selected
direction, the reasoning from the existing Grove skills and nullsec workflow,
open design questions, and the next useful experiment.

The intended experience combines linked work, questions, research, project
knowledge, and evidence. A CLI serves agents and humans; a TUI can make the
project visible and eventually launch agent sessions. An optional local web UI
can operate on the same records. Core project operations require no hosted
service or persistent daemon.

## Local restart

On 2026-09-18, the first Grove application checkout was preserved intact at
`../grove-archive-2026-09-18/` and this fresh repository replaced `../grove/`.
The archived checkout retains its Git history and ignored local data. Its last
commit was `be40e46`. Existing deployments and external data were not changed.
The archive's instructions and service-based planning authority belong to that
old application; they do not govern this restart.

`../skills/` remains the working Grove skill suite and CLI. `../nullsec/` remains
the real project providing workflow evidence. Neither is migrated by this reset.

## Resume this conversation

Give a new agent this prompt:

> Read AGENTS.md and docs/restart-brief.md. Continue shaping the file-backed
> Grove CLI and interactive workspace described there. The old application was
> archived, and the new implementation has not started. Begin with the first
> end-to-end experience and resolve where project records, worktree ownership,
> and agent run results live. Treat the brief's proposals as proposals. Inspect
> the sibling skills and nullsec projects through their Grove CLI when evidence
> is needed. Preserve this direction and update the brief as choices settle.

