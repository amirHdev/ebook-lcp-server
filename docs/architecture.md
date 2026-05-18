# Architecture

## Components

- HTTP server: standard-library `net/http` router in `cmd/server`.
- Auth middleware: HS256 JWT validation and RBAC in `internal/auth`.
- REST adapter: contract endpoints in `internal/adapter/rest`.
- GraphQL adapter: legacy publication and license operations in `internal/adapter/graphql`.
- Use cases: publication upload/encryption and license issuing in `internal/usecase/lcp`.
- Repositories: in-memory by default, JSON-backed metadata when `DATA_DIR` is configured.
- Storage: filesystem-backed encrypted content under `LCP_STORAGE_FS_DIR`, or S3-compatible storage when `LCP_STORAGE_MODE=s3`.
- Webhooks: optional HTTP event delivery for publication and license lifecycle events.
- Audit log: local or JSON-backed action history for publication and license changes.
- Tenancy: publications and licenses carry a tenant ID and reads are filtered by the authenticated tenant.

## Data Flow

```text
Client
  -> JWT/RBAC middleware
  -> REST or GraphQL handler
  -> LCP use case
  -> encrypter/license service
  -> metadata repository
  -> publication storage
```

When S3 storage is enabled, publication download routes return a short-lived signed URL instead of streaming the object through the API process.

## Production Notes

The local JSON repository is suitable for a single-writer deployment and acceptance testing. For production metadata persistence, use PostgreSQL and run the included migrations in `migrations/`.
