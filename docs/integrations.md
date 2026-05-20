# Integrations

The integration package has three entry points:

- `integrations/lcp_forwarder.py` forwards one EPUB/PDF file to the LCP server.
- `integrations/lcp_library_watcher.py` scans a shared library folder and forwards new EPUB/PDF files.
- `integrations/calibre_plugin` builds into a Calibre plugin that forwards selected books from the desktop app.

All entry points use the same LCP credentials:

| Variable | Default |
| --- | --- |
| `LCP_BASE_URL` | `http://localhost:8080` |
| `LCP_USERNAME` | `publisher` |
| `LCP_PASSWORD` | `publisher` |
| `LCP_2FA_CODE` | empty |

## Single-file forwarder

```bash
python3 integrations/lcp_forwarder.py "/library/book.epub" --title "Book Title"
```

Use this from any scriptable import hook.

## Calibre desktop plugin

Build the plugin ZIP:

```bash
sh scripts/build-calibre-plugin.sh
```

Install `dist/lcp-send-calibre-plugin.zip` in Calibre from:

```text
Preferences -> Plugins -> Load plugin from file
```

Configure the plugin with the LCP server URL, publisher credentials, and the local path to `integrations/lcp_forwarder.py`. Select one or more EPUB/PDF books, then run `Send to LCP Server`.

## calibre-web and Kavita sidecar

calibre-web and Kavita deployments usually keep imported books on a mounted library directory. Run the watcher as a sidecar against that shared directory:

```bash
export LCP_BASE_URL=http://localhost:8080
export LCP_USERNAME=publisher
export LCP_PASSWORD=publisher
python3 integrations/lcp_library_watcher.py /library --state /data/lcp-forwarded.sqlite3 --interval 30
```

For a one-shot import pass, omit `--interval`:

```bash
python3 integrations/lcp_library_watcher.py /library --state /data/lcp-forwarded.sqlite3
```

The watcher stores forwarded file paths and modification times in SQLite so repeated scans do not upload the same book again.

## Docker sidecar example

| Variable | Default |
| --- | --- |
| `LIBRARY_PATH` | `/library` |
| `STATE_PATH` | `/data/lcp-forwarded.sqlite3` |
| `SCAN_INTERVAL` | `30` |

```yaml
services:
  lcp-library-watcher:
    image: python:3.12-alpine
    working_dir: /app
    command: >
      sh -c "python3 integrations/lcp_library_watcher.py
      $${LIBRARY_PATH:-/library}
      --state $${STATE_PATH:-/data/lcp-forwarded.sqlite3}
      --interval $${SCAN_INTERVAL:-30}"
    environment:
      LCP_BASE_URL: http://lcp-server:8080
      LCP_USERNAME: publisher
      LCP_PASSWORD: publisher
    volumes:
      - ./:/app:ro
      - /path/to/your/library:/library:ro
      - lcp-watcher-data:/data
```

Mount the same library folder used by calibre-web or Kavita at `/library`.
