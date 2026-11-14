# Remove needs_human

## Problem

A cycle file may name members under `needs_human`, each with the owner's
reason (`grove/specs/cycles/spec.md` at 999ab56, "Ask about work that needs
a person"). Plan and propose ask the owner about it for every member that
might need one, `grove claim` refuses a cycle that names one, `grove check`
reports a marked cycle that names one, and `grove serve` gives the key and
a "needs a person" chip.

The owner's analysis of 76 Claude Code sessions in another project, from
2026-10-24 to 2026-11-14 (a brief kept outside this repository), counted
140 structured questions; the most frequent mechanical kind was "does this
change need a person", asked 17 times, and the owner took the recommended
answer to plan's questions 6 times of 6. The same brief asks the cycle
title 5 times. Its direction: agents complete work without human input
and every implement run ends at a pull request for review, so a change
needing a person is settled in review, not by a gate. On `needs_human`
the owner wrote: "I think we need to revisit this concept."
`CONSTITUTION.md` rule 10, amended on 2026-11-14 (#158), no longer lists
it among `claim`'s gates.

## Intended outcome

A cycle file has no `needs_human`. Plan writes the cycle's title itself and
asks no one whether a change needs a person; propose, revising a member of
a cycle, does not ask it either. `grove claim` claims a cycle whatever its
members call for, the readiness `grove serve` gives has no such verdict,
and the project index names no member as needing a person. A cycle file
that still holds `needs_human` is reported by `grove check` as having an
unknown field, like any other; the project deletes the key.

## Scope

### In scope

- The `needs_human` field: parsing, the cycle template, the project index
  line, the JSON key `grove serve` gives and its chip on the page.
- `grove claim`'s refusal and `grove check`'s report of a marked cycle
  naming a member under `needs_human`.
- The question in the plan and propose instructions, plan's verdict
  obstacle for it, and the sentences of the implement instructions that
  name the field.
- Plan titling the cycle itself, with no question for it.
- The maintained specs, `docs/`, `dev/walkthrough.md`, the tests, and
  Grove's own cycle file that holds the field.

### Out of scope

- Implement no longer stopping at checkpoints, and the owner's and
  `--unattended` instructions saying one thing: a separate change (Part C1
  of the brief), which rewrites `requires/implement-owner.md` and
  `requires/implement-unattended.md`. This change only deletes their
  sentences naming `needs_human`.
- `unattended`, `depends_on` and the stale-plan gate: unchanged
  (`CONSTITUTION.md` rule 10).
- How `grove check` prints a file changed since base (Part C2).

## Decisions

The owner pre-answered this change's round in the brief: "The owner has
answered the questions here. Do not ask a round. Take the recommended
choice, write it down where the reviewer will see it, and go on." Each
decision below is the recommended answer taken under that direction.

1. Remove the field entirely, not deprecate it: the brief says "Remove
   `needs_human`: the field, the questions in `plan.md` and `propose.md`,
   `claim`'s refusal, the `check` diagnostic and the `serve` key."
2. A cycle file that still holds `needs_human` gets `grove check`'s
   existing unknown-field report, "unknown field "needs_human"; this file
   takes …", and no message of its own; `grove claim` reads the file as if
   the key were absent and claims the cycle when nothing else stops it.
   Reason: rule 8 gives no documented step that needs a dedicated message,
   and the unknown-field report already names the line.
3. Breaking, with this footer: "`needs_human` is gone from cycle files;
   `grove check` reports it as an unknown field; delete the key and its
   entries." The brief: "Breaking: say what a project with the field in
   its cycle files must do (memex has six)."
4. Plan writes the cycle's title from its members without asking, and asks
   about a one-line goal only by taking one the owner gives; with one
   change named, its title stays the cycle's unless the owner gives
   another. The brief: "Plan titles the cycle itself."
5. Instruction text drops on net (`CONSTITUTION.md` rule 12; the brief:
   "Part C must cut words on net. Do not raise the budget.").
6. Routine: the archived changes keep their text; `grove/cycles/
   261023-9kfxn-smaller-and-plainer.yaml`, done, loses its `needs_human`
   so `grove check` stays clean; tests change in place in
   `internal/cli/walkthrough_test.go`.

## Open questions

None.

## Related documents

- Requirements: [cycles](specs/cycles/spec.md), [runs](specs/runs/spec.md),
  [serve](specs/serve/spec.md), [propose](specs/propose/spec.md),
  [implement](specs/implement/spec.md)
- Tasks: [tasks.md](tasks.md)
