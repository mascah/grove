# Unattended implement edits a file a task needs

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

## Intended outcome

- Unattended, the agent edits a file that a task needs changed to be built
  as the change says, or that the task makes false, even when no `Files:`
  line names it. It names the file and why in that task's evidence in
  `implementation.md` and writes no `Also changed:` line. The task is not
  blocked for it.
- Unattended, the agent follows a task's `Delegated:` line, as it does
  with the owner.
- `grove review` shows such a file as `named nowhere`, as it already does
  for any file no `Files:` or `Also changed:` line names, so the owner sees
  it before merging.

## Scope

### In scope

- The unattended implement instructions, the tasks guide's sentence about
  `Delegated:`, and the wording test that checks the instructions.
- The maintained requirement "Change only the files the tasks name" and
  its scenarios.
- `docs/coordinating.md`, `docs/getting-started.md` and the walkthrough
  step "Delegated with me, ignored unattended".
- Trial case (c)'s unattended good run, in `dev/trials/c.sh` and
  `dev/trials.md`.

### Out of scope

- `CONSTITUTION.md`: rule 3 ties an `Also changed:` line to the owner's
  agreement, and this change writes none unattended, so no amendment.
- Go code: `grove review` already marks such a file `named nowhere`.
- Implementing with the owner: unchanged.
- The trials' "Edits outside" count from 261028-3ntsa: unchanged.

## Decisions

- **Unattended edits the file and names it in the task's evidence, with no
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
- **Trial case (c), unattended:** a good run changes `ledger/records.py`,
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
