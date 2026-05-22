# Contributing to Ebook LCP Server

Thank you for your interest in contributing! This project welcomes developers of all experience levels.

## How you can help

You can contribute by:

- Reporting bugs
- Suggesting features
- Improving documentation
- Writing tests
- Fixing issues
- Improving Docker/Kubernetes deployment
- Improving the admin UI
- Reviewing pull requests

## Developer setup

```bash
cp .env.example .env
docker compose up --build
sh scripts/demo-local.sh
make lint
go test ./...
make coverage
cd frontend && npm ci && npm run build
```

For local demos, the placeholders in `.env.example` are enough to start. The values most contributors
change first are auth (`JWT_SECRET`, admin credentials), storage (`LCP_STORAGE_MODE` and `LCP_S3_*`),
and service URLs (`PUBLIC_BASE_URL`, `STATUS_BASE_URL`, `LCP_CORE_URL`).

Before opening a pull request, include a short summary, verification notes, and docs updates when public APIs or deployment flows change.
