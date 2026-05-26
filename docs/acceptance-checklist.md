# Acceptance Checklist

Use this checklist for final client sign-off.

Local automation:

- `sh scripts/acceptance-smoke.sh` validates the API/auth checks against a live stack and fails if any recorded check is unsuccessful.
- `sh scripts/e2e-readiness.sh` brings up the minimal compose stack, runs acceptance plus certification smoke flows, proves the license status lifecycle, and writes artifacts under `artifacts/e2e-readiness/`.

| Item | Status |
| --- | --- |
| `/api/v1/lcp/process` encrypts a valid authenticated EPUB request | Enforced by E2E Readiness |
| `/api/v1/lcp/status` returns status for user and guest roles | Enforced by E2E Readiness |
| `/api/v1/admin/metrics` rejects non-admin users with 403 | Enforced by E2E Readiness |
| Admin metrics require `X-2FA-Code` when configured | Enforced by E2E Readiness |
| License issuance registers in LSD, supports extension, and maps revocation to LSD `cancelled` status | Enforced by E2E Readiness |
| PostgreSQL migrations apply successfully | Evidenced by E2E Readiness artifacts |
| Docker backend image builds | Enforced by E2E Readiness |
| Docker frontend image builds | Enforced by E2E Readiness |
| Kubernetes Deployment, Service, Ingress, ConfigMap, Secret apply | Pending cluster |
| HPA scales backend pods by CPU/memory | Pending cluster |
| NetworkPolicy is active | Pending cluster |
| Prometheus scrapes `/metrics` | Pending cluster |
| Daily backup CronJob exists | Pending cluster |
| Trivy scan has no critical findings | Enforced by CI |
| k6 p95 response time is below 200 ms with 100 VUs | Pending load environment |

The automated statuses above describe repository gates. Production certificate validation, reader validation, cluster checks, load testing, and official EDRLab certification remain external acceptance work.
