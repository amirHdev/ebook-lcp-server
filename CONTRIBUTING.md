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

Run the full local stack when you want API, storage, Readium services, and the admin UI together:

```bash
git clone https://github.com/amirHdev/ebook-lcp-server.git
cd ebook-lcp-server
cp .env.example .env
docker compose up --build
sh scripts/demo-local.sh
```

For frontend-only work, install the admin UI dependencies and start Vite:

```bash
cd frontend
npm ci
npm run dev
```

The frontend dev server is available at `http://localhost:5173` and expects the API at
`http://localhost:8080`. Start the API with Docker Compose or run it directly from the repo root:

```bash
cp .env.example .env
export $(grep -v '^#' .env | xargs)
go run ./cmd/server
```

Before opening a pull request, run the checks that match the files you changed:

```bash
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
