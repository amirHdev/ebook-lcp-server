#!/bin/sh
set -eu

base_url="${BASE_URL:-http://localhost:8080}"
out_dir="${1:-certification-packet}"
book_path="${BOOK_PATH:-examples/pride-and-prejudice/pride-and-prejudice.epub}"
admin_username="${ADMIN_USERNAME:-admin}"
admin_password="${ADMIN_PASSWORD:-admin}"
admin_2fa="${ADMIN_2FA_CODE:-123456}"

mkdir -p "$out_dir/responses" "$out_dir/config" "$out_dir/docs"

admin_token="$(
  curl -fsS "$base_url/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d "$(jq -cn \
      --arg username "$admin_username" \
      --arg password "$admin_password" \
      --arg twoFactor "$admin_2fa" \
      '{username: $username, password: $password, twoFactor: $twoFactor}')" |
    jq -r .token
)"

cp docs/certification-blueprint.md "$out_dir/docs/"
cp docs/reader-compatibility.md "$out_dir/docs/"
cp docs/acceptance-checklist.md "$out_dir/docs/"
cp docker-compose.yml "$out_dir/config/docker-compose.yml"

if command -v docker >/dev/null 2>&1; then
  if docker compose config > "$out_dir/config/docker-compose.resolved.yaml" 2>/dev/null; then
    :
  else
    rm -f "$out_dir/config/docker-compose.resolved.yaml"
  fi
fi

curl -fsS "$base_url/healthz" > "$out_dir/responses/healthz.json"
curl -fsS "$base_url/readyz" > "$out_dir/responses/readyz.json"
curl -fsS "$base_url/api/v1/lcp/status" -H "Authorization: Bearer $admin_token" > "$out_dir/responses/lcp-status.json"
curl -fsS "$base_url/api/v1/admin/licenses" -H "Authorization: Bearer $admin_token" -H "X-2FA-Code: $admin_2fa" > "$out_dir/responses/admin-licenses.json"
curl -fsS "$base_url/api/v1/admin/audit" -H "Authorization: Bearer $admin_token" -H "X-2FA-Code: $admin_2fa" > "$out_dir/responses/admin-audit.json"

demo_output="$(BASE_URL="$base_url" BOOK_PATH="$book_path" sh scripts/demo-local.sh)"
printf '%s\n' "$demo_output" > "$out_dir/responses/demo.txt"

publication_id="$(printf '%s\n' "$demo_output" | awk -F= '/^publication_id=/{print $2}')"
license_id="$(printf '%s\n' "$demo_output" | awk -F= '/^license_id=/{print $2}')"
license_url="$(printf '%s\n' "$demo_output" | awk -F= '/^license_url=/{print $2}')"

if [ -z "$publication_id" ] || [ -z "$license_id" ] || [ -z "$license_url" ]; then
  echo "demo flow did not return publication_id, license_id, and license_url" >&2
  exit 1
fi

curl -fsS "$license_url" > "$out_dir/responses/license.lcpl"
curl -fsS "$base_url/licenses/$license_id/status" > "$out_dir/responses/license-status.json"
curl -fsSL "$base_url/publications/$publication_id/content" -H "Authorization: Bearer $admin_token" -o "$out_dir/responses/publication-content.bin"
sh scripts/certification-smoke.sh "$out_dir/certification-report.json" >/dev/null

cert_subject=""
cert_issuer=""
cert_expires=""
cert_fingerprint=""
cert_path="${LCP_CERTIFICATE:-${CERT_PATH:-}}"

if [ -n "$cert_path" ] && [ -f "$cert_path" ] && command -v openssl >/dev/null 2>&1; then
  cert_subject="$(openssl x509 -in "$cert_path" -noout -subject | sed 's/^subject=//')"
  cert_issuer="$(openssl x509 -in "$cert_path" -noout -issuer | sed 's/^issuer=//')"
  cert_expires="$(openssl x509 -in "$cert_path" -noout -enddate | sed 's/^notAfter=//')"
  cert_fingerprint="$(openssl x509 -in "$cert_path" -noout -fingerprint -sha256 | sed 's/^sha256 Fingerprint=//')"
fi

git_commit=""
if command -v git >/dev/null 2>&1; then
  git_commit="$(git rev-parse HEAD 2>/dev/null || true)"
fi

jq -n \
  --arg generatedAt "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  --arg baseURL "$base_url" \
  --arg bookPath "$book_path" \
  --arg publicationID "$publication_id" \
  --arg licenseID "$license_id" \
  --arg licenseURL "$license_url" \
  --arg gitCommit "$git_commit" \
  --arg certPath "$cert_path" \
  --arg certSubject "$cert_subject" \
  --arg certIssuer "$cert_issuer" \
  --arg certExpires "$cert_expires" \
  --arg certFingerprint "$cert_fingerprint" \
  '{
    generatedAt: $generatedAt,
    baseURL: $baseURL,
    gitCommit: $gitCommit,
    sampleBook: $bookPath,
    demo: {
      publicationID: $publicationID,
      licenseID: $licenseID,
      licenseURL: $licenseURL
    },
    certificate: {
      path: $certPath,
      subject: $certSubject,
      issuer: $certIssuer,
      expiresAt: $certExpires,
      sha256Fingerprint: $certFingerprint
    },
    files: {
      smokeReport: "certification-report.json",
      health: "responses/healthz.json",
      ready: "responses/readyz.json",
      lcpStatus: "responses/lcp-status.json",
      adminLicenses: "responses/admin-licenses.json",
      adminAudit: "responses/admin-audit.json",
      demo: "responses/demo.txt",
      lcpl: "responses/license.lcpl",
      licenseStatus: "responses/license-status.json",
      encryptedPublication: "responses/publication-content.bin"
    }
  }' > "$out_dir/manifest.json"

printf '%s\n' "$out_dir"
