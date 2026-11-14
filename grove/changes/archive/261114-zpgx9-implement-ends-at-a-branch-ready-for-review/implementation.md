# Implementation: Implement ends at a branch ready for review

<!-- Progress on change 261114-zpgx9, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — One set of implement instructions, ending at a branch

Status: done

`requires/implement.md` now holds the one set of rules; `implement-owner.md`
and `implement-unattended.md` are deleted; `implement-cycle.md` asks no one
and reports `Check:` lines first. `instruct.Implement` and `instruct.Cycle`
lost their `unattended` argument, `cli.go` the `--unattended` case, and the
implementation guide and template the `to check` status. `just check`:
tests pass; instruction words 8,613 → 7,935. `TestImplement` checks the new
text, the absence of the old ("checkpoint", "to check", "--unattended",
"with the owner, in their own working tree") and that
`instructions implement <id> --unattended` prints the usage with status 2.

Also changed: `.claude/skills/grove-implement/SKILL.md` — `grove init` rewrites the skill description T-01 changes
Also changed: `.agents/skills/grove-implement/SKILL.md` — the same, for Codex
Also changed: `go.mod` — the pre-commit `mod-tidy` hook fails on main; `go mod tidy` marks Lip Gloss indirect, still a dependency through Glamour (rule 16)

Check: the flag is removed, not kept as a no-op: the brief says keep it "only if a documented step still needs it (rule 8)", and once the instructions are one, none does. A timer still passing `--unattended` now gets the usage and exit 2 until the flag is dropped from its prompt (the BREAKING CHANGE footer says so).
Check: the line recording a choice, a moved file or a revision is the existing `Check:` line under a done task, not a new kind of line, so the reviewer finds every judgement in one place and the report lists them first.

### T-02 — Docs, walkthrough and specs say the same

Status: done

`docs/concepts.md`, `coordinating.md`, `getting-started.md`, `index.md`,
`reference.md`, `running.md`, `guides/run-on-a-timer.md`, `dev/walkthrough.md`
and the Purpose and Boundaries of `grove/specs/implement`, `cycles` and
`runs` rewritten. `git grep -- --unattended` outside `CHANGELOG.md`, archived
changes and this change finds only the requirement text archive folds.

Check: `docs/planning.md` and the `needs_human` sentences are left to part C3, its own change; the walkthrough's serve step keeps a hand-written `Status: to check` task, since serve shows any status.

### T-03 — Skill trials expect no stop

Status: done

Trial (d)'s `good` now fails a run that got the scripted answer, and wants
`ledger/report.py` changed and a `Check:` line naming `report` under T-01;
trial (c) has one run; `run.sh` loses the `--unattended` case and the unused
`kept` helper. `bash -n` passes on the three scripts. `just trials` table:
in the pull request.

## Scenarios

- implement "Record each choice for the reviewer", "A screen to look at", "An open question, low cost": TestImplement checks the instruction text; behavior unverified here, trial runs in the pull request
- implement "Block only a costly choice", "A format at stake", "A task that builds on a blocked one": TestImplement checks the text; unverified as behavior, no trial case
- implement "Take a blocked task's answer", "An answer in the words that start a session": TestImplement checks the `Answer:` text; unverified as behavior
- implement "Build a change, a draft or a cycle", "An approved change", "A draft", "A change that waits", "A cycle", "A member left out": TestImplement and TestDependencies, passed
- implement "The flag is gone": TestImplement, `--unattended` gives status 2 and the usage, passed
- implement "Build a cycle's changes in their own contexts", all scenarios: TestImplement checks the cycle text; parallel building unverified, a person runs dev/walkthrough.md step 8
- implement "Name the change's documents", "A large change", "A document that differs": TestImplement, passed
- implement "Settle what the change does not", "The approach rests on something false", "Feedback changes the change": text checked by TestImplement; behavior unverified
- implement "Change only the files the tasks name", all scenarios: trial (c) in the pull request
- implement "A task's file changed on the target", "Changed since base": trial (d) in the pull request; "Its dependency landed": unverified; "A member of a cycle": TestMovedFiles, passed
- implement "Verify before finishing", all scenarios: text checked by TestImplement; trial (e) in the pull request
- implement "Rules for implement", "A rule": TestImplement, passed; "The owner agrees", "No one in the session": unverified, a person runs dev/walkthrough.md "A line for next time"
- config "Give context to every skill", "Building a claimed cycle", "No context": TestImplement and TestPropose, passed
- runs "Start one attempt", "Nothing to apply", all scenarios: TestClaim, passed
- propose "Ask a round before writing", all scenarios: trial (b) in the pull request

## Reconciliation
- implement: applied to specs/implement/spec.md
- config: applied to specs/config/spec.md
- runs: applied to specs/runs/spec.md
- propose: applied to specs/propose/spec.md
