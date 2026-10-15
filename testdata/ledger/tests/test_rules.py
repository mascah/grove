import unittest

from ledger.records import Record
from ledger.rules import validate


class RulesTest(unittest.TestCase):
    def test_valid_record(self):
        validate(Record(1, "2026-10-01", 12.5, "food", "Lunch"))

    def test_amount_must_be_positive(self):
        with self.assertRaisesRegex(ValueError, "amount must be positive"):
            validate(Record(1, "2026-10-01", 0, "food"))

    def test_retired_category(self):
        with self.assertRaisesRegex(ValueError, "category 'misc' is not allowed"):
            validate(Record(1, "2026-10-01", 5, "misc"))

    def test_long_note(self):
        with self.assertRaisesRegex(ValueError, "longer than 60"):
            validate(Record(1, "2026-10-01", 5, "food", "x" * 61))


if __name__ == "__main__":
    unittest.main()
