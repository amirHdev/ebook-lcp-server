# Foundational Task: Publication Catalog Service API Design

This document proposes the first incremental step toward the full digital bookstore with LCP DRM: define the Publication Catalog service API. The goal is to provide a stable contract for managing protected eBooks, which future services (storefront, orders, payments) can integrate with.

## Scope
- Manage publications (metadata + encrypted assets) stored via the existing LCP stack.
- Provide endpoints for ingestion, retrieval, search, and lifecycle operations (activate/deactivate).
- Use JWT bearer authentication; authorization rules can be refined later.

## API Endpoints
| Method | Path | Purpose |
| --- | --- | --- |
| `POST` | `/publications` | Ingest a new publication (metadata + encrypted file location). |
| `GET` | `/publications/{id}` | Retrieve a single publication by ID. |
| `GET` | `/publications` | List/search publications (filter by status, title, author, tags). |
| `PATCH` | `/publications/{id}` | Update metadata fields (title, authors, subjects, tags). |
| `POST` | `/publications/{id}/activate` | Mark a publication as available for sale/delivery. |
| `POST` | `/publications/{id}/deactivate` | Mark a publication as unavailable. |

### Sample Payloads
**Create Publication (ingest)**
```json
{
  "title": "The Go Programming Language",
  "authors": ["Alan A. A. Donovan", "Brian W. Kernighan"],
  "language": "en",
  "subjects": ["programming", "golang"],
  "encrypted_uri": "s3://bucket/path/book.lcp",
  "checksum": "<sha256>",
  "tags": ["bestseller"],
  "license_duration_days": 30
}
```

**Publication Response**
```json
{
  "id": "pub_123",
  "title": "The Go Programming Language",
  "authors": ["Alan A. A. Donovan", "Brian W. Kernighan"],
  "language": "en",
  "subjects": ["programming", "golang"],
  "encrypted_uri": "s3://bucket/path/book.lcp",
  "status": "active",
  "tags": ["bestseller"],
  "license_duration_days": 30,
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-02T12:00:00Z"
}
```

## Integration Notes
- **LCP Encryption:** The `encrypted_uri` should reference content processed through `lcpencrypt` with certificates configured in `config.LCP`. Ingestion flow should validate checksum and store resulting encryption metadata alongside publication records.
- **GraphQL Compatibility:** Expose the same operations through the existing GraphQL handler to avoid duplicating business logic. REST paths map to GraphQL mutations/queries that reuse the publication use case.
- **Validation & Errors:** Use consistent error structure (e.g., `{ "error": "message" }`) with HTTP status codes: `400` for validation, `401/403` for auth, `404` for missing records, `409` for checksum conflicts, `500` for server issues.

## Next Steps
1. Add handlers for the above endpoints in the Fiber router, wiring them to the `publication` use case and repository.
2. Extend the publication repository to store additional metadata (subjects, tags, activation status, license duration) if not already present in migrations.
3. Update Swagger/OpenAPI docs in `docs/swagger.yaml` and regenerate `swagger.json` to reflect the new routes.
4. Add request/response validation tests and integration tests against the PostgreSQL-backed repository.
