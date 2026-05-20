#!/usr/bin/env python3
"""Forward EPUB/PDF files from self-hosted library tools into the LCP server."""

from __future__ import annotations

import argparse
import base64
import dataclasses
import json
import os
import pathlib
import sys
import urllib.error
import urllib.request

SUPPORTED_SUFFIXES = {".epub", ".pdf"}


@dataclasses.dataclass(frozen=True)
class LCPConfig:
    base_url: str = "http://localhost:8080"
    username: str = "publisher"
    password: str = "publisher"
    two_factor: str = ""

    @classmethod
    def from_env(cls) -> "LCPConfig":
        return cls(
            base_url=os.getenv("LCP_BASE_URL", cls.base_url).rstrip("/"),
            username=os.getenv("LCP_USERNAME", cls.username),
            password=os.getenv("LCP_PASSWORD", cls.password),
            two_factor=os.getenv("LCP_2FA_CODE", cls.two_factor),
        )


def request_json(url: str, payload: dict, token: str | None = None) -> dict:
    body = json.dumps(payload).encode()
    req = urllib.request.Request(url, data=body, method="POST")
    req.add_header("Content-Type", "application/json")
    if token:
        req.add_header("Authorization", f"Bearer {token}")
    try:
        with urllib.request.urlopen(req, timeout=30) as response:
            return json.loads(response.read().decode())
    except urllib.error.HTTPError as exc:
        detail = exc.read().decode()
        raise SystemExit(f"{exc.code} {exc.reason}: {detail}") from exc


def login(base_url: str, username: str, password: str, two_factor: str) -> str:
    result = request_json(
        f"{base_url}/api/v1/auth/login",
        {"username": username, "password": password, "twoFactor": two_factor},
    )
    return result["token"]


def upload(base_url: str, token: str, path: pathlib.Path, title: str) -> dict:
    encoded = base64.b64encode(path.read_bytes()).decode()
    result = request_json(
        f"{base_url}/graphql",
        {
            "query": (
                "mutation UploadPublication($title: String!, $file: Upload!) "
                "{ uploadPublication(title: $title, file: $file) { id title } }"
            ),
            "variables": {"title": title, "file": encoded},
        },
        token,
    )
    if result.get("errors"):
        raise SystemExit(result["errors"][0]["message"])
    return result["data"]["uploadPublication"]

def forward_file(path: pathlib.Path, title: str | None, config: LCPConfig) -> dict:
    if not path.is_file():
        raise SystemExit(f"file not found: {path}")
    if path.suffix.lower() not in SUPPORTED_SUFFIXES:
        raise SystemExit(f"unsupported file type: {path.suffix}")
    token = login(config.base_url.rstrip("/"), config.username, config.password, config.two_factor)
    return upload(config.base_url.rstrip("/"), token, path, title or path.stem)


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("path", type=pathlib.Path, help="EPUB or PDF file to forward")
    parser.add_argument("--title", help="Catalog title override")
    env_config = LCPConfig.from_env()
    parser.add_argument("--base-url", default=env_config.base_url)
    parser.add_argument("--username", default=env_config.username)
    parser.add_argument("--password", default=env_config.password)
    parser.add_argument("--two-factor", default=env_config.two_factor)
    args = parser.parse_args()
    config = LCPConfig(args.base_url.rstrip("/"), args.username, args.password, args.two_factor)
    print(json.dumps(forward_file(args.path, args.title, config)))
    return 0


if __name__ == "__main__":
    sys.exit(main())
