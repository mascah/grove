# Implementation: Implement edits a file a task needs and says why

<!-- Progress on change 261101-mg8hd, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Implement edits a file a task needs and says why

Status: done

- `internal/instruct/requires/implement-owner.md` and
  `implement-unattended.md`: both say to make the choices and change the
  paths a `Delegated:` line names and, beyond them and the `Files:` lines,
  to edit any other file a task needs changed or makes false, without
  asking and without blocking the task, writing
  `` Also changed: `<path>` — <why> `` under it, and never `grove.yaml`
  for it. The owner text no longer says to ask first; the unattended text
  no longer says such a file stays as it is, to ignore `Delegated:` lines,
  or to write no `Also changed:` line. The bullet saying not to edit the
  change's files other than `implementation.md` stays.
- `internal/instruct/schemas/default/documents/implementation/guide.md`:
  an `Also changed:` line records a file the implementing agent changed,
  and why. `.../tasks/guide.md`: `Delegated:` is honored in both modes,
  and a file a task needs or makes false is recorded with `Also changed:`.
- `internal/cli/review.go`: `mark` labels such a file
  `also changed <change-id>:T-NN: <why>`.
- `internal/cli/walkthrough_test.go`: `TestImplementInstructions` checks
  the new text in both modes and the absence of the old;
  `TestPropose` checks the tasks guide; `TestPlanned` checks
  `also changed` and that `Cargo.lock` is still `named nowhere`.
- `docs/coordinating.md`, `docs/getting-started.md`, `docs/reference.md`:
  both modes and the label, with no "once you agree".
- `dev/walkthrough.md`: "Delegated with me, ignored unattended" is now
  "Delegated, with me and unattended", expecting `grove run <cycle-id>` to
  create `tests/test_records.py`; "A file no task names, with me" expects
  the edit without asking and review's `also changed`; "Planned against
  actual" shows `also changed`.
- Would have changed `internal/instruct/requires/implement.md`, which no
  `Files:` line names: its shared line "Change only the files the tasks
  name on their `Files:` lines." still prints above both modes' new
  bullet, which reads as the exception ("Beyond them and the `Files:`
  lines, ..."). Left as it is.
- `just check`: passed (4031 Go lines, 1697 test lines, 8048 words).

### T-02 — Trial case (c) expects the edit in both modes

Status: done

- `dev/trials/c.sh`: no scripted answers; `good` checks, on both
  `unknown-id-impl` and `unknown-id-unattended`, that `ledger/records.py`
  changed, that T-01 is `Status: done` with an
  `` Also changed: `ledger/records.py` `` line (in the change folder or,
  once archived, in `changes/archive/`), and that `python3 -m unittest`
  passes. The check that the file stayed unchanged before an answer is
  gone.
- `dev/trials.md` case (c): "Owner's answers: none"; "A good run" and
  "Checked by `good`" say both modes edit `ledger/records.py` without
  asking and write an `Also changed:` line for it.
- Checked on case (c)'s setup (`SETUP_ONLY=1 dev/trials/run.sh c`, clean),
  with end states written by hand: with both branches changing
  `ledger/records.py` and T-01 done with the line (the unattended one
  archived), `good` printed nothing; with `unknown-id-unattended` leaving
  the file, T-01 `Status: blocked` and no line, it printed one line for
  each of the three checks.
- The trials themselves (`just trials`) were not run in this session.

## Scenarios

- implement "A delegated path": internal/cli/walkthrough_test.go TestImplementInstructions shows both modes' text says to change the paths a `Delegated:` line names, passed; the agent creating the file is the walkthrough step "Delegated, with me and unattended", not run here: unverified
- implement "A file nothing names": internal/cli/walkthrough_test.go TestImplementInstructions shows the owner text says to edit it without asking and write the line, passed; an agent doing so is the walkthrough step "A file no task names, with me", not run here: unverified
- implement "Unattended": internal/cli/walkthrough_test.go TestImplementInstructions (unattended text) and TestPlanned (`also changed` with the reason), passed; an agent doing it unattended was not run here: unverified
- implement "Unattended, a file the task needs": trial case (c)'s `good` checks it; the trial was not run here, only `good` on hand-written end states: unverified
- implement "A draft, before and after building": internal/cli/walkthrough_test.go TestApprove (future work) and TestOneMerge (with its code), passed
- implement "Planned, also changed and named nowhere": internal/cli/walkthrough_test.go TestPlanned, passed

## Reconciliation

- implement: applied to grove/specs/implement/spec.md
