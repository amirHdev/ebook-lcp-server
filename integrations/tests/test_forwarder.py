import pathlib
import sys
import tempfile
import unittest
from unittest import mock

sys.path.insert(0, str(pathlib.Path(__file__).resolve().parents[1]))

import lcp_forwarder
import lcp_library_watcher


class ForwarderTests(unittest.TestCase):
    def test_forward_file_rejects_unsupported_suffix(self):
        with tempfile.TemporaryDirectory() as tmp:
            path = pathlib.Path(tmp) / "book.txt"
            path.write_text("not a book")
            with self.assertRaises(SystemExit):
                lcp_forwarder.forward_file(path, None, lcp_forwarder.LCPConfig())

    def test_forward_file_uploads_supported_book(self):
        with tempfile.TemporaryDirectory() as tmp:
            path = pathlib.Path(tmp) / "book.epub"
            path.write_bytes(b"book")
            with mock.patch.object(lcp_forwarder, "login", return_value="token") as login:
                with mock.patch.object(lcp_forwarder, "upload", return_value={"id": "pub1"}) as upload:
                    result = lcp_forwarder.forward_file(path, None, lcp_forwarder.LCPConfig())
            self.assertEqual(result["id"], "pub1")
            login.assert_called_once()
            upload.assert_called_once()


class LibraryWatcherTests(unittest.TestCase):
    def test_scan_once_marks_books_and_skips_second_scan(self):
        with tempfile.TemporaryDirectory() as tmp:
            root = pathlib.Path(tmp) / "library"
            root.mkdir()
            state = pathlib.Path(tmp) / "state.sqlite3"
            book = root / "Example.epub"
            book.write_bytes(b"book")
            with mock.patch.object(lcp_library_watcher, "forward_file", return_value={"id": "pub1"}) as forward:
                first = lcp_library_watcher.scan_once(root, state, lcp_forwarder.LCPConfig())
                second = lcp_library_watcher.scan_once(root, state, lcp_forwarder.LCPConfig())
            self.assertEqual(len(first), 1)
            self.assertEqual(second, [])
            forward.assert_called_once()


if __name__ == "__main__":
    unittest.main()
