from __future__ import annotations

import json
import os

from lcp_sdk import LcpClient


def main() -> None:
    client = LcpClient(
        base_url=os.getenv("LCP_BASE_URL", "http://127.0.0.1:8080"),
        two_factor_code=os.getenv("LCP_2FA_CODE", "123456"),
    )
    username = os.getenv("LCP_USERNAME", "admin")
    password = os.getenv("LCP_PASSWORD", "admin")
    book_path = os.getenv("BOOK_PATH", "examples/pride-and-prejudice/pride-and-prejudice.epub")

    login = client.login(username, password, os.getenv("LCP_2FA_CODE", "123456"))
    health = client.health()
    ready = client.ready()
    publication = client.upload_publication_from_file(book_path, "Python SDK Smoke Book")
    license_doc = client.create_license(publication["id"], "python-sdk-reader-01", "open-sesame", "demo")
    lcpl = client.download_lcpl(license_doc["id"])
    status = client.get_license_status_document(license_doc["id"])
    licenses = client.list_admin_licenses()

    print(
        json.dumps(
            {
                "loginRole": login.get("role"),
                "health": health,
                "ready": ready,
                "publicationID": publication["id"],
                "licenseID": license_doc["id"],
                "lcplBytes": len(lcpl),
                "status": status.get("status"),
                "adminLicenseCount": len(licenses.get("licenses", [])),
            },
            indent=2,
        )
    )


if __name__ == "__main__":
    main()
