## T-01 — Report a change's moved files in one line

Result: `grove check` reports each change in no open cycle whose named
files changed on the target since `base` once, on its `tasks.md`, naming
each changed file once.

Files:
- `internal/cli/check.go`
- `internal/cli/walkthrough_test.go`
- `docs/coordinating.md`
- `dev/walkthrough.md`

Interfaces: `check`'s line `<change>/tasks.md: <file> and <file> changed
on <target> since base <sha7>; revise the change with propose, which also
re-bases it after a dependency lands`.

Acceptance:
- On the fixture, a change whose T-01 names `ledger/` and
  `tests/test_cli.py`, with `ledger/report.py` and `tests/test_cli.py`
  changed on main since `base`, gets one problem naming both files, as in
  "Report a named file that moved" and "Say that revising re-bases".
- A file two tasks name is named once.
- The docs and walkthrough quote the new line.
