#!/bin/sh
set -eu

base_dir="${1:-/private/tmp/lcp-reader-demos}"

swift_dir="$base_dir/swift"
android_dir="$base_dir/android"

sh scripts/export-reader-demo.sh swift "$swift_dir" >/dev/null
sh scripts/export-reader-demo.sh android "$android_dir" >/dev/null

for dir in "$swift_dir" "$android_dir"; do
  test -f "$dir/license.lcpl"
  test -f "$dir/license-status.json"
  test -f "$dir/metadata.json"
  jq -e '.licenseID and .licenseURL and .readerCredentials.passphrase' "$dir/metadata.json" >/dev/null
done

printf '%s\n' "$base_dir"
