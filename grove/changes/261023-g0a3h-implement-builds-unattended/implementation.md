# Implementation: Implement builds unattended, and apply goes

<!-- Progress on change 261023-g0a3h, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Implement takes `--unattended`, and apply goes

Status: in progress

- `internal/cli/cli.go`: `instructions implement <id> [--unattended]`; the
  `instructions apply` case and usage rows are gone, so it prints the usage
  and exits 2.
- `internal/cli/instructions.go`: `instructionsApply` deleted;
  `instructionsImplement` takes `unattended` and finds the open cycle the
  target holds that lists a change, so a member gets its plan with the owner
  or not; `instructionsChange` and `instructionsCycle` pass `unattended`
  instead of a skill. Drafts and waiting changes are refused by the same
  code in both modes.
- `internal/instruct/instruct.go`: `Apply` deleted; `Implement` and `Cycle`
  take `unattended`, and append `requires/implement-owner.md` or
  `requires/implement-unattended.md` after the shared `requires/implement.md`
  or `requires/implement-cycle.md`. The cycle form names `instructions
  implement <change-id> --unattended` for each member when unattended.
- `requires/implement.md` and `requires/implement-cycle.md` hold the rules
  both modes share, 261023-m2she's verify rule among them;
  `requires/implement-owner.md` holds asking the owner, revising with their
  agreement, a file no task names, a moved file, and telling the owner how
  to try a task (moved from `skills/implement.md`);
  `requires/implement-unattended.md` holds asking no one, leaving the
  change's files as approved, a file no task names, a moved file blocking
  its task, and building whatever `unattended` and `needs_human` say.
- `requires/apply.md`, `requires/apply-cycle.md`, `skills/apply.md` and the
  apply rows of `skills.yaml` and `schema.yaml` are deleted; implement's
  description in `skills.yaml` names `--unattended`.
- `documents/proposal/guide.md`: `/grove-implement --unattended` leaves a
  task an open question decides blocked, `/grove-implement` asks.
- `internal/cli/check.go`: a `rules.skills.apply` key is reported, saying
  its lines go under `rules.skills.implement`.
- `internal/cli/init.go`: the template comment's example is under
  `implement`, and its skills list has no apply. Init already removes skill
  files Grove no longer writes, so `grove-apply` goes with no new code.
- `grove.yaml`: the apply lines merged into `rules.skills.implement` (the
  pull-request line keeps apply's wording, listing archived changes and
  blocked tasks); `grove init` regenerated `grove-implement` and removed
  `grove-apply` from `.claude/skills` and `.agents/skills`; `grove check`
  prints "No problems in grove/.".
- `internal/cli/walkthrough_test.go`: `TestImplementInstructions` (replacing
  `TestApplyInstructions`) checks init removing the `grove-apply` skill and
  command, `instructions apply` exiting 2, drafts refused in both modes, the
  shared rules and `rules.skills.implement` in both modes, the owner's part
  only without `--unattended` and the unattended part only with it, the
  context before "What Grove requires of implement" unattended, a member's
  plan in both modes for a cycle not marked and with `needs_human`, and the
  cycle form naming `--unattended` for each member. `TestCheck` checks the
  `rules.skills.apply` report; `TestCoordinating` refuses a waiting change
  and builds the second member and a cycle unattended; `TestMovedFiles`
  names the moved file in both modes.
