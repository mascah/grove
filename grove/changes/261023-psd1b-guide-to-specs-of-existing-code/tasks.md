## T-01 — The tasks guide says when a change has no tasks

Result: the tasks guide says a change that only records behavior the code
already has has no tasks, as
[specs/propose/spec.md](specs/propose/spec.md) "A change that records
existing behavior has no tasks" says.

Files: `internal/instruct/schemas/default/documents/tasks/guide.md`

Interfaces: none.

Acceptance:
- `grove instructions propose` prints the sentence in the tasks guide.
- `just check` passes, the instruction-text word count included.

## T-02 — The guide, trialed on the fixture

Result: `docs/guides/spec-existing-code.md` is a recipe for proposing a
change that only records behavior the code already has, trialed on the
ledger fixture, and linked from the guides list, the docs site and
getting started.

Files: `docs/guides/spec-existing-code.md`, `docs/index.md`,
`docs/getting-started.md`, `mkdocs.yml`

Interfaces: none.

Constraints: the guide uses only Grove's existing features and links
the maintained specs rather than restating their rules (rule 8).

Acceptance:
- The guide opens as the other guides do, saying it is a tried recipe and
  not part of what Grove does. It says to put the agreed rules under
  ADDED, what the code does today under "Observed current behavior", and
  each difference the owner does not accept under the maintained spec's
  "Known discrepancies"; that the change has no tasks; and that it is
  built on its own branch with `/grove-implement <change-id>` on the
  draft, which verifies and archives it, so one merge approves and
  records it. It does not offer `/grove-archive` alone.
- On a fresh ledger fixture, in Claude Code, following only the guide: a
  change recording what `ledger edit` does is proposed with no tasks,
  with the traceback `ledger edit 99 --note x` ends in left as a
  difference the owner does not accept; `/grove-implement <change-id>`
  on the same branch builds nothing and archives it; `grove/specs/record-editing/spec.md` holds the rules and lists
  the traceback under "Known discrepancies"; `grove check` prints "No
  problems in grove/."; and `grove review` names the change as approved
  by the merge.
- The implementing pull request's description says what the trial did
  and showed, and the guide has what the trial taught.
- `docs/index.md`'s guides list and the `Guides` section of `mkdocs.yml`
  link the guide, and `docs/getting-started.md` links it after "A new
  project has no specs, and that is fine".
