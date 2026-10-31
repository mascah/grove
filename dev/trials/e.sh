# (e) A fresh session picking up half-done work: T-01 built and recorded
# done in implementation.md, uncommitted, on delete-impl; T-02 not started.

# delete_test is T-01's test as the setup writes it, for good to find unchanged.
delete_test='    def test_delete_record(self):
        self.assertEqual(self.run_cli("delete", "6"), (0, "deleted 6\n"))
        self.assertEqual([r.id for r in store.load(self.file)], [1, 2, 3, 4, 5])'

setup() {
  local c=grove/changes/261015-e0000-add-ledger-delete-id
  id=261015-e0000
  put $c/.grove.yaml << 'EOF'
title: Add ledger delete ID
created: 2026-10-15T09:00:00Z
schema: default
base: @BASE@
EOF
  put $c/proposal.md << 'EOF'
# Add ledger delete ID

## Problem

There is no way to remove a record (README.md, "Planned for 0.4").

## Intended outcome

`ledger delete ID` removes the record and prints `deleted <id>`; an
unknown ID prints `error: no record <id>` and exits 2, as the README
says.

## Decisions

- The owner: the README and CHANGELOG say so in the same change.

## Open questions

None.

## Related documents

- Requirements: [specs/record-deletion/spec.md](specs/record-deletion/spec.md)
- Tasks: [tasks.md](tasks.md)
EOF
  put $c/specs/record-deletion/spec.md << 'EOF'
# Changes to record-deletion

There is no maintained spec for `record-deletion` at @BASE@.

## ADDED Requirements

### Requirement: Delete a record

When a user runs `ledger delete <id>` for an existing record, the system
shall remove it and print `deleted <id>`; if no record has that ID, it
shall print `error: no record <id>` and exit 2.

#### Scenario: Delete record 6

- GIVEN the six records of `records.json`
- WHEN the user runs `ledger delete 6`
- THEN it prints `deleted 6`, and `ledger list` shows records 1 to 5

## Observed current behavior

- `ledger` has no `delete` command (`ledger/cli.py` at @BASE@).
EOF
  put $c/tasks.md << 'EOF'
## T-01 — The delete command

Result: `ledger delete ID` works as the change spec says.

Files: `ledger/cli.py`, `tests/test_cli.py`

Interfaces: the `delete ID` command.

Acceptance:
- A test deletes record 6 and finds records 1 to 5 left.

## T-02 — README and CHANGELOG

Result: the README lists `delete` among the commands and no longer plans
it, and the CHANGELOG has it under a new 0.4.0 entry.

Files: `README.md`, `CHANGELOG.md`

Interfaces: none.

Acceptance:
- The README's usage shows `ledger delete ID`.
EOF
  commit "Propose ledger delete"
  git switch -q -c delete-impl
  edit ledger/cli.py '    report = sub.add_parser("report")' '    delete = sub.add_parser("delete")
    delete.add_argument("id", type=int)
    report = sub.add_parser("report")'
  edit ledger/cli.py '        elif args.command == "report":' '        elif args.command == "delete":
            record = next((r for r in records if r.id == args.id), None)
            if record is None:
                raise ValueError(f"no record {args.id}")
            store.save(args.file, [r for r in records if r is not record])
            print(f"deleted {record.id}")
        elif args.command == "report":'
  edit tests/test_cli.py '    def test_edit_current_record_note(self):' "$delete_test

    def test_edit_current_record_note(self):"
  put $c/implementation.md << 'EOF'
# Implementation: Add ledger delete ID

## Tasks

### T-01 — The delete command

Status: done

`ledger/cli.py` has `delete ID`; `tests/test_cli.py` test_delete_record.
`python3 -m unittest`: passed.

### T-02 — README and CHANGELOG

Status: not started

## Scenarios

## Reconciliation
EOF
}

runs=(
  "delete-impl => /grove-implement {id}"
)
answers=()

# good checks that the session reran the tests, left T-01's test as it was,
# and built T-02: done in implementation.md and the README's usage.
good() {
  local c=grove/changes/261015-e0000-add-ledger-delete-id
  grep -q '^TOOL: .*unittest' "$dir/transcript.txt" || echo "python3 -m unittest was not run in the session"
  python3 -c 'import sys; sys.exit(sys.argv[1] + "\n\n" not in open("tests/test_cli.py").read())' "$delete_test" ||
    echo "test_delete_record in tests/test_cli.py is not as the setup wrote it"
  grep -qx 'Status: done' <<< "$(task T-02 < "$c/implementation.md")" || echo "T-02 is not Status: done in implementation.md"
  # Usage only: the README already plans ledger delete under "Planned for 0.4".
  grep -q 'ledger delete' <<< "$(awk '/^## Usage/{on=1; next} /^## /{on=0} on' README.md)" || echo "README.md's usage does not show ledger delete"
  python3 -m unittest -q 2> /dev/null || echo "python3 -m unittest fails"
}
