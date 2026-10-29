## T-01 — The Good run column

Result: `just trials` prints a "Good run" column from each case's `good`
function, and takes a snapshot after every turn, as the proposal's
decisions say.

Files:
- `dev/trials/run.sh`
- `dev/trials.md`

Interfaces: a case file's `good` function, run in the fixture once the
case's runs end, printing one line per failed check; the snapshot after
each turn, and the turn each scripted answer was given in, for `good` to
read.

Constraints: rule 17; no Go, so rule 12's budget is unchanged. No model
call decides the column.

Acceptance:
- A case whose `good` prints nothing shows `pass`; one whose `good` prints
  two lines shows `fail`, and both lines are listed below the table under
  its case, harness and run.
- `just trials` exits 0 when a row shows `fail`.
- `dev/trials.md`'s counts table has a "Good run" row saying how it is
  taken, and "Results" shows the column.

## T-02 — Each case's checks

Result: each case file's `good` checks what the proposal's "Intended
outcome" lists for that case, and `dev/trials.md`'s case says so.

Files:
- `dev/trials/a.sh`
- `dev/trials/b.sh`
- `dev/trials/c.sh`
- `dev/trials/d.sh`
- `dev/trials/e.sh`
- `dev/trials.md`

Interfaces: none.

Acceptance:
- For each case, on its setup with a run's end state written by hand, a
  state meeting every check passes, and breaking any one check fails
  with a line naming it.
- In case (c), a records.py edit in the turn before the scripted answer
  fails the run with the owner; in case (d), a `ledger/report.py` edit in
  the first turn fails it.
- Each case of `dev/trials.md` lists its checks under its "A good run".
