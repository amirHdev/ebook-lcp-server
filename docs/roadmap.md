# Roadmap

The core server work is in place now: local Compose stack, Readium sidecars, PostgreSQL, S3/MinIO storage, signed URLs, webhooks, audit logs, request tracing, metrics, tenant support, and live end-to-end checks.

## Active roadmap

| Priority | Item | Status | Why it matters |
| --- | --- | --- | --- |
| P4 | Run official EDRLab validation with production certificates | Ready for external run | The repo-owned packet and evidence workflow are complete; the remaining work is the real external certification run. |
| P6 | Publish SDKs after the API settles | Planned | Makes client integration easier once the public surface is stable enough to support long-lived packages. |
| P7 | Invest in public project polish and community growth | Planned | Improves discoverability and trust through stronger docs, screenshots, release cadence, fast issue response, and regular community posts. |

## Notes

- Completed phases P0 through P5 are now delivered in the repo and omitted here to keep the roadmap focused on open work.
- Deployment guides are complete under `docs/deploy.md`, `docs/deploy-flyio.md`, and `docs/deploy-railway.md`.
- The immediate product focus is adoption first, then operator ergonomics, then certification, broader client support, and visible community momentum.
