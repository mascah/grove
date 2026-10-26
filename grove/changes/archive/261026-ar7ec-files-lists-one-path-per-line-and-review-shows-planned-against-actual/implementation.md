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

Status: done

Evidence:
- `internal/instruct/requires/implement-owner.md`: an agreed file beyond
  `Files:` and `Delegated:` is edited and recorded as
  `` Also changed: `<path>` — <reason> `` under the task; a revision edits
  no `Files:` line, and one that changes a task's files or adds a task
  goes through `grove instructions propose <change-id>`.
- `internal/instruct/requires/implement-unattended.md`: writes no
  `Also changed:` line.
- `internal/instruct/schemas/default/documents/implementation/guide.md`
  describes the line, one per file.
- `docs/getting-started.md`, `docs/coordinating.md` and
  `dev/walkthrough.md` ("Delegated with me, ignored unattended" and a new
  "A file no task names, with me") say the same.
- `internal/cli/walkthrough_test.go` `TestImplementInstructions` checks
  that the owner's instructions name `Also changed:`, the frozen `Files:`
  line and the propose path, and that the unattended ones say to write no
  `Also changed:` line (and the owner's do not).
- The on-fixture acceptance (an agreed file ending as an `Also changed:`
  line, a new task through propose) needs an agent session with the owner:
  not run here, see Scenarios.
- `just check` passed.

### T-03 — Review shows planned against actual

Status: done

Evidence:
- Rule 3 names `Also changed:` among what Grove reads in `review`
  (`CONSTITUTION.md` at 2c6de73), so the constraint holds.
- `internal/doc/doc.go`: `Task.Also` holds the path and reason of each
  `Also changed:` line under the task's heading in `implementation.md`,
  not in the project data.
- `internal/cli/review.go`: while the branch holds a change open, or
  archived where the target does not, `grove review` lists before the
  documents each file `<target>...HEAD` changes outside the root, as
  `planned <id>:T-NN`, `agreed <id>:T-NN: <reason>` or `named nowhere`,
  then the path.
- `docs/reference.md` and `dev/walkthrough.md` (step 6 of the main line
  and "Planned against actual") show the list.
- `internal/cli/walkthrough_test.go` `TestPlanned` covers the three marks
  and a folder entry, before and after archiving, and a branch holding no
  change.
- `just check` passed: 3982 Go lines, 1663 test lines.

### T-04 — Trials read bullets and Also changed

Status: done

Evidence:
- `dev/trials/run.sh`: `allowed` adds the paths of every
  `implementation.md`'s `Also changed:` lines to those of the `Files:` and
  `Delegated:` lines (whose bullets the awk range already took in), and
  `named` counts a path under an allowed folder ending in `/`.
- `dev/trials.md` describes the count so.
- Exercised `allowed` and `named`, taken from the script, on a scratch
  folder whose `tasks.md` lists `ledger/records.py` and `tests/fixtures/`
  as bullets and whose `implementation.md` has
  `` Also changed: `README.md` ``: `ledger/records.py`,
  `tests/fixtures/a.json` and `README.md` count as inside, `Cargo.lock`
  and a path only an `Acceptance:` bullet names as outside. `bash -n`
  passes. A full `just trials` run, which drives Claude Code and Codex, was
  not run in this session.

## Scenarios

- implement "Planned, agreed and named nowhere": `internal/cli/walkthrough_test.go` TestPlanned, passed
- implement "After archiving": TestPlanned, passed
- implement "No change on the branch": TestPlanned, passed
- implement "Feedback changes a requirement mid-build": unverified; it needs an agent session with the owner (the instruction text is checked by TestImplementInstructions)
- implement "The owner does not agree": unverified; it needs an agent session with the owner
- implement "A revision that needs a new task": unverified; it needs an agent session with the owner (TestImplementInstructions checks the owner's instructions name the propose path and the frozen `Files:` line)
- implement "A named file changed since base": TestMovedFiles, passed
- implement "A file no task names needs a change": unverified; it needs an agent session with the owner (TestImplementInstructions checks the instructions name `Also changed:`)
- implement "A member of a cycle": TestMovedFiles, passed
- implement "A delegated path, unattended": unverified; it needs an unattended agent run (TestImplementInstructions checks the unattended instructions say to write no `Also changed:` line)
- propose "A path outside the repository": TestMovedFiles (`../ledger/cli.py`), passed
- propose "A folder": TestMovedFiles (`ledger/` checks clean) and TestPlanned (`tests/fixtures/`), passed
- propose "A task that fixes a function": unverified; it is how a proposing agent writes a task
- propose "A task that updates documentation": unverified; it is how a proposing agent writes a task
- propose "One path per bullet": TestMovedFiles reads the bullet form clean and TestPropose checks the template teaches it, passed; that an agent writes it is unverified
- propose "A comma line still reads": TestMovedFiles (a wrapped comma line reports both files), passed
- propose "A file changed since base": TestMovedFiles, passed
- propose "A file under a named folder": TestMovedFiles (`ledger/report.py` under `ledger/`), passed
- propose "A member of a cycle": TestMovedFiles, passed
- propose "A file this change's own branch touched": TestMovedFiles, passed

## Reconciliation
- implement: applied to grove/specs/implement/spec.md
- propose: applied to grove/specs/propose/spec.md
