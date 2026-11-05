# Implement edits a file a task needs and says why

## Problem

With `--unattended`, the implement instructions forbid any file no task's
`Files:` line names: "A file no task's `Files:` line names stays as it is,
even if it looks like it should change, and even if a `Delegated:` line
names it" (`internal/instruct/requires/implement-unattended.md` at
1a24167). The maintained requirement "Change only the files the tasks
name" in [implement](../../specs/implement/spec.md) says the same. So a
task that needs one more file goes unbuilt or is built around the file.
Trial case (c) in `dev/trials.md` at 1a24167 counts that as the good run:
T-01 needs `find` in `ledger/records.py`, and the unattended run must leave
it alone and mark T-01 blocked.

An agent in another project, running Grove's unattended text alongside
that project's own rule allowing such edits, reported: "Grove's built-in
unattended text still says a file on no Files line 'stays as it is', so
an agent sees both lines" (owner's session, 2026-10-31).

The rule dates from 2026-10-23 (#31, #63). Change 261026-ar7ec kept it
with "unattended never adds one", but its round, in the owner's session of
2026-10-26, asked only how a file the owner agrees to is recorded. It never
asked about unattended runs.
With the owner, the agent stops instead: it tells the owner what the task
makes false in such a file and edits it only "with their agreement"
(`internal/instruct/requires/implement-owner.md` at 6090cb3), and `grove
review` labels the file `agreed` (`internal/cli/review.go`, same commit).

## Intended outcome

- In both modes, with the owner and with `--unattended`, the agent edits a
  file that a task needs changed to be built as the change says, or that
  the task makes false, even when no `Files:` line names it. It does not
  ask first, and the task is not blocked for it. It writes
  `` Also changed: `<path>` — <why> `` under that task in
  `implementation.md`.
- Unattended, the agent follows a task's `Delegated:` line, as it does
  with the owner.
- `grove review` shows such a file as `also changed <change-id>:T-NN:
  <why>`, so the reviewer reads the reason beside the file. A file it
  shows as `named nowhere` is one the agent changed without saying why.

## Scope

### In scope

- The implement instructions, with the owner and unattended, the
  `implementation.md` guide's sentence about `Also changed:`, the tasks
  guide's sentence about `Delegated:`, and the wording test.
- `grove review`'s label for a file an `Also changed:` line names.
- The maintained requirements "Change only the files the tasks name" and
  "Review names what the branch approves and builds", with their
  scenarios.
- `docs/coordinating.md`, `docs/getting-started.md`, `docs/reference.md`
  and the walkthrough steps that build such a file and review it.
- Trial case (c)'s good runs, in `dev/trials/c.sh` and `dev/trials.md`.

### Out of scope

- `CONSTITUTION.md`: rule 3 was amended on 2026-11-05 (#125) so that an
  `Also changed:` line records each file the implementing agent changed
  beyond the task's files, and why.
- The trials' "Edits outside Files/Also changed" count: unchanged, so a
  good run of case (c) now counts 0.

## Decisions

Revised 2026-11-05. The owner said a stop for a person to allow an extra
file "is counter intuitive. The agent should comeplete the work and the
reviewer can evaluate the implementation", and chose, of two options,
that the agent writes the `Also changed:` line itself ("B: agent writes
Also changed", 2026-11-05). Rule 3 was amended first (#125). The owner
then asked to "reconcile everything", which gives the revision below; the
first two items replace the first decision of 2026-10-31.

- **In both modes the agent edits a file a task needs or makes false
  without asking, and writes `` Also changed: `<path>` — <why> `` under
  the task.** Reason: the reviewer judges the implementation, and the
  reason sits where `grove review` shows it.
- **`grove review` labels such a file `also changed <change-id>:T-NN:
  <why>` instead of `agreed ...`.** Reason: after the amendment, the line
  records the agent's change, not the owner's agreement. This label was
  chosen while reconciling, not asked in a round.
- **Trial case (c), with the owner and unattended alike:** a good run
  asks nothing about `ledger/records.py`, changes it, marks T-01
  `Status: done` with an `Also changed:` line naming it, and the tests
  pass; the scripted answer about it goes.

Decided before the revision, as they stand:


- *(Replaced by the revision above.)* **Unattended edits the file and names it in the task's evidence, with no
  `Also changed:` line; `grove review` shows it as `named nowhere`.** The
  owner chose this over amending rule 3 to let unattended write an
  `Also changed:` line ("yes", 2026-10-31). Reason: review already shows
  every changed file from Git, so the edit is visible at merge without a
  constitution change.
- **Unattended follows `Delegated:` lines instead of ignoring them**
  ("yes", 2026-11-01). Reason: the owner agreed to them in the round, and
  ignoring them while allowing edits to unnamed files would be backwards.
- **The files it may edit:** any file the task needs changed to be built
  as the change says, or that the task makes false; never the change's own
  documents other than `implementation.md`, nor `grove.yaml`, as today
  ("yes", 2026-11-01).
- **What still blocks a task:** an open question that decides it, or
  building it would need a change to what the change says. Needing an
  unnamed file is no longer "cannot be built as the change says" ("yes",
  2026-11-01).
- *(Replaced by the revision above.)* **Trial case (c), unattended:** a good run changes `ledger/records.py`,
  marks T-01 `Status: done` with evidence naming `ledger/records.py`, and
  the tests pass. The "Edits outside" count shows 1 for it and stays a
  measurement, with no exemption ("yes", 2026-11-01).
- **`depends_on: [261029-wtbc0]`** ("yes", 2026-11-01). Reason: this change
  edits `dev/trials/c.sh` and `dev/trials.md`, which cycle 261031-yzcwm's
  tasks name, so landing first would make that cycle's plan stale.
- **One change, two tasks: Grove's behavior, then the trials; the
  instruction text stays within rule 12's budget** ("yes", 2026-11-01).

## Open questions

None.

## Related documents

- Requirements: [implement](specs/implement/spec.md)
- Tasks: [tasks.md](tasks.md)
