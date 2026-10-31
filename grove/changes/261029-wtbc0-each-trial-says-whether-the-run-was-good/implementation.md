# Implementation: Each trial says whether the run was good

<!-- Progress on change 261029-wtbc0, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — The Good run column

Status: in progress

`dev/trials/run.sh`: after each turn the script appends `<turn>\t<branch>\t<owner's
message>` to `turns.txt` and writes `snapshot-<turn>.txt` (every file of
the fixture with its hash); before each run's first turn, `start-<turn>.txt`.
Helpers for the case files: `answered <text>` prints the turn whose
owner's message holds the text, so the turn a scripted answer was given;
`kept <path> <branch> [<turn>]` says whether the path stayed as at the
start of the branch's run after each of its turns before that turn, or
after all; `task <T-NN>` prints a task's section of the `implementation.md`
on stdin. Once the case's runs end and what they left is committed, the
case's `good` runs in the fixture into `good.txt`, one line per failed
check; no model is called. The summary prints a "Good run" column after
Harness, `passed/k`, a run passing when it finished and `good.txt` is
empty, and lists each failed check below the table as `(<case>) <harness>
run <r>: <check>`. It exits 0 whatever the column shows. `dev/trials.md`:
a "Good run" row in the counts table, the column and lines in "Results",
and the new files in the run's folder.

Checks so far:
- The summary on hand-written folders, `RUNS=3`: (c) Claude Code, every
  `good.txt` empty, `3/3`; (c) Codex, run 2's `good.txt` with two lines,
  `2/3`, and both lines below the table as `(c) Codex run 2: …`; exit 0.

### T-02 — Each case's checks

Status: in progress

Each case file has a `good` beside `setup`, checking what the proposal's
"Intended outcome" lists for it, and `dev/trials.md` lists them under each
case's "A good run". `e.sh` keeps T-01's test in `delete_test`, which its
setup writes and its `good` looks for unchanged. How two checks were made
concrete:
- (a) runs `python3 -m ledger edit 3 --note x` with `--file` on a copy of
  `records.json`, so the check does not change the fixture it reads.
- (e)'s "`README.md` shows `ledger delete`" looks in the README's
  "Usage" only: its "Planned for 0.4" already says `ledger delete ID`, so
  a check of the whole file could not fail.

Checks: a scratch script (not in the repository) copies each case's
`SETUP_ONLY=1` fixture, writes a run's end state, `turns.txt` and the
snapshots by hand, and runs the case's `good` with `run.sh`'s helpers.
Every passing state printed nothing, and each broken check printed one
line naming it:
- (a): not archived; schema `default`; a `specs/` folder; proposal without
  `record-editing`; a failing test; no change proposed. With the fix
  undone, both `python3 -m unittest fails` and `python3 -m ledger edit 3
  --note x fails`, as both need the fix.
- (b): the fallback given instead of the ID answer; `ledger/cli.py`
  changed on `delete`.
- (c): `ledger/records.py` edited in turn 1, the turn before the scripted
  answer in turn 2; with the owner, unchanged at the end; unattended,
  edited; unattended T-01 `in progress`; its T-01 not naming the file;
  a failing test on `unknown-id-impl`.
- (d): `ledger/report.py` edited in turn 1, with the answer in turn 2;
  edited in turn 2, with the answer in turn 3; no answer given; a failing
  test.
- (e): no `unittest` TOOL line in the transcript; a line added to
  `test_delete_record`; T-02 `in progress`; no `ledger delete` in Usage; a
  failing test.
