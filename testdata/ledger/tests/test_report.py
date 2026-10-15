import unittest
from pathlib import Path

from ledger import records as store
from ledger.report import totals

RECORDS = Path(__file__).resolve().parent.parent / "records.json"


class ReportTest(unittest.TestCase):
    def test_all_records(self):
        sums = totals(store.load(RECORDS))
        self.assertEqual(sums["misc"], 60.7)
        self.assertEqual(sorted(sums), ["food", "misc", "rent", "travel", "utilities"])

    def test_one_month(self):
        self.assertEqual(totals(store.load(RECORDS), "2026-09"), {"utilities": 88.4})


if __name__ == "__main__":
    unittest.main()
