# Ledger fixture

The project every walkthrough step and every test runs on (rule 17). A
small Python package, about 290 lines, standard library only, tests with
`python3 -m unittest`. Described here; built in phase 1 under
`testdata/ledger/`.

`just fixture DIR` copies it to `DIR`, which must not exist, runs
`git init -b main` and makes one commit, "Ledger 0.3.0", as "Fixture Owner".
To reset, delete `DIR` and run it again.

## Files

| File | Lines | Contents |
| --- | --- | --- |
| `pyproject.toml` | 12 | Name `ledger`, version 0.3.0, script `ledger = "ledger.cli:main"`, no dependencies. |
| `README.md` | 31 | What the ledger is, the four commands, "Every record uses one of the categories food, rent, travel, utilities.", and "Planned for 0.4": `ledger delete ID`, an error instead of a traceback for an unknown ID, and monthly budgets in `report`. |
| `CHANGELOG.md` | 14 | 0.1.0 records and `list`; 0.2.0 `report`; 0.3.0 (2026-01-05) "Retired the `misc` category. New records must use one of food, rent, travel, utilities." |
| `.gitignore` | 2 | `__pycache__/` |
| `records.json` | 8 | The six records below, one per line. |
| `ledger/__init__.py`, `__main__.py` | 4 | Version; `python3 -m ledger` runs the CLI. |
| `ledger/records.py` | 35 | `Record` dataclass (id, date, amount, category, note); `load`, `save`; `find(records, id)` uses `next(...)` with no default. |
| `ledger/rules.py` | 25 | `CATEGORIES = ("food", "rent", "travel", "utilities")`; `validate(record)` raises `ValueError` for a non-positive amount, a category not in the list ("category 'misc' is not allowed"), or a note over 60 characters. |
| `ledger/report.py` | 30 | Totals by category, optionally for one month; `misc` appears as its own line. |
| `ledger/cli.py` | 55 | `list`, `add DATE AMOUNT CATEGORY [NOTE]`, `edit ID [--date] [--amount] [--category] [--note]`, `report [--month YYYY-MM]`; `--file` defaults to `records.json`. `add` and `edit` validate the whole record; a `ValueError` prints `error: <message>` and exits 2. |
| `tests/test_rules.py` | 30 | Valid record; bad amount; retired category; long note. |
| `tests/test_cli.py` | 35 | Add a valid record; add with `misc` is rejected; edit a current record's note; list. Each on a temporary copy of `records.json`. |
| `tests/test_report.py` | 15 | Totals for all records and for one month. |

## Records

| id | date | amount | category | note |
| --- | --- | --- | --- | --- |
| 1 | 2025-02-03 | 42.50 | misc | Bike repair |
| 2 | 2025-04-01 | 1200.00 | rent | April rent |
| 3 | 2025-06-07 | 18.20 | misc | Library fine |
| 4 | 2026-01-12 | 63.75 | food | Groceries |
| 5 | 2026-03-02 | 310.00 | travel | Train to Calgary |
| 6 | 2026-09-30 | 88.40 | utilities | Power bill |

Records 1 and 3 predate 0.3.0.

## What the walkthrough needs from it

- **First change.** `ledger edit 3 --note "Library fine, paid"` prints
  `error: category 'misc' is not allowed` and exits 2, because `edit`
  validates the whole record. The README ("every record") and the CHANGELOG
  ("new records") disagree on whether that is intended, so the exploration has
  something to ask. No test covers editing an older record.
- **Second change.** `ledger edit 99 --note x` ends in an uncaught
  `StopIteration` traceback. Observed behavior, agreed nowhere.
- All tests pass as shipped, so a failing test after implementation is the
  implementer's, not the fixture's.

## What the coordinating walkthrough needs

The README's "Planned for 0.4" gives three changes to coordinate, after the
main line has run.

- **Delete waits on the unknown-ID fix.** `delete` would find its record
  with `find`, which ends in `StopIteration` for an unknown ID, as
  `ledger edit 99 --note x` does today. So the delete change waits on the
  change that makes an unknown ID print an error.
- **Budgets wait on nothing.** They touch `report`, which neither of the
  other two changes.
- **Two changes to one requirement.** After the main line, the
  `record-editing` spec holds the requirement for editing an older record.
  Changing an older record's category, and printing what an edit saved, both
  modify it, so whichever is archived second is proposed from text the
  first rewrote.
