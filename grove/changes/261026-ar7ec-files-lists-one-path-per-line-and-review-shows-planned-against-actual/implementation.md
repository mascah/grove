# Implementation: Files lists one path per line, and review shows planned against actual

## Tasks

### T-01 — Files takes bullets and folders

Status: done

Evidence:
- `internal/doc/doc.go`: a `Files:` line with nothing after it, followed
  by bullets each holding one backticked path, is read as the bullet form;
  the comma line reads as before; a path may end in `/`. `doc.Names`
  reports whether a named file is a path or a folder above it. The
  malformed-line message names both forms and says a folder ends in `/`.
- `internal/cli/check.go`: `naming` (check's moved-file report and
  implement's changed-on-target line) and `changedSincePlan` (a cycle's
  staleness) match through `doc.Names`; a folder entry reports each
  changed file under it.
- Tasks guide and template teach one path per bullet and folders;
  `docs/reference.md` and `dev/walkthrough.md` say the same.
- `internal/cli/walkthrough_test.go` `TestMovedFiles` now names
  `ledger/` and `tests/test_cli.py` as bullets: moved since base, the
  implement line and the cycle's staleness all hold through the folder;
  then the wrapped comma line reports both files, and three malformed
  lines (no backticks, `../ledger/cli.py`, a comma line in a bullet) get
  the new message. `TestPropose` expects the bullet template.
- `just check` passed: 3921 Go lines, 1625 test lines, 7938 instruction
  words.

### T-02 — Implement records agreed files on Also changed lines

Status: not started

### T-03 — Review shows planned against actual

Status: not started

### T-04 — Trials read bullets and Also changed

Status: not started

## Scenarios

## Reconciliation
