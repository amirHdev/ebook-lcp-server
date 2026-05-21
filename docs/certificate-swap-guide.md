# Production Certificate Swap Guide

Use this guide when moving from the repository's local test certificate setup to a real EDRLab-issued production certificate.

## What changes

The local Compose flow uses the bundled test material that is suitable for development and smoke checks. Production deployments must replace that material with:

1. the EDRLab-issued certificate
2. the matching private key
3. the correct public provider URI used in issued licenses

## Environment variables

Set these values in the production environment:

```bash
LCP_CERTIFICATE=/run/secrets/lcp-cert.pem
LCP_PRIVATE_KEY=/run/secrets/lcp-key.pem
LCP_PROVIDER_URI=https://yourdomain.example
PUBLIC_BASE_URL=https://yourdomain.example
STATUS_BASE_URL=https://status.yourdomain.example
```

Keep the certificate and private key outside the repo and mount them through your runtime secret mechanism.

## Docker and Kubernetes guidance

- Docker: mount the certificate and key as read-only files or secrets.
- Kubernetes: store them in a Secret and mount them into the container filesystem.
- Do not bake production certificate material into images.

## Verification before an EDRLab run

1. Generate a local evidence bundle:

```bash
sh scripts/generate-certification-packet.sh
```

2. Confirm the packet includes:
   - health and readiness responses
   - a sample uploaded publication
   - a generated LCPL file
   - a license status document
   - admin audit and license snapshots
   - certificate metadata when `LCP_CERTIFICATE` is available

3. Check the certificate details:
   - expected subject
   - expected issuer
   - correct expiration date
   - expected provider URI in generated licenses

## Key handling expectations

- Keep the private key in a secret store or HSM-backed workflow where possible.
- Never commit the certificate keypair to the repo.
- Never print private key material into logs.
- Restrict filesystem and runtime access to the signing key.

## What remains external

This repo can prepare the packet and local evidence. It still cannot self-issue official certification. The final EDRLab run must happen with production certificate material and official external validation.
