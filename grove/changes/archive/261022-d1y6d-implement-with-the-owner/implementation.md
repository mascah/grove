# Implementation: Implement a change with the owner

<!-- Progress on change 261022-d1y6d, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — The implement skill

Status: done

Built on branch `implement-with-the-owner`:

- `internal/instruct/skills.yaml`: the `implement` skill, so `grove init`
  writes it for each harness and `rules.skills.implement` is accepted.
- `internal/instruct/requires/implement.md`: ask the owner instead of
  blocking; revise proposal, change specs and tasks only with the owner's
  agreement, never `.grove.yaml`, never renumbering a task; `grove check`
  clean before building; blocked only when the owner says to leave it.
- `internal/instruct/schemas/default/skills/implement.md` and `schema.yaml`:
  the default schema's guide; implement writes proposal, change-spec, tasks
  and implementation.
- `internal/instruct/instruct.go` (`Implement`), `internal/cli/cli.go` and
  `instructions.go`: `grove instructions implement <change-id>`, sharing
  apply's refusals of a draft and of a change that waits.
- `requires/apply.md` is unchanged: it already says apply asks no one.
- Docs: `docs/getting-started.md` §5 builds with `/grove-implement`;
  `docs/reference.md` lists the skill and the command; `docs/concepts.md`
  says implement asks and may revise on its branch, apply asks no one.
  `dev/walkthrough.md` runs those steps, with a new "Implement asks, and
  revises with the owner" and "Implement refuses a draft".

Checks:

- `internal/cli/walkthrough_test.go`: `TestInit` checks the implement skill
  for both harnesses; `TestApplyInstructions` checks apply and implement
  refuse a draft, and the implement instructions; `TestCoordinating` checks
  both refuse a change that waits. `just check` passes (2826 Go lines, 874
  test lines).
- `grove check` on this repository: no problems.
- On the ledger fixture, `grove init` writes `.claude/skills/grove-implement/SKILL.md`
  running `grove instructions implement <args>`.

Acceptance run in Claude Code 2.1.308 on the ledger fixture, headless
(`claude -p`, resumed per turn), with the implementer playing the owner at
the owner's request:

- `/grove-propose Editing a missing record should print an error`, the
  exit-code question deferred; the proposal names it as blocking T-01.
  Merged to main.
- `/grove-apply <id>`: T-01 `Status: blocked` with the question, nobody
  asked, the change's files unchanged; `grove archive` refuses.
- `/grove-implement <id>` from main: it creates a branch and asks the exit
  code question. Answer "Exit 2"; it lists its edits to the proposal, change
  spec and tasks and waits. On agreement it edits them, `grove check` is
  clean, it builds T-01 with a test (11 tests pass), archives, and
  `.grove.yaml` is unchanged. `grove review` shows the change spec's
  "exit with status 2" and the tasks edit; the code is in the same branch.
- Declining the edits in a second session: the change's files are
  unchanged and T-01 is built as the change says.
- `grove instructions implement` for a change only on a branch: "grove: main
  holds no approved change 261022-zzzzz".
- The first run built T-01 before the owner agreed to the edits and edited
  one proposal line beyond what was agreed. `requires/implement.md` now says
  to build nothing that depends on the edit yet and to edit only as agreed;
  the rerun waited for agreement before building.

Also, at the owner's request: `README.md` and `docs/index.md` name
`/grove-implement` as the way to build a change.

## Reconciliation

- implement: applied to specs/implement/spec.md (new)
