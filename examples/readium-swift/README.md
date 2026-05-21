# Readium Swift Demo

This example shows how to generate a fresh `.lcpl` package from the local server and import it into a Readium Swift-based reader app.

## Generate demo assets

From the repo root:

```bash
docker compose up -d postgres redis minio lcp-core lsd-core lcp-server
sh scripts/export-reader-demo.sh swift /private/tmp/lcp-reader-demo-swift
```

The output directory contains:

- `license.lcpl`
- `license-status.json`
- `metadata.json`
- `demo.json`

## Import flow

1. Open your Readium Swift demo app on device or simulator.
2. Use the app's import or fulfillment entry point for an `.lcpl` file.
3. Select `/private/tmp/lcp-reader-demo-swift/license.lcpl`.
4. When prompted, use the passphrase from `/private/tmp/lcp-reader-demo-swift/metadata.json`.
5. Confirm the app downloads and opens the protected publication.

## What to verify

- the `.lcpl` imports successfully
- the passphrase unlock works
- the publication downloads fully
- the book opens without license errors
- the status document at `license-status.json` reports the expected state

## Notes

- The local Compose setup uses `127.0.0.1` in generated public links because that has been the safest loopback host for reader import flows.
- The repo generates the server-side assets. Final device or simulator validation still depends on your local Readium Swift app setup.
