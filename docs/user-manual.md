# LCP Server User Manual

## Sign In

`POST /api/v1/auth/login`

Body:

```json
{
  "username": "amiradmin",
  "password": "<admin password>",
  "twoFactor": "<2fa code>"
}
```

Response:

```json
{
  "token": "<jwt>",
  "role": "admin",
  "subject": "amiradmin",
  "expiresAt": "2026-05-16T00:00:00Z"
}
```

## Upload And Process

`POST /api/v1/lcp/process`

Headers:

```http
Authorization: Bearer <jwt>
Content-Type: application/json
```

Body:

```json
{
  "title": "Example Book",
  "file": "base64-encoded-content"
}
```

Response:

```json
{
  "id": "process-id",
  "status": "completed",
  "publicationId": "publication-id",
  "createdAt": "2026-05-09T00:00:00Z",
  "updatedAt": "2026-05-09T00:00:00Z"
}
```

The dashboard now provides a file picker. It reads the selected publication file in the browser, base64-encodes it, and sends it to this endpoint.

## Check Status

`GET /api/v1/lcp/status?id=<process-id>`

Omit `id` to list known process statuses and service uptime.

## Admin Metrics

`GET /api/v1/admin/metrics`

Headers:

```http
Authorization: Bearer <admin-jwt>
X-2FA-Code: <configured-code>
```

The response includes uptime, process count, and request counters.
