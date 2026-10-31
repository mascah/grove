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
