#!/bin/sh
set -eu

base_url="${BASE_URL:-http://localhost:8080}"
out="${1:-acceptance-report.json}"
admin_username="${ADMIN_USERNAME:-admin}"
admin_password="${ADMIN_PASSWORD:-admin}"
admin_2fa="${ADMIN_2FA_CODE:-123456}"
publisher_username="${PUBLISHER_USERNAME:-publisher}"
publisher_password="${PUBLISHER_PASSWORD:-publisher}"
tenant_id="${TENANT_ID:-default}"
guest_api_key="${GUEST_API_KEY:-guest-smoke-key}"
book_path="${BOOK_PATH:-examples/pride-and-prejudice/pride-and-prejudice.epub}"

tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT INT TERM

json_request() {
  method="$1"
  path="$2"
  body="${3:-}"
  token="${4:-}"
  api_key="${5:-}"
  two_factor="${6:-}"
  response_file="$7"

  set -- curl -sS -o "$response_file" -w "%{http_code}" -X "$method" "$base_url$path" -H "Content-Type: application/json"
  if [ -n "$token" ]; then
    set -- "$@" -H "Authorization: Bearer $token"
  fi
  if [ -n "$api_key" ]; then
    set -- "$@" -H "X-API-Key: $api_key"
  fi
  if [ -n "$two_factor" ]; then
    set -- "$@" -H "X-2FA-Code: $two_factor"
  fi
  if [ -n "$body" ]; then
    set -- "$@" --data "$body"
  fi

  "$@"
}

login() {
  username="$1"
  password="$2"
  two_factor="$3"
  response_file="$tmpdir/login-$username.json"
  payload="$(jq -cn \
    --arg username "$username" \
    --arg password "$password" \
    --arg twoFactor "$two_factor" \
    '{username: $username, password: $password, twoFactor: $twoFactor}')"
  code="$(json_request POST /api/v1/auth/login "$payload" "" "" "" "$response_file")"
  if [ "$code" != "200" ]; then
    echo "login failed for $username: $(cat "$response_file")" >&2
    exit 1
  fi
  jq -r '.token' "$response_file"
}

admin_token="$(login "$admin_username" "$admin_password" "$admin_2fa")"
publisher_token="$(login "$publisher_username" "$publisher_password" "")"

tenant_file="$tmpdir/tenant.json"
tenant_code="$(json_request GET "/api/v1/admin/tenants/$tenant_id" "" "$admin_token" "" "$admin_2fa" "$tenant_file")"
if [ "$tenant_code" != "200" ]; then
  echo "failed to load tenant $tenant_id: $(cat "$tenant_file")" >&2
  exit 1
fi

updated_tenant="$tmpdir/tenant-updated.json"
jq \
  --arg tenantID "$tenant_id" \
  --arg apiKey "$guest_api_key" \
  '.id = $tenantID
  | .name = (if (.name // "") == "" then $tenantID else .name end)
  | .apiKeys = ((.apiKeys // []) | map(select(.key != $apiKey)) + [{key: $apiKey, subject: "guest-smoke", role: "guest"}])' \
  "$tenant_file" > "$updated_tenant"

tenant_put_file="$tmpdir/tenant-put.json"
tenant_put_code="$(json_request PUT "/api/v1/admin/tenants/$tenant_id" "$(cat "$updated_tenant")" "$admin_token" "" "$admin_2fa" "$tenant_put_file")"
if [ "$tenant_put_code" != "200" ]; then
  echo "failed to update tenant $tenant_id: $(cat "$tenant_put_file")" >&2
  exit 1
fi

book_b64_file="$tmpdir/acceptance-book.b64"
process_body_file="$tmpdir/process-body.json"
base64 < "$book_path" | tr -d '\n' > "$book_b64_file"
jq -cn --rawfile file "$book_b64_file" \
  '{title: "Acceptance Smoke", file: $file}' > "$process_body_file"
process_file="$tmpdir/process.json"
process_code="$(
  curl -sS -o "$process_file" -w "%{http_code}" \
    -X POST "$base_url/api/v1/lcp/process" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $publisher_token" \
    --data-binary "@$process_body_file"
)"

status_publisher_file="$tmpdir/status-publisher.json"
status_publisher_code="$(json_request GET /api/v1/lcp/status "" "$publisher_token" "" "" "$status_publisher_file")"

status_guest_file="$tmpdir/status-guest.json"
status_guest_code="$(json_request GET /api/v1/lcp/status "" "" "$guest_api_key" "" "$status_guest_file")"

metrics_publisher_file="$tmpdir/metrics-publisher.json"
metrics_publisher_code="$(json_request GET /api/v1/admin/metrics "" "$publisher_token" "" "" "$metrics_publisher_file")"

metrics_admin_no_2fa_file="$tmpdir/metrics-admin-no-2fa.json"
metrics_admin_no_2fa_code="$(json_request GET /api/v1/admin/metrics "" "$admin_token" "" "" "$metrics_admin_no_2fa_file")"

metrics_admin_file="$tmpdir/metrics-admin.txt"
metrics_admin_code="$(json_request GET /api/v1/admin/metrics "" "$admin_token" "" "$admin_2fa" "$metrics_admin_file")"

jq -n \
  --arg generatedAt "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  --arg baseURL "$base_url" \
  --arg tenantID "$tenant_id" \
  --arg guestAPIKey "$guest_api_key" \
  --arg processCode "$process_code" \
  --arg statusPublisherCode "$status_publisher_code" \
  --arg statusGuestCode "$status_guest_code" \
  --arg metricsPublisherCode "$metrics_publisher_code" \
  --arg metricsAdminNo2FACode "$metrics_admin_no_2fa_code" \
  --arg metricsAdminCode "$metrics_admin_code" \
  --slurpfile process "$process_file" \
  --slurpfile statusPublisher "$status_publisher_file" \
  --slurpfile statusGuest "$status_guest_file" \
  --slurpfile metricsPublisher "$metrics_publisher_file" \
  --slurpfile metricsAdminNo2FA "$metrics_admin_no_2fa_file" \
  '{
    generatedAt: $generatedAt,
    baseURL: $baseURL,
    tenantID: $tenantID,
    guestAPIKey: $guestAPIKey,
    checks: {
      lcpProcess: {
        ok: ($processCode == "200"),
        statusCode: ($processCode | tonumber),
        response: ($process[0] // null)
      },
      lcpStatusPublisher: {
        ok: ($statusPublisherCode == "200"),
        statusCode: ($statusPublisherCode | tonumber),
        response: ($statusPublisher[0] // null)
      },
      lcpStatusGuest: {
        ok: ($statusGuestCode == "200"),
        statusCode: ($statusGuestCode | tonumber),
        response: ($statusGuest[0] // null)
      },
      adminMetricsRejectsPublisher: {
        ok: ($metricsPublisherCode == "403"),
        statusCode: ($metricsPublisherCode | tonumber),
        response: ($metricsPublisher[0] // null)
      },
      adminMetricsRequiresTwoFactor: {
        ok: ($metricsAdminNo2FACode == "403"),
        statusCode: ($metricsAdminNo2FACode | tonumber),
        response: ($metricsAdminNo2FA[0] // null)
      },
      adminMetricsWithTwoFactor: {
        ok: ($metricsAdminCode == "200"),
        statusCode: ($metricsAdminCode | tonumber)
      }
    }
  }' > "$out"

if ! jq -e '.checks | all(.[]; .ok == true)' "$out" >/dev/null; then
  jq '.checks | to_entries[] | select(.value.ok != true)' "$out" >&2
  exit 1
fi

echo "$out"
