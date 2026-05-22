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
git clone https://github.com/amirHdev/ebook-lcp-server.git
cd ebook-lcp-server
docker compose up --build
sh scripts/demo-local.sh
make lint
go test ./...
make coverage
```

## Frontend development

The admin UI is a React + Vite + TypeScript app. To run it locally:

```bash
cd frontend
npm ci
npm run dev
```

The frontend talks to the API at `http://localhost:8080` by default. Make sure the API is running (or use the full Compose stack to start everything together). The frontend is served at `http://localhost:5173`.

To build the frontend:

```bash
cd frontend
npm ci
npm run build
```


Before opening a pull request, include a short summary, verification notes, and docs updates when public APIs or deployment flows change.
