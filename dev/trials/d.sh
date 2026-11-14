# (d) A stale proposal after related code changes: the budgets task keeps
# the report in category order, and main then lists categories by amount.

setup() {
  local c=grove/changes/261015-d0000-monthly-budgets-in-report
  id=261015-d0000
  put $c/.grove.yaml << 'EOF'
title: Monthly budgets in report
created: 2026-10-15T09:00:00Z
schema: default
base: @BASE@
EOF
  put $c/proposal.md << 'EOF'
# Monthly budgets in report

## Problem

`ledger report --month YYYY-MM` shows each category's total but no
budget to compare it with (README.md, "Planned for 0.4").

## Intended outcome

With a `budgets.json` beside `records.json`, `report` shows each
category's budget after its total, in the same category order as today.

## Scope

### In scope

- Reading `budgets.json` and showing a budget beside each total.

### Out of scope

- Setting budgets from the command line.

## Decisions

- The owner: budgets live in `budgets.json`, a map of category to amount.
- The owner: the report keeps its order, one line per category.

## Open questions

None.

## Related documents

- Requirements: [specs/reporting/spec.md](specs/reporting/spec.md)
- Tasks: [tasks.md](tasks.md)
EOF
  put $c/specs/reporting/spec.md << 'EOF'
# Changes to reporting

There is no maintained spec for `reporting` at @BASE@.

## ADDED Requirements

### Requirement: Show a budget beside each total

Where `budgets.json` gives a category a budget, the system shall show
that budget after the category's total in `ledger report`.

#### Scenario: Food under budget

- GIVEN `budgets.json` holds `{"food": 200}`
- WHEN the user runs `ledger report --month 2026-01`
- THEN the food line shows 63.75 and then 200.00

## Observed current behavior

- `format_report` in `ledger/report.py` lists categories in the order
  `totals` gives them, sorted by name (at @BASE@).
EOF
  put $c/tasks.md << 'EOF'
## T-01 — Budgets in the report

Result: the report shows each category's budget after its total, as the
change spec says.

Files: `ledger/report.py`, `ledger/cli.py`, `tests/test_report.py`

Interfaces: `format_report(sums, budgets=None)` in `ledger/report.py`.

Acceptance:
- With `budgets.json`, each line with a budget shows it after the total,
  the lines in category order as today.
- Without it, the report is unchanged.
EOF
  commit "Propose monthly budgets in report"
  edit ledger/report.py 'for c, amount in sums.items()]' 'for c, amount in sorted(sums.items(), key=lambda kv: -kv[1])]'
  commit "List report categories by amount, largest first"
}

check_ok="tasks.md: ledger/report.py changed on main since base"
runs=(
  "budgets-impl => /grove-implement {id}"
)
answers=(
  "report\.py|changed|since base|stale|order => Build on the file as it is now, listing by amount."
)

# good checks that the agent did not ask about the change to
# ledger/report.py, so never got the scripted answer, built on the file,
# wrote a Check: line under T-01 naming it, and the tests pass.
good() {
  local c=grove/changes/261015-d0000-monthly-budgets-in-report t01
  [ -z "$(answered "Build on the file as it is now")" ] || echo "the agent asked about the change to ledger/report.py"
  ! git diff --quiet main budgets-impl -- ledger/report.py || echo "ledger/report.py is unchanged"
  t01=$( (git show "budgets-impl:$c/implementation.md" || git show "budgets-impl:grove/changes/archive/${c#grove/changes/}/implementation.md") 2> /dev/null | task T-01)
  grep -q '^Check:.*report' <<< "$t01" || echo "no Check: line under T-01 names the change to ledger/report.py"
  python3 -m unittest -q 2> /dev/null || echo "python3 -m unittest fails"
}
