"""Validation rules for records."""

CATEGORIES = ("food", "rent", "travel", "utilities")
MAX_NOTE = 60


def validate(record):
    """Raise ValueError if record breaks a rule."""
    if record.amount <= 0:
        raise ValueError("amount must be positive")
    if record.category not in CATEGORIES:
        raise ValueError(f"category {record.category!r} is not allowed")
    if len(record.note) > MAX_NOTE:
        raise ValueError(f"note is longer than {MAX_NOTE} characters")
