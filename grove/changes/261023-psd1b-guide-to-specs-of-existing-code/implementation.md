# Implementation: A guide to writing specs of existing code

## Tasks

### T-01 — The tasks guide says when a change has no tasks

Status: done

Added "A change that only records behavior the code already has has no
tasks." as the tasks guide's second bullet
(`internal/instruct/schemas/default/documents/tasks/guide.md`).
`grove instructions propose` prints it; `just check` passed, 7,711
instruction words.

### T-02 — The guide, trialed on the fixture

Status: done

Wrote `docs/guides/spec-existing-code.md`, linked from `docs/index.md`'s
guides list, the `Guides` section of `mkdocs.yml`, and
`docs/getting-started.md` after "A new project has no specs, and that is
fine".

Trial, 2026-10-24, Claude Code headless (`claude -p`, the tool list of
`dev/trials/run.sh`), on a fresh ledger fixture on branch `record-edit`,
the implementer answering as owner from the guide:

- `/grove-propose Record what \`ledger edit\` does today, ...` wrote
  change 261024-3tcgj with five requirements under ADDED, what the code
  does under "Observed current behavior", and a `tasks.md` holding a note
  and no `T-NN` task. Its first answer for the `edit 99` traceback was
  to record it only as observed; told to keep a rule and the traceback as
  a known discrepancy, it did. `grove check`: "No problems in grove/."
- `/grove-implement 261024-3tcgj`, in a new session on the same branch,
  built nothing, checked each scenario by running the code once, wrote
  `grove/specs/ledger-edit/spec.md` with the traceback under "Known
  discrepancies", and archived the change.
- `grove check`: "No problems in grove/."; `grove review`: "Merging
  record-edit into main approves change 261024-3tcgj with its code."
- The spec is `ledger-edit`, not the acceptance's `record-editing`: the
  agent chose the ID and the owner's answer took its default. The guide
  now says to name the spec in the round.

The guide's "What to expect" holds what the trial taught.

## Scenarios

## Reconciliation
