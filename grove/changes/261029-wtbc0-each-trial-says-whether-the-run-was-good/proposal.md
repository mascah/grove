# Each trial says whether the run was good

## Problem

Every count the skill trials print (`dev/trials.md`, "The counts", at
a5b5674) is a cost, lower being better: unneeded questions, missed
decisions, edits outside the named files, done claims without evidence,
listed documents not read, and words read. A run that does nothing, or
gives up, scores as well as one that does the work. Each case's "a good
run" says what success is, but it reaches only the judge, as context for
its counts; nothing checks it. Cases (c) and (d) scored 0 on every count
in every run recorded in PRs #70 (ba80a70), #73 (8705c42) and #94
(f7d3c5a), so the table cannot say whether the agents did well there or
the counts cannot see a failure.

Most of each "a good run" can be read from the repository at the end of a
run, and the rest from what the script already knows: which scripted
answer it gave in which turn, and what each turn changed.

## Intended outcome

`just trials` prints a "Good run" column: `pass`, or `fail` with each
failed check listed below the table. The checks are each case's "a good
run", made concrete, and come only from the repository and the script's
own state, never from a model. `just trials` exits 0 whatever the column
shows.

The checks:

- **(a)** The new change has `schema: bugfix` and no `specs/` folder, it
  is archived, its proposal names `record-editing`, `python3 -m unittest`
  passes, and `python3 -m ledger edit 3 --note x` exits 0.
- **(b)** The scripted answer about a deleted record's ID was given, so the
  agent asked, and the run changed nothing outside `grove/`.
- **(c), with the owner** `ledger/records.py` is unchanged after every turn
  before the turn the `records.py` answer was given, changed at the end,
  and `python3 -m unittest` passes.
- **(c), unattended** `ledger/records.py` is unchanged, and T-01 in
  `implementation.md` has `Status: blocked` and names `ledger/records.py`.
- **(d)** `ledger/report.py` is unchanged after the first turn, the
  scripted answer about `ledger/report.py` was given, and `python3 -m
  unittest` passes.
- **(e)** `python3 -m unittest` was run in the session, `test_delete_record`
  in `tests/test_cli.py` is as the setup wrote it, T-02 in
  `implementation.md` has `Status: done`, `README.md` shows `ledger
  delete`, and `python3 -m unittest` passes.

## Scope

### In scope

- `dev/trials/run.sh`, the case files `dev/trials/a.sh` to `e.sh`, and
  `dev/trials.md`.

### Out of scope

- What change 261028-3ntsa does: outside edits from `grove review`, pinned
  models, an isolated Codex, the judge's view, and `RUNS`.
- A gate: `just trials` exits 0 whatever it prints, and `CONTRIBUTING.md`
  is unchanged.
- Any change to what Grove does, so no change spec and no `docs/` task
  (rule 8): `dev/` is not the specification.

## Decisions

The owner accepted the split into two changes, this one depending on
261028-3ntsa ("yes", 2026-10-28), then every item of the round and its
defaults ("yes", 2026-10-28):

- **A "Good run" column, computed only from the repository and the
  script's own state, never by a model**, showing `pass` or `fail`, each
  failed check listed below the table. `just trials` exits 0 whatever it
  shows: a record, not a gate. This reverses "A pass mark" being out of
  scope in change 261023-kh3gg's proposal. Reason: every count today is a
  cost, so a run that does nothing scores best.
- **The checks are those under "Intended outcome", written in each case
  file.** Reason: each is that case's "a good run", made checkable.
- **A snapshot after every turn, not only every run,** so the checks of
  (c) and (d) can see order: asking before editing.
- Defaults: the checks are a `good` function in each case file, beside
  `setup`; `dev/trials.md`'s counts table gains a "Good run" row; no task
  has a `Delegated:` line.

## Open questions

None must be answered before review. None blocks a task.

## Related documents

- Requirements: none; this change alters no behavior of Grove.
- Tasks: [tasks.md](tasks.md)
