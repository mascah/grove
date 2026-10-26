# Re-base a change after its dependency lands

## Problem

The review of Grove 0.5.0 on memex (`dev/memex-review-2026-10-26.md`,
finding 4) observed that a change with `depends_on` is proposed before its
dependency lands, so every file the dependency creates that its tasks name
is reported as "changed on main since base <base>; revise the change with
propose" from the moment the dependency merges until the change is
archived. Seen on 2c149 (4 files) and epcxe (7 files, 11 notices). Both
times the owner said to build on; the notices kept firing through `check`
and `archive`, and one agent filtered them with `grep -v`.

At ea7fbaf the message, in `movedSinceBase` in `internal/cli/check.go`,
reads as an error rather than as the expected step. Propose's
instructions (`internal/instruct/requires/propose.md`) set `base` only for
a change reported stale on a requirement ("Change `base` only so"), so
revising a change for a moved file does not clear the report.

## Intended outcome

Once a change it depends on lands, revising the change with propose
re-bases it: the agent brings the tasks up to date with the named files as
the target holds them and sets `base` to the target's commit, and `grove
check` is clean again. The moved-file message, the owner's implement
instructions and the docs say so. Grove's own logic for what counts as
moved does not change.

## Scope

### In scope

- The moved-file message of `grove check`.
- Propose's instructions for revising a change reported for a moved file.
- The owner's implement instructions for a task whose file changed since
  `base`.
- `docs/coordinating.md` and the walkthrough.

### Out of scope

- Silencing the report for files only a dependency's archive changed: it
  needs per-file `git log` and, for changes, an amendment to rule 10's
  exemption (`dev/memex-review-2026-10-26.md`, finding 4).
- `grove check`'s exit status on the report.
- Unattended implement, which still blocks such a task.

## Decisions

- "Re-basing with propose's revise path is the normal step after a
  dependency lands; the message, the implement instructions and the docs
  say so. No new logic." (owner, `dev/memex-review-2026-10-26.md`,
  decision 4)
- One wording for every change, with no condition on `depends_on`:
  "T-01's x changed on main since base b; revise the change with propose,
  which also re-bases it after a dependency lands"; `check` still exits
  non-zero (round, "yes").
- The owner's implement instructions recommend propose's revise path when
  the moved file came with a dependency that landed; unattended still
  blocks the task (round, "yes").
- Propose's revise path sets `base` for a change reported for a moved
  file, as well as for one stale on a requirement; without it decision 4
  cannot clear the report (owner, after the change was written: "Regarding
  D: Confirmed").
- This change depends on 261026-804yw, which edits the same implement
  instructions (owner, in this conversation: "yes").
- The instruction text stays at or under 8,000 words (rule 12) by
  tightening the text this change touches (round, "yes").

## Open questions

None.

## Related documents

- Requirements: [propose](specs/propose/spec.md),
  [implement](specs/implement/spec.md)
- Tasks: [tasks.md](tasks.md)
