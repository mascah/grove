# Implementation: Remove needs_human

<!-- Progress on change 261114-yx2b7, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Grove reads, gates and asks nothing about needs_human

Status: done

`doc.Meta.NeedsHuman`, its parsing and `members`' report are gone, so
`needs_human` is an unknown field; `runnable` in `claim.go` no longer
refuses for it; `pageCycle` has no `needs_human` key and the page no chip;
the project index and the cycle template name no `needs_human`. The plan
instructions have the agent title the cycle itself and ask nothing about a
person; propose, re-planning, asks nothing either; the implement
instructions name no `needs_human`. Tests in
`internal/cli/walkthrough_test.go`: `TestCycles` (a leftover
`needs_human` is reported as an unknown field, line 5), `TestClaim`
(`grove claim --retry` claims a cycle whose file still holds it),
`TestServe` (cycle L's verdict names its member waiting outside it),
`TestImplementInstructions` (the implement and cycle instructions). `just
check` passes: 8,358 instruction words, from 8,613.

Also changed: `.claude/skills/grove-plan/SKILL.md` — `grove init` rewrites it from `skills.yaml`'s plan description, which no longer says it asks which changes need a person.
Also changed: `.agents/skills/grove-plan/SKILL.md` — the same, for Codex.

### T-02 — Docs, walkthrough and Grove's own cycle file

Status: not started

## Scenarios

## Reconciliation
