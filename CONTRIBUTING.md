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
docker compose up --build
sh scripts/demo-local.sh
make lint
go test ./...
make coverage
cd frontend && npm ci && npm run build
```

Before opening a pull request, include a short summary, verification notes, and docs updates when public APIs or deployment flows change.
