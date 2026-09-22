---
id: "G-076"
type: work
title: "Filter grove list by status"
status: proposed
created: "2026-09-22T15:01:38Z"
updated: "2026-09-22T15:02:11Z"
kind: feature
size: small
relates_to: ["G-039"]
---

## Outcome

`grove list` can show only the records in a chosen status, so a person or an
agent finds open work in one command instead of piping the whole table
through `grep`. The owner asked for this on 2026-09-22, in the `/grove-shape`
session that wrote this record, as the real change for the
[G-039](G-039-interactive-loop.md) interactive-loop trial (its plan G-075,
step 1, exists only on branch `worktree-G-039` until that work is
integrated, so it is not linked here).

## Constraints

Observed at main `5f07c94`, 2026-09-22:

- `grove list --status proposed` exits 2 with `unknown option --status`.
  `list` takes no option of its own; it prints every record of the checkout
  as `ID TYPE STATUS TITLE`, a page's status as `-`
  (`internal/cli/cli.go:156`). Options are one flat table in `parseArgs`
  in the same file, with `--include` already repeatable.
- The record model gives each type a closed status vocabulary
  (`docs/record-model.md`, "Statuses"), and the loader rejects a value
  outside it. This checkout holds 7 distinct values across 74 records; the
  same word can belong to several types (`proposed` is work, question,
  decision and term; `current` is plan and review).
- The shaping and work guides read the unfiltered `list` at the start of a
  session, and the board has its own fixed status columns
  (`internal/tui/model.go:37`); neither depends on this change.

The owner decided on 2026-09-22, in that session, among three offered scopes:
**status only, and `list` without the flag keeps printing everything.** Out
of scope, each a separate proposal if wanted: a `--type` filter, hiding
finished work by default, filtering the board, and `--json` for `list`.

Proposed design, binding nobody: a repeatable `--status VALUE` (also
`--status=VALUE`), accepted only after `list`; a record is shown when its
status equals any given value exactly; a value outside the union of the
model's status vocabulary is a usage error (exit 2) rather than an empty
result, since the vocabulary is closed; matching nothing prints the header
only and exits 0. The usage text, the README's `list` line and the record
model's `list` contract say so.

## Acceptance

1. `grove list --status proposed` prints the same header and rows, in the
   same order and format, as `grove list` minus every row whose status is not
   `proposed`; `--status active --status review` unions the two.
2. `grove list` without the option prints exactly what it printed before.
3. A value outside the model's vocabulary, an empty value, and `--status` on
   any other command are refused as usage errors (exit 2, message naming the
   option) before any project is read.
4. Filtering to a status no record holds prints the header and exits 0.
5. Usage text, README and the record model describe the option and agree
   with the behavior; `go vet`, `gofmt -l`, `grove check` and the suite pass
   within the AGENTS.md budget.
6. The owner runs it on this repository and judges from the output that it
   answers "what is open" without further piping.

## Next

Assign: `/grove-work G-076` in a fresh session, per G-075 step 2. No plan
record is expected for a change of this size; say so in Next when preparing.
