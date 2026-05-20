# Roadmap

The core server work is in place now: local Compose stack, Readium sidecars, PostgreSQL, S3/MinIO storage, signed URLs, webhooks, audit logs, request tracing, metrics, tenant support, and live end-to-end checks.

## Next steps

| Priority | Item | Status | Why it matters |
| --- | --- | --- | --- |
| P0 | Complete one clean frontend dependency install and production build verification | Done | Confirmed with `npm ci` and `npm run build` on May 18, 2026. |
| P1 | Build native integrations for Calibre, calibre-web, and Kavita | Done | Adds a Calibre plugin package, a single-file forwarder, and a polling sidecar for calibre-web/Kavita library folders. |
| P2 | Expand the admin dashboard into a fuller operations console | Planned | Adds the day-to-day controls operators expect: license search, filters, active/expired loan views, audit visibility, and clearer states. |
| P3 | Add a small operator CLI for common tasks | Planned | Gives self-hosters a fast terminal path for upload, license creation, revocation, health checks, and smoke tests. |
| P4 | Prepare the final certification packet and run official EDRLab validation | External | Turns the current blueprint into commercial proof once production certificates and official EDRLab testing are available. |
| P5 | Add real reader demos for Readium Swift and Android | Planned | Proves the encrypted publications work in concrete mobile client flows, not only on the server side. |
| P6 | Publish SDKs after the API settles | Planned | Makes client integration easier once the public surface is stable enough to support long-lived packages. |
| P7 | Invest in public project polish and community growth | Planned | Improves discoverability and trust through stronger docs, screenshots, release cadence, fast issue response, and regular community posts. |

## Notes

- Deployment guides are complete under `docs/deploy.md`, `docs/deploy-flyio.md`, and `docs/deploy-railway.md`.
- The immediate product focus is adoption first, then operator ergonomics, then certification, broader client support, and visible community momentum.
