# LCP Server - EBook Manager v1

Lightweight License Content Protection (LCP) server that exposes REST and GraphQL APIs for processing encrypted publications and issuing licenses. The repository includes the DevOps assets needed to run the service on self-hosted K3s with Docker, Kubernetes, GitLab CI/CD, GitHub Actions, and ArgoCD.

## Features

- Contract REST endpoints at `/api/v1/lcp/process`, `/api/v1/lcp/status`, and `/api/v1/admin/metrics`.
- JWT authentication with RBAC roles (`admin`, `publisher`, `user`, `guest`) and optional admin 2FA via `X-2FA-Code`.
- GraphQL endpoint at `/graphql` for managing publications and licenses.
- Pluggable encryption interface backed by the upstream Readium `lcpencrypt` tool for real LCP publication processing.
- In-memory repositories by default and JSON-backed metadata persistence when `DATA_DIR` is configured.
- Download endpoint at `/publications/{id}/content` for clients to retrieve encrypted assets using the URLs returned on licenses.
- Deployment assets for Docker, Kubernetes probes, PVCs, HPA, NetworkPolicy, backup CronJob, and ArgoCD GitOps flows.
- Self-hosted registry manifests in `deploy/registry`.
- React/TypeScript admin UI in `frontend` for processing content, viewing status, and checking admin metrics.
- CI pipelines that format-check, vet, test, scan with Trivy, build, and deploy the container image.

## Configuration

Set the following environment variables (see `.env.example & .env.local` for defaults):

- `DB_DSN`: Database connection string (used by adapters that expect persistent storage).
- `LCP_PROFILE`: `basic` or `production`.
- `LCP_CERTIFICATE` / `LCP_PRIVATE_KEY`: Paths to DRM keys/certificates.
- `LCP_STORAGE_MODE`: `fs` (default) or `s3`.
- `LCP_STORAGE_FS_DIR`: Target directory for encrypted assets.
- `LCP_S3_REGION`, `LCP_S3_BUCKET`, `LCP_S3_ACCESS_KEY`, `LCP_S3_SECRET_KEY`: S3 storage settings when `LCP_STORAGE_MODE=s3`.
- `JWT_SECRET`: Secret for future JWT-protected endpoints.
- `ADMIN_2FA_CODE`: Optional code required for admin role requests.
- `SERVER_PORT`: Listen address (defaults to `:8080`).
- `PUBLIC_BASE_URL`: Public base URL used to generate download links (defaults to `http://localhost:PORT`).
- `LCP_PROVIDER_URI`: Provider URI embedded in generated LCP licenses.
- `LCP_CORE_URL`: Internal Readium LCP core URL.
- `STATUS_BASE_URL`: Public Readium status server URL.
- `DATA_DIR`: Directory used for JSON-backed metadata persistence.

## Local development

```bash
# Install dependencies and start the API
cp .env.example .env  # then edit values
export $(grep -v '^#' .env | xargs)
go run ./cmd/server
```

The GraphQL playground will be available at `http://localhost:8080/graphql`.

### REST API

All contract endpoints require a Bearer JWT signed with `JWT_SECRET`. Admin calls also require `X-2FA-Code` when `ADMIN_2FA_CODE` is set.

```bash
curl -X POST http://localhost:8080/api/v1/lcp/process \
  -H "Authorization: Bearer $JWT" \
  -H "Content-Type: application/json" \
  -d '{"title":"Example","file":"aGVsbG8="}'

curl http://localhost:8080/api/v1/lcp/status \
  -H "Authorization: Bearer $JWT"

curl http://localhost:8080/api/v1/admin/metrics \
  -H "Authorization: Bearer $ADMIN_JWT" \
  -H "X-2FA-Code: $ADMIN_2FA_CODE"
```

### GraphQL upload notes

The `uploadPublication` mutation expects the `file` argument as a **base64-encoded string**. Example variables:

```json
{
  "title": "My Book",
  "file": "<base64 encoded content>"
}
```

## Docker

```bash
docker build -t lcp-server:local .
docker run -p 8080:8080 --env-file .env lcp-server:local
```

The multi-stage Dockerfile compiles the Go binary and ships a minimal distroless runtime image suitable for production.

For the full local stack:

```bash
docker compose up --build
```

This starts PostgreSQL, Redis, the backend, the admin UI, Prometheus, and Grafana.

## Frontend

```bash
cd frontend
npm ci
npm run dev
```

The admin UI is available at `http://localhost:5173` in development.

## Kubernetes

Apply the manifests with Kustomize:

```bash
kubectl apply -k deploy/overlays/prod
```

The deployment uses K3s-friendly defaults: Traefik ingress, local-path storage, and built-in cluster components. Update the overlay image references and hostnames for your environment.
Images are expected to live in the self-hosted registry at `registry.testmedical.ir:5000`.

## ArgoCD

The ArgoCD root application at `deploy/argocd/root-application.yaml` continuously syncs the environment apps from this repo:

```bash
kubectl apply -n argocd -f deploy/argocd/root-application.yaml
```

Namespaces `lcp-dev`, `lcp-staging`, and `lcp-prod` will be created automatically, and changes merged to the repo will propagate to the cluster.

## GitLab CI/CD

`.gitlab-ci.yml` defines four stages:

1. **lint**: runs `gofmt -l` and `go vet`.
2. **test**: executes `go test ./...`.
3. **build**: builds the Docker image and optionally pushes it when registry credentials exist.
4. **deploy**: applies the Kustomize overlays to the cluster (expects `kubectl` credentials in CI variables).

Set `CI_REGISTRY`, `CI_REGISTRY_USER`, `CI_REGISTRY_PASSWORD`, and `KUBECONFIG` (or in-cluster service account variables) in GitLab to enable full automation.

## Repository layout

- `cmd/server`: HTTP server wiring, GraphQL handler, and LCP use cases.
- `internal/auth`: JWT validation, RBAC, and admin 2FA middleware.
- `internal/adapter/rest`: REST endpoints required by the contract.
- `internal/usecase/lcp`: Business logic for publications and licenses.
- `internal/adapter/graphql`: GraphQL schema and resolvers.
- `deploy/k8s`: Production manifests with Kustomize.
- `frontend`: React/TypeScript admin dashboard.
- `deploy/argocd`: GitOps application definition.
- `deploy/registry`: in-cluster image registry.
- `.gitlab-ci.yml`: Pipeline definition for GitLab.

## Documentation

- `docs/deployment-guide.md`
- `docs/security-policy.md`
- `docs/user-manual.md`
- `docs/architecture.md`
- `docs/openapi-rest.yaml`
- Swagger/OpenAPI is also exposed at runtime on `/swagger.yaml` and `/swagger.json`.
- `docs/acceptance-checklist.md`
- `docs/support-and-knowledge-transfer.md`

## Load Test

```bash
k6 run -e BASE_URL=http://localhost:8080 -e JWT="$JWT" tests/k6/lcp-status.js
```
