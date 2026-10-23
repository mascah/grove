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

### T-02 — `grove run` launches implement unattended

Status: in progress

- `internal/cli/run.go`: the session's prompt is `/grove-implement
  <cycle-id> --unattended` where Claude Code is listed as skills and
  `/grove:implement <cycle-id> --unattended` as commands; comments name the
  implement skill. "Nothing to apply" stays, as the runs spec's "Nothing to
  apply" says.
- `internal/cli/cli.go`: `run`'s usage line says it implements the cycle
  asking no one, naming no apply skill.
- `internal/cli/walkthrough_test.go`: `TestRun` checks the stub `claude`
  given `-p /grove-implement <k> --unattended` for skills, and, after the
  owner lists Claude Code as commands, `-p /grove:implement <m>
  --unattended`.

### T-03 — Document implement unattended

Status: in progress

- `docs/guides/run-apply-on-a-timer.md` moved to
  `docs/guides/run-on-a-timer.md`, titled "Run on a timer"; `mkdocs.yml`,
  `docs/index.md` and `docs/running.md` link the new path, and no page or
  nav entry names the old one. `docs/running.md` is "Running unattended",
  and the links to it in `README.md`, `docs/index.md`, `docs/planning.md`,
  `docs/reference.md` and the timer guide follow.
- `docs/coordinating.md`: "Applying a cycle" is "Building a cycle
  unattended": `/grove-implement <cycle-id> --unattended`, each member by
  `grove instructions implement <change-id> --unattended`, the cycle's list
  for every member with the owner or not, one change with `--unattended`,
  a moved file blocking its task, and rules under `rules.skills.implement`.
- `docs/running.md`: the prompt `grove run` types, `/grove-implement
  <cycle-id> --unattended` for skills and `/grove:implement <cycle-id>
  --unattended` for commands; the rules example is under `implement`, and
  says those lines are given with the owner too.
- `docs/reference.md`: `--unattended` rows for a change and a cycle under
  `instructions implement`, a member's cycle list, no `instructions apply`
  rows, six skills with implement in the table of how each harness types a
  skill, `init` removing `grove-apply`, the `rules.skills` list without
  apply, and `check`'s report of `rules.skills.apply`.
- `docs/index.md`, `docs/concepts.md`, `docs/getting-started.md`: implement
  with `--unattended` in place of apply; `docs/guides/glossary-and-decisions.md`
  puts its building rule under `implement`.
- `dev/walkthrough.md`: init lists six skills; coordinating step 7 is
  "Building a cycle unattended"; "Running unattended" puts its rules under
  `rules.skills.implement` and expects `/grove-implement <cycle-id>
  --unattended` and `grove instructions implement <change-id>
  --unattended` in the output; other paths add "One change, unattended" and
  "Upgrading from apply", and the notes speak of implement.
- `grep -rn "grove-apply\|instructions apply\|rules.skills.apply"` over
  `README.md`, `docs/`, `mkdocs.yml` and `dev/walkthrough.md` finds only
  the reference's and walkthrough's upgrade notes, which tell no reader to
  use them.
