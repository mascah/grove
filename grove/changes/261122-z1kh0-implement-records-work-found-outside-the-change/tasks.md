<!-- Each task is a section headed "## T-NN — title". IDs are never
renumbered or reused. -->

## T-01 — Implement writes and reports `Found:` lines

Result: `grove instructions implement` has the agent record work outside
the change as a `Found:` line under the task it was building, change
nothing for it, and list those lines in its report after what the
reviewer must judge, as "Record work found outside the change" says.

Files:
- `internal/instruct/requires/implement.md`
- `internal/instruct/schemas/default/documents/implementation/guide.md`
- `internal/instruct/schemas/default/documents/implementation/template.md`
- `internal/cli/walkthrough_test.go`

Delegated: test layout and the exact wording of the instruction text.

Constraints:
- `Found:` is distinct from `Check:` and `Also changed:` in the text: a
  file the task needs or makes false is still `Also changed:`, and a wrong
  approach is still a revision.
- No Go outside the test changes. Grove reads nothing from the line.
- The added instruction text keeps `just check` under the 9,000-word
  budget.

Acceptance:
- The output of `grove instructions implement <change-id>` holds the rule:
  write `Found: <what is wrong or missing, where, and how it was seen>`
  under the task, change nothing for it, and list it in the report after
  every `Check:` line, blocked task's question and unverified scenario,
  and before what was done. The implement test checks for that text.
- The implementation template's comment and guide name the `Found:` line
  beside `Check:` and `Also changed:`.
- `grove changes --json` and `grove review` on a change whose
  `implementation.md` has a `Found:` line print what they print without it.
- `just check` passes.

## T-02 — Docs, the tracker guide and the walkthrough

Result: the docs say what a `Found:` line is and where the report lists
it, the push-to-a-tracker guide opens one issue per `Found:` line, and
the walkthrough has a step that sees one written.

Files:
- `docs/getting-started.md`
- `docs/reference.md`
- `docs/concepts.md`
- `docs/guides/push-to-a-tracker.md`
- `dev/walkthrough.md`

Delegated: the exact wording, and where in each page the text goes.

Acceptance:
- `docs/getting-started.md`, where it describes what implement writes and
  reports, says that work outside the change is written as a `Found:` line
  and not done, and that the report lists those lines after the `Check:`
  lines, blocked questions and unverified scenarios.
- `docs/reference.md`, where it lists the lines of `implementation.md`,
  names `Found:` and says Grove does not read it.
- `docs/concepts.md`'s glossary has "Found line".
- `docs/guides/push-to-a-tracker.md`'s prompt and commands open one issue
  for each `Found:` line in the target's `implementation.md` files,
  labelled `grove-found`, whose body is the line and a link to its change,
  ending with `<!-- grove-found:<change-id>:<T-NN>:<n> -->`. The prompt
  finds an existing issue by that marker and never closes one.
- `dev/walkthrough.md` has a step on the fixture: approve a change that
  adds `ledger delete ID`, run `/grove-implement <change-id>`, and see the
  session, having seen `ledger edit 999` print a traceback, leave `edit`
  unchanged and write a `Found:` line under T-01, which its report lists
  after the `Check:` lines.
