from calibre.customize import InterfaceActionBase


class LCPSendPlugin(InterfaceActionBase):
    name = "Send to LCP Server"
    description = "Forward selected EPUB/PDF books to an ebook-lcp-server instance"
    supported_platforms = ["windows", "osx", "linux"]
    author = "ebook-lcp-server"
    version = (0, 1, 0)
    minimum_calibre_version = (6, 0, 0)
    actual_plugin = "calibre_plugins.lcp_send.ui:LCPSendAction"

    def is_customizable(self):
        return True

    def config_widget(self):
        from calibre_plugins.lcp_send.config import ConfigWidget

        return ConfigWidget()

    def save_settings(self, config_widget):
        config_widget.save()
