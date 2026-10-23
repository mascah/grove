# Implement builds unattended, and apply goes

## Problem

Two skills build a change, and most of what they require is written
twice. At 57af2e3, implement and apply each have their own requires text
(`internal/instruct/requires/implement.md`, `requires/apply.md`,
`requires/implement-cycle.md`, `requires/apply-cycle.md`) and their own
schema guide (`internal/instruct/schemas/default/skills/implement.md`,
`skills/apply.md`). The branch rule, the `Status:` rule, "build what it
says, nothing more" and archiving at the end appear in both. Change
261023-m2she adds "Verify before finishing" to implement and "Verify
before finishing a member" to apply, the same rule twice again.

Apply takes only a cycle (`instructionsApply` in
`internal/cli/instructions.go`). Building one change without being asked
means planning a cycle of one first. Only apply's member form, `grove
instructions apply <cycle> <change>`, prints the cycle's plan;
`instructionsImplement` passes no cycle for a member, so an owner
building a member with `/grove-implement` does not see its place in the
cycle.

The constitution, amended in PR #40, says `run` launches the implement
skill with no owner in the session (`CONSTITUTION.md` rule 2), and the
amendment folds apply into implement. `grove run` still types
`/grove-apply` (`internal/cli/run.go`, `runAttempt`). Grove's own
`grove.yaml` keeps near-identical lines under `rules.skills.implement` and
`rules.skills.apply`.

## Intended outcome

`/grove-implement <id> --unattended` builds a change or a cycle as
`/grove-implement <id>` does, but asks no one: a task an open question
decides, or that cannot be built as the change says, is marked blocked,
and the change's files stay as approved. A cycle is built whether or not
it is marked `unattended` and whatever its `needs_human` says; only
`grove run` reads those. `grove run` types `/grove-implement <cycle-id>
--unattended`, or `/grove:implement <cycle-id> --unattended` for the
commands form.

A member of a cycle that is not done always gets the cycle's plan with
its instructions, with the owner or not. Implement's instructions are one
shared requires text and a short part for each mode, so a shared rule is
written once. `rules.skills.implement` applies in both modes.

The apply skill, `grove instructions apply` and the `apply` capability
are gone. `grove init` removes the old skill files, and `grove check`
reports a `rules.skills.apply` key, pointing to `implement`.

The rules are in [specs/implement/spec.md](specs/implement/spec.md),
[specs/apply/spec.md](specs/apply/spec.md),
[specs/runs/spec.md](specs/runs/spec.md),
[specs/config/spec.md](specs/config/spec.md) and
[specs/propose/spec.md](specs/propose/spec.md).

## Scope

### In scope

- `--unattended` on `grove instructions implement` and the implement
  skill, for a change and a cycle; the cycle's plan for every member.
- Removing apply: its skill, `grove instructions apply`, its requires
  text and schema guide, and its row in `skills.yaml` and `schema.yaml`.
- `grove check` reporting `rules.skills.apply`; `grove init`'s template
  comment; Grove's own `grove.yaml` and its generated skill files.
- `grove run` typing the implement skill.
- Retiring the `apply` capability: archiving this change deletes
  `grove/specs/apply/` and points the Purpose and Boundaries of the
  maintained `implement`, `runs`, `cycles` and `propose` specs that name
  apply at implement. The `cycles` spec has no requirement that names
  apply, so it has no change spec here.
- `README.md`, the pages in `docs/` that name apply, `mkdocs.yml` and
  `dev/walkthrough.md`.

### Out of scope

- What `grove run` refuses, and how it selects, records and pauses: they
  stay as `grove/specs/runs/spec.md` says.
- `run` for harnesses other than Claude Code (rule 11).
- Approving a change drafted on its implementing branch with one merge:
  change 261023-6yafg.

## Decisions

The owner agreed to the direction on 2026-10-23 ("yes to all") and to
every answer below in one round ("yes").

- **The form is `/grove-implement <id> --unattended`,** and `grove run`
  types `/grove-implement <cycle-id> --unattended`, or `/grove:implement`
  for the commands form. `/grove-apply` and `grove instructions apply` are
  removed, and `grove init` prunes the old skill files as it already does
  for skills Grove no longer has. Reason: one skill builds a change, as
  rule 2 now says, and a flag says who is present.
- **`--unattended` takes a change ID as well as a cycle ID.** By hand it
  means "build this, don't ask me". Only `grove run` enforces cycles,
  `unattended` and `needs_human`. Reason: planning a cycle of one just to
  build without being asked is ceremony.
- **A member of a cycle that is not done always gets the cycle's plan
  with its instructions,** with the owner or without. Reason: today only
  apply's member form includes it, so the owner building a member does
  not see its place.
- **Implement's instructions are one shared requires file plus a short
  part per mode,** with the owner and unattended. Reason: shared rules are
  written once, so the two modes cannot drift.
- **`rules.skills.implement` applies in both modes; `rules.skills.apply`
  becomes an unknown key that `grove check` reports, pointing to
  `implement`.** Grove's own `grove.yaml` moves its apply lines. Reason:
  one skill, one place for its rules, and a project that kept the old key
  is told where its lines go.
- **The `apply` capability is retired.** Its change spec removes each of
  its requirements, the rules move into `implement`, and archiving this
  change deletes `grove/specs/apply/`. Requirements of `runs` and other
  capabilities that name the apply skill are modified. Reason: a
  maintained spec for a skill that does not exist would say agreed what
  nothing does.
- **The change is breaking,** `feat!` when implemented. The migration:
  replace `/grove-apply <cycle-id>` with `/grove-implement <cycle-id>
  --unattended`, in prompts, scripts and timers; move the lines under
  `rules.skills.apply` to `rules.skills.implement`; and run `grove init`
  again. A project schema that names `apply` under `skills`, or overrides
  `skills/apply.md`, is then reported by `grove check` as an unknown
  skill or an unused file (`internal/instruct/instruct.go`, `problems` and
  `Check`, at 57af2e3), until what it says moves to implement's guide.
  Reason: `CONTRIBUTING.md` counts a removed skill and a removed
  `grove.yaml` key as breaking.

## Open questions

None must be answered before review.

- **A moved file under `--unattended`.** When `/grove-implement <change-id>
  --unattended` is run for a change in no cycle and a task names a file
  the target changed since the change's `base`, implement's instructions
  today have the agent ask the owner ("Tell the owner of a moved file").
  Does the unattended agent mark that task `Status: blocked`, naming the
  file, or build on the file as it is now? Recommended: blocked. It
  decides the unattended part of the requires text in T-01; until it is
  answered, T-01 leaves that line unwritten.
- **The guide's path.** Does `docs/guides/run-apply-on-a-timer.md` keep
  its path, with only its text changed, or move to a path without
  "apply"? It can wait; it decides T-03's `mkdocs.yml` and link edits.

## Related documents

- Requirements: [specs/implement/spec.md](specs/implement/spec.md),
  [specs/apply/spec.md](specs/apply/spec.md),
  [specs/runs/spec.md](specs/runs/spec.md),
  [specs/config/spec.md](specs/config/spec.md),
  [specs/propose/spec.md](specs/propose/spec.md)
- Tasks: [tasks.md](tasks.md)
