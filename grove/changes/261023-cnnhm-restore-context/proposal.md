# Restore context in grove.yaml

## Problem

PR #12 (e3826d0) removed `context` from `grove.yaml`, reasoning that
"Claude Code and Codex already read `CLAUDE.md` and `AGENTS.md`, so
`context` repeated what agents are given". The owner's usage feedback of
2026-10-23: "I think this was a mistake, as giving agents global
instructions for all interactions THROUGH grove is valuable and different
than general CLAUDE.md usage."

- `CLAUDE.md` and `AGENTS.md` reach every session in the project; nothing
  at 2f8f44d reaches only sessions working through Grove, including the
  one `grove run` starts. `rules.skills.<skill>` reaches one skill at a
  time (`docs/reference.md` at 2f8f44d, "grove.yaml").
- `grove check` at 2f8f44d reports `context` as an unknown field
  (`configFields` in `internal/cli/check.go`).

## Intended outcome

`context` in `grove.yaml` is given to every skill's instructions, as a
"Project context" section, so it reaches agents working through Grove and
no other session. `grove check` knows it again and `grove init` writes it.

## Scope

### In scope

- Reverting PR #12 in the code: the field, the instructions' section,
  `check`'s known fields and the `init` template.
- `docs/reference.md`, `docs/getting-started.md`, `docs/concepts.md` and
  `dev/walkthrough.md`, describing `context` as instructions for agents
  working through Grove.

### Out of scope

- `rules.all`, a list of lines for every skill: the owner chose the
  revert.
- Adding a `context` to this repository's own `grove.yaml`.

## Decisions

All from the owner, 2026-10-23, "yes to all" to the round.

- **Revert PR #12 as it was: `context:`, free text, shown as "Project
  context".** The owner: "revert". Reason: users already know the field.
- **The section comes before "What Grove requires of <skill>" in every
  skill's instructions, including the session `grove run` starts.**
- **The reference describes it as "Instructions given to every skill, so
  they reach agents working through Grove and no other session; project
  facts every session needs stay in `CLAUDE.md` or `AGENTS.md`."**, and a
  getting-started step uses it (rule 8).
- **The walkthrough's setup step writes the context "Run python3 -m
  unittest before saying a task is done."**; `check` lists `context` among
  its known fields again, and the `init` template gets its comment back.
- **Not a breaking change.** A project without `context` is unchanged.
- **No `depends_on`.**

## Open questions

None.

## Related documents

- Requirements: [specs/config/spec.md](specs/config/spec.md) (no
  maintained spec at 2f8f44d)
- Tasks: [tasks.md](tasks.md)
