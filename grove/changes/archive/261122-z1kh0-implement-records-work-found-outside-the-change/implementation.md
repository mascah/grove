# Implementation: Implement records work it finds outside the change

<!-- Progress on change 261122-z1kh0, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Implement writes and reports `Found:` lines

Status: done

- `internal/instruct/requires/implement.md`: a new rule after the
  `Also changed:` one: work the change does not ask for and no task needs
  is left as it is and written as
  `Found: <what is wrong or missing, where, and how it was seen>` under the
  task being built; a file a task needs or makes false stays
  `Also changed:`, a wrong approach a revision; Grove reads nothing from
  it. The finish line lists every `Found:` line after the `Check:` lines,
  blocked questions and unverified scenarios, before what was done.
- The implementation guide names `Found:` beside `Also changed:`; the
  template's comment shows `Also changed:` and `Found:` examples beside
  `Check:` and `Answer:`.
- `internal/cli/walkthrough_test.go`: `TestImplementInstructions` checks
  the rule, the report order and the template example in
  `grove instructions implement`; `TestArchive` adds a `Found:` line to an
  archived change on the branch and checks that `grove review`'s readout
  and `grove changes --json` stay the same.
- `bin/grove instructions implement 261122-z1kh0` prints the rule and the
  finish line. `go test ./internal/cli/ -run 'TestImplement|TestArchive'`
  passed. Instruction text is 8,382 words of 9,000.

Check: `grove review` prints the branch's diff after its readout, and the diff shows any line added to a file, a `Found:` line included. The test compares the readout above the diff, and `grove changes --json` without the time it was read and the branch tip's hash, date and subject, which a new commit always changes. The `Found:` line itself changes neither.
Check: `internal/cli/walkthrough_test.go` changed on main after this change was proposed (#189, #194: explore and propose tests). Those edits are to other tests; this task's edits sit beside them on the file as main holds it.

Also changed: `internal/instruct/requires/implement-cycle.md` — a cycle is built by `/grove-implement` too, and its report must list `Found:` lines as the requirement says: each change's subagent now reports them, and the finish line lists them after the `Check:` lines and blocked questions; `TestImplementInstructions` checks both (from review).

### T-02 — Docs, the tracker guide and the walkthrough

Status: done

- `docs/getting-started.md`: after the `Also changed:` paragraph, work
  outside the change is left alone and written as a `Found:` line, with
  an example; the report lists those lines after the `Check:` lines,
  blocked questions and unverified scenarios, before what was done; the
  tracker guide turns each into an issue; Grove reads nothing from them.
- `docs/reference.md`, "What a schema cannot change", beside
  `Also changed:`: names `Found:` and says Grove does not read it.
- `docs/concepts.md`: a glossary entry, "Found line".
- `docs/guides/push-to-a-tracker.md`: the prompt opens one issue per
  `Found:` line in the target's `implementation.md` files, labelled
  `grove-found`, its body the line, a link to its change's folder and
  `<!-- grove-found:<change-id>:<T-NN>:<n> -->`; it finds an existing one
  by that marker, open or closed, and never closes one. The commands show
  `git grep` finding the lines and `gh issue create` with the marker; "What
  to expect" says the issues stay open until the owner closes them.
- `dev/walkthrough.md`: a step "Work found outside the change" beside "A
  file no task names": approve `ledger delete ID`, run
  `/grove-implement`, see `ledger edit 999` print a traceback, `edit`
  unchanged, a `Found:` line under T-01, listed after the `Check:` lines.
- On the fixture, `python3 -m ledger edit 999 --note x` prints a
  `StopIteration` traceback, as the step and the example say.

Check: the walkthrough step was written, not run: running it needs a live `claude` session on the fixture. Its last sentence also covers the "Nothing found" scenario.
Check: the tracker guide's prompt and commands were not run against GitHub; `git grep -n '^Found:' main -- 'grove/changes/**/implementation.md'` was run on this repository and finds this change's line once committed. The issue title is what the line says is wrong, which the change left open.
Check: `docs/getting-started.md`, `docs/reference.md`, `docs/concepts.md` and `dev/walkthrough.md` changed on main after this change was proposed (#189, #194: explore and propose wording). None of those edits touch implement's text; this task's paragraphs sit beside them as main holds them.

## Scenarios

- implement "A bug next to the work": unverified; the agent's behavior shows only in a live session, the walkthrough's "Work found outside the change" step, not run. `TestImplementInstructions` checks the instructions hold the rule and the report order, passed.
- implement "A file the task makes false": unverified; agent behavior, the walkthrough's "A file no task names" step, not run. `TestImplementInstructions` checks the instructions keep such a file an `Also changed:` line, passed.
- implement "Nothing found": unverified; agent behavior, the last sentence of the walkthrough's "Work found outside the change" step, not run.
- implement "Grove does not read it": internal/cli/walkthrough_test.go TestArchive, "Grove reads nothing from a Found: line", passed.

## Reconciliation

- implement: applied to specs/implement/spec.md
