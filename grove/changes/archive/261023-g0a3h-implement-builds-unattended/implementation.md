# Implementation: Implement builds unattended, and apply goes

<!-- Progress on change 261023-g0a3h, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Implement takes `--unattended`, and apply goes

Status: done

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
- Trial A, in Claude Code 2.1.310 on a fresh ledger fixture with this
  build first on PATH, headless: an approved change in no cycle, "Unknown
  ID error", whose open question decides T-02, and `claude -p
  "/grove-implement <change-id> --unattended"`. The session ran
  `grove instructions implement <change-id> --unattended`, created a
  branch, built T-01 with `test_edit_unknown_id` (11 tests pass), marked
  T-02 `Status: blocked` with the question under it, asked nothing (no
  question tool call, none in its report), left proposal, specs and tasks
  unchanged (`git diff main` empty over them), and did not archive.
- `TestPrototype` also checks that a project with no `context` gets no
  "Project context" section.

### T-02 — `grove run` launches implement unattended

Status: done

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
- Trial B: the same change in a cycle of one on the fixture, `grove run
  <cycle-id>`. The session's first call was `grove instructions implement
  <cycle-id> --unattended`, and the output showed `› Bash grove
  instructions implement <change-id> --unattended`; T-01 done, T-02
  blocked with its question, nothing asked, nothing archived.
- Trial C: cycle `[B, A]` on the fixture with B depending on A, no
  `origin`, `rules.skills.implement` asking for a commit per task. `grove
  run <cycle-id>` ran `grove instructions implement <A> --unattended`,
  archived A, then `grove instructions implement <B> --unattended`, and
  archived B: four commits, one per task named by its qualified ID and one
  per archive; nothing asked; 12 tests pass; `grove check` clean on the
  run branch.

### T-03 — Document implement unattended

Status: done

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
- Walked on the ledger fixture in Claude Code by trials A, B and C above:
  "One change, unattended" (A), "Running unattended" step 2's prompt and
  per-member instructions with rules under `rules.skills.implement` (B, C),
  and coordinating step 7's build order, archiving each member before the
  next and asking nothing (C, through `grove run`, which types the same
  `/grove-implement <cycle-id> --unattended`). Not walked by hand: step 7's
  `needs_human` member and "Without a remote"'s pause; `TestImplementInstructions`
  and `TestRun` cover them.
- Recheck before finishing: T-01 to T-03's evidence holds against the code
  as it now is. `just check` passes (3705 Go lines, 1283 test lines);
  `grove check` prints "No problems in grove/.".

## Scenarios

- implement "One change, by hand": trial A, passed (branch created, T-01 built; with no open question T-02 would be built the same way). `TestImplementInstructions` checks the instructions.
- implement "An open question decides a task": trials A and B, passed.
- implement "A named file moved since base": `walkthrough_test.go` `TestMovedFiles`, passed, for the instructions naming T-01's file with `--unattended`; the agent marking it blocked is unverified, no trial staged a moved file.
- implement "Two members, one waiting on the other": trial C, passed; `TestCoordinating`, passed.
- implement "A cycle that needs a person, by hand": `TestImplementInstructions`, passed (instructions printed for a cycle not marked, `needs_human` naming its member; the unattended part says to build whatever they say). Building it in a session is unverified.
- implement "A change that waits" (Build unattended): `TestCoordinating`, passed.
- implement "The second member, unattended": `TestCoordinating`, passed.
- implement "A member built with the owner": `TestImplementInstructions`, passed.
- implement "A rule for both modes": `TestImplementInstructions`, passed.
- implement "A project that kept its apply rules": `TestCheck`, passed.
- implement "Init after upgrading": `TestImplementInstructions`, passed.
- implement "The old command": `TestImplementInstructions`, passed.
- implement "A change with two tasks": unverified with the owner in this run (no owner here); the instructions are `TestImplementInstructions`'s, passed.
- implement "A change that is only a draft": `TestImplementInstructions`, passed, with and without `--unattended`.
- implement "A change that waits" (Implement an approved change): `TestCoordinating`, passed.
- implement "An open question decides a task" (Ask the owner): unverified, needs an owner in the session; `TestImplementInstructions` checks the owner's part says to ask and wait.
- implement "The owner leaves it open": unverified, needs an owner in the session.
- implement "A cycle the owner builds": unverified, needs an owner; `TestImplementInstructions` checks the cycle form with the owner's part.
- implement "Every scenario has a test": trials A and B, passed (`## Scenarios` named the test and "passed").
- implement "A done task later broken": unverified, no trial staged it; the rule is in the shared requires text, checked by `TestImplementInstructions`.
- implement "A scenario only a person can judge": unverified, no trial staged one.
- implement "A scenario no test shows, unattended": trial A, passed (the delete scenario marked unverified, named in the report, nobody asked).
- config "Building unattended": `TestImplementInstructions`, passed.
- config "No context": `TestPrototype`, passed.
- propose "The owner accepts the round", "The owner declines one", "Nothing to assume": unverified here, no propose session was run; the proposal guide's line on `--unattended` is shown by `grove instructions propose`, and the round rule is unchanged.
- runs "Local target has commits of its own": `TestRun`, passed.
- runs "A marked cycle of two changes, one waiting on the other": `TestRun`, passed (stub `claude` given `-p /grove-implement <k> --unattended`).
- runs "Claude Code listed as commands": `TestRun`, passed.
- runs "Two marked cycles": `TestRun`, passed.
- runs "A given cycle not marked": `TestRun`, passed; trials B and C, passed.
- runs "A given cycle while paused": `TestRun`, passed.
- runs "A project that does not list Claude Code": `TestRun`, passed.
- runs "A checkout with no remote": `TestRun`, passed; trial C, passed.
- runs "The branch merged locally": `TestRun`, passed.
- runs "A blocked change revised": `TestRun` (`--retry` starts a session on the implement skill), passed; the full revise-and-re-plan path is unverified here.

## Reconciliation

- apply: retired; `grove/specs/apply/` deleted.
- config: applied to specs/config/spec.md
- implement: applied to specs/implement/spec.md; Purpose, Boundaries and description name `--unattended` instead of apply.
- propose: applied to specs/propose/spec.md; Boundaries no longer name apply.
- runs: applied to specs/runs/spec.md; Purpose, Boundaries and description name the implement skill, unattended.
- cycles: no change spec; Boundaries point building a cycle at implement.
