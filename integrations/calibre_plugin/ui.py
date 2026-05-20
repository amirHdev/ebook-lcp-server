import os
import subprocess

from calibre.gui2 import error_dialog, info_dialog
from calibre.gui2.actions import InterfaceAction
from calibre_plugins.lcp_send.config import prefs
from qt.core import QAction


class LCPSendAction(InterfaceAction):
    name = "Send to LCP Server"
    action_spec = ("Send to LCP Server", None, "Forward selected EPUB/PDF books to LCP", None)

    def genesis(self):
        self.qaction.triggered.connect(self.forward_selected)

    def forward_selected(self):
        db = self.gui.current_db
        book_ids = self.gui.library_view.get_selected_ids()
        if not book_ids:
            error_dialog(self.gui, "No books selected", "Select at least one EPUB or PDF.", show=True)
            return

        sent = 0
        errors = []
        for book_id in book_ids:
            title = db.title(book_id, index_is_id=True)
            for fmt in ("EPUB", "PDF"):
                path = db.format_abspath(book_id, fmt, index_is_id=True)
                if path and os.path.exists(path):
                    try:
                        self.run_forwarder(path, title)
                        sent += 1
                    except Exception as exc:
                        errors.append(f"{title}: {exc}")
                    break

        if errors:
            error_dialog(self.gui, "LCP forwarding errors", "\n".join(errors), show=True)
            return
        info_dialog(self.gui, "LCP forwarding complete", f"Forwarded {sent} book(s).", show=True)

    def run_forwarder(self, path, title):
        env = os.environ.copy()
        env["LCP_BASE_URL"] = prefs["base_url"]
        env["LCP_USERNAME"] = prefs["username"]
        env["LCP_PASSWORD"] = prefs["password"]
        env["LCP_2FA_CODE"] = prefs["two_factor"]
        subprocess.run(
            ["python3", prefs["forwarder"], path, "--title", title],
            env=env,
            check=True,
            capture_output=True,
            text=True,
        )
