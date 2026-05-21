# LCP Certification Blueprint

This repository cannot self-certify an installation for EDRLab, but it can prepare a repeatable certification packet before an official test run.

## Repo-owned preparation

The repository can now generate a local evidence bundle with:

- service health and readiness responses
- LCP status API output
- admin audit and license snapshots
- a sample upload and license issuance trace
- a generated `.lcpl`
- a generated license status document
- an encrypted publication artifact
- deployment config snapshots
- certificate metadata when a production certificate path is supplied

Generate the packet:

```bash
sh scripts/generate-certification-packet.sh
```

The command writes a `certification-packet/` directory with a manifest and captured responses.

## Evidence to collect for the official run

1. Build SHA and deployment configuration snapshot
2. Public provider URI and certificate chain used for signing
3. Encrypted EPUB, PDF, and manifest samples
4. License create, download, status, extension, and revocation traces
5. Reader validation notes
6. Official EDRLab test output once run against production certificates

## Local report

```bash
sh scripts/demo-local.sh
sh scripts/certification-smoke.sh
```

The smoke script writes `certification-report.json` with machine-readable readiness checks.

## Production certificate swap

See [docs/certificate-swap-guide.md](/Users/mehrbod/projects/golang/lcp/docs/certificate-swap-guide.md) for the operational handoff from the bundled test setup to a real production certificate and private key.

## External step that still remains

Official EDRLab certification still requires:

1. real EDRLab-issued production certificate material
2. the final reader and publisher validation run
3. the official external test execution and resulting report
