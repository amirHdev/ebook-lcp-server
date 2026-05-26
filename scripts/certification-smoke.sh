#!/bin/sh
set -eu

base_url="${BASE_URL:-http://localhost:8080}"
out="${1:-certification-report.json}"
admin_username="${ADMIN_USERNAME:-admin}"
admin_password="${ADMIN_PASSWORD:-admin}"
admin_2fa="${ADMIN_2FA_CODE:-123456}"

token="$(
  curl -fsS "$base_url/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d "$(jq -cn \
      --arg username "$admin_username" \
      --arg password "$admin_password" \
      --arg twoFactor "$admin_2fa" \
      '{username: $username, password: $password, twoFactor: $twoFactor}')" |
    jq -r .token
)"

health="$(curl -fsS "$base_url/healthz")"
ready="$(curl -fsS "$base_url/readyz")"
status="$(curl -fsS "$base_url/api/v1/lcp/status" -H "Authorization: Bearer $token")"
licenses="$(curl -fsS "$base_url/api/v1/admin/licenses" -H "Authorization: Bearer $token" -H "X-2FA-Code: $admin_2fa")"
audit="$(curl -fsS "$base_url/api/v1/admin/audit" -H "Authorization: Bearer $token" -H "X-2FA-Code: $admin_2fa")"

jq -n \
  --arg generatedAt "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  --arg baseURL "$base_url" \
  --argjson health "$health" \
  --argjson ready "$ready" \
  --argjson status "$status" \
  --argjson licenses "$licenses" \
  --argjson audit "$audit" \
  '{
    generatedAt: $generatedAt,
    baseURL: $baseURL,
    checks: {
      healthz: {
        ok: ($health.status == "ok"),
        response: $health
      },
      readyz: {
        ok: ($ready.status == "ready"),
        response: $ready
      },
      lcpStatus: {
        ok: (($status.status == "ok") and (($status.processes // []) | all(.status != "failed"))),
        response: $status
      },
      adminLicenses: {
        ok: true,
        count: (($licenses.licenses // []) | length)
      },
      adminAudit: {
        ok: true,
        count: (($audit.entries // []) | length)
      }
    }
  }' > "$out"

if ! jq -e '.checks | all(.[]; .ok == true)' "$out" >/dev/null; then
  jq '.checks | to_entries[] | select(.value.ok != true)' "$out" >&2
  exit 1
fi

echo "$out"
