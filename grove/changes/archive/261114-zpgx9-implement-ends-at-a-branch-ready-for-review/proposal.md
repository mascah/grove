# Implement ends at a branch ready for review

## Problem

The owner's analysis of 76 Claude Code sessions in memex, 2026-10-24 to
2026-11-14 (the owner's brief `memex/docs/efficiency-2026-11-14.md`, not in
this repository), found:

- In `/grove-implement` sessions of Nov 8-14, 1.2 h of agent work against
  16.2 h waiting on the owner. 112 of 472 end-of-turn stops asked the
  owner to try something in the app; 42 `Check:` lines got 29 answers and
  none rejected what was built.
- Asked "Why can't you check them yourself", the agent checked two
  checkpoints itself in 17 minutes and found two bugs.
- Base drift was asked about 4 times and answered "build on current" 4
  times. Unattended, two tasks of one change were blocked over a drifted
  file though nothing conflicted; told to go on, both were built in 8
  minutes.

The instructions cause this. `internal/instruct/requires/implement.md` at
999ab56 gives "how something looks" as the example of a checkpoint;
`requires/implement-owner.md` stops once to put checkpoints to the owner,
asks before building on a drifted file, and edits the change only with the
owner's agreement; `requires/implement-unattended.md` blocks a drifted
task and never revises the change. Since #158, `CONSTITUTION.md` rule 4
says a change revised on its branch by implement is approved as revised
once the target holds it, so the merge is the agreement.

## Intended outcome

`/grove-implement` always ends at a branch ready for review and never
waits on anyone. It looks at the app or command before calling acceptance
unseeable, builds its best choice where the change leaves one open and
records it under the task as a `Check:` line, builds on a drifted file as
it is now, and revises the change on the branch through propose where it
must. Only a choice whose wrong answer would cost data, break an interface
other tasks build on, or be hard to undo is left `Status: blocked`, and
every other task is still built. The finish report lists every `Check:`
line and blocked task first, for the reviewer. With the owner or not, the
instructions are the same, so `--unattended` goes.

## Scope

### In scope

- One set of implement instructions (`requires/implement.md`), replacing
  `requires/implement-owner.md` and `requires/implement-unattended.md`;
  the cycle instructions (`requires/implement-cycle.md`) to match.
- `Status: to check` goes from the implementation template and guide: a
  built choice is `done` with a `Check:` line.
- `--unattended` removed from `grove instructions implement`, the skill
  description, the proposal guide, the specs, the docs, the walkthrough
  and the skill trials.
- Trial (d) now expects the agent to build on the drifted file without
  asking; trial (c) loses its `--unattended` run.

### Out of scope

- How `grove check` prints drift (part C2) and `needs_human` (part C3):
  their own changes.
- Opening a pull request: Grove opens none (rule 11); a project's `rules`
  line has its agent do it.
- `grove claim`'s gates, including the stale plan (rule 10).

## Decisions

The owner answered in the brief and said "Do not ask a round. Take the
recommended choice, write it down where the reviewer will see it, and go
on." Each decision below is that brief's recommendation, part C1.

1. No stop at a checkpoint: "For a choice the change leaves open, take the
   best one, build it, and record it under the task. Archive goes ahead.
   The finish report lists every such guess first, for the reviewer." The
   record is the existing `Check:` line under a `Status: done` task, so
   `Status: to check` goes; a project's old `to check` tasks still stop
   archive, as any status but `done` does.
2. "Before calling acceptance unseeable, run the project's app or command
   and look", replacing "how something looks" as a checkpoint's example.
3. A task's file changed on the target: "read the diff from the base.
   Build on the file as it is and record a line under the task. If the
   diff changes what the change says, revise the change on the branch
   through `instructions propose <change-id>`, which re-bases it, and go
   on." The line is a `Check:` line, so the reviewer finds every judgement
   in one place.
4. "`Status: blocked` is kept for a wrong guess that would cost data, break
   an interface other tasks build on, or be hard to undo. Even then build
   every other task and hand the branch off with the question; do not
   wait."
5. "The owner's mode and `--unattended` say the same thing. Prefer one set
   of instructions and delete the other; keep the flag only if a
   documented step still needs it (rule 8)." No step does once the
   instructions are one, so the flag goes: breaking. A prompt, script or
   timer drops ` --unattended` from `/grove-implement <id> --unattended`.
   A draft is then built by anyone who runs `/grove-implement` on its
   branch; `grove claim` still claims only cycles the target holds.
6. Implement revises the change on its branch, through propose, without
   asking: rule 4 as amended in #158.
7. Instruction words drop on net (rule 12, at 8,613 of 9,000).

## Open questions

None.

## Related documents

- Requirements: [implement](specs/implement/spec.md),
  [config](specs/config/spec.md), [runs](specs/runs/spec.md),
  [propose](specs/propose/spec.md)
- Tasks: [tasks.md](tasks.md)
