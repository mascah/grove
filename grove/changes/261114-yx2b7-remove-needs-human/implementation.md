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

Status: done

`docs/concepts.md` loses "Needs a person"; `docs/planning.md` has the
agent title the cycle itself, ask nothing about a person, and its example
cycle and verdicts lose `needs_human`; `docs/running.md`,
`docs/watching.md` and `docs/coordinating.md` lose it from claim's tests,
claim's refusals, readiness and `--unattended`; `docs/reference.md` lists
no `needs_human` field and says `grove check` reports a leftover one as an
unknown field, to delete. `dev/walkthrough.md`: planning steps 1, 3 and 4
(a leftover `needs_human` reported as an unknown field), coordinating
step 5, running steps 2 and 6, and watching step 2 (L's member waits on
D). Grove's own done cycle `261023-9kfxn` loses its `needs_human`. `git
grep` finds `needs_human` outside the archive and this change only in
`docs/reference.md`, the walkthrough's leftover-field step, the tests of
it, and the maintained specs this change's archive rewrites. `grove
check`: No problems in grove/.

## Scenarios

## Reconciliation
