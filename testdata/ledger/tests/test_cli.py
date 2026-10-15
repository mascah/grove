import contextlib
import io
import shutil
import tempfile
import unittest
from pathlib import Path

from ledger import records as store
from ledger.cli import main

RECORDS = Path(__file__).resolve().parent.parent / "records.json"


class CliTest(unittest.TestCase):
    def setUp(self):
        self.dir = tempfile.mkdtemp()
        self.file = str(Path(self.dir) / "records.json")
        shutil.copy(RECORDS, self.file)

    def tearDown(self):
        shutil.rmtree(self.dir)

    def run_cli(self, *args):
        out, err = io.StringIO(), io.StringIO()
        with contextlib.redirect_stdout(out), contextlib.redirect_stderr(err):
            code = main(["--file", self.file, *args])
        return code, out.getvalue() + err.getvalue()

    def test_add_valid_record(self):
        self.assertEqual(self.run_cli("add", "2026-10-01", "12.5", "food", "Lunch"), (0, "added 7\n"))

    def test_add_retired_category_is_rejected(self):
        code, out = self.run_cli("add", "2026-10-01", "5", "misc")
        self.assertEqual(code, 2)
        self.assertIn("category 'misc' is not allowed", out)

    def test_edit_current_record_note(self):
        self.assertEqual(self.run_cli("edit", "4", "--note", "Groceries, weekly"), (0, "edited 4\n"))
        self.assertEqual(store.find(store.load(self.file), 4).note, "Groceries, weekly")

    def test_list(self):
        code, out = self.run_cli("list")
        self.assertEqual(code, 0)
        self.assertEqual(len(out.splitlines()), 6)


if __name__ == "__main__":
    unittest.main()
