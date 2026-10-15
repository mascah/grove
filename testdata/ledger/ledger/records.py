"""Records and the JSON file they are kept in."""

import json
from dataclasses import asdict, dataclass
from pathlib import Path


@dataclass
class Record:
    id: int
    date: str
    amount: float
    category: str
    note: str = ""


def load(path):
    """Read every record from the JSON file at path."""
    with Path(path).open(encoding="utf-8") as f:
        return [Record(**item) for item in json.load(f)]


def save(path, records):
    """Write records to the JSON file at path, one record per line."""
    lines = [json.dumps(asdict(r)) for r in records]
    Path(path).write_text("[\n  " + ",\n  ".join(lines) + "\n]\n", encoding="utf-8")


def next_id(records):
    return max((r.id for r in records), default=0) + 1


def find(records, record_id):
    """Return the record with record_id."""
    return next(r for r in records if r.id == record_id)
