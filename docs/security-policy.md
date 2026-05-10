# Security Policy

## Authentication

All contract API endpoints require `Authorization: Bearer <jwt>`.

JWTs must be signed with HS256 using `JWT_SECRET` and include one of:

```json
{"sub":"user-id","role":"admin","exp":1893456000}
```

or:

```json
{"sub":"user-id","roles":["user"],"exp":1893456000}
```

## Authorization

Roles:

- `admin`: all endpoints, including `/api/v1/admin/metrics` and publication catalog management.
- `publisher`: publication catalog management, GraphQL publication workflows, and processing.
- `user`: `/api/v1/lcp/process`, `/api/v1/lcp/status`, GraphQL operations, and publication downloads.
- `guest`: `/api/v1/lcp/status` and publication downloads.

Admin requests must also include `X-2FA-Code` matching `ADMIN_2FA_CODE` when that environment variable is configured.

## Public API docs

The running service exposes the API spec at:

- `/swagger.yaml`
- `/swagger.json`
- `/docs/openapi.yaml`
- `/docs/swagger.json`

## Kubernetes

The production overlay includes:

- Secrets for JWT and 2FA values.
- NetworkPolicy restricting ingress to HTTP and egress to DNS/PostgreSQL.
- Liveness and readiness probes.
- HPA for CPU and memory scaling.
- PVCs for metadata and encrypted assets.

On self-hosted K3s, enable encryption at rest for Secrets if you need at-rest secret protection, or replace the Secret manifest with Vault/External Secrets.

## Vulnerability Management

CI runs Go vet, tests, and Trivy filesystem scanning for critical vulnerabilities. Container image scanning should be enabled in the target registry as well.
