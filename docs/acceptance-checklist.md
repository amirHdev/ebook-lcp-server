# Acceptance Checklist

Use this checklist for final client sign-off.

Local automation:

- `sh scripts/acceptance-smoke.sh` validates the API/auth checks that can run against a live local stack.
- `sh scripts/e2e-readiness.sh` brings up the minimal compose stack, runs acceptance plus certification smoke flows, and writes artifacts under `artifacts/e2e-readiness/`.

| Item | Status |
| --- | --- |
| `/api/v1/lcp/process` returns 200 for valid authenticated requests | Pending client test |
| `/api/v1/lcp/status` returns status for user and guest roles | Pending client test |
| `/api/v1/admin/metrics` rejects non-admin users with 403 | Pending client test |
| Admin metrics require `X-2FA-Code` when configured | Pending client test |
| PostgreSQL migrations apply successfully | Pending environment |
| Docker backend image builds | Pending CI |
| Docker frontend image builds | Pending CI |
| Kubernetes Deployment, Service, Ingress, ConfigMap, Secret apply | Pending cluster |
| HPA scales backend pods by CPU/memory | Pending cluster |
| NetworkPolicy is active | Pending cluster |
| Prometheus scrapes `/metrics` | Pending cluster |
| Daily backup CronJob exists | Pending cluster |
| Trivy scan has no critical findings | Pending CI |
| k6 p95 response time is below 200 ms with 100 VUs | Pending load environment |
