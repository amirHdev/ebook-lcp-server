#!/bin/sh
set -eu

platform="${1:-swift}"
out_dir="${2:-/private/tmp/lcp-reader-demo-$platform}"
base_url="${BASE_URL:-http://127.0.0.1:8080}"
book_path="${BOOK_PATH:-examples/pride-and-prejudice/pride-and-prejudice.epub}"
publisher_username="${PUBLISHER_USERNAME:-publisher}"
publisher_password="${PUBLISHER_PASSWORD:-publisher}"
passphrase="${LCP_PASSPHRASE:-open-sesame}"
hint="${LCP_HINT:-demo}"
user_id="${LCP_USER_ID:-reader-01}"

case "$platform" in
  swift|android)
    ;;
  *)
    echo "platform must be swift or android" >&2
    exit 1
    ;;
esac

mkdir -p "$out_dir"

demo_json="$(
  go run ./cmd/lcpctl demo \
    --base-url "$base_url" \
    --username "$publisher_username" \
    --password "$publisher_password" \
    --file "$book_path" \
    --user-id "$user_id" \
    --passphrase "$passphrase" \
    --hint "$hint"
)"

printf '%s\n' "$demo_json" > "$out_dir/demo.json"

publication_id="$(printf '%s\n' "$demo_json" | jq -r '.publication_id')"
license_id="$(printf '%s\n' "$demo_json" | jq -r '.license_id')"
license_url="$(printf '%s\n' "$demo_json" | jq -r '.license_url')"

if [ -z "$publication_id" ] || [ "$publication_id" = "null" ] || [ -z "$license_id" ] || [ "$license_id" = "null" ] || [ -z "$license_url" ] || [ "$license_url" = "null" ]; then
  echo "demo command did not return publication_id, license_id, and license_url" >&2
  exit 1
fi

curl -fsS "$license_url" > "$out_dir/license.lcpl"
curl -fsS "$base_url/licenses/$license_id/status" > "$out_dir/license-status.json"

jq -n \
  --arg platform "$platform" \
  --arg baseURL "$base_url" \
  --arg publicationID "$publication_id" \
  --arg licenseID "$license_id" \
  --arg licenseURL "$license_url" \
  --arg passphrase "$passphrase" \
  --arg hint "$hint" \
  --arg userID "$user_id" \
  --arg sampleBook "$book_path" \
  --arg generatedAt "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  '{
    platform: $platform,
    generatedAt: $generatedAt,
    baseURL: $baseURL,
    sampleBook: $sampleBook,
    publicationID: $publicationID,
    licenseID: $licenseID,
    licenseURL: $licenseURL,
    readerCredentials: {
      passphrase: $passphrase,
      hint: $hint,
      userID: $userID
    },
    files: {
      license: "license.lcpl",
      status: "license-status.json",
      demo: "demo.json"
    }
  }' > "$out_dir/metadata.json"

printf '%s\n' "$out_dir"
