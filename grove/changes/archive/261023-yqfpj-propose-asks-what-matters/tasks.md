## T-01 — Propose orders its round, offers to split and writes delegation

Result: `/grove-propose` behaves as
[specs/propose/spec.md](specs/propose/spec.md) "Ask a round before
writing", "Offer to split a change" and "Delegate choices in a task" say,
and the tasks guide and template show the optional `Delegated:` line.

Files: `internal/instruct/requires/propose.md`,
`internal/instruct/schemas/default/skills/propose.md`,
`internal/instruct/schemas/default/documents/tasks/guide.md`,
`internal/instruct/schemas/default/documents/tasks/template.md`,
`docs/concepts.md`, `docs/getting-started.md`, `docs/reference.md`,
`dev/walkthrough.md`, `internal/cli/walkthrough_test.go`

Interfaces: a task's optional `Delegated:` line, free text read by people
and agents; Grove parses nothing new (rule 3).

Constraints: the round still asks every decision the agent would
otherwise assume; the `Files:` line and how Grove reads it stay as they
are.

Acceptance:
- `docs/getting-started.md`, step 3, says the round comes in order of
  consequence with routine choices on one last line, that propose offers
  to split a big idea into changes linked with `depends_on`, and that a
  task may carry a `Delegated:` line you agreed to; `docs/concepts.md`'s
  `tasks.md` row and `docs/reference.md`'s description of a task mention
  the optional line (rule 8).
- On the ledger fixture, in Claude Code: `/grove-propose Editing an
  unknown ID prints an error` asks the exit status before where the check
  goes, and lists the test's name and place on a last line of defaults;
  when the round recommends delegating test layout and
  `tests/test_records.py` and the owner says "yes", T-01 has that
  `Delegated:` line, and otherwise none;
  `/grove-propose Add ledger delete ID, and an error for an unknown ID`
  offers two changes and, on "yes", writes both on the branch, the delete
  change naming the other in `depends_on`. `dev/walkthrough.md` records
  both.
- `internal/cli/walkthrough_test.go` checks that `instructions propose`
  prints the order of the round, the offer to split and the
  `Delegated:` line in the tasks guide, and that `grove check` reports no
  problem on a task with a `Delegated:` line right after its `Files:`
  line, and no moved file for a path only that line names.

## T-02 — Implement honors delegation and stops on a discovery, and ignores delegation unattended

Result: implement behaves as
[specs/implement/spec.md](specs/implement/spec.md) "Honor a task's
delegation", "Stop on a discovery that breaks the approach" and "Ignore a
task's delegation unattended" say.

Files: `internal/instruct/requires/implement-owner.md`,
`internal/instruct/requires/implement-unattended.md`,
`docs/getting-started.md`, `dev/walkthrough.md`,
`internal/cli/walkthrough_test.go`

Interfaces: none; it is instruction text.

Constraints: the owner's rules go in the owner's part and the unattended
rule in the unattended part, the files change 261023-g0a3h creates; the
shared `requires/implement.md` is left as that change leaves it.

Acceptance:
- `docs/getting-started.md`, step 5, says implement makes the choices and
  changes the paths a task delegates without asking, stops a task when it
  finds the approach wrong and says what revision the change needs, and
  that `/grove-implement --unattended` ignores `Delegated:` (rule 8).
- On the ledger fixture, in Claude Code: a change whose T-01 delegates
  `tests/test_records.py` is built by `/grove-implement`, which creates
  that file without asking, and by `/grove-implement <cycle-id>
  --unattended` on a cycle of it, which leaves the file uncreated and
  names it in T-01's evidence. `dev/walkthrough.md` records both.
- `internal/cli/walkthrough_test.go` checks that `instructions implement`
  prints the delegation rule and the stop and not the rule to ignore
  `Delegated:`, and that `instructions implement --unattended` prints the
  rule to ignore it and not the other two.
