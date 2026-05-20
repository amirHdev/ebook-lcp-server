#!/usr/bin/env python3
"""Poll a Calibre/calibre-web/Kavita library folder and forward new books to LCP."""

from __future__ import annotations

import argparse
import json
import pathlib
import sqlite3
import sys
import time

from lcp_forwarder import LCPConfig, SUPPORTED_SUFFIXES, forward_file


def ensure_state(path: pathlib.Path) -> sqlite3.Connection:
    path.parent.mkdir(parents=True, exist_ok=True)
    db = sqlite3.connect(path)
    db.execute(
        "CREATE TABLE IF NOT EXISTS forwarded_books (path TEXT PRIMARY KEY, mtime REAL NOT NULL, publication_id TEXT NOT NULL)"
    )
    return db


def seen(db: sqlite3.Connection, path: pathlib.Path, mtime: float) -> bool:
    row = db.execute(
        "SELECT mtime FROM forwarded_books WHERE path = ?", (str(path),)
    ).fetchone()
    return bool(row and row[0] == mtime)


def mark_seen(db: sqlite3.Connection, path: pathlib.Path, mtime: float, publication_id: str) -> None:
    db.execute(
        "INSERT OR REPLACE INTO forwarded_books (path, mtime, publication_id) VALUES (?, ?, ?)",
        (str(path), mtime, publication_id),
    )
    db.commit()


def iter_books(root: pathlib.Path) -> list[pathlib.Path]:
    return sorted(
        path
        for path in root.rglob("*")
        if path.is_file() and path.suffix.lower() in SUPPORTED_SUFFIXES
    )


def scan_once(root: pathlib.Path, state_path: pathlib.Path, config: LCPConfig) -> list[dict]:
    if not root.is_dir():
        raise SystemExit(f"library directory not found: {root}")
    db = ensure_state(state_path)
    forwarded: list[dict] = []
    try:
        for book in iter_books(root):
            mtime = book.stat().st_mtime
            if seen(db, book, mtime):
                continue
            publication = forward_file(book, book.stem, config)
            mark_seen(db, book, mtime, publication["id"])
            forwarded.append({"path": str(book), "publication": publication})
    finally:
        db.close()
    return forwarded


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("library", type=pathlib.Path, help="Library folder to scan")
    parser.add_argument("--state", type=pathlib.Path, default=pathlib.Path(".lcp-forwarded.sqlite3"))
    parser.add_argument("--interval", type=int, default=0, help="Polling interval in seconds; 0 scans once")
    args = parser.parse_args()

    config = LCPConfig.from_env()
    while True:
        forwarded = scan_once(args.library, args.state, config)
        for item in forwarded:
            print(json.dumps(item), flush=True)
        if args.interval <= 0:
            return 0
        time.sleep(args.interval)


if __name__ == "__main__":
    sys.exit(main())
