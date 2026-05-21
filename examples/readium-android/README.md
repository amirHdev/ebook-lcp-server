# Readium Android Demo

This example shows how to generate a fresh `.lcpl` package from the local server and import it into a Readium Android-based reader app.

## Generate demo assets

From the repo root:

```bash
docker compose up -d postgres redis minio lcp-core lsd-core lcp-server
sh scripts/export-reader-demo.sh android /private/tmp/lcp-reader-demo-android
```

The output directory contains:

- `license.lcpl`
- `license-status.json`
- `metadata.json`
- `demo.json`

## Import flow

1. Open your Readium Android demo app or integration build.
2. Use the app flow that imports or fulfills an `.lcpl` license.
3. Select `/private/tmp/lcp-reader-demo-android/license.lcpl`.
4. Enter the passphrase from `/private/tmp/lcp-reader-demo-android/metadata.json`.
5. Confirm the protected publication downloads and opens.

## What to verify

- the `.lcpl` is accepted by the app
- the passphrase unlock succeeds
- fulfillment downloads the encrypted publication
- the reader opens the book without LCP errors
- the associated status document reflects the expected state

## Notes

- The repo-generated bundle is intended to be copied onto an emulator, test device, or into your app's import fixture flow.
- The repo can generate and validate the server-side assets locally. Final Android reader validation still depends on the app environment you use.
