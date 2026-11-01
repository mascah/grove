# (a) A clear bug an existing requirement covers. The main line's end
# state: the record-editing spec and its archived change, with the fix and
# its test; then ledger/ restored, so the test fails again.

setup() {
  spec record-editing "Record editing" << 'EOF'
# Record editing

## Purpose

Editing a saved record keeps its category when the edit leaves it
unchanged, even when that category is retired.

## Boundaries

Adding records, and changing an older record's category, are outside this
spec.

## Requirements

### Requirement: Edit older records

When a user edits an existing record and leaves its category unchanged,
the system shall save the edit, even if that category is no longer
allowed for new records.

#### Scenario: Update the note on a record with a retired category

- GIVEN record 3, whose category `misc` was retired in 0.3.0
- WHEN the user runs `ledger edit 3 --note "Library fine, paid"`
- THEN the note is saved, the category stays `misc`, and the command
  exits 0

## Known discrepancies

None recorded.

## Open questions

- What happens when an older record's category is changed.
EOF
  local a=grove/changes/archive/261015-k3m9q-allow-editing-older-records
  put $a/.grove.yaml << 'EOF'
title: Allow editing older records
created: 2026-10-15T09:00:00Z
schema: default
base: @BASE@
EOF
  put $a/proposal.md << 'EOF'
# Allow editing older records

## Problem

`ledger edit 3 --note "Library fine, paid"` prints `error: category 'misc'
is not allowed` and exits 2, because `edit` validates the whole record
(`ledger/cli.py`).

## Intended outcome

An older record can be edited while its category is left unchanged.

## Decisions

- The owner: older records keep their category when it is unchanged.

## Open questions

- What happens when an older record's category is changed; it blocks no
  task.

## Related documents

- Requirements: [specs/record-editing/spec.md](specs/record-editing/spec.md)
- Tasks: [tasks.md](tasks.md)
EOF
  put $a/specs/record-editing/spec.md << 'EOF'
# Changes to record-editing

## ADDED Requirements

### Requirement: Edit older records

When a user edits an existing record and leaves its category unchanged,
the system shall save the edit, even if that category is no longer
allowed for new records.

#### Scenario: Update the note on a record with a retired category

- GIVEN record 3, whose category `misc` was retired in 0.3.0
- WHEN the user runs `ledger edit 3 --note "Library fine, paid"`
- THEN the note is saved, the category stays `misc`, and the command
  exits 0
EOF
  put $a/tasks.md << 'EOF'
## T-01 — Edit a record whose category is retired

Result: an older record's edit is saved when its category is unchanged.

Files: `ledger/cli.py`, `tests/test_cli.py`

Interfaces: none.

Acceptance:
- `ledger edit 3 --note "Library fine, paid"` exits 0 and keeps `misc`.
EOF
  put $a/implementation.md << 'EOF'
# Implementation: Allow editing older records

## Tasks

### T-01 — Edit a record whose category is retired

Status: done

`ledger/cli.py` skips the category rule when the category is unchanged;
`tests/test_cli.py` test_edit_older_record_note. `python3 -m unittest`:
passed.

## Reconciliation

- record-editing: applied to specs/record-editing/spec.md
EOF
  edit ledger/cli.py '            record = store.find(records, args.id)
' '            record = store.find(records, args.id)
            old = record.category
'
  edit ledger/cli.py '            validate(record)
            store.save(args.file, records)' '            try:
                validate(record)
            except ValueError as e:
                if record.category != old or "category" not in str(e):
                    raise
            store.save(args.file, records)'
  edit tests/test_cli.py '    def test_edit_current_record_note(self):' '    def test_edit_older_record_note(self):
        self.assertEqual(self.run_cli("edit", "3", "--note", "Library fine, paid"), (0, "edited 3\n"))
        self.assertEqual(store.find(store.load(self.file), 3).category, "misc")

    def test_edit_current_record_note(self):'
  commit "Allow editing older records"
  git checkout ':/^Set up Grove' -- ledger
  commit "Restore ledger"
}

runs=(
  "fix => /grove-propose Editing record 3's note fails again with \"category 'misc' is not allowed\""
  "fix => /grove-implement {id}"
)
answers=(
  "bugfix => Yes, use bugfix."
)

# good checks the end of a run: the new change on fix is a bugfix, archived,
# with no change spec, citing record-editing, and the fix works.
good() {
  git switch -q fix
  local d
  d=$(ls -d grove/changes/archive/"$id"-* grove/changes/"$id"-* 2> /dev/null | head -1)
  if [ -z "$id" ] || [ -z "$d" ]; then
    echo "the run proposed no change"
  else
    [[ $d == grove/changes/archive/* ]] || echo "the new change is not archived"
    grep -qx 'schema: bugfix' "$d/.grove.yaml" || echo "the new change does not have schema: bugfix"
    [ ! -e "$d/specs" ] || echo "the new change has a specs/ folder"
    grep -q record-editing "$d/proposal.md" || echo "the new change's proposal does not name record-editing"
  fi
  python3 -m unittest -q 2> /dev/null || echo "python3 -m unittest fails"
  # On a copy, so the check leaves records.json as the run did.
  cp records.json "$dir/records.json"
  python3 -m ledger --file "$dir/records.json" edit 3 --note x > /dev/null 2>&1 || echo "python3 -m ledger edit 3 --note x fails"
}
