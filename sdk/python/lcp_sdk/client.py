from __future__ import annotations

import base64
import json
from pathlib import Path
from typing import Any
from urllib import request
from urllib.error import HTTPError, URLError


class LcpClient:
    def __init__(
        self,
        base_url: str = "http://127.0.0.1:8080",
        token: str | None = None,
        two_factor_code: str | None = None,
    ) -> None:
        self.base_url = base_url.rstrip("/")
        self.token = token
        self.two_factor_code = two_factor_code

    def login(self, username: str, password: str, two_factor: str | None = None) -> dict[str, Any]:
        payload: dict[str, Any] = {
            "username": username,
            "password": password,
        }
        if two_factor:
            payload["twoFactor"] = two_factor
        response = self._request_json("/api/v1/auth/login", method="POST", body=payload)
        self.token = response["token"]
        return response

    def health(self) -> dict[str, Any]:
        return self._request_json("/healthz")

    def ready(self) -> dict[str, Any]:
        return self._request_json("/readyz")

    def status(self) -> dict[str, Any]:
        return self._request_json("/api/v1/lcp/status", headers=self._auth_headers())

    def upload_publication_from_file(self, file_path: str, title: str | None = None) -> dict[str, Any]:
        path = Path(file_path)
        return self.upload_publication(path.read_bytes(), title or path.stem)

    def upload_publication(self, data: bytes, title: str) -> dict[str, Any]:
        payload = {
            "query": "mutation UploadPublication($title: String!, $file: Upload!) { uploadPublication(title: $title, file: $file) { id title downloadURL } }",
            "variables": {
                "title": title,
                "file": base64.b64encode(data).decode("ascii"),
            },
        }
        return self._graphql("uploadPublication", payload)

    def create_license(
        self,
        publication_id: str,
        user_id: str,
        passphrase: str,
        hint: str,
    ) -> dict[str, Any]:
        payload = {
            "query": "mutation CreateLicense($publicationID: ID!, $userID: ID!, $passphrase: String!, $hint: String!) { createLicense(publicationID: $publicationID, userID: $userID, passphrase: $passphrase, hint: $hint) { id publicationID userID publicationURL passphrase hint } }",
            "variables": {
                "publicationID": publication_id,
                "userID": user_id,
                "passphrase": passphrase,
                "hint": hint,
            },
        }
        return self._graphql("createLicense", payload)

    def revoke_license(self, license_id: str) -> dict[str, str]:
        self._request_json(
            f"/api/v1/admin/licenses/{license_id}/revoke",
            method="POST",
            headers=self._admin_headers(),
        )
        return {"status": "revoked", "licenseID": license_id}

    def list_admin_licenses(self) -> dict[str, Any]:
        return self._request_json("/api/v1/admin/licenses", headers=self._admin_headers())

    def list_audit(self, limit: int | None = None) -> dict[str, Any]:
        path = "/api/v1/admin/audit"
        if limit is not None:
            path = f"{path}?limit={limit}"
        return self._request_json(path, headers=self._admin_headers())

    def download_lcpl(self, license_id: str) -> str:
        return self._request_text(f"/api/v1/licenses/{license_id}/lcpl", headers=self._auth_headers())

    def get_license_status_document(self, license_id: str) -> dict[str, Any]:
        return self._request_json(f"/licenses/{license_id}/status")

    def _graphql(self, field: str, payload: dict[str, Any]) -> dict[str, Any]:
        response = self._request_json(
            "/graphql",
            method="POST",
            headers={
                **self._auth_headers(),
                "Content-Type": "application/json",
            },
            body=payload,
        )
        errors = response.get("errors") or []
        if errors:
            raise RuntimeError(errors[0].get("message", "GraphQL request failed"))
        data = response.get("data") or {}
        if field not in data:
            raise RuntimeError(f"GraphQL response missing {field}")
        return data[field]

    def _auth_headers(self) -> dict[str, str]:
        if not self.token:
            raise RuntimeError("Client token is not set. Call login() first.")
        return {"Authorization": f"Bearer {self.token}"}

    def _admin_headers(self) -> dict[str, str]:
        headers = self._auth_headers()
        if self.two_factor_code:
            headers["X-2FA-Code"] = self.two_factor_code
        return headers

    def _request_json(
        self,
        path: str,
        *,
        method: str = "GET",
        headers: dict[str, str] | None = None,
        body: dict[str, Any] | None = None,
    ) -> dict[str, Any]:
        raw = self._request(path, method=method, headers=headers, body=body)
        return json.loads(raw)

    def _request_text(
        self,
        path: str,
        *,
        method: str = "GET",
        headers: dict[str, str] | None = None,
    ) -> str:
        return self._request(path, method=method, headers=headers)

    def _request(
        self,
        path: str,
        *,
        method: str,
        headers: dict[str, str] | None = None,
        body: dict[str, Any] | None = None,
    ) -> str:
        payload: bytes | None = None
        request_headers = dict(headers or {})
        if body is not None:
            payload = json.dumps(body).encode("utf-8")
            request_headers.setdefault("Content-Type", "application/json")

        req = request.Request(
            f"{self.base_url}{path}",
            data=payload,
            headers=request_headers,
            method=method,
        )

        try:
            with request.urlopen(req) as response:
                return response.read().decode("utf-8")
        except HTTPError as exc:
            details = exc.read().decode("utf-8", errors="replace")
            raise RuntimeError(f"{path} returned {exc.code} {exc.reason}: {details}") from exc
        except URLError as exc:
            raise RuntimeError(f"request to {path} failed: {exc.reason}") from exc
