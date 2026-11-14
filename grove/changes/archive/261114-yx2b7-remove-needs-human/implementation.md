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

All tests are in `internal/cli/walkthrough_test.go`; `just check` passed
after both tasks. Where a scenario is what an agent does, the test checks
the instructions it follows; no agent session was run for those, and the
walkthrough's planning steps were not run by hand.

- cycles "Two approved changes": `TestCycles`, plan instructions and template, passed; the agent's run unverified (no session).
- cycles "The agent titles the cycle": `TestCycles`, plan instructions say "write its title yourself from its members, without asking" and name no person, passed; the agent's run unverified (no session).
- cycles "A draft change named": `TestCycles`, plan instructions list only approved changes, passed; the agent's run unverified (no session).
- cycles "A member stale on its base": `TestCycles`, plan instructions bring a stale member up to date, passed; the agent's run unverified (no session).
- cycles "Nothing in the way": `TestCycles`, plan's finish in order, passed; the agent's run unverified (no session).
- cycles "A member waits on a change outside the cycle": `TestCycles`, plan's finish names that obstacle, passed; the agent's run unverified (no session).
- cycles "The cycle waits on a cycle": `TestCycles`, plan's finish names that obstacle, passed; the agent's run unverified (no session).
- cycles "The runner needs something set up": unverified: only the instruction line "Something the runner's checkout must have" shows it, and no test names it.
- cycles "Marked and stale": `TestCycles`, check reports `planned: <a> changed on main since this plan`, passed.
- cycles "Marked and a member's file moved": `TestMovedFiles`, passed.
- cycles "One change in two cycles": `TestCycles`, passed.
- cycles "A needs_human left over": `TestCycles`, check reports line 5 as `unknown field "needs_human"`, passed.
- cycles "A dependency on no cycle": `TestCycles`, passed.
- cycles "Cycles that wait on each other": `TestCycles`, passed.
- runs "A needs_human left over": `TestClaim`, `grove claim --retry <n>` prints n's ID with `needs_human` in its file, passed.
- runs "A stale plan": `TestReplan`, passed.
- runs "A cycle it follows not done": `TestClaim`, passed.
- runs "No such cycle": `TestClaim`, passed.
- serve "A marked cycle that is fresh": `TestServe`, cycle K, passed.
- serve "A member waits outside the cycle": `TestServe`, cycle L's verdict starts `<c> waits on <d>`, passed.
- serve "A member changed since the plan": `TestServe`, cycle M, passed.
- serve "A claimed cycle": `TestServe`, passed.
- serve "Readiness at a glance": `TestServe` checks the value the page draws; how the page shows it unverified (no browser run).
- serve "A refresh while reading": unverified: the page's script, unchanged here but for the chip, was not run in a browser.
- propose "A blocked member answered": `TestCycles`, the revise instructions print the cycle, keep `unattended` and name no person, passed; the agent's run unverified (no session).
- propose "A change in no cycle": `TestCycles`, passed.
- implement "An approved change": `TestImplementInstructions`, passed.
- implement "A draft": `TestImplementInstructions`, passed.
- implement "A change that waits": existing waiting tests in `walkthrough_test.go`, passed.
- implement "A cycle": `TestImplementInstructions`, a member of a cycle not marked gets its plan, and the unattended instructions say only `grove claim` reads `unattended`, passed.
- implement "A member left out": existing test of a member waiting outside its cycle, passed.

## Reconciliation

- cycles: applied to `grove/specs/cycles/spec.md`, with the Purpose edit the change spec names
- runs: applied to `grove/specs/runs/spec.md`, with the Purpose edit the change spec names
- serve: applied to `grove/specs/serve/spec.md`
- propose: applied to `grove/specs/propose/spec.md`
- implement: applied to `grove/specs/implement/spec.md`
