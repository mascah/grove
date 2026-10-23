## T-01 — Implement takes `--unattended`, and apply goes

Result: `grove instructions implement` and the implement skill behave as
[specs/implement/spec.md](specs/implement/spec.md) says, the apply skill
and `grove instructions apply` are gone, and the other skills'
instructions say what [specs/config/spec.md](specs/config/spec.md) and
[specs/propose/spec.md](specs/propose/spec.md) say.

Files: `internal/cli/cli.go`, `internal/cli/instructions.go`,
`internal/cli/check.go`, `internal/cli/init.go`,
`internal/instruct/instruct.go`, `internal/instruct/skills.yaml`,
`internal/instruct/requires/implement.md`,
`internal/instruct/requires/implement-cycle.md`,
`internal/instruct/requires/implement-owner.md`,
`internal/instruct/requires/implement-unattended.md`,
`internal/instruct/requires/apply.md`,
`internal/instruct/requires/apply-cycle.md`,
`internal/instruct/schemas/default/schema.yaml`,
`internal/instruct/schemas/default/skills/implement.md`,
`internal/instruct/schemas/default/skills/apply.md`,
`internal/instruct/schemas/default/documents/proposal/guide.md`,
`grove.yaml`, `.claude/skills/grove-implement/SKILL.md`,
`.agents/skills/grove-implement/SKILL.md`,
`.claude/skills/grove-apply/SKILL.md`,
`.agents/skills/grove-apply/SKILL.md`, `internal/cli/walkthrough_test.go`

Interfaces: `grove instructions implement <change-id> | <cycle-id>
[--unattended]`, which `grove run` (T-02) types through the skill.
`requires/implement.md` and `requires/implement-cycle.md` hold the
shared rules for a change and a cycle; `requires/implement-owner.md` or
`requires/implement-unattended.md` follows them, for a change and a cycle
alike. `instruct.Implement` and `instruct.Cycle` take whether the owner
is present; `instruct.Apply`, `instructionsApply`, `requires/apply.md`,
`requires/apply-cycle.md` and `skills/apply.md` are deleted.

Constraints: each rule of today's apply and implement requires text, and
261023-m2she's verify rule, is written once, in the shared files or in
one mode's part (decision 4); the schema guide holds only what both modes
do, so telling the owner how to try a task moves to the owner's part. The
`Status:` lines stay exactly as the template has them (rule 3).

Acceptance:
- `walkthrough_test.go`, on the fixture: `instructions implement` for a
  change and for a cycle, with `--unattended`, prints the shared rules,
  261023-m2she's verify rule among them, and the unattended part and not
  the owner's; without it, the owner's part and not the unattended one.
- With `context` in `grove.yaml`, `instructions implement <change-id>
  --unattended` prints it before "What Grove requires of implement".
- The proposal guide says that `/grove-implement --unattended` leaves a
  task an open question decides blocked, while `/grove-implement` asks.
- The cycle form with `--unattended` tells the agent to run `instructions
  implement <change-id> --unattended` for each member, and prints
  instructions for a cycle not marked `unattended` and for one whose
  `needs_human` names a member.
- `instructions implement B`, with and without `--unattended`, for the
  second member of an open cycle whose first is archived on the branch,
  prints the cycle's plan naming B as the second of two.
- `--unattended` refuses a draft and a change that waits, with the same
  messages as without it.
- `instructions apply` prints the usage and exits 2; `grove init` on the
  fixture removes the `grove-apply` skill and command files and names them
  among what it removed.
- `rules.skills.implement` lines are printed in both modes; `grove check`
  reports a `rules.skills.apply` line, pointing to `implement`.
- Grove's own `grove.yaml` holds its apply lines under
  `rules.skills.implement`, its skill files are regenerated with `grove
  init`, and `grove check` prints "No problems in grove/.".
- On the ledger fixture, in Claude Code, `/grove-implement <change-id>
  --unattended` builds a change whose open question decides a task,
  marking that task blocked without asking.
- `walkthrough_test.go`: `instructions implement <change-id> --unattended`
  for a change in no cycle whose task names a file moved since base names
  that task and file, and its unattended part says to mark the task
  blocked with the reason and not build on the file.
- Archiving this change deletes `grove/specs/apply/` and points the
  Purpose and Boundaries of the maintained `implement`, `runs`, `cycles`
  and `propose` specs that name apply at implement.
- The implementing pull request is titled `feat!` and its
  `BREAKING CHANGE:` footer states the proposal's migration.

## T-02 — `grove run` launches implement unattended

Result: `grove run` behaves as [specs/runs/spec.md](specs/runs/spec.md)
"Start one attempt" and "Retry" say.

Files: `internal/cli/run.go`, `internal/cli/cli.go`,
`internal/cli/walkthrough_test.go`

Interfaces: the session's prompt, `/grove-implement <cycle-id>
--unattended` for skills or `/grove:implement <cycle-id> --unattended`
for commands.

Acceptance:
- `TestRun`: the stub `claude` is given `-p` with `/grove-implement <cycle-id>
  --unattended` for a project listing Claude Code as skills, and
  `/grove:implement <cycle-id> --unattended` for one listing it as
  commands.
- `grove run`'s line in the usage names no apply skill.
- On the ledger fixture, the output of `grove run <cycle-id>` shows the
  session running `grove instructions implement <change-id> --unattended`
  for each member.

## T-03 — Document implement unattended

Result: the documentation describes `/grove-implement <id>
--unattended`, the cycle's plan for every member and `grove run`'s use of
the implement skill, and tells no reader to use apply (rule 8).

Files: `README.md`, `docs/index.md`, `docs/concepts.md`,
`docs/getting-started.md`, `docs/coordinating.md`, `docs/planning.md`,
`docs/running.md`, `docs/reference.md`,
`docs/guides/run-apply-on-a-timer.md`, `docs/guides/run-on-a-timer.md`,
`docs/guides/glossary-and-decisions.md`, `mkdocs.yml`,
`dev/walkthrough.md`

Interfaces: none.

Acceptance:
- `docs/guides/run-apply-on-a-timer.md` is moved to
  `docs/guides/run-on-a-timer.md`, and `mkdocs.yml`, `docs/index.md`,
  `docs/running.md` and any other page link the new path; no page or nav
  entry names the old one.
- No page or walkthrough step tells the reader to type `/grove-apply`, run
  `grove instructions apply` or write lines under `rules.skills.apply`.
- `docs/reference.md` has `--unattended` in the `instructions implement`
  rows, no `instructions apply` rows, implement in the table of how each
  harness types a skill, the `rules.skills` list without apply, and
  `check`'s report of `rules.skills.apply`.
- `docs/coordinating.md` describes building a cycle, or one change, with
  `--unattended`, and the plan a member gets; `docs/running.md` names the
  prompt `grove run` types, for skills and for commands.
- `dev/walkthrough.md`'s coordinating and running steps use
  `/grove-implement <id> --unattended` and put their rules under
  `rules.skills.implement`, and are walked on the ledger fixture in
  Claude Code.
