from calibre.gui2 import error_dialog
from calibre.utils.config import JSONConfig
from qt.core import QFormLayout, QLineEdit, QWidget

prefs = JSONConfig("plugins/lcp_send")
prefs.defaults["base_url"] = "http://localhost:8080"
prefs.defaults["username"] = "publisher"
prefs.defaults["password"] = "publisher"
prefs.defaults["two_factor"] = ""
prefs.defaults["forwarder"] = "lcp_forwarder.py"


class ConfigWidget(QWidget):
    def __init__(self):
        QWidget.__init__(self)
        layout = QFormLayout()
        self.setLayout(layout)

        self.base_url = QLineEdit(prefs["base_url"])
        self.username = QLineEdit(prefs["username"])
        self.password = QLineEdit(prefs["password"])
        self.password.setEchoMode(QLineEdit.EchoMode.Password)
        self.two_factor = QLineEdit(prefs["two_factor"])
        self.forwarder = QLineEdit(prefs["forwarder"])

        layout.addRow("LCP base URL", self.base_url)
        layout.addRow("Username", self.username)
        layout.addRow("Password", self.password)
        layout.addRow("2FA code", self.two_factor)
        layout.addRow("Forwarder script path", self.forwarder)

    def save(self):
        if not self.base_url.text().strip():
            error_dialog(self, "Missing LCP URL", "LCP base URL is required", show=True)
            return
        prefs["base_url"] = self.base_url.text().strip()
        prefs["username"] = self.username.text().strip()
        prefs["password"] = self.password.text()
        prefs["two_factor"] = self.two_factor.text().strip()
        prefs["forwarder"] = self.forwarder.text().strip() or "lcp_forwarder.py"
