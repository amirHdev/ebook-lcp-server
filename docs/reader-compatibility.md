# Reader compatibility

This project emits Readium LCP licenses and protected publications. The server side is meant to stay reader-neutral, but compatibility still needs to be proven reader by reader.

| Reader | Status | Notes |
| --- | --- | --- |
| Thorium Reader | Verified | Local end-to-end flow tested with a generated `.lcpl`, signed download URL, and passphrase unlock. |
| Readium Swift Toolkit | Demo kit ready | The repo now exports a repeatable `.lcpl`, status document, and metadata bundle for a Readium Swift import flow. Final simulator or device validation still depends on the local app setup. |
| Readium Kotlin Toolkit / Android | Demo kit ready | The repo now exports a repeatable `.lcpl`, status document, and metadata bundle for an Android import flow. Final emulator or device validation still depends on the local app setup. |

## Thorium Reader

The current local flow has been checked with Thorium Reader:

1. Start the local stack with `docker compose up --build`.
2. Run `sh scripts/demo-local.sh`.
3. Import the generated `.lcpl` into Thorium Reader.
4. Use the license passphrase from the demo flow.

The local Compose setup uses `127.0.0.1` in generated public links because Thorium rejects `localhost` URLs while importing a license.

## Readium Swift

Use the shared export script to generate a fresh import bundle:

```bash
sh scripts/export-reader-demo.sh swift /private/tmp/lcp-reader-demo-swift
```

Then follow [examples/readium-swift/README.md](/Users/mehrbod/projects/golang/lcp/examples/readium-swift/README.md) for the import and verification flow.

## Android

Use the shared export script to generate a fresh import bundle:

```bash
sh scripts/export-reader-demo.sh android /private/tmp/lcp-reader-demo-android
```

Then follow [examples/readium-android/README.md](/Users/mehrbod/projects/golang/lcp/examples/readium-android/README.md) for the import and verification flow.
