# ledger

A tiny personal ledger. Records live in `records.json`, one entry per
expense: an id, a date, an amount, a category and a short note.

## Usage

    python3 -m ledger list
    python3 -m ledger add 2026-10-01 12.50 food "Lunch"
    python3 -m ledger edit 4 --note "Groceries, weekly"
    python3 -m ledger report --month 2026-09

Every command takes `--file PATH` to use another records file.

## Rules

Every record uses one of the categories food, rent, travel, utilities.
Amounts are positive. Notes are at most 60 characters.

## Planned for 0.4

- `ledger delete ID` removes a record.
- An unknown ID prints an error instead of a traceback, in `edit` and, once
  it exists, `delete`.
- Monthly budgets in `report`.

## Development

    python3 -m unittest

See CHANGELOG.md for what changed between versions.
