# Roadmap

This document tracks the next work needed to move the project from a usable implementation to a production-ready LCP service.

## Current state

| Item | Status |
| --- | --- |
| REST and GraphQL APIs | Done |
| PostgreSQL support | Done |
| Docker Compose stack | Done |
| Readium LCP and LSD services in the local stack | Done |
| MinIO/S3 storage support | Done |
| Signed download URLs | Done |
| Webhook events | Done |
| Audit logs | Done |
| Rate limiting | Done |
| Tenant-aware publications and licenses | Done |
| OpenAPI files | Done |
| Admin UI | Done |

## Immediate proof work

Before calling the current build release-ready, the full local flow still needs to be proven end to end:

| Item | Status |
| --- | --- |
| Boot the full Docker Compose stack | Next |
| Upload the sample EPUB | Next |
| Create a real license | Next |
| Download the protected file through the signed URL flow | Next |
| Open the result in Thorium Reader | Next |
| Verify MinIO objects, webhook delivery, audit records, rate limits, and tenant isolation | Next |
| Fix any bugs found during that run | Next |

## Next implementation batches

### Batch 1 - Make the current features durable

| Item | Status |
| --- | --- |
| Store audit logs in PostgreSQL | Planned |
| Add webhook retries with backoff | Planned |
| Add a failed-delivery record or dead-letter path for webhooks | Planned |
| Add request IDs to API responses and logs | Planned |
| Switch server logging to structured output | Planned |

### Batch 2 - Make operations boring

| Item | Status |
| --- | --- |
| Improve readiness checks beyond PostgreSQL | Planned |
| Add backup and restore verification steps | Planned |
| Add metrics for webhook delivery, S3 operations, license generation, and auth failures | Planned |
| Add load-test baselines and publish expected numbers | Planned |
| Document upgrade and rollback steps | Planned |

### Batch 3 - Finish tenant support

| Item | Status |
| --- | --- |
| Tenant CRUD/admin API | Planned |
| Tenant-aware API keys or external auth integration | Planned |
| Per-tenant webhook config | Planned |
| Per-tenant storage/prefix options | Planned |
| Tenant-aware audit queries | Planned |
| Quota and rate-limit settings per tenant | Planned |

## Adoption work

These are not blockers for correctness, but they matter for repo growth:

| Item | Status |
| --- | --- |
| Real screenshots and one short demo GIF | Planned |
| Hosted OpenAPI docs | Planned |
| One-click deployment guides | Planned |
| Reader demos for Thorium, Readium Swift, and Android | Planned |
| Small CLI for common tasks | Planned |
| SDKs after the API settles | Planned |

## Release line

The next strong public milestone should be:

**v0.2 - self-hosted local stack**

Expected before that tag:

| Item | Status |
| --- | --- |
| Full Docker demo passes | Next |
| README demo works as written | Next |
| MinIO/S3 path verified | Next |
| Signed URLs verified | Next |
| Webhooks verified | Next |
| Audit log endpoint verified | Next |
| Tenant isolation verified | Next |

The next production-minded milestone after that:

**v0.3 - operations hardening**

Expected before that tag:

| Item | Status |
| --- | --- |
| PostgreSQL audit logs | Planned |
| Webhook retry handling | Planned |
| Structured logs and request IDs | Planned |
| Richer metrics | Planned |
| Backup/restore exercise documented | Planned |
