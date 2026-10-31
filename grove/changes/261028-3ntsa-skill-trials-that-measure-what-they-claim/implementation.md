# Implementation: Skill trials that measure what they claim

<!-- Progress on change 261028-3ntsa, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Outside edits counted by grove review

Status: done

`dev/trials/run.sh`: `allowed` and `named` (the awk copy of Grove's rule)
and the per-run `snapshot` are gone. After the case's last run the script
commits what it left (`save`), and `outside` lists, on each branch the
case used, each file outside `grove/` the branch changes against main that
`grove review` marks neither `planned` nor `agreed`, as `branch: path` in
`outside.txt`. A branch whose review lists nothing, as one holding no
change, counts every such file. The column is "Edits outside Files/Also
changed"; `dev/trials.md`'s counts table says how it is taken.

Checks, on `SETUP_ONLY=1 HARNESSES=claude dev/trials/run.sh a c` with the
end state written by hand and `outside` run in each fixture:
- (a): on `fix` from main, a new change whose T-01 names only
  `ledger/rules.py`, and `ledger/cli.py` changed: `fix: ledger/cli.py`,
  1. The awk of a5b5674 on the same state lists `ledger/cli.py`,
  `ledger/rules.py`, `tests/test_cli.py` as allowed, so 0.
- (c): on `unknown-id-impl`, `ledger/records.py` and `ledger/cli.py`
  changed: `unknown-id-impl: ledger/records.py`, 1. With `Also changed:
  \`ledger/records.py\`` under T-01 in `implementation.md`: none, 0
  (`grove review` marks it `agreed 261015-c0000:T-01`).
