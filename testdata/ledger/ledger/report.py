"""Totals by category."""

from collections import defaultdict


def totals(records, month=None):
    """Return {category: total} for records, optionally only one YYYY-MM month."""
    sums = defaultdict(float)
    for r in records:
        if month is None or r.date.startswith(month + "-"):
            sums[r.category] += r.amount
    return dict(sorted(sums.items()))


def format_report(sums):
    """Return the report as lines of text."""
    if not sums:
        return ["No records."]
    width = max(len(c) for c in sums)
    lines = [f"{c:<{width}}  {amount:>9.2f}" for c, amount in sums.items()]
    lines.append(f"{'total':<{width}}  {sum(sums.values()):>9.2f}")
    return lines
