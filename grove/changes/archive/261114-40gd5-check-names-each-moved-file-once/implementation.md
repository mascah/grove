# Implementation: Check names each moved file once per change

<!-- Progress on change 261114-40gd5, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Report a change's moved files in one line

Status: done

`movedSinceBase` in `internal/cli/check.go` gathers the moved files of every
task, sorts them, drops repeats and appends one diagnostic on the change's
`tasks.md`, the files joined by "and". `TestMovedFiles` now expects one
line naming `ledger/report.py and tests/test_cli.py` where it expected two,
and adds a T-02 naming `tests/test_cli.py` to show a file two tasks name is
named once. `docs/coordinating.md` and `dev/walkthrough.md` quote the new
line. `just check` passes: 4409 Go lines, 2101 test lines, 8613 words.

Also changed: `dev/trials/d.sh` — its `check_ok` pattern matched the old line, `T-01's ledger/report.py changed on main since base`, so the setup of trial (d) would have failed on the new one.
Also changed: `go.mod` — the pre-commit hook's `go mod tidy` check failed, as it does on main at 999ab56, since no Go file imports Lip Gloss directly; tidy marks it `// indirect`, so the commit could be made. It stays a dependency, through Glamour (rule 16).

## Scenarios

- propose "A file changed since base": `internal/cli/walkthrough_test.go` TestMovedFiles, `ledger/report.py and tests/test_cli.py changed on main since base`, passed
- propose "A file under a named folder": TestMovedFiles, `ledger/report.py` reported through T-01's `ledger/`, passed
- propose "Several files, several tasks": TestMovedFiles, T-01 and T-02 naming `tests/test_cli.py`, one line naming it once, passed
- propose "A member of a cycle": TestMovedFiles, reported against the plan, not base, passed
- propose "A file this change's own branch touched": TestMovedFiles, `older-edits-impl` branch checks clean, passed
- propose "The message": TestMovedFiles, the sentence ends "; revise the change with propose, which also re-bases it after a dependency lands", passed

## Reconciliation

- propose: applied to specs/propose/spec.md
