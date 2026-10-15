"""Command-line interface: list, add, edit, report."""

import argparse
import sys

from ledger import records as store
from ledger.records import Record
from ledger.report import format_report, totals
from ledger.rules import validate


def parser():
    p = argparse.ArgumentParser(prog="ledger")
    p.add_argument("--file", default="records.json")
    sub = p.add_subparsers(dest="command", required=True)
    sub.add_parser("list")
    add = sub.add_parser("add")
    add.add_argument("date")
    add.add_argument("amount", type=float)
    add.add_argument("category")
    add.add_argument("note", nargs="?", default="")
    edit = sub.add_parser("edit")
    edit.add_argument("id", type=int)
    for field in ("date", "category", "note"):
        edit.add_argument("--" + field)
    edit.add_argument("--amount", type=float)
    report = sub.add_parser("report")
    report.add_argument("--month")
    return p


def main(argv=None):
    args = parser().parse_args(argv)
    records = store.load(args.file)
    try:
        if args.command == "list":
            for r in records:
                print(f"{r.id:>3}  {r.date}  {r.amount:>9.2f}  {r.category:<9}  {r.note}")
        elif args.command == "add":
            record = Record(store.next_id(records), args.date, args.amount, args.category, args.note)
            validate(record)
            store.save(args.file, records + [record])
            print(f"added {record.id}")
        elif args.command == "edit":
            record = store.find(records, args.id)
            for field in ("date", "amount", "category", "note"):
                if getattr(args, field) is not None:
                    setattr(record, field, getattr(args, field))
            validate(record)
            store.save(args.file, records)
            print(f"edited {record.id}")
        elif args.command == "report":
            print("\n".join(format_report(totals(records, args.month))))
    except ValueError as e:
        print(f"error: {e}", file=sys.stderr)
        return 2
    return 0
